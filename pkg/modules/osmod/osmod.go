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

package osmod

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates the osmod Lua module
//
// @luamodule osmod
func Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), exports)
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"getenv":   getenv,
	"setenv":   setenv,
	"unsetenv": unsetenv,
	"hostname": hostname,
	"tmpdir":   tmpdir,
}

// getenv: gets the value of an environment variable
//
// @luafunc getenv
// @luaparam name string The environment variable name
// @luareturn string The value of the environment variable (empty string if not set)
func getenv(L *lua.LState) int {
	name := L.CheckString(1)
	value := os.Getenv(name)
	L.Push(lua.LString(value))
	return 1
}

// setenv: sets the value of an environment variable
//
// @luafunc setenv
// @luaparam name string The environment variable name
// @luaparam value string The value to set
// @luareturn string|nil Error message if operation failed
func setenv(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckString(2)

	err := os.Setenv(name, value)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// unsetenv: unsets an environment variable
//
// @luafunc unsetenv
// @luaparam name string The environment variable name
// @luareturn string|nil Error message if operation failed
func unsetenv(L *lua.LState) int {
	name := L.CheckString(1)

	err := os.Unsetenv(name)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// hostname: gets the system hostname
//
// @luafunc hostname
// @luareturn string The hostname
// @luareturn string|nil Error message if operation failed
func hostname(L *lua.LState) int {
	hostname, err := os.Hostname()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(hostname))
	L.Push(lua.LNil)
	return 2
}

// tmpdir: gets the system temporary directory path
//
// @luafunc tmpdir
// @luareturn string The temporary directory path
func tmpdir(L *lua.LState) int {
	tmpdir := os.TempDir()
	L.Push(lua.LString(tmpdir))
	return 1
}
