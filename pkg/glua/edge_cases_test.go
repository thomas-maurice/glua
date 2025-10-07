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

package glua

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestEdgeCases_NilValues: tests handling of nil values in various contexts
func TestEdgeCases_NilValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Nil pointer", func(t *testing.T) {
		var ptr *string
		result, err := translator.ToLua(L, ptr)
		if err != nil {
			t.Fatalf("ToLua(nil pointer) failed: %v", err)
		}
		if result != lua.LNil {
			t.Errorf("Expected LNil, got %T", result)
		}
	})

	t.Run("Nil map", func(t *testing.T) {
		var m map[string]string
		result, err := translator.ToLua(L, m)
		if err != nil {
			t.Fatalf("ToLua(nil map) failed: %v", err)
		}
		if result != lua.LNil {
			t.Errorf("Expected LNil, got %T", result)
		}
	})

	t.Run("Nil slice", func(t *testing.T) {
		var s []string
		result, err := translator.ToLua(L, s)
		if err != nil {
			t.Fatalf("ToLua(nil slice) failed: %v", err)
		}
		if result != lua.LNil {
			t.Errorf("Expected LNil, got %T", result)
		}
	})

	t.Run("Struct with nil fields", func(t *testing.T) {
		type TestStruct struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
			Meta   *string  `json:"meta"`
		}

		obj := TestStruct{Name: "test"}
		result, err := translator.ToLua(L, obj)
		if err != nil {
			t.Fatalf("ToLua(struct with nils) failed: %v", err)
		}

		if _, ok := result.(*lua.LTable); !ok {
			t.Errorf("Expected LTable, got %T", result)
		}
	})
}

// TestEdgeCases_EmptyCollections: tests handling of empty maps, slices, arrays
func TestEdgeCases_EmptyCollections(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Empty slice", func(t *testing.T) {
		input := []string{}
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(empty slice) failed: %v", err)
		}

		table, ok := result.(*lua.LTable)
		if !ok {
			t.Fatalf("Expected LTable, got %T", result)
		}

		if table.Len() != 0 {
			t.Errorf("Expected empty table, got length %d", table.Len())
		}
	})

	t.Run("Empty map", func(t *testing.T) {
		input := map[string]string{}
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(empty map) failed: %v", err)
		}

		table, ok := result.(*lua.LTable)
		if !ok {
			t.Fatalf("Expected LTable, got %T", result)
		}

		count := 0
		table.ForEach(func(k, v lua.LValue) {
			count++
		})

		if count != 0 {
			t.Errorf("Expected 0 entries, got %d", count)
		}
	})

	t.Run("Empty string", func(t *testing.T) {
		input := ""
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(empty string) failed: %v", err)
		}

		str, ok := result.(lua.LString)
		if !ok {
			t.Fatalf("Expected LString, got %T", result)
		}

		if string(str) != "" {
			t.Errorf("Expected empty string, got %q", str)
		}
	})
}

// TestEdgeCases_SpecialCharacters: tests handling of special characters in strings and keys
func TestEdgeCases_SpecialCharacters(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Unicode strings", func(t *testing.T) {
		input := "Hello 世界 🚀"
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(unicode) failed: %v", err)
		}

		str, ok := result.(lua.LString)
		if !ok {
			t.Fatalf("Expected LString, got %T", result)
		}

		if string(str) != input {
			t.Errorf("Unicode not preserved: got %q, want %q", str, input)
		}
	})

	t.Run("Special characters in map keys", func(t *testing.T) {
		input := map[string]string{
			"key.with.dots":    "value1",
			"key/with/slashes": "value2",
			"key~with~tilde":   "value3",
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(special keys) failed: %v", err)
		}

		var output map[string]string
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if output["key.with.dots"] != "value1" {
			t.Errorf("Dots not preserved in key")
		}
		if output["key/with/slashes"] != "value2" {
			t.Errorf("Slashes not preserved in key")
		}
		if output["key~with~tilde"] != "value3" {
			t.Errorf("Tildes not preserved in key")
		}
	})

	t.Run("Newlines and tabs", func(t *testing.T) {
		input := "Line1\nLine2\tTabbed"
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(newlines/tabs) failed: %v", err)
		}

		var output string
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if output != input {
			t.Errorf("Newlines/tabs not preserved: got %q, want %q", output, input)
		}
	})
}

