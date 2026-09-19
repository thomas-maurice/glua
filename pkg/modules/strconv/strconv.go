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

// Package strconv binds Go's strconv package for Lua scripts: numeric
// parsing/formatting and Go-syntax string quoting.
//
// Lua 5.1 already has tonumber(s) and tonumber(s, base), so parse_int and
// parse_float are admittedly thin wrappers — the value this module adds over
// them is narrow but real:
//
//   - tonumber returns nil with no explanation on failure. parse_int and
//     parse_float raise with Go's own message (e.g. `strconv.ParseInt:
//     parsing "12a": invalid syntax`), matching this repo's fail-loud
//     convention (a trailing error return aborts the script).
//   - tonumber silently returns a float for an integer literal that doesn't
//     fit. parse_int and format_int instead raise when the value's magnitude
//     exceeds 2^53, the largest integer a float64 (what every Lua number is)
//     can represent exactly. Silently handing back a rounded value would be
//     worse than raising.
//   - format_int(n, base) has no Lua 5.1 equivalent at all — string.format
//     has no base-N integer conversion.
//   - format_float(f, fmt, prec) with prec == -1 gives shortest
//     round-trip formatting, which string.format cannot do.
//   - quote/unquote (Go-syntax string literals with escapes) have no
//     equivalent in Lua 5.1.
//   - parse_bool accepting Go's full spelling set (1/t/T/TRUE/true/True/
//     0/f/F/FALSE/false/False) is genuinely missing from Lua.
//
// There is no itoa: it is exactly format_int(n, 10), and this module picks
// one name for one job rather than shipping two.
package strconv

import (
	"fmt"
	"math"
	"strconv"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// maxSafeInt: the largest magnitude (2^53) an integer can have and still
// round-trip through a float64 (what every Lua number is) exactly. parse_int
// and format_int raise rather than silently handing back — or formatting — a
// value that Lua could not have held in the first place.
const maxSafeInt = 1 << 53

// parseInt: parses s as an integer in the given base. base == 0 infers the
// base from s's prefix (0x/0o/0b/leading 0), exactly like Go's
// strconv.ParseInt — including that digit-separator underscores are only
// permitted when base == 0. Raises (via the returned error) on a syntax
// error, a value that overflows int64, or a value whose magnitude exceeds
// 2^53 and therefore cannot be represented exactly as a Lua number.
func parseInt(s string, base int) (float64, error) {
	v, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, err
	}
	if v > maxSafeInt || v < -maxSafeInt {
		return 0, fmt.Errorf("strconv.parse_int: value %d exceeds the exact integer range of a float64 (±2^53)", v)
	}
	return float64(v), nil
}

// atoi: base-10 shorthand for parseInt(s, 10).
func atoi(s string) (float64, error) {
	return parseInt(s, 10)
}

// parseFloat: parses s as a float64, accepting the same syntax as Go's
// strconv.ParseFloat — including "Inf", "+Inf", "-Inf" and "NaN" spellings.
// Raises on a syntax error.
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// parseBool: parses s as a boolean. Accepts exactly the set Go's
// strconv.ParseBool accepts: "1", "t", "T", "TRUE", "true", "True", "0",
// "f", "F", "FALSE", "false", "False". Raises on anything else.
func parseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// formatInt: formats n as a string in the given base (2-36, lowercase
// digits). Raises if n is not an exact integer, if |n| exceeds 2^53 (the
// value could not have arrived from a genuine Lua integer), or if base is
// outside [2, 36] — Go's own strconv.FormatInt panics on an invalid base
// rather than erroring, so that check is done here first to keep this a
// raise, not a crash.
func formatInt(n float64, base int) (string, error) {
	if base < 2 || base > 36 {
		return "", fmt.Errorf("strconv.format_int: base must be in [2, 36], got %d", base)
	}
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return "", fmt.Errorf("strconv.format_int: n must be a finite integer, got %v", n)
	}
	if n != math.Trunc(n) {
		return "", fmt.Errorf("strconv.format_int: n must be an integer, got %v", n)
	}
	if n > maxSafeInt || n < -maxSafeInt {
		return "", fmt.Errorf("strconv.format_int: n %v exceeds the exact integer range of a float64 (±2^53)", n)
	}
	return strconv.FormatInt(int64(n), base), nil
}

