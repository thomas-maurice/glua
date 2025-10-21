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
	Name        string
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
	classMap := make(map[string]*LuaClass)

	// Process function declarations
	currentModule = a.processFuncDeclarations(file, classMap)

	// Process standalone comments
	a.processStandaloneComments(file, currentModule, classMap)

	return nil
}

// processFuncDeclarations: processes all function declarations in a file
func (a *Analyzer) processFuncDeclarations(file *ast.File, classMap map[string]*LuaClass) *LuaModule {
	var currentModule *LuaModule

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Doc == nil {
			continue
		}

		comment := funcDecl.Doc.Text()

		// Try to extract module, function, const, or method
		if module := a.tryExtractModule(comment); module != nil {
			currentModule = module
			a.modules[module.Name] = module
		} else if currentModule != nil {
			a.tryAddFunctionOrConstOrMethod(comment, currentModule, classMap)
		}
	}

	return currentModule
}

// tryExtractModule: tries to extract a module from a comment
func (a *Analyzer) tryExtractModule(comment string) *LuaModule {
	if moduleName := a.extractModuleName(comment); moduleName != "" {
		return &LuaModule{
			Name:              moduleName,
			Functions:         make([]*LuaFunction, 0),
			Classes:           make([]*LuaClass, 0),
			CustomAnnotations: a.extractCustomAnnotations(comment),
		}
	}
	return nil
}

// tryAddFunctionOrConstOrMethod: tries to add a function, constant, or method to the current module
func (a *Analyzer) tryAddFunctionOrConstOrMethod(comment string, module *LuaModule, classMap map[string]*LuaClass) {
	if luaFunc := a.extractFunction(comment); luaFunc != nil {
		module.Functions = append(module.Functions, luaFunc)
		return
	}

	if luaConst := a.extractConst(comment); luaConst != nil {
		module.Constants = append(module.Constants, luaConst)
		return
	}

	if luaMethod := a.extractMethod(comment); luaMethod != nil {
		a.addMethodToClass(luaMethod, module, classMap)
	}
}

// addMethodToClass: adds a method to a class, creating the class if needed
func (a *Analyzer) addMethodToClass(luaMethod *methodResult, module *LuaModule, classMap map[string]*LuaClass) {
	className := luaMethod.className
	class, ok := classMap[className]
	if !ok {
		class = &LuaClass{
			Name:              className,
			Methods:           make([]*LuaMethod, 0),
			Fields:            make([]*LuaField, 0),
			CustomAnnotations: make([]string, 0),
		}
		classMap[className] = class
		module.Classes = append(module.Classes, class)
	}
	class.Methods = append(class.Methods, luaMethod.method)
}

// processStandaloneComments: processes standalone comments for constants and class definitions
func (a *Analyzer) processStandaloneComments(file *ast.File, module *LuaModule, classMap map[string]*LuaClass) {
	if module == nil {
		return
	}

	for _, commentGroup := range file.Comments {
		comment := commentGroup.Text()

		if luaConst := a.extractConst(comment); luaConst != nil {
			module.Constants = append(module.Constants, luaConst)
		}

		if classInfo := a.extractClassDefinition(comment); classInfo != nil {
			a.mergeOrAddClass(classInfo, module, classMap)
		}
	}
}

