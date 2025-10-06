package glua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestTranslator_ToLua_Primitives(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool", true, "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tr.ToLua(L, tt.input)
			if err != nil {
				t.Fatalf("ToLua failed: %v", err)
			}

			if result.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result.String())
			}
		})
	}
}

func TestTranslator_ToLua_Struct(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	input := Person{Name: "Alice", Age: 30}
	result, err := tr.ToLua(L, input)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	table, ok := result.(*lua.LTable)
	if !ok {
		t.Fatalf("expected LTable, got %T", result)
	}

	name := table.RawGetString("name")
	if name.String() != "Alice" {
		t.Errorf("expected name 'Alice', got %s", name.String())
	}

	age := table.RawGetString("age")
	if age.String() != "30" {
		t.Errorf("expected age '30', got %s", age.String())
	}
}

func TestTranslator_ToLua_Slice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	input := []string{"a", "b", "c"}
	result, err := tr.ToLua(L, input)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	table, ok := result.(*lua.LTable)
	if !ok {
		t.Fatalf("expected LTable, got %T", result)
	}

	if table.Len() != 3 {
		t.Errorf("expected length 3, got %d", table.Len())
	}
}

func TestTranslator_FromLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Create a Lua table
	table := L.NewTable()
	table.RawSetString("name", lua.LString("Bob"))
	table.RawSetString("age", lua.LNumber(25))

	var output Person
	err := tr.FromLua(L, table, &output)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if output.Name != "Bob" {
		t.Errorf("expected name 'Bob', got %s", output.Name)
	}

	if output.Age != 25 {
		t.Errorf("expected age 25, got %d", output.Age)
	}
}

func TestTranslator_RoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}

	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	original := Person{
		Name: "Charlie",
		Age:  35,
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	// Convert to Lua
	luaVal, err := tr.ToLua(L, original)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	table, ok := luaVal.(*lua.LTable)
	if !ok {
		t.Fatalf("expected LTable, got %T", luaVal)
	}

	// Convert back to Go
	var result Person
	err = tr.FromLua(L, table, &result)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if result.Name != original.Name {
		t.Errorf("expected name %s, got %s", original.Name, result.Name)
	}

	if result.Age != original.Age {
		t.Errorf("expected age %d, got %d", original.Age, result.Age)
	}

	if result.Address.Street != original.Address.Street {
		t.Errorf("expected street %s, got %s", original.Address.Street, result.Address.Street)
	}

	if result.Address.City != original.Address.City {
		t.Errorf("expected city %s, got %s", original.Address.City, result.Address.City)
	}
}

func TestTranslator_FromLua_Primitives(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	tests := []struct {
		name     string
		luaValue lua.LValue
		expected interface{}
	}{
		{"string", lua.LString("hello"), "hello"},
		{"number", lua.LNumber(42), float64(42)},
		{"bool", lua.LBool(true), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output interface{}
			err := tr.FromLua(L, tt.luaValue, &output)
			if err != nil {
				t.Fatalf("FromLua failed: %v", err)
			}

			switch expected := tt.expected.(type) {
			case string:
				if output != expected {
					t.Errorf("expected %v, got %v", expected, output)
				}
			case float64:
				if output != expected {
					t.Errorf("expected %v, got %v", expected, output)
				}
			case bool:
				if output != expected {
					t.Errorf("expected %v, got %v", expected, output)
				}
			}
		})
	}
}

func TestTranslator_FromLua_PrimitiveToStruct(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	// Test converting Lua string to Go string
	var strOutput string
	err := tr.FromLua(L, lua.LString("test"), &strOutput)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}
	if strOutput != "test" {
		t.Errorf("expected 'test', got %s", strOutput)
	}

	// Test converting Lua number to Go int
	var intOutput int
	err = tr.FromLua(L, lua.LNumber(123), &intOutput)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}
	if intOutput != 123 {
		t.Errorf("expected 123, got %d", intOutput)
	}

	// Test converting Lua bool to Go bool
	var boolOutput bool
	err = tr.FromLua(L, lua.LBool(false), &boolOutput)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}
	if boolOutput != false {
		t.Errorf("expected false, got %v", boolOutput)
	}
}

func TestTranslator_RoundTrip_KubernetesPod(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	// Import the sample package to get a real pod
	pod := getSamplePod()

	// Convert Go Pod to Lua
	luaVal, err := tr.ToLua(L, pod)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	table, ok := luaVal.(*lua.LTable)
	if !ok {
		t.Fatalf("expected LTable, got %T", luaVal)
	}

	// Verify some key fields are present in Lua
	metadata := table.RawGetString("metadata")
	if metadata == lua.LNil {
		t.Fatal("metadata field is nil")
	}

	kind := table.RawGetString("kind")
	if kind.String() != "Pod" {
		t.Errorf("expected kind 'Pod', got %s", kind.String())
	}

	// Convert back to Go
	var reconstructedPod podStruct
	err = tr.FromLua(L, table, &reconstructedPod)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	// Verify key fields
	if reconstructedPod.Kind != pod.Kind {
		t.Errorf("kind mismatch: expected %s, got %s", pod.Kind, reconstructedPod.Kind)
	}

	if reconstructedPod.APIVersion != pod.APIVersion {
		t.Errorf("apiVersion mismatch: expected %s, got %s", pod.APIVersion, reconstructedPod.APIVersion)
	}

	if reconstructedPod.Metadata.Name != pod.Metadata.Name {
		t.Errorf("name mismatch: expected %s, got %s", pod.Metadata.Name, reconstructedPod.Metadata.Name)
	}

	if reconstructedPod.Metadata.Namespace != pod.Metadata.Namespace {
		t.Errorf("namespace mismatch: expected %s, got %s", pod.Metadata.Namespace, reconstructedPod.Metadata.Namespace)
	}
}

