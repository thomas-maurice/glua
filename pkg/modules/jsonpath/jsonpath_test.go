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

package jsonpath

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/, covering the
// Lua-visible behaviour of every function: syntax coverage, cardinality,
// render vs query/first/exists semantics, type fidelity and malformed
// input (see the package doc for the design decisions these pin down).
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("jsonpath", Loader)
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

// TestNormalisePath: pins the auto-brace decision -- a path that already
// starts with "{" is passed through untouched (so a full multi-segment
// render template is never double-wrapped); anything else gets wrapped.
func TestNormalisePath(t *testing.T) {
	assert.Equal(t, "{.spec.replicas}", normalisePath(".spec.replicas"))
	assert.Equal(t, "{$.spec.replicas}", normalisePath("$.spec.replicas"))
	assert.Equal(t, "{.spec.replicas}", normalisePath("{.spec.replicas}"))
	assert.Equal(t, "{.a}{.b}", normalisePath("{.a}{.b}"),
		"a template that already starts with { must not be re-wrapped")
}

// TestEvalPath_EmptyDataIsNotAMatch: querying against Go nil (e.g. a Lua nil
// document) must behave like a miss under AllowMissingKeys, not raise --
// this is a boundary the reflection walk could plausibly panic on.
func TestEvalPath_EmptyDataIsNotAMatch(t *testing.T) {
	results, err := evalPath("jsonpath.query", nil, ".spec.replicas", true)
	require.NoError(t, err)
	assert.Empty(t, results)
}

// TestQueryFn_EmptyResultIsEmptyTableNotNil: the nil-slice rule (see the
// working-here skill) applies to jsonpath.query's *lua.LTable return too --
// a zero-match query must give Lua an empty table it can safely call #/
// ipairs on, never a nil the Translator would otherwise have produced.
func TestQueryFn_EmptyResultIsEmptyTableNotNil(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := L.NewTable()
	result, err := queryFn(L, data, ".missing")
	require.NoError(t, err)
	require.NotNil(t, result, "queryFn must never return a nil *lua.LTable")
	assert.Equal(t, 0, result.MaxN())
}
