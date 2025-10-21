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
	"encoding/json"
	"fmt"
	"os"

	"github.com/neilotoole/jsoncolor"
	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the spew module for Lua.
// This function should be registered with L.PreloadModule("spew", spew.Loader)
//
// @luamodule spew
//
// Example usage in Lua:
//
//	local spew = require("spew")
//	spew.dump({name="John", items={1,2,3}})
//	local str = spew.sdump({key="value"})
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"dump":  dump,
	"sdump": sdump,
}

// dump: prints the contents of a Lua value to stdout with colored JSON formatting.
// This is useful for debugging and inspecting complex table structures.
//
// @luafunc dump
// @luaparam value any The Lua value to dump (table, string, number, etc.)
//
// Example:
//
//	local spew = require("spew")
//	spew.dump({name="John", age=30, items={1,2,3}})
//	-- Prints colored JSON representation to stdout
func dump(L *lua.LState) int {
	value := L.CheckAny(1)

	// Convert Lua value to Go
	goValue := luaToGo(L, value)

	// Marshal to JSON with indentation
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		return 0
	}

	// Create encoder with default colors for terminal output
	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetColors(jsoncolor.DefaultColors())

	// Decode and re-encode with colors
	var v interface{}
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
		return 0
	}

	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding colored JSON: %v\n", err)
		return 0
	}

	return 0
}

// sdump: returns a JSON string representation of a Lua value with indentation.
// Unlike dump, this returns the string instead of printing to stdout (no colors).
//
// @luafunc sdump
// @luaparam value any The Lua value to dump (table, string, number, etc.)
// @luareturn string A JSON string representation of the value
//
// Example:
//
//	local spew = require("spew")
//	local str = spew.sdump({name="John", age=30})
//	print(str)  -- Prints the JSON representation
func sdump(L *lua.LState) int {
	value := L.CheckAny(1)

	// Convert Lua value to Go
	goValue := luaToGo(L, value)

	// Marshal to JSON with indentation
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("Error: %v", err)))
		return 1
	}

	L.Push(lua.LString(string(jsonBytes)))
	return 1
}

// luaToGo: converts a Lua value to a Go value for spew dumping
func luaToGo(L *lua.LState, value lua.LValue) interface{} {
	switch v := value.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		return convertLuaTable(L, v)
	case *lua.LFunction:
		return "<function>"
	case *lua.LUserData:
		return fmt.Sprintf("<userdata: %v>", v.Value)
	default:
		return fmt.Sprintf("<%v>", v.Type().String())
	}
}

// convertLuaTable: converts a Lua table to either a Go slice or map based on key structure
func convertLuaTable(L *lua.LState, table *lua.LTable) interface{} {
	maxN, isArray, hasElements := analyzeTableStructure(table)

	if isArray && maxN > 0 && hasElements {
		return convertTableToArray(L, table, maxN)
	}

	return convertTableToMap(L, table)
}

// analyzeTableStructure: determines if a Lua table should be treated as an array or map
func analyzeTableStructure(table *lua.LTable) (maxN int, isArray bool, hasElements bool) {
	isArray = true
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		hasElements = true
		if keyNum, ok := key.(lua.LNumber); ok {
			if n := int(keyNum); n > 0 && float64(n) == float64(keyNum) {
				if n > maxN {
					maxN = n
				}
			} else {
				isArray = false
			}
		} else {
			isArray = false
		}
	})
	return
}

// convertTableToArray: converts a Lua table with numeric indices to a Go slice
func convertTableToArray(L *lua.LState, table *lua.LTable, maxN int) []interface{} {
	arr := make([]interface{}, maxN)
	for i := 1; i <= maxN; i++ {
		arr[i-1] = luaToGo(L, table.RawGetInt(i))
	}
	return arr
}

// convertTableToMap: converts a Lua table with string keys to a Go map
func convertTableToMap(L *lua.LState, table *lua.LTable) map[string]interface{} {
	obj := make(map[string]interface{})
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		var keyStr string
		if ks, ok := key.(lua.LString); ok {
			keyStr = string(ks)
		} else {
			keyStr = fmt.Sprintf("%v", luaToGo(L, key))
		}
		obj[keyStr] = luaToGo(L, val)
	})
	return obj
}
