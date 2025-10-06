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

// LuaModule represents a discovered Lua module
type LuaModule struct {
	Name               string
	Functions          []*LuaFunction
	CustomAnnotations  []string // Module-level custom annotations
}

// LuaFunction represents a Lua function exported by a module
type LuaFunction struct {
	Name               string
	Description        string
	Params             []*LuaParam
	Returns            []*LuaReturn
	CustomAnnotations  []string // Function-level custom annotations
}

// LuaParam represents a function parameter
type LuaParam struct {
	Name        string
	Type        string
	Description string
}

// LuaReturn represents a function return value
type LuaReturn struct {
	Type        string
	Description string
}

// Analyzer scans Go source files and extracts Lua module definitions
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
				Name:               moduleName,
				Functions:          make([]*LuaFunction, 0),
				CustomAnnotations:  a.extractCustomAnnotations(comment),
			}
			a.modules[moduleName] = currentModule
			continue
		}

		// Check if this is a Lua function
		if luaFunc := a.extractFunction(comment); luaFunc != nil && currentModule != nil {
			currentModule.Functions = append(currentModule.Functions, luaFunc)
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
			// Function description
			if fn.Description != "" {
				sb.WriteString(fmt.Sprintf("--- %s\n", fn.Description))
			}

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

		sb.WriteString(fmt.Sprintf("return %s\n", moduleName))
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

	// Add meta annotation for Lua LSP
	sb.WriteString("---@meta\n\n")

	// Module-level custom annotations
	for _, annotation := range module.CustomAnnotations {
		sb.WriteString(fmt.Sprintf("---%s\n", annotation))
	}

	// Generate module class
	sb.WriteString(fmt.Sprintf("---@class %s\n", moduleName))
	sb.WriteString(fmt.Sprintf("local %s = {}\n\n", moduleName))

	// Generate function stubs
	for _, fn := range module.Functions {
		// Function description
		if fn.Description != "" {
			sb.WriteString(fmt.Sprintf("--- %s\n", fn.Description))
		}

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

	sb.WriteString(fmt.Sprintf("return %s\n", moduleName))

	return sb.String(), nil
}
