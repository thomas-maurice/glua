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

package stubgen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// LuaModule: represents a discovered Lua module
type LuaModule struct {
	Name              string
	Functions         []*LuaFunction
	Classes           []*LuaClass
	Constants         []*LuaConst
	CustomAnnotations []string // Module-level custom annotations
}

// LuaClass: represents a Lua class (UserData type with methods)
type LuaClass struct {
	Name              string
	Description       string
	Methods           []*LuaMethod
	Fields            []*LuaField
	CustomAnnotations []string
}

// LuaField: represents a field in a Lua class
type LuaField struct {
	Name        string
	Type        string
	Description string
}

// LuaMethod: represents a method on a Lua class
type LuaMethod struct {
	Name              string
	Description       string
	Params            []*LuaParam
	Returns           []*LuaReturn
	CustomAnnotations []string
}

// LuaFunction: represents a Lua function exported by a module
type LuaFunction struct {
	Name              string
	Description       string
	Params            []*LuaParam
	Returns           []*LuaReturn
	CustomAnnotations []string // Function-level custom annotations
}

// LuaParam: represents a function parameter
type LuaParam struct {
	Name        string
	Type        string
	Description string
}

// LuaReturn: represents a function return value
type LuaReturn struct {
	Type        string
	Description string
}

// LuaConst: represents a constant exported by a module
type LuaConst struct {
	Name        string
	Type        string
	Description string
}

// Analyzer: scans Go source files and extracts Lua module definitions
type Analyzer struct {
	modules map[string]*LuaModule
}

// NewAnalyzer: creates a new code analyzer instance
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		modules: make(map[string]*LuaModule),
	}
}

// ScanDirectory: recursively scans a directory for Go files and extracts Lua module metadata.
// It parses comment annotations like @luamodule, @luafunc, @luaparam, @luareturn.
func (a *Analyzer) ScanDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		return a.parseFile(path)
	})
}

// parseFile: parses a single Go file and extracts Lua module information
func (a *Analyzer) parseFile(filename string) error {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filename, err)
	}

	var currentModule *LuaModule
	classMap := make(map[string]*LuaClass) // Track classes by name

	// Iterate through all declarations
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Skip if no doc comment
		if funcDecl.Doc == nil {
			continue
		}

		comment := funcDecl.Doc.Text()

		// Check if this is a module loader
		if moduleName := a.extractModuleName(comment); moduleName != "" {
			currentModule = &LuaModule{
				Name:              moduleName,
				Functions:         make([]*LuaFunction, 0),
				Classes:           make([]*LuaClass, 0),
				CustomAnnotations: a.extractCustomAnnotations(comment),
			}
			a.modules[moduleName] = currentModule
			continue
		}

		// Check if this is a Lua function
		if luaFunc := a.extractFunction(comment); luaFunc != nil && currentModule != nil {
			currentModule.Functions = append(currentModule.Functions, luaFunc)
			continue
		}

		// Check if this is a Lua constant
		if luaConst := a.extractConst(comment); luaConst != nil && currentModule != nil {
			currentModule.Constants = append(currentModule.Constants, luaConst)
			continue
		}

		// Check if this is a Lua method
		if luaMethod := a.extractMethod(comment); luaMethod != nil && currentModule != nil {
			className := luaMethod.className
			class, ok := classMap[className]
			if !ok {
				// Create new class
				class = &LuaClass{
					Name:              className,
					Methods:           make([]*LuaMethod, 0),
					Fields:            make([]*LuaField, 0),
					CustomAnnotations: make([]string, 0),
				}
				classMap[className] = class
				currentModule.Classes = append(currentModule.Classes, class)
			}
			class.Methods = append(class.Methods, luaMethod.method)
		}
	}

	// After processing function declarations, scan all comments for constants and standalone classes
	// This allows finding @luaconst and @luaclass annotations anywhere in the file
	if currentModule != nil {
		for _, commentGroup := range file.Comments {
			comment := commentGroup.Text()

			// Check for constants
			if luaConst := a.extractConst(comment); luaConst != nil {
				currentModule.Constants = append(currentModule.Constants, luaConst)
			}

			// Check for standalone class definitions
			if classInfo := a.extractClassDefinition(comment); classInfo != nil {
				// Check if class already exists (e.g., from methods)
				existingClass, exists := classMap[classInfo.Name]
				if exists {
					// Merge fields and description into existing class
					existingClass.Fields = append(existingClass.Fields, classInfo.Fields...)
					if existingClass.Description == "" {
						existingClass.Description = classInfo.Description
					}
				} else {
					// Add new class
					classMap[classInfo.Name] = classInfo
					currentModule.Classes = append(currentModule.Classes, classInfo)
				}
			}
		}
	}

	return nil
}

