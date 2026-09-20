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

// Package filepath provides file path manipulation utilities for Lua scripts,
// wrapping Go's stdlib path/filepath.
package filepath

import (
	"path/filepath"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// join: joins a slice of path elements. Wraps the variadic filepath.Join so
// luareg can handle it as a []string parameter.
func join(elem []string) string {
	return filepath.Join(elem...)
}

// splitPath: returns (dir, file) components of a path.
func splitPath(path string) (string, string) {
	return filepath.Split(path)
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("filepath", "file path manipulation utilities")
	m.Fn("join", join, "joins a table of path elements into a single path",
		luareg.Args("elem"),
		luareg.ArgDoc("elem", "table (array) of path segments to join with the OS path separator"),
		luareg.ReturnDoc(0, "path", "the joined and Clean-ed path"))
	m.Fn("split", splitPath, "splits a path into directory and file components",
		luareg.Args("path"),
		luareg.ArgDoc("path", "the path to split, e.g. \"/a/b/c.txt\""),
		luareg.ReturnDoc(0, "dir", "everything up to and including the final separator, e.g. \"/a/b/\""),
		luareg.ReturnDoc(1, "file", "everything after the final separator, e.g. \"c.txt\""))
	m.Fn("abs", filepath.Abs, "returns the absolute form of the path, raises on error",
		luareg.Args("path"),
		luareg.ArgDoc("path", "relative or absolute path to resolve against the process working directory"),
		luareg.ReturnDoc(0, "path", "the absolute, Clean-ed form of path"))
	m.Fn("ext", filepath.Ext, "returns the file extension including the dot",
		luareg.Args("path"),
		luareg.ArgDoc("path", "the path whose extension to extract"),
		luareg.ReturnDoc(0, "ext", "the file extension including the leading dot, or empty string if none"))
	m.Fn("base", filepath.Base, "returns the last element of the path",
		luareg.Args("path"),
		luareg.ArgDoc("path", "the path whose last element to extract"),
		luareg.ReturnDoc(0, "name", "the last path element, with trailing separators removed"))
	m.Fn("dir", filepath.Dir, "returns all but the last element of the path",
		luareg.Args("path"),
		luareg.ArgDoc("path", "the path whose parent directory to extract"),
		luareg.ReturnDoc(0, "dir", "all but the last element of path"))
	m.Fn("clean", filepath.Clean, "returns the shortest path equivalent to path",
		luareg.Args("path"),
		luareg.ArgDoc("path", "the path to simplify, e.g. containing \"..\" or repeated separators"),
		luareg.ReturnDoc(0, "path", "the shortest path lexically equivalent to path"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("filepath", filepath.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
