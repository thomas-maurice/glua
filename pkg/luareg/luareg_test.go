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

package luareg_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// newState: returns a fresh *lua.LState with the test module pre-loaded under
// the given global name. Caller must close it.
func newState(t *testing.T, globalName string, mod *luareg.Module) *lua.LState {
	t.Helper()
	L := lua.NewState()
	n := mod.PushTo(L)
	require.Equal(t, 1, n)
	L.SetGlobal(globalName, L.Get(-1))
	L.Pop(1)
	return L
}

// eval: runs a Lua snippet and returns any error.
func eval(L *lua.LState, snippet string) error {
	return L.DoString(snippet)
}

// evalGet: runs a snippet then retrieves a global as a string for assertion.
func evalGet(t *testing.T, L *lua.LState, snippet, globalName string) lua.LValue {
	t.Helper()
	require.NoError(t, L.DoString(snippet))
	return L.GetGlobal(globalName)
}

// ---- 1. Per-primitive-type round-trip tests --------------------------------

func TestPrimitiveRoundTrips(t *testing.T) {
	type testCase struct {
		name    string
		regFn   func(m *luareg.Module)
		snippet string
		want    lua.LValue
	}

	cases := []testCase{
		{
			name:    "string",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s string) string { return s }, "") },
			snippet: `result = m.f("hello")`,
			want:    lua.LString("hello"),
		},
		{
			name:    "bool",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(b bool) bool { return b }, "") },
			snippet: `result = m.f(true)`,
			want:    lua.LTrue,
		},
		{
			name:    "int",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n int) int { return n }, "") },
			snippet: `result = m.f(42)`,
			want:    lua.LNumber(42),
		},
		{
			name:    "int8",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n int8) int8 { return n }, "") },
			snippet: `result = m.f(7)`,
			want:    lua.LNumber(7),
		},
		{
			name:    "int16",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n int16) int16 { return n }, "") },
			snippet: `result = m.f(1000)`,
			want:    lua.LNumber(1000),
		},
		{
			name:    "int32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n int32) int32 { return n }, "") },
			snippet: `result = m.f(100000)`,
			want:    lua.LNumber(100000),
		},
		{
			name:    "int64",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n int64) int64 { return n }, "") },
			snippet: `result = m.f(9000000)`,
			want:    lua.LNumber(9000000),
		},
		{
			name:    "uint",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n uint) uint { return n }, "") },
			snippet: `result = m.f(3)`,
			want:    lua.LNumber(3),
		},
		{
			name:    "uint8",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n uint8) uint8 { return n }, "") },
			snippet: `result = m.f(255)`,
			want:    lua.LNumber(255),
		},
		{
			name:    "uint16",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n uint16) uint16 { return n }, "") },
			snippet: `result = m.f(65000)`,
			want:    lua.LNumber(65000),
		},
		{
			name:    "uint32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n uint32) uint32 { return n }, "") },
			snippet: `result = m.f(123456)`,
			want:    lua.LNumber(123456),
		},
		{
			name:    "uint64",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n uint64) uint64 { return n }, "") },
			snippet: `result = m.f(999)`,
			want:    lua.LNumber(999),
		},
		{
			name:    "float32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n float32) float32 { return n }, "") },
			snippet: `result = m.f(1.5)`,
			want:    lua.LNumber(1.5),
		},
		{
			name:    "float64",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(n float64) float64 { return n }, "") },
			snippet: `result = m.f(3.14)`,
			want:    lua.LNumber(3.14),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mod := luareg.NewModule("m", "test")
			tc.regFn(mod)
			L := newState(t, "m", mod)
			defer L.Close()

			got := evalGet(t, L, tc.snippet, "result")
			assert.Equal(t, tc.want, got)
		})
	}
}

// ---- 2. Slice round-trip tests ---------------------------------------------

