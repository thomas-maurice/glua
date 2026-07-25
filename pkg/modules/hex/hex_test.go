// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hex

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
			L.PreloadModule("hex", Loader)
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

// TestDecodeRaisesOnInvalid: decode raises a Lua error on invalid input.
func TestDecodeRaisesOnInvalid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local ok, err = pcall(hex.decode, "zzzz")
		assert(not ok, "Expected error for invalid hex")
	`
	require.NoError(t, L.DoString(code))
}

// TestEncodeDecodeRoundTrip: encode followed by decode returns the original.
func TestEncodeDecodeRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local s = "hello world"
		local rt = hex.decode(hex.encode(s))
		assert(rt == s, "Round-trip failed, got: " .. rt)
	`
	require.NoError(t, L.DoString(code))
}
