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

package glua

import (
	"fmt"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestFromLua_SelfReferentialTable: a table that references itself directly
// (`t.self = t`) must produce an ordinary, catchable Go error instead of
// recursing forever. Before the depth/ancestor guard, this pattern crashed
// the entire process with an unrecoverable "fatal error: stack overflow" —
// verified separately with a throwaway program against the pre-fix code
// (see the chunk handoff notes), since a real stack overflow cannot be
// caught inside a Go test without killing the test binary too.
func TestFromLua_SelfReferentialTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		t = {}
		t.self = t
	`); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()
	lv := L.GetGlobal("t")

	var out map[string]interface{}
	err := tr.FromLua(L, lv, &out)
	if err == nil {
		t.Fatal("expected an error converting a self-referential table, got nil")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected error to mention a cycle, got: %v", err)
	}
}

// TestFromLua_SelfReferentialTable_ViaLuaPcall: the same self-reference as
// above, but proven at the level a script author actually experiences it —
// wrapped in a Lua-level pcall through a real module-style Go function
// registered via *lua.LState so the error surfaces the same way a real
// module function's error would. Asserts BOTH that pcall catches it (it is
// an ordinary error, not fatal) AND that the Lua state keeps running
// afterwards, which is the whole point of the fix: the process survives.
func TestFromLua_SelfReferentialTable_ViaLuaPcall(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tr := NewTranslator()
	L.SetGlobal("convert", L.NewFunction(func(L *lua.LState) int {
		var out map[string]interface{}
		if err := tr.FromLua(L, L.CheckTable(1), &out); err != nil {
			L.RaiseError("%s", err.Error())
			return 0
		}
		L.Push(lua.LTrue)
		return 1
	}))

	err := L.DoString(`
		local t = {}
		t.self = t

		local ok, e = pcall(convert, t)
		assert(ok == false, "expected pcall to report failure")
		assert(string.find(tostring(e), "cycle") ~= nil, "expected error to mention 'cycle', got: " .. tostring(e))

		-- Prove the state is still alive and usable after the caught error.
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
		t.Errorf("expected script to keep executing correctly, got result=%s", got.String())
	}
}

// TestFromLua_MutuallyRecursiveTables: two tables that reference each other
// (a.b = b, b.a = a) form a cycle just as much as direct self-reference, and
// must be rejected the same way.
func TestFromLua_MutuallyRecursiveTables(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		a = {}
		b = {}
		a.b = b
		b.a = a
	`); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()

	var out map[string]interface{}
	err := tr.FromLua(L, L.GetGlobal("a"), &out)
	if err == nil {
		t.Fatal("expected an error converting mutually-recursive tables, got nil")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected error to mention a cycle, got: %v", err)
	}
}

// TestFromLua_CycleNestedSeveralLevelsDeep: the cycle guard must catch a
// self-reference that only appears several levels below the root, not just
// at the top — otherwise a check that only inspects the root table would
// pass this case straight through into the same fatal recursion.
func TestFromLua_CycleNestedSeveralLevelsDeep(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		root = { level1 = { level2 = { level3 = {} } } }
		root.level1.level2.level3.backref = root.level1
	`); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()

	var out map[string]interface{}
	err := tr.FromLua(L, L.GetGlobal("root"), &out)
	if err == nil {
		t.Fatal("expected an error converting a table with a cycle several levels down, got nil")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected error to mention a cycle, got: %v", err)
	}
}

