// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hash

import (
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test files in testdata directory
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	if err != nil {
		t.Fatalf("Failed to glob test files: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No test files found in testdata/")
	}

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