// extractModuleName: extracts module name from @luamodule annotation
func (a *Analyzer) extractModuleName(comment string) string {
	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@luamodule ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "@luamodule "))
		}
	}
	return ""
}

// extractCustomAnnotations: extracts custom Lua annotations from @luaannotation lines
func (a *Analyzer) extractCustomAnnotations(comment string) []string {
	lines := strings.Split(comment, "\n")
	var annotations []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "@luaannotation ") {
			// Extract everything after @luaannotation
			annotation := strings.TrimSpace(strings.TrimPrefix(line, "@luaannotation "))
			if annotation != "" {
				annotations = append(annotations, annotation)
			}
		}
	}

	return annotations
}

// extractFunction: extracts function information from comment annotations
func (a *Analyzer) extractFunction(comment string) *LuaFunction {
	lines := strings.Split(comment, "\n")

	var (
		fn          *LuaFunction
		description strings.Builder
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "@luafunc ") {
			fn = &LuaFunction{
				Name:              strings.TrimSpace(strings.TrimPrefix(line, "@luafunc ")),
				Params:            make([]*LuaParam, 0),
				Returns:           make([]*LuaReturn, 0),
				CustomAnnotations: make([]string, 0),
			}
			continue
		}

		if fn == nil {
			// Collect description before any annotations
			if !strings.HasPrefix(line, "@") && line != "" {
				if description.Len() > 0 {
					description.WriteString(" ")
				}
				description.WriteString(line)
			}
			continue
		}

		if strings.HasPrefix(line, "@luaparam ") {
			param := a.parseParam(line)
			if param != nil {
				fn.Params = append(fn.Params, param)
			}
		} else if strings.HasPrefix(line, "@luareturn ") {
			ret := a.parseReturn(line)
			if ret != nil {
				fn.Returns = append(fn.Returns, ret)
			}
		} else if strings.HasPrefix(line, "@luaannotation ") {
			annotation := strings.TrimSpace(strings.TrimPrefix(line, "@luaannotation "))
			if annotation != "" {
				fn.CustomAnnotations = append(fn.CustomAnnotations, annotation)
			}
		}
	}

	if fn != nil {
		fn.Description = description.String()
	}

	return fn
}

// extractConst: extracts constant information from comment annotations
func (a *Analyzer) extractConst(comment string) *LuaConst {
	lines := strings.Split(comment, "\n")

	var cnst *LuaConst

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "@luaconst ") {
			// Parse: @luaconst NAME TYPE DESCRIPTION
			parts := strings.Fields(strings.TrimPrefix(line, "@luaconst "))
			if len(parts) >= 2 {
				cnst = &LuaConst{
					Name: parts[0],
					Type: parts[1],
				}
				if len(parts) > 2 {
					cnst.Description = strings.Join(parts[2:], " ")
				}
			}
			break
		}
	}

	return cnst
}

