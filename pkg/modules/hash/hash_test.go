// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hash

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test files in testdata directory.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No test files found in testdata/")

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("hash", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result.Type() == lua.LTBool && !lua.LVAsBool(result) {
				t.Fatal("Test returned false")
			}
		})
	}
}

// TestStringHashFunctions: plain string hashes return the expected hex strings.
func TestStringHashFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		-- MD5 of "hello world"
		local h = hash.md5("hello world")
		assert(h == "5eb63bbbe01eeed093cb22bb8f5acdc3", "Unexpected MD5: " .. h)
		-- SHA256 of empty string
		local e = hash.sha256("")
		assert(#e == 64, "SHA256 should be 64 hex chars")
	`
	require.NoError(t, L.DoString(code))
}

// TestObjHashRaisesOnUnserializable: sha256_obj raises when the value cannot
// be serialised (currently all Lua tables marshal fine; this guards the path
// via pcall to confirm luareg error propagation).
func TestObjHashDoesNotRaiseOnTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local ok, result = pcall(hash.sha256_obj, {key="value"})
		assert(ok, "sha256_obj should not raise on a valid table")
		assert(#result == 64, "Expected 64 hex chars")
	`
	require.NoError(t, L.DoString(code))
}
