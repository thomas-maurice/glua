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
