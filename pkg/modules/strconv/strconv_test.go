// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package strconv

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/ directory.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("strconv", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("Test script returned %v, expected true", result)
			}
		})
	}
}

// TestFormatIntInvalidBasePanicsGuarded: proves formatInt raises (does not
// panic) for a base outside [2, 36], where Go's own strconv.FormatInt would
// panic and crash the whole process.
func TestFormatIntInvalidBasePanicsGuarded(t *testing.T) {
	_, err := formatInt(420, 40)
	require.Error(t, err)
	_, err = formatInt(420, 1)
	require.Error(t, err)
}

// TestFormatFloatUnknownFmtGuarded: proves formatFloat raises for an unknown
// format verb, where Go's own strconv.FormatFloat silently returns a bogus
// "%!z(...)"-shaped string instead of erroring.
func TestFormatFloatUnknownFmtGuarded(t *testing.T) {
	_, err := formatFloat(1.5, "z", -1)
	require.Error(t, err)
	_, err = formatFloat(1.5, "ff", -1)
	require.Error(t, err)
}

// TestParseIntOversizedRaises: an integer whose magnitude exceeds 2^53
// raises rather than silently rounding, because it cannot be represented
// exactly as a Lua number (float64).
func TestParseIntOversizedRaises(t *testing.T) {
	// 2^53 exactly is still representable.
	v, err := parseInt("9007199254740992", 10)
	require.NoError(t, err)
	require.Equal(t, float64(9007199254740992), v)

	// 2^53 + 1 is not.
	_, err = parseInt("9007199254740993", 10)
	require.Error(t, err)
}

// TestFormatIntOversizedRaises: mirrors TestParseIntOversizedRaises on the
// formatting side. Uses 2^54 rather than 2^53+1: a float64 literal that
// close to 2^53 would round to 2^53 itself before the function ever saw it.
func TestFormatIntOversizedRaises(t *testing.T) {
	_, err := formatInt(9007199254740992, 10) // exactly 2^53: still safe
	require.NoError(t, err)

	_, err = formatInt(math.Pow(2, 54), 10)
	require.Error(t, err)
}

// TestUnquoteRoundTrip: quote followed by unquote returns the original
// string, including embedded quotes, newlines, and a non-UTF-8 byte — the
// exact cases the S2 spec calls out.
func TestUnquoteRoundTrip(t *testing.T) {
	cases := []string{
		"",
		`hello "world"`,
		"line one\nline two",
		string([]byte{0xff, 0x00, 0x41}),
	}
	for _, s := range cases {
		q := quote(s)
		rt, err := unquote(q)
		require.NoError(t, err)
		require.Equal(t, s, rt)
	}
}