// TestFromLua_DeepButFiniteTable_AtLimit_Succeeds: a legitimately deep,
// acyclic table nested right up to maxTableDepth must still convert
// successfully — the guard exists to stop unbounded recursion, not to
// arbitrarily reject realistic (if unusually deep) payloads.
func TestFromLua_DeepButFiniteTable_AtLimit_Succeeds(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Build a chain of maxTableDepth nested tables: {inner = {inner = ...}}.
	// The outermost table is depth 1, so nesting maxTableDepth-1 more levels
	// inside it lands exactly at maxTableDepth total, which must still pass.
	var b strings.Builder
	b.WriteString("root = {}\n")
	b.WriteString("local cur = root\n")
	fmt.Fprintf(&b, "for i = 1, %d do\n", maxTableDepth-1)
	b.WriteString("  cur.inner = {}\n")
	b.WriteString("  cur = cur.inner\n")
	b.WriteString("end\n")
	b.WriteString("cur.value = \"leaf\"\n")

	if err := L.DoString(b.String()); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()
	var out map[string]interface{}
	if err := tr.FromLua(L, L.GetGlobal("root"), &out); err != nil {
		t.Fatalf("expected a table at the depth limit to convert successfully, got error: %v", err)
	}
}

// TestFromLua_TooDeepTable_ErrorsCleanly: a table nested one level beyond
// maxTableDepth must be rejected with a clear, actionable error rather than
// overflowing the stack.
func TestFromLua_TooDeepTable_ErrorsCleanly(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	var b strings.Builder
	b.WriteString("root = {}\n")
	b.WriteString("local cur = root\n")
	fmt.Fprintf(&b, "for i = 1, %d do\n", maxTableDepth+10)
	b.WriteString("  cur.inner = {}\n")
	b.WriteString("  cur = cur.inner\n")
	b.WriteString("end\n")

	if err := L.DoString(b.String()); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()
	var out map[string]interface{}
	err := tr.FromLua(L, L.GetGlobal("root"), &out)
	if err == nil {
		t.Fatal("expected an error converting a table nested beyond the depth limit, got nil")
	}
	if !strings.Contains(err.Error(), "maximum nesting depth") {
		t.Errorf("expected error to mention the depth limit, got: %v", err)
	}
}

// TestFromLua_DiamondReference_NotACycle_Succeeds: the same table reachable
// from two different fields of its parent (a "diamond") is NOT a cycle — it
// never appears as its own ancestor on either path — and must still convert.
// This pins the deliberate design choice to use PATH-scoped ancestor
// tracking rather than a global "already seen" set: a global set would
// reject this shape on the second occurrence, which would be a regression
// relative to encoding/json (which marshals a shared reference fine, just
// as two independent copies, since JSON itself has no notion of aliasing).
func TestFromLua_DiamondReference_NotACycle_Succeeds(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		shared = { value = "shared-value" }
		root = { a = shared, b = shared }
	`); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()
	var out map[string]interface{}
	if err := tr.FromLua(L, L.GetGlobal("root"), &out); err != nil {
		t.Fatalf("expected a diamond (non-cyclic shared reference) to convert successfully, got error: %v", err)
	}

	a, ok := out["a"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected out[\"a\"] to be a map, got %T", out["a"])
	}
	bVal, ok := out["b"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected out[\"b\"] to be a map, got %T", out["b"])
	}
	if a["value"] != "shared-value" || bVal["value"] != "shared-value" {
		t.Errorf("expected both branches of the diamond to carry the shared value, got a=%v b=%v", a, bVal)
	}
}

// TestFromLua_CycleInsideArray: a cycle reached through the array portion of
// a table (rather than the map/ForEach portion) must also be caught — the
// array and map branches of fromLuaValueGuarded share the same guard, but
// this pins that the array branch actually exercises it.
func TestFromLua_CycleInsideArray(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoString(`
		t = {}
		t[1] = t
	`); err != nil {
		t.Fatalf("failed to set up Lua state: %v", err)
	}

	tr := NewTranslator()
	var out []interface{}
	err := tr.FromLua(L, L.GetGlobal("t"), &out)
	if err == nil {
		t.Fatal("expected an error converting an array containing itself, got nil")
	}
	if !strings.Contains(err.Error(), "cycle") {
		t.Errorf("expected error to mention a cycle, got: %v", err)
	}
}