// TestEdgeCases_NumericBoundaries: tests handling of numeric edge cases
func TestEdgeCases_NumericBoundaries(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Zero values", func(t *testing.T) {
		type TestStruct struct {
			Int    int     `json:"int"`
			Float  float64 `json:"float"`
			Bool   bool    `json:"bool"`
			String string  `json:"string"`
		}

		input := TestStruct{}
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(zero values) failed: %v", err)
		}

		var output TestStruct
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if output.Int != 0 {
			t.Errorf("Int: got %d, want 0", output.Int)
		}
		if output.Float != 0.0 {
			t.Errorf("Float: got %f, want 0.0", output.Float)
		}
		if output.Bool != false {
			t.Errorf("Bool: got %v, want false", output.Bool)
		}
		if output.String != "" {
			t.Errorf("String: got %q, want empty string", output.String)
		}
	})

	t.Run("Large numbers", func(t *testing.T) {
		type TestStruct struct {
			Int64 int64 `json:"int64"`
		}

		input := TestStruct{Int64: 9007199254740991} // Max safe integer in JSON (2^53 - 1)
		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(large number) failed: %v", err)
		}

		var output TestStruct
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if output.Int64 != input.Int64 {
			t.Errorf("Large int64 not preserved: got %d, want %d", output.Int64, input.Int64)
		}
	})

	t.Run("Negative numbers", func(t *testing.T) {
		input := map[string]int{
			"negative": -42,
			"zero":     0,
			"positive": 42,
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(negative numbers) failed: %v", err)
		}

		var output map[string]int
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if output["negative"] != -42 {
			t.Errorf("Negative number not preserved: got %d", output["negative"])
		}
	})
}

// TestEdgeCases_DeeplyNested: tests handling of deeply nested structures
func TestEdgeCases_DeeplyNested(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Nested maps", func(t *testing.T) {
		input := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": map[string]interface{}{
						"level4": map[string]interface{}{
							"value": "deep",
						},
					},
				},
			},
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(nested maps) failed: %v", err)
		}

		var output map[string]interface{}
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		// Navigate to deep value
		l1, ok := output["level1"].(map[string]interface{})
		if !ok {
			t.Fatal("Level 1 not a map")
		}
		l2, ok := l1["level2"].(map[string]interface{})
		if !ok {
			t.Fatal("Level 2 not a map")
		}
		l3, ok := l2["level3"].(map[string]interface{})
		if !ok {
			t.Fatal("Level 3 not a map")
		}
		l4, ok := l3["level4"].(map[string]interface{})
		if !ok {
			t.Fatal("Level 4 not a map")
		}
		value, ok := l4["value"].(string)
		if !ok || value != "deep" {
			t.Errorf("Deep value not preserved: got %v", l4["value"])
		}
	})

	t.Run("Nested slices", func(t *testing.T) {
		input := []interface{}{
			[]interface{}{
				[]interface{}{
					"deep",
				},
			},
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(nested slices) failed: %v", err)
		}

		var output []interface{}
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if len(output) != 1 {
			t.Fatalf("Expected 1 element, got %d", len(output))
		}
	})
}