// formatFloat: formats f according to format (one of "b" "e" "E" "f" "g" "G"
// "x" "X") and prec. prec == -1 uses the shortest representation that
// round-trips exactly through parseFloat. Raises on an unrecognised format
// or prec < -1 — Go's own strconv.FormatFloat neither panics nor errors on
// an unknown format byte, it just silently emits a bogus "%!x(...)"-style
// string, so this is validated here explicitly rather than trusted.
func formatFloat(f float64, format string, prec int) (string, error) {
	if len(format) != 1 {
		return "", fmt.Errorf("strconv.format_float: fmt must be a single character, one of b e E f g G x X, got %q", format)
	}
	switch format[0] {
	case 'b', 'e', 'E', 'f', 'g', 'G', 'x', 'X':
	default:
		return "", fmt.Errorf("strconv.format_float: unknown fmt %q, must be one of b e E f g G x X", format)
	}
	if prec < -1 {
		return "", fmt.Errorf("strconv.format_float: prec must be >= -1, got %d", prec)
	}
	return strconv.FormatFloat(f, format[0], prec, 64), nil
}

// quote: returns a Go-syntax double-quoted string literal for s, escaping
// every non-printable and non-ASCII-printable byte. Never raises.
func quote(s string) string {
	return strconv.Quote(s)
}

// unquote: the inverse of quote. Accepts a Go-syntax double-quoted literal
// ("..."), a raw string literal (`...`), or a single-character literal
// ('c'). Raises on anything that is not a valid Go string/rune literal.
func unquote(s string) (string, error) {
	return strconv.Unquote(s)
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("strconv", "numeric parsing/formatting and Go-syntax string quoting with useful errors")

	m.Fn("parse_int", parseInt, "parses a string as an integer, raising on syntax/range errors instead of returning nil",
		luareg.Args("s", "base"),
		luareg.ArgDoc("s", "the string to parse"),
		luareg.ArgDoc("base", "numeric base 2-36, or 0 to infer from s's prefix (0x/0o/0b/leading 0); digit-separator underscores are only accepted when base is 0"),
		luareg.ReturnDoc(0, "n", "the parsed integer; raises if it cannot fit exactly in a float64 (magnitude > 2^53)"))
	m.Fn("atoi", atoi, "parses a string as a base-10 integer; shorthand for parse_int(s, 10)",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to parse"),
		luareg.ReturnDoc(0, "n", "the parsed integer; raises if it cannot fit exactly in a float64 (magnitude > 2^53)"))
	m.Fn("parse_float", parseFloat, "parses a string as a floating-point number, raising on syntax errors instead of returning nil",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to parse; accepts decimal, exponent, hex-float, and Inf/NaN spellings"),
		luareg.ReturnDoc(0, "f", "the parsed float"))
	m.Fn("parse_bool", parseBool, "parses a string as a boolean, accepting Go's spelling set",
		luareg.Args("s"),
		luareg.ArgDoc("s", "one of 1 t T TRUE true True 0 f F FALSE false False"),
		luareg.ReturnDoc(0, "b", "the parsed boolean"))
	m.Fn("format_int", formatInt, "formats an integer in the given base; has no Lua 5.1 equivalent",
		luareg.Args("n", "base"),
		luareg.ArgDoc("n", "the value to format; must be an exact integer with magnitude <= 2^53"),
		luareg.ArgDoc("base", "numeric base, 2-36; digits above 9 are lowercase letters"),
		luareg.ReturnDoc(0, "s", "n rendered in the given base"))
	m.Fn("format_float", formatFloat, "formats a float with an explicit format and precision; has no Lua 5.1 equivalent",
		luareg.Args("f", "fmt", "prec"),
		luareg.ArgDoc("f", "the value to format"),
		luareg.ArgDoc("fmt", "one of b e E f g G x X, matching Go's strconv.FormatFloat verbs"),
		luareg.ArgDoc("prec", "digits of precision; -1 selects the shortest representation that round-trips exactly through parse_float"),
		luareg.ReturnDoc(0, "s", "f rendered per fmt and prec"))
	m.Fn("quote", quote, "returns a Go-syntax double-quoted string literal for s, with non-printables escaped",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to quote"),
		luareg.ReturnDoc(0, "quoted", "a double-quoted Go string literal representing s"))
	m.Fn("unquote", unquote, "parses a Go-syntax string, raw string, or rune literal back into its raw value",
		luareg.Args("s"),
		luareg.ArgDoc("s", `a "..." double-quoted, `+"`...`"+" raw, or 'c' rune literal"),
		luareg.ReturnDoc(0, "s", "the decoded raw string"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("strconv", strconv.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
