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

package strings

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestRep_CountCapped is the exact regression this chunk fixes (security
// review LOW 3): strings.rep("a", 1e8) previously allocated 100 MB
// silently, and scaling the count up further reaches process OOM, which is
// a Go FATAL error a caller's pcall cannot recover from -- unlike random's
// length arguments, which already cap at maxRepCount for the same reason.
func TestRep_CountCapped(t *testing.T) {
	// A normal use well under the cap must still work.
	out, err := rep("ab", 1000)
	require.NoError(t, err)
	assert.Len(t, out, 2000)

	// At the cap: still accepted.
	out, err = rep("a", maxRepCount)
	require.NoError(t, err)
	assert.Len(t, out, maxRepCount)

	// Just over the cap: rejected, naming the bound.
	_, err = rep("a", maxRepCount+1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "1048576")
}

// TestLuaScripts: runs all Lua test scripts in testdata/ directory
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	if err != nil {
		t.Fatalf("Failed to glob testdata: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No Lua test files found in testdata/")
	}

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("strings", Loader)

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
