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

	script := `
		local spew = require("spew")
		local result = spew.sdump({name="John", age=30})

		-- Check that result contains expected values
		if not string.find(result, "name") then
			error("Expected result to contain 'name'")
		end

		if not string.find(result, "John") then
			error("Expected result to contain 'John'")
		end

		if not string.find(result, "age") then
			error("Expected result to contain 'age'")
		end

		if not string.find(result, "30") then
			error("Expected result to contain '30'")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
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

	script := `
		local spew = require("spew")
		local result = spew.sdump({1, 2, 3, 4, 5})

		-- Check that result is non-empty
		if #result == 0 then
			error("Expected non-empty result")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
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

	script := `
		local spew = require("spew")
		local data = {
			person = {
				name = "Alice",
				address = {
					city = "NYC",
					zip = "10001"
				}
			},
			items = {1, 2, 3}
		}

		local result = spew.sdump(data)

		-- Check nested values are present
		if not string.find(result, "Alice") then
			error("Expected result to contain 'Alice'")
		end

		if not string.find(result, "NYC") then
			error("Expected result to contain 'NYC'")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
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

	script := `
		local spew = require("spew")

		-- Test string
		local str = spew.sdump("hello")
		if #str == 0 then
			error("Expected non-empty result for string")
		end

		-- Test number
		local num = spew.sdump(42)
		if #num == 0 then
			error("Expected non-empty result for number")
		end

		-- Test boolean
		local bool = spew.sdump(true)
		if #bool == 0 then
			error("Expected non-empty result for boolean")
		end

		-- Test nil
		local nil_val = spew.sdump(nil)
		if #nil_val == 0 then
			error("Expected non-empty result for nil")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_EmptyTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	script := `
		local spew = require("spew")
		local result = spew.sdump({})

		if #result == 0 then
			error("Expected non-empty result for empty table")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_MixedTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	script := `
		local spew = require("spew")
		local data = {
			name = "test",
			count = 42,
			active = true,
			tags = {"a", "b", "c"},
			nested = {
				key = "value"
			}
		}

		local result = spew.sdump(data)

		-- Verify various types are present
		if not string.find(result, "test") then
			error("Expected result to contain string value")
		end

		if not string.find(result, "42") then
			error("Expected result to contain number value")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestDump_OutputsToStdout(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	// dump() prints to stdout and returns nothing
	script := `
		local spew = require("spew")
		spew.dump({name="Bob", age=25})
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}

	// Just verify it doesn't crash
	// Actual stdout capture would be more complex
}

func TestSdump_ComplexNesting(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	script := `
		local spew = require("spew")
		local data = {
			level1 = {
				level2 = {
					level3 = {
						level4 = {
							deep_value = "found"
						}
					}
				}
			}
		}

		local result = spew.sdump(data)

		if not string.find(result, "found") then
			error("Expected result to contain deeply nested value")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
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

	script := `
		local spew = require("spew")
		local users = {
			{name="Alice", age=30},
			{name="Bob", age=25},
			{name="Charlie", age=35}
		}

		local result = spew.sdump(users)

		-- Check all names are present
		if not string.find(result, "Alice") then
			error("Expected result to contain 'Alice'")
		end

		if not string.find(result, "Bob") then
			error("Expected result to contain 'Bob'")
		end

		if not string.find(result, "Charlie") then
			error("Expected result to contain 'Charlie'")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestSdump_NumberKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	script := `
		local spew = require("spew")
		-- Non-consecutive numeric keys should be treated as map
		local data = {
			[1] = "first",
			[5] = "fifth",
			[10] = "tenth"
		}

		local result = spew.sdump(data)

		if not string.find(result, "first") then
			error("Expected result to contain 'first'")
		end

		if not string.find(result, "fifth") then
			error("Expected result to contain 'fifth'")
		end

		return result
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestModuleLoading(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("spew", Loader)

	script := `
		local spew = require("spew")

		-- Verify module has expected functions
		if type(spew.dump) ~= "function" then
			error("Expected spew.dump to be a function")
		end

		if type(spew.sdump) ~= "function" then
			error("Expected spew.sdump to be a function")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}
