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

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// setenv: sets an environment variable; raises on error, returns true on success.
func setenv(name, value string) (bool, error) {
	if err := os.Setenv(name, value); err != nil {
		return false, err
	}
	return true, nil
}

// unsetenv: unsets an environment variable; raises on error, returns true on success.
func unsetenv(name string) (bool, error) {
	if err := os.Unsetenv(name); err != nil {
		return false, err
	}
	return true, nil
}

// hostname: returns the system hostname; raises on error.
func hostname() (string, error) {
	return os.Hostname()
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("osmod", "operating system utilities")
	m.Fn("getenv", os.Getenv, "returns the value of an environment variable",
		luareg.Args("name"))
	m.Fn("setenv", setenv, "sets an environment variable, raises on error",
		luareg.Args("name", "value"))
	m.Fn("unsetenv", unsetenv, "unsets an environment variable, raises on error",
		luareg.Args("name"))
	m.Fn("hostname", hostname, "returns the system hostname, raises on error")
	m.Fn("tmpdir", os.TempDir, "returns the default temporary directory path")
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("osmod", osmod.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
