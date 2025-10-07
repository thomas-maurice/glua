package glua

import (
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestPrimitiveConversions: tests conversion of basic types in Lua context
func TestPrimitiveConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	tests := []struct {
		name     string
		input    interface{}
		luaCheck string // Lua code to verify the value
		wantErr  bool
	}{
		{
			name:     "string conversion",
			input:    "hello world",
			luaCheck: `assert(type(value) == "string" and value == "hello world", "string mismatch")`,
		},
		{
			name:     "int conversion",
			input:    42,
			luaCheck: `assert(type(value) == "number" and value == 42, "int mismatch")`,
		},
		{
			name:     "float conversion",
			input:    3.14159,
			luaCheck: `assert(type(value) == "number" and math.abs(value - 3.14159) < 0.00001, "float mismatch")`,
		},
		{
			name:     "bool true",
			input:    true,
			luaCheck: `assert(type(value) == "boolean" and value == true, "bool true mismatch")`,
		},
		{
			name:     "bool false",
			input:    false,
			luaCheck: `assert(type(value) == "boolean" and value == false, "bool false mismatch")`,
		},
		{
			name:     "nil value",
			input:    nil,
			luaCheck: `assert(value == nil, "nil mismatch")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv, err := translator.ToLua(L, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToLua() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Set global and verify in Lua
			L.SetGlobal("value", lv)
			if err := L.DoString(tt.luaCheck); err != nil {
				t.Errorf("Lua verification failed: %v", err)
			}
		})
	}
}

// TestSliceConversions: tests array/slice conversions in Lua context
func TestSliceConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	tests := []struct {
		name     string
		input    interface{}
		luaCheck string
	}{
		{
			name:  "string slice",
			input: []string{"a", "b", "c"},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(#value == 3, "wrong length")
				assert(value[1] == "a" and value[2] == "b" and value[3] == "c", "values mismatch")
			`,
		},
		{
			name:  "int slice",
			input: []int{1, 2, 3, 4, 5},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(#value == 5, "wrong length")
				for i = 1, 5 do
					assert(value[i] == i, "value at index " .. i .. " is wrong")
				end
			`,
		},
		{
			name:  "empty slice",
			input: []string{},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(#value == 0, "should be empty")
			`,
		},
		{
			name:  "nested slices",
			input: [][]int{{1, 2}, {3, 4}, {5, 6}},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(#value == 3, "wrong outer length")
				assert(value[1][1] == 1 and value[1][2] == 2, "first nested array wrong")
				assert(value[2][1] == 3 and value[2][2] == 4, "second nested array wrong")
				assert(value[3][1] == 5 and value[3][2] == 6, "third nested array wrong")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv, err := translator.ToLua(L, tt.input)
			if err != nil {
				t.Fatalf("ToLua() error = %v", err)
			}

			L.SetGlobal("value", lv)
			if err := L.DoString(tt.luaCheck); err != nil {
				t.Errorf("Lua verification failed: %v", err)
			}
		})
	}
}

// TestMapConversions: tests map conversions in Lua context
func TestMapConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	tests := []struct {
		name     string
		input    interface{}
		luaCheck string
	}{
		{
			name:  "string map",
			input: map[string]string{"key1": "value1", "key2": "value2"},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(value.key1 == "value1", "key1 mismatch")
				assert(value.key2 == "value2", "key2 mismatch")
			`,
		},
		{
			name:  "int values map",
			input: map[string]int{"count": 42, "total": 100},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(value.count == 42, "count mismatch")
				assert(value.total == 100, "total mismatch")
			`,
		},
		{
			name:  "empty map",
			input: map[string]string{},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				local count = 0
				for k, v in pairs(value) do
					count = count + 1
				end
				assert(count == 0, "should be empty")
			`,
		},
		{
			name: "nested maps",
			input: map[string]map[string]int{
				"group1": {"a": 1, "b": 2},
				"group2": {"c": 3, "d": 4},
			},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(value.group1.a == 1, "group1.a mismatch")
				assert(value.group1.b == 2, "group1.b mismatch")
				assert(value.group2.c == 3, "group2.c mismatch")
				assert(value.group2.d == 4, "group2.d mismatch")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv, err := translator.ToLua(L, tt.input)
			if err != nil {
				t.Fatalf("ToLua() error = %v", err)
			}

			L.SetGlobal("value", lv)
			if err := L.DoString(tt.luaCheck); err != nil {
				t.Errorf("Lua verification failed: %v", err)
			}
		})
	}
}

// TestStructConversions: tests struct conversions in Lua context
func TestStructConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	type Company struct {
		Name      string   `json:"name"`
		Employees []Person `json:"employees"`
	}

	tests := []struct {
		name     string
		input    interface{}
		luaCheck string
	}{
		{
			name:  "simple struct",
			input: Person{Name: "Alice", Age: 30},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(value.name == "Alice", "name mismatch")
				assert(value.age == 30, "age mismatch")
			`,
		},
		{
			name: "nested struct",
			input: Company{
				Name: "ACME Corp",
				Employees: []Person{
					{Name: "Bob", Age: 25},
					{Name: "Charlie", Age: 35},
				},
			},
			luaCheck: `
				assert(type(value) == "table", "not a table")
				assert(value.name == "ACME Corp", "company name mismatch")
				assert(#value.employees == 2, "employees count mismatch")
				assert(value.employees[1].name == "Bob", "first employee name mismatch")
				assert(value.employees[1].age == 25, "first employee age mismatch")
				assert(value.employees[2].name == "Charlie", "second employee name mismatch")
				assert(value.employees[2].age == 35, "second employee age mismatch")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lv, err := translator.ToLua(L, tt.input)
			if err != nil {
				t.Fatalf("ToLua() error = %v", err)
			}

			L.SetGlobal("value", lv)
			if err := L.DoString(tt.luaCheck); err != nil {
				t.Errorf("Lua verification failed: %v", err)
			}
		})
	}
}

// TestKubernetesTypeConversions: tests Kubernetes-specific type conversions
func TestKubernetesTypeConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	// Create a sample Pod
	timestamp := metav1.Time{Time: time.Date(2025, 10, 3, 16, 39, 0, 0, time.UTC)}
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-pod",
			Namespace:         "default",
			CreationTimestamp: timestamp,
			Labels: map[string]string{
				"app": "test",
				"env": "prod",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.21",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
					Env: []corev1.EnvVar{
						{Name: "ENV_VAR_1", Value: "value1"},
						{Name: "ENV_VAR_2", Value: "value2"},
					},
				},
			},
		},
	}

	t.Run("Kubernetes Pod structure", func(t *testing.T) {
		lv, err := translator.ToLua(L, pod)
		if err != nil {
			t.Fatalf("ToLua() error = %v", err)
		}

		L.SetGlobal("pod", lv)

		luaCheck := `
			assert(type(pod) == "table", "pod is not a table")
			assert(pod.kind == "Pod", "kind mismatch")
			assert(pod.apiVersion == "v1", "apiVersion mismatch")

			-- Check metadata
			assert(pod.metadata.name == "test-pod", "name mismatch")
			assert(pod.metadata.namespace == "default", "namespace mismatch")
			assert(pod.metadata.creationTimestamp == "2025-10-03T16:39:00Z", "timestamp mismatch")

			-- Check labels
			assert(pod.metadata.labels.app == "test", "app label mismatch")
			assert(pod.metadata.labels.env == "prod", "env label mismatch")

			-- Check container
			assert(#pod.spec.containers == 1, "container count mismatch")
			assert(pod.spec.containers[1].name == "nginx", "container name mismatch")
			assert(pod.spec.containers[1].image == "nginx:1.21", "container image mismatch")

			-- Check resources (these are preserved as strings)
			assert(pod.spec.containers[1].resources.limits.cpu == "500m", "cpu limit mismatch")
			assert(pod.spec.containers[1].resources.limits.memory == "512Mi", "memory limit mismatch")
			assert(pod.spec.containers[1].resources.requests.cpu == "100m", "cpu request mismatch")
			assert(pod.spec.containers[1].resources.requests.memory == "128Mi", "memory request mismatch")

			-- Check environment variables
			assert(#pod.spec.containers[1].env == 2, "env var count mismatch")
			assert(pod.spec.containers[1].env[1].name == "ENV_VAR_1", "env var 1 name mismatch")
			assert(pod.spec.containers[1].env[1].value == "value1", "env var 1 value mismatch")
			assert(pod.spec.containers[1].env[2].name == "ENV_VAR_2", "env var 2 name mismatch")
			assert(pod.spec.containers[1].env[2].value == "value2", "env var 2 value mismatch")
		`

		if err := L.DoString(luaCheck); err != nil {
			t.Errorf("Lua verification failed: %v", err)
		}
	})
}

// TestLuaToGoConversions: tests Lua to Go conversions
func TestLuaToGoConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Lua table to Go struct", func(t *testing.T) {
		type TestStruct struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Email string `json:"email"`
		}

		// Create Lua table
		luaCode := `
			return {
				name = "John Doe",
				age = 42,
				email = "john@example.com"
			}
		`

		if err := L.DoString(luaCode); err != nil {
			t.Fatalf("Lua code failed: %v", err)
		}

		lv := L.Get(-1)
		L.Pop(1)

		var result TestStruct
		if err := translator.FromLua(L, lv, &result); err != nil {
			t.Fatalf("FromLua() error = %v", err)
		}

		if result.Name != "John Doe" {
			t.Errorf("Name = %v, want John Doe", result.Name)
		}
		if result.Age != 42 {
			t.Errorf("Age = %v, want 42", result.Age)
		}
		if result.Email != "john@example.com" {
			t.Errorf("Email = %v, want john@example.com", result.Email)
		}
	})

	t.Run("Lua table with nested arrays to Go", func(t *testing.T) {
		type Item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		type Container struct {
			Title string `json:"title"`
			Items []Item `json:"items"`
		}

		luaCode := `
			return {
				title = "My Container",
				items = {
					{id = 1, name = "Item 1"},
					{id = 2, name = "Item 2"},
					{id = 3, name = "Item 3"}
				}
			}
		`

		if err := L.DoString(luaCode); err != nil {
			t.Fatalf("Lua code failed: %v", err)
		}

		lv := L.Get(-1)
		L.Pop(1)

		var result Container
		if err := translator.FromLua(L, lv, &result); err != nil {
			t.Fatalf("FromLua() error = %v", err)
		}

		if result.Title != "My Container" {
			t.Errorf("Title = %v, want My Container", result.Title)
		}
		if len(result.Items) != 3 {
			t.Fatalf("len(Items) = %v, want 3", len(result.Items))
		}
		if result.Items[0].ID != 1 || result.Items[0].Name != "Item 1" {
			t.Errorf("Items[0] = %+v, want {ID:1, Name:Item 1}", result.Items[0])
		}
		if result.Items[1].ID != 2 || result.Items[1].Name != "Item 2" {
			t.Errorf("Items[1] = %+v, want {ID:2, Name:Item 2}", result.Items[1])
		}
		if result.Items[2].ID != 3 || result.Items[2].Name != "Item 3" {
			t.Errorf("Items[2] = %+v, want {ID:3, Name:Item 3}", result.Items[2])
		}
	})
}

// TestRoundTripConversions: tests complete round-trip conversions
func TestRoundTripConversions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	type NestedStruct struct {
		Value string `json:"value"`
	}

	type ComplexStruct struct {
		StringField  string         `json:"stringField"`
		IntField     int            `json:"intField"`
		BoolField    bool           `json:"boolField"`
		SliceField   []string       `json:"sliceField"`
		MapField     map[string]int `json:"mapField"`
		NestedStruct *NestedStruct  `json:"nestedStruct"`
	}

	original := ComplexStruct{
		StringField: "test",
		IntField:    42,
		BoolField:   true,
		SliceField:  []string{"a", "b", "c"},
		MapField:    map[string]int{"x": 1, "y": 2},
		NestedStruct: &NestedStruct{
			Value: "nested",
		},
	}

	t.Run("Complete round-trip", func(t *testing.T) {
		// Go -> Lua
		lv, err := translator.ToLua(L, original)
		if err != nil {
			t.Fatalf("ToLua() error = %v", err)
		}

		// Verify in Lua
		L.SetGlobal("value", lv)
		luaCheck := `
			assert(value.stringField == "test", "stringField mismatch")
			assert(value.intField == 42, "intField mismatch")
			assert(value.boolField == true, "boolField mismatch")
			assert(#value.sliceField == 3, "sliceField length mismatch")
			assert(value.mapField.x == 1, "mapField.x mismatch")
			assert(value.nestedStruct.value == "nested", "nested value mismatch")
		`
		if err := L.DoString(luaCheck); err != nil {
			t.Errorf("Lua verification failed: %v", err)
		}

		// Lua -> Go
		var reconstructed ComplexStruct
		if err := translator.FromLua(L, lv, &reconstructed); err != nil {
			t.Fatalf("FromLua() error = %v", err)
		}

		// Verify reconstruction
		if reconstructed.StringField != original.StringField {
			t.Errorf("StringField = %v, want %v", reconstructed.StringField, original.StringField)
		}
		if reconstructed.IntField != original.IntField {
			t.Errorf("IntField = %v, want %v", reconstructed.IntField, original.IntField)
		}
		if reconstructed.BoolField != original.BoolField {
			t.Errorf("BoolField = %v, want %v", reconstructed.BoolField, original.BoolField)
		}
		if len(reconstructed.SliceField) != len(original.SliceField) {
			t.Errorf("len(SliceField) = %v, want %v", len(reconstructed.SliceField), len(original.SliceField))
		}
		if reconstructed.NestedStruct.Value != original.NestedStruct.Value {
			t.Errorf("NestedStruct.Value = %v, want %v", reconstructed.NestedStruct.Value, original.NestedStruct.Value)
		}
	})
}
