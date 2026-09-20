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

package kubernetes

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestEnsureMetadata_SelfReferentialTable_ViaLuaPcall: proves the
// pkg/glua.Translator cycle guard at a real module boundary, not just
// against the Translator directly. ensure_metadata takes a `map[string]any`
// Go parameter, which pkg/luareg's argument conversion (luaToGo, reflect.Map
// case) routes through Translator.FromLua — the same path k8sclient writes
// and every other module function taking a table/map/struct argument uses.
//
// Before the pkg/glua fix, handing a self-referential table to ANY such
// function recursed forever and killed the whole process with an
// unrecoverable "fatal error: stack overflow" (gopher-lua's pcall cannot
// catch that, nor can Go's recover()). This asserts the call now surfaces as
// an ordinary catchable Lua error, AND that the Lua state keeps running
// afterwards.
//
// Note: json.stringify / yaml.stringify / template.render do NOT exercise
// this fix — they have their own independent, unguarded recursive
// converters (pkg/modules/json, pkg/modules/yaml, pkg/modules/template) that
// never call pkg/glua.Translator at all, because their Lua parameter type is
// `lua.LValue`/`*lua.LTable` taken as a pass-through escape hatch rather than
// a typed Go struct/map. That gap is out of scope for this chunk (which
// targets pkg/glua/translator.go) and is called out separately in the
// handoff notes; ensure_metadata is used here specifically because it DOES
// go through the Translator.
func TestEnsureMetadata_SelfReferentialTable_ViaLuaPcall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	err := L.DoString(`
		local k8s = require("kubernetes")

		local obj = {}
		obj.self = obj

		local ok, e = pcall(k8s.ensure_metadata, obj)
		assert(ok == false, "expected pcall to report failure for a self-referential table")
		assert(string.find(tostring(e), "cycle") ~= nil, "expected error to mention 'cycle', got: " .. tostring(e))

		-- Prove the process/state survived the caught error.
		survived = true
		result = 1 + 1
	`)
	if err != nil {
		t.Fatalf("unexpected top-level Lua error: %v", err)
	}

	if L.GetGlobal("survived") != lua.LTrue {
		t.Fatal("expected script execution to continue past the caught error")
	}
	if got := L.GetGlobal("result"); got.String() != "2" {
		t.Errorf("expected script to keep executing correctly after the caught error, got result=%s", got.String())
	}
}
