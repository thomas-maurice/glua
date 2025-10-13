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
	"fmt"
	"os"
	"strings"

	"github.com/thomas-maurice/glua/pkg/glua"
)

// Generator: combines function stub generation (from annotations) and type stub generation (from TypeRegistry).
// This provides a simple API for modules to generate complete Lua stubs.
type Generator struct {
	analyzer     *Analyzer
	typeRegistry *glua.TypeRegistry
}

// NewGenerator: creates a new Generator instance
func NewGenerator() *Generator {
	return &Generator{
		analyzer:     NewAnalyzer(),
		typeRegistry: glua.NewTypeRegistry(),
	}
}

// ScanDirectory: scans a directory for Go files with @luafunc annotations
func (g *Generator) ScanDirectory(dir string) error {
	return g.analyzer.ScanDirectory(dir)
}

// RegisterType: registers a Go type for Lua stub generation
func (g *Generator) RegisterType(obj interface{}) error {
	return g.typeRegistry.Register(obj)
}

// ProcessTypes: processes all registered types to discover dependencies
func (g *Generator) ProcessTypes() error {
	return g.typeRegistry.Process()
}

// GenerateModule: generates complete Lua stubs for a module (functions + types merged into one file).
// Returns the combined stub content or an error.
func (g *Generator) GenerateModule(moduleName string) (string, error) {
	// Generate type stubs first
	typeStubs, err := g.typeRegistry.GenerateStubs()
	if err != nil {
		return "", fmt.Errorf("failed to generate type stubs: %w", err)
	}

	// Generate function stubs
	functionStubs, err := g.analyzer.GenerateModuleStub(moduleName)
	if err != nil {
		return "", fmt.Errorf("failed to generate function stubs: %w", err)
	}

	// Merge: remove "return {}" from type stubs, then append function stubs
	typeStubs = strings.TrimSuffix(strings.TrimSpace(typeStubs), "return {}")
	typeStubs = strings.TrimSpace(typeStubs)

	// Combine with a blank line separator
	var sb strings.Builder
	if typeStubs != "" {
		sb.WriteString(typeStubs)
		sb.WriteString("\n\n")
	}
	sb.WriteString(functionStubs)

	return sb.String(), nil
}

// GenerateConfig: configuration for stub generation
type GenerateConfig struct {
	// ScanDir is the directory to scan for Go files with @luafunc annotations
	ScanDir string
	// OutputDir is the output directory for generated Lua files
	OutputDir string
	// ModuleName is the name of the module for Lua code generation
	ModuleName string
	// OutputFile is the name of the generated Lua file (e.g., "k8sclient.gen.lua")
	OutputFile string
	// Types is an optional list of types to register for stub generation
	Types []interface{}
}

// Generate: generates Lua stubs for a module based on the provided configuration.
// This method orchestrates the entire stub generation process:
// 1. Scans the specified directory for annotated functions
// 2. Registers any provided types
// 3. Processes type dependencies
// 4. Generates combined stubs
// 5. Writes output to the specified file
//
// Returns the path to the generated file or an error.
func (g *Generator) Generate(config GenerateConfig) (string, error) {
	// Scan directory for function annotations
	if err := g.ScanDirectory(config.ScanDir); err != nil {
		return "", fmt.Errorf("error scanning directory: %w", err)
	}

	// Register types if provided
	for _, t := range config.Types {
		if err := g.RegisterType(t); err != nil {
			return "", fmt.Errorf("error registering types: %w", err)
		}
	}

	// Process types to discover dependencies
	if err := g.ProcessTypes(); err != nil {
		return "", fmt.Errorf("error processing types: %w", err)
	}

	// Generate combined stub
	stubs, err := g.GenerateModule(config.ModuleName)
	if err != nil {
		return "", fmt.Errorf("error generating stubs: %w", err)
	}

	// Write output
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("error creating output directory: %w", err)
	}

	outputPath := fmt.Sprintf("%s/%s", config.OutputDir, config.OutputFile)
	if err := os.WriteFile(outputPath, []byte(stubs), 0644); err != nil {
		return "", fmt.Errorf("error writing output: %w", err)
	}

	return outputPath, nil
}
