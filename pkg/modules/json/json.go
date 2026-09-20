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

// Package json provides JSON serialisation and deserialisation utilities for
// Lua scripts.
package json

import (
	"encoding/json"
	"fmt"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// parse: parses a JSON string and returns the corresponding Lua value; raises
// on malformed JSON. Uses *lua.LState escape hatch to push structured Lua values.
func parse(L *lua.LState, jsonStr string) (lua.LValue, error) {
	var data interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return lua.LNil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return goToLua(L, data), nil
}

// stringify: converts a Lua value to a JSON string; raises on marshal failure
// or on a cyclic/pathologically deep table (see luaToGo).
func stringify(L *lua.LState, value lua.LValue) (string, error) {
	goValue, err := luaToGo(L, value, glua.NewTableGuard())
	if err != nil {
		return "", fmt.Errorf("failed to convert Lua value: %w", err)
	}
	b, err := json.Marshal(goValue)
	if err != nil {
		return "", fmt.Errorf("failed to stringify to JSON: %w", err)
	}
	return string(b), nil
}

// goToLua: converts a Go value (from json.Unmarshal) to a Lua value.
func goToLua(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		tbl := L.NewTable()
		for i, item := range v {
			tbl.RawSetInt(i+1, goToLua(L, item))
		}
		return tbl
	case map[string]interface{}:
		tbl := L.NewTable()
		for key, val := range v {
			tbl.RawSetString(key, goToLua(L, val))
		}
		return tbl
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

// luaToGo: converts a Lua value to a Go value for json.Marshal.
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
		if err := guard.Enter(v); err != nil {
			return nil, err
		}
		defer guard.Leave(v)

		maxN := 0
		isArray := true
		v.ForEach(func(key, _ lua.LValue) {
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
		if isArray && maxN > 0 {
			arr := make([]interface{}, maxN)
			for i := 1; i <= maxN; i++ {
				item, err := luaToGo(L, v.RawGetInt(i), guard)
				if err != nil {
					return nil, err
				}
				arr[i-1] = item
			}
			return arr, nil
		}
		obj := make(map[string]interface{})
		var forEachErr error
		v.ForEach(func(key lua.LValue, val lua.LValue) {
			if forEachErr != nil {
				return
			}
			item, err := luaToGo(L, val, guard)
			if err != nil {
				forEachErr = err
				return
			}
			if keyStr, ok := key.(lua.LString); ok {
				obj[string(keyStr)] = item
			} else {
				obj[fmt.Sprintf("%v", key)] = item
			}
		})
		if forEachErr != nil {
			return nil, forEachErr
		}
		return obj, nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("json", "JSON serialisation and deserialisation utilities")
	m.Fn("parse", parse, "parses a JSON string into a Lua value, raises on invalid JSON",
		luareg.Args("jsonstr"),
		luareg.ArgDoc("jsonstr", "a JSON-encoded string"),
		luareg.ReturnDoc(0, "value", "the decoded value: table, string, number, boolean or nil"))
	m.Fn("stringify", stringify, "converts a Lua value to a JSON string, raises on error",
		luareg.Args("value"),
		luareg.ArgDoc("value", "the Lua value to encode; tables become JSON objects or arrays"),
		luareg.ReturnDoc(0, "jsonstr", "the JSON-encoded string"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("json", json.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
