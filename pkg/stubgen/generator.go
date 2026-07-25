// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package stubgen generates Lua LSP annotation stubs from a luareg.Registry.
// The output format matches the library/*.gen.lua files consumed by lua-language-server.
//
// Basic usage:
//
//	gen := stubgen.NewGenerator()
//	files, err := gen.GenerateFromRegistry(reg, "./library")
package stubgen

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

var (
	errorType    = reflect.TypeOf((*error)(nil)).Elem()
	luaStateType = reflect.TypeOf((*lua.LState)(nil))
)

// Generator: builds Lua LSP stubs from a luareg.Registry and an optional set of
// glua type registrations (for K8s-style struct shape stubs that aren't in
// function signatures).
type Generator struct {
	typeRegistry *glua.TypeRegistry
}

// NewGenerator: creates a new Generator instance.
func NewGenerator() *Generator {
	return &Generator{
		typeRegistry: glua.NewTypeRegistry(),
	}
}

// RegisterType: queues a Go type for stub emission alongside the registry-driven
// output. Used for explicit K8s-style "I want stubs for corev1.Pod even though
// no function takes it as a parameter."
func (g *Generator) RegisterType(obj any) error {
	return g.typeRegistry.Register(obj)
}

// GenerateFromRegistry: emits one .gen.lua file per module in reg into outDir,
// plus an annotations.gen.lua for any types registered via RegisterType.
// Returns the list of files written.
func (g *Generator) GenerateFromRegistry(reg *luareg.Registry, outDir string) ([]string, error) {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, fmt.Errorf("stubgen: create output dir: %w", err)
	}

	var written []string

	for _, m := range reg.Modules() {
		content, err := g.GenerateModule(m)
		if err != nil {
			return written, fmt.Errorf("stubgen: module %q: %w", m.Name(), err)
		}

		path := filepath.Join(outDir, m.Name()+".gen.lua")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return written, fmt.Errorf("stubgen: write %s: %w", path, err)
		}
		written = append(written, path)
	}

	// Emit annotations.gen.lua only if types were registered.
	annotations, err := g.GenerateTypeStubs()
	if err != nil {
		return written, fmt.Errorf("stubgen: annotations: %w", err)
	}
	if annotations != "" {
		path := filepath.Join(outDir, "annotations.gen.lua")
		if err := os.WriteFile(path, []byte(annotations), 0644); err != nil {
			return written, fmt.Errorf("stubgen: write %s: %w", path, err)
		}
		written = append(written, path)
	}

	return written, nil
}