// extractClassDefinition: extracts standalone class definition with fields
// Format:
//
//	@luaclass <ClassName>
//	@luafield <fieldName> <type> <description>
//	@luafield <fieldName> <type> <description>
func (a *Analyzer) extractClassDefinition(comment string) *LuaClass {
	lines := strings.Split(comment, "\n")

	var class *LuaClass
	var className string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for @luaclass annotation
		if strings.HasPrefix(line, "@luaclass ") {
			className = strings.TrimSpace(strings.TrimPrefix(line, "@luaclass "))
			if className != "" {
				class = &LuaClass{
					Name:              className,
					Methods:           make([]*LuaMethod, 0),
					Fields:            make([]*LuaField, 0),
					CustomAnnotations: make([]string, 0),
				}
			}
			continue
		}

		// Look for @luafield annotations (only if we found a class)
		if class != nil && strings.HasPrefix(line, "@luafield ") {
			// Parse: @luafield FIELDNAME TYPE DESCRIPTION
			parts := strings.Fields(strings.TrimPrefix(line, "@luafield "))
			if len(parts) >= 2 {
				field := &LuaField{
					Name: parts[0],
					Type: parts[1],
				}
				if len(parts) > 2 {
					field.Description = strings.Join(parts[2:], " ")
				}
				class.Fields = append(class.Fields, field)
			}
		}
	}

	return class
}

// methodResult: temporary struct to hold method extraction results
type methodResult struct {
	className string
	method    *LuaMethod
}

// extractMethod: extracts method information from comment annotations
// Format: @luamethod <ClassName> <methodName>
func (a *Analyzer) extractMethod(comment string) *methodResult {
	lines := strings.Split(comment, "\n")

	var (
		result      *methodResult
		description strings.Builder
	)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// @luamethod ClassName methodName
		if strings.HasPrefix(line, "@luamethod ") {
			parts := strings.Fields(strings.TrimPrefix(line, "@luamethod "))
			if len(parts) >= 2 {
				result = &methodResult{
					className: parts[0],
					method: &LuaMethod{
						Name:              parts[1],
						Params:            make([]*LuaParam, 0),
						Returns:           make([]*LuaReturn, 0),
						CustomAnnotations: make([]string, 0),
					},
				}
			}
			continue
		}

		if result == nil {
			// Collect description before any annotations
			if !strings.HasPrefix(line, "@") && line != "" {
				if description.Len() > 0 {
					description.WriteString(" ")
				}
				description.WriteString(line)
			}
			continue
		}

		if strings.HasPrefix(line, "@luaparam ") {
			param := a.parseParam(line)
			if param != nil {
				result.method.Params = append(result.method.Params, param)
			}
		} else if strings.HasPrefix(line, "@luareturn ") {
			ret := a.parseReturn(line)
			if ret != nil {
				result.method.Returns = append(result.method.Returns, ret)
			}
		} else if strings.HasPrefix(line, "@luaannotation ") {
			annotation := strings.TrimSpace(strings.TrimPrefix(line, "@luaannotation "))
			if annotation != "" {
				result.method.CustomAnnotations = append(result.method.CustomAnnotations, annotation)
			}
		}
	}

	if result != nil {
		result.method.Description = description.String()
	}

	return result
}

// parseParam: parses @luaparam annotation
// Format: @luaparam <name> <type> <description>
func (a *Analyzer) parseParam(line string) *LuaParam {
	line = strings.TrimSpace(strings.TrimPrefix(line, "@luaparam "))
	parts := strings.SplitN(line, " ", 3)

	if len(parts) < 2 {
		return nil
	}

	param := &LuaParam{
		Name: parts[0],
		Type: parts[1],
	}

	if len(parts) >= 3 {
		param.Description = parts[2]
	}

	return param
}

// parseReturn: parses @luareturn annotation
// Format: @luareturn <type> <description>
func (a *Analyzer) parseReturn(line string) *LuaReturn {
	line = strings.TrimSpace(strings.TrimPrefix(line, "@luareturn "))
	parts := strings.SplitN(line, " ", 2)

	if len(parts) < 1 {
		return nil
	}

	ret := &LuaReturn{
		Type: parts[0],
	}

	if len(parts) >= 2 {
		ret.Description = parts[1]
	}

	return ret
}

