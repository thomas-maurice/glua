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

package json

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestParse_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/parse_simple_object.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/parse_array.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_NestedObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/parse_nested_object.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/parse_invalid_json.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/stringify_simple_object.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/stringify_array.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_NestedStructure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/stringify_nested_structure.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_BooleanAndNull(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/stringify_boolean_and_null.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestRoundTrip_ComplexData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/roundtrip_complex_data.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_EmptyObjectAndArray(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/parse_empty.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_Numbers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	if err := L.DoFile("testdata/stringify_numbers.lua"); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}
