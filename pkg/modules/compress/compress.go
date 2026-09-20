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

// Package compress provides gzip, zlib and raw DEFLATE compression and
// decompression for Lua scripts, over 8-bit-clean strings ([]byte at the top
// level maps to a raw Lua string, not base64 — see pkg/luareg's F3 rule), so
// it composes directly with the base64 module:
// base64.encode(compress.gzip_compress(s, compress.BEST_COMPRESSION)).
//
// Every decompress function takes a REQUIRED max_bytes argument, with
// deliberately no "0 means unlimited" escape hatch. glua is embedded in
// admission controllers and policy engines that routinely decompress input
// from outside the process (a base64 annotation, an HTTP body, a ConfigMap).
// An unbounded read of a decompression stream turns a small input into a
// decompression bomb that OOMs the *host Go process* — strictly worse than a
// Lua error, which a caller can pcall and handle. The limit is enforced by
// reading at most max_bytes+1 bytes from the decompressor, so memory use is
// bounded by the limit regardless of how large the compressed input claims
// to decompress to. max_bytes is also capped at maxBytesCeiling (4 GiB) and
// must be finite: an unvalidated max_bytes = math.huge (gopher-lua's
// math.huge is math.MaxFloat64) previously overflowed the internal
// max_bytes+1 arithmetic and produced an empty result with NO error --
// silently worse than either a correct decompress or a raised error.
package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// Compression level constants, mirrored as module constants below. All three
// codecs are backed by compress/flate, so the accepted range is identical:
// -2 (HuffmanOnly) through 9 (BestCompression), with -1 as the default and 0
// meaning "framing only, no compression".
const (
	noCompression      = 0
	bestSpeed          = 1
	bestCompression    = 9
	defaultCompression = -1
	huffmanOnly        = -2

	// maxBytesDefault: a generous but finite default ceiling for decompress
	// calls, exposed as compress.MAX_BYTES_DEFAULT so the common call site
	// is only one token longer than an unbounded call would have been.
	maxBytesDefault = 64 * 1024 * 1024

	// maxBytesCeiling: the largest max_bytes a caller may request. Chosen
	// far above any legitimate decompress budget (4 GiB vs. the 64 MiB
	// default) while staying many orders of magnitude below MaxInt64, so
	// drainLimited's "maxBytes+1" arithmetic can never overflow regardless
	// of what a caller passes here.
	maxBytesCeiling = 4 * 1024 * 1024 * 1024
)

// validateLevel: rejects a compression level outside the range every codec
// in this module accepts. fnName is the Lua-visible function name, used to
// make the error greppable back to its call site.
func validateLevel(fnName string, level int) error {
	if level < huffmanOnly || level > bestCompression {
		return fmt.Errorf("%s: invalid level %d (want %d..%d)", fnName, level, huffmanOnly, bestCompression)
	}
	return nil
}

// validateMaxBytes: rejects a max_bytes argument that would reopen the
// unbounded-decompression footgun, or that cannot be safely turned into a
// bounded byte count. There is deliberately no "0 means unlimited" escape
// hatch.
//
// maxBytes is taken as float64 -- the raw Lua number, BEFORE any conversion
// to an integer type -- and validated as a float, not converted first. This
// matters because Go's float64->int64 conversion for an out-of-range value
// is architecture-defined: it saturates to MaxInt64 on arm64 and produces
// MinInt64 on amd64. gopher-lua's math.huge is math.MaxFloat64 (~1.8e308),
// which is exactly such an out-of-range value once multiplied through by
// this function's ceiling math -- if this validated an already-converted
// int, the check itself would silently differ between arm64 and amd64
// (which is what CI runs, so it would not have caught the arm64 bug).
// Confirmed: passing math.huge as max_bytes previously made
// gzip_decompress return an empty string with NO error for a 100 KB
// payload, because int64(maxBytes)+1 overflowed negative and io.CopyN with
// a negative count reads nothing.
func validateMaxBytes(fnName string, maxBytes float64) (int64, error) {
	if math.IsNaN(maxBytes) || math.IsInf(maxBytes, 0) {
		return 0, fmt.Errorf("%s: max_bytes must be a finite number, got %v", fnName, maxBytes)
	}
	if maxBytes < 1 || maxBytes > maxBytesCeiling {
		return 0, fmt.Errorf("%s: max_bytes must be in [1, %d], got %v", fnName, maxBytesCeiling, maxBytes)
	}
	return int64(maxBytes), nil
}

