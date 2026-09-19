// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package password

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// testCost: bcrypt's minimum allowed cost, used throughout this file to
// keep the suite fast. bcrypt's own runtime is exponential in cost — cost
// 12+ takes on the order of hundreds of milliseconds per call, which would
// make this package's tests crawl if run at a production-realistic cost.
// Only TestHash_DefaultCost below exercises DEFAULT_COST, and it is skipped
// under -short.
const testCost = 4

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
			L.PreloadModule("password", Loader)
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

// TestHash_ProducesBcryptFormat: hash output must start with the bcrypt
// "$2a$<cost>$" prefix at the requested cost, so callers can rely on the
// standard bcrypt storage format.
func TestHash_ProducesBcryptFormat(t *testing.T) {
	h, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(h, "$2a$04$"), "expected $2a$04$ prefix, got %q", h)
}

// TestHash_RandomSalt: two hashes of the same password must differ, because
// bcrypt draws a fresh random salt each call. A test that only checked
// verify would not catch a regression that reused a fixed salt (which would
// make two identical passwords produce identical stored hashes, leaking
// that fact to anyone with database access).
func TestHash_RandomSalt(t *testing.T) {
	h1, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)
	h2, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)
	assert.NotEqual(t, h1, h2, "expected distinct hashes for the same password due to random salt")
}

// TestVerify_CorrectAndWrongPassword: verify must return true for the
// correct password and false — not an error — for a wrong one, since the
// plaintext is attacker-controlled and a wrong guess is a normal outcome.
func TestVerify_CorrectAndWrongPassword(t *testing.T) {
	h, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)

	ok, err := verifyPassword("hunter2", h)
	require.NoError(t, err)
	assert.True(t, ok)

	ok, err = verifyPassword("wrong-password", h)
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestVerify_EmptyAttempt: an empty attempt against a non-empty hash must
// return false, not raise or panic — it is still just "the wrong password".
func TestVerify_EmptyAttempt(t *testing.T) {
	h, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)

	ok, err := verifyPassword("", h)
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestVerify_RaisesOnMalformedHash: hash is the application's OWN stored
// data, not attacker input, so a malformed value must raise — this is the
// deliberate mirror image of hmac.verify_* (pkg/modules/hmac), which never
// raises because its tag argument IS attacker input. A malformed hash here
// means a storage or migration bug, and silently treating it as "wrong
// password" would hide that bug forever.
func TestVerify_RaisesOnMalformedHash(t *testing.T) {
	_, err := verifyPassword("hunter2", "not-a-hash")
	assert.Error(t, err)

	_, err = verifyPassword("hunter2", "$2a$04$tooshort")
	assert.Error(t, err)
}

// TestCost_ReturnsEmbeddedCost: cost(hash) must return the cost the hash was
// created with — this is the primitive rehash-on-login logic depends on.
func TestCost_ReturnsEmbeddedCost(t *testing.T) {
	h, err := hashPassword("hunter2", testCost)
	require.NoError(t, err)

	c, err := costOf(h)
	require.NoError(t, err)
	assert.Equal(t, testCost, c)
}

// TestCost_RaisesOnMalformedHash: same reasoning as TestVerify_RaisesOnMalformedHash.
func TestCost_RaisesOnMalformedHash(t *testing.T) {
	_, err := costOf("garbage")
	assert.Error(t, err)
}

// TestHash_InvalidCostRaises: cost outside bcrypt's [4, 31] range must raise
// with a message useful enough to debug — x/crypto/bcrypt itself silently
// clamps a too-low cost up to its own internal default instead of erroring,
// which would otherwise let an out-of-range cost through unnoticed.
func TestHash_InvalidCostRaises(t *testing.T) {
	_, err := hashPassword("hunter2", 3)
	assert.Error(t, err, "cost below bcrypt.MinCost must raise")

	_, err = hashPassword("hunter2", 32)
	assert.Error(t, err, "cost above bcrypt.MaxCost must raise")
}

// TestHash_OverLengthPasswordRaises: bcrypt silently truncates input beyond
// 72 bytes internally, which means two distinct passwords sharing the first
// 72 bytes would otherwise verify identically. This module deliberately
// raises instead of truncating (see package doc) — pin that decision with a
// test at the boundary.
func TestHash_OverLengthPasswordRaises(t *testing.T) {
	ok72 := strings.Repeat("a", 72)
	_, err := hashPassword(ok72, testCost)
	assert.NoError(t, err, "a 72-byte password must succeed (boundary)")

	tooLong := strings.Repeat("a", 73)
	_, err = hashPassword(tooLong, testCost)
	assert.Error(t, err, "a 73-byte password must raise")
}

// TestHash_TruncationWouldOtherwiseCollide: demonstrates directly why the
// over-length rejection matters. Two 73-byte passwords sharing the same
// first 72 bytes would, if we allowed bcrypt's silent truncation through,
// verify against each other. We instead confirm both are rejected at
// hash-time so that collision can never be reached.
func TestHash_TruncationWouldOtherwiseCollide(t *testing.T) {
	prefix := strings.Repeat("a", 72)
	_, err1 := hashPassword(prefix+"1", testCost)
	_, err2 := hashPassword(prefix+"2", testCost)
	assert.Error(t, err1)
	assert.Error(t, err2)
}

// TestHash_DefaultCost: exercises the real, production-recommended
// DEFAULT_COST end to end. Skipped under -short because bcrypt at cost 12
// takes roughly a quarter of a second per call.
func TestHash_DefaultCost(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping DEFAULT_COST bcrypt round trip in -short mode")
	}
	h, err := hashPassword("hunter2", defaultCost)
	require.NoError(t, err)

	ok, err := verifyPassword("hunter2", h)
	require.NoError(t, err)
	assert.True(t, ok)

	c, err := costOf(h)
	require.NoError(t, err)
	assert.Equal(t, defaultCost, c)
}