func TestSliceRoundTrip(t *testing.T) {
	t.Run("string slice non-empty", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(ss []string) []string { return ss }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.f({"a","b","c"})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok, "expected table")
		assert.Equal(t, lua.LString("a"), tbl.RawGetInt(1))
		assert.Equal(t, lua.LString("b"), tbl.RawGetInt(2))
		assert.Equal(t, lua.LString("c"), tbl.RawGetInt(3))
	})

	t.Run("string slice empty", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(ss []string) []string { return ss }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.f({})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, 0, tbl.Len())
	})

	t.Run("int slice", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(ns []int) []int { return ns }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.f({1,2,3})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, lua.LNumber(1), tbl.RawGetInt(1))
		assert.Equal(t, lua.LNumber(3), tbl.RawGetInt(3))
	})

	t.Run("float64 slice", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(fs []float64) []float64 { return fs }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.f({1.1, 2.2})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		first, _ := tbl.RawGetInt(1).(lua.LNumber)
		assert.InDelta(t, 1.1, float64(first), 1e-9)
	})
}

// ---- 3. Struct param + struct return via Translator ------------------------

type testPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestStructRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("identity", func(p testPerson) testPerson { return p }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.identity({name="Alice", age=30})`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("Alice"), tbl.RawGetString("name"))
	assert.Equal(t, lua.LNumber(30), tbl.RawGetString("age"))
}

func TestStructPtrRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("identity", func(p *testPerson) testPerson {
		if p == nil {
			return testPerson{}
		}
		return *p
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.identity({name="Bob", age=25})`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("Bob"), tbl.RawGetString("name"))
}

// ---- 4. Map params ---------------------------------------------------------

func TestMapParam(t *testing.T) {
	t.Run("map[string]string", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string]string) map[string]string { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({key="val"})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, lua.LString("val"), tbl.RawGetString("key"))
	})

	t.Run("map[string]int via float64 json unmarshalling", func(t *testing.T) {
		// The Translator round-trips via JSON so integers become float64.
		// We verify the value is numeric and correct.
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string]interface{}) map[string]interface{} { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({x=99})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, lua.LNumber(99), tbl.RawGetString("x"))
	})
}

// ---- 5. *lua.LState first-param escape hatch --------------------------------

