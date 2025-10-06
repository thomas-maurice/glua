package kubernetes

import (
	"fmt"
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
			script := `
				local k8s = require("kubernetes")
				local result, err = k8s.parse_memory("` + tt.input + `")
				return result, err
			`

			if err := L.DoString(script); err != nil {
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
			script := `
				local k8s = require("kubernetes")
				local result, err = k8s.parse_cpu("` + tt.input + `")
				return result, err
			`

			if err := L.DoString(script); err != nil {
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
			script := `
				local k8s = require("kubernetes")
				local result, err = k8s.parse_time("` + tt.input + `")
				return result, err
			`

			if err := L.DoString(script); err != nil {
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
			script := fmt.Sprintf(`
				local k8s = require("kubernetes")
				local result, err = k8s.format_time(%d)
				return result, err
			`, tt.timestamp)

			if err := L.DoString(script); err != nil {
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
			script := fmt.Sprintf(`
				local k8s = require("kubernetes")

				-- Format timestamp to string
				local timestr, err1 = k8s.format_time(%d)
				assert(err1 == nil, "format_time failed: " .. tostring(err1))

				-- Parse string back to timestamp
				local timestamp, err2 = k8s.parse_time(timestr)
				assert(err2 == nil, "parse_time failed: " .. tostring(err2))

				-- Should match original
				assert(timestamp == %d, string.format("Round-trip failed: %%d != %d", timestamp))

				return true
			`, tt.timestamp, tt.timestamp, tt.timestamp)

			if err := L.DoString(script); err != nil {
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

func TestModuleIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	script := `
		local k8s = require("kubernetes")

		-- Test parse_memory
		local mem_bytes = k8s.parse_memory("1Gi")
		assert(mem_bytes == 1073741824, "Memory parsing failed")

		-- Test parse_cpu
		local cpu_millis = k8s.parse_cpu("100m")
		assert(cpu_millis == 100, "CPU parsing failed")

		-- Test parse_time
		local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")
		assert(timestamp > 0, "Time parsing failed")

		-- Test format_time
		local timestr = k8s.format_time(1759509540)
		assert(timestr == "2025-10-03T16:39:00Z", "Time formatting failed")

		-- Test round-trip
		local formatted = k8s.format_time(timestamp)
		local parsed = k8s.parse_time(formatted)
		assert(parsed == timestamp, "Round-trip failed")

		return true
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Errorf("Expected true, got %v", result)
	}
}
