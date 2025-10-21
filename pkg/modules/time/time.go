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

package time

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the time module for Lua.
// This function should be registered with L.PreloadModule("time", time.Loader)
//
// @luamodule time
//
// Example usage in Lua:
//
//	local time = require("time")
//	local now = time.now()
//	local formatted = time.format(now, "2006-01-02 15:04:05")
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"now":           now,
	"parse":         parse,
	"format":        format,
	"add":           add,
	"diff":          diff,
	"sleep":         sleep,
	"to_osdate":     toOsdate,
	"from_osdate":   fromOsdate,
	"parse_rfc3339": parseRFC3339,
}

// now: returns the current Unix timestamp.
//
// @luafunc now
// @luareturn number timestamp Current Unix timestamp (seconds since epoch)
//
// Example:
//
//	local now = time.now()
//	print(now)  -- prints current timestamp
func now(L *lua.LState) int {
	L.Push(lua.LNumber(time.Now().Unix()))
	return 1
}

// parse: parses a time string using Go time format layout.
//
// @luafunc parse
// @luaparam timestr string The time string to parse
// @luaparam layout string The Go time layout format (e.g., "2006-01-02 15:04:05")
// @luareturn number timestamp Unix timestamp, or nil on error
// @luareturn string|nil err Error message if parsing failed
//
// Example:
//
//	local ts, err = time.parse("2024-03-15 14:30:00", "2006-01-02 15:04:05")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print("Timestamp: " .. ts)
//	end
func parse(L *lua.LState) int {
	timeStr := L.CheckString(1)
	layout := L.CheckString(2)

	t, err := time.Parse(layout, timeStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse time: %v", err)))
		return 2
	}

	L.Push(lua.LNumber(t.Unix()))
	L.Push(lua.LNil)
	return 2
}

// parseRFC3339: parses an RFC3339 time string (common in Kubernetes).
//
// @luafunc parse_rfc3339
// @luaparam timestr string The RFC3339 time string (e.g., "2024-03-15T14:30:00Z")
// @luareturn number timestamp Unix timestamp, or nil on error
// @luareturn string|nil err Error message if parsing failed
//
// Example:
//
//	local ts, err = time.parse_rfc3339("2024-03-15T14:30:00Z")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print("Timestamp: " .. ts)
//	end
func parseRFC3339(L *lua.LState) int {
	timeStr := L.CheckString(1)

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse RFC3339 time: %v", err)))
		return 2
	}

	L.Push(lua.LNumber(t.Unix()))
	L.Push(lua.LNil)
	return 2
}

// format: formats a Unix timestamp using Go time format layout.
//
// @luafunc format
// @luaparam timestamp number Unix timestamp
// @luaparam layout string The Go time layout format (e.g., "2006-01-02 15:04:05")
// @luareturn string formatted Formatted time string
//
// Example:
//
//	local formatted = time.format(1710512400, "2006-01-02 15:04:05")
//	print(formatted)  -- prints "2024-03-15 14:30:00"
func format(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	layout := L.CheckString(2)

	t := time.Unix(int64(timestamp), 0).UTC()
	formatted := t.Format(layout)

	L.Push(lua.LString(formatted))
	return 1
}

// add: adds seconds to a Unix timestamp.
//
// @luafunc add
// @luaparam timestamp number Unix timestamp
// @luaparam seconds number Number of seconds to add (can be negative)
// @luareturn number new_timestamp New Unix timestamp
//
// Example:
//
//	local tomorrow = time.add(time.now(), 86400)  -- add 24 hours
//	local yesterday = time.add(time.now(), -86400)  -- subtract 24 hours
func add(L *lua.LState) int {
	timestamp := L.CheckNumber(1)
	seconds := L.CheckNumber(2)

	t := time.Unix(int64(timestamp), 0)
	newTime := t.Add(time.Duration(seconds) * time.Second)

	L.Push(lua.LNumber(newTime.Unix()))
	return 1
}

// diff: calculates the difference between two Unix timestamps.
//
// @luafunc diff
// @luaparam time1 number First Unix timestamp
// @luaparam time2 number Second Unix timestamp
// @luareturn number seconds Difference in seconds (time1 - time2)
//
// Example:
//
//	local age = time.diff(time.now(), pod_creation_time)
//	print("Pod is " .. age .. " seconds old")
func diff(L *lua.LState) int {
	time1 := L.CheckNumber(1)
	time2 := L.CheckNumber(2)

	diff := int64(time1) - int64(time2)

	L.Push(lua.LNumber(diff))
	return 1
}

// sleep: pauses execution for the specified number of seconds.
//
// @luafunc sleep
// @luaparam seconds number Number of seconds to sleep
//
// Example:
//
//	print("Starting...")
//	time.sleep(2)
//	print("2 seconds later")
func sleep(L *lua.LState) int {
	seconds := L.CheckNumber(1)
	time.Sleep(time.Duration(float64(seconds) * float64(time.Second)))
	return 0
}

// toOsdate: converts a Unix timestamp to Lua os.date compatible table.
//
// @luafunc to_osdate
// @luaparam timestamp number Unix timestamp
// @luareturn table date_table Table with year, month, day, hour, min, sec, wday, yday, isdst
//
// Example:
//
//	local dt = time.to_osdate(time.now())
//	print(dt.year .. "-" .. dt.month .. "-" .. dt.day)
func toOsdate(L *lua.LState) int {
	timestamp := L.CheckNumber(1)

	t := time.Unix(int64(timestamp), 0).UTC()

	tbl := L.NewTable()
	tbl.RawSetString("year", lua.LNumber(t.Year()))
	tbl.RawSetString("month", lua.LNumber(t.Month()))
	tbl.RawSetString("day", lua.LNumber(t.Day()))
	tbl.RawSetString("hour", lua.LNumber(t.Hour()))
	tbl.RawSetString("min", lua.LNumber(t.Minute()))
	tbl.RawSetString("sec", lua.LNumber(t.Second()))
	tbl.RawSetString("wday", lua.LNumber(t.Weekday()+1)) // Lua uses 1=Sunday
	tbl.RawSetString("yday", lua.LNumber(t.YearDay()))
	tbl.RawSetString("isdst", lua.LBool(false)) // Go doesn't track DST in time.Time

	L.Push(tbl)
	return 1
}

// fromOsdate: converts a Lua os.date compatible table to Unix timestamp.
//
// @luafunc from_osdate
// @luaparam date_table table Table with year, month, day, hour, min, sec (other fields optional)
// @luareturn number Unix timestamp
//
// Example:
//
//	local ts = time.from_osdate({year=2024, month=3, day=15, hour=14, min=30, sec=0})
//	print(ts)
func fromOsdate(L *lua.LState) int {
	tbl := L.CheckTable(1)

	year := int(tbl.RawGetString("year").(lua.LNumber))
	month := time.Month(tbl.RawGetString("month").(lua.LNumber))
	day := int(tbl.RawGetString("day").(lua.LNumber))
	hour := int(tbl.RawGetString("hour").(lua.LNumber))
	min := int(tbl.RawGetString("min").(lua.LNumber))
	sec := int(tbl.RawGetString("sec").(lua.LNumber))

	t := time.Date(year, month, day, hour, min, sec, 0, time.UTC)

	L.Push(lua.LNumber(t.Unix()))
	return 1
}
