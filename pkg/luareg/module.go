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

// Package luareg provides reflection-based registration of idiomatic Go
// functions as Lua module functions. It converts Go function signatures into
// lua.LGFunction wrappers automatically and records typed metadata for stub
// generation (see pkg/stubgen, chunk A3).
//
// Basic usage:
//
//	func Loader(L *lua.LState) int {
//	    m := luareg.NewModule("strings", "string manipulation utilities")
//	    m.Fn("has_prefix", strings.HasPrefix, "checks if string has prefix")
//	    m.Fn("trim", strings.Trim, "removes cutset from both ends",
//	        luareg.Args("s", "cutset"),
//	        luareg.ArgDoc("s", "the string to trim"),
//	        luareg.ArgDoc("cutset", "the characters to remove"),
//	        luareg.ReturnDoc(0, "out", "the trimmed string"),
//	    )
//	    return m.PushTo(L)
//	}
package luareg

import (
	"fmt"
	"reflect"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
)

// ArgDocEntry: name and description for a single parameter.
// Consumed by stub generators (A3).
type ArgDocEntry struct {
	Name string
	Doc  string
}

// ReturnDocEntry: name and description for a single return value.
// Index is 0-based, excluding the trailing error return if present.
// Consumed by stub generators (A3).
type ReturnDocEntry struct {
	Index int
	Name  string
	Doc   string
}

// FnMeta: all metadata for a single registered function.
// Returned by Module.Funcs() for use by stub generators (A3).
type FnMeta struct {
	LuaName    string
	GoFn       any
	Doc        string
	ArgNames   []string         // positional names; may be shorter than param count
	ArgDocs    []ArgDocEntry    // per-arg descriptions keyed by name
	ReturnDocs []ReturnDocEntry // per-return descriptions keyed by index
}

// ArgName: returns the positional name for arg i (0-indexed). Falls back to
// "arg<n>" when no explicit name was set.
func (f *FnMeta) ArgName(i int) string {
	if i < len(f.ArgNames) {
		return f.ArgNames[i]
	}
	return fmt.Sprintf("arg%d", i+1)
}

// FnOpt: option for decorating a function registration with names/descriptions.
type FnOpt func(*FnMeta)

// Args: sets positional parameter names. Names are used instead of the default
// "arg1", "arg2", … in generated stubs and error messages.
//
// Example:
//
//	m.Fn("trim", strings.Trim, "...", luareg.Args("s", "cutset"))
func Args(names ...string) FnOpt {
	return func(m *FnMeta) {
		m.ArgNames = names
	}
}

// ArgDoc: attaches a description to a parameter identified by name.
// The name must match one previously set via Args, or the stub generator
// will use the positional name.
func ArgDoc(name, doc string) FnOpt {
	return func(m *FnMeta) {
		m.ArgDocs = append(m.ArgDocs, ArgDocEntry{Name: name, Doc: doc})
	}
}

// ReturnDoc: attaches a name and description to the Nth return value
// (0-indexed, excluding the trailing error return if present).
func ReturnDoc(index int, name, doc string) FnOpt {
	return func(m *FnMeta) {
		m.ReturnDocs = append(m.ReturnDocs, ReturnDocEntry{Index: index, Name: name, Doc: doc})
	}
}

// ConstMeta: metadata for a single module-level constant value.
// Used by stub generators to emit ---@field annotations on the module class block.
type ConstMeta struct {
	Name    string
	Value   any
	LuaType string // pre-resolved type annotation; empty = infer at gen time
	Doc     string
}

// Module: a Lua module being assembled. Holds per-function metadata and the
// built lua.LGFunction wrappers. Create with NewModule.
type Module struct {
	name        string
	doc         string
	funcs       []*FnMeta                 // ordered list of registered functions
	byName      map[string]*FnMeta        // index for duplicate detection
	wrappers    map[string]lua.LGFunction // cached wrappers built by PushTo/Register
	classes     []AnyClass                // registered classes in registration order
	classByName map[string]AnyClass       // index for duplicate class detection
	classByType map[reflect.Type]AnyClass // type → class for auto-wrap lookup
	consts      []ConstMeta               // ordered list of module-level constants
	constByName map[string]struct{}       // index for duplicate const detection
}

// NewModule: creates a new module with the given Lua module name and doc string.
//
// Example:
//
//	m := luareg.NewModule("strings", "string manipulation utilities")
func NewModule(name, doc string) *Module {
	return &Module{
		name:        name,
		doc:         doc,
		byName:      make(map[string]*FnMeta),
		wrappers:    make(map[string]lua.LGFunction),
		classByName: make(map[string]AnyClass),
		classByType: make(map[reflect.Type]AnyClass),
		constByName: make(map[string]struct{}),
	}
}

// Name: returns the module's Lua name.
func (m *Module) Name() string { return m.name }

// Doc: returns the module's documentation string.
func (m *Module) Doc() string { return m.doc }

// Funcs: returns all registered function metadata in registration order.
// Intended for use by stub generators (A3).
func (m *Module) Funcs() []*FnMeta { return m.funcs }

