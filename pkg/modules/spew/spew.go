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

// Package spew provides debug dump utilities for Lua values, rendering them
// as indented (optionally colored) JSON.
package spew

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/neilotoole/jsoncolor"
	glua "github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// dump: prints a Lua value to stdout as colored indented JSON; raises on a
// cyclic/pathologically deep table (see luaToGo). A JSON marshal/encode
// failure is reported to stderr and swallowed, matching this function's
// pre-existing behaviour for those unrelated failure modes.
func dump(L *lua.LState, value lua.LValue) error {
	goValue, err := luaToGo(L, value, glua.NewTableGuard())
	if err != nil {
		return fmt.Errorf("failed to convert Lua value: %w", err)
	}
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		return nil
	}

	enc := jsoncolor.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetColors(jsoncolor.DefaultColors())

	var v interface{}
	if err := json.Unmarshal(jsonBytes, &v); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling JSON: %v\n", err)
		return nil
	}
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding colored JSON: %v\n", err)
	}
	return nil
}

// sdump: returns a JSON string representation of a Lua value with
// indentation; raises on a cyclic/pathologically deep table (see luaToGo). A
// JSON marshal failure is embedded in the returned string, matching this
// function's pre-existing behaviour for that unrelated failure mode.
func sdump(L *lua.LState, value lua.LValue) (string, error) {
	goValue, err := luaToGo(L, value, glua.NewTableGuard())
	if err != nil {
		return "", fmt.Errorf("failed to convert Lua value: %w", err)
	}
	jsonBytes, err := json.MarshalIndent(goValue, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: %v", err), nil
	}
	return string(jsonBytes), nil
}

// luaToGo: converts a Lua value to a Go value for JSON marshalling.
//
// guard bounds the *lua.LTable recursion below against a cyclic or
// pathologically deep table — see glua.TableGuard's doc comment. Callers at
// the top of a conversion pass glua.NewTableGuard(); recursive calls MUST
// pass the same instance through, not a fresh one, or the guard cannot see
// the whole path.
func luaToGo(L *lua.LState, value lua.LValue, guard *glua.TableGuard) (interface{}, error) {
	switch v := value.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LBool:
		return bool(v), nil
	case lua.LNumber:
		return float64(v), nil
	case lua.LString:
		return string(v), nil
	case *lua.LTable:
		return convertLuaTable(L, v, guard)
	case *lua.LFunction:
		return "<function>", nil
	case *lua.LUserData:
		return fmt.Sprintf("<userdata: %v>", v.Value), nil
	default:
		return fmt.Sprintf("<%v>", v.Type().String()), nil
	}
}

// convertLuaTable: converts a Lua table to a Go slice or map. guard is
// documented on luaToGo.
func convertLuaTable(L *lua.LState, table *lua.LTable, guard *glua.TableGuard) (interface{}, error) {
	if err := guard.Enter(table); err != nil {
		return nil, err
	}
	defer guard.Leave(table)

	maxN, isArray, hasElements := analyzeTableStructure(table)
	if isArray && maxN > 0 && hasElements {
		return convertTableToArray(L, table, maxN, guard)
	}
	return convertTableToMap(L, table, guard)
}

// analyzeTableStructure: determines if a Lua table is an array or map.
func analyzeTableStructure(table *lua.LTable) (maxN int, isArray bool, hasElements bool) {
	isArray = true
	table.ForEach(func(key, _ lua.LValue) {
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

// convertTableToArray: converts an array-like Lua table to a Go slice. guard
// is documented on luaToGo; table has already been entered into it by
// convertLuaTable.
func convertTableToArray(L *lua.LState, table *lua.LTable, maxN int, guard *glua.TableGuard) ([]interface{}, error) {
	arr := make([]interface{}, maxN)
	for i := 1; i <= maxN; i++ {
		item, err := luaToGo(L, table.RawGetInt(i), guard)
		if err != nil {
			return nil, err
		}
		arr[i-1] = item
	}
	return arr, nil
}

// convertTableToMap: converts a map-like Lua table to a Go map. guard is
// documented on luaToGo; table has already been entered into it by
// convertLuaTable.
func convertTableToMap(L *lua.LState, table *lua.LTable, guard *glua.TableGuard) (map[string]interface{}, error) {
	obj := make(map[string]interface{})
	var forEachErr error
	table.ForEach(func(key lua.LValue, val lua.LValue) {
		if forEachErr != nil {
			return
		}
		item, err := luaToGo(L, val, guard)
		if err != nil {
			forEachErr = err
			return
		}
		var keyStr string
		if ks, ok := key.(lua.LString); ok {
			keyStr = string(ks)
		} else {
			keyGo, err := luaToGo(L, key, guard)
			if err != nil {
				forEachErr = err
				return
			}
			keyStr = fmt.Sprintf("%v", keyGo)
		}
		obj[keyStr] = item
	})
	if forEachErr != nil {
		return nil, forEachErr
	}
	return obj, nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("spew", "debug dump utilities for Lua values")
	m.Fn("dump", dump, "prints a Lua value to stdout as colored indented JSON",
		luareg.Args("value"),
		luareg.ArgDoc("value", "the Lua value to dump; tables are walked recursively"))
	m.Fn("sdump", sdump, "returns a JSON string representation of a Lua value",
		luareg.Args("value"),
		luareg.ArgDoc("value", "the Lua value to dump; tables are walked recursively"),
		luareg.ReturnDoc(0, "s", "the indented JSON representation of value"))
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
