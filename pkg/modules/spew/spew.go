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
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// dump: prints a Lua value to stdout as colored indented JSON.
// Uses *lua.LState escape hatch to accept any Lua value type.
func dump(L *lua.LState, value lua.LValue) {
	goValue := luaToGo(L, value)
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		return
	}

	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetColors(jsoncolor.DefaultColors())

	var v interface{}
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
		return
	}
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding colored JSON: %v\n", err)
	}
}

// sdump: returns a JSON string representation of a Lua value with indentation.
func sdump(L *lua.LState, value lua.LValue) string {
	goValue := luaToGo(L, value)
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(jsonBytes)
}

// luaToGo: converts a Lua value to a Go value for JSON marshalling.
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

// convertLuaTable: converts a Lua table to a Go slice or map.
func convertLuaTable(L *lua.LState, table *lua.LTable) interface{} {
	maxN, isArray, hasElements := analyzeTableStructure(table)
	if isArray && maxN > 0 && hasElements {
		return convertTableToArray(L, table, maxN)
	}
	return convertTableToMap(L, table)
}

// analyzeTableStructure: determines if a Lua table is an array or map.
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

// convertTableToArray: converts an array-like Lua table to a Go slice.
func convertTableToArray(L *lua.LState, table *lua.LTable, maxN int) []interface{} {
	arr := make([]interface{}, maxN)
	for i := 1; i <= maxN; i++ {
		arr[i-1] = luaToGo(L, table.RawGetInt(i))
	}
	return arr
}

// convertTableToMap: converts a map-like Lua table to a Go map.
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

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("spew", "debug dump utilities for Lua values")
	m.Fn("dump", dump, "prints a Lua value to stdout as colored indented JSON",
		luareg.Args("value"))
	m.Fn("sdump", sdump, "returns a JSON string representation of a Lua value",
		luareg.Args("value"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("spew", spew.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
