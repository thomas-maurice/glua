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

// Package yaml provides YAML serialisation and deserialisation utilities for
// Lua scripts.
package yaml

import (
	"fmt"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

// parse: parses a YAML string into a Lua value; raises on malformed YAML.
func parse(L *lua.LState, yamlStr string) (lua.LValue, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlStr), &data); err != nil {
		return lua.LNil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return goToLua(L, data), nil
}

// stringify: converts a Lua value to a YAML string; raises on marshal failure.
func stringify(L *lua.LState, value lua.LValue) (string, error) {
	goValue := luaToGo(L, value)
	b, err := yaml.Marshal(goValue)
	if err != nil {
		return "", fmt.Errorf("failed to stringify to YAML: %w", err)
	}
	return string(b), nil
}

// goToLua: converts a Go value (from yaml.Unmarshal) to a Lua value.
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
	case map[interface{}]interface{}:
		// YAML can produce non-string keys; convert to string.
		tbl := L.NewTable()
		for key, val := range v {
			tbl.RawSetString(fmt.Sprintf("%v", key), goToLua(L, val))
		}
		return tbl
	default:
		return lua.LString(fmt.Sprintf("%v", v))
	}
}

// luaToGo: converts a Lua value to a Go value for yaml.Marshal.
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
				arr[i-1] = luaToGo(L, v.RawGetInt(i))
			}
			return arr
		}
		obj := make(map[string]interface{})
		v.ForEach(func(key lua.LValue, val lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				obj[string(keyStr)] = luaToGo(L, val)
			} else {
				obj[fmt.Sprintf("%v", key)] = luaToGo(L, val)
			}
		})
		return obj
	default:
		return fmt.Sprintf("%v", v)
	}
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("yaml", "YAML serialisation and deserialisation utilities")
	m.Fn("parse", parse, "parses a YAML string into a Lua value, raises on invalid YAML",
		luareg.Args("yamlstr"),
		luareg.ArgDoc("yamlstr", "a YAML-encoded string"),
		luareg.ReturnDoc(0, "value", "the decoded value: table, string, number, boolean or nil"))
	m.Fn("stringify", stringify, "converts a Lua value to a YAML string, raises on error",
		luareg.Args("value"),
		luareg.ArgDoc("value", "the Lua value to encode; tables become YAML mappings or sequences"),
		luareg.ReturnDoc(0, "yamlstr", "the YAML-encoded string"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("yaml", yaml.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