// mergeOrAddClass: merges class info into existing class or adds new class
func (a *Analyzer) mergeOrAddClass(classInfo *LuaClass, module *LuaModule, classMap map[string]*LuaClass) {
	existingClass, exists := classMap[classInfo.Name]
	if exists {
		existingClass.Fields = append(existingClass.Fields, classInfo.Fields...)
		if existingClass.Description == "" {
			existingClass.Description = classInfo.Description
		}
	} else {
		classMap[classInfo.Name] = classInfo
		module.Classes = append(module.Classes, classInfo)
	}
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
// Format: @luareturn <type> <name> <description>
// Or: @luareturn <type> <description> (name will be empty)
func (a *Analyzer) parseReturn(line string) *LuaReturn {
	line = strings.TrimSpace(strings.TrimPrefix(line, "@luareturn "))
	parts := strings.SplitN(line, " ", 3)

	if len(parts) < 1 {
		return nil
	}

	ret := &LuaReturn{
		Type: parts[0],
	}

	// If we have 3 parts, check if middle part is actually a name or part of description
	// If we have 2 parts, it could be: <type> <description> OR <type> <name>
	// We detect based on whether the second part starts with uppercase (description)
	// or lowercase (likely a name followed by description)
	if len(parts) >= 3 {
		// Check if parts[1] looks like a variable name (starts with lowercase)
		potentialName := parts[1]
		if len(potentialName) > 0 && potentialName[0] >= 'a' && potentialName[0] <= 'z' {
			// Looks like a name
			ret.Name = parts[1]
			ret.Description = parts[2]
		} else {
			// Starts with uppercase - it's part of the description, not a name
			ret.Description = parts[1] + " " + parts[2]
		}
	} else if len(parts) == 2 {
		secondPart := parts[1]

		// If it starts with uppercase or is empty, it's a description
		if len(secondPart) == 0 || (secondPart[0] >= 'A' && secondPart[0] <= 'Z') {
			ret.Description = secondPart
		} else {
			// Starts with lowercase - check if first word looks like a variable name
			words := strings.Fields(secondPart)
			if len(words) == 0 {
				ret.Description = secondPart
			} else {
				firstWord := words[0]
				// Check if first word is a valid identifier (lowercase start, alphanumeric+underscore)
				isValidIdentifier := true
				for i, ch := range firstWord {
					if i == 0 {
						if (ch < 'a' || ch > 'z') && ch != '_' {
							isValidIdentifier = false
							break
						}
					} else {
						if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') && ch != '_' {
							isValidIdentifier = false
							break
						}
					}
				}

				if isValidIdentifier && len(words) > 1 {
					// First word is a valid identifier and there are more words - treat as name
					ret.Name = firstWord
					ret.Description = strings.Join(words[1:], " ")
				} else {
					// Either not a valid identifier or only one word - treat whole thing as description
					ret.Description = secondPart
				}
			}
		}
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
					sb.WriteString(fmt.Sprintf("---@param %s %s %s\n", param.Type, param.Name, param.Description))
				} else {
					sb.WriteString(fmt.Sprintf("---@param %s %s\n", param.Type, param.Name))
				}
			}

			// Return annotations
			for _, ret := range fn.Returns {
				if ret.Name != "" && ret.Description != "" {
					sb.WriteString(fmt.Sprintf("---@return %s %s %s\n", ret.Type, ret.Name, ret.Description))
				} else if ret.Name != "" {
					sb.WriteString(fmt.Sprintf("---@return %s %s\n", ret.Type, ret.Name))
				} else if ret.Description != "" {
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
	a.writeAnnotations(&sb, module.CustomAnnotations)

	// Generate class stubs FIRST (must come before module for proper LSP resolution)
	namespacedClasses := a.generateClassStubs(&sb, module, moduleName)

	// Generate module declaration
	a.generateModuleDeclaration(&sb, module, moduleName)

	// Generate function stubs
	a.generateFunctionStubs(&sb, module, moduleName)

	// Generate constant declarations
	a.generateConstantStubs(&sb, module, moduleName)

	// Assign namespaced classes to module fields (e.g., log.Logger = Logger)
	for _, fieldName := range namespacedClasses {
		sb.WriteString(fmt.Sprintf("%s.%s = %s\n\n", moduleName, fieldName, fieldName))
	}

	// Return the module
	sb.WriteString(fmt.Sprintf("return %s\n", moduleName))

	return sb.String(), nil
}

// writeAnnotations: writes custom annotations to the string builder
func (a *Analyzer) writeAnnotations(sb *strings.Builder, annotations []string) {
	for _, annotation := range annotations {
		fmt.Fprintf(sb, "---%s\n", annotation)
	}
}

// generateClassStubs: generates Lua annotation stubs for all classes in a module.
// Returns a list of namespaced class names for later assignment.
func (a *Analyzer) generateClassStubs(sb *strings.Builder, module *LuaModule, moduleName string) []string {
	var namespacedClasses []string

	for _, class := range module.Classes {
		localVarName := a.generateSingleClassStub(sb, class, moduleName)
		if strings.HasPrefix(class.Name, moduleName+".") {
			namespacedClasses = append(namespacedClasses, localVarName)
		}
	}

	return namespacedClasses
}

// generateSingleClassStub: generates stub for a single class and returns its local variable name
func (a *Analyzer) generateSingleClassStub(sb *strings.Builder, class *LuaClass, moduleName string) string {
	// Class declaration
	fmt.Fprintf(sb, "---@class %s\n", class.Name)

	// Class-level custom annotations
	a.writeAnnotations(sb, class.CustomAnnotations)

	// Field annotations
	a.writeFieldAnnotations(sb, class.Fields)

	// Determine local variable name
	localVarName := a.getClassLocalName(class.Name, moduleName)

	// Only generate local variable for classes with methods
	if len(class.Methods) > 0 {
		fmt.Fprintf(sb, "local %s = {}\n\n", localVarName)
		a.generateMethodStubs(sb, class.Methods, localVarName)
	} else {
		sb.WriteString("\n")
	}

	return localVarName
}

// writeFieldAnnotations: writes field annotations for a class
func (a *Analyzer) writeFieldAnnotations(sb *strings.Builder, fields []*LuaField) {
	for _, field := range fields {
		if field.Description != "" {
			fmt.Fprintf(sb, "---@field %s %s %s\n", field.Name, field.Type, field.Description)
		} else {
			fmt.Fprintf(sb, "---@field %s %s\n", field.Name, field.Type)
		}
	}
}

// getClassLocalName: determines the local variable name for a class
func (a *Analyzer) getClassLocalName(className, moduleName string) string {
	if strings.HasPrefix(className, moduleName+".") {
		return strings.TrimPrefix(className, moduleName+".")
	}
	return className
}

// generateMethodStubs: generates method stub declarations for a class
func (a *Analyzer) generateMethodStubs(sb *strings.Builder, methods []*LuaMethod, localVarName string) {
	for _, method := range methods {
		a.writeParamAnnotations(sb, method.Params, true)
		a.writeReturnAnnotations(sb, method.Returns)
		a.writeAnnotations(sb, method.CustomAnnotations)

		paramNames := a.extractParamNames(method.Params, true)
		fmt.Fprintf(sb, "function %s:%s(%s) end\n\n", localVarName, method.Name, strings.Join(paramNames, ", "))
	}
}

// writeParamAnnotations: writes parameter annotations
func (a *Analyzer) writeParamAnnotations(sb *strings.Builder, params []*LuaParam, skipSelf bool) {
	for _, param := range params {
		if skipSelf && param.Name == "self" {
			continue
		}
		if param.Description != "" {
			fmt.Fprintf(sb, "---@param %s %s %s\n", param.Type, param.Name, param.Description)
		} else {
			fmt.Fprintf(sb, "---@param %s %s\n", param.Type, param.Name)
		}
	}
}

// writeReturnAnnotations: writes return value annotations
func (a *Analyzer) writeReturnAnnotations(sb *strings.Builder, returns []*LuaReturn) {
	for _, ret := range returns {
		if ret.Name != "" && ret.Description != "" {
			fmt.Fprintf(sb, "---@return %s %s %s\n", ret.Type, ret.Name, ret.Description)
		} else if ret.Name != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", ret.Type, ret.Name)
		} else if ret.Description != "" {
			fmt.Fprintf(sb, "---@return %s %s\n", ret.Type, ret.Description)
		} else {
			fmt.Fprintf(sb, "---@return %s\n", ret.Type)
		}
	}
}

