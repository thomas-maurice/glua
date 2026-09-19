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

// Package strings provides string manipulation utilities for Lua scripts,
// wrapping Go's stdlib strings package.
package strings

import (
	"strings"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// join: joins array values with separator, coercing all elements to strings
// via Lua's tostring semantics. Uses *lua.LState escape hatch so non-string
// table elements (numbers, booleans) are coerced rather than rejected.
func join(_ *lua.LState, partsTable *lua.LTable, sep string) string {
	n := partsTable.Len()
	parts := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		parts = append(parts, partsTable.RawGetInt(i).String())
	}
	return strings.Join(parts, sep)
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("strings", "string manipulation utilities")
	m.Fn("has_prefix", strings.HasPrefix, "checks if string has prefix",
		luareg.Args("s", "prefix"),
		luareg.ArgDoc("s", "the string to test"),
		luareg.ArgDoc("prefix", "the prefix to look for"),
		luareg.ReturnDoc(0, "ok", "true if s starts with prefix"))
	m.Fn("has_suffix", strings.HasSuffix, "checks if string has suffix",
		luareg.Args("s", "suffix"),
		luareg.ArgDoc("s", "the string to test"),
		luareg.ArgDoc("suffix", "the suffix to look for"),
		luareg.ReturnDoc(0, "ok", "true if s ends with suffix"))
	m.Fn("trim", strings.Trim, "removes cutset characters from both ends of a string",
		luareg.Args("s", "cutset"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ArgDoc("cutset", "the set of characters to remove from both ends"),
		luareg.ReturnDoc(0, "out", "s with leading and trailing cutset characters removed"))
	m.Fn("trim_left", strings.TrimLeft, "removes cutset characters from the left end of a string",
		luareg.Args("s", "cutset"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ArgDoc("cutset", "the set of characters to remove from the left"),
		luareg.ReturnDoc(0, "out", "s with leading cutset characters removed"))
	m.Fn("trim_right", strings.TrimRight, "removes cutset characters from the right end of a string",
		luareg.Args("s", "cutset"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ArgDoc("cutset", "the set of characters to remove from the right"),
		luareg.ReturnDoc(0, "out", "s with trailing cutset characters removed"))
	m.Fn("split", strings.Split, "splits a string by separator into a table",
		luareg.Args("s", "sep"),
		luareg.ArgDoc("s", "the string to split"),
		luareg.ArgDoc("sep", "the separator; if empty, splits after every UTF-8 character"),
		luareg.ReturnDoc(0, "parts", "table (array) of substrings between occurrences of sep"))
	m.Fn("join", join, "joins a table of strings with a separator, coercing values to strings",
		luareg.Args("parts", "sep"),
		luareg.ArgDoc("parts", "table (array) of values to join; non-string values are coerced via tostring"),
		luareg.ArgDoc("sep", "the separator to place between elements"),
		luareg.ReturnDoc(0, "out", "the joined string"))
	m.Fn("to_upper", strings.ToUpper, "converts a string to uppercase",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to convert"),
		luareg.ReturnDoc(0, "out", "s with all letters mapped to upper case"))
	m.Fn("to_lower", strings.ToLower, "converts a string to lowercase",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to convert"),
		luareg.ReturnDoc(0, "out", "s with all letters mapped to lower case"))
	m.Fn("contains", strings.Contains, "checks if a string contains a substring",
		luareg.Args("s", "substr"),
		luareg.ArgDoc("s", "the string to search"),
		luareg.ArgDoc("substr", "the substring to look for"),
		luareg.ReturnDoc(0, "ok", "true if substr appears anywhere in s"))
	m.Fn("count", strings.Count, "counts occurrences of substr in s",
		luareg.Args("s", "substr"),
		luareg.ArgDoc("s", "the string to search"),
		luareg.ArgDoc("substr", "the non-overlapping substring to count; empty string counts UTF-8 characters plus one"),
		luareg.ReturnDoc(0, "n", "the number of non-overlapping instances of substr in s"))
	m.Fn("replace", strings.Replace, "replaces occurrences of old with new in s",
		luareg.Args("s", "old", "new", "n"),
		luareg.ArgDoc("s", "the string to search"),
		luareg.ArgDoc("old", "the substring to replace"),
		luareg.ArgDoc("new", "the replacement substring"),
		luareg.ArgDoc("n", "maximum number of replacements; a negative value replaces all occurrences"),
		luareg.ReturnDoc(0, "out", "s with up to n occurrences of old replaced by new"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("strings", strings.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
