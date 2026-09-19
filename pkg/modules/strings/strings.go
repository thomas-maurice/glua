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
	"fmt"
	"strings"
	"unicode"

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

// index: returns the 1-based position of the first occurrence of substr in
// s, or 0 if substr is not present. Deliberately NOT Go's strings.Index
// convention (0-based, -1 for absent): Lua's own string.find is 1-based, so
// this return value composes directly with string.sub without an off-by-one
// adjustment at every call site.
func index(s, substr string) int {
	i := strings.Index(s, substr)
	if i < 0 {
		return 0
	}
	return i + 1
}

// lastIndex: like index, but for the last occurrence. Same 1-based/0-absent
// deviation from Go's strings.LastIndex (0-based, -1 for absent).
func lastIndex(s, substr string) int {
	i := strings.LastIndex(s, substr)
	if i < 0 {
		return 0
	}
	return i + 1
}

// rep: repeats s count times. Named rep, not repeat: repeat is a reserved
// word in Lua, so strings.repeat(s, n) is a syntax error at the call site,
// not a runtime one. Go's strings.Repeat panics on a negative count; that is
// validated here and turned into a raised Lua error instead, per this
// module's fail-loud convention.
func rep(s string, count int) (string, error) {
	if count < 0 {
		return "", fmt.Errorf("strings.rep: count must be >= 0, got %d", count)
	}
	return strings.Repeat(s, count), nil
}

// splitN: like strings.SplitN, but always returns a non-nil slice. Go's
// SplitN returns a nil slice when n == 0 (zero substrings requested); the
// Translator marshals a nil slice as Lua nil rather than an empty table,
// which would break "#parts" and "ipairs(parts)" in a caller that did not
// special-case n == 0. Normalizing to a non-nil empty slice here keeps the
// Lua-visible contract uniform: split_n always returns a table.
func splitN(s, sep string, n int) []string {
	parts := strings.SplitN(s, sep, n)
	if parts == nil {
		return []string{}
	}
	return parts
}

// title: returns s with the first rune of every word title-cased and every
// other rune left untouched. A "word" starts at the beginning of the string
// or immediately after a run of unicode.IsSpace characters; the original
// whitespace (including width and kind) is preserved verbatim in the output.
//
// This deliberately does NOT call the standard library's strings.Title,
// which is documented as deprecated because it mishandles Unicode special
// casing (e.g. it upper-cases both letters of the Dutch "ij" digraph).
// unicode.ToTitle is used instead of unicode.ToUpper because a handful of
// Unicode digraph code points (e.g. U+01C4..U+01CC, the "DZ"/"Dz"/"dz"
// family) have a titlecase form distinct from their uppercase form, and
// title-casing is the semantically correct mapping for the first letter of
// a word.
//
// This is first-rune title casing only, not a language-aware word-breaking
// or casing algorithm: it has no notion of locale-specific exceptions (e.g.
// English "of"/"the" staying lowercase in real title case).
func title(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	atWordStart := true
	for _, r := range s {
		if atWordStart && !unicode.IsSpace(r) {
			b.WriteRune(unicode.ToTitle(r))
		} else {
			b.WriteRune(r)
		}
		atWordStart = unicode.IsSpace(r)
	}
	return b.String()
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
	m.Fn("trim_space", strings.TrimSpace, "removes leading and trailing whitespace from a string",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ReturnDoc(0, "out", "s with leading and trailing Unicode whitespace removed"))
	m.Fn("trim_prefix", strings.TrimPrefix, "removes a leading prefix from a string, if present",
		luareg.Args("s", "prefix"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ArgDoc("prefix", "the prefix to remove"),
		luareg.ReturnDoc(0, "out", "s with prefix removed, or s unchanged if it does not start with prefix"))
	m.Fn("trim_suffix", strings.TrimSuffix, "removes a trailing suffix from a string, if present",
		luareg.Args("s", "suffix"),
		luareg.ArgDoc("s", "the string to trim"),
		luareg.ArgDoc("suffix", "the suffix to remove"),
		luareg.ReturnDoc(0, "out", "s with suffix removed, or s unchanged if it does not end with suffix"))
	m.Fn("fields", strings.Fields, "splits a string around runs of whitespace",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to split"),
		luareg.ReturnDoc(0, "parts", "table (array) of substrings between runs of Unicode whitespace; leading/trailing whitespace produces no empty entries"))
	m.Fn("rep", rep, "repeats a string count times",
		luareg.Args("s", "count"),
		luareg.ArgDoc("s", "the string to repeat"),
		luareg.ArgDoc("count", "number of repetitions; must be >= 0"),
		luareg.ReturnDoc(0, "out", "s repeated count times"))
	m.Fn("index", index, "returns the 1-based position of the first occurrence of substr in s, or 0 if absent",
		luareg.Args("s", "substr"),
		luareg.ArgDoc("s", "the string to search"),
		luareg.ArgDoc("substr", "the substring to look for"),
		luareg.ReturnDoc(0, "pos", "1-based index of the first occurrence, or 0 if not found; NOTE this deviates from Go's strings.Index, which is 0-based and returns -1 for not found — Lua's own string.find convention is used instead so the result composes with string.sub"))
	m.Fn("last_index", lastIndex, "returns the 1-based position of the last occurrence of substr in s, or 0 if absent",
		luareg.Args("s", "substr"),
		luareg.ArgDoc("s", "the string to search"),
		luareg.ArgDoc("substr", "the substring to look for"),
		luareg.ReturnDoc(0, "pos", "1-based index of the last occurrence, or 0 if not found; same 1-based/0-absent convention as index, deviating from Go's strings.LastIndex"))
	m.Fn("equal_fold", strings.EqualFold, "reports whether two strings are equal under simple Unicode case-folding",
		luareg.Args("a", "b"),
		luareg.ArgDoc("a", "the first string"),
		luareg.ArgDoc("b", "the second string"),
		luareg.ReturnDoc(0, "ok", "true if a and b are equal under case-insensitive comparison"))
	m.Fn("title", title, "title-cases the first rune of every whitespace-separated word",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to title-case"),
		luareg.ReturnDoc(0, "out", "s with the first rune of each word title-cased; ASCII/first-rune only, not language-aware (does not implement locale casing exceptions)"))
	m.Fn("cut", strings.Cut, "slices s around the first instance of sep",
		luareg.Args("s", "sep"),
		luareg.ArgDoc("s", "the string to slice"),
		luareg.ArgDoc("sep", "the separator to find"),
		luareg.ReturnDoc(0, "before", "the portion of s before the first occurrence of sep, or all of s if sep is not found"),
		luareg.ReturnDoc(1, "after", "the portion of s after the first occurrence of sep, or empty if sep is not found"),
		luareg.ReturnDoc(2, "found", "true if sep occurs in s"))
	m.Fn("split_n", splitN, "splits a string by separator into at most n substrings",
		luareg.Args("s", "sep", "n"),
		luareg.ArgDoc("s", "the string to split"),
		luareg.ArgDoc("sep", "the separator; if empty, splits after every UTF-8 character"),
		luareg.ArgDoc("n", "maximum number of substrings: n > 0 limits the result to at most n elements (the last one unsplit); n == 0 returns an empty table; n < 0 splits all occurrences, like split"),
		luareg.ReturnDoc(0, "parts", "table (array) of at most n substrings"))
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
