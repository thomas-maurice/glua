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
	"path/filepath"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/ directory
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	if err != nil {
		t.Fatalf("Failed to glob testdata: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No Lua test files found in testdata/")
	}

	// Map of test files to their specific validation logic
	validators := map[string]func(*testing.T, *lua.LState){
		"sdump_simple_object.lua": func(t *testing.T, L *lua.LState) {
			result := L.Get(-1)
			if result.Type() != lua.LTString {
				t.Fatalf("Expected string result, got %v", result.Type())
			}
			resultStr := result.String()
			if !strings.Contains(resultStr, "name") {
				t.Errorf("Expected result to contain 'name', got: %s", resultStr)
			}
		},
		"sdump_array.lua": func(t *testing.T, L *lua.LState) {
			result := L.Get(-1).String()
			if len(result) == 0 {
				t.Error("Expected non-empty spew output for array")
			}
		},
		"sdump_nested_structure.lua": func(t *testing.T, L *lua.LState) {
			result := L.Get(-1).String()
			if !strings.Contains(result, "Alice") || !strings.Contains(result, "NYC") {
				t.Errorf("Expected result to contain nested values")
			}
		},
		"sdump_complex_nesting.lua": func(t *testing.T, L *lua.LState) {
			result := L.Get(-1).String()
			if !strings.Contains(result, "found") {
				t.Error("Expected deeply nested value to be present in output")
			}
		},
	}

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("spew", Loader)

			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}

			// Run validator if one exists for this test
			if validator, ok := validators[testName]; ok {
				validator(t, L)
			}
		})
	}
}
