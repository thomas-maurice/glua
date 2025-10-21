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

package hex

import (
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the hex module for Lua.
// This function should be registered with L.PreloadModule("hex", hex.Loader)
//
// @luamodule hex
//
// Example usage in Lua:
//
//	local hex = require("hex")
//	local encoded = hex.encode("hello")
//	local decoded = hex.decode(encoded)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"encode": encode,
	"decode": decode,
}

// encode: encodes a string to hexadecimal.
//
// @luafunc encode
// @luaparam str string The string to encode
// @luareturn string encoded The hex encoded string
//
// Example:
//
//	local encoded = hex.encode("hello")
//	print(encoded)  -- prints "68656c6c6f"
func encode(L *lua.LState) int {
	str := L.CheckString(1)
	encoded := hex.EncodeToString([]byte(str))
	L.Push(lua.LString(encoded))
	return 1
}

// decode: decodes a hexadecimal string.
//
// @luafunc decode
// @luaparam encoded string The hex encoded string
// @luareturn string decoded The decoded string, or nil on error
// @luareturn string|nil err Error message if decoding failed
//
// Example:
//
//	local decoded, err = hex.decode("68656c6c6f")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(decoded)  -- prints "hello"
//	end
func decode(L *lua.LState) int {
	encoded := L.CheckString(1)
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to decode hex: %v", err)))
		return 2
	}
	L.Push(lua.LString(string(decoded)))
	L.Push(lua.LNil)
	return 2
}
