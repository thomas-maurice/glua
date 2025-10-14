// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package fs

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello World"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local content, err = fs.read_file("` + tmpFile + `")
		assert(err == nil, "Expected no error")
		assert(content == "Hello World", "Expected correct content")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "output.txt")

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local err = fs.write_file("` + tmpFile + `", "Test Content")
		assert(err == nil, "Expected no error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}

	// Verify file was written
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}
	if string(content) != "Test Content" {
		t.Errorf("Expected 'Test Content', got '%s'", string(content))
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		assert(fs.exists("` + existingFile + `"), "File should exist")
		assert(not fs.exists("` + filepath.Join(tmpDir, "nonexistent.txt") + `"), "File should not exist")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestMkdir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir")

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local err = fs.mkdir("` + newDir + `")
		assert(err == nil, "Expected no error")
		assert(fs.exists("` + newDir + `"), "Directory should exist")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "path", "to", "nested")

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local err = fs.mkdir_all("` + nestedDir + `")
		assert(err == nil, "Expected no error")
		assert(fs.exists("` + nestedDir + `"), "Nested directory should exist")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRemove(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "remove.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		assert(fs.exists("` + tmpFile + `"), "File should exist before removal")
		local err = fs.remove("` + tmpFile + `")
		assert(err == nil, "Expected no error")
		assert(not fs.exists("` + tmpFile + `"), "File should not exist after removal")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRemoveAll(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "parent", "child")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}
	testFile := filepath.Join(nestedDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	parentDir := filepath.Join(tmpDir, "parent")
	code := `
		local fs = require("fs")
		assert(fs.exists("` + parentDir + `"), "Parent dir should exist")
		local err = fs.remove_all("` + parentDir + `")
		assert(err == nil, "Expected no error")
		assert(not fs.exists("` + parentDir + `"), "Parent dir should not exist after removal")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local entries, err = fs.list("` + tmpDir + `")
		assert(err == nil, "Expected no error")
		assert(#entries == 3, "Expected 3 entries, got " .. #entries)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestStat(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "stat.txt")
	content := "Hello World"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local info, err = fs.stat("` + tmpFile + `")
		assert(err == nil, "Expected no error")
		assert(info.name == "stat.txt", "Expected correct name")
		assert(info.size == 11, "Expected size 11, got " .. info.size)
		assert(not info.is_dir, "Should not be a directory")
		assert(info.mode > 0, "Mode should be set")
		assert(info.mod_time > 0, "Modification time should be set")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestStatDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local info, err = fs.stat("` + tmpDir + `")
		assert(err == nil, "Expected no error")
		assert(info.is_dir, "Should be a directory")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "roundtrip.txt")

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local original = "Test Content\nWith Multiple Lines\n"

		-- Write file
		local err = fs.write_file("` + tmpFile + `", original)
		assert(err == nil, "Write should succeed")

		-- Read file back
		local content, err = fs.read_file("` + tmpFile + `")
		assert(err == nil, "Read should succeed")
		assert(content == original, "Content should match")

		-- Cleanup
		err = fs.remove("` + tmpFile + `")
		assert(err == nil, "Remove should succeed")
		assert(not fs.exists("` + tmpFile + `"), "File should be gone")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
