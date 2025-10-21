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

package kubernetes

import (
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestParseMemory(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"1024Mi", "1024Mi", 1073741824, false},
		{"1Gi", "1Gi", 1073741824, false},
		{"512M", "512M", 512000000, false},
		{"1Ki", "1Ki", 1024, false},
		{"invalid", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_input", lua.LString(tt.input))

			if err := L.DoFile("testdata/test_parse_memory.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if tt.wantErr {
				if errVal == lua.LNil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
			} else {
				if errVal != lua.LNil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, errVal)
				}

				if num, ok := result.(lua.LNumber); ok {
					if int64(num) != tt.expected {
						t.Errorf("Expected %d, got %d", tt.expected, int64(num))
					}
				} else {
					t.Errorf("Expected LNumber, got %T", result)
				}
			}
		})
	}
}

func TestParseCPU(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"100m", "100m", 100, false},
		{"1", "1", 1000, false},
		{"2000m", "2000m", 2000, false},
		{"500m", "500m", 500, false},
		{"invalid", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_input", lua.LString(tt.input))

			if err := L.DoFile("testdata/test_parse_cpu.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if tt.wantErr {
				if errVal == lua.LNil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
			} else {
				if errVal != lua.LNil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, errVal)
				}

				if num, ok := result.(lua.LNumber); ok {
					if int64(num) != tt.expected {
						t.Errorf("Expected %d, got %d", tt.expected, int64(num))
					}
				} else {
					t.Errorf("Expected LNumber, got %T", result)
				}
			}
		})
	}
}

func TestParseTime(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"RFC3339", "2025-10-03T16:39:00Z", 1759509540, false},
		{"RFC3339 with timezone", "2025-10-03T16:39:00+00:00", 1759509540, false},
		{"invalid", "not-a-time", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_input", lua.LString(tt.input))

			if err := L.DoFile("testdata/test_parse_time.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if tt.wantErr {
				if errVal == lua.LNil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
			} else {
				if errVal != lua.LNil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, errVal)
				}

				if num, ok := result.(lua.LNumber); ok {
					if int64(num) != tt.expected {
						t.Errorf("Expected %d, got %d", tt.expected, int64(num))
					}
				} else {
					t.Errorf("Expected LNumber, got %T", result)
				}
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name      string
		timestamp int64
		expected  string
	}{
		{"epoch", 0, "1970-01-01T00:00:00Z"},
		{"specific time", 1759509540, "2025-10-03T16:39:00Z"},
		{"negative timestamp", -1, "1969-12-31T23:59:59Z"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_timestamp", lua.LNumber(tt.timestamp))

			if err := L.DoFile("testdata/test_format_time.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if errVal != lua.LNil {
				t.Errorf("Unexpected error: %v", errVal)
			}

			if str, ok := result.(lua.LString); ok {
				if string(str) != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, string(str))
				}
			} else {
				t.Errorf("Expected LString, got %T", result)
			}
		})
	}
}

func TestFormatParseRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name      string
		timestamp int64
	}{
		{"epoch", 0},
		{"specific time", 1759509540},
		{"recent time", 1696347540}, // 2023-10-03T16:39:00Z
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_timestamp", lua.LNumber(tt.timestamp))

			if err := L.DoFile("testdata/test_roundtrip_time.lua"); err != nil {
				t.Fatalf("Round-trip test failed: %v", err)
			}

			result := L.Get(-1)
			L.Pop(1)

			if result != lua.LTrue {
				t.Errorf("Expected true, got %v", result)
			}
		})
	}
}

func TestInitDefaults(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		filename string
	}{
		{"nil labels and annotations", "testdata/init_defaults_nil.lua"},
		{"existing labels and annotations", "testdata/init_defaults_existing.lua"},
		{"no metadata", "testdata/init_defaults_no_metadata.lua"},
		{"returns same object", "testdata/init_defaults_returns_same.lua"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoFile(tt.filename); err != nil {
				t.Fatalf("Test failed: %v", err)
			}

			result := L.Get(-1)
			L.Pop(1)

			if result != lua.LTrue {
				t.Errorf("Expected true, got %v", result)
			}
		})
	}
}

// TestLuaIntegrationScripts: runs integration test scripts (those not starting with test_)
func TestLuaIntegrationScripts(t *testing.T) {
	patterns := []string{
		"testdata/init_*.lua",
		"testdata/integration_*.lua",
		"testdata/module_*.lua",
	}

	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatalf("Failed to glob %s: %v", pattern, err)
		}
		files = append(files, matches...)
	}

	if len(files) == 0 {
		t.Fatal("No integration test files found in testdata/")
	}

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("kubernetes", Loader)

			if err := L.DoFile(file); err != nil {
				t.Fatalf("Test failed: %v", err)
			}

			result := L.Get(-1)
			L.Pop(1)

			if result != lua.LTrue {
				t.Errorf("Expected true, got %v", result)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{"5 minutes", "5m", 300, false},
		{"1 hour", "1h", 3600, false},
		{"1h30m", "1h30m", 5400, false},
		{"10 seconds", "10s", 10, false},
		{"invalid", "invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_input", lua.LString(tt.input))

			if err := L.DoFile("testdata/test_parse_duration.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if tt.wantErr {
				if errVal == lua.LNil {
					t.Errorf("Expected error for input %s, got nil", tt.input)
				}
			} else {
				if errVal != lua.LNil {
					t.Errorf("Unexpected error for input %s: %v", tt.input, errVal)
				}

				if num, ok := result.(lua.LNumber); ok {
					if float64(num) != tt.expected {
						t.Errorf("Expected %f, got %f", tt.expected, float64(num))
					}
				} else {
					t.Errorf("Expected LNumber, got %T", result)
				}
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{"5 minutes", 300, "5m0s"},
		{"1 hour", 3600, "1h0m0s"},
		{"90 seconds", 90, "1m30s"},
		{"10 seconds", 10, "10s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_seconds", lua.LNumber(tt.seconds))

			if err := L.DoFile("testdata/test_format_duration.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			errVal := L.Get(-1)
			L.Pop(2)

			if errVal != lua.LNil {
				t.Errorf("Unexpected error: %v", errVal)
			}

			if str, ok := result.(lua.LString); ok {
				if string(str) != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, string(str))
				}
			} else {
				t.Errorf("Expected LString, got %T", result)
			}
		})
	}
}

// TestLuaMatchGVKScripts: runs all match_gvk test scripts using file discovery
func TestLuaMatchGVKScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/test_match_gvk_*.lua")
	if err != nil {
		t.Fatalf("Failed to glob testdata: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No match_gvk test files found in testdata/")
	}

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("kubernetes", Loader)

			if err := L.DoFile(file); err != nil {
				t.Fatalf("Test failed: %v", err)
			}

			result := L.Get(-1)
			L.Pop(1)

			if result != lua.LTrue {
				t.Errorf("Expected true, got %v", result)
			}
		})
	}
}