func TestLStateEscapeHatch(t *testing.T) {
	// The function takes *lua.LState + string and builds a custom table.
	fn := func(L *lua.LState, key string) *lua.LTable {
		tbl := L.NewTable()
		tbl.RawSetString(key, lua.LString("injected"))
		return tbl
	}

	mod := luareg.NewModule("m", "test")
	mod.Fn("mk", fn, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.mk("greeting")`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("injected"), tbl.RawGetString("greeting"))
}

// ---- 6. Error raise --------------------------------------------------------

func TestErrorRaise(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("boom", func(n int) (int, error) {
		if n < 0 {
			return 0, errors.New("negative number")
		}
		return n * 2, nil
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	// Non-nil error path: pcall catches it.
	require.NoError(t, L.DoString(`ok, msg = pcall(m.boom, -1)`))
	assert.Equal(t, lua.LFalse, L.GetGlobal("ok"))
	msg := L.GetGlobal("msg").String()
	assert.True(t, strings.Contains(msg, "negative number"), "expected error message, got: %s", msg)

	// Nil error path: value is returned.
	require.NoError(t, L.DoString(`val = m.boom(5)`))
	assert.Equal(t, lua.LNumber(10), L.GetGlobal("val"))
}

// ---- 7. Multiple returns (no error) ----------------------------------------

func TestMultipleReturns(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("pair", func(s string) (string, int) {
		return s + "!", len(s)
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`a, b = m.pair("hi")`))
	assert.Equal(t, lua.LString("hi!"), L.GetGlobal("a"))
	assert.Equal(t, lua.LNumber(2), L.GetGlobal("b"))
}

// ---- 8. Arity mismatch -----------------------------------------------------

func TestArityMismatch(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(a, b string) string { return a + b }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	// Too few args.
	err := eval(L, `m.f("only_one")`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected 2 argument(s), got 1")

	// Too many args.
	err = eval(L, `m.f("a","b","c")`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected 2 argument(s), got 3")
}

// ---- 9. Type mismatch ------------------------------------------------------

func TestTypeMismatch(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(s string) string { return s }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	// Pass a number where a string is expected; L.CheckString raises a Lua error.
	err := eval(L, `m.f(42)`)
	require.Error(t, err)
}

// ---- 10. Duplicate registration panic --------------------------------------

func TestDuplicateRegistrationPanic(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("x", func() {}, "first")
	assert.PanicsWithValue(t,
		`luareg: module "m": duplicate registration of function "x"`,
		func() { mod.Fn("x", func() {}, "second") },
	)
}

// ---- 11. Variadic support --------------------------------------------------

// TestVariadicAllString: a pure-variadic ...string function collects its tail
// into a []string and returns it concatenated. Exercises empty/one/many cases.
func TestVariadicAllString(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("join", func(parts ...string) string {
		return strings.Join(parts, ",")
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`a = m.join(); b = m.join("x"); c = m.join("x","y","z")`))
	assert.Equal(t, lua.LString(""), L.GetGlobal("a"))
	assert.Equal(t, lua.LString("x"), L.GetGlobal("b"))
	assert.Equal(t, lua.LString("x,y,z"), L.GetGlobal("c"))
}

// TestVariadicWithFixedPrefix: a function with one fixed param plus a
// variadic tail. Verifies the fixed arg is bound correctly and the tail is
// collected starting at the right Lua position.
func TestVariadicWithFixedPrefix(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("prefix", func(sep string, parts ...string) string {
		return sep + ":" + strings.Join(parts, ",")
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`r0 = m.prefix("|"); r1 = m.prefix("|","a","b")`))
	assert.Equal(t, lua.LString("|:"), L.GetGlobal("r0"))
	assert.Equal(t, lua.LString("|:a,b"), L.GetGlobal("r1"))
}

// TestVariadicIntSum: numeric variadic, exercises the int element conversion path.
func TestVariadicIntSum(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("sum", func(ns ...int) int {
		total := 0
		for _, n := range ns {
			total += n
		}
		return total
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`r = m.sum(1,2,3,4)`))
	assert.Equal(t, lua.LNumber(10), L.GetGlobal("r"))
}

// TestVariadicMissingFixedArg: a function with a required fixed param plus a
// variadic tail must still enforce the fixed arg minimum.
func TestVariadicMissingFixedArg(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("need", func(first string, rest ...string) string { return first }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := eval(L, `m.need()`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected at least 1 argument(s)")
}

// TestVariadicElementTypeMismatch: a wrong-typed value in the variadic tail
// raises an ArgError at the variadic position, not at a fixed position.
func TestVariadicElementTypeMismatch(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("join", func(parts ...string) string { return strings.Join(parts, "") }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := eval(L, `m.join("a", 2, "c")`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected string")
}

// TestVariadicWithError: variadic plus trailing error return still surfaces
// errors via the standard raise path.
func TestVariadicWithError(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("sumPos", func(ns ...int) (int, error) {
		total := 0
		for _, n := range ns {
			if n < 0 {
				return 0, errors.New("negative not allowed")
			}
			total += n
		}
		return total, nil
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`ok, msg = pcall(m.sumPos, 1, -1, 2)`))
	assert.Equal(t, lua.LFalse, L.GetGlobal("ok"))
	assert.Contains(t, L.GetGlobal("msg").String(), "negative not allowed")

	require.NoError(t, L.DoString(`r = m.sumPos(1,2,3)`))
	assert.Equal(t, lua.LNumber(6), L.GetGlobal("r"))
}

// ---- 12. Non-func rejection ------------------------------------------------

func TestNonFuncRejectionPanic(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	assert.PanicsWithValue(t,
		`luareg: Fn "f": goFn must be a function, got int`,
		func() { mod.Fn("f", 42, "") },
	)
}

// ---- 13. Options populate metadata correctly --------------------------------

func TestOptionMetadata(t *testing.T) {
	reg := luareg.NewRegistry()
	mod := luareg.NewModule("m", "my module")
	mod.Fn("greet", func(name string) string { return "hi " + name }, "returns greeting",
		luareg.Args("name"),
		luareg.ArgDoc("name", "the person to greet"),
		luareg.ReturnDoc(0, "greeting", "the rendered greeting"),
	)
	mod.Register(reg)

	mods := reg.Modules()
	require.Len(t, mods, 1)
	assert.Equal(t, "m", mods[0].Name())
	assert.Equal(t, "my module", mods[0].Doc())

	fns := mods[0].Funcs()
	require.Len(t, fns, 1)

	// Verify arg name override.
	assert.Equal(t, "name", fns[0].ArgName(0))
	// Fallback for unset names.
	assert.Equal(t, "arg2", fns[0].ArgName(1))
}

// ---- 14. PushTo idempotency ------------------------------------------------

func TestPushToIdempotency(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(s string) string { return s + "!" }, "")

	L := lua.NewState()
	defer L.Close()

	// First push.
	n := mod.PushTo(L)
	assert.Equal(t, 1, n)
	L.SetGlobal("m1", L.Get(-1))
	L.Pop(1)

	// Second push.
	n = mod.PushTo(L)
	assert.Equal(t, 1, n)
	L.SetGlobal("m2", L.Get(-1))
	L.Pop(1)

	// Both tables should work identically.
	require.NoError(t, L.DoString(`r1 = m1.f("x"); r2 = m2.f("x")`))
	assert.Equal(t, lua.LString("x!"), L.GetGlobal("r1"))
	assert.Equal(t, lua.LString("x!"), L.GetGlobal("r2"))
}

// ---- Additional: nil return handling ----------------------------------------

func TestNilLValueReturn(t *testing.T) {
	// Return a nil *lua.LTable — should produce LNil on the Lua stack.
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func() *lua.LTable { return nil }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.f()`))
	assert.Equal(t, lua.LNil, L.GetGlobal("result"))
}