// Fn: registers a Go function as a Lua function in this module.
//
// luaName is the name exposed in Lua. goFn is any Go function whose parameter
// and return types are drawn from the supported set: string, bool,
// int/int8/16/32/64, uint/uint8/16/32/64, float32/64, []string, []int,
// []float64, map[string]T, struct and *struct (via pkg/glua.Translator).
// Variadic functions are supported — trailing Lua args are collected into a
// slice of the variadic element type. The special first parameter type
// *lua.LState is passed the Lua state directly and not counted as a Lua
// argument — use it as an escape hatch when reflection cannot satisfy a need.
//
// doc is a one-line summary of the function. opts decorate parameter names and
// descriptions for stub generation.
//
// Panics if luaName was already registered on this module, or if goFn is nil
// or not a function.
//
// Returns the receiver for chaining.
func (m *Module) Fn(luaName string, goFn any, doc string, opts ...FnOpt) *Module {
	if _, exists := m.byName[luaName]; exists {
		panic(fmt.Sprintf("luareg: module %q: duplicate registration of function %q", m.name, luaName))
	}

	validateGoFn(luaName, goFn)

	meta := &FnMeta{
		LuaName: luaName,
		GoFn:    goFn,
		Doc:     doc,
	}
	for _, opt := range opts {
		opt(meta)
	}

	m.funcs = append(m.funcs, meta)
	m.byName[luaName] = meta
	// The classLookup closure captures m.classByType by reference so classes
	// registered after this Fn call are also visible during execution.
	m.wrappers[luaName] = buildWrapper(meta, func(t reflect.Type) AnyClass {
		return m.classByType[t]
	})

	return m
}

// RegisterClass: adds a class to this module. The class's metatable is
// constructed when the module is PushTo'd. The class type is also indexed for
// auto-wrap so functions returning instances of T have their returns wrapped
// automatically.
//
// Panics if a class with the same Lua name was already registered on this module.
//
// Returns the receiver for chaining.
func (m *Module) RegisterClass(c AnyClass) *Module {
	if _, exists := m.classByName[c.Name()]; exists {
		panic(fmt.Sprintf("luareg: module %q: duplicate class registration %q", m.name, c.Name()))
	}
	m.classes = append(m.classes, c)
	m.classByName[c.Name()] = c
	m.classByType[c.goType()] = c
	return m
}

// Classes: returns all registered classes in registration order.
// Intended for use by stub generators (A3).
func (m *Module) Classes() []AnyClass { return m.classes }

// Consts: returns all registered constant metadata in registration order.
// Intended for use by stub generators.
func (m *Module) Consts() []ConstMeta { return m.consts }

// Const: records a module-level constant value. At PushTo time the value is
// translated to a Lua value and set on the module table. At stub-gen time the
// constant becomes a ---@field <name> <luaType> annotation on the module's
// class block.
//
// luaType is the LSP type annotation (e.g. "kubernetes.GVKMatcher" or "table"
// or "string"). Pass an empty string to have the generator infer it from the
// Go value's reflect type.
//
// doc is a one-line summary for the field annotation.
//
// Panics if name was already registered as a const on this module.
func (m *Module) Const(name string, value any, luaType, doc string) *Module {
	if _, exists := m.constByName[name]; exists {
		panic(fmt.Sprintf("luareg: module %q: duplicate registration of const %q", m.name, name))
	}
	m.constByName[name] = struct{}{}
	m.consts = append(m.consts, ConstMeta{
		Name:    name,
		Value:   value,
		LuaType: luaType,
		Doc:     doc,
	})
	return m
}

// PushTo: builds the Lua module table from all registered functions and pushes
// it onto L's stack. Also registers metatables for all classes and sets any
// constant values registered via Const. Returns 1, the Lua loader convention.
//
// Calling PushTo multiple times is safe and produces equivalent tables each
// time (wrappers are cached).
func (m *Module) PushTo(L *lua.LState) int {
	// Build classLookup closure for class metatables.
	classLookup := func(t reflect.Type) AnyClass {
		return m.classByType[t]
	}

	// Register class metatables first so method wrappers can wrap return values.
	for _, c := range m.classes {
		c.register(L, classLookup)
	}

	exports := make(map[string]lua.LGFunction, len(m.wrappers))
	for name, fn := range m.wrappers {
		exports[name] = fn
	}
	mod := L.SetFuncs(L.NewTable(), exports)

	// Set constant values on the module table via the Translator.
	if len(m.consts) > 0 {
		tr := glua.NewTranslator()
		for _, c := range m.consts {
			lv, err := tr.ToLua(L, c.Value)
			if err != nil {
				// Panic at PushTo time so the developer sees the error immediately.
				panic(fmt.Sprintf("luareg: module %q: const %q: ToLua failed: %v", m.name, c.Name, err))
			}
			L.SetField(mod, c.Name, lv)
		}
	}

	L.Push(mod)
	return 1
}

// Register: records this module's metadata in reg for later stub generation.
// The Lua-runtime side is handled by PushTo; Register is for the A3 stub path.
func (m *Module) Register(reg *Registry) {
	reg.add(m)
}