// nonNilBytes: a nil []byte at the top level crosses into Lua as nil, not an
// empty string (pkg/luareg's []byte special case mirrors Go's json.Marshal
// nil-slice behaviour). bytes.Buffer.Bytes() returns nil for a buffer that
// was never written to, which happens whenever a decompressed payload is
// legitimately empty. Normalise so an empty result is an empty Lua string,
// not nil.
func nonNilBytes(b []byte) []byte {
	if b == nil {
		return []byte{}
	}
	return b
}

// drainLimited: copies at most maxBytes+1 bytes from r into memory and
// raises if more than maxBytes were available. It never buffers more than
// maxBytes+1 bytes, which is what makes it safe against a decompression
// bomb regardless of how large the compressed input claims to decompress
// to. fnName is used to prefix any error raised. maxBytes must already be
// validated (validateMaxBytes) to be within [1, maxBytesCeiling], so
// maxBytes+1 cannot overflow int64 here.
func drainLimited(fnName string, r io.Reader, maxBytes int64) ([]byte, error) {
	var buf bytes.Buffer
	n, err := io.CopyN(&buf, r, maxBytes+1)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%s: %w", fnName, err)
	}
	if n > maxBytes {
		return nil, fmt.Errorf("%s: output exceeds max_bytes (%d)", fnName, maxBytes)
	}
	return nonNilBytes(buf.Bytes()), nil
}

// gzipCompress: compresses data as a gzip stream at the given level.
func gzipCompress(data []byte, level int) ([]byte, error) {
	if err := validateLevel("compress.gzip_compress", level); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("compress.gzip_compress: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("compress.gzip_compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress.gzip_compress: %w", err)
	}
	return nonNilBytes(buf.Bytes()), nil
}

// gzipDecompress: decompresses a gzip stream, raising if the decompressed
// output would exceed maxBytes.
func gzipDecompress(data []byte, maxBytes float64) ([]byte, error) {
	mb, err := validateMaxBytes("compress.gzip_decompress", maxBytes)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("compress.gzip_decompress: %w", err)
	}
	defer func() { _ = r.Close() }()
	return drainLimited("compress.gzip_decompress", r, mb)
}

// zlibCompress: compresses data as a zlib stream at the given level.
func zlibCompress(data []byte, level int) ([]byte, error) {
	if err := validateLevel("compress.zlib_compress", level); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w, err := zlib.NewWriterLevel(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("compress.zlib_compress: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("compress.zlib_compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress.zlib_compress: %w", err)
	}
	return nonNilBytes(buf.Bytes()), nil
}

// zlibDecompress: decompresses a zlib stream, raising if the decompressed
// output would exceed maxBytes.
func zlibDecompress(data []byte, maxBytes float64) ([]byte, error) {
	mb, err := validateMaxBytes("compress.zlib_decompress", maxBytes)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("compress.zlib_decompress: %w", err)
	}
	defer func() { _ = r.Close() }()
	return drainLimited("compress.zlib_decompress", r, mb)
}

// flateCompress: compresses data as a raw DEFLATE stream (no gzip/zlib
// framing) at the given level.
func flateCompress(data []byte, level int) ([]byte, error) {
	if err := validateLevel("compress.flate_compress", level); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, level)
	if err != nil {
		return nil, fmt.Errorf("compress.flate_compress: %w", err)
	}
	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("compress.flate_compress: %w", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("compress.flate_compress: %w", err)
	}
	return nonNilBytes(buf.Bytes()), nil
}

