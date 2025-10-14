// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package yaml

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestParse_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local tbl, err = yaml.parse("name: John\nage: 30")
		assert(err == nil, "Expected no error")
		assert(tbl.name == "John", "Expected name to be John")
		assert(tbl.age == 30, "Expected age to be 30")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestParse_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local tbl, err = yaml.parse("- apple\n- banana\n- cherry")
		assert(err == nil, "Expected no error")
		assert(tbl[1] == "apple", "Expected first item to be apple")
		assert(tbl[2] == "banana", "Expected second item to be banana")
		assert(tbl[3] == "cherry", "Expected third item to be cherry")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local tbl, err = yaml.parse("invalid: [yaml")
		assert(tbl == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestStringify_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local str, err = yaml.stringify({name="Jane", age=25})
		assert(err == nil, "Expected no error")
		assert(type(str) == "string", "Expected string result")
		assert(string.find(str, "name:") ~= nil, "Expected yaml to contain 'name:'")
		assert(string.find(str, "Jane") ~= nil, "Expected yaml to contain 'Jane'")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local original = {name="Test", value=42, items={"a", "b", "c"}}
		local str, err = yaml.stringify(original)
		assert(err == nil, "Expected no error on stringify")

		local parsed, err2 = yaml.parse(str)
		assert(err2 == nil, "Expected no error on parse")
		assert(parsed.name == "Test", "Expected name to match")
		assert(parsed.value == 42, "Expected value to match")
		assert(parsed.items[1] == "a", "Expected first item to match")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
