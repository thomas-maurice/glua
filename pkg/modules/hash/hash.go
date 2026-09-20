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

// Package hash provides cryptographic hash utilities for Lua scripts: MD5,
// SHA1, SHA256, SHA512 and HMAC-SHA256, over strings or arbitrary Lua values
// serialised to JSON.
//
// hmac_sha256 is deprecated in favour of the hmac module
// (pkg/modules/hmac), which additionally provides HMAC-SHA1/SHA512 and,
// more importantly, verify_* functions that compare tags in constant time
// inside Go. hmac_sha256 is kept here byte-identical in behaviour so
// existing scripts do not break; it will be removed at the next deliberate
// breaking release. Do not build new HMAC verification against == on this
// function's result — see the hmac package doc for why that is a timing
// oracle.
package hash

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// md5Hash: computes the hex-encoded MD5 hash of a string.
func md5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// sha1Hash: computes the hex-encoded SHA1 hash of a string.
func sha1Hash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// sha256Hash: computes the hex-encoded SHA256 hash of a string.
func sha256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// sha512Hash: computes the hex-encoded SHA512 hash of a string.
func sha512Hash(s string) string {
	h := sha512.Sum512([]byte(s))
	return hex.EncodeToString(h[:])
}

// hmacSHA256: computes the hex-encoded HMAC-SHA256 of message with key.
func hmacSHA256(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// tableToJSON: converts a Lua value to its JSON representation for hashing.
func tableToJSON(L *lua.LState, value lua.LValue) ([]byte, error) {
	translator := glua.NewTranslator()
	var goValue interface{}
	if err := translator.FromLua(L, value, &goValue); err != nil {
		return nil, err
	}
	return json.Marshal(goValue)
}

// md5HashObj: computes the MD5 hash of a Lua value serialised to JSON.
// Uses *lua.LState escape hatch to accept any Lua value.
func md5HashObj(L *lua.LState, tbl lua.LValue) (string, error) {
	b, err := tableToJSON(L, tbl)
	if err != nil {
		return "", fmt.Errorf("failed to serialise value: %w", err)
	}
	h := md5.Sum(b)
	return hex.EncodeToString(h[:]), nil
}

// sha1HashObj: computes the SHA1 hash of a Lua value serialised to JSON.
func sha1HashObj(L *lua.LState, tbl lua.LValue) (string, error) {
	b, err := tableToJSON(L, tbl)
	if err != nil {
		return "", fmt.Errorf("failed to serialise value: %w", err)
	}
	h := sha1.Sum(b)
	return hex.EncodeToString(h[:]), nil
}

// sha256HashObj: computes the SHA256 hash of a Lua value serialised to JSON.
func sha256HashObj(L *lua.LState, tbl lua.LValue) (string, error) {
	b, err := tableToJSON(L, tbl)
	if err != nil {
		return "", fmt.Errorf("failed to serialise value: %w", err)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// sha512HashObj: computes the SHA512 hash of a Lua value serialised to JSON.
func sha512HashObj(L *lua.LState, tbl lua.LValue) (string, error) {
	b, err := tableToJSON(L, tbl)
	if err != nil {
		return "", fmt.Errorf("failed to serialise value: %w", err)
	}
	h := sha512.Sum512(b)
	return hex.EncodeToString(h[:]), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("hash", "cryptographic hash utilities")
	m.Fn("md5", md5Hash, "computes the hex-encoded MD5 hash of a string",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to hash"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded MD5 digest"))
	m.Fn("sha1", sha1Hash, "computes the hex-encoded SHA1 hash of a string",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to hash"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA1 digest"))
	m.Fn("sha256", sha256Hash, "computes the hex-encoded SHA256 hash of a string",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to hash"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA256 digest"))
	m.Fn("sha512", sha512Hash, "computes the hex-encoded SHA512 hash of a string",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to hash"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA512 digest"))
	m.Fn("hmac_sha256", hmacSHA256, "computes the hex-encoded HMAC-SHA256 of a message with a key (deprecated: use hmac.sha256; and never compare its result with == against attacker-supplied input, see the hmac module)",
		luareg.Args("message", "key"),
		luareg.ArgDoc("message", "the message to authenticate"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded HMAC-SHA256 tag"))
	m.Fn("md5_obj", md5HashObj, "computes the MD5 hash of a Lua value serialised to JSON",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the Lua value to hash; it is JSON-marshalled before hashing"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded MD5 digest of the JSON encoding"))
	m.Fn("sha1_obj", sha1HashObj, "computes the SHA1 hash of a Lua value serialised to JSON",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the Lua value to hash; it is JSON-marshalled before hashing"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA1 digest of the JSON encoding"))
	m.Fn("sha256_obj", sha256HashObj, "computes the SHA256 hash of a Lua value serialised to JSON",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the Lua value to hash; it is JSON-marshalled before hashing"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA256 digest of the JSON encoding"))
	m.Fn("sha512_obj", sha512HashObj, "computes the SHA512 hash of a Lua value serialised to JSON",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the Lua value to hash; it is JSON-marshalled before hashing"),
		luareg.ReturnDoc(0, "hex", "the lowercase hex-encoded SHA512 digest of the JSON encoding"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("hash", hash.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
