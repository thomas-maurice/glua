// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package compress

import (
	"bytes"
	"math"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/modules/base64"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/. Both compress and
// base64 are preloaded so scripts can exercise the composition the package
// doc promises.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("compress", Loader)
			L.PreloadModule("base64", base64.Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("test script returned %v, expected true", result)
			}
		})
	}
}

// codecs enumerates the three codec function pairs under test, so the
// round-trip/level tests below run identically over all of them instead of
// being copy-pasted three times.
type codec struct {
	name       string
	compress   func([]byte, int) ([]byte, error)
	decompress func([]byte, float64) ([]byte, error)
}

var codecs = []codec{
	{"gzip", gzipCompress, gzipDecompress},
	{"zlib", zlibCompress, zlibDecompress},
	{"flate", flateCompress, flateDecompress},
}

// TestRoundTrip_AllLevels: every codec must round-trip ASCII, all 256 byte
// values (the 8-bit-cleanliness proof) and the empty string, at every named
// compression level. This is the test that would catch a codec silently
// corrupting binary data or mishandling a boundary level.
func TestRoundTrip_AllLevels(t *testing.T) {
	full256 := make([]byte, 256)
	for i := range full256 {
		full256[i] = byte(i)
	}

	inputs := map[string][]byte{
		"ascii":  []byte("the quick brown fox jumps over the lazy dog"),
		"binary": full256,
		"empty":  []byte(""),
	}

	levels := []int{huffmanOnly, defaultCompression, noCompression, bestSpeed, bestCompression}

	for _, c := range codecs {
		for inputName, data := range inputs {
			for _, level := range levels {
				t.Run(c.name+"/"+inputName+"/level", func(t *testing.T) {
					compressed, err := c.compress(data, level)
					require.NoError(t, err)

					decompressed, err := c.decompress(compressed, maxBytesDefault)
					require.NoError(t, err)

					assert.True(t, bytes.Equal(decompressed, data),
						"round-trip mismatch for %s at level %d", c.name, level)
				})
			}
		}
	}
}

// TestInvalidLevel_Raises: a level outside [-2, 9] must raise a useful
// error for every codec's compress function.
func TestInvalidLevel_Raises(t *testing.T) {
	for _, c := range codecs {
		for _, level := range []int{-3, 10, 100} {
			_, err := c.compress([]byte("data"), level)
			assert.Error(t, err, "%s: expected error for level %d", c.name, level)
		}
	}
}

// TestDecompress_MaxBytesLessThanOne_Raises: max_bytes < 1 must always
// raise, naming the limit — there is deliberately no "0 means unlimited"
// escape hatch (SPECS.md S7).
func TestDecompress_MaxBytesLessThanOne_Raises(t *testing.T) {
	blob, err := gzipCompress([]byte("hello"), defaultCompression)
	require.NoError(t, err)

	for _, c := range codecs {
		for _, maxBytes := range []float64{0, -1, -100} {
			_, err := c.decompress(blob, maxBytes)
			assert.Error(t, err, "%s: expected error for max_bytes=%v", c.name, maxBytes)
		}
	}
}

// TestDecompress_MaxBytesOutOfRange_Raises is the exact regression this
// chunk fixes (security review MEDIUM 3): math.huge (gopher-lua's
// math.huge is math.MaxFloat64), a value near MaxInt64, and NaN must all be
// rejected explicitly rather than silently truncating during the
// float->int64 conversion and producing an empty result with no error.
// Confirmed before the fix: math.huge made a 100 KB gzip payload decompress
// to length 0 with err == nil, because int64(math.MaxFloat64)+1 overflows
// (saturating differently per architecture) and io.CopyN with a resulting
// negative count reads nothing.
func TestDecompress_MaxBytesOutOfRange_Raises(t *testing.T) {
	const payload = "a payload that must not silently vanish"
	blob, err := gzipCompress([]byte(payload), defaultCompression)
	require.NoError(t, err)

	cases := map[string]float64{
		"math.MaxFloat64 (gopher-lua's math.huge)": math.MaxFloat64,
		"near MaxInt64": 9.2e18,
		"NaN":           math.NaN(),
		"+Inf":          math.Inf(1),
	}
	for _, c := range codecs {
		for name, maxBytes := range cases {
			t.Run(c.name+"/"+name, func(t *testing.T) {
				out, err := c.decompress(blob, maxBytes)
				require.Error(t, err, "max_bytes=%v must be rejected, not silently return a truncated/empty result", maxBytes)
				assert.Nil(t, out, "no partial output should be returned alongside the error")
			})
		}
	}
}

// TestDecompress_MaxBytesNearCeiling_Works: a normal large-but-valid limit
// (well above the 64 MiB default, at the module's ceiling) must still work
// -- the fix must reject only what is unsafe, not shrink the usable range.
func TestDecompress_MaxBytesNearCeiling_Works(t *testing.T) {
	blob, err := gzipCompress([]byte("hello, ceiling"), defaultCompression)
	require.NoError(t, err)

	out, err := gzipDecompress(blob, maxBytesCeiling)
	require.NoError(t, err)
	assert.Equal(t, "hello, ceiling", string(out))
}

// TestCrossCodecRejection: feeding a zlib stream to gzip_decompress must
// raise rather than silently misparse — gzip validates its magic header
// before ever touching the DEFLATE payload, so this fails fast.
func TestCrossCodecRejection(t *testing.T) {
	zlibBlob, err := zlibCompress([]byte("hello world"), defaultCompression)
	require.NoError(t, err)

	_, err = gzipDecompress(zlibBlob, maxBytesDefault)
	assert.Error(t, err, "expected gzip_decompress to reject a zlib stream")
}

// TestDecompressionBomb_LimitFires: compresses several MB of zeros (which
// deflate reduces to a tiny stream) and proves that decompressing it with a
// small max_bytes raises rather than allocating the full decompressed size.
// This is the test that justifies the mandatory max_bytes argument: without
// the limit, this call would allocate the full uncompressed size in the
// host Go process regardless of how small the caller's budget was.
func TestDecompressionBomb_LimitFires(t *testing.T) {
	const bombSize = 8 * 1024 * 1024 // 8 MiB of zeros
	zeros := make([]byte, bombSize)

	bomb, err := gzipCompress(zeros, bestCompression)
	require.NoError(t, err)
	// Confirm this really is a bomb: tiny compressed size, huge decompressed
	// size. If this assumption ever stops holding the test would pass for
	// the wrong reason, so pin it explicitly.
	require.Less(t, len(bomb), bombSize/100, "fixture is not actually a decompression bomb")

	const limit = 1024
	_, err = gzipDecompress(bomb, limit)
	require.Error(t, err, "expected the decompression bomb to be rejected")
	assert.Contains(t, err.Error(), "1024", "error should name the limit that fired")
}

// TestEmptyDecompress_NotNilString: a nil []byte at the top level crosses
// into Lua as nil rather than an empty string (pkg/luareg's []byte special
// case). Decompressing an empty payload is a real case (compressed empty
// input) and must come back as an empty Lua string, not nil, or callers
// that do `#result == 0` or string concatenation break.
func TestEmptyDecompress_NotNilString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("compress", Loader)

	code := `
		local compress = require("compress")
		local blob = compress.gzip_compress("", compress.DEFAULT_COMPRESSION)
		local out = compress.gzip_decompress(blob, compress.MAX_BYTES_DEFAULT)
		assert(type(out) == "string", "expected a string, got " .. type(out))
		assert(out == "", "expected empty string, got: " .. tostring(out))
	`
	require.NoError(t, L.DoString(code))
}
