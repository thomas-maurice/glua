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
	"fmt"
	"regexp"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// match: reports whether pattern matches text; raises on invalid pattern.
func match(pattern, text string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, fmt.Errorf("invalid pattern: %w", err)
	}
	return re.MatchString(text), nil
}

// find: returns the first match of pattern in text, or empty string; raises on
// invalid pattern.
func find(pattern, text string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid pattern: %w", err)
	}
	return re.FindString(text), nil
}

// findAll: returns up to n matches of pattern in text as a []string; raises on
// invalid pattern.
func findAll(pattern, text string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	matches := re.FindAllString(text, n)
	if matches == nil {
		return []string{}, nil
	}
	return matches, nil
}

// replace: replaces the first match of pattern in text with replacement; raises
// on invalid pattern.
func replace(pattern, text, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid pattern: %w", err)
	}
	// Replace only the first match.
	loc := re.FindStringIndex(text)
	if loc == nil {
		return text, nil
	}
	return text[:loc[0]] + replacement + text[loc[1]:], nil
}

// replaceAll: replaces all matches of pattern in text with replacement; raises
// on invalid pattern.
func replaceAll(pattern, text, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid pattern: %w", err)
	}
	return re.ReplaceAllLiteralString(text, replacement), nil
}

// split: splits text by pattern into at most n parts; raises on invalid pattern.
func split(pattern, text string, n int) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}
	return re.Split(text, n), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("regexp", "regular expression utilities")
	m.Fn("match", match, "reports whether pattern matches text, raises on invalid pattern",
		luareg.Args("pattern", "text"))
	m.Fn("find", find, "returns the first match of pattern in text, raises on invalid pattern",
		luareg.Args("pattern", "text"))
	m.Fn("find_all", findAll, "returns all matches of pattern in text up to n, raises on invalid pattern",
		luareg.Args("pattern", "text", "n"))
	m.Fn("replace", replace, "replaces the first match of pattern with replacement, raises on invalid pattern",
		luareg.Args("pattern", "text", "replacement"))
	m.Fn("replace_all", replaceAll, "replaces all matches of pattern with replacement, raises on invalid pattern",
		luareg.Args("pattern", "text", "replacement"))
	m.Fn("split", split, "splits text by pattern into at most n parts, raises on invalid pattern",
		luareg.Args("pattern", "text", "n"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("regexp", regexp.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
