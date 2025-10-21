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

package base64

import (
	"encoding/base64"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the base64 module for Lua.
// This function should be registered with L.PreloadModule("base64", base64.Loader)
//
// @luamodule base64
//
// Example usage in Lua:
//
//	local base64 = require("base64")
//	local encoded = base64.encode("hello world")
//	local decoded = base64.decode(encoded)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"encode":     encode,
	"decode":     decode,
	"encode_url": encodeURL,
	"decode_url": decodeURL,
}

// encode: encodes a string to base64.
//
// @luafunc encode
// @luaparam str string The string to encode
// @luareturn string encoded The base64 encoded string
//
// Example:
//
//	local encoded = base64.encode("hello world")
//	print(encoded)  -- prints "aGVsbG8gd29ybGQ="
func encode(L *lua.LState) int {
	str := L.CheckString(1)
	encoded := base64.StdEncoding.EncodeToString([]byte(str))
	L.Push(lua.LString(encoded))
	return 1
}

// decode: decodes a base64 string.
//
// @luafunc decode
// @luaparam encoded string The base64 encoded string
// @luareturn string decoded The decoded string, or nil on error
// @luareturn string|nil err Error message if decoding failed
//
// Example:
//
//	local decoded, err = base64.decode("aGVsbG8gd29ybGQ=")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(decoded)  -- prints "hello world"
//	end
func decode(L *lua.LState) int {
	encoded := L.CheckString(1)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to decode base64: %v", err)))
		return 2
	}
	L.Push(lua.LString(string(decoded)))
	L.Push(lua.LNil)
	return 2
}

// encode_url: encodes a string to URL-safe base64.
//
// @luafunc encode_url
// @luaparam str string The string to encode
// @luareturn string encoded The URL-safe base64 encoded string
//
// Example:
//
//	local encoded = base64.encode_url("hello world")
func encodeURL(L *lua.LState) int {
	str := L.CheckString(1)
	encoded := base64.URLEncoding.EncodeToString([]byte(str))
	L.Push(lua.LString(encoded))
	return 1
}

// decode_url: decodes a URL-safe base64 string.
//
// @luafunc decode_url
// @luaparam encoded string The URL-safe base64 encoded string
// @luareturn string decoded The decoded string, or nil on error
// @luareturn string|nil err Error message if decoding failed
//
// Example:
//
//	local decoded, err = base64.decode_url("aGVsbG8gd29ybGQ")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(decoded)  -- prints "hello world"
//	end
func decodeURL(L *lua.LState) int {
	encoded := L.CheckString(1)
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to decode URL-safe base64: %v", err)))
		return 2
	}
	L.Push(lua.LString(string(decoded)))
	L.Push(lua.LNil)
	return 2
}
