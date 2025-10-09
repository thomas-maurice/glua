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

package spew

import (
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestSdump_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_simple_object.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	// Check the returned value
	result := L.Get(-1)
	if result.Type() != lua.LTString {
		t.Fatalf("Expected string result, got %v", result.Type())
	}

	resultStr := result.String()
	if !strings.Contains(resultStr, "name") {
		t.Errorf("Expected result to contain 'name', got: %s", resultStr)
	}
}

func TestSdump_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_array.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	result := L.Get(-1).String()
	if len(result) == 0 {
		t.Error("Expected non-empty spew output for array")
	}
}

func TestSdump_NestedStructure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_nested_structure.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	result := L.Get(-1).String()
	if !strings.Contains(result, "Alice") || !strings.Contains(result, "NYC") {
		t.Errorf("Expected result to contain nested values")
	}
}

func TestSdump_PrimitiveTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_primitive_types.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_EmptyTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_empty_table.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_MixedTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_mixed_table.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestDump_OutputsToStdout(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/dump_outputs_to_stdout.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	// Just verify it doesn't crash
	// Actual stdout capture would be more complex
}

func TestSdump_ComplexNesting(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_complex_nesting.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	result := L.Get(-1).String()
	if !strings.Contains(result, "found") {
		t.Error("Expected deeply nested value to be present in output")
	}
}

func TestSdump_ArrayOfObjects(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_array_of_objects.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_NumberKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/sdump_number_keys.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestModuleLoading(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	if err := L.DoFile("testdata/module_loading.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}
