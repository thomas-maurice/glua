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

package regexp

import (
	"regexp"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates the regexp Lua module
//
// @luamodule regexp
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"match":       matchFunc,
	"find":        findFunc,
	"find_all":    findAllFunc,
	"replace":     replaceFunc,
	"replace_all": replaceAllFunc,
	"split":       splitFunc,
}

// matchFunc: checks if pattern matches text
//
// @luafunc match
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to match against
// @luareturn boolean True if pattern matches, false otherwise
// @luareturn string|nil Error message if pattern is invalid
func matchFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	matched := re.MatchString(text)
	L.Push(lua.LBool(matched))
	L.Push(lua.LNil)
	return 2
}

// findFunc: finds first match of pattern in text
//
// @luafunc find
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to search
// @luareturn string The first match (empty string if no match)
// @luareturn string|nil Error message if pattern is invalid
func findFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LString(""))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	match := re.FindString(text)
	L.Push(lua.LString(match))
	L.Push(lua.LNil)
	return 2
}

// findAllFunc: finds all matches of pattern in text
//
// @luafunc find_all
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to search
// @luaparam limit number Maximum number of matches (-1 for all)
// @luareturn table Array of matches
// @luareturn string|nil Error message if pattern is invalid
func findAllFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)
	limit := L.CheckInt(3)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(L.NewTable())
		L.Push(lua.LString(err.Error()))
		return 2
	}

	matches := re.FindAllString(text, limit)

	table := L.NewTable()
	for i, match := range matches {
		table.RawSetInt(i+1, lua.LString(match))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}

// replaceFunc: replaces first match of pattern with replacement
//
// @luafunc replace
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to search
// @luaparam replacement string The replacement string
// @luareturn string The text with first match replaced
// @luareturn string|nil Error message if pattern is invalid
func replaceFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)
	replacement := L.CheckString(3)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LString(text))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	// For single replacement, we need to limit to 1
	result := re.ReplaceAllStringFunc(text, func(s string) string {
		return replacement
	})

	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// replaceAllFunc: replaces all matches of pattern with replacement
//
// @luafunc replace_all
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to search
// @luaparam replacement string The replacement string
// @luareturn string The text with all matches replaced
// @luareturn string|nil Error message if pattern is invalid
func replaceAllFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)
	replacement := L.CheckString(3)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LString(text))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	result := re.ReplaceAllLiteralString(text, replacement)
	L.Push(lua.LString(result))
	L.Push(lua.LNil)
	return 2
}

// splitFunc: splits text by pattern
//
// @luafunc split
// @luaparam pattern string The regular expression pattern
// @luaparam text string The text to split
// @luaparam limit number Maximum number of splits (-1 for all)
// @luareturn table Array of split parts
// @luareturn string|nil Error message if pattern is invalid
func splitFunc(L *lua.LState) int {
	pattern := L.CheckString(1)
	text := L.CheckString(2)
	limit := L.CheckInt(3)

	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(L.NewTable())
		L.Push(lua.LString(err.Error()))
		return 2
	}

	parts := re.Split(text, limit)

	table := L.NewTable()
	for i, part := range parts {
		table.RawSetInt(i+1, lua.LString(part))
	}

	L.Push(table)
	L.Push(lua.LNil)
	return 2
}
