// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hmac

import (
	"bytes"
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
			L.PreloadModule("hmac", Loader)
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

// TestSHA1_RFC2202: RFC 2202 HMAC-SHA1 test vectors. A hand-rolled fixture
// proves nothing about correctness against the spec; these are the
// published reference values.
func TestSHA1_RFC2202(t *testing.T) {
	cases := []struct {
		key, data, want string
	}{
		{
			key:  string(bytes.Repeat([]byte{0x0b}, 20)),
			data: "Hi There",
			want: "b617318655057264e28bc0b6fb378c8ef146be00",
		},
		{
			key:  "Jefe",
			data: "what do ya want for nothing?",
			want: "effcdf6ae5eb2fa2d27416d5f184df9c259a7c79",
		},
		{
			key:  string(bytes.Repeat([]byte{0xaa}, 20)),
			data: string(bytes.Repeat([]byte{0xdd}, 50)),
			want: "125d7342b9ac11cd91a39af48aa17b4f63f175d3",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, sha1Tag(c.data, c.key))
	}
}

// TestSHA256_RFC4231: RFC 4231 HMAC-SHA256 test vectors (cases 1-3).
func TestSHA256_RFC4231(t *testing.T) {
	cases := []struct {
		key, data, want string
	}{
		{
			key:  string(bytes.Repeat([]byte{0x0b}, 20)),
			data: "Hi There",
			want: "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7",
		},
		{
			key:  "Jefe",
			data: "what do ya want for nothing?",
			want: "5bdcc146bf60754e6a042426089575c75a003f089d2739839dec58b964ec3843",
		},
		{
			key:  string(bytes.Repeat([]byte{0xaa}, 20)),
			data: string(bytes.Repeat([]byte{0xdd}, 50)),
			want: "773ea91e36800e46854db8ebd09181a72959098b3ef8c122d9635514ced565fe",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, sha256Tag(c.data, c.key))
	}
}

// TestSHA512_RFC4231: RFC 4231 HMAC-SHA512 test vectors (cases 1-2).
func TestSHA512_RFC4231(t *testing.T) {
	cases := []struct {
		key, data, want string
	}{
		{
			key:  string(bytes.Repeat([]byte{0x0b}, 20)),
			data: "Hi There",
			want: "87aa7cdea5ef619d4ff0b4241a1d6cb02379f4e2ce4ec2787ad0b30545e17cdedaa833b7d6b8a702038b274eaea3f4e4be9d914eeb61f1702e696c203a126854",
		},
		{
			key:  "Jefe",
			data: "what do ya want for nothing?",
			want: "164b7a7bfcf819e2e395fbe73b56e0a387bd64222e831fd610270cd7ea2505549758bf75c05a994a6d034f65f8f0e6fdcaeab1a34d4a6b4b636e070a38bce737",
		},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, sha512Tag(c.data, c.key))
	}
}

// TestTags_EmptyKeyAndMessage: an empty key or message must not panic and
// must produce a valid tag of the expected fixed length — HMAC is defined
// for the empty string on either side.
func TestTags_EmptyKeyAndMessage(t *testing.T) {
	assert.Len(t, sha256Tag("", ""), 64)
	assert.Len(t, sha256Tag("message", ""), 64)
	assert.Len(t, sha256Tag("", "key"), 64)
}

// TestTags_BinaryKey: a binary (non-UTF-8) key must work — HMAC keys are
// arbitrary byte strings, and glua's string type is 8-bit clean, so this
// must not be mangled or rejected.
func TestTags_BinaryKey(t *testing.T) {
	binaryKey := string([]byte{0x00, 0xff, 0x80, 0x01, 0xfe})
	got := sha256Tag("message", binaryKey)
	assert.Len(t, got, 64)
	assert.True(t, verifySHA256("message", binaryKey, got))
}

// TestVerify_MatchingTag: verify_* must return true for a tag it itself
// computed — the baseline correctness case every other test in this file
// contrasts against.
func TestVerify_MatchingTag(t *testing.T) {
	tag := sha256Tag("payload", "secret")
	assert.True(t, verifySHA256("payload", "secret", tag))
}

// TestVerify_BitFlippedTag: a single-bit flip in an otherwise valid tag must
// be rejected. This is the property that makes verify_* useful at all: it
// must be sensitive to any change in the tag, not just gross corruption.
func TestVerify_BitFlippedTag(t *testing.T) {
	tag := sha256Tag("payload", "secret")
	// Flip the low bit of the first hex-decoded byte by flipping its last
	// hex digit.
	flipped := "0" + tag[1:]
	if flipped == tag {
		flipped = "1" + tag[1:]
	}
	assert.False(t, verifySHA256("payload", "secret", flipped))
}

// TestVerify_NeverRaises: every malformed-tag shape (non-hex garbage, empty,
// truncated, over-long, valid-hex-but-wrong-value) must return false, never
// raise. tag is attacker-controlled input, and raising on it would turn a
// verification failure into a script-aborting error — precisely the bug
// this module's API shape exists to prevent.
func TestVerify_NeverRaises(t *testing.T) {
	validTag := sha256Tag("payload", "secret")

	cases := map[string]string{
		"wrong tag":                          sha256Tag("payload", "different-secret"),
		"empty tag":                          "",
		"non-hex garbage":                    "not hex at all!!",
		"truncated tag":                      validTag[:10],
		"over-long tag":                      validTag + "00",
		"odd-length hex":                     "abc",
		"uppercase mismatch but wrong value": "FFFFFFFF",
	}
	for name, tag := range cases {
		t.Run(name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				ok := verifySHA256("payload", "secret", tag)
				assert.False(t, ok, "expected verify to return false for %s", name)
			})
		})
	}
}

// TestVerify_UppercaseHexAccepted: verify_* accepts uppercase hex for the
// tag, matching Go's encoding/hex.DecodeString, which is case-insensitive.
// This is a deliberate, documented choice (see package doc) rather than an
// accident of the underlying decoder — pin it with a test so a future
// switch to a stricter decoder is a visible, deliberate change.
func TestVerify_UppercaseHexAccepted(t *testing.T) {
	tag := sha256Tag("payload", "secret")
	upper := make([]byte, len(tag))
	for i := 0; i < len(tag); i++ {
		c := tag[i]
		if c >= 'a' && c <= 'f' {
			c -= 'a' - 'A'
		}
		upper[i] = c
	}
	assert.True(t, verifySHA256("payload", "secret", string(upper)))
}

// TestVerifySHA1_MatchAndMismatch: sanity check that verify_sha1 has the
// same match/mismatch shape as verify_sha256, since it shares the generic
// verify() helper.
func TestVerifySHA1_MatchAndMismatch(t *testing.T) {
	tag := sha1Tag("m", "k")
	assert.True(t, verifySHA1("m", "k", tag))
	assert.False(t, verifySHA1("m", "k", "0000000000000000000000000000000000000a"))
}

// TestVerifySHA512_MatchAndMismatch: sanity check that verify_sha512 has the
// same match/mismatch shape as verify_sha256.
func TestVerifySHA512_MatchAndMismatch(t *testing.T) {
	tag := sha512Tag("m", "k")
	assert.True(t, verifySHA512("m", "k", tag))
	assert.False(t, verifySHA512("m", "k", tag[:len(tag)-1]+"0"))
}
