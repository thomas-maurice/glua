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

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// now: returns the current Unix timestamp.
func now() int64 {
	return time.Now().Unix()
}

// parse: parses a time string with the given Go layout; raises on failure.
func parse(timeStr, layout string) (int64, error) {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %w", err)
	}
	return t.Unix(), nil
}

// parseRFC3339: parses an RFC3339 time string; raises on failure.
func parseRFC3339(timeStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse RFC3339 time: %w", err)
	}
	return t.Unix(), nil
}

// format: formats a Unix timestamp using the given Go layout.
func format(timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).UTC().Format(layout)
}

// add: adds seconds to a Unix timestamp.
func add(timestamp, seconds int64) int64 {
	return time.Unix(timestamp, 0).Add(time.Duration(seconds) * time.Second).Unix()
}

// diff: returns the difference in seconds (t1 - t2).
func diff(t1, t2 int64) int64 {
	return t1 - t2
}

// sleep: pauses execution for the given number of seconds (fractional OK).
func sleep(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}

// toOsdate: converts a Unix timestamp to a Lua os.date-compatible table.
func toOsdate(L *lua.LState, timestamp int64) *lua.LTable {
	t := time.Unix(timestamp, 0).UTC()
	tbl := L.NewTable()
	tbl.RawSetString("year", lua.LNumber(t.Year()))
	tbl.RawSetString("month", lua.LNumber(t.Month()))
	tbl.RawSetString("day", lua.LNumber(t.Day()))
	tbl.RawSetString("hour", lua.LNumber(t.Hour()))
	tbl.RawSetString("min", lua.LNumber(t.Minute()))
	tbl.RawSetString("sec", lua.LNumber(t.Second()))
	tbl.RawSetString("wday", lua.LNumber(t.Weekday()+1)) // Lua uses 1=Sunday
	tbl.RawSetString("yday", lua.LNumber(t.YearDay()))
	tbl.RawSetString("isdst", lua.LBool(false))
	return tbl
}

// fromOsdate: converts an os.date-compatible table to a Unix timestamp.
// Required fields: year, month, day. Optional: hour, min, sec (default 0).
// Raises on missing or wrong-typed required fields.
func fromOsdate(L *lua.LState, tbl *lua.LTable) (int64, error) {
	year, err := requireNumberField(tbl, "year")
	if err != nil {
		return 0, err
	}
	month, err := requireNumberField(tbl, "month")
	if err != nil {
		return 0, err
	}
	day, err := requireNumberField(tbl, "day")
	if err != nil {
		return 0, err
	}
	hour := optionalNumberField(tbl, "hour", 0)
	min := optionalNumberField(tbl, "min", 0)
	sec := optionalNumberField(tbl, "sec", 0)

	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC)
	return t.Unix(), nil
}

// requireNumberField: extracts a numeric field from a Lua table; errors if
// missing or not a number.
func requireNumberField(tbl *lua.LTable, name string) (int, error) {
	v := tbl.RawGetString(name)
	n, ok := v.(lua.LNumber)
	if !ok {
		return 0, fmt.Errorf("field %q is required and must be a number, got %s", name, v.Type())
	}
	return int(n), nil
}

// optionalNumberField: extracts an optional numeric field, returning def if
// absent or wrong-typed.
func optionalNumberField(tbl *lua.LTable, name string, def int) int {
	v := tbl.RawGetString(name)
	if n, ok := v.(lua.LNumber); ok {
		return int(n)
	}
	return def
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("time", "time and date utilities")
	m.Fn("now", now, "returns the current Unix timestamp")
	m.Fn("parse", parse, "parses a time string with a Go layout, raises on error",
		luareg.Args("timestr", "layout"))
	m.Fn("parse_rfc3339", parseRFC3339, "parses an RFC3339 time string, raises on error",
		luareg.Args("timestr"))
	m.Fn("format", format, "formats a Unix timestamp with a Go layout",
		luareg.Args("timestamp", "layout"))
	m.Fn("add", add, "adds seconds to a Unix timestamp",
		luareg.Args("timestamp", "seconds"))
	m.Fn("diff", diff, "returns the difference in seconds between two timestamps (t1 - t2)",
		luareg.Args("t1", "t2"))
	m.Fn("sleep", sleep, "pauses execution for the given number of seconds",
		luareg.Args("seconds"))
	m.Fn("to_osdate", toOsdate, "converts a Unix timestamp to an os.date-compatible table",
		luareg.Args("timestamp"))
	m.Fn("from_osdate", fromOsdate, "converts an os.date-compatible table to a Unix timestamp, raises on invalid input",
		luareg.Args("date_table"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("time", time.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