// extractParamNames: extracts parameter names from a list of parameters
func (a *Analyzer) extractParamNames(params []*LuaParam, skipSelf bool) []string {
	var paramNames []string
	for i, param := range params {
		if skipSelf && i == 0 && param.Name == "self" {
			continue
		}
		paramNames = append(paramNames, param.Name)
	}
	return paramNames
}

// generateModuleDeclaration: generates the module class declaration
func (a *Analyzer) generateModuleDeclaration(sb *strings.Builder, module *LuaModule, moduleName string) {
	fmt.Fprintf(sb, "---@class %s\n", moduleName)

	// Add field annotations for namespaced classes
	for _, class := range module.Classes {
		if strings.HasPrefix(class.Name, moduleName+".") {
			fieldName := strings.TrimPrefix(class.Name, moduleName+".")
			fmt.Fprintf(sb, "---@field %s %s\n", fieldName, class.Name)
		}
	}

	fmt.Fprintf(sb, "local %s = {}\n\n", moduleName)
}

// generateFunctionStubs: generates function stub declarations for a module
func (a *Analyzer) generateFunctionStubs(sb *strings.Builder, module *LuaModule, moduleName string) {
	for _, fn := range module.Functions {
		a.writeParamAnnotations(sb, fn.Params, false)
		a.writeReturnAnnotations(sb, fn.Returns)
		a.writeAnnotations(sb, fn.CustomAnnotations)

		paramNames := make([]string, len(fn.Params))
		for i, param := range fn.Params {
			paramNames[i] = param.Name
		}
		fmt.Fprintf(sb, "function %s.%s(%s) end\n\n", moduleName, fn.Name, strings.Join(paramNames, ", "))
	}
}

// generateConstantStubs: generates constant declarations for a module
func (a *Analyzer) generateConstantStubs(sb *strings.Builder, module *LuaModule, moduleName string) {
	for _, cnst := range module.Constants {
		if cnst.Description != "" {
			fmt.Fprintf(sb, "---@type %s %s\n", cnst.Type, cnst.Description)
		} else {
			fmt.Fprintf(sb, "---@type %s\n", cnst.Type)
		}
		fmt.Fprintf(sb, "%s.%s = nil\n\n", moduleName, cnst.Name)
	}
}
