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

package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// render: renders a Go text/template with the provided data table; raises on
// parse or execution failure. Uses *lua.LTable escape hatch for the data arg.
func render(L *lua.LState, tmplStr string, data *lua.LTable) (string, error) {
	goData := luaTableToGoMap(L, data)
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
// table; raises on read, parse, or execution failure.
func renderFile(L *lua.LState, path string, data *lua.LTable) (string, error) {
	tmplBytes, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}
	goData := luaTableToGoMap(L, data)
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
func luaTableToGoMap(L *lua.LState, tbl *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	tbl.ForEach(func(key lua.LValue, val lua.LValue) {
		var keyStr string
		if k, ok := key.(lua.LString); ok {
			keyStr = string(k)
		} else {
			keyStr = fmt.Sprintf("%v", key)
		}
		result[keyStr] = luaValueToGo(L, val)
	})
	return result
}

// luaValueToGo: converts a Lua value to a Go value for template data.
func luaValueToGo(L *lua.LState, val lua.LValue) interface{} {
	switch v := val.(type) {
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
		if isArray && maxN > 0 {
			arr := make([]interface{}, maxN)
			for i := 1; i <= maxN; i++ {
				arr[i-1] = luaValueToGo(L, v.RawGetInt(i))
			}
			return arr
		}
		return luaTableToGoMap(L, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("template", "Go text/template rendering utilities")
	m.Fn("render", render, "renders a template string with data, raises on error",
		luareg.Args("tmpl", "data"))
	m.Fn("render_file", renderFile, "renders a template file with data, raises on error",
		luareg.Args("path", "data"))
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
