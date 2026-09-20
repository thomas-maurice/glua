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

// Package bit32 provides 32-bit unsigned bitwise operations for Lua scripts.
//
// gopher-lua implements Lua 5.1, which has no bitwise operators and no `bit`
// or `bit32` library at all — today a glua script cannot AND, OR or shift two
// numbers. This module closes that gap.
//
// Why 32-bit, and why unsigned: Lua numbers are float64. A float64 can only
// represent integers exactly up to 2^53; above that, values silently round.
// A 64-bit bitwise result (e.g. `bor(1<<62, 1)`) would routinely exceed that
// range and corrupt the result without any error — a bitwise library that
// silently loses bits is worse than no bitwise library at all. Lua numbers
// also cannot *carry* a 64-bit mask in the first place, so 64-bit semantics
// would need a boxed integer type, which is out of scope here. Restricting
// results to [0, 2^32) keeps every result exactly representable (2^32 <
// 2^53) and therefore safe by construction. Unsigned (rather than the
// LuaJIT BitOp convention of signed 32-bit) means every result is
// non-negative — `bnot(0)` is `4294967295`, never `-1` — which removes a
// common class of "why is my mask negative" confusion, and matches Lua
// 5.2's own `bit32` library, the naming precedent for this module.
//
// The name is deliberately `bit32`, not `bit`: the width is spelled out so
// nobody assumes 64-bit semantics and gets silently wrong results.
// LuaJIT users coming from `require("bit")` should note the different name
// and the unsigned (rather than signed) result convention.
//
// 64-bit masks, rotations and bit32.extract/replace are out of scope — see
// SPECS.md S1 for the full rationale.
package bit32

import (
	"fmt"
	"math"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// maxSafeInt: the largest magnitude (2^53) an integer can have and still be
// represented exactly by a float64. Operands outside [-maxSafeInt,
// maxSafeInt] are rejected before any bitwise reduction: silently truncating
// them via a raw int64 conversion would let a caller believe a value it never
// actually held was ANDed/ORed/XORed.
const maxSafeInt = 1 << 53

// toUint32: validates that x is an exact, finite integer within float64's
// safe range, then reduces it to its 32-bit two's-complement representation
// modulo 2^32. Negative values wrap the way a two's-complement CPU register
// would: -1 becomes 0xFFFFFFFF, so `band(-1, 0xFF)` is 255.
func toUint32(x float64) (uint32, error) {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return 0, fmt.Errorf("bit32: operand must be a finite integer, got %v", x)
	}
	if x != math.Trunc(x) {
		return 0, fmt.Errorf("bit32: operand must be an integer, got %v", x)
	}
	if x < -maxSafeInt || x > maxSafeInt {
		return 0, fmt.Errorf("bit32: operand %v exceeds the exact integer range of a float64 (±2^53)", x)
	}
	// Two's-complement reduction modulo 2^32: reinterpret the int64 bit
	// pattern as uint64 and keep the low 32 bits.
	return uint32(uint64(int64(x)) & 0xFFFFFFFF), nil //nolint:gosec
}

// shiftCount: validates a shift amount n for lshift/rshift/arshift. Shift
// counts must be non-negative exact integers. There is deliberately no upper
// bound check here (unlike bitIndex) — lshift/rshift/arshift explicitly
// define behaviour for n >= 32 by clamping the result instead of raising.
// Negative counts always raise: silently shifting the other way (as Lua
// 5.2's bit32 does) is a bug magnet, and an explicit call to the opposite
// shift function is clearer.
func shiftCount(n float64) (int, error) {
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("bit32: shift count must be a finite integer, got %v", n)
	}
	if n != math.Trunc(n) {
		return 0, fmt.Errorf("bit32: shift count must be an integer, got %v", n)
	}
	if n < 0 {
		return 0, fmt.Errorf("bit32: shift count must not be negative, got %v", n)
	}
	return int(n), nil
}

