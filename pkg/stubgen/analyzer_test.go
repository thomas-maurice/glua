package stubgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzer_ExtractModuleName(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "valid module",
			comment:  "@luamodule mymodule",
			expected: "mymodule",
		},
		{
			name:     "module with description",
			comment:  "Some description\n@luamodule testmodule\nMore text",
			expected: "testmodule",
		},
		{
			name:     "no module",
			comment:  "Just a regular comment",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := a.extractModuleName(tt.comment)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAnalyzer_ParseParam(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		name     string
		line     string
		expected *LuaParam
	}{
		{
			name: "param with description",
			line: "@luaparam name string The parameter name",
			expected: &LuaParam{
				Name:        "name",
				Type:        "string",
				Description: "The parameter name",
			},
		},
		{
			name: "param without description",
			line: "@luaparam count number",
			expected: &LuaParam{
				Name: "count",
				Type: "number",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := a.parseParam(tt.line)
			if result == nil {
				t.Fatal("Expected param, got nil")
			}

			if result.Name != tt.expected.Name {
				t.Errorf("Expected name %q, got %q", tt.expected.Name, result.Name)
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Expected type %q, got %q", tt.expected.Type, result.Type)
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Expected description %q, got %q", tt.expected.Description, result.Description)
			}
		})
	}
}

func TestAnalyzer_ParseReturn(t *testing.T) {
	a := NewAnalyzer()

	tests := []struct {
		name     string
		line     string
		expected *LuaReturn
	}{
		{
			name: "return with description",
			line: "@luareturn number The result value",
			expected: &LuaReturn{
				Type:        "number",
				Description: "The result value",
			},
		},
		{
			name: "return without description",
			line: "@luareturn boolean",
			expected: &LuaReturn{
				Type: "boolean",
			},
		},
		{
			name: "return with complex type",
			line: "@luareturn string|nil Error message or nil",
			expected: &LuaReturn{
				Type:        "string|nil",
				Description: "Error message or nil",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := a.parseReturn(tt.line)
			if result == nil {
				t.Fatal("Expected return, got nil")
			}

			if result.Type != tt.expected.Type {
				t.Errorf("Expected type %q, got %q", tt.expected.Type, result.Type)
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Expected description %q, got %q", tt.expected.Description, result.Description)
			}
		})
	}
}

