// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, L.DoString(code))
}

func TestParse(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.parse("2024-03-15 14:30:00", "2006-01-02 15:04:05")
		assert(type(ts) == "number", "Expected number")
		assert(ts == 1710513000, "Expected specific timestamp, got " .. ts)
	`
	require.NoError(t, L.DoString(code))
}

func TestParseRFC3339(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.parse_rfc3339("2024-03-15T14:30:00Z")
		assert(type(ts) == "number", "Expected number")
		assert(ts == 1710513000, "Expected specific timestamp, got " .. ts)
	`
	require.NoError(t, L.DoString(code))
}

// TestParseInvalid: invalid time string raises a Lua error.
func TestParseInvalid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ok, err = pcall(time.parse, "invalid", "2006-01-02")
		assert(not ok, "Expected error for invalid time string")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
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
	require.NoError(t, L.DoString(code))
}

func TestAdd(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local base = 1710513000
		local future = time.add(base, 3600)
		local past = time.add(base, -3600)
		assert(future == base + 3600, "Expected timestamp + 1 hour")
		assert(past == base - 3600, "Expected timestamp - 1 hour")
	`
	require.NoError(t, L.DoString(code))
}

func TestDiff(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local diff = time.diff(1710513000, 1710509400)
		assert(diff == 3600, "Expected 3600 seconds difference")
	`
	require.NoError(t, L.DoString(code))
}

func TestSleep(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	start := time.Now()
	code := `
		local time = require("time")
		time.sleep(0.1)
	`
	require.NoError(t, L.DoString(code))
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
		local dt = time.to_osdate(1710513000)
		assert(dt.year == 2024, "Expected year 2024")
		assert(dt.month == 3, "Expected month 3")
		assert(dt.day == 15, "Expected day 15")
		assert(dt.hour == 14, "Expected hour 14")
		assert(dt.min == 30, "Expected minute 30")
		assert(dt.sec == 0, "Expected second 0")
	`
	require.NoError(t, L.DoString(code))
}

func TestFromOsdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.from_osdate({year=2024, month=3, day=15, hour=14, min=30, sec=0})
		assert(ts == 1710513000, "Expected timestamp 1710513000, got " .. ts)
	`
	require.NoError(t, L.DoString(code))
}

// TestFromOsdateDefaults: hour/min/sec are optional and default to 0.
func TestFromOsdateDefaults(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.from_osdate({year=2024, month=3, day=15})
		assert(ts == 1710460800, "Expected midnight UTC timestamp, got " .. tostring(ts))
	`
	require.NoError(t, L.DoString(code))
}

// TestFromOsdateMissingField: a missing required field raises a Lua error.
func TestFromOsdateMissingField(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ok, err = pcall(time.from_osdate, {year=2024, month=3})
		assert(not ok, "Expected error for missing field")
		assert(type(err) == "string", "Error should be a string")
		assert(string.find(err, "day") ~= nil, "Error should mention missing field, got: " .. err)
	`
	require.NoError(t, L.DoString(code))
}

// TestFromOsdateWrongType: wrong-typed required field raises a Lua error.
func TestFromOsdateWrongType(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ok, err = pcall(time.from_osdate, {year="twenty-twenty-four", month=3, day=15})
		assert(not ok, "Expected error for wrong-typed field")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}

func TestRoundTripOsdate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local original = 1710513000
		local dt = time.to_osdate(original)
		local ts = time.from_osdate(dt)
		assert(ts == original, "Round-trip should preserve timestamp")
	`
	require.NoError(t, L.DoString(code))
}

func TestRoundTripParseFormat(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local original = "2024-03-15 14:30:00"
		local layout = "2006-01-02 15:04:05"
		local ts = time.parse(original, layout)
		local formatted = time.format(ts, layout)
		assert(formatted == original, "Round-trip should preserve format")
	`
	require.NoError(t, L.DoString(code))
}

func TestKubernetesTimestamp(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("time", Loader)

	code := `
		local time = require("time")
		local ts = time.parse_rfc3339("2024-10-03T16:39:00Z")
		local now = time.now()
		local age = time.diff(now, ts)
		assert(age > 0, "Pod should have positive age")
		local dt = time.to_osdate(ts)
		assert(dt.year == 2024, "Correct year")
		assert(dt.month == 10, "Correct month")
	`
	require.NoError(t, L.DoString(code))
}
