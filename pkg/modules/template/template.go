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

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the template module for Lua.
// This function should be registered with L.PreloadModule("template", template.Loader)
//
// @luamodule template
//
// Example usage in Lua:
//
//	local template = require("template")
//	local result = template.render("Hello {{.Name}}", {Name = "World"})
//	print(result)  -- prints "Hello World"
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"render":      render,
	"render_file": renderFile,
}

// render: renders a Go text/template with the provided data.
//
// @luafunc render
// @luaparam tmpl string The template string
// @luaparam data table The data to render with
// @luareturn result string The rendered template, or nil on error
// @luareturn err string|nil Error message if rendering failed
//
// Example:
//
//	local result, err = template.render("Hello {{.Name}}", {Name = "World"})
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(result)  -- prints "Hello World"
//	end
func render(L *lua.LState) int {
	tmplStr := L.CheckString(1)
	data := L.CheckTable(2)

	// Convert Lua table to Go map
	goData := luaTableToGoMap(L, data)

	// Parse template
	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse template: %v", err)))
		return 2
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, goData); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to execute template: %v", err)))
		return 2
	}

	L.Push(lua.LString(buf.String()))
	L.Push(lua.LNil)
	return 2
}

// render_file: renders a Go text/template from a file with the provided data.
//
// @luafunc render_file
// @luaparam path string The path to the template file
// @luaparam data table The data to render with
// @luareturn result string The rendered template, or nil on error
// @luareturn err string|nil Error message if rendering failed
//
// Example:
//
//	local result, err = template.render_file("/path/to/template.tmpl", {Name = "World"})
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(result)
//	end
func renderFile(L *lua.LState) int {
	path := L.CheckString(1)
	data := L.CheckTable(2)

	// Read template file
	tmplBytes, err := os.ReadFile(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read template file: %v", err)))
		return 2
	}

	// Convert Lua table to Go map
	goData := luaTableToGoMap(L, data)

	// Parse template
	tmpl, err := template.New("tmpl").Parse(string(tmplBytes))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse template: %v", err)))
		return 2
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, goData); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to execute template: %v", err)))
		return 2
	}

	L.Push(lua.LString(buf.String()))
	L.Push(lua.LNil)
	return 2
}

// luaTableToGoMap: converts a Lua table to a Go map for template rendering
func luaTableToGoMap(L *lua.LState, tbl *lua.LTable) map[string]interface{} {
	result := make(map[string]interface{})
	tbl.ForEach(func(key lua.LValue, val lua.LValue) {
		keyStr := ""
		if k, ok := key.(lua.LString); ok {
			keyStr = string(k)
		} else {
			keyStr = fmt.Sprintf("%v", key)
		}
		result[keyStr] = luaValueToGo(L, val)
	})
	return result
}

// luaValueToGo: converts a Lua value to a Go value
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
		// Check if it's an array or map
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