// GenerateStubs: generates Lua annotation stubs for all discovered modules
func (a *Analyzer) GenerateStubs() (string, error) {
	var sb strings.Builder

	// Add meta annotation for Lua LSP
	sb.WriteString("---@meta\n\n")

	// Sort module names for consistent output
	var moduleNames []string
	for name := range a.modules {
		moduleNames = append(moduleNames, name)
	}
	sort.Strings(moduleNames)

	for _, moduleName := range moduleNames {
		module := a.modules[moduleName]

		// Generate module comment
		sb.WriteString(fmt.Sprintf("--- %s module\n", moduleName))

		// Module-level custom annotations
		for _, annotation := range module.CustomAnnotations {
			sb.WriteString(fmt.Sprintf("---%s\n", annotation))
		}

		sb.WriteString(fmt.Sprintf("---@class %s\n", moduleName))
		sb.WriteString(fmt.Sprintf("local %s = {}\n\n", moduleName))

		// Generate function stubs
		for _, fn := range module.Functions {
			// Parameter annotations
			for _, param := range fn.Params {
				if param.Description != "" {
					sb.WriteString(fmt.Sprintf("---@param %s %s %s\n", param.Name, param.Type, param.Description))
				} else {
					sb.WriteString(fmt.Sprintf("---@param %s %s\n", param.Name, param.Type))
				}
			}

			// Return annotations
			for _, ret := range fn.Returns {
				if ret.Description != "" {
					sb.WriteString(fmt.Sprintf("---@return %s %s\n", ret.Type, ret.Description))
				} else {
					sb.WriteString(fmt.Sprintf("---@return %s\n", ret.Type))
				}
			}

			// Function-level custom annotations
			for _, annotation := range fn.CustomAnnotations {
				sb.WriteString(fmt.Sprintf("---%s\n", annotation))
			}

			// Function signature
			paramNames := make([]string, len(fn.Params))
			for i, param := range fn.Params {
				paramNames[i] = param.Name
			}
			sb.WriteString(fmt.Sprintf("function %s.%s(%s) end\n\n", moduleName, fn.Name, strings.Join(paramNames, ", ")))
		}
	}

	// Add single return statement at the end of the file
	if len(moduleNames) > 0 {
		sb.WriteString("return {}\n")
	}

	return sb.String(), nil
}

// ModuleCount: returns the number of modules discovered
func (a *Analyzer) ModuleCount() int {
	return len(a.modules)
}

// GetModules: returns a map of module names to module definitions
func (a *Analyzer) GetModules() map[string]*LuaModule {
	return a.modules
}

