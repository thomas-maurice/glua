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

// glua-gen regenerates the library/*.gen.lua stubs for every module shipped
// with glua. It is the CLI used by `make gen-stubs` and by CI to keep the
// checked-in stubs in sync with the Go source.
//
// Downstream projects that embed glua and add their own modules should NOT
// invoke this binary. Instead they should ship their own tools/stubgen/main.go
// that imports their modules + glua's modules package and calls
// stubgen.GenerateFromRegistry directly. See README.md for the template.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/thomas-maurice/glua/pkg/luareg"
	"github.com/thomas-maurice/glua/pkg/modules"
	"github.com/thomas-maurice/glua/pkg/stubgen"
)

// main: parses flags, registers every glua module, and writes stubs to outDir.
func main() {
	outDir := flag.String("out", "library", "output directory for .gen.lua files")
	flag.Parse()

	reg := luareg.NewRegistry()
	modules.RegisterAll(reg)

	gen := stubgen.NewGenerator()
	files, err := gen.GenerateFromRegistry(reg, *outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "glua-gen: %v\n", err)
		os.Exit(1)
	}

	for _, f := range files {
		fmt.Println(f)
	}
	fmt.Fprintf(os.Stderr, "glua-gen: wrote %d file(s) to %s\n", len(files), *outDir)
}
