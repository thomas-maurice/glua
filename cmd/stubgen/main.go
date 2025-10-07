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

	"github.com/thomas-maurice/glua/pkg/stubgen"
)

func main() {
	var (
		dir       = flag.String("dir", ".", "Directory to scan for Go module files")
		output    = flag.String("output", "module_stubs.gen.lua", "Output file for generated Lua stubs")
		outputDir = flag.String("output-dir", "", "Output directory for per-module stub files (recommended for LSP)")
	)

	flag.Parse()

	analyzer := stubgen.NewAnalyzer()

	if err := analyzer.ScanDirectory(*dir); err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	// If output-dir is specified, generate per-module files
	if *outputDir != "" {
		if err := os.MkdirAll(*outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			os.Exit(1)
		}

		for moduleName := range analyzer.GetModules() {
			stub, err := analyzer.GenerateModuleStub(moduleName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating stub for %s: %v\n", moduleName, err)
				os.Exit(1)
			}

			outputFile := fmt.Sprintf("%s/%s.gen.lua", *outputDir, moduleName)
			if err := os.WriteFile(outputFile, []byte(stub), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outputFile, err)
				os.Exit(1)
			}
			fmt.Printf("Generated %s\n", outputFile)
		}

		fmt.Printf("Generated Lua stubs for %d module(s) in %s/\n", analyzer.ModuleCount(), *outputDir)
		return
	}

	// Otherwise, generate combined file
	stubs, err := analyzer.GenerateStubs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating stubs: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(*output, []byte(stubs), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated Lua stubs for %d module(s) in %s\n", analyzer.ModuleCount(), *output)
}
