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

// Package random provides crypto/rand-backed randomness for Lua scripts:
// raw bytes, hex, URL-safe tokens, charset-drawn strings and uniform
// integers.
//
// This is a CSPRNG, not gopher-lua's math.random. math.random is a seeded,
// deterministic PRNG — its output is predictable from the seed and it must
// never be used for secrets, tokens, keys or anything security-sensitive.
// Everything in this module is backed by crypto/rand, suitable for API
// keys, salts, nonces and similar. There is deliberately no seeding
// function here, so random.* is never reproducible; tests that need
// deterministic output must inject their own values rather than seed this
// module.
//
// No function in this module writes with modulo bias: every draw from a
// range uses crypto/rand.Int against a big.Int span, which rejection-samples
// internally, rather than `% n` on a fixed-width random value.
package random

import (
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// maxLen: the largest n accepted by bytes/hex/token/string. A generous but
// finite bound — large enough for any realistic token/salt/string use, small
// enough that a typo'd argument (e.g. a byte count meant to be a bit count)
// cannot turn into an unbounded allocation.
const maxLen = 1 << 20

// maxSafeSpan: the largest magnitude a random.int span (max - min) may have
// and still guarantee every representable result is an exact float64
// integer (2^53 is the largest integer float64 represents exactly).
const maxSafeSpan = 1 << 53

// alphanum: the value of the ALPHANUM module constant.
const alphanum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// readRandom: reads n cryptographically random bytes, validating n against
// [1, maxLen]. fnName prefixes any error for greppability back to the
// call site.
func readRandom(fnName string, n int) ([]byte, error) {
	if n < 1 || n > maxLen {
		return nil, fmt.Errorf("%s: n must be in [1, %d], got %d", fnName, maxLen, n)
	}
	b := make([]byte, n)
	// crypto/rand.Read cannot fail on any platform Go modern supports, but
	// the error is still checked and raised rather than ignored.
	if _, err := crand.Read(b); err != nil {
		return nil, fmt.Errorf("%s: %w", fnName, err)
	}
	return b, nil
}

// randomBytes: returns n cryptographically random bytes as a raw string.
func randomBytes(n int) ([]byte, error) {
	return readRandom("random.bytes", n)
}

// randomHex: returns n random bytes, hex-encoded (2n hex characters).
func randomHex(n int) (string, error) {
	b, err := readRandom("random.hex", n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// randomToken: returns n random bytes, base64url-encoded without padding —
// safe to use directly in a URL, header or filename, unlike standard
// base64's '+', '/' and '=' characters.
func randomToken(n int) (string, error) {
	b, err := readRandom("random.token", n)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// randomString: returns a string of n characters drawn uniformly (and
// bias-free) from charset. charset is treated as a sequence of runes, not
// bytes, so a UTF-8 charset cannot produce broken UTF-8 output. Duplicate
// runes in charset are not de-duplicated: a repeated character is simply
// more likely to be drawn, which is the caller's responsibility to avoid if
// unweighted output is required.
func randomString(n int, charset string) (string, error) {
	if n < 1 || n > maxLen {
		return "", fmt.Errorf("random.string: n must be in [1, %d], got %d", maxLen, n)
	}
	runes := []rune(charset)
	if len(runes) == 0 {
		return "", fmt.Errorf("random.string: charset must not be empty")
	}

	span := big.NewInt(int64(len(runes)))
	out := make([]rune, n)
	for i := range out {
		idx, err := crand.Int(crand.Reader, span)
		if err != nil {
			return "", fmt.Errorf("random.string: %w", err)
		}
		out[i] = runes[idx.Int64()]
	}
	return string(out), nil
}

// randInt: returns a uniformly distributed, bias-free random integer in
// [min, max] — inclusive of BOTH ends, matching Lua's math.random(m, n)
// rather than Go's half-open range convention. A Go reader should not
// assume half-open semantics here: this is a deliberate deviation to stay
// consistent with the target language.
//
// min and max arrive as float64 (Lua has no separate integer type), so
// non-integral input is rejected explicitly rather than silently truncated.
func randInt(lo, hi float64) (float64, error) {
	if lo != math.Trunc(lo) || hi != math.Trunc(hi) {
		return 0, fmt.Errorf("random.int: min and max must be integral, got min=%v max=%v", lo, hi)
	}
	if lo > hi {
		return 0, fmt.Errorf("random.int: min (%v) must be <= max (%v)", lo, hi)
	}
	span := hi - lo
	if span >= maxSafeSpan {
		return 0, fmt.Errorf("random.int: max - min must be < 2^53, got %v", span)
	}

	// span+1 possible values in [min, max] inclusive; crypto/rand.Int
	// rejection-samples internally, so this draw carries no modulo bias.
	n, err := crand.Int(crand.Reader, big.NewInt(int64(span)+1))
	if err != nil {
		return 0, fmt.Errorf("random.int: %w", err)
	}
	return lo + float64(n.Int64()), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("random", "crypto/rand-backed randomness: bytes, hex, tokens, strings and integers")

	m.Fn("bytes", randomBytes, "returns n cryptographically random bytes as a raw 8-bit-clean string",
		luareg.Args("n"),
		luareg.ArgDoc("n", fmt.Sprintf("number of bytes to generate; must be in [1, %d]", maxLen)),
		luareg.ReturnDoc(0, "data", "n random bytes"))

	m.Fn("hex", randomHex, "returns n random bytes, hex-encoded",
		luareg.Args("n"),
		luareg.ArgDoc("n", fmt.Sprintf("number of underlying random bytes; must be in [1, %d]", maxLen)),
		luareg.ReturnDoc(0, "hex", "lowercase hex encoding, 2n characters long"))

	m.Fn("token", randomToken, "returns n random bytes as an unpadded, URL-safe base64 token",
		luareg.Args("n"),
		luareg.ArgDoc("n", fmt.Sprintf("number of underlying random bytes; must be in [1, %d]", maxLen)),
		luareg.ReturnDoc(0, "token", "base64url encoding without padding, safe to use directly in a URL or header"))

	m.Fn("string", randomString, "returns a string of n characters drawn uniformly from charset",
		luareg.Args("n", "charset"),
		luareg.ArgDoc("n", fmt.Sprintf("number of characters to generate; must be in [1, %d]", maxLen)),
		luareg.ArgDoc("charset", "the runes to draw from, as a UTF-8 string; must not be empty. Duplicate runes are not de-duplicated and are therefore weighted"),
		luareg.ReturnDoc(0, "s", "n runes drawn from charset"))

	m.Fn("int", randInt, "returns a uniform, bias-free random integer in [min, max], inclusive of both ends (like Lua's math.random, NOT Go's half-open convention)",
		luareg.Args("min", "max"),
		luareg.ArgDoc("min", "inclusive lower bound; must be an integral value"),
		luareg.ArgDoc("max", "inclusive upper bound; must be an integral value, >= min, with max - min < 2^53"),
		luareg.ReturnDoc(0, "n", "a random integer in [min, max]"))

	m.Const("ALPHANUM", alphanum, "string", "uppercase+lowercase letters and digits, for use as a random.string charset")

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("random", random.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