func TestNilStructPtrReturn(t *testing.T) {
	// Return a nil *testPerson — should produce LNil.
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func() *testPerson { return nil }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.f()`))
	assert.Equal(t, lua.LNil, L.GetGlobal("result"))
}

// ---- Additional: slice element type coverage --------------------------------

func TestSliceElementTypes(t *testing.T) {
	type sliceCase struct {
		name    string
		regFn   func(m *luareg.Module)
		snippet string
		check   func(t *testing.T, tbl *lua.LTable)
	}

	cases := []sliceCase{
		{
			name:    "[]bool",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []bool) []bool { return s }, "") },
			snippet: `result = m.f({true, false, true})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LTrue, tbl.RawGetInt(1))
				assert.Equal(t, lua.LFalse, tbl.RawGetInt(2))
			},
		},
		{
			name:    "[]int8",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []int8) []int8 { return s }, "") },
			snippet: `result = m.f({1, 2})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(1), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]int16",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []int16) []int16 { return s }, "") },
			snippet: `result = m.f({100})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(100), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]int32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []int32) []int32 { return s }, "") },
			snippet: `result = m.f({200})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(200), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]int64",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []int64) []int64 { return s }, "") },
			snippet: `result = m.f({300})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(300), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]uint",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []uint) []uint { return s }, "") },
			snippet: `result = m.f({4})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(4), tbl.RawGetInt(1))
			},
		},
		// Note: []uint8 == []byte; encoding/json marshals it as a base64 string,
		// so the Translator returns LString, not LTable. This is a known limitation
		// of the JSON round-trip path in pkg/glua.Translator and is out of A1 scope.
		{
			name:    "[]uint16",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []uint16) []uint16 { return s }, "") },
			snippet: `result = m.f({16})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(16), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]uint32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []uint32) []uint32 { return s }, "") },
			snippet: `result = m.f({32})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(32), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]uint64",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []uint64) []uint64 { return s }, "") },
			snippet: `result = m.f({64})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				assert.Equal(t, lua.LNumber(64), tbl.RawGetInt(1))
			},
		},
		{
			name:    "[]float32",
			regFn:   func(m *luareg.Module) { m.Fn("f", func(s []float32) []float32 { return s }, "") },
			snippet: `result = m.f({1.5})`,
			check: func(t *testing.T, tbl *lua.LTable) {
				n, _ := tbl.RawGetInt(1).(lua.LNumber)
				assert.InDelta(t, 1.5, float64(n), 1e-6)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mod := luareg.NewModule("m", "test")
			tc.regFn(mod)
			L := newState(t, "m", mod)
			defer L.Close()
			require.NoError(t, L.DoString(tc.snippet))
			tbl, ok := L.GetGlobal("result").(*lua.LTable)
			require.True(t, ok)
			tc.check(t, tbl)
		})
	}
}

