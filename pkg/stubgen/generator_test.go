// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package stubgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewGenerator: tests that NewGenerator creates a valid generator instance
func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if gen.analyzer == nil {
		t.Error("Generator has nil analyzer")
	}
	if gen.typeRegistry == nil {
		t.Error("Generator has nil typeRegistry")
	}
}

// TestScanDirectory: tests that ScanDirectory works correctly
func TestScanDirectory(t *testing.T) {
	gen := NewGenerator()

	// Use testdata directory which should exist
	err := gen.ScanDirectory("testdata")
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	// Verify that modules were found
	if gen.analyzer.ModuleCount() == 0 {
		t.Error("No modules discovered in testdata")
	}
}

// TestRegisterType: tests that RegisterType works correctly
func TestRegisterType(t *testing.T) {
	gen := NewGenerator()

	type TestStruct struct {
		Field1 string
		Field2 int
	}

	err := gen.RegisterType(TestStruct{})
	if err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}
}

// TestProcessTypes: tests that ProcessTypes works correctly
func TestProcessTypes(t *testing.T) {
	gen := NewGenerator()

	type TestStruct struct {
		Field1 string
		Field2 int
	}

	if err := gen.RegisterType(TestStruct{}); err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}

	if err := gen.ProcessTypes(); err != nil {
		t.Fatalf("ProcessTypes failed: %v", err)
	}
}

// TestGenerateModule: tests that GenerateModule produces valid output
func TestGenerateModule(t *testing.T) {
	gen := NewGenerator()

	// Scan testdata to get some modules
	if err := gen.ScanDirectory("testdata"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	// Get the first module
	modules := gen.analyzer.GetModules()
	if len(modules) == 0 {
		t.Skip("No modules found in testdata")
	}

	var moduleName string
	for name := range modules {
		moduleName = name
		break
	}

	output, err := gen.GenerateModule(moduleName)
	if err != nil {
		t.Fatalf("GenerateModule failed: %v", err)
	}

	if output == "" {
		t.Error("GenerateModule returned empty output")
	}

	// Verify output contains module name
	if !strings.Contains(output, moduleName) {
		t.Errorf("Generated output doesn't contain module name %q", moduleName)
	}
}

// TestGenerate: tests the full Generate workflow
func TestGenerate(t *testing.T) {
	gen := NewGenerator()

	// Create temporary output directory
	tmpDir := t.TempDir()

	// First scan to find an actual module
	if err := gen.ScanDirectory("testdata"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	modules := gen.analyzer.GetModules()
	if len(modules) == 0 {
		t.Skip("No modules found in testdata")
	}

	var moduleName string
	for name := range modules {
		moduleName = name
		break
	}

	// Reset generator to test full workflow
	gen = NewGenerator()

	config := GenerateConfig{
		ScanDir:    "testdata",
		OutputDir:  tmpDir,
		ModuleName: moduleName,
		OutputFile: moduleName + ".gen.lua",
		Types:      nil,
	}

	outputPath, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputPath)
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Output file is empty")
	}

	// Verify it's valid Lua syntax
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---@meta") {
		t.Error("Output doesn't start with Lua meta annotation")
	}
}

// TestGenerateWithTypes: tests Generate with type registration
func TestGenerateWithTypes(t *testing.T) {
	gen := NewGenerator()

	type TestStruct struct {
		Field1 string
		Field2 int
	}

	tmpDir := t.TempDir()

	// First scan to find an actual module
	if err := gen.ScanDirectory("testdata"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	modules := gen.analyzer.GetModules()
	if len(modules) == 0 {
		t.Skip("No modules found in testdata")
	}

	var moduleName string
	for name := range modules {
		moduleName = name
		break
	}

	// Reset generator
	gen = NewGenerator()

	config := GenerateConfig{
		ScanDir:    "testdata",
		OutputDir:  tmpDir,
		ModuleName: moduleName,
		OutputFile: "test_with_types.gen.lua",
		Types:      []interface{}{TestStruct{}},
	}

	outputPath, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("Generate with types failed: %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputPath)
	}
}

// TestGenerateInvalidDir: tests Generate with invalid directory
func TestGenerateInvalidDir(t *testing.T) {
	gen := NewGenerator()

	config := GenerateConfig{
		ScanDir:    "/nonexistent/directory",
		OutputDir:  t.TempDir(),
		ModuleName: "test",
		OutputFile: "test.gen.lua",
		Types:      nil,
	}

	_, err := gen.Generate(config)
	if err == nil {
		t.Error("Expected error for invalid scan directory, got nil")
	}
}

// TestGenerateInvalidModule: tests GenerateModule with non-existent module
func TestGenerateInvalidModule(t *testing.T) {
	gen := NewGenerator()

	_, err := gen.GenerateModule("nonexistent_module")
	if err == nil {
		t.Error("Expected error for non-existent module, got nil")
	}
}

// TestGenerateOutputCreation: tests that Generate creates output directory if needed
func TestGenerateOutputCreation(t *testing.T) {
	gen := NewGenerator()

	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "nested", "output", "dir")

	// First scan to find an actual module
	if err := gen.ScanDirectory("testdata"); err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	modules := gen.analyzer.GetModules()
	if len(modules) == 0 {
		t.Skip("No modules found in testdata")
	}

	var moduleName string
	for name := range modules {
		moduleName = name
		break
	}

	// Reset generator
	gen = NewGenerator()

	config := GenerateConfig{
		ScanDir:    "testdata",
		OutputDir:  outputDir,
		ModuleName: moduleName,
		OutputFile: "test.gen.lua",
		Types:      nil,
	}

	outputPath, err := gen.Generate(config)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify nested directory was created
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		t.Error("Output directory was not created")
	}

	// Verify output file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file %s was not created", outputPath)
	}
}
