package kubernetes

import (
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

func TestInitDefaultsFullWorkflow(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	if err := L.DoFile("testdata/init_defaults_full_workflow.lua"); err != nil {
		t.Fatalf("Test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestModuleIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	if err := L.DoFile("testdata/module_integration.lua"); err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Errorf("Expected true, got %v", result)
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

func TestParseIntOrString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name      string
		input     lua.LValue
		expectVal interface{}
		isString  bool
	}{
		{"number", lua.LNumber(8080), 8080.0, false},
		{"string", lua.LString("http"), "http", true},
		{"string number", lua.LString("8080"), "8080", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("test_input", tt.input)

			if err := L.DoFile("testdata/test_parse_int_or_string.lua"); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}

			result := L.Get(-2)
			isStr := L.Get(-1)
			L.Pop(2)

			if tt.isString {
				if isStr != lua.LTrue {
					t.Errorf("Expected is_string to be true")
				}
				if str, ok := result.(lua.LString); ok {
					if string(str) != tt.expectVal.(string) {
						t.Errorf("Expected %s, got %s", tt.expectVal.(string), string(str))
					}
				} else {
					t.Errorf("Expected LString, got %T", result)
				}
			} else {
				if isStr != lua.LFalse {
					t.Errorf("Expected is_string to be false")
				}
				if num, ok := result.(lua.LNumber); ok {
					if float64(num) != tt.expectVal.(float64) {
						t.Errorf("Expected %f, got %f", tt.expectVal.(float64), float64(num))
					}
				} else {
					t.Errorf("Expected LNumber, got %T", result)
				}
			}
		})
	}
}

func TestMatchesSelector(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		filename string
	}{
		{"selector matches", "testdata/test_matches_selector_true.lua"},
		{"selector does not match", "testdata/test_matches_selector_false.lua"},
		{"missing label", "testdata/test_matches_selector_missing.lua"},
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

func TestTolerationMatches(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		filename string
	}{
		{"equal operator match", "testdata/test_toleration_matches_equal.lua"},
		{"exists operator match", "testdata/test_toleration_matches_exists.lua"},
		{"no match", "testdata/test_toleration_no_match.lua"},
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

func TestNewFunctionsIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	if err := L.DoFile("testdata/integration_new_functions.lua"); err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestMatchGVK(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		filename string
	}{
		{"pod matches v1/Pod", "testdata/test_match_gvk_pod.lua"},
		{"deployment matches apps/v1/Deployment", "testdata/test_match_gvk_deployment.lua"},
		{"wrong kind does not match", "testdata/test_match_gvk_wrong_kind.lua"},
		{"wrong version does not match", "testdata/test_match_gvk_wrong_version.lua"},
		{"configmap matches v1/ConfigMap", "testdata/test_match_gvk_configmap.lua"},
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

func TestMatchGVKIntegration(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	if err := L.DoFile("testdata/integration_match_gvk.lua"); err != nil {
		t.Fatalf("Integration test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Errorf("Expected true, got %v", result)
	}
}