// GenerateModuleStub: generates a Lua annotation stub for a single module.
// This is the recommended format for per-module files in a library/ directory.
func (a *Analyzer) GenerateModuleStub(moduleName string) (string, error) {
	module, exists := a.modules[moduleName]
	if !exists {
		return "", fmt.Errorf("module %s not found", moduleName)
	}

	var sb strings.Builder

	// Add meta annotation for Lua LSP with module name
	sb.WriteString(fmt.Sprintf("---@meta %s\n\n", moduleName))

	// Module-level custom annotations
	for _, annotation := range module.CustomAnnotations {
		sb.WriteString(fmt.Sprintf("---%s\n", annotation))
	}

	// Generate class stubs FIRST (must come before module for proper LSP resolution)
	var namespacedClasses []string
	for _, class := range module.Classes {
		// Class declaration
		sb.WriteString(fmt.Sprintf("---@class %s\n", class.Name))

		// Class-level custom annotations
		for _, annotation := range class.CustomAnnotations {
			sb.WriteString(fmt.Sprintf("---%s\n", annotation))
		}

		// Field annotations
		for _, field := range class.Fields {
			if field.Description != "" {
				sb.WriteString(fmt.Sprintf("---@field %s %s %s\n", field.Name, field.Type, field.Description))
			} else {
				sb.WriteString(fmt.Sprintf("---@field %s %s\n", field.Name, field.Type))
			}
		}

		// Determine local variable name
		var localVarName string
		if strings.HasPrefix(class.Name, moduleName+".") {
			// Namespaced class (e.g., "log.Logger" -> "Logger")
			localVarName = strings.TrimPrefix(class.Name, moduleName+".")
			namespacedClasses = append(namespacedClasses, localVarName)
		} else {
			// Simple class name
			localVarName = class.Name
		}

		// Only generate local variable for classes with methods
		// Pure data classes (only fields, no methods) don't need a local variable
		if len(class.Methods) > 0 {
			sb.WriteString(fmt.Sprintf("local %s = {}\n\n", localVarName))
		} else {
			sb.WriteString("\n")
		}

		// Generate method stubs
		for _, method := range class.Methods {
			// Parameter annotations (skip self since it's in the method signature comment)
			for _, param := range method.Params {
				if param.Name == "self" {
					continue // Skip self parameter in annotations
				}
				if param.Description != "" {
					sb.WriteString(fmt.Sprintf("---@param %s %s %s\n", param.Name, param.Type, param.Description))
				} else {
					sb.WriteString(fmt.Sprintf("---@param %s %s\n", param.Name, param.Type))
				}
			}

			// Return annotations
			for _, ret := range method.Returns {
				if ret.Description != "" {
					sb.WriteString(fmt.Sprintf("---@return %s %s\n", ret.Type, ret.Description))
				} else {
					sb.WriteString(fmt.Sprintf("---@return %s\n", ret.Type))
				}
			}

			// Method-level custom annotations
			for _, annotation := range method.CustomAnnotations {
				sb.WriteString(fmt.Sprintf("---%s\n", annotation))
			}

			// Method signature (using : for method notation, using local variable name)
			// Skip the first parameter if it's 'self' since : notation provides it implicitly
			var paramNames []string
			for i, param := range method.Params {
				if i == 0 && param.Name == "self" {
					continue // Skip self parameter
				}
				paramNames = append(paramNames, param.Name)
			}
			sb.WriteString(fmt.Sprintf("function %s:%s(%s) end\n\n", localVarName, method.Name, strings.Join(paramNames, ", ")))
		}
	}

	// Generate module class
	sb.WriteString(fmt.Sprintf("---@class %s\n", moduleName))

	// Add field annotations for namespaced classes (e.g., log.Logger)
	for _, class := range module.Classes {
		if strings.HasPrefix(class.Name, moduleName+".") {
			// Extract the field name (e.g., "Logger" from "log.Logger")
			fieldName := strings.TrimPrefix(class.Name, moduleName+".")
			sb.WriteString(fmt.Sprintf("---@field %s %s\n", fieldName, class.Name))
		}
	}

	sb.WriteString(fmt.Sprintf("local %s = {}\n\n", moduleName))

	// Generate function stubs
	for _, fn := range module.Functions {
		// Parameter annotations
		for _, param := range fn.Params {
			if param.Description != "" {
				sb.WriteString(fmt.Sprintf("---@param %s %s %s\n", param.Name, param.Type, param.Description))
			} else {
				sb.WriteString(fmt.Sprintf("---@param %s %s\n", param.Name, param.Type))
			}
		}

		// Return annotations
		for _, ret := range fn.Returns {
			if ret.Description != "" {
				sb.WriteString(fmt.Sprintf("---@return %s %s\n", ret.Type, ret.Description))
			} else {
				sb.WriteString(fmt.Sprintf("---@return %s\n", ret.Type))
			}
		}

		// Function-level custom annotations
		for _, annotation := range fn.CustomAnnotations {
			sb.WriteString(fmt.Sprintf("---%s\n", annotation))
		}

		// Function signature
		paramNames := make([]string, len(fn.Params))
		for i, param := range fn.Params {
			paramNames[i] = param.Name
		}
		sb.WriteString(fmt.Sprintf("function %s.%s(%s) end\n\n", moduleName, fn.Name, strings.Join(paramNames, ", ")))
	}

	// Generate constant declarations
	for _, cnst := range module.Constants {
		// Constant annotation
		if cnst.Description != "" {
			sb.WriteString(fmt.Sprintf("---@type %s %s\n", cnst.Type, cnst.Description))
		} else {
			sb.WriteString(fmt.Sprintf("---@type %s\n", cnst.Type))
		}
		sb.WriteString(fmt.Sprintf("%s.%s = nil\n\n", moduleName, cnst.Name))
	}

	// Assign namespaced classes to module fields (e.g., log.Logger = Logger)
	for _, fieldName := range namespacedClasses {
		sb.WriteString(fmt.Sprintf("%s.%s = %s\n\n", moduleName, fieldName, fieldName))
	}

	// Return the module
	sb.WriteString(fmt.Sprintf("return %s\n", moduleName))

	return sb.String(), nil
}
