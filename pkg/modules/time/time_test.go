// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package time

import (
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestNow(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local now = time.now()
		assert(type(now) == "number", "Expected number")
		assert(now > 0, "Timestamp should be positive")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestParse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts, err = time.parse("2024-03-15 14:30:00", "2006-01-02 15:04:05")
		assert(err == nil, "Expected no error")
		assert(type(ts) == "number", "Expected number")
		assert(ts == 1710513000, "Expected specific timestamp")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestParseRFC3339(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts, err = time.parse_rfc3339("2024-03-15T14:30:00Z")
		assert(err == nil, "Expected no error")
		assert(type(ts) == "number", "Expected number")
		assert(ts == 1710513000, "Expected specific timestamp")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestParseInvalid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts, err = time.parse("invalid", "2006-01-02")
		assert(ts == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestFormat(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local formatted = time.format(1710513000, "2006-01-02 15:04:05")
		assert(formatted == "2024-03-15 14:30:00", "Expected correct format")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestAdd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local base = 1710513000
		local future = time.add(base, 3600)  -- add 1 hour
		local past = time.add(base, -3600)  -- subtract 1 hour

		assert(future == base + 3600, "Expected timestamp + 1 hour")
		assert(past == base - 3600, "Expected timestamp - 1 hour")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDiff(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local t1 = 1710513000
		local t2 = 1710509400
		local diff = time.diff(t1, t2)

		assert(diff == 3600, "Expected 3600 seconds difference")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestSleep(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	start := time.Now()

	code := `
		local time = require("time")
		time.sleep(0.1)  -- sleep 100ms
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}

	elapsed := time.Since(start)
	if elapsed < 100*time.Millisecond {
		t.Errorf("Sleep did not wait long enough: %v", elapsed)
	}
}

func TestToOsdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local dt = time.to_osdate(1710513000)  -- 2024-03-15 14:30:00 UTC

		assert(dt.year == 2024, "Expected year 2024")
		assert(dt.month == 3, "Expected month 3")
		assert(dt.day == 15, "Expected day 15")
		assert(dt.hour == 14, "Expected hour 14")
		assert(dt.min == 30, "Expected minute 30")
		assert(dt.sec == 0, "Expected second 0")
		assert(type(dt.wday) == "number", "wday should be number")
		assert(type(dt.yday) == "number", "yday should be number")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestFromOsdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.from_osdate({
			year = 2024,
			month = 3,
			day = 15,
			hour = 14,
			min = 30,
			sec = 0
		})

		assert(ts == 1710513000, "Expected timestamp 1710513000")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTripOsdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local original = 1710513000

		-- Convert to osdate and back
		local dt = time.to_osdate(original)
		local ts = time.from_osdate(dt)

		assert(ts == original, "Round-trip should preserve timestamp")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTripParseFomat(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local original = "2024-03-15 14:30:00"
		local layout = "2006-01-02 15:04:05"

		-- Parse and format back
		local ts, err = time.parse(original, layout)
		assert(err == nil, "Parse should succeed")

		local formatted = time.format(ts, layout)
		assert(formatted == original, "Round-trip should preserve format")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestKubernetesTimestamp(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")

		-- Parse Kubernetes-style timestamp
		local k8s_time = "2024-10-03T16:39:00Z"
		local ts, err = time.parse_rfc3339(k8s_time)
		assert(err == nil, "Should parse K8s timestamp")

		-- Calculate age
		local now = time.now()
		local age = time.diff(now, ts)
		assert(age > 0, "Pod should have positive age")

		-- Format for display
		local dt = time.to_osdate(ts)
		assert(dt.year == 2024, "Correct year")
		assert(dt.month == 10, "Correct month")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
