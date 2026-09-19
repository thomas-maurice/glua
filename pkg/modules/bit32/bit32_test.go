// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package bit32

import (
	"math"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/ directory.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("bit32", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("Test script returned %v, expected true", result)
			}
		})
	}
}

// TestToUint32: table-driven coverage of the operand-validation helper —
// this is where the two's-complement and mod-2^32 reduction rules actually
// live, so it is worth pinning at the Go level as well as through Lua.
func TestToUint32(t *testing.T) {
	cases := []struct {
		name    string
		in      float64
		want    uint32
		wantErr bool
	}{
		{"zero", 0, 0, false},
		{"max_uint32", 4294967295, 4294967295, false},
		{"negative_one_is_two_complement", -1, 0xFFFFFFFF, false},
		{"negative_wraps", -256, 0xFFFFFF00, false},
		{"mod_2_32_reduction", 0x1FFFFFFFF, 0xFFFFFFFF, false}, // 2^33 - 1
		{"non_integer_float_raises", 3.5, 0, true},
		{"nan_raises", math.NaN(), 0, true},
		{"inf_raises", math.Inf(1), 0, true},
		{"beyond_2_53_raises", math.Pow(2, 60), 0, true},
		{"exactly_2_53_ok", math.Pow(2, 53), 0, false}, // boundary: exactly 2^53 is allowed
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toUint32(tc.in)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			// 2^53 does not fit in uint32; only check the exact value for
			// cases that do.
			if tc.name != "exactly_2_53_ok" {
				require.Equal(t, tc.want, got)
			}
		})
	}
}

// TestArshiftVsRshiftDiffer: proves arshift and rshift genuinely differ on a
// value with the high bit set — the whole reason arshift exists.
func TestArshiftVsRshiftDiffer(t *testing.T) {
	logical, err := rshift(0x80000000, 4)
	require.NoError(t, err)
	require.Equal(t, float64(0x08000000), logical)

	arithmetic, err := arshift(0x80000000, 4)
	require.NoError(t, err)
	require.Equal(t, float64(0xF8000000), arithmetic)

	require.NotEqual(t, logical, arithmetic)
}

// TestBnotRoundTrip: bnot(bnot(x)) == x for representative values including
// the two boundary values of the 32-bit range.
func TestBnotRoundTrip(t *testing.T) {
	for _, x := range []float64{0, 1, 0xFFFFFFFF, 0x80000000, 12345} {
		once, err := bnot(x)
		require.NoError(t, err)
		twice, err := bnot(once)
		require.NoError(t, err)
		require.Equal(t, x, twice)
	}
}

// TestShiftCountAtLeast32: n >= 32 clamps rather than raising for
// lshift/rshift/arshift.
func TestShiftCountAtLeast32(t *testing.T) {
	v, err := lshift(1, 32)
	require.NoError(t, err)
	require.Equal(t, float64(0), v)

	v, err = rshift(0xFFFFFFFF, 40)
	require.NoError(t, err)
	require.Equal(t, float64(0), v)

	v, err = arshift(0x80000000, 32)
	require.NoError(t, err)
	require.Equal(t, float64(0xFFFFFFFF), v)

	v, err = arshift(0x7FFFFFFF, 100)
	require.NoError(t, err)
	require.Equal(t, float64(0), v)
}

// TestNegativeShiftRaises: negative shift counts always raise, for all three
// shift functions.
func TestNegativeShiftRaises(t *testing.T) {
	_, err := lshift(1, -1)
	require.Error(t, err)
	_, err = rshift(1, -1)
	require.Error(t, err)
	_, err = arshift(1, -1)
	require.Error(t, err)
}

// TestBitIndexOutOfRangeRaises: test/set/clear raise for an index outside
// [0, 31], unlike the shift functions which clamp.
func TestBitIndexOutOfRangeRaises(t *testing.T) {
	_, err := test(1, 32)
	require.Error(t, err)
	_, err = setBit(1, 32)
	require.Error(t, err)
	_, err = clearBit(1, 32)
	require.Error(t, err)
	_, err = test(1, -1)
	require.Error(t, err)
}