// TestSliceElementTypeMismatchVariants: non-matching values in typed slices raise errors.
func TestSliceElementTypeMismatchVariants(t *testing.T) {
	t.Run("bool slice gets number", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(s []bool) []bool { return s }, "")
		L := newState(t, "m", mod)
		defer L.Close()
		err := eval(L, `m.f({1, 2})`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected boolean")
	})

	t.Run("int slice gets string", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(s []int) []int { return s }, "")
		L := newState(t, "m", mod)
		defer L.Close()
		err := eval(L, `m.f({"hello"})`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected number")
	})

	t.Run("float64 slice gets string", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("f", func(s []float64) []float64 { return s }, "")
		L := newState(t, "m", mod)
		defer L.Close()
		err := eval(L, `m.f({"oops"})`)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected number")
	})
}

// ---- Additional: typed-map coverage (string→primitive values) --------------

// TestTypedMapRoundTrips: round-trips Go map types with primitive value
// kinds through the Translator. Locks in that the map branch in luaToGo /
// goToLua works for the common typed-value cases (not just map[string]any).
func TestTypedMapRoundTrips(t *testing.T) {
	t.Run("map[string]int", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string]int) map[string]int { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({a=1, b=2})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, lua.LNumber(1), tbl.RawGetString("a"))
		assert.Equal(t, lua.LNumber(2), tbl.RawGetString("b"))
	})

	t.Run("map[string]bool", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string]bool) map[string]bool { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({on=true, off=false})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		assert.Equal(t, lua.LTrue, tbl.RawGetString("on"))
		assert.Equal(t, lua.LFalse, tbl.RawGetString("off"))
	})

	t.Run("map[string]float64", func(t *testing.T) {
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string]float64) map[string]float64 { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({pi=3.14, e=2.71})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		pi, _ := tbl.RawGetString("pi").(lua.LNumber)
		assert.InDelta(t, 3.14, float64(pi), 1e-9)
	})

	t.Run("map[string][]string", func(t *testing.T) {
		// Map whose value is itself a slice — exercises nested Translator paths.
		mod := luareg.NewModule("m", "test")
		mod.Fn("echo", func(m map[string][]string) map[string][]string { return m }, "")
		L := newState(t, "m", mod)
		defer L.Close()

		require.NoError(t, L.DoString(`result = m.echo({tags={"a","b","c"}})`))
		tbl, ok := L.GetGlobal("result").(*lua.LTable)
		require.True(t, ok)
		inner, ok := tbl.RawGetString("tags").(*lua.LTable)
		require.True(t, ok, "expected nested table")
		assert.Equal(t, lua.LString("a"), inner.RawGetInt(1))
		assert.Equal(t, lua.LString("c"), inner.RawGetInt(3))
	})
}

// ---- Additional: composite slice / nested types -----------------------------

// nestedItem: small struct used to test []struct round-trips.
type nestedItem struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// TestSliceOfStructRoundTrip: a Go function taking []<struct> and returning
// the same slice. Composite slice element types (struct, *struct, slice, map)
// are dispatched to the Translator inside luaToGoValue, so K8s-style list
// shapes work on the input path.
func TestSliceOfStructRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("echo", func(items []nestedItem) []nestedItem { return items }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`
		result = m.echo({
			{name="alice", score=10},
			{name="bob",   score=20},
		})
	`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)

	first, ok := tbl.RawGetInt(1).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("alice"), first.RawGetString("name"))
	assert.Equal(t, lua.LNumber(10), first.RawGetString("score"))

	second, ok := tbl.RawGetInt(2).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("bob"), second.RawGetString("name"))
	assert.Equal(t, lua.LNumber(20), second.RawGetString("score"))
}

// TestSliceOfPtrStructRoundTrip: []*<struct> on the input path. K8s typed
// clients often hand back pointer slices (e.g. []*corev1.Pod), so authors
// writing reverse-direction helpers want to accept that shape.
func TestSliceOfPtrStructRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("count", func(items []*nestedItem) int {
		// Touch each pointer to confirm they're real instances, not nil.
		sum := 0
		for _, it := range items {
			sum += it.Score
		}
		return sum
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`
		result = m.count({
			{name="alice", score=10},
			{name="bob",   score=20},
			{name="carol", score=5},
		})
	`))
	assert.Equal(t, lua.LNumber(35), L.GetGlobal("result"))
}