// GenerateModule: emits a single module's stub as a string. Struct types found
// in function and method signatures are auto-discovered and emitted as
// ---@class blocks at the top of the file.
func (g *Generator) GenerateModule(m *luareg.Module) (string, error) {
	modName := m.Name()

	// Build class lookup: go type → AnyClass, for detecting class return types.
	classLookup := buildClassLookup(m)

	// Discover struct types from all function and method signatures.
	structReg := glua.NewTypeRegistry()
	if err := discoverStructTypes(m, classLookup, structReg); err != nil {
		return "", fmt.Errorf("discover struct types: %w", err)
	}
	// Also register explicit stub types (Module.RegisterStubType). These are
	// types the author wants stubbed even though no function signature
	// references them — typical for K8s modules where the API surface is
	// `table<string, any>` but autocomplete on the wrapped Go shapes is wanted.
	for _, v := range m.StubTypes() {
		if v == nil {
			continue
		}
		if err := structReg.Register(v); err != nil {
			return "", fmt.Errorf("register stub type %T: %w", v, err)
		}
	}
	if err := structReg.Process(); err != nil {
		return "", fmt.Errorf("process struct types: %w", err)
	}
	structStubs, err := structReg.GenerateStubs()
	if err != nil {
		return "", fmt.Errorf("generate struct stubs: %w", err)
	}
	// Strip the trailing "return {}\n" that TypeRegistry appends.
	structStubs = stripTypeRegistryTrailer(structStubs)

	var sb strings.Builder

	// Header.
	fmt.Fprintf(&sb, "---@meta %s\n", modName)
	sb.WriteString("\n")

	// ---@alias declarations (Module.RegisterStubAlias). Emit before any
	// ---@class blocks so subsequent field type references resolve.
	if aliases := m.StubAliases(); len(aliases) > 0 {
		for _, a := range aliases {
			if a.Doc != "" {
				fmt.Fprintf(&sb, "---@alias %s %s %s\n", a.Name, a.Def, a.Doc)
			} else {
				fmt.Fprintf(&sb, "---@alias %s %s\n", a.Name, a.Def)
			}
		}
		sb.WriteString("\n")
	}

	// Struct ---@class blocks (if any).
	if structStubs != "" {
		sb.WriteString(structStubs)
		sb.WriteString("\n")
	}

	// Class blocks (in registration order).
	for _, cls := range m.Classes() {
		writeClassBlock(&sb, cls, modName, classLookup)
	}

	// Module class declaration + class fields + const fields.
	fmt.Fprintf(&sb, "---@class %s\n", modName)
	for _, cls := range m.Classes() {
		localName := classLocalName(cls.Name(), modName)
		fmt.Fprintf(&sb, "---@field %s %s\n", localName, cls.Name())
	}
	// Emit ---@field annotations for module-level constants.
	for _, c := range m.Consts() {
		luaType := c.LuaType
		if luaType == "" {
			luaType = goTypeToLua(reflect.TypeOf(c.Value), classLookup)
		}
		if c.Doc != "" {
			fmt.Fprintf(&sb, "---@field %s %s %s\n", c.Name, luaType, c.Doc)
		} else {
			fmt.Fprintf(&sb, "---@field %s %s\n", c.Name, luaType)
		}
	}
	fmt.Fprintf(&sb, "local %s = {}\n", modName)
	sb.WriteString("\n")

	// Module-level function stubs (in registration order).
	for _, fn := range m.Funcs() {
		writeFuncStub(&sb, fn, modName, classLookup)
	}

	// Assign namespaced class variables: modname.ClassName = ClassName
	for _, cls := range m.Classes() {
		localName := classLocalName(cls.Name(), modName)
		fmt.Fprintf(&sb, "%s.%s = %s\n", modName, localName, localName)
		sb.WriteString("\n")
	}

	fmt.Fprintf(&sb, "return %s\n", modName)

	return sb.String(), nil
}

// GenerateTypeStubs: emits the ---@class blocks for any types registered via
// RegisterType. Returns a non-empty string (suitable for annotations.gen.lua)
// only if types were registered; returns "" if the registry is empty.
func (g *Generator) GenerateTypeStubs() (string, error) {
	if err := g.typeRegistry.Process(); err != nil {
		return "", fmt.Errorf("process types: %w", err)
	}
	stubs, err := g.typeRegistry.GenerateStubs()
	if err != nil {
		return "", err
	}
	// If the registry had nothing, GenerateStubs returns just "return {}\n".
	// Treat that as empty — caller should not write the file.
	stripped := stripTypeRegistryTrailer(stubs)
	if stripped == "" {
		return "", nil
	}
	// Return as a valid Lua file with ---@meta header.
	return "---@meta annotations\n\n" + stripped + "return {}\n", nil
}

// buildClassLookup: builds a map from reflect.Type to AnyClass for all classes
// registered in the module. This is used during type inference so that return
// values matching a registered class are annotated with the class's Lua name
// rather than being treated as anonymous structs.
func buildClassLookup(m *luareg.Module) map[reflect.Type]luareg.AnyClass {
	lookup := make(map[reflect.Type]luareg.AnyClass, len(m.Classes()))
	for _, cls := range m.Classes() {
		lookup[cls.GoType()] = cls
	}
	return lookup
}

