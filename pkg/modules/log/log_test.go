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

package log

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"
	lua "github.com/yuin/gopher-lua"
)

// TestModuleInfo: tests that module-level info logging works
func TestModuleInfo(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		log.info("Test message", "key", "value", "number", 42)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected log output to contain 'Test message', got: %s", output)
	}
	if !strings.Contains(output, "key") {
		t.Errorf("Expected log output to contain 'key', got: %s", output)
	}
}

// TestLoggerObject: tests using the logger object directly
func TestLoggerObject(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		local logger = log.logger()
		logger:info("Logger object test", "key", "value")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Logger object test") {
		t.Errorf("Expected log output to contain 'Logger object test', got: %s", output)
	}
}

// TestLoggerWith: tests that logger:with() returns a new logger with fields
func TestLoggerWith(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		local logger = log.logger()
		local contextLogger = logger:with("request_id", "abc123", "user", "john")
		contextLogger:info("First message")
		contextLogger:info("Second message")
		logger:info("Third message without context")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 log lines, got %d: %s", len(lines), output)
	}

	// Check first two lines contain context fields
	for i := 0; i < 2; i++ {
		line := lines[i]
		if !strings.Contains(line, "request_id") {
			t.Errorf("Expected line %d to contain 'request_id', got: %s", i+1, line)
		}
		if !strings.Contains(line, "abc123") {
			t.Errorf("Expected line %d to contain 'abc123', got: %s", i+1, line)
		}
	}

	// Third line should not have context
	if strings.Contains(lines[2], "request_id") {
		t.Errorf("Expected line 3 to NOT contain 'request_id', got: %s", lines[2])
	}
}

// TestLoggerWithChaining: tests chaining logger:with() calls
func TestLoggerWithChaining(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		local logger = log.logger()
		local reqLogger = logger:with("request_id", "req-123")
		local userLogger = reqLogger:with("user_id", 456)

		reqLogger:info("Request started")
		userLogger:info("User action")
		logger:info("System message")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 log lines, got %d", len(lines))
	}

	// First line: request_id only
	if !strings.Contains(lines[0], "request_id") {
		t.Errorf("Expected line 1 to contain 'request_id'")
	}
	if strings.Contains(lines[0], "user_id") {
		t.Errorf("Expected line 1 to NOT contain 'user_id'")
	}

	// Second line: both request_id and user_id
	if !strings.Contains(lines[1], "request_id") {
		t.Errorf("Expected line 2 to contain 'request_id'")
	}
	if !strings.Contains(lines[1], "user_id") {
		t.Errorf("Expected line 2 to contain 'user_id'")
	}

	// Third line: no extra fields
	if strings.Contains(lines[2], "request_id") {
		t.Errorf("Expected line 3 to NOT contain 'request_id'")
	}
}

// TestDefaultLogger: tests that default logger is used when none is injected
func TestDefaultLogger(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		log.info("Using default logger")
	`

	// Should not error even without injecting a logger
	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}
}

// TestAllLogLevels: tests all log levels work correctly
func TestAllLogLevels(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
		Level:           log.DebugLevel,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		log.debug("Debug message")
		log.info("Info message")
		log.warn("Warn message")
		log.error("Error message")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Debug message") {
		t.Errorf("Expected output to contain 'Debug message'")
	}
	if !strings.Contains(output, "Info message") {
		t.Errorf("Expected output to contain 'Info message'")
	}
	if !strings.Contains(output, "Warn message") {
		t.Errorf("Expected output to contain 'Warn message'")
	}
	if !strings.Contains(output, "Error message") {
		t.Errorf("Expected output to contain 'Error message'")
	}
}

// TestLoggerObjectMethods: tests all logger object methods
func TestLoggerObjectMethods(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
		Level:           log.DebugLevel,
	})
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		local logger = log.logger()
		logger:debug("Debug via object")
		logger:info("Info via object")
		logger:warn("Warn via object")
		logger:error("Error via object")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Debug via object") {
		t.Errorf("Expected output to contain 'Debug via object'")
	}
	if !strings.Contains(output, "Info via object") {
		t.Errorf("Expected output to contain 'Info via object'")
	}
}

// TestInjectLoggerWithFields: tests injecting a logger with pre-set fields from Go
func TestInjectLoggerWithFields(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var buf bytes.Buffer
	logger := log.NewWithOptions(&buf, log.Options{
		TimeFormat:      time.RFC3339,
		ReportTimestamp: false,
		ReportCaller:    false,
		Formatter:       log.JSONFormatter,
	})
	logger = logger.With("app", "test-app", "version", "1.0.0")
	InjectLogger(L, logger)

	L.PreloadModule("log", Loader)

	script := `
		local log = require("log")
		log.info("Service started")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute Lua script: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "app") || !strings.Contains(output, "test-app") {
		t.Errorf("Expected output to contain pre-set fields")
	}
}
