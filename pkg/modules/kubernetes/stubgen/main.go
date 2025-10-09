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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	"github.com/thomas-maurice/glua/pkg/stubgen"
)

func main() {
	outputDir := flag.String("output", "library", "Output directory for generated stubs")
	flag.Parse()

	// Create generator
	gen := stubgen.NewGenerator()

	// Scan module directory for function annotations
	// Get the directory where this source file lives
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error determining source directory\n")
		os.Exit(1)
	}
	moduleDir := filepath.Dir(filepath.Dir(filename))

	if err := gen.ScanDirectory(moduleDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	// Register types used by this module
	if err := gen.RegisterType(kubernetes.GVKMatcher{}); err != nil {
		fmt.Fprintf(os.Stderr, "Error registering types: %v\n", err)
		os.Exit(1)
	}

	// Process types to discover dependencies
	if err := gen.ProcessTypes(); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing types: %v\n", err)
		os.Exit(1)
	}

	// Generate combined stub
	stubs, err := gen.GenerateModule("kubernetes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating stubs: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	outputFile := filepath.Join(*outputDir, "kubernetes.gen.lua")
	if err := os.WriteFile(outputFile, []byte(stubs), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outputFile)
}
