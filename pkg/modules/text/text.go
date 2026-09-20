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

// Package text provides presentation/layout string operations that Go's
// standard library does not offer: word wrapping, indent/dedent, truncation
// with an ellipsis, and rune padding.
//
// Split rule versus the strings module (see SPECS.md, S3): strings binds
// Go's strings package one-to-one; text is for operations Go's stdlib does
// not provide at all. If pkg/strings already has it, it belongs in strings,
// not here.
//
// IMPORTANT — every function in this module is rune-oriented, not
// display-width-oriented. "Width" and "length" below always mean a count of
// Unicode code points (runes), obtained by counting []rune(s) — never bytes,
// and never terminal display cells. This means:
//
//   - A CJK string (each rune commonly rendering two terminal columns) will
//     NOT align in a monospace terminal even though wrap/pad/truncate agree
//     it is "N runes wide".
//   - Combining-character sequences and multi-rune emoji (which are commonly
//     one visual glyph made of several code points) are NOT collapsed to one
//     unit of width.
//
// This is a deliberate scope decision, not an oversight: getting display
// width right needs a dependency like mattn/go-runewidth, which is present
// transitively in this module's graph but not promoted to direct for this
// batch. Revisit only if a real caller hits the misalignment.
package text

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// maxPadWidth: the largest width accepted by pad_left/pad_right. Same
// reasoning and same bound as pkg/modules/strings' maxRepCount and
// pkg/modules/random's maxLen: generous enough for any realistic layout
// use, small enough that a typo'd width argument cannot turn into an
// unbounded allocation that OOMs the host Go process.
const maxPadWidth = 1 << 20

// wrapParagraph: greedily word-wraps a single line (no embedded newlines) to
// at most width runes per output line. Words are whitespace-separated
// (strings.Fields semantics); a single word longer than width is placed on
// its own line and is NOT split — it overflows. An input with no words
// (empty or all-whitespace) produces exactly one empty output line, so blank
// lines round-trip through wrap unchanged.
func wrapParagraph(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}

	lines := make([]string, 0, 1)
	var cur []string
	curLen := 0
	for _, w := range words {
		wLen := utf8.RuneCountInString(w)
		switch {
		case len(cur) == 0:
			cur = append(cur, w)
			curLen = wLen
		case curLen+1+wLen > width:
			lines = append(lines, strings.Join(cur, " "))
			cur = []string{w}
			curLen = wLen
		default:
			cur = append(cur, w)
			curLen += 1 + wLen
		}
	}
	lines = append(lines, strings.Join(cur, " "))
	return lines
}

// wrap: greedily word-wraps s to at most width runes per line. Existing "\n"
// characters in s are preserved as hard paragraph breaks — each paragraph is
// wrapped independently, so a blank line in the input produces a blank line
// in the output. A single word longer than width is not broken; it overflows
// its line rather than being hyphenated or truncated. Raises if width < 1,
// since a non-positive width has no sensible greedy-wrap behaviour (every
// word, even a one-rune one, would "overflow").
func wrap(s string, width int) (string, error) {
	if width < 1 {
		return "", fmt.Errorf("text.wrap: width must be >= 1, got %d", width)
	}

	paragraphs := strings.Split(s, "\n")
	var outLines []string
	for _, p := range paragraphs {
		outLines = append(outLines, wrapParagraph(p, width)...)
	}
	return strings.Join(outLines, "\n"), nil
}

// indent: prefixes every line of s with prefix. A trailing empty line — i.e.
// the empty string produced after the final "\n" when s ends with a newline
// — is deliberately NOT prefixed, so indenting a string and then comparing
// "does it end with a newline" still works the same way it did before
// indenting. Never raises.
func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// leadingWhitespace: returns the run of spaces/tabs at the start of s.
// Comparison is byte-literal — a tab and a space are different characters,
// never expanded or normalized against each other.
func leadingWhitespace(s string) string {
	i := 0
	for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
		i++
	}
	return s[:i]
}

// commonPrefix: returns the longest literal byte-prefix shared by a and b.
func commonPrefix(a, b string) string {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return a[:i]
}

// dedent: removes the longest common leading-whitespace prefix shared by
// every non-blank line of s. A line is "blank" if it is empty or contains
// only whitespace; blank lines do not participate in computing the common
// prefix, and are always normalized to the empty string in the output
// (matching Python's textwrap.dedent — a line of pure whitespace carries no
// useful indentation information and is not worth preserving as noise).
// Leading whitespace is compared byte-literally: tabs and spaces are
// distinct characters, never expanded or treated as equivalent. Never
// raises — a string with no non-blank lines, or with no common prefix
// (margin == ""), is returned with blank lines normalized and non-blank
// lines unchanged.
func dedent(s string) string {
	lines := strings.Split(s, "\n")

	var margin string
	haveMargin := false
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		prefix := leadingWhitespace(line)
		if !haveMargin {
			margin = prefix
			haveMargin = true
			continue
		}
		margin = commonPrefix(margin, prefix)
	}

	out := make([]string, len(lines))
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			out[i] = ""
			continue
		}
		out[i] = strings.TrimPrefix(line, margin)
	}
	return strings.Join(out, "\n")
}

