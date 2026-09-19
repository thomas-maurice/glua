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

// Package hex provides hexadecimal encoding and decoding utilities for Lua
// scripts.
package hex

import (
	"encoding/hex"
	"fmt"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// encode: hex-encodes a string, returning the hex string.
func encode(s string) string {
	return hex.EncodeToString([]byte(s))
}

// decode: hex-decodes a string, returning the decoded string or an error.
func decode(s string) (string, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex: %w", err)
	}
	return string(b), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("hex", "hexadecimal encoding and decoding utilities")
	m.Fn("encode", encode, "encodes a string to hexadecimal",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to encode"),
		luareg.ReturnDoc(0, "encoded", "the lowercase hexadecimal encoding of s"))
	m.Fn("decode", decode, "decodes a hexadecimal string, raises on invalid input",
		luareg.Args("encoded"),
		luareg.ArgDoc("encoded", "a hexadecimal string with an even number of digits"),
		luareg.ReturnDoc(0, "s", "the decoded raw string"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("hex", hex.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