// discoverStructTypes: walks every function and method signature in the module
// and registers any struct (or *struct) param/return types that are not
// registered classes into structReg for ---@class generation.
func discoverStructTypes(m *luareg.Module, classLookup map[reflect.Type]luareg.AnyClass, structReg *glua.TypeRegistry) error {
	register := func(t reflect.Type) error {
		// Unwrap pointer.
		base := t
		for base.Kind() == reflect.Pointer {
			base = base.Elem()
		}
		if base.Kind() != reflect.Struct {
			return nil
		}
		// Skip if it's a registered class.
		if _, isClass := classLookup[t]; isClass {
			return nil
		}
		if _, isClass := classLookup[base]; isClass {
			return nil
		}
		// Create a zero-value instance to register.
		return structReg.Register(reflect.New(base).Elem().Interface())
	}

	walkFn := func(fn *luareg.FnMeta, skipFirstParam bool) error {
		if fn.GoFn == nil {
			return nil
		}
		ft := reflect.TypeOf(fn.GoFn)
		startParam := 0
		if skipFirstParam && ft.NumIn() > 0 {
			startParam = 1 // skip receiver
		}
		for i := startParam; i < ft.NumIn(); i++ {
			pt := ft.In(i)
			// Skip *lua.LState escape hatch.
			if isLStateType(pt) {
				continue
			}
			if err := register(pt); err != nil {
				return err
			}
		}
		for i := 0; i < ft.NumOut(); i++ {
			rt := ft.Out(i)
			if rt.Implements(errorType) {
				continue
			}
			if err := register(rt); err != nil {
				return err
			}
		}
		return nil
	}

	for _, fn := range m.Funcs() {
		if err := walkFn(fn, false); err != nil {
			return err
		}
	}
	for _, cls := range m.Classes() {
		for _, method := range cls.MethodsMeta() {
			if err := walkFn(method, true); err != nil {
				return err
			}
		}
	}
	// Discover struct types from const values (e.g. GVKMatcher).
	for _, c := range m.Consts() {
		if c.Value == nil {
			continue
		}
		if err := register(reflect.TypeOf(c.Value)); err != nil {
			return err
		}
	}
	return nil
}

// writeClassBlock: writes the ---@class block for cls, then method stubs.
func writeClassBlock(sb *strings.Builder, cls luareg.AnyClass, modName string, classLookup map[reflect.Type]luareg.AnyClass) {
	localName := classLocalName(cls.Name(), modName)

	fmt.Fprintf(sb, "---@class %s\n", cls.Name())
	fmt.Fprintf(sb, "local %s = {}\n", localName)
	sb.WriteString("\n")

	for _, method := range cls.MethodsMeta() {
		writeMethodStub(sb, method, localName, classLookup)
	}
}

// writeMethodStub: writes ---@param / ---@return / function stub for a class method.
func writeMethodStub(sb *strings.Builder, fn *luareg.FnMeta, localName string, classLookup map[reflect.Type]luareg.AnyClass) {
	if fn.GoFn == nil {
		fmt.Fprintf(sb, "function %s:%s() end\n\n", localName, fn.LuaName)
		return
	}

	ft := reflect.TypeOf(fn.GoFn)

	// Build param list, skipping receiver (index 0) and *lua.LState.
	params := buildParamList(ft, fn, true, classLookup)
	rets := buildReturnList(ft, fn, classLookup)

	if fn.Doc != "" {
		fmt.Fprintf(sb, "--- %s\n", fn.Doc)
	}
	for _, p := range params {
		if p.doc != "" {
			fmt.Fprintf(sb, "---@param %s %s %s\n", p.name, p.luaType, p.doc)
		} else {
			fmt.Fprintf(sb, "---@param %s %s\n", p.name, p.luaType)
		}
	}
	for _, r := range rets {
		if r.name != "" && r.doc != "" {
			fmt.Fprintf(sb, "---@return %s %s %s\n", r.luaType, r.name, r.doc)
		} else if r.name != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", r.luaType, r.name)
		} else if r.doc != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", r.luaType, r.doc)
		} else {
			fmt.Fprintf(sb, "---@return %s\n", r.luaType)
		}
	}

	paramNames := make([]string, len(params))
	for i, p := range params {
		paramNames[i] = p.name
	}
	fmt.Fprintf(sb, "function %s:%s(%s) end\n\n", localName, fn.LuaName, strings.Join(paramNames, ", "))
}

