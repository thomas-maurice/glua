// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package template

import (
	"os"
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestRender_Simple(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render("Hello {{.Name}}", {Name = "World"})
		assert(err == nil, "Expected no error")
		assert(result == "Hello World", "Expected 'Hello World', got: " .. result)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRender_Multiple(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render("{{.First}} {{.Last}}", {First = "John", Last = "Doe"})
		assert(err == nil, "Expected no error")
		assert(result == "John Doe", "Expected 'John Doe', got: " .. result)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRender_Array(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render("{{range .Items}}{{.}} {{end}}", {Items = {"a", "b", "c"}})
		assert(err == nil, "Expected no error")
		assert(result == "a b c ", "Expected 'a b c ', got: " .. result)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRender_InvalidTemplate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render("{{.Missing", {Name = "Test"})
		assert(result == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRenderFile(t *testing.T) {
	// Create temporary template file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.tmpl")
	content := []byte("Hello {{.Name}}")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render_file("` + tmpFile + `", {Name = "File"})
		assert(err == nil, "Expected no error: " .. tostring(err))
		assert(result == "Hello File", "Expected 'Hello File', got: " .. result)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRenderFile_NotFound(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("template", Loader)

	code := `
		local template = require("template")
		local result, err = template.render_file("/nonexistent/file.tmpl", {Name = "Test"})
		assert(result == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
