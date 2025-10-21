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

package filepath

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates the filepath Lua module
//
// @luamodule filepath
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"join":  join,
	"split": split,
	"abs":   abs,
	"ext":   ext,
	"base":  base,
	"dir":   dir,
	"clean": clean,
}

// join: joins path elements into a single path
//
// @luafunc join
// @luaparam ... string Path elements to join
// @luareturn string The joined path
func join(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LString(""))
		return 1
	}

	parts := make([]string, n)
	for i := 1; i <= n; i++ {
		parts[i-1] = L.CheckString(i)
	}

	result := filepath.Join(parts...)
	L.Push(lua.LString(result))
	return 1
}

// split: splits path into directory and file
//
// @luafunc split
// @luaparam path string The path to split
// @luareturn string The directory part
// @luareturn string The file part
func split(L *lua.LState) int {
	path := L.CheckString(1)
	dir, file := filepath.Split(path)
	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))
	return 2
}

// abs: returns absolute path
//
// @luafunc abs
// @luaparam path string The path to make absolute
// @luareturn string The absolute path
// @luareturn string|nil Error message if operation failed
func abs(L *lua.LState) int {
	path := L.CheckString(1)
	absPath, err := filepath.Abs(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(absPath))
	L.Push(lua.LNil)
	return 2
}

// ext: returns the file extension
//
// @luafunc ext
// @luaparam path string The file path
// @luareturn string The file extension (including the dot)
func ext(L *lua.LState) int {
	path := L.CheckString(1)
	extension := filepath.Ext(path)
	L.Push(lua.LString(extension))
	return 1
}

// base: returns the last element of path
//
// @luafunc base
// @luaparam path string The file path
// @luareturn string The base name
func base(L *lua.LState) int {
	path := L.CheckString(1)
	baseName := filepath.Base(path)
	L.Push(lua.LString(baseName))
	return 1
}

// dir: returns all but the last element of path
//
// @luafunc dir
// @luaparam path string The file path
// @luareturn string The directory path
func dir(L *lua.LState) int {
	path := L.CheckString(1)
	dirPath := filepath.Dir(path)
	L.Push(lua.LString(dirPath))
	return 1
}

// clean: returns the shortest path equivalent to path
//
// @luafunc clean
// @luaparam path string The path to clean
// @luareturn string The cleaned path
func clean(L *lua.LState) int {
	path := L.CheckString(1)
	cleanPath := filepath.Clean(path)
	L.Push(lua.LString(cleanPath))
	return 1
}
