// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thomas-maurice/glua/pkg/stubgen"
)

func main() {
	outputDir := flag.String("output", "library", "Output directory for generated stubs")
	flag.Parse()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error determining source directory\n")
		os.Exit(1)
	}
	moduleDir := filepath.Dir(filepath.Dir(filename))

	gen := stubgen.NewGenerator()
	outputFile, err := gen.Generate(stubgen.GenerateConfig{
		ScanDir:    moduleDir,
		OutputDir:  *outputDir,
		ModuleName: "hash",
		OutputFile: "hash.gen.lua",
		Types:      nil,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", outputFile)
}
