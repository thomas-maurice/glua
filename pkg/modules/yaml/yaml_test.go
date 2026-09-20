// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package yaml

import (
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestParse_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local tbl = yaml.parse("name: John\nage: 30")
		assert(tbl.name == "John", "Expected name to be John")
		assert(tbl.age == 30, "Expected age to be 30")
	`
	require.NoError(t, L.DoString(code))
}

func TestParse_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local tbl = yaml.parse("- apple\n- banana\n- cherry")
		assert(tbl[1] == "apple", "Expected first item to be apple")
		assert(tbl[3] == "cherry", "Expected third item to be cherry")
	`
	require.NoError(t, L.DoString(code))
}

// TestParse_InvalidYAML: parse raises on malformed YAML.
func TestParse_InvalidYAML(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local ok, err = pcall(yaml.parse, "invalid: [yaml")
		assert(not ok, "Expected error for invalid YAML")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}

func TestStringify_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local str = yaml.stringify({name="Jane", age=25})
		assert(type(str) == "string", "Expected string result")
		assert(string.find(str, "name:") ~= nil, "Expected yaml to contain 'name:'")
		assert(string.find(str, "Jane") ~= nil, "Expected yaml to contain 'Jane'")
	`
	require.NoError(t, L.DoString(code))
}

func TestRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")
		local original = {name="Test", value=42, items={"a", "b", "c"}}
		local str = yaml.stringify(original)
		local parsed = yaml.parse(str)
		assert(parsed.name == "Test", "Expected name to match")
		assert(parsed.value == 42, "Expected value to match")
		assert(parsed.items[1] == "a", "Expected first item to match")
	`
	require.NoError(t, L.DoString(code))
}

// TestStringify_CyclicTable_ViaLuaPcall: a self-referential table raises a
// catchable error instead of crashing the process, and the Lua state keeps
// working afterwards.
//
// Regression test: yaml.stringify used to hand-roll its own unguarded
// recursive Lua->Go conversion (luaToGo), so handing it a cyclic table
// (`t.self = t`) recursed forever and killed the whole host process with an
// unrecoverable "fatal error: stack overflow" -- not something pcall could
// ever catch. It now shares pkg/glua's TableGuard, so the same input raises
// an ordinary, catchable Lua error.
func TestStringify_CyclicTable_ViaLuaPcall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("yaml", Loader)

	code := `
		local yaml = require("yaml")

		local t = {}
		t.self = t

		local ok, err = pcall(yaml.stringify, t)
		assert(not ok, "Expected error for a self-referential table")
		assert(string.find(err, "cycle") ~= nil, "Expected error to mention 'cycle', got: " .. err)

		-- Prove the Lua state is still usable after the caught error.
		local result = yaml.stringify({hello = "world"})
		assert(string.find(result, "hello") ~= nil, "Expected normal stringify to still work, got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}