// truncate: returns s unchanged if it is at most width runes long. Otherwise
// returns the first (width - runelen(ellipsis)) runes of s followed by
// ellipsis, so the result is always exactly width runes long. Raises if
// width is smaller than ellipsis's own rune length, since there would be no
// room left for any content plus the ellipsis.
func truncate(s string, width int, ellipsis string) (string, error) {
	ellipsisLen := utf8.RuneCountInString(ellipsis)
	if width < ellipsisLen {
		return "", fmt.Errorf("text.truncate: width (%d) must be >= the rune length of ellipsis (%d)", width, ellipsisLen)
	}

	runes := []rune(s)
	if len(runes) <= width {
		return s, nil
	}
	keep := width - ellipsisLen
	return string(runes[:keep]) + ellipsis, nil
}

// singlePadRune: validates that pad is exactly one rune and returns it.
func singlePadRune(fnName, pad string) (rune, error) {
	padRunes := []rune(pad)
	if len(padRunes) != 1 {
		return 0, fmt.Errorf("text.%s: pad must be exactly one rune, got %d", fnName, len(padRunes))
	}
	return padRunes[0], nil
}

// validatePadWidth: rejects a width argument that could turn pad_left/
// pad_right into an unbounded allocation (e.g. a typo'd width meant to be a
// column count). fnName is the Lua-visible function name, used to make the
// error greppable back to its call site.
func validatePadWidth(fnName string, width int) error {
	if width > maxPadWidth {
		return fmt.Errorf("text.%s: width must be <= %d, got %d", fnName, maxPadWidth, width)
	}
	return nil
}

// padLeft: left-pads s with copies of pad until it is width runes long.
// Never truncates: if s is already width runes or wider, it is returned
// unchanged. Raises if pad is not exactly one rune, or if width exceeds
// maxPadWidth.
func padLeft(s string, width int, pad string) (string, error) {
	r, err := singlePadRune("pad_left", pad)
	if err != nil {
		return "", err
	}
	if err := validatePadWidth("pad_left", width); err != nil {
		return "", err
	}
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s, nil
	}
	return strings.Repeat(string(r), width-n) + s, nil
}

// padRight: right-pads s with copies of pad until it is width runes long.
// Never truncates: if s is already width runes or wider, it is returned
// unchanged. Raises if pad is not exactly one rune, or if width exceeds
// maxPadWidth.
func padRight(s string, width int, pad string) (string, error) {
	r, err := singlePadRune("pad_right", pad)
	if err != nil {
		return "", err
	}
	if err := validatePadWidth("pad_right", width); err != nil {
		return "", err
	}
	n := utf8.RuneCountInString(s)
	if n >= width {
		return s, nil
	}
	return s + strings.Repeat(string(r), width-n), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("text", "presentation/layout string operations not provided by Go's stdlib: wrap, indent, dedent, truncate, pad")
	m.Fn("wrap", wrap, "greedily word-wraps a string to a maximum rune width per line",
		luareg.Args("s", "width"),
		luareg.ArgDoc("s", "the string to wrap; existing newlines are preserved as hard paragraph breaks and each paragraph wraps independently"),
		luareg.ArgDoc("width", "maximum runes per output line, must be >= 1; a single word longer than width is not split and overflows its line"),
		luareg.ReturnDoc(0, "out", "s re-flowed to width runes per line, using rune counts, not display width"))
	m.Fn("indent", indent, "prefixes every line of a string",
		luareg.Args("s", "prefix"),
		luareg.ArgDoc("s", "the string to indent"),
		luareg.ArgDoc("prefix", "the string prepended to every line; a trailing empty line produced by a final newline in s is not prefixed"),
		luareg.ReturnDoc(0, "out", "s with prefix prepended to each line"))
	m.Fn("dedent", dedent, "removes the common leading-whitespace prefix shared by every non-blank line",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to dedent; leading whitespace is compared byte-literally (tabs and spaces are distinct, never expanded), and blank/whitespace-only lines do not participate in computing the common prefix and are normalized to empty in the output"),
		luareg.ReturnDoc(0, "out", "s with the longest common leading whitespace removed from every non-blank line"))
	m.Fn("truncate", truncate, "truncates a string to a maximum rune width, appending an ellipsis if it was shortened",
		luareg.Args("s", "width", "ellipsis"),
		luareg.ArgDoc("s", "the string to truncate; rune-counted, not display-width-counted"),
		luareg.ArgDoc("width", "maximum runes of the result, including the ellipsis; must be >= the rune length of ellipsis"),
		luareg.ArgDoc("ellipsis", "the marker appended when s is shortened, e.g. \"...\" or the single rune \"…\"; required, there is no default"),
		luareg.ReturnDoc(0, "out", "s unchanged if it already fits in width runes, otherwise a prefix of s plus ellipsis, exactly width runes long"))
	m.Fn("pad_left", padLeft, "left-pads a string with a single rune to a minimum rune width",
		luareg.Args("s", "width", "pad"),
		luareg.ArgDoc("s", "the string to pad"),
		luareg.ArgDoc("width", "minimum runes of the result; s is never truncated if it is already this wide or wider; must be <= 1048576"),
		luareg.ArgDoc("pad", "the single rune to pad with; exactly one rune is required"),
		luareg.ReturnDoc(0, "out", "s left-padded with pad to width runes, or s unchanged if already >= width runes"))
	m.Fn("pad_right", padRight, "right-pads a string with a single rune to a minimum rune width",
		luareg.Args("s", "width", "pad"),
		luareg.ArgDoc("s", "the string to pad"),
		luareg.ArgDoc("width", "minimum runes of the result; s is never truncated if it is already this wide or wider; must be <= 1048576"),
		luareg.ArgDoc("pad", "the single rune to pad with; exactly one rune is required"),
		luareg.ReturnDoc(0, "out", "s right-padded with pad to width runes, or s unchanged if already >= width runes"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("text", text.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
