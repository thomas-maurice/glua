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
