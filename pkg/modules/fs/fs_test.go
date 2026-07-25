// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	require.NoError(t, os.WriteFile(tmpFile, []byte("Hello World"), 0644))

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local content = fs.read_file("` + tmpFile + `")
		assert(content == "Hello World", "Expected correct content")
	`
	require.NoError(t, L.DoString(code))
}

// TestReadFileRaises: reading a nonexistent file raises a Lua error.
func TestReadFileRaises(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local ok, err = pcall(fs.read_file, "/nonexistent/file.txt")
		assert(not ok, "Expected error for missing file")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "output.txt")

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		fs.write_file("` + tmpFile + `", "Test Content")
	`
	require.NoError(t, L.DoString(code))

	content, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Equal(t, "Test Content", string(content))
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	require.NoError(t, os.WriteFile(existingFile, []byte("test"), 0644))

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		assert(fs.exists("` + existingFile + `"), "File should exist")
		assert(not fs.exists("` + filepath.Join(tmpDir, "nonexistent.txt") + `"), "File should not exist")
	`
	require.NoError(t, L.DoString(code))
}

func TestMkdir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir")

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		fs.mkdir("` + newDir + `")
		assert(fs.exists("` + newDir + `"), "Directory should exist")
	`
	require.NoError(t, L.DoString(code))
}

func TestMkdirAll(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "path", "to", "nested")

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		fs.mkdir_all("` + nestedDir + `")
		assert(fs.exists("` + nestedDir + `"), "Nested directory should exist")
	`
	require.NoError(t, L.DoString(code))
}

func TestRemove(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "remove.txt")
	require.NoError(t, os.WriteFile(tmpFile, []byte("test"), 0644))

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		assert(fs.exists("` + tmpFile + `"), "File should exist before removal")
		fs.remove("` + tmpFile + `")
		assert(not fs.exists("` + tmpFile + `"), "File should not exist after removal")
	`
	require.NoError(t, L.DoString(code))
}

func TestRemoveAll(t *testing.T) {
	tmpDir := t.TempDir()
	nestedDir := filepath.Join(tmpDir, "parent", "child")
	require.NoError(t, os.MkdirAll(nestedDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nestedDir, "file.txt"), []byte("test"), 0644))

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	parentDir := filepath.Join(tmpDir, "parent")
	code := `
		local fs = require("fs")
		assert(fs.exists("` + parentDir + `"), "Parent dir should exist")
		fs.remove_all("` + parentDir + `")
		assert(not fs.exists("` + parentDir + `"), "Parent dir should not exist after removal")
	`
	require.NoError(t, L.DoString(code))
}

func TestList(t *testing.T) {
	tmpDir := t.TempDir()
	for _, f := range []string{"file1.txt", "file2.txt", "file3.txt"} {
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644))
	}

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local entries = fs.list("` + tmpDir + `")
		assert(#entries == 3, "Expected 3 entries, got " .. #entries)
	`
	require.NoError(t, L.DoString(code))
}

func TestStat(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "stat.txt")
	require.NoError(t, os.WriteFile(tmpFile, []byte("Hello World"), 0644))

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local info = fs.stat("` + tmpFile + `")
		assert(info.name == "stat.txt", "Expected correct name")
		assert(info.size == 11, "Expected size 11, got " .. info.size)
		assert(not info.is_dir, "Should not be a directory")
		assert(info.mode > 0, "Mode should be set")
		assert(info.mod_time > 0, "Modification time should be set")
	`
	require.NoError(t, L.DoString(code))
}

// TestStatRaises: stat on a nonexistent path raises a Lua error.
func TestStatRaises(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("fs", Loader)

	code := `
		local fs = require("fs")
		local ok, err = pcall(fs.stat, "/nonexistent/path")
		assert(not ok, "Expected error for nonexistent path")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
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
		fs.write_file("` + tmpFile + `", original)
		local content = fs.read_file("` + tmpFile + `")
		assert(content == original, "Content should match")
		fs.remove("` + tmpFile + `")
		assert(not fs.exists("` + tmpFile + `"), "File should be gone")
	`
	require.NoError(t, L.DoString(code))
}
