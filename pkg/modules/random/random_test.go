// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package random

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("random", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("test script returned %v, expected true", result)
			}
		})
	}
}

// TestBytes_LengthAndDistinct: bytes(n) returns exactly n bytes, and two
// calls do not return the same value — this is the CSPRNG property the
// whole module exists to guarantee (math.random would also pass a length
// check, but is predictable from its seed).
func TestBytes_LengthAndDistinct(t *testing.T) {
	a, err := randomBytes(32)
	require.NoError(t, err)
	assert.Len(t, a, 32)

	b, err := randomBytes(32)
	require.NoError(t, err)
	assert.NotEqual(t, a, b, "two calls to random.bytes must not collide")
}

// TestBytes_InvalidN_Raises: n <= 0 and n beyond the cap both raise.
func TestBytes_InvalidN_Raises(t *testing.T) {
	for _, n := range []int{0, -1, maxLen + 1} {
		_, err := randomBytes(n)
		assert.Error(t, err, "expected error for n=%d", n)
	}
}

// TestHex_LengthAndAlphabet: hex(n) returns 2n lowercase hex characters.
func TestHex_LengthAndAlphabet(t *testing.T) {
	s, err := randomHex(16)
	require.NoError(t, err)
	require.Len(t, s, 32)
	for _, c := range s {
		assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'),
			"unexpected character %q in hex output", c)
	}
}

// TestToken_URLSafeUnpadded: token(n) must never contain standard base64's
// '+', '/' or '=' — those are exactly the characters that make plain
// base64 unsafe to drop into a URL or header, which is the entire reason
// this function exists instead of telling callers to use base64.encode.
func TestToken_URLSafeUnpadded(t *testing.T) {
	s, err := randomToken(32)
	require.NoError(t, err)
	assert.NotContains(t, s, "+")
	assert.NotContains(t, s, "/")
	assert.NotContains(t, s, "=")
}

// TestString_OnlyFromCharset_MultiByte: every rune of the output must come
// from charset, including when charset contains multi-byte UTF-8 runes —
// proving the module is rune-oriented, not byte-oriented, so it cannot
// produce broken UTF-8.
func TestString_OnlyFromCharset_MultiByte(t *testing.T) {
	charset := "aXbYcZ日本語"
	runes := []rune(charset)
	allowed := make(map[rune]bool, len(runes))
	for _, r := range runes {
		allowed[r] = true
	}

	s, err := randomString(200, charset)
	require.NoError(t, err)
	for _, r := range s {
		assert.True(t, allowed[r], "unexpected rune %q not in charset", r)
	}
}

// TestString_EmptyCharset_Raises: an empty charset cannot produce any
// output, so it must raise rather than silently return an empty string
// (which would be indistinguishable from n=0, itself already rejected).
func TestString_EmptyCharset_Raises(t *testing.T) {
	_, err := randomString(4, "")
	assert.Error(t, err)
}

// TestString_InvalidN_Raises: n <= 0 and n beyond the cap both raise.
func TestString_InvalidN_Raises(t *testing.T) {
	for _, n := range []int{0, -1, maxLen + 1} {
		_, err := randomString(n, "abc")
		assert.Error(t, err, "expected error for n=%d", n)
	}
}

// TestInt_RangeAndEndpoints: over many samples, random.int(min, max) never
// escapes [min, max] and hits both endpoints. This is a distribution
// smoke test, not a statistical bias proof — it would only reliably catch
// a badly broken implementation (off-by-one, wrong inclusivity), not a
// subtle bias.
func TestInt_RangeAndEndpoints(t *testing.T) {
	const min, max = 1, 6
	seenMin, seenMax := false, false
	for i := 0; i < 10000; i++ {
		n, err := randInt(min, max)
		require.NoError(t, err)
		require.GreaterOrEqual(t, n, float64(min))
		require.LessOrEqual(t, n, float64(max))
		if n == min {
			seenMin = true
		}
		if n == max {
			seenMax = true
		}
	}
	assert.True(t, seenMin, "min endpoint was never drawn in 10000 samples")
	assert.True(t, seenMax, "max endpoint was never drawn in 10000 samples")
}

// TestInt_MinEqualsMax: the single-value edge case always returns that
// value.
func TestInt_MinEqualsMax(t *testing.T) {
	n, err := randInt(42, 42)
	require.NoError(t, err)
	assert.Equal(t, float64(42), n)
}

// TestInt_MinGreaterThanMax_Raises: an inverted range raises rather than
// silently swapping or returning a nonsensical value.
func TestInt_MinGreaterThanMax_Raises(t *testing.T) {
	_, err := randInt(5, 1)
	assert.Error(t, err)
}

// TestInt_NonIntegral_Raises: min/max are float64 (Lua has no separate
// integer type); a fractional bound must raise rather than being silently
// truncated, since the caller's intent (which integer did they mean?) is
// ambiguous.
func TestInt_NonIntegral_Raises(t *testing.T) {
	_, err := randInt(1.5, 6)
	assert.Error(t, err, "expected non-integral min to raise")

	_, err = randInt(1, 6.5)
	assert.Error(t, err, "expected non-integral max to raise")
}

// TestInt_SpanTooLarge_Raises: max - min >= 2^53 would overflow float64's
// exact integer range, so it must raise rather than silently return an
// imprecise result.
func TestInt_SpanTooLarge_Raises(t *testing.T) {
	_, err := randInt(0, maxSafeSpan)
	assert.Error(t, err)

	// One below the boundary must succeed.
	_, err = randInt(0, maxSafeSpan-1)
	assert.NoError(t, err)
}