// TestEdgeCases_InvalidFromLua: tests error handling for invalid Lua to Go conversions
func TestEdgeCases_InvalidFromLua(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Non-pointer output", func(t *testing.T) {
		if err := L.DoString("return {name = 'test'}"); err != nil {
			t.Fatalf("DoString failed: %v", err)
		}
		lv := L.Get(-1)
		L.Pop(1)

		var output struct {
			Name string `json:"name"`
		}

		err := translator.FromLua(L, lv, output) // Not a pointer!
		if err == nil {
			t.Error("Expected error for non-pointer output, got nil")
		}
	})

	t.Run("Nil LValue", func(t *testing.T) {
		var output string
		err := translator.FromLua(L, lua.LNil, &output)
		if err != nil {
			t.Errorf("FromLua(nil) should handle nil gracefully: %v", err)
		}
	})

	t.Run("Type mismatch", func(t *testing.T) {
		if err := L.DoString("return 'not a number'"); err != nil {
			t.Fatalf("DoString failed: %v", err)
		}
		lv := L.Get(-1)
		L.Pop(1)

		var output int
		err := translator.FromLua(L, lv, &output)
		// This might not error due to JSON marshaling flexibility
		// but we're testing the behavior is defined
		if err != nil {
			// Error is acceptable
			t.Logf("Type mismatch error (acceptable): %v", err)
		}
	})
}

// TestEdgeCases_ConcurrentAccess: tests thread safety (Lua states are not thread-safe)
func TestEdgeCases_ConcurrentAccess(t *testing.T) {
	// Note: This test documents that we should NOT share Lua states across goroutines
	// Each goroutine should have its own Lua state

	translator := NewTranslator()

	type TestStruct struct {
		Value string `json:"value"`
	}

	// Create separate Lua states for each goroutine (correct usage)
	t.Run("Separate states (safe)", func(t *testing.T) {
		done := make(chan bool, 2)

		for i := 0; i < 2; i++ {
			go func(id int) {
				defer func() { done <- true }()

				L := lua.NewState()
				defer L.Close()

				input := TestStruct{Value: string(rune('A' + id))}
				result, err := translator.ToLua(L, input)
				if err != nil {
					t.Errorf("Goroutine %d: ToLua failed: %v", id, err)
					return
				}

				var output TestStruct
				err = translator.FromLua(L, result, &output)
				if err != nil {
					t.Errorf("Goroutine %d: FromLua failed: %v", id, err)
					return
				}

				if output.Value != input.Value {
					t.Errorf("Goroutine %d: value mismatch", id)
				}
			}(i)
		}

		<-done
		<-done
	})
}

// TestEdgeCases_LargeData: tests handling of large data structures
func TestEdgeCases_LargeData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	translator := NewTranslator()

	t.Run("Large slice", func(t *testing.T) {
		input := make([]int, 1000)
		for i := range input {
			input[i] = i
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(large slice) failed: %v", err)
		}

		var output []int
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if len(output) != len(input) {
			t.Errorf("Length mismatch: got %d, want %d", len(output), len(input))
		}

		for i := range input {
			if output[i] != input[i] {
				t.Errorf("Value mismatch at index %d: got %d, want %d", i, output[i], input[i])
				break
			}
		}
	})

	t.Run("Large map", func(t *testing.T) {
		input := make(map[string]int)
		for i := 0; i < 100; i++ {
			input[string(rune('A'+i%26))+string(rune('0'+i/26))] = i
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(large map) failed: %v", err)
		}

		var output map[string]int
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if len(output) != len(input) {
			t.Errorf("Size mismatch: got %d, want %d", len(output), len(input))
		}
	})

	t.Run("Long string", func(t *testing.T) {
		input := string(make([]byte, 10000))
		for i := range []byte(input) {
			input = input[:i] + string(rune('A'+i%26)) + input[i+1:]
		}

		result, err := translator.ToLua(L, input)
		if err != nil {
			t.Fatalf("ToLua(long string) failed: %v", err)
		}

		var output string
		err = translator.FromLua(L, result, &output)
		if err != nil {
			t.Fatalf("FromLua failed: %v", err)
		}

		if len(output) != len(input) {
			t.Errorf("Length mismatch: got %d, want %d", len(output), len(input))
		}
	})
}