// TestSliceOfStructInsideStructField: same shape inside a struct field.
// This used to be the only working path; it still works after the input-side
// fix and serves as a sanity check that the Translator path was not broken.
func TestSliceOfStructInsideStructField(t *testing.T) {
	type bag struct {
		Items []nestedItem `json:"items"`
	}
	mod := luareg.NewModule("m", "test")
	mod.Fn("echo", func(b bag) bag { return b }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`
		result = m.echo({items={
			{name="alice", score=10},
			{name="bob",   score=20},
		}})
	`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	items, ok := tbl.RawGetString("items").(*lua.LTable)
	require.True(t, ok)
	first, ok := items.RawGetInt(1).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("alice"), first.RawGetString("name"))
	assert.Equal(t, lua.LNumber(10), first.RawGetString("score"))
}

// TestNestedSliceOfSliceRoundTrip: [][]string on the input path. The outer
// slice is walked directly; each inner []string element is now dispatched to
// the Translator.
func TestNestedSliceOfSliceRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("echo", func(rows [][]string) [][]string { return rows }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.echo({ {"a","b"}, {"c","d","e"} })`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)

	row1, ok := tbl.RawGetInt(1).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("a"), row1.RawGetInt(1))
	assert.Equal(t, lua.LString("b"), row1.RawGetInt(2))

	row2, ok := tbl.RawGetInt(2).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("e"), row2.RawGetInt(3))
}

// TestSliceOfMapRoundTrip: []map[string]int on the input path. Element type
// is a map; goes through the Translator branch in luaToGoValue.
func TestSliceOfMapRoundTrip(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("echo", func(rows []map[string]int) []map[string]int { return rows }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.echo({ {a=1, b=2}, {c=3} })`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	first, ok := tbl.RawGetInt(1).(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LNumber(1), first.RawGetString("a"))
	assert.Equal(t, lua.LNumber(2), first.RawGetString("b"))
}

// TestSliceOfPtrStructNilEntry: a nil entry in a []*struct materializes as
// a typed nil pointer rather than raising. This matches typical Lua-side
// idioms where users may sparsely populate a list.
func TestSliceOfPtrStructNilEntry(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("nils", func(items []*nestedItem) int {
		count := 0
		for _, it := range items {
			if it == nil {
				count++
			}
		}
		return count
	}, "")
	L := newState(t, "m", mod)
	defer L.Close()

	// Lua arrays don't store nils contiguously by default; explicitly set
	// index 2 to nil after building, so MaxN still reports 3.
	require.NoError(t, L.DoString(`
		local t = { {name="a", score=1}, false, {name="c", score=3} }
		t[2] = nil
		t[3] = {name="c", score=3}  -- ensure MaxN reaches 3
		result = m.nils(t)
	`))
	assert.Equal(t, lua.LNumber(1), L.GetGlobal("result"))
}

// TestStructWithNestedSliceField: a struct field containing []string round-trips.
// This is the common case in K8s-like modules (e.g. a Pod with .Containers[]).
func TestStructWithNestedSliceField(t *testing.T) {
	type withSlice struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	mod := luareg.NewModule("m", "test")
	mod.Fn("echo", func(s withSlice) withSlice { return s }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.echo({name="x", tags={"a","b"}})`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("x"), tbl.RawGetString("name"))

	tags, ok := tbl.RawGetString("tags").(*lua.LTable)
	require.True(t, ok)
	assert.Equal(t, lua.LString("a"), tags.RawGetInt(1))
	assert.Equal(t, lua.LString("b"), tags.RawGetInt(2))
}

// ---- Additional: []byte / []uint8 quirks -----------------------------------

// TestByteSliceReturnIsBase64String: encoding/json marshals []byte to a base64
// string, so a Go function returning []byte yields an LString on the Lua side,
// not a table of numbers. This locks in the known limitation of the JSON
// round-trip path used by pkg/glua.Translator on the return side.
func TestByteSliceReturnIsBase64String(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("payload", func() []byte { return []byte("hi") }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.payload()`))
	got, ok := L.GetGlobal("result").(lua.LString)
	require.True(t, ok, "expected LString (base64-encoded), got %T", L.GetGlobal("result"))
	// base64("hi") == "aGk="
	assert.Equal(t, lua.LString("aGk="), got)
}