// bitIndex: validates a bit position n for test/set/clear. Unlike shift
// counts, a position outside a 32-bit word is meaningless (there is no
// analogous "clamp" behaviour), so it raises rather than silently doing
// nothing.
func bitIndex(n float64) (int, error) {
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, fmt.Errorf("bit32: bit index must be a finite integer, got %v", n)
	}
	if n != math.Trunc(n) {
		return 0, fmt.Errorf("bit32: bit index must be an integer, got %v", n)
	}
	if n < 0 || n > 31 {
		return 0, fmt.Errorf("bit32: bit index must be in [0, 31], got %v", n)
	}
	return int(n), nil
}

// band: bitwise AND of x and every value in rest, reduced modulo 2^32.
func band(x float64, rest ...float64) (float64, error) {
	acc, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	for _, r := range rest {
		ru, err := toUint32(r)
		if err != nil {
			return 0, err
		}
		acc &= ru
	}
	return float64(acc), nil
}

// bor: bitwise OR of x and every value in rest, reduced modulo 2^32.
func bor(x float64, rest ...float64) (float64, error) {
	acc, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	for _, r := range rest {
		ru, err := toUint32(r)
		if err != nil {
			return 0, err
		}
		acc |= ru
	}
	return float64(acc), nil
}

// bxor: bitwise XOR of x and every value in rest, reduced modulo 2^32.
func bxor(x float64, rest ...float64) (float64, error) {
	acc, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	for _, r := range rest {
		ru, err := toUint32(r)
		if err != nil {
			return 0, err
		}
		acc ^= ru
	}
	return float64(acc), nil
}

// bnot: bitwise NOT of x, reduced modulo 2^32. bnot(0) is 4294967295, not -1
// — results are always unsigned.
func bnot(x float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	return float64(^ux), nil
}

// lshift: logical left shift of x by n bits. n >= 32 yields 0.
func lshift(x, n float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	sc, err := shiftCount(n)
	if err != nil {
		return 0, err
	}
	if sc >= 32 {
		return 0, nil
	}
	return float64(ux << uint(sc)), nil //nolint:gosec
}

// rshift: logical right shift of x by n bits — the vacated high bits are
// always filled with zero, regardless of x's sign bit. n >= 32 yields 0.
func rshift(x, n float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	sc, err := shiftCount(n)
	if err != nil {
		return 0, err
	}
	if sc >= 32 {
		return 0, nil
	}
	return float64(ux >> uint(sc)), nil //nolint:gosec
}

// arshift: arithmetic right shift of x by n bits — the vacated high bits are
// filled with a copy of bit 31 (the sign bit), so a value with the high bit
// set stays "negative" as it shifts. This is what makes it genuinely differ
// from rshift: arshift(0x80000000, 4) is 0xF8000000, while rshift(0x80000000,
// 4) is 0x08000000. n >= 32 yields 0xFFFFFFFF if bit 31 was set, else 0.
func arshift(x, n float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	sc, err := shiftCount(n)
	if err != nil {
		return 0, err
	}
	signed := int32(ux) //nolint:gosec // reinterpret bits as signed for sign propagation
	if sc >= 32 {
		if signed < 0 {
			return float64(uint32(0xFFFFFFFF)), nil
		}
		return 0, nil
	}
	// Go's >> on a signed integer is an arithmetic (sign-extending) shift.
	result := signed >> uint(sc)        //nolint:gosec
	return float64(uint32(result)), nil //nolint:gosec
}

// test: reports whether bit n (0 = least significant) of x is set.
func test(x, n float64) (bool, error) {
	ux, err := toUint32(x)
	if err != nil {
		return false, err
	}
	idx, err := bitIndex(n)
	if err != nil {
		return false, err
	}
	return (ux>>uint(idx))&1 == 1, nil //nolint:gosec
}