// writeFuncStub: writes ---@param / ---@return / function stub for a module-level function.
func writeFuncStub(sb *strings.Builder, fn *luareg.FnMeta, modName string, classLookup map[reflect.Type]luareg.AnyClass) {
	if fn.GoFn == nil {
		fmt.Fprintf(sb, "function %s.%s() end\n\n", modName, fn.LuaName)
		return
	}

	ft := reflect.TypeOf(fn.GoFn)

	params := buildParamList(ft, fn, false, classLookup)
	rets := buildReturnList(ft, fn, classLookup)

	if fn.Doc != "" {
		fmt.Fprintf(sb, "--- %s\n", fn.Doc)
	}
	for _, p := range params {
		if p.doc != "" {
			fmt.Fprintf(sb, "---@param %s %s %s\n", p.name, p.luaType, p.doc)
		} else {
			fmt.Fprintf(sb, "---@param %s %s\n", p.name, p.luaType)
		}
	}
	for _, r := range rets {
		if r.name != "" && r.doc != "" {
			fmt.Fprintf(sb, "---@return %s %s %s\n", r.luaType, r.name, r.doc)
		} else if r.name != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", r.luaType, r.name)
		} else if r.doc != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", r.luaType, r.doc)
		} else {
			fmt.Fprintf(sb, "---@return %s\n", r.luaType)
		}
	}

	paramNames := make([]string, len(params))
	for i, p := range params {
		paramNames[i] = p.name
	}
	fmt.Fprintf(sb, "function %s.%s(%s) end\n\n", modName, fn.LuaName, strings.Join(paramNames, ", "))
}

// paramInfo: collected info for one function parameter.
type paramInfo struct {
	name    string
	luaType string
	doc     string
}

// returnInfo: collected info for one function return value.
type returnInfo struct {
	luaType string
	name    string
	doc     string
}

// buildParamList: reflects on ft to produce a paramInfo slice.
// If skipReceiver is true, ft.In(0) (the method receiver) is skipped.
// *lua.LState parameters are always skipped. A trailing variadic param is
// emitted with name "..." and the variadic element type (i.e. the slice's
// Elem()), so consumers can render `---@param ... <type>` and `function(..., ...)`.
func buildParamList(ft reflect.Type, fn *luareg.FnMeta, skipReceiver bool, classLookup map[reflect.Type]luareg.AnyClass) []paramInfo {
	var params []paramInfo
	luaIdx := 0 // index into ArgNames / ArgDocs (counts only visible Lua params)

	startParam := 0
	if skipReceiver && ft.NumIn() > 0 {
		startParam = 1
	}

	lastParam := ft.NumIn() - 1
	isVariadic := ft.IsVariadic()

	for i := startParam; i < ft.NumIn(); i++ {
		pt := ft.In(i)
		// Skip *lua.LState escape hatch.
		if isLStateType(pt) {
			continue
		}

		variadicTail := isVariadic && i == lastParam

		var name, luaType string
		if variadicTail {
			// LuaLS convention: variadic param name is literally "..." and the
			// type is the element type (string for `...string`), not the slice.
			name = "..."
			luaType = goTypeToLua(pt.Elem(), classLookup)
		} else {
			name = fn.ArgName(luaIdx)
			luaType = goTypeToLua(pt, classLookup)
		}
		// Doc lookup uses the user-visible name. For variadic, the author can
		// attach a doc by passing "..." to ArgDoc.
		doc := findArgDoc(fn, name)

		params = append(params, paramInfo{
			name:    name,
			luaType: luaType,
			doc:     doc,
		})
		luaIdx++
	}
	return params
}

// buildReturnList: reflects on ft to produce a returnInfo slice.
// Trailing error returns are omitted.
func buildReturnList(ft reflect.Type, fn *luareg.FnMeta, classLookup map[reflect.Type]luareg.AnyClass) []returnInfo {
	numOut := ft.NumOut()
	if numOut == 0 {
		return nil
	}
	// Determine if last return is error.
	hasError := ft.Out(numOut - 1).Implements(errorType)
	nonErrCount := numOut
	if hasError {
		nonErrCount--
	}

	// Index ReturnDocs by their .Index field rather than registration order.
	// Authors may register them out-of-order (e.g. ReturnDoc(1,...) before
	// ReturnDoc(0,...)) and the doc must still land on the right return.
	docByIndex := make(map[int]luareg.ReturnDocEntry, len(fn.ReturnDocs))
	for _, d := range fn.ReturnDocs {
		docByIndex[d.Index] = d
	}

	var rets []returnInfo
	for i := 0; i < nonErrCount; i++ {
		rt := ft.Out(i)
		luaType := goTypeToLua(rt, classLookup)

		var name, doc string
		if d, ok := docByIndex[i]; ok {
			name = d.Name
			doc = d.Doc
		}

		rets = append(rets, returnInfo{
			luaType: luaType,
			name:    name,
			doc:     doc,
		})
	}
	return rets
}