func TestTranslator_RoundTrip_Timestamp(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	// Create a struct with a timestamp field
	type TimestampTest struct {
		Name      string `json:"name"`
		CreatedAt string `json:"createdAt"`
	}

	original := TimestampTest{
		Name:      "test",
		CreatedAt: "2025-10-03T16:39:00Z",
	}

	// Convert to Lua
	luaVal, err := tr.ToLua(L, original)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	// Convert back to Go
	var result TimestampTest
	err = tr.FromLua(L, luaVal, &result)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if result.CreatedAt != original.CreatedAt {
		t.Errorf("timestamp mismatch: expected %s, got %s", original.CreatedAt, result.CreatedAt)
	}
}

func TestTranslator_RoundTrip_ResourceQuantities(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	// Create a struct with resource quantity strings
	type ResourceTest struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	}

	original := ResourceTest{
		CPU:    "100m",
		Memory: "100Mi",
	}

	// Convert to Lua
	luaVal, err := tr.ToLua(L, original)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	table, ok := luaVal.(*lua.LTable)
	if !ok {
		t.Fatalf("expected LTable, got %T", luaVal)
	}

	// Verify values in Lua
	cpu := table.RawGetString("cpu")
	if cpu.String() != "100m" {
		t.Errorf("expected cpu '100m', got %s", cpu.String())
	}

	memory := table.RawGetString("memory")
	if memory.String() != "100Mi" {
		t.Errorf("expected memory '100Mi', got %s", memory.String())
	}

	// Convert back to Go
	var result ResourceTest
	err = tr.FromLua(L, table, &result)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if result.CPU != original.CPU {
		t.Errorf("CPU mismatch: expected %s, got %s", original.CPU, result.CPU)
	}

	if result.Memory != original.Memory {
		t.Errorf("Memory mismatch: expected %s, got %s", original.Memory, result.Memory)
	}
}

func TestTranslator_RoundTrip_Maps(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	type MapTest struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	}

	original := MapTest{
		Name: "test",
		Labels: map[string]string{
			"app":    "myapp",
			"env":    "prod",
			"mutate": "true",
		},
	}

	// Convert to Lua
	luaVal, err := tr.ToLua(L, original)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	// Convert back to Go
	var result MapTest
	err = tr.FromLua(L, luaVal, &result)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if result.Name != original.Name {
		t.Errorf("name mismatch: expected %s, got %s", original.Name, result.Name)
	}

	if len(result.Labels) != len(original.Labels) {
		t.Errorf("labels length mismatch: expected %d, got %d", len(original.Labels), len(result.Labels))
	}

	for key, val := range original.Labels {
		if result.Labels[key] != val {
			t.Errorf("label %s mismatch: expected %s, got %s", key, val, result.Labels[key])
		}
	}
}

func TestTranslator_RoundTrip_NestedArrays(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := &Translator{}

	type EnvVar struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	type Container struct {
		Name  string   `json:"name"`
		Image string   `json:"image"`
		Env   []EnvVar `json:"env"`
	}

	original := Container{
		Name:  "mycontainer",
		Image: "alpine",
		Env: []EnvVar{
			{Name: "FOO", Value: "bar"},
			{Name: "BAZ", Value: "qux"},
		},
	}

	// Convert to Lua
	luaVal, err := tr.ToLua(L, original)
	if err != nil {
		t.Fatalf("ToLua failed: %v", err)
	}

	// Convert back to Go
	var result Container
	err = tr.FromLua(L, luaVal, &result)
	if err != nil {
		t.Fatalf("FromLua failed: %v", err)
	}

	if result.Name != original.Name {
		t.Errorf("name mismatch: expected %s, got %s", original.Name, result.Name)
	}

	if len(result.Env) != len(original.Env) {
		t.Fatalf("env length mismatch: expected %d, got %d", len(original.Env), len(result.Env))
	}

	for i, env := range original.Env {
		if result.Env[i].Name != env.Name {
			t.Errorf("env[%d].name mismatch: expected %s, got %s", i, env.Name, result.Env[i].Name)
		}
		if result.Env[i].Value != env.Value {
			t.Errorf("env[%d].value mismatch: expected %s, got %s", i, env.Value, result.Env[i].Value)
		}
	}
}

// Helper types and functions for Kubernetes Pod testing
type podStruct struct {
	Kind       string       `json:"kind"`
	APIVersion string       `json:"apiVersion"`
	Metadata   metadataStruct `json:"metadata"`
}

type metadataStruct struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
}

func getSamplePod() *podStruct {
	return &podStruct{
		Kind:       "Pod",
		APIVersion: "v1",
		Metadata: metadataStruct{
			Name:              "test-pod",
			Namespace:         "default",
			CreationTimestamp: "2025-10-03T16:39:00Z",
			Labels: map[string]string{
				"app": "test",
			},
		},
	}
}