// setBit: returns x with bit n set. Named setBit rather than set purely to
// stay consistent with clearBit below (set itself is not a Go builtin, but
// pairing the two names makes the module easier to scan).
func setBit(x, n float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	idx, err := bitIndex(n)
	if err != nil {
		return 0, err
	}
	return float64(ux | (1 << uint(idx))), nil //nolint:gosec
}

// clearBit: returns x with bit n cleared. Named clearBit, not clear, because
// clear is a Go builtin (clear(map|slice)) and shadowing it trips
// revive's redefines-builtin-id lint rule.
func clearBit(x, n float64) (float64, error) {
	ux, err := toUint32(x)
	if err != nil {
		return 0, err
	}
	idx, err := bitIndex(n)
	if err != nil {
		return 0, err
	}
	return float64(ux &^ (1 << uint(idx))), nil //nolint:gosec
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("bit32", "32-bit unsigned bitwise operations; see the package doc for why 32-bit and unsigned")

	m.Fn("band", band, "bitwise AND of x and every value in rest, reduced modulo 2^32",
		luareg.Args("x", "..."),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)"),
		luareg.ArgDoc("...", "additional operands, same constraints as x"),
		luareg.ReturnDoc(0, "result", "the AND of all operands, in [0, 2^32)"))
	m.Fn("bor", bor, "bitwise OR of x and every value in rest, reduced modulo 2^32",
		luareg.Args("x", "..."),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)"),
		luareg.ArgDoc("...", "additional operands, same constraints as x"),
		luareg.ReturnDoc(0, "result", "the OR of all operands, in [0, 2^32)"))
	m.Fn("bxor", bxor, "bitwise XOR of x and every value in rest, reduced modulo 2^32",
		luareg.Args("x", "..."),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)"),
		luareg.ArgDoc("...", "additional operands, same constraints as x"),
		luareg.ReturnDoc(0, "result", "the XOR of all operands, in [0, 2^32)"))
	m.Fn("bnot", bnot, "bitwise NOT of x, reduced modulo 2^32",
		luareg.Args("x"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ReturnDoc(0, "result", "the bitwise complement of x, in [0, 2^32); bnot(0) is 4294967295, never -1"))
	m.Fn("lshift", lshift, "logical left shift of x by n bits",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "shift amount; must be a non-negative integer. n >= 32 yields 0"),
		luareg.ReturnDoc(0, "result", "x shifted left by n bits, in [0, 2^32)"))
	m.Fn("rshift", rshift, "logical (zero-filling) right shift of x by n bits",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "shift amount; must be a non-negative integer. n >= 32 yields 0"),
		luareg.ReturnDoc(0, "result", "x shifted right by n bits with zero fill, in [0, 2^32)"))
	m.Fn("arshift", arshift, "arithmetic (sign-propagating) right shift of x by n bits",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "shift amount; must be a non-negative integer. n >= 32 yields 0xFFFFFFFF if bit 31 of x was set, else 0"),
		luareg.ReturnDoc(0, "result", "x shifted right by n bits with bit-31 (sign) fill, in [0, 2^32)"))
	m.Fn("test", test, "reports whether bit n of x is set",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "bit position, least significant bit is 0; must be in [0, 31]"),
		luareg.ReturnDoc(0, "set", "true if bit n of x is 1"))
	m.Fn("set", setBit, "returns x with bit n set",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "bit position, least significant bit is 0; must be in [0, 31]"),
		luareg.ReturnDoc(0, "result", "x with bit n forced to 1, in [0, 2^32)"))
	m.Fn("clear", clearBit, "returns x with bit n cleared",
		luareg.Args("x", "n"),
		luareg.ArgDoc("x", "an exact integer in [-2^53, 2^53]"),
		luareg.ArgDoc("n", "bit position, least significant bit is 0; must be in [0, 31]"),
		luareg.ReturnDoc(0, "result", "x with bit n forced to 0, in [0, 2^32)"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("bit32", bit32.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
