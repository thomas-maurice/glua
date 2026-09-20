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

// Package base64 provides base64 encoding and decoding utilities for Lua
// scripts (standard and URL-safe alphabets).
package base64

import (
	"encoding/base64"
	"fmt"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// encode: base64-encodes a string using standard encoding.
func encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// decode: decodes a standard base64 string; raises on invalid input.
func decode(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}
	return string(b), nil
}

// encodeURL: base64-encodes a string using URL-safe encoding.
func encodeURL(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// decodeURL: decodes a URL-safe base64 string; raises on invalid input.
func decodeURL(s string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode URL-safe base64: %w", err)
	}
	return string(b), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("base64", "base64 encoding and decoding utilities")
	m.Fn("encode", encode, "encodes a string to standard base64",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to encode"),
		luareg.ReturnDoc(0, "encoded", "the standard-alphabet base64 encoding of s"))
	m.Fn("decode", decode, "decodes a standard base64 string, raises on invalid input",
		luareg.Args("encoded"),
		luareg.ArgDoc("encoded", "a standard-alphabet base64 string to decode"),
		luareg.ReturnDoc(0, "s", "the decoded raw string"))
	m.Fn("encode_url", encodeURL, "encodes a string to URL-safe base64",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to encode"),
		luareg.ReturnDoc(0, "encoded", "the URL-safe base64 encoding of s"))
	m.Fn("decode_url", decodeURL, "decodes a URL-safe base64 string, raises on invalid input",
		luareg.Args("encoded"),
		luareg.ArgDoc("encoded", "a URL-safe base64 string to decode"),
		luareg.ReturnDoc(0, "s", "the decoded raw string"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("base64", base64.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
