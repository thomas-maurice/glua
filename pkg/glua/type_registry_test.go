package glua

import (
	"strings"
	"testing"
)

func TestTypeRegistry_SimpleStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Person{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that the stub contains the Person class
	if !strings.Contains(stubs, "---@class glua.Person") {
		t.Errorf("Expected stub to contain '---@class glua.Person', got:\n%s", stubs)
	}

	// Check that the stub contains the name field
	if !strings.Contains(stubs, "---@field name string") {
		t.Errorf("Expected stub to contain '---@field name string', got:\n%s", stubs)
	}

	// Check that the stub contains the age field
	if !strings.Contains(stubs, "---@field age number") {
		t.Errorf("Expected stub to contain '---@field age number', got:\n%s", stubs)
	}
}

func TestTypeRegistry_NestedStructs(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Address Address `json:"address"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Person{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that both types are registered
	if !strings.Contains(stubs, "---@class glua.Person") {
		t.Errorf("Expected stub to contain '---@class glua.Person', got:\n%s", stubs)
	}

	if !strings.Contains(stubs, "---@class glua.Address") {
		t.Errorf("Expected stub to contain '---@class glua.Address', got:\n%s", stubs)
	}

	// Check that Person has Address field with correct type
	if !strings.Contains(stubs, "---@field address glua.Address") {
		t.Errorf("Expected stub to contain '---@field address glua.Address', got:\n%s", stubs)
	}
}

func TestTypeRegistry_Slices(t *testing.T) {
	type Container struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	type Pod struct {
		Name       string      `json:"name"`
		Containers []Container `json:"containers"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Pod{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that both types are registered
	if !strings.Contains(stubs, "---@class glua.Pod") {
		t.Errorf("Expected stub to contain '---@class glua.Pod', got:\n%s", stubs)
	}

	if !strings.Contains(stubs, "---@class glua.Container") {
		t.Errorf("Expected stub to contain '---@class glua.Container', got:\n%s", stubs)
	}

	// Check that containers field is an array of Container
	if !strings.Contains(stubs, "---@field containers glua.Container[]") {
		t.Errorf("Expected stub to contain '---@field containers glua.Container[]', got:\n%s", stubs)
	}
}

func TestTypeRegistry_CircularDependency(t *testing.T) {
	type Node struct {
		Value    string  `json:"value"`
		Children []*Node `json:"children"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Node{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that Node class is registered
	if !strings.Contains(stubs, "---@class glua.Node") {
		t.Errorf("Expected stub to contain '---@class glua.Node', got:\n%s", stubs)
	}

	// Check that children field references Node (circular reference)
	if !strings.Contains(stubs, "---@field children glua.Node[]") {
		t.Errorf("Expected stub to contain '---@field children glua.Node[]', got:\n%s", stubs)
	}

	// Ensure we don't have duplicate Node definitions
	count := strings.Count(stubs, "---@class glua.Node")
	if count != 1 {
		t.Errorf("Expected exactly 1 Node class definition, got %d:\n%s", count, stubs)
	}
}

func TestTypeRegistry_PrimitiveTypes(t *testing.T) {
	type Types struct {
		StringVal string  `json:"stringVal"`
		IntVal    int     `json:"intVal"`
		FloatVal  float64 `json:"floatVal"`
		BoolVal   bool    `json:"boolVal"`
		Int64Val  int64   `json:"int64Val"`
		Uint32Val uint32  `json:"uint32Val"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Types{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	tests := []string{
		"---@field stringVal string",
		"---@field intVal number",
		"---@field floatVal number",
		"---@field boolVal boolean",
		"---@field int64Val number",
		"---@field uint32Val number",
	}

	for _, expected := range tests {
		if !strings.Contains(stubs, expected) {
			t.Errorf("Expected stub to contain '%s', got:\n%s", expected, stubs)
		}
	}
}

func TestTypeRegistry_Maps(t *testing.T) {
	type Config struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Config{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that labels field is a map
	if !strings.Contains(stubs, "---@field labels table<string, string>") {
		t.Errorf("Expected stub to contain '---@field labels table<string, string>', got:\n%s", stubs)
	}
}

func TestTypeRegistry_SkipUnexportedFields(t *testing.T) {
	type Private struct {
		Public  string `json:"public"`
		private string //nolint:unused
	}

	registry := NewTypeRegistry()
	err := registry.Register(Private{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that public field is included
	if !strings.Contains(stubs, "---@field public string") {
		t.Errorf("Expected stub to contain '---@field public string', got:\n%s", stubs)
	}

	// Check that private field is NOT included
	if strings.Contains(stubs, "---@field private") {
		t.Errorf("Did not expect stub to contain '---@field private', got:\n%s", stubs)
	}
}

func TestTypeRegistry_SkipIgnoredJSONFields(t *testing.T) {
	type Ignored struct {
		Included string `json:"included"`
		Skipped  string `json:"-"`
		NoTag    string
	}

	registry := NewTypeRegistry()
	err := registry.Register(Ignored{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that included field is present
	if !strings.Contains(stubs, "---@field included string") {
		t.Errorf("Expected stub to contain '---@field included string', got:\n%s", stubs)
	}

	// Check that skipped field is NOT present
	if strings.Contains(stubs, "---@field Skipped") || strings.Contains(stubs, "---@field skipped") {
		t.Errorf("Did not expect stub to contain skipped field, got:\n%s", stubs)
	}

	// Check that NoTag field is NOT present (no JSON tag)
	if strings.Contains(stubs, "---@field NoTag") {
		t.Errorf("Did not expect stub to contain NoTag field, got:\n%s", stubs)
	}
}

func TestTypeRegistry_Pointers(t *testing.T) {
	type Inner struct {
		Value string `json:"value"`
	}

	type Outer struct {
		Name     string `json:"name"`
		InnerPtr *Inner `json:"innerPtr"`
	}

	registry := NewTypeRegistry()
	err := registry.Register(Outer{})
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	err = registry.Process()
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	stubs, err := registry.GenerateStubs()
	if err != nil {
		t.Fatalf("GenerateStubs failed: %v", err)
	}

	// Check that pointer types are handled correctly
	if !strings.Contains(stubs, "---@field innerPtr glua.Inner") {
		t.Errorf("Expected stub to contain '---@field innerPtr glua.Inner', got:\n%s", stubs)
	}
}
