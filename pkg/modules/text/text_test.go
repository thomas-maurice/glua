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

package text

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestDedentCommonPrefix: table-driven coverage of the margin computation
// that dedent relies on — this is where the "shortest common indentation,
// literal byte comparison, blank lines excluded" rules actually live.
func TestDedentCommonPrefix(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"uniform_4_space_margin", "    a\n    b\n", "a\nb\n"},
		{"shortest_wins", "    a\n      b\n    c", "a\n  b\nc"},
		{"no_common_prefix_is_noop", "a\n b", "a\n b"},
		{"blank_lines_excluded_from_margin_and_normalized", "    a\n\t\n    b", "a\n\nb"},
		{"tabs_and_spaces_are_literal_not_expanded", "\ta\n    b", "\ta\n    b"},
		{"single_line_no_margin", "hello", "hello"},
		{"all_blank_normalizes_every_line", "  \n\t\n", "\n\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, dedent(tc.in))
		})
	}
}

// TestWrapParagraphLongWordOverflows: a word longer than width occupies its
// own line unmodified rather than being split — the specific behaviour the
// chunk spec calls out as a hard requirement.
func TestWrapParagraphLongWordOverflows(t *testing.T) {
	lines := wrapParagraph("a supercalifragilisticexpialidocious b", 5)
	require.Equal(t, []string{"a", "supercalifragilisticexpialidocious", "b"}, lines)
}

// TestWrapParagraphEmptyProducesOneBlankLine: an empty or all-whitespace
// paragraph must still produce exactly one output line (empty), so that
// blank lines in the input round-trip through wrap rather than disappearing.
func TestWrapParagraphEmptyProducesOneBlankLine(t *testing.T) {
	require.Equal(t, []string{""}, wrapParagraph("", 10))
	require.Equal(t, []string{""}, wrapParagraph("   ", 10))
}

// TestTruncateBoundary: width exactly equal to the input's rune length
// leaves it unchanged (no ellipsis appended); width exactly equal to the
// ellipsis's rune length returns just the ellipsis.
func TestTruncateBoundary(t *testing.T) {
	out, err := truncate("hello", 5, "...")
	require.NoError(t, err)
	require.Equal(t, "hello", out)

	out, err = truncate("hello world", 3, "...")
	require.NoError(t, err)
	require.Equal(t, "...", out)

	_, err = truncate("hello world", 2, "...")
	require.Error(t, err)
}

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

			L.PreloadModule("text", Loader)

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