// flateDecompress: decompresses a raw DEFLATE stream, raising if the
// decompressed output would exceed maxBytes. Raw DEFLATE has no header, so
// unlike gzip/zlib, feeding it framed input does not necessarily fail fast —
// it may simply fail (or in rare cases silently misparse) partway through
// the stream. This is an inherent limitation of the format, not a bug here.
func flateDecompress(data []byte, maxBytes float64) ([]byte, error) {
	mb, err := validateMaxBytes("compress.flate_decompress", maxBytes)
	if err != nil {
		return nil, err
	}
	r := flate.NewReader(bytes.NewReader(data))
	defer func() { _ = r.Close() }()
	return drainLimited("compress.flate_decompress", r, mb)
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("compress", "gzip, zlib and raw DEFLATE compression and decompression")

	m.Fn("gzip_compress", gzipCompress, "compresses data as a gzip stream",
		luareg.Args("data", "level"),
		luareg.ArgDoc("data", "the raw bytes to compress"),
		luareg.ArgDoc("level", "compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants"),
		luareg.ReturnDoc(0, "compressed", "the gzip stream"))
	m.Fn("gzip_decompress", gzipDecompress, "decompresses a gzip stream, raising if the output would exceed max_bytes",
		luareg.Args("data", "max_bytes"),
		luareg.ArgDoc("data", "the gzip stream to decompress"),
		luareg.ArgDoc("max_bytes", "the maximum number of decompressed bytes to allow; must be finite and in [1, 4 GiB], no unlimited option"),
		luareg.ReturnDoc(0, "data", "the decompressed bytes"))

	m.Fn("zlib_compress", zlibCompress, "compresses data as a zlib stream",
		luareg.Args("data", "level"),
		luareg.ArgDoc("data", "the raw bytes to compress"),
		luareg.ArgDoc("level", "compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants"),
		luareg.ReturnDoc(0, "compressed", "the zlib stream"))
	m.Fn("zlib_decompress", zlibDecompress, "decompresses a zlib stream, raising if the output would exceed max_bytes",
		luareg.Args("data", "max_bytes"),
		luareg.ArgDoc("data", "the zlib stream to decompress"),
		luareg.ArgDoc("max_bytes", "the maximum number of decompressed bytes to allow; must be finite and in [1, 4 GiB], no unlimited option"),
		luareg.ReturnDoc(0, "data", "the decompressed bytes"))

	m.Fn("flate_compress", flateCompress, "compresses data as a raw DEFLATE stream (no gzip/zlib framing)",
		luareg.Args("data", "level"),
		luareg.ArgDoc("data", "the raw bytes to compress"),
		luareg.ArgDoc("level", "compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants"),
		luareg.ReturnDoc(0, "compressed", "the raw DEFLATE stream"))
	m.Fn("flate_decompress", flateDecompress, "decompresses a raw DEFLATE stream, raising if the output would exceed max_bytes",
		luareg.Args("data", "max_bytes"),
		luareg.ArgDoc("data", "the raw DEFLATE stream to decompress"),
		luareg.ArgDoc("max_bytes", "the maximum number of decompressed bytes to allow; must be finite and in [1, 4 GiB], no unlimited option"),
		luareg.ReturnDoc(0, "data", "the decompressed bytes"))

	m.Const("NO_COMPRESSION", noCompression, "number", "no compression, framing only")
	m.Const("BEST_SPEED", bestSpeed, "number", "fastest compression level")
	m.Const("BEST_COMPRESSION", bestCompression, "number", "smallest output, slowest")
	m.Const("DEFAULT_COMPRESSION", defaultCompression, "number", "the codec's own default level")
	m.Const("HUFFMAN_ONLY", huffmanOnly, "number", "Huffman-only compression, very fast, weak ratio")
	m.Const("MAX_BYTES_DEFAULT", maxBytesDefault, "number", "a generous default max_bytes ceiling (64 MiB)")

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("compress", compress.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
