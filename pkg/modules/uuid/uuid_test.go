// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package uuid

import (
	"path/filepath"
	"sort"
	"testing"
	"time"

	googleuuid "github.com/google/uuid"
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
			L.PreloadModule("uuid", Loader)
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

// TestV4_VersionAndVariant: a v4 UUID must report version 4 and the RFC4122
// variant — this is what distinguishes it from a malformed or non-UUID
// string that happens to be 36 characters long.
func TestV4_VersionAndVariant(t *testing.T) {
	id, err := v4()
	require.NoError(t, err)

	info, err := parse(id)
	require.NoError(t, err)
	assert.Equal(t, 4, info.Version)
	assert.Equal(t, "RFC4122", info.Variant)
	assert.Equal(t, float64(0), info.Timestamp, "v4 has no embedded timestamp")
}

// TestV4_NoCollisions: many generations must never repeat. This is the
// property that makes v4 usable as a unique identifier at all.
func TestV4_NoCollisions(t *testing.T) {
	const n = 5000
	seen := make(map[string]bool, n)
	for i := 0; i < n; i++ {
		id, err := v4()
		require.NoError(t, err)
		require.False(t, seen[id], "v4 collision at iteration %d: %s", i, id)
		seen[id] = true
	}
}

// TestV7_VersionAndTimestamp: a v7 UUID must report version 7 and a
// timestamp within a few seconds of "now" — this is v7's whole point: the
// UUID itself carries a meaningful creation time.
func TestV7_VersionAndTimestamp(t *testing.T) {
	before := float64(time.Now().Unix())
	id, err := v7()
	require.NoError(t, err)
	after := float64(time.Now().Unix())

	info, err := parse(id)
	require.NoError(t, err)
	assert.Equal(t, 7, info.Version)
	assert.GreaterOrEqual(t, info.Timestamp, before-1)
	assert.LessOrEqual(t, info.Timestamp, after+1)
}

// TestV7_Monotonic: minting 1000 v7 UUIDs in a tight loop must produce a
// non-decreasing string sequence. This is the test that justifies taking
// the google/uuid dependency at all — a hand-rolled v7 without the
// monotonic counter would fail exactly this test under a fast loop where
// multiple UUIDs land in the same millisecond.
func TestV7_Monotonic(t *testing.T) {
	const n = 1000
	ids := make([]string, n)
	for i := range ids {
		id, err := v7()
		require.NoError(t, err)
		ids[i] = id
	}

	assert.True(t, sort.StringsAreSorted(ids), "v7 UUIDs must sort in generation order")
}

// TestV5_Deterministic: identical inputs must produce identical output,
// twice, and it must be version 5.
func TestV5_Deterministic(t *testing.T) {
	a, err := v5(googleuuid.NameSpaceDNS.String(), "example.com")
	require.NoError(t, err)
	b, err := v5(googleuuid.NameSpaceDNS.String(), "example.com")
	require.NoError(t, err)
	assert.Equal(t, a, b)

	info, err := parse(a)
	require.NoError(t, err)
	assert.Equal(t, 5, info.Version)
}

// TestV5_InvalidNamespace_Raises: a malformed namespace must raise rather
// than silently hashing the garbage string.
func TestV5_InvalidNamespace_Raises(t *testing.T) {
	_, err := v5("not-a-uuid", "example.com")
	assert.Error(t, err)
}

// TestParse_AllAcceptedForms: parse must accept canonical, plain (no
// hyphens), urn:uuid:, and {braced} forms of the same UUID, and report the
// same canonical value for all of them.
func TestParse_AllAcceptedForms(t *testing.T) {
	id, err := v4()
	require.NoError(t, err)

	plain := id[0:8] + id[9:13] + id[14:18] + id[19:23] + id[24:]
	forms := map[string]string{
		"canonical": id,
		"plain":     plain,
		"urn":       "urn:uuid:" + id,
		"braced":    "{" + id + "}",
	}

	for name, s := range forms {
		info, err := parse(s)
		require.NoError(t, err, "form %s should parse", name)
		assert.Equal(t, id, info.UUID, "form %s should report the canonical uuid", name)
	}
}

// TestParse_Invalid_Raises: garbage input must raise with a useful error,
// not panic or return a zero-value Info silently.
func TestParse_Invalid_Raises(t *testing.T) {
	_, err := parse("definitely-not-a-uuid")
	assert.Error(t, err)
}

// TestIsValid_MatchesParse: is_valid must return false exactly where parse
// raises, and true exactly where it succeeds — the two must never
// disagree, or a caller who checks is_valid before calling parse could
// still see it raise.
func TestIsValid_MatchesParse(t *testing.T) {
	id, err := v4()
	require.NoError(t, err)
	assert.True(t, isValid(id))

	_, err = parse(id)
	assert.NoError(t, err)

	assert.False(t, isValid("garbage"))
	_, err = parse("garbage")
	assert.Error(t, err)
}

// TestFormat_AllStyles_RoundTrip: formatting into every style and parsing
// it back must report the same canonical UUID.
func TestFormat_AllStyles_RoundTrip(t *testing.T) {
	id, err := v4()
	require.NoError(t, err)

	for _, style := range []string{"canonical", "plain", "urn", "braced"} {
		formatted, err := format(id, style)
		require.NoError(t, err, "style %s", style)

		info, err := parse(formatted)
		require.NoError(t, err, "re-parsing style %s", style)
		assert.Equal(t, id, info.UUID, "style %s did not round-trip", style)
	}
}

// TestFormat_UnknownStyle_Raises: an unrecognised style must raise a
// useful error rather than silently falling back to canonical.
func TestFormat_UnknownStyle_Raises(t *testing.T) {
	id, err := v4()
	require.NoError(t, err)

	_, err = format(id, "nope")
	assert.Error(t, err)
}

// TestFormat_InvalidUUID_Raises: format must validate its input before
// touching style at all.
func TestFormat_InvalidUUID_Raises(t *testing.T) {
	_, err := format("garbage", "canonical")
	assert.Error(t, err)
}
