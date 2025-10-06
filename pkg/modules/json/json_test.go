package json

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestParse_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local tbl, err = json.parse('{"name":"John","age":30,"active":true}')

		if err then
			error("Parse failed: " .. err)
		end

		if tbl.name ~= "John" then
			error("Expected name to be 'John', got: " .. tostring(tbl.name))
		end

		if tbl.age ~= 30 then
			error("Expected age to be 30, got: " .. tostring(tbl.age))
		end

		if tbl.active ~= true then
			error("Expected active to be true, got: " .. tostring(tbl.active))
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local tbl, err = json.parse('[1,2,3,4,5]')

		if err then
			error("Parse failed: " .. err)
		end

		if #tbl ~= 5 then
			error("Expected array length 5, got: " .. #tbl)
		end

		if tbl[1] ~= 1 or tbl[5] ~= 5 then
			error("Array values incorrect")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_NestedObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local tbl, err = json.parse('{"person":{"name":"Jane","address":{"city":"NYC"}}}')

		if err then
			error("Parse failed: " .. err)
		end

		if tbl.person.name ~= "Jane" then
			error("Expected nested name to be 'Jane', got: " .. tostring(tbl.person.name))
		end

		if tbl.person.address.city ~= "NYC" then
			error("Expected nested city to be 'NYC', got: " .. tostring(tbl.person.address.city))
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local tbl, err = json.parse('{invalid json}')

		if tbl ~= nil then
			error("Expected tbl to be nil for invalid JSON")
		end

		if err == nil then
			error("Expected error for invalid JSON")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_SimpleObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local str, err = json.stringify({name="Alice", age=25})

		if err then
			error("Stringify failed: " .. err)
		end

		-- Parse it back to verify
		local tbl, parseErr = json.parse(str)
		if parseErr then
			error("Failed to parse stringified JSON: " .. parseErr)
		end

		if tbl.name ~= "Alice" then
			error("Expected name to be 'Alice' after round-trip")
		end

		if tbl.age ~= 25 then
			error("Expected age to be 25 after round-trip")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local str, err = json.stringify({1, 2, 3, 4, 5})

		if err then
			error("Stringify failed: " .. err)
		end

		-- Parse it back to verify
		local tbl, parseErr = json.parse(str)
		if parseErr then
			error("Failed to parse stringified JSON: " .. parseErr)
		end

		if #tbl ~= 5 then
			error("Expected array length 5 after round-trip, got: " .. #tbl)
		end

		if tbl[1] ~= 1 or tbl[5] ~= 5 then
			error("Array values incorrect after round-trip")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_NestedStructure(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local data = {
			user = {
				name = "Bob",
				tags = {"admin", "user"}
			},
			count = 42
		}

		local str, err = json.stringify(data)
		if err then
			error("Stringify failed: " .. err)
		end

		-- Parse it back to verify
		local tbl, parseErr = json.parse(str)
		if parseErr then
			error("Failed to parse stringified JSON: " .. parseErr)
		end

		if tbl.user.name ~= "Bob" then
			error("Expected nested name to be 'Bob' after round-trip")
		end

		if tbl.user.tags[1] ~= "admin" then
			error("Expected first tag to be 'admin' after round-trip")
		end

		if tbl.count ~= 42 then
			error("Expected count to be 42 after round-trip")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_BooleanAndNull(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local data = {
			active = true,
			disabled = false,
			value = nil
		}

		local str, err = json.stringify(data)
		if err then
			error("Stringify failed: " .. err)
		end

		-- Parse it back to verify
		local tbl, parseErr = json.parse(str)
		if parseErr then
			error("Failed to parse stringified JSON: " .. parseErr)
		end

		if tbl.active ~= true then
			error("Expected active to be true after round-trip")
		end

		if tbl.disabled ~= false then
			error("Expected disabled to be false after round-trip")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestRoundTrip_ComplexData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")

		-- Original JSON
		local original = '{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}],"active":true,"count":2}'

		-- Parse to table
		local tbl, err1 = json.parse(original)
		if err1 then
			error("Parse failed: " .. err1)
		end

		-- Stringify back to JSON
		local stringified, err2 = json.stringify(tbl)
		if err2 then
			error("Stringify failed: " .. err2)
		end

		-- Parse again to verify
		local final, err3 = json.parse(stringified)
		if err3 then
			error("Second parse failed: " .. err3)
		end

		-- Verify structure
		if final.users[1].name ~= "Alice" then
			error("Round-trip failed: name mismatch")
		end

		if final.users[2].age ~= 25 then
			error("Round-trip failed: age mismatch")
		end

		if final.active ~= true then
			error("Round-trip failed: active mismatch")
		end

		if final.count ~= 2 then
			error("Round-trip failed: count mismatch")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestParse_EmptyObjectAndArray(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")

		-- Empty object
		local obj, err1 = json.parse('{}')
		if err1 then
			error("Parse empty object failed: " .. err1)
		end

		-- Empty array
		local arr, err2 = json.parse('[]')
		if err2 then
			error("Parse empty array failed: " .. err2)
		end

		if #arr ~= 0 then
			error("Expected empty array length to be 0, got: " .. #arr)
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}

func TestStringify_Numbers(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", Loader)

	script := `
		local json = require("json")
		local data = {
			integer = 42,
			negative = -10,
			float = 3.14,
			zero = 0
		}

		local str, err = json.stringify(data)
		if err then
			error("Stringify failed: " .. err)
		end

		-- Parse it back to verify
		local tbl, parseErr = json.parse(str)
		if parseErr then
			error("Failed to parse stringified JSON: " .. parseErr)
		end

		if tbl.integer ~= 42 then
			error("Expected integer to be 42 after round-trip")
		end

		if tbl.negative ~= -10 then
			error("Expected negative to be -10 after round-trip")
		end

		if math.abs(tbl.float - 3.14) > 0.01 then
			error("Expected float to be ~3.14 after round-trip")
		end

		if tbl.zero ~= 0 then
			error("Expected zero to be 0 after round-trip")
		end
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua script failed: %v", err)
	}
}
