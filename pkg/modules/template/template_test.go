// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package template

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestRender_Simple(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result = template.render("Hello {{.Name}}", {Name = "World"})
		assert(result == "Hello World", "Expected 'Hello World', got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}

func TestRender_Multiple(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result = template.render("{{.First}} {{.Last}}", {First = "John", Last = "Doe"})
		assert(result == "John Doe", "Expected 'John Doe', got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}

func TestRender_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result = template.render("{{range .Items}}{{.}} {{end}}", {Items = {"a", "b", "c"}})
		assert(result == "a b c ", "Expected 'a b c ', got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}

// TestRender_InvalidTemplate: invalid template syntax raises a Lua error.
func TestRender_InvalidTemplate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local ok, err = pcall(template.render, "{{.Missing", {Name = "Test"})
		assert(not ok, "Expected error for invalid template")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}

// TestRender_CyclicData_ViaLuaPcall: a self-referential data table raises a
// catchable error instead of crashing the process, and the Lua state keeps
// working afterwards.
//
// Regression test: template.render/render_file used to hand-roll their own
// unguarded recursive Lua->Go conversion (luaTableToGoMap/luaValueToGo), so
// handing them a cyclic data table (`t.self = t`) recursed forever and
// killed the whole host process with an unrecoverable "fatal error: stack
// overflow" -- not something pcall could ever catch. They now share
// pkg/glua's TableGuard, so the same input raises an ordinary, catchable
// Lua error.
func TestRender_CyclicData_ViaLuaPcall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")

		local t = {}
		t.self = t

		local ok, err = pcall(template.render, "{{.Name}}", t)
		assert(not ok, "Expected error for a self-referential data table")
		assert(string.find(err, "cycle") ~= nil, "Expected error to mention 'cycle', got: " .. err)

		-- Prove the Lua state is still usable after the caught error.
		local result = template.render("Hello {{.Name}}", {Name = "World"})
		assert(result == "Hello World", "Expected normal render to still work, got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}

func TestRenderFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.tmpl")
	if err := os.WriteFile(tmpFile, []byte("Hello {{.Name}}"), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result = template.render_file("` + tmpFile + `", {Name = "File"})
		assert(result == "Hello File", "Expected 'Hello File', got: " .. result)
	`
	require.NoError(t, L.DoString(code))
}

// TestRenderFile_NotFound: a missing template file raises a Lua error.
func TestRenderFile_NotFound(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local ok, err = pcall(template.render_file, "/nonexistent/file.tmpl", {Name = "Test"})
		assert(not ok, "Expected error for missing file")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}
