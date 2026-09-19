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

// Package hmac provides HMAC tag computation and verification for Lua
// scripts, over SHA1, SHA256 and SHA512.
//
// The central design point: a script that does
// `if hash.hmac_sha256(msg, key) == tag then ... end` has written a timing
// oracle without realising it. Go's `==` on strings is an early-exit
// comparison, so the time it takes leaks how many leading bytes of the
// computed tag matched the attacker-supplied one, and a patient attacker can
// recover a valid tag byte by byte over repeated requests.
//
// The fix here is not documentation, it is API shape: there is deliberately
// NO way to obtain a comparison primitive from this module. There is no
// hmac.equal and no hmac.constant_time_compare. Verification is done
// entirely inside Go via crypto/hmac.Equal, which compares in constant time,
// and the only thing Lua ever sees is the boolean result. If a reviewer sees
// a comparison primitive added to this module later, that is a regression of
// the whole point of the module.
//
// Verification is also per-algorithm (verify_sha1/verify_sha256/verify_sha512)
// rather than a single verify(message, key, tag, algorithm) taking an
// algorithm name. A string-typed algorithm argument would let the algorithm
// arrive from the same attacker-controlled data the tag comes from (e.g. a
// JSON body with {"alg": "...), which is exactly the class of bug this
// module exists to design out.
//
// Tags are lowercase hex, matching hash.hmac_sha256's existing encoding.
// verify_* accepts both lowercase and uppercase hex for the tag argument
// (encoding/hex.DecodeString is case-insensitive) — this is simplest and
// avoids a surprising rejection of an otherwise-valid tag; it is not a
// security-relevant choice since hex casing carries no secret information.
//
// verify_* NEVER raises, not even on a malformed tag: tag is
// attacker-controlled input (e.g. a request header), so a non-hex or
// wrong-length tag is a normal "verification failed" outcome (false), not a
// script-aborting error — raising here would turn an authentication failure
// into a 500. Contrast this with password.verify (pkg/modules/password),
// which DOES raise on a malformed stored hash, because that hash is the
// application's own data and a malformed one means the database or a
// migration is broken, not that the user typed the wrong password.
//
// hash.hmac_sha256 (pkg/modules/hash) remains unchanged and is not removed —
// removing a published function would break every script already using it.
// It is doc-marked deprecated in favour of hmac.sha256.
package hmac

import (
	"crypto/hmac"
	// SHA1 is offered for HMAC interop with legacy systems, not for new
	// designs. HMAC-SHA1 remains unbroken as a MAC even though SHA1's
	// collision resistance is broken (HMAC only relies on second-preimage
	// resistance of the compression function, not collision resistance).
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// tag: computes the hex-encoded HMAC tag of message under key using the
// given hash constructor.
func tag(newHash func() hash.Hash, message, key string) string {
	mac := hmac.New(newHash, []byte(key))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// verify: decodes wantTag as hex and compares it in constant time against
// the HMAC of message under key. Never raises: a decode failure (bad hex,
// wrong length) is treated as "does not match" rather than an error, because
// wantTag is attacker-controlled. The hex decode itself is not constant
// time, but it only leaks whether the tag is valid hex and its length, both
// of which are public information carried by the request itself, not secret
// material.
func verify(newHash func() hash.Hash, message, key, wantTag string) bool {
	decoded, err := hex.DecodeString(wantTag)
	if err != nil {
		return false
	}
	mac := hmac.New(newHash, []byte(key))
	mac.Write([]byte(message))
	return hmac.Equal(mac.Sum(nil), decoded)
}

// sha1Tag: computes the hex-encoded HMAC-SHA1 tag of message under key.
func sha1Tag(message, key string) string {
	return tag(sha1.New, message, key)
}

// sha256Tag: computes the hex-encoded HMAC-SHA256 tag of message under key.
func sha256Tag(message, key string) string {
	return tag(sha256.New, message, key)
}

// sha512Tag: computes the hex-encoded HMAC-SHA512 tag of message under key.
func sha512Tag(message, key string) string {
	return tag(sha512.New, message, key)
}

// verifySHA1: reports whether tagHex is the correct HMAC-SHA1 tag for
// message under key. Never raises; a malformed tagHex returns false.
func verifySHA1(message, key, tagHex string) bool {
	return verify(sha1.New, message, key, tagHex)
}

// verifySHA256: reports whether tagHex is the correct HMAC-SHA256 tag for
// message under key. Never raises; a malformed tagHex returns false.
func verifySHA256(message, key, tagHex string) bool {
	return verify(sha256.New, message, key, tagHex)
}

// verifySHA512: reports whether tagHex is the correct HMAC-SHA512 tag for
// message under key. Never raises; a malformed tagHex returns false.
func verifySHA512(message, key, tagHex string) bool {
	return verify(sha512.New, message, key, tagHex)
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("hmac", "HMAC tag computation and constant-time verification")

	m.Fn("sha1", sha1Tag, "computes the hex-encoded HMAC-SHA1 tag of a message with a key",
		luareg.Args("message", "key"),
		luareg.ArgDoc("message", "the message to authenticate"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ReturnDoc(0, "tag", "the lowercase hex-encoded HMAC-SHA1 tag"))
	m.Fn("sha256", sha256Tag, "computes the hex-encoded HMAC-SHA256 tag of a message with a key",
		luareg.Args("message", "key"),
		luareg.ArgDoc("message", "the message to authenticate"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ReturnDoc(0, "tag", "the lowercase hex-encoded HMAC-SHA256 tag"))
	m.Fn("sha512", sha512Tag, "computes the hex-encoded HMAC-SHA512 tag of a message with a key",
		luareg.Args("message", "key"),
		luareg.ArgDoc("message", "the message to authenticate"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ReturnDoc(0, "tag", "the lowercase hex-encoded HMAC-SHA512 tag"))

	m.Fn("verify_sha1", verifySHA1,
		"verifies an HMAC-SHA1 tag in constant time; never raises, a malformed tag simply returns false",
		luareg.Args("message", "key", "tag"),
		luareg.ArgDoc("message", "the message that was authenticated"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ArgDoc("tag", "the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false"),
		luareg.ReturnDoc(0, "ok", "true if tag is the correct HMAC-SHA1 tag for message under key"))
	m.Fn("verify_sha256", verifySHA256,
		"verifies an HMAC-SHA256 tag in constant time; never raises, a malformed tag simply returns false",
		luareg.Args("message", "key", "tag"),
		luareg.ArgDoc("message", "the message that was authenticated"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ArgDoc("tag", "the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false"),
		luareg.ReturnDoc(0, "ok", "true if tag is the correct HMAC-SHA256 tag for message under key"))
	m.Fn("verify_sha512", verifySHA512,
		"verifies an HMAC-SHA512 tag in constant time; never raises, a malformed tag simply returns false",
		luareg.Args("message", "key", "tag"),
		luareg.ArgDoc("message", "the message that was authenticated"),
		luareg.ArgDoc("key", "the shared secret key"),
		luareg.ArgDoc("tag", "the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false"),
		luareg.ReturnDoc(0, "ok", "true if tag is the correct HMAC-SHA512 tag for message under key"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("hmac", hmac.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