// ---- Additional: goToLua coverage for primitives returned ------------------

func TestReturnBool(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func() bool { return false }, "")
	L := newState(t, "m", mod)
	defer L.Close()
	require.NoError(t, L.DoString(`result = m.f()`))
	assert.Equal(t, lua.LFalse, L.GetGlobal("result"))
}

func TestReturnUint(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func() uint { return 7 }, "")
	L := newState(t, "m", mod)
	defer L.Close()
	require.NoError(t, L.DoString(`result = m.f()`))
	assert.Equal(t, lua.LNumber(7), L.GetGlobal("result"))
}

// ---- Additional: unsupported param type raises error at call time ----------

func TestUnsupportedParamType(t *testing.T) {
	// chan int is not in the supported set; should raise a Lua error on call.
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(ch chan int) string { return "ok" }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := eval(L, `m.f(42)`) // 42 not a table, but error is "unsupported type"
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

func TestUnsupportedSliceElementType(t *testing.T) {
	// []chan int — slice of unsupported element type; hits luaToGoValue default branch.
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(s []chan int) string { return "ok" }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := eval(L, `m.f({1})`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported")
}

// ---- Additional: no-arg / no-return functions work -------------------------

func TestNoArgs(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("ping", func() string { return "pong" }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.ping()`))
	assert.Equal(t, lua.LString("pong"), L.GetGlobal("result"))
}

func TestNoReturn(t *testing.T) {
	called := false
	mod := luareg.NewModule("m", "test")
	mod.Fn("noop", func(s string) { called = true }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`m.noop("x")`))
	assert.True(t, called)
}

// ---- Additional: slice type-mismatch (string element, number given) --------

func TestSliceElementTypeMismatch(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(ss []string) []string { return ss }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := eval(L, `m.f({1, 2, 3})`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected string")
}

// ---- Additional: nil goFn panics -------------------------------------------

func TestNilGoFnPanic(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	assert.PanicsWithValue(t,
		`luareg: Fn "f": goFn must not be nil`,
		func() { mod.Fn("f", nil, "") },
	)
}

// ---- Regression: arity-mismatch uses RaiseError (no "bad argument" prefix) -

// TestArityErrorMessageNoBadArgumentPrefix: the arity mismatch should produce
// a plain "expected N argument(s), got M" message, not "bad argument #1" — the
// previous implementation used ArgError which is misleading for arity issues.
func TestArityErrorMessageNoBadArgumentPrefix(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("f", func(a string) string { return a }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := L.DoString(`m.f("a", "extra")`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected 1 argument(s), got 2")
	assert.NotContains(t, err.Error(), "bad argument")
}

// ---- Regression: bare struct AND *struct params/returns work --------------

type roundTripStruct struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

// TestStructByValueParamAndReturn: bare-struct param (no pointer) round-trips
// through reflection + Translator. The struct branch in luaToGo previously had
// dead `isPtr` logic; this test locks in that the simplified single-branch
// path still works for value receivers.
func TestStructByValueParamAndReturn(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("identity", func(s roundTripStruct) roundTripStruct { return s }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`out = m.identity({name = "alice", n = 42})`))
	out := L.GetGlobal("out").(*lua.LTable)
	assert.Equal(t, lua.LString("alice"), out.RawGetString("name"))
	assert.Equal(t, lua.LNumber(42), out.RawGetString("n"))
}

// TestStructByPointerParamAndReturn: *struct param round-trips. The Ptr branch
// in luaToGo was rewritten to handle only pointer-to-struct (anything else is
// rejected); this verifies the happy path.
func TestStructByPointerParamAndReturn(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("identity_ptr", func(s *roundTripStruct) *roundTripStruct { return s }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`out = m.identity_ptr({name = "bob", n = 7})`))
	out := L.GetGlobal("out").(*lua.LTable)
	assert.Equal(t, lua.LString("bob"), out.RawGetString("name"))
	assert.Equal(t, lua.LNumber(7), out.RawGetString("n"))
}

// TestPointerToNonStructRejected: pointer-to-non-struct surfaces as an arg
// conversion error at call time (the wrapper is built fine, but the call fails
// the unsupported-type guard).
func TestPointerToNonStructRejected(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("bad", func(p *int) int { return *p }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	err := L.DoString(`m.bad(7)`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported pointer type")
}

// ---- Regression: returning a bare lua.LValue concrete type doesn't panic --

// TestReturnBareLuaStringConcreteType: returning lua.LString directly (a value
// type that implements lua.LValue) previously triggered IsNil on a non-nilable
// kind and panicked. The fix gates IsNil to nilable kinds only.
func TestReturnBareLuaStringConcreteType(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("hello", func() lua.LString { return lua.LString("hi") }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`out = m.hello()`))
	assert.Equal(t, lua.LString("hi"), L.GetGlobal("out"))
}

// TestReturnBareLuaNumberConcreteType: same regression for LNumber.
func TestReturnBareLuaNumberConcreteType(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("forty_two", func() lua.LNumber { return lua.LNumber(42) }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`out = m.forty_two()`))
	assert.Equal(t, lua.LNumber(42), L.GetGlobal("out"))
}

// TestReturnNilLuaTablePointer: a nil *lua.LTable still returns LNil cleanly
// (the nilable kinds branch is exercised).
func TestReturnNilLuaTablePointer(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("nil_table", func() *lua.LTable { return nil }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`out = m.nil_table()`))
	assert.Equal(t, lua.LNil, L.GetGlobal("out"))
}

// ---- Const API tests ---------------------------------------------------------

// constPoint: a small struct used to test Const with struct values.
type constPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// TestConstString: a string constant is accessible on the module table.
func TestConstString(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("VERSION", "1.2.3", "string", "the version string")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.VERSION`))
	assert.Equal(t, lua.LString("1.2.3"), L.GetGlobal("result"))
}

// TestConstNumber: a numeric constant is accessible on the module table.
func TestConstNumber(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("PI", 3.14, "number", "pi approximation")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.PI`))
	got := L.GetGlobal("result").(lua.LNumber)
	assert.InDelta(t, 3.14, float64(got), 1e-9)
}

// TestConstStruct: a struct constant is serialized via the Translator and
// accessible as a table with the expected fields.
func TestConstStruct(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("ORIGIN", constPoint{X: 0, Y: 0}, "stubgen.constPoint", "origin point")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`result = m.ORIGIN`))
	tbl, ok := L.GetGlobal("result").(*lua.LTable)
	require.True(t, ok, "expected table for struct const")
	assert.Equal(t, lua.LNumber(0), tbl.RawGetString("x"))
}

// TestConstAndFnCoexist: a module with both Const and Fn has both available.
func TestConstAndFnCoexist(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("LABEL", "hello", "string", "a label constant")
	mod.Fn("identity", func(s string) string { return s }, "")
	L := newState(t, "m", mod)
	defer L.Close()

	require.NoError(t, L.DoString(`c = m.LABEL; r = m.identity("world")`))
	assert.Equal(t, lua.LString("hello"), L.GetGlobal("c"))
	assert.Equal(t, lua.LString("world"), L.GetGlobal("r"))
}

// TestConstDuplicatePanic: registering the same const name twice panics.
func TestConstDuplicatePanic(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("X", "first", "string", "first")
	assert.PanicsWithValue(t,
		`luareg: module "m": duplicate registration of const "X"`,
		func() { mod.Const("X", "second", "string", "second") },
	)
}

// TestConstMetadata: Consts() returns metadata in registration order.
func TestConstMetadata(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Const("A", "aval", "string", "first const")
	mod.Const("B", 42, "number", "second const")

	consts := mod.Consts()
	require.Len(t, consts, 2)
	assert.Equal(t, "A", consts[0].Name)
	assert.Equal(t, "aval", consts[0].Value)
	assert.Equal(t, "string", consts[0].LuaType)
	assert.Equal(t, "first const", consts[0].Doc)
	assert.Equal(t, "B", consts[1].Name)
}
