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

package strings

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates the strings Lua module
//
// @luamodule strings
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"has_prefix": hasPrefix,
	"has_suffix": hasSuffix,
	"trim":       trim,
	"trim_left":  trimLeft,
	"trim_right": trimRight,
	"split":      split,
	"join":       join,
	"to_upper":   toUpper,
	"to_lower":   toLower,
	"contains":   contains,
	"count":      count,
	"replace":    replace,
}

// hasPrefix: checks if string has prefix
//
// @luafunc has_prefix
// @luaparam s string The string to check
// @luaparam prefix string The prefix to look for
// @luareturn boolean True if string has prefix
func hasPrefix(L *lua.LState) int {
	s := L.CheckString(1)
	prefix := L.CheckString(2)
	L.Push(lua.LBool(strings.HasPrefix(s, prefix)))
	return 1
}

// hasSuffix: checks if string has suffix
//
// @luafunc has_suffix
// @luaparam s string The string to check
// @luaparam suffix string The suffix to look for
// @luareturn boolean True if string has suffix
func hasSuffix(L *lua.LState) int {
	s := L.CheckString(1)
	suffix := L.CheckString(2)
	L.Push(lua.LBool(strings.HasSuffix(s, suffix)))
	return 1
}

// trim: removes cutset from both ends of string
//
// @luafunc trim
// @luaparam s string The string to trim
// @luaparam cutset string The characters to remove
// @luareturn string The trimmed string
func trim(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)
	result := strings.Trim(s, cutset)
	L.Push(lua.LString(result))
	return 1
}

// trimLeft: removes cutset from left end of string
//
// @luafunc trim_left
// @luaparam s string The string to trim
// @luaparam cutset string The characters to remove
// @luareturn string The trimmed string
func trimLeft(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)
	result := strings.TrimLeft(s, cutset)
	L.Push(lua.LString(result))
	return 1
}

// trimRight: removes cutset from right end of string
//
// @luafunc trim_right
// @luaparam s string The string to trim
// @luaparam cutset string The characters to remove
// @luareturn string The trimmed string
func trimRight(L *lua.LState) int {
	s := L.CheckString(1)
	cutset := L.CheckString(2)
	result := strings.TrimRight(s, cutset)
	L.Push(lua.LString(result))
	return 1
}

// split: splits string by separator
//
// @luafunc split
// @luaparam s string The string to split
// @luaparam sep string The separator
// @luareturn table Array of split parts
func split(L *lua.LState) int {
	s := L.CheckString(1)
	sep := L.CheckString(2)
	parts := strings.Split(s, sep)

	table := L.NewTable()
	for i, part := range parts {
		table.RawSetInt(i+1, lua.LString(part))
	}

	L.Push(table)
	return 1
}

// join: joins array of strings with separator
//
// @luafunc join
// @luaparam parts table Array of strings to join
// @luaparam sep string The separator
// @luareturn string The joined string
func join(L *lua.LState) int {
	partsTable := L.CheckTable(1)
	sep := L.CheckString(2)

	var parts []string
	partsTable.ForEach(func(k, v lua.LValue) {
		if str, ok := v.(lua.LString); ok {
			parts = append(parts, string(str))
		}
	})

	result := strings.Join(parts, sep)
	L.Push(lua.LString(result))
	return 1
}

// toUpper: converts string to uppercase
//
// @luafunc to_upper
// @luaparam s string The string to convert
// @luareturn string The uppercase string
func toUpper(L *lua.LState) int {
	s := L.CheckString(1)
	L.Push(lua.LString(strings.ToUpper(s)))
	return 1
}

// toLower: converts string to lowercase
//
// @luafunc to_lower
// @luaparam s string The string to convert
// @luareturn string The lowercase string
func toLower(L *lua.LState) int {
	s := L.CheckString(1)
	L.Push(lua.LString(strings.ToLower(s)))
	return 1
}

// contains: checks if string contains substring
//
// @luafunc contains
// @luaparam s string The string to search
// @luaparam substr string The substring to find
// @luareturn boolean True if string contains substring
func contains(L *lua.LState) int {
	s := L.CheckString(1)
	substr := L.CheckString(2)
	L.Push(lua.LBool(strings.Contains(s, substr)))
	return 1
}

// count: counts occurrences of substring
//
// @luafunc count
// @luaparam s string The string to search
// @luaparam substr string The substring to count
// @luareturn number The count of occurrences
func count(L *lua.LState) int {
	s := L.CheckString(1)
	substr := L.CheckString(2)
	L.Push(lua.LNumber(strings.Count(s, substr)))
	return 1
}

// replace: replaces occurrences of old with new
//
// @luafunc replace
// @luaparam s string The string to search
// @luaparam old string The substring to replace
// @luaparam new string The replacement string
// @luaparam n number Number of replacements (-1 for all)
// @luareturn string The string with replacements
func replace(L *lua.LState) int {
	s := L.CheckString(1)
	old := L.CheckString(2)
	new := L.CheckString(3)
	n := L.CheckInt(4)

	result := strings.Replace(s, old, new, n)
	L.Push(lua.LString(result))
	return 1
}
