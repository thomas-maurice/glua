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

package collections

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/, covering the
// Lua-visible behaviour of every function (see the package doc for the
// array-vs-map mode split each script exercises).
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("collections", Loader)
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

// runWithTimeout: runs fn and fails the test if it has not returned within
// timeout. Used to turn "hangs forever" into a deterministic test failure
// instead of relying on go test's global per-package timeout, which is a
// much coarser and slower signal.
func runWithTimeout(t *testing.T, timeout time.Duration, fn func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		fn()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("operation did not complete within %s -- likely an infinite loop on a cyclic table", timeout)
	}
}

// TestDeepEqualCycleDoesNotHang: a table that references itself must not send
// deep_equal into infinite recursion. This pins the design decision to guard
// with a visited (*lua.LTable, *lua.LTable) set rather than depth-limiting or
// simply not handling the case -- an unguarded implementation would hang this
// test forever rather than fail it.
func TestDeepEqualCycleDoesNotHang(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	cyclic := L.NewTable()
	cyclic.RawSetString("self", cyclic)

	var result bool
	runWithTimeout(t, 2*time.Second, func() {
		result = deepEqualFn(cyclic, cyclic)
	})
	require.True(t, result, "a cyclic table must deep_equal itself")

	// Two independently-built, structurally-equal cyclic tables must also
	// compare equal rather than hang.
	other := L.NewTable()
	other.RawSetString("self", other)
	var crossResult bool
	runWithTimeout(t, 2*time.Second, func() {
		crossResult = deepEqualFn(cyclic, other)
	})
	require.True(t, crossResult, "two structurally-equal cyclic tables must compare equal")
}

// TestDeepCopyCycleDoesNotHang: deep_copy on a self-referencing table must
// terminate and preserve the cycle shape (the copy's self-reference points
// back to the copy, not the original, and not a second unbounded expansion).
func TestDeepCopyCycleDoesNotHang(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	cyclic := L.NewTable()
	cyclic.RawSetString("self", cyclic)

	var copyVal lua.LValue
	runWithTimeout(t, 2*time.Second, func() {
		copyVal = deepCopyFn(L, cyclic)
	})

	copyTbl, ok := copyVal.(*lua.LTable)
	require.True(t, ok, "deep_copy of a table must return a table")
	require.NotSame(t, cyclic, copyTbl, "deep_copy must allocate a new table, not alias the original")
	require.Same(t, copyTbl, copyTbl.RawGetString("self"), "deep_copy must preserve the cycle: the copy's self field must point back to the copy itself")
}

// TestLookupPathSegmentNumberBeforeString: pins the documented ambiguity
// rule for get_path -- a segment that parses as an integer is tried as a
// number key first, falling back to a string key only if the number key is
// absent.
func TestLookupPathSegmentNumberBeforeString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	both := L.NewTable()
	both.RawSetInt(1, lua.LString("number-key"))
	both.RawSetString("1", lua.LString("string-key"))
	require.Equal(t, lua.LString("number-key"), lookupPathSegment(both, "1"),
		"a numeric-looking segment must prefer the number key when both exist")

	stringOnly := L.NewTable()
	stringOnly.RawSetString("1", lua.LString("string-key"))
	require.Equal(t, lua.LString("string-key"), lookupPathSegment(stringOnly, "1"),
		"a numeric-looking segment must fall back to the string key when no number key exists")

	require.Equal(t, lua.LNil, lookupPathSegment(L.NewTable(), "missing"),
		"a segment with no corresponding key must resolve to nil")
}
