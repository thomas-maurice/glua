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
	"fmt"
	"os"

	logmodule "github.com/thomas-maurice/glua/pkg/modules/log"
	lua "github.com/yuin/gopher-lua"
)

func main() {
	fmt.Println("=== Log Module Examples ===")

	// Example 1: Module-level functions (simplest)
	fmt.Println("=== Example 1: Module-level logging (simplest) ===")
	runModuleLevelExample()

	fmt.Println("\n=== Example 2: Logger object ===")
	runLoggerObjectExample()

	fmt.Println("\n=== Example 3: Logger object with context (logger:with) ===")
	runLoggerWithExample()

	fmt.Println("\n=== Example 4: Injecting fields from Go (advanced) ===")
	runFieldInjectionExample()

	fmt.Println("\n=== Done! ===")
	fmt.Println("\nTry running with different formats:")
	fmt.Println("  LOG_FORMAT=json go run ./example/log   # JSON format")
	fmt.Println("  go run ./example/log                    # Text format (default)")
}

// runModuleLevelExample: demonstrates the simplest usage - module-level functions
func runModuleLevelExample() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("log", logmodule.Loader)

	if err := L.DoString(`
		local log = require("log")
		log.info("Application started")
		log.warn("This is a warning")
		log.error("An error occurred")
	`); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runLoggerObjectExample: demonstrates using logger objects
func runLoggerObjectExample() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("log", logmodule.Loader)

	if err := L.DoString(`
		local log = require("log")
		local logger = log.logger()
		logger:info("Using logger object")
		logger:warn("Warning via object")
	`); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runLoggerWithExample: demonstrates logger:with() for adding context
func runLoggerWithExample() {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("log", logmodule.Loader)

	if err := L.DoString(`
		local log = require("log")
		local logger = log.logger()

		-- Create logger with context
		local contextLogger = logger:with("session_id", "sess-789", "tenant", "acme-corp")
		contextLogger:info("Processing request")
		contextLogger:info("Request completed", "duration_ms", 150)

		-- Original logger unchanged
		logger:info("Other action without context")

		-- Can chain with() calls
		local userLogger = contextLogger:with("user_id", 12345)
		userLogger:info("User action")  -- has session_id, tenant, AND user_id
	`); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// runFieldInjectionExample: demonstrates injecting fields from Go
func runFieldInjectionExample() {
	L := lua.NewState()
	defer L.Close()

	// ADVANCED: Inject fields from Go before running Lua
	// Use GetDefaultLogger() to get the configured default logger, then add fields
	logger := logmodule.GetDefaultLogger().With(
		"app", "log-example",
		"version", "1.0.0",
		"request_id", "abc-123-xyz",
	)
	logmodule.InjectLogger(L, logger)

	L.PreloadModule("log", logmodule.Loader)

	if err := L.DoString(`
		local log = require("log")
		-- These logs automatically include app, version, and request_id
		log.info("Processing user request")
		log.info("Request completed successfully")

		-- Logger object also has injected fields
		local logger = log.logger()
		logger:info("Using logger object with injected fields")
	`); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
