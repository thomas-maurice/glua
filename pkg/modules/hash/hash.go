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

package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"

	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the hash module for Lua.
// This function should be registered with L.PreloadModule("hash", hash.Loader)
//
// @luamodule hash
//
// Example usage in Lua:
//
//	local hash = require("hash")
//	local h = hash.sha256("hello world")
//	print(h)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"md5":         md5Hash,
	"sha1":        sha1Hash,
	"sha256":      sha256Hash,
	"sha512":      sha512Hash,
	"hmac_sha256": hmacSha256,
	"md5_obj":     md5HashObj,
	"sha1_obj":    sha1HashObj,
	"sha256_obj":  sha256HashObj,
	"sha512_obj":  sha512HashObj,
}

// md5Hash: computes the MD5 hash of a string.
//
// @luafunc md5
// @luaparam str string The string to hash
// @luareturn string hash The hex-encoded MD5 hash
//
// Example:
//
//	local h = hash.md5("hello world")
//	print(h)  -- prints "5eb63bbbe01eeed093cb22bb8f5acdc3"
func md5Hash(L *lua.LState) int {
	str := L.CheckString(1)
	h := md5.Sum([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	return 1
}

// sha1Hash: computes the SHA1 hash of a string.
//
// @luafunc sha1
// @luaparam str string The string to hash
// @luareturn string hash The hex-encoded SHA1 hash
//
// Example:
//
//	local h = hash.sha1("hello world")
//	print(h)  -- prints "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"
func sha1Hash(L *lua.LState) int {
	str := L.CheckString(1)
	h := sha1.Sum([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	return 1
}

// sha256Hash: computes the SHA256 hash of a string.
//
// @luafunc sha256
// @luaparam str string The string to hash
// @luareturn string hash The hex-encoded SHA256 hash
//
// Example:
//
//	local h = hash.sha256("hello world")
//	print(h)  -- prints "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
func sha256Hash(L *lua.LState) int {
	str := L.CheckString(1)
	h := sha256.Sum256([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	return 1
}

// sha512Hash: computes the SHA512 hash of a string.
//
// @luafunc sha512
// @luaparam str string The string to hash
// @luareturn string hash The hex-encoded SHA512 hash
//
// Example:
//
//	local h = hash.sha512("hello world")
func sha512Hash(L *lua.LState) int {
	str := L.CheckString(1)
	h := sha512.Sum512([]byte(str))
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	return 1
}

// hmacSha256: computes the HMAC-SHA256 of a message with a key.
//
// @luafunc hmac_sha256
// @luaparam message string The message to authenticate
// @luaparam key string The secret key
// @luareturn string hash The hex-encoded HMAC-SHA256
//
// Example:
//
//	local h = hash.hmac_sha256("message", "secret_key")
//	print(h)
func hmacSha256(L *lua.LState) int {
	message := L.CheckString(1)
	key := L.CheckString(2)

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	h := mac.Sum(nil)

	L.Push(lua.LString(hex.EncodeToString(h)))
	return 1
}

// tableToJSON: converts a Lua table to JSON string
func tableToJSON(L *lua.LState, table lua.LValue) ([]byte, error) {
	translator := glua.NewTranslator()
	var goValue interface{}
	err := translator.FromLua(L, table, &goValue)
	if err != nil {
		return nil, err
	}
	return json.Marshal(goValue)
}

// md5HashObj: computes the MD5 hash of a Lua table (converted to JSON).
//
// @luafunc md5_obj
// @luaparam obj table The table to hash
// @luareturn string hash The hex-encoded MD5 hash
// @luareturn string|nil err Error message if conversion fails
//
// Example:
//
//	local h, err = hash.md5_obj({name="John", age=30})
//	if err then error(err) end
//	print(h)
func md5HashObj(L *lua.LState) int {
	table := L.CheckAny(1)
	jsonBytes, err := tableToJSON(L, table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	h := md5.Sum(jsonBytes)
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	L.Push(lua.LNil)
	return 2
}

// sha1HashObj: computes the SHA1 hash of a Lua table (converted to JSON).
//
// @luafunc sha1_obj
// @luaparam obj table The table to hash
// @luareturn string hash The hex-encoded SHA1 hash
// @luareturn string|nil err Error message if conversion fails
//
// Example:
//
//	local h, err = hash.sha1_obj({name="John", age=30})
//	if err then error(err) end
//	print(h)
func sha1HashObj(L *lua.LState) int {
	table := L.CheckAny(1)
	jsonBytes, err := tableToJSON(L, table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	h := sha1.Sum(jsonBytes)
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	L.Push(lua.LNil)
	return 2
}

// sha256HashObj: computes the SHA256 hash of a Lua table (converted to JSON).
//
// @luafunc sha256_obj
// @luaparam obj table The table to hash
// @luareturn string hash The hex-encoded SHA256 hash
// @luareturn string|nil err Error message if conversion fails
//
// Example:
//
//	local h, err = hash.sha256_obj({name="John", age=30})
//	if err then error(err) end
//	print(h)
func sha256HashObj(L *lua.LState) int {
	table := L.CheckAny(1)
	jsonBytes, err := tableToJSON(L, table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	h := sha256.Sum256(jsonBytes)
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	L.Push(lua.LNil)
	return 2
}

// sha512HashObj: computes the SHA512 hash of a Lua table (converted to JSON).
//
// @luafunc sha512_obj
// @luaparam obj table The table to hash
// @luareturn string hash The hex-encoded SHA512 hash
// @luareturn string|nil err Error message if conversion fails
//
// Example:
//
//	local h, err = hash.sha512_obj({name="John", age=30})
//	if err then error(err) end
//	print(h)
func sha512HashObj(L *lua.LState) int {
	table := L.CheckAny(1)
	jsonBytes, err := tableToJSON(L, table)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	h := sha512.Sum512(jsonBytes)
	L.Push(lua.LString(hex.EncodeToString(h[:])))
	L.Push(lua.LNil)
	return 2
}