func TestAnalyzer_ScanDirectory(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")

	testCode := `package test

// Loader: creates the test module
//
// @luamodule testmodule
func Loader(L *lua.LState) int {
	return 1
}

// testFunc: a test function
//
// @luafunc test_function
// @luaparam input string The input value
// @luareturn number The result
func testFunc(L *lua.LState) int {
	return 1
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	a := NewAnalyzer()
	if err := a.ScanDirectory(tmpDir); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if a.ModuleCount() != 1 {
		t.Errorf("Expected 1 module, got %d", a.ModuleCount())
	}

	module, exists := a.modules["testmodule"]
	if !exists {
		t.Fatal("Expected testmodule to be registered")
	}

	if len(module.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(module.Functions))
	}

	fn := module.Functions[0]
	if fn.Name != "test_function" {
		t.Errorf("Expected function name 'test_function', got %q", fn.Name)
	}

	if len(fn.Params) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(fn.Params))
	}

	if len(fn.Returns) != 1 {
		t.Errorf("Expected 1 return value, got %d", len(fn.Returns))
	}
}

func TestAnalyzer_GenerateStubs(t *testing.T) {
	a := NewAnalyzer()

	// Manually add a test module
	a.modules["testmod"] = &LuaModule{
		Name: "testmod",
		Functions: []*LuaFunction{
			{
				Name:        "add",
				Description: "Adds two numbers",
				Params: []*LuaParam{
					{Name: "a", Type: "number", Description: "First number"},
					{Name: "b", Type: "number", Description: "Second number"},
				},
				Returns: []*LuaReturn{
					{Type: "number", Description: "The sum"},
				},
			},
		},
	}

	stubs, err := a.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that output contains expected elements
	expectedStrings := []string{
		"---@class testmod",
		"function testmod.add(a, b) end",
		"---@param a number First number",
		"---@param b number Second number",
		"---@return number The sum",
		"return {}",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(stubs, expected) {
			t.Errorf("Expected stubs to contain %q, got:\n%s", expected, stubs)
		}
	}
}

func TestAnalyzer_RealModule(t *testing.T) {
	a := NewAnalyzer()

	// Test with the actual kubernetes module
	if err := a.ScanDirectory("../modules/kubernetes"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if a.ModuleCount() != 1 {
		t.Errorf("Expected 1 module, got %d", a.ModuleCount())
	}

	module, exists := a.modules["kubernetes"]
	if !exists {
		t.Fatal("Expected kubernetes module to be registered")
	}

	// Should have 11 functions
	if len(module.Functions) != 11 {
		t.Errorf("Expected 11 functions, got %d", len(module.Functions))
	}

	// Verify function names
	fnNames := make(map[string]bool)
	for _, fn := range module.Functions {
		fnNames[fn.Name] = true
	}

	expectedFuncs := []string{
		"parse_memory", "parse_cpu", "parse_time", "format_time", "init_defaults",
		"parse_duration", "format_duration", "parse_int_or_string",
		"matches_selector", "toleration_matches", "match_gvk",
	}
	for _, name := range expectedFuncs {
		if !fnNames[name] {
			t.Errorf("Expected function %q not found", name)
		}
	}

	// Generate stubs
	stubs, err := a.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Verify output
	if !strings.Contains(stubs, "---@class kubernetes") {
		t.Error("Expected kubernetes class annotation")
	}

	if !strings.Contains(stubs, "function kubernetes.parse_memory(quantity) end") {
		t.Error("Expected parse_memory function")
	}
}

func TestAnalyzer_CustomAnnotations(t *testing.T) {
	a := NewAnalyzer()

	// Test with the custom_annotations testdata
	if err := a.ScanDirectory("testdata"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	module, exists := a.modules["custom_annotations"]
	if !exists {
		t.Fatal("Expected custom_annotations module to be registered")
	}

	// Check module-level custom annotations
	if len(module.CustomAnnotations) != 2 {
		t.Errorf("Expected 2 module-level annotations, got %d", len(module.CustomAnnotations))
	}

	expectedModuleAnnotations := []string{
		"@alias ID string|number",
		"@alias Handler fun(id: ID): boolean",
	}

	for i, expected := range expectedModuleAnnotations {
		if i >= len(module.CustomAnnotations) {
			t.Errorf("Missing module annotation: %q", expected)
			continue
		}
		if module.CustomAnnotations[i] != expected {
			t.Errorf("Expected module annotation %q, got %q", expected, module.CustomAnnotations[i])
		}
	}

	// Check function-level custom annotations
	if len(module.Functions) != 2 {
		t.Fatalf("Expected 2 functions, got %d", len(module.Functions))
	}

	// First function should have deprecation and nodiscard annotations
	fn1 := module.Functions[0]
	if len(fn1.CustomAnnotations) != 2 {
		t.Errorf("Expected 2 annotations on %s, got %d", fn1.Name, len(fn1.CustomAnnotations))
	}

	expectedFn1Annotations := []string{
		"@deprecated Use process_typed_id instead",
		"@nodiscard",
	}

	for i, expected := range expectedFn1Annotations {
		if i >= len(fn1.CustomAnnotations) {
			t.Errorf("Missing function annotation: %q", expected)
			continue
		}
		if fn1.CustomAnnotations[i] != expected {
			t.Errorf("Expected function annotation %q, got %q", expected, fn1.CustomAnnotations[i])
		}
	}

	// Second function should have generic annotations
	fn2 := module.Functions[1]
	if len(fn2.CustomAnnotations) != 3 {
		t.Errorf("Expected 3 annotations on %s, got %d", fn2.Name, len(fn2.CustomAnnotations))
	}

	// Generate stubs and verify they contain custom annotations
	stubs, err := a.GenerateModuleStub("custom_annotations")
	if err != nil {
		t.Fatalf("GenerateModuleStub failed: %v", err)
	}

	// Check that module-level annotations are present
	if !strings.Contains(stubs, "---@alias ID string|number") {
		t.Error("Expected module-level @alias annotation for ID")
	}

	if !strings.Contains(stubs, "---@alias Handler fun(id: ID): boolean") {
		t.Error("Expected module-level @alias annotation for Handler")
	}

	// Check that function-level annotations are present
	if !strings.Contains(stubs, "---@deprecated Use process_typed_id instead") {
		t.Error("Expected @deprecated annotation")
	}

	if !strings.Contains(stubs, "---@nodiscard") {
		t.Error("Expected @nodiscard annotation")
	}

	if !strings.Contains(stubs, "---@generic T") {
		t.Error("Expected @generic annotation")
	}
}
