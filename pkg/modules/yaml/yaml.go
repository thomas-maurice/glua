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

package yaml

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// Loader: creates and returns the yaml module for Lua.
// This function should be registered with L.PreloadModule("yaml", yaml.Loader)
//
// @luamodule yaml
//
// Example usage in Lua:
//
//	local yaml = require("yaml")
//	local tbl = yaml.parse('name: John\nage: 30')
//	local str = yaml.stringify({name="Jane", age=25})
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"parse":     parse,
	"stringify": stringify,
}

// parse: parses a YAML string and returns a Lua table.
// Returns nil and error message on failure.
//
// @luafunc parse
// @luaparam yamlstr string The YAML string to parse
// @luareturn table The parsed YAML as a Lua table, or nil on error
// @luareturn string|nil Error message if parsing failed
//
// Example:
//
//	local tbl, err = yaml.parse('name: John\nage: 30')
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(tbl.name)  -- prints "John"
//	end
func parse(L *lua.LState) int {
	yamlStr := L.CheckString(1)

	// Parse YAML into a generic map
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse YAML: %v", err)))
		return 2
	}

	// Convert to Lua value
	luaValue := goToLua(L, data)
	L.Push(luaValue)
	L.Push(lua.LNil)
	return 2
}

// stringify: converts a Lua table to a YAML string.
// Returns nil and error message on failure.
//
// @luafunc stringify
// @luaparam tbl table The Lua table to convert to YAML
// @luareturn string The YAML string, or nil on error
// @luareturn string|nil Error message if conversion failed
//
// Example:
//
//	local str, err = yaml.stringify({name="Jane", age=25})
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(str)  -- prints 'age: 25\nname: Jane\n'
//	end
func stringify(L *lua.LState) int {
	luaValue := L.CheckAny(1)

	// Convert Lua value to Go
	goValue := luaToGo(L, luaValue)

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(goValue)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stringify to YAML: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(yamlBytes)))
	L.Push(lua.LNil)
	return 2
}

// goToLua: converts a Go value (from yaml.Unmarshal) to a Lua value
func goToLua(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}

	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		// Convert array to Lua table (1-indexed)
		tbl := L.NewTable()
		for i, item := range v {
			tbl.RawSetInt(i+1, goToLua(L, item))
		}
		return tbl
	case map[string]interface{}:
		// Convert object to Lua table
		tbl := L.NewTable()
		for key, val := range v {
			tbl.RawSetString(key, goToLua(L, val))
		}
		return tbl
	case map[interface{}]interface{}:
		// YAML can have non-string keys, convert to string keys
		tbl := L.NewTable()
		for key, val := range v {
			keyStr := fmt.Sprintf("%v", key)
			tbl.RawSetString(keyStr, goToLua(L, val))
		}
		return tbl
	default:
		// Fallback: convert to string
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

// luaToGo: converts a Lua value to a Go value (for yaml.Marshal)
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
		// Determine if table is an array or object
		maxN := 0
		isArray := true
		v.ForEach(func(key lua.LValue, val lua.LValue) {
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

		// If it's an array (consecutive integer keys starting from 1)
		if isArray && maxN > 0 {
			arr := make([]interface{}, maxN)
			for i := 1; i <= maxN; i++ {
				arr[i-1] = luaToGo(L, v.RawGetInt(i))
			}
			return arr
		}

		// Otherwise, treat as object
		obj := make(map[string]interface{})
		v.ForEach(func(key lua.LValue, val lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				obj[string(keyStr)] = luaToGo(L, val)
			} else {
				// Convert non-string keys to strings
				obj[fmt.Sprintf("%v", key)] = luaToGo(L, val)
			}
		})
		return obj
	default:
		// Fallback: convert to string
		return fmt.Sprintf("%v", v)
	}
}
