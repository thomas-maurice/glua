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

// Package template provides Go text/template rendering utilities for Lua
// scripts, from an inline string or a file, with a Lua table as the data
// context.
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// render: renders a Go text/template with the provided data table; raises on
// parse or execution failure, or on a cyclic/pathologically deep data table
// (see luaValueToGo). Uses *lua.LTable escape hatch for the data arg.
func render(L *lua.LState, tmplStr string, data *lua.LTable) (string, error) {
	goData, err := luaTableToGoMap(L, data, glua.NewTableGuard())
	if err != nil {
		return "", fmt.Errorf("failed to convert template data: %w", err)
	}
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, goData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// renderFile: renders a Go text/template from a file with the provided data
// table; raises on read, parse, or execution failure, or on a
// cyclic/pathologically deep data table (see luaValueToGo).
func renderFile(L *lua.LState, path string, data *lua.LTable) (string, error) {
	tmplBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	goData, err := luaTableToGoMap(L, data, glua.NewTableGuard())
	if err != nil {
		return "", fmt.Errorf("failed to convert template data: %w", err)
	}
	tmpl, err := template.New("tmpl").Parse(string(tmplBytes))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, goData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// luaTableToGoMap: converts a Lua table to a Go map for template rendering.
//
// guard bounds the recursion below against a cyclic or pathologically deep
// table — see glua.TableGuard's doc comment. Callers at the top of a
// conversion pass glua.NewTableGuard(); recursive calls (via luaValueToGo)
// MUST pass the same instance through, not a fresh one, or the guard cannot
// see the whole path.
func luaTableToGoMap(L *lua.LState, tbl *lua.LTable, guard *glua.TableGuard) (map[string]interface{}, error) {
	if err := guard.Enter(tbl); err != nil {
		return nil, err
	}
	defer guard.Leave(tbl)

	result := make(map[string]interface{})
	var forEachErr error
	tbl.ForEach(func(key lua.LValue, val lua.LValue) {
		if forEachErr != nil {
			return
		}
		var keyStr string
		if k, ok := key.(lua.LString); ok {
			keyStr = string(k)
		} else {
			keyStr = fmt.Sprintf("%v", key)
		}
		item, err := luaValueToGo(L, val, guard)
		if err != nil {
			forEachErr = err
			return
		}
		result[keyStr] = item
	})
	if forEachErr != nil {
		return nil, forEachErr
	}
	return result, nil
}

// luaValueToGo: converts a Lua value to a Go value for template data. guard
// is documented on luaTableToGoMap.
func luaValueToGo(L *lua.LState, val lua.LValue, guard *glua.TableGuard) (interface{}, error) {
	switch v := val.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LBool:
		return bool(v), nil
	case lua.LNumber:
		return float64(v), nil
	case lua.LString:
		return string(v), nil
	case *lua.LTable:
		// analyzeAsArray inspects v without opening it in the guard again;
		// only the actual recursion below (via RawGetInt/luaTableToGoMap)
		// enters the guard, so a table is only counted once per level
		// whether it turns out to be an array or a map.
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
			if err := guard.Enter(v); err != nil {
				return nil, err
			}
			defer guard.Leave(v)

			arr := make([]interface{}, maxN)
			for i := 1; i <= maxN; i++ {
				item, err := luaValueToGo(L, v.RawGetInt(i), guard)
				if err != nil {
					return nil, err
				}
				arr[i-1] = item
			}
			return arr, nil
		}
		return luaTableToGoMap(L, v, guard)
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("template", "Go text/template rendering utilities")
	m.Fn("render", render, "renders a template string with data, raises on error",
		luareg.Args("tmpl", "data"),
		luareg.ArgDoc("tmpl", "a Go text/template source string"),
		luareg.ArgDoc("data", "table exposed to the template as the root context (dot)"),
		luareg.ReturnDoc(0, "out", "the rendered template output"))
	m.Fn("render_file", renderFile, "renders a template file with data, raises on error",
		luareg.Args("path", "data"),
		luareg.ArgDoc("path", "path to a file containing Go text/template source"),
		luareg.ArgDoc("data", "table exposed to the template as the root context (dot)"),
		luareg.ReturnDoc(0, "out", "the rendered template output"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("template", template.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