// goTypeToLua: maps a Go reflect.Type to its Lua LSP type annotation string.
// Struct or *struct types that are registered classes use the class's Lua name.
// Other struct types use the TypeRegistry naming convention (pkgname.TypeName).
func goTypeToLua(t reflect.Type, classLookup map[reflect.Type]luareg.AnyClass) string {
	// Check if it's a registered class first (before dereferencing pointer).
	if cls, ok := classLookup[t]; ok {
		return cls.Name()
	}

	// Dereference pointer.
	if t.Kind() == reflect.Pointer {
		inner := t.Elem()
		if cls, ok := classLookup[inner]; ok {
			return cls.Name()
		}
		if inner.Kind() == reflect.Struct {
			return structLuaName(inner)
		}
		// Recurse for other pointer types.
		return goTypeToLua(inner, classLookup)
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice:
		elemLua := goTypeToLua(t.Elem(), classLookup)
		return elemLua + "[]"
	case reflect.Map:
		keyLua := goTypeToLua(t.Key(), classLookup)
		valLua := goTypeToLua(t.Elem(), classLookup)
		return fmt.Sprintf("table<%s, %s>", keyLua, valLua)
	case reflect.Struct:
		return structLuaName(t)
	case reflect.Interface:
		return "any"
	default:
		return "any"
	}
}

// structLuaName: derives the Lua class name for a struct type using the last
// path segment of its package as the prefix (e.g. "pkg/modules/widget.Config"
// → "widget.Config"). Mirrors the convention in glua.TypeRegistry.getTypeName.
func structLuaName(t reflect.Type) string {
	pkgPath := t.PkgPath()
	if pkgPath == "" {
		return t.Name()
	}
	parts := strings.Split(pkgPath, "/")
	pkgName := parts[len(parts)-1]

	// Kubernetes special case: k8s.io/api/core/v1 → corev1
	if strings.Contains(pkgPath, "k8s.io/api/") {
		apiParts := strings.Split(pkgPath, "/")
		for i, part := range apiParts {
			if part == "api" && i+2 < len(apiParts) {
				group := apiParts[i+1]
				version := apiParts[i+2]
				return group + version + "." + t.Name()
			}
		}
	}

	return pkgName + "." + t.Name()
}

// classLocalName: strips the "modname." prefix from a fully-qualified class name
// to get the local variable name. E.g. "log.Logger" with modName "log" → "Logger".
// If the class name doesn't start with modname+".", returns the full class name.
func classLocalName(className, modName string) string {
	prefix := modName + "."
	if strings.HasPrefix(className, prefix) {
		return strings.TrimPrefix(className, prefix)
	}
	return className
}

// isLStateType: returns true if t is *lua.LState.
func isLStateType(t reflect.Type) bool {
	// luaStateType is reflect.TypeOf((*lua.LState)(nil)) which is the *lua.LState pointer type.
	return t == luaStateType
}

// findArgDoc: looks up the doc string for a named parameter in fn.ArgDocs.
func findArgDoc(fn *luareg.FnMeta, name string) string {
	for _, ad := range fn.ArgDocs {
		if ad.Name == name {
			return ad.Doc
		}
	}
	return ""
}

// stripTypeRegistryTrailer: removes the trailing "return {}\n" that
// TypeRegistry.GenerateStubs appends, and trims surrounding whitespace.
// Returns "" if nothing remained.
func stripTypeRegistryTrailer(s string) string {
	s = strings.TrimSuffix(strings.TrimSpace(s), "return {}")
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return s + "\n"
}
