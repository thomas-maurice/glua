// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package base64

import (
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
			L.PreloadModule("base64", Loader)
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

// TestDecodeRaisesOnInvalid: decode raises a Lua error on invalid base64.
func TestDecodeRaisesOnInvalid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local ok, err = pcall(base64.decode, "not!valid!base64!")
		assert(not ok, "Expected error for invalid base64")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}

// TestEncodeDecodeRoundTrip: encode then decode recovers the original string.
func TestEncodeDecodeRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local s = "The quick brown fox jumps over the lazy dog"
		local rt = base64.decode(base64.encode(s))
		assert(rt == s, "Round-trip failed, got: " .. rt)
	`
	require.NoError(t, L.DoString(code))
}

// TestURLRoundTrip: URL-safe encode/decode round-trips correctly.
func TestURLRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local s = "test data with /+= chars"
		local rt = base64.decode_url(base64.encode_url(s))
		assert(rt == s, "URL round-trip failed")
	`
	require.NoError(t, L.DoString(code))
}
