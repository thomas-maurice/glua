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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// ---------------------------------------------------------------------------
// Test fixtures
// ---------------------------------------------------------------------------

// Counter: simple stateful type used across class tests.
type Counter struct {
	value int
}

// NewCounter: factory function — returns a *Counter.
func NewCounter() *Counter { return &Counter{} }

// Increment: adds delta to the counter. Returns the receiver for chaining tests.
func (c *Counter) Increment(delta int) { c.value += delta }

// Get: returns the current value.
func (c *Counter) Get() int { return c.value }

// Reset: sets the counter to zero.
func (c *Counter) Reset() { c.value = 0 }

// GetWithState: escape-hatch method that uses *lua.LState as second param.
// Returns a Lua table {value=N} built with L to verify the hatch works.
func (c *Counter) GetWithState(L *lua.LState) *lua.LTable {
	tbl := L.NewTable()
	L.SetField(tbl, "value", lua.LNumber(c.value))
	return tbl
}

// IncrementBoom: returns a *luareg.Error when delta < 0.
func (c *Counter) IncrementBoom(delta int) error {
	if delta < 0 {
		return &luareg.Error{Kind: "NegativeDelta", Message: "delta must be non-negative"}
	}
	c.value += delta
	return nil
}

// IncrementPlainError: returns a plain errors.New error for backward-compat test.
func (c *Counter) IncrementPlainError(delta int) error {
	if delta < 0 {
		return fmt.Errorf("plain: delta must be non-negative")
	}
	c.value += delta
	return nil
}

// Chainer: a second class used to verify multi-class modules.
type Chainer struct {
	label string
}

// NewChainer: factory.
func NewChainer(label string) *Chainer { return &Chainer{label: label} }

// Label: returns the label.
func (ch *Chainer) Label() string { return ch.label }

// WithLabel: returns a new Chainer with a different label (tests *Chainer return auto-wrap).
func (ch *Chainer) WithLabel(label string) *Chainer { return &Chainer{label: label} }

// ---------------------------------------------------------------------------
// Value-receiver fixture
// ---------------------------------------------------------------------------

// ValCounter: a value-type counter (not a pointer). Used to test NewClass[ValCounter].
type ValCounter struct {
	V int
}

func (vc ValCounter) GetVal() int { return vc.V }

// ---------------------------------------------------------------------------
// Helper: build a module with Counter and push it.
// ---------------------------------------------------------------------------

func counterModule() *luareg.Module {
	m := luareg.NewModule("counter", "counter module")

	c := luareg.NewClass[*Counter]("counter.Counter", "a simple counter")
	c.Method("increment", (*Counter).Increment, "add delta to counter", luareg.Args("delta"))
	c.Method("get", (*Counter).Get, "return current value")
	c.Method("reset", (*Counter).Reset, "reset to zero")
	c.Method("get_with_state", (*Counter).GetWithState, "return value in a table (lstate escape hatch)")
	c.Method("increment_boom", (*Counter).IncrementBoom, "increment or raise structured error", luareg.Args("delta"))
	c.Method("increment_plain_error", (*Counter).IncrementPlainError, "increment or raise plain error", luareg.Args("delta"))
	m.RegisterClass(c)

	m.Fn("new", NewCounter, "create a new counter")
	return m
}

func newCounterState(t *testing.T) *lua.LState {
	t.Helper()
	L := lua.NewState()
	n := counterModule().PushTo(L)
	require.Equal(t, 1, n)
	L.SetGlobal("counter", L.Get(-1))
	L.Pop(1)
	return L
}

// ---------------------------------------------------------------------------
// 1. Basic class lifecycle
// ---------------------------------------------------------------------------

func TestClassBasicLifecycle(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	// increment(1) + increment(5) = 6; get; reset; get again
	require.NoError(t, L.DoString(`
		local c = counter.new()
		c:increment(1)
		c:increment(5)
		assert(c:get() == 6, "expected 6")
		c:reset()
		assert(c:get() == 0, "expected 0 after reset")
	`))
}

// ---------------------------------------------------------------------------
// 2. Method chaining via class return (auto-wrap *Chainer)
// ---------------------------------------------------------------------------

func TestClassMethodChaining(t *testing.T) {
	m := luareg.NewModule("ch", "chainer module")
	c := luareg.NewClass[*Chainer]("ch.Chainer", "a chainer")
	c.Method("with_label", (*Chainer).WithLabel, "return chainer with new label", luareg.Args("label"))
	c.Method("label", (*Chainer).Label, "return label")
	m.RegisterClass(c)
	m.Fn("new", NewChainer, "create chainer", luareg.Args("label"))

	L := lua.NewState()
	defer L.Close()
	n := m.PushTo(L)
	require.Equal(t, 1, n)
	L.SetGlobal("ch", L.Get(-1))
	L.Pop(1)

	// with_label returns *Chainer — must be auto-wrapped so :label() works.
	require.NoError(t, L.DoString(`
		local c = ch.new("first")
		local c2 = c:with_label("second"):with_label("third")
		assert(c2:label() == "third", "expected third, got " .. c2:label())
		-- original unchanged
		assert(c:label() == "first", "expected first")
	`))
}

// ---------------------------------------------------------------------------
// 3. Multiple instances are independent
// ---------------------------------------------------------------------------

func TestClassMultipleInstancesIndependent(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	require.NoError(t, L.DoString(`
		local a = counter.new()
		local b = counter.new()
		a:increment(10)
		b:increment(3)
		assert(a:get() == 10, "a expected 10")
		assert(b:get() == 3, "b expected 3")
	`))
}

// ---------------------------------------------------------------------------
// 4. Wrong-type receiver raises error mentioning class name
// ---------------------------------------------------------------------------

func TestClassWrongTypeReceiver(t *testing.T) {
	m := luareg.NewModule("mod", "test")

	c1 := luareg.NewClass[*Counter]("mod.Counter", "counter")
	c1.Method("get", (*Counter).Get, "get")
	m.RegisterClass(c1)

	c2 := luareg.NewClass[*Chainer]("mod.Chainer", "chainer")
	c2.Method("label", (*Chainer).Label, "label")
	m.RegisterClass(c2)

	m.Fn("new_counter", NewCounter, "new counter")
	m.Fn("new_chainer", func() *Chainer { return &Chainer{label: "x"} }, "new chainer")

	L := lua.NewState()
	defer L.Close()
	m.PushTo(L)
	L.SetGlobal("mod", L.Get(-1))
	L.Pop(1)

	// Reach into Counter's metatable to call its :get method with a Chainer
	// userdata as the receiver. The method wrapper's Check(L, 1) must surface
	// a Lua error naming the expected class.
	err := L.DoString(`
		local counter_ud = mod.new_counter()
		local chainer_ud = mod.new_chainer()
		-- invoke counter's "get" method directly with chainer_ud as self
		local mt = getmetatable(counter_ud)
		local getfn = mt.__index.get
		local ok, msg = pcall(getfn, chainer_ud)
		assert(ok == false, "expected error")
		assert(string.find(msg, "mod.Counter") ~= nil, "expected class name in error, got: " .. msg)
	`)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// 5. Bare table passed where userdata expected
// ---------------------------------------------------------------------------

func TestClassBareTableRaisesError(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	err := L.DoString(`
		local ok, msg = pcall(function()
			-- Get the increment method from a real counter's metatable,
			-- then call it with a plain table instead of userdata.
			local real = counter.new()
			local mt = getmetatable(real)
			local inc = mt.__index.increment
			inc({}, 1)  -- {} is a plain table, not userdata
		end)
		assert(ok == false, "expected error")
		-- gopher-lua CheckUserData raises "userdata expected" for non-userdata
	`)
	require.NoError(t, err)
}

// ---------------------------------------------------------------------------
// 6. Factory function returning the class is auto-wrapped
// ---------------------------------------------------------------------------

func TestClassFactoryAutoWrapped(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	// counter.new() returns *Counter; the module-level Fn auto-wraps it.
	require.NoError(t, L.DoString(`
		local c = counter.new()
		-- If auto-wrap works, c is userdata and :increment/:get work.
		c:increment(7)
		assert(c:get() == 7, "expected 7, got " .. tostring(c:get()))
	`))
}

// ---------------------------------------------------------------------------
// 7. Method returning a non-class type goes through normal goToLua path
// ---------------------------------------------------------------------------

func TestClassNonClassReturnNormalPath(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	// Get() returns int — should be lua.LNumber, not userdata.
	require.NoError(t, L.DoString(`
		local c = counter.new()
		c:increment(42)
		local v = c:get()
		assert(type(v) == "number", "expected number, got " .. type(v))
		assert(v == 42)
	`))
}

// ---------------------------------------------------------------------------
// 8. Plain string error from a method — raises string, catchable by pcall
// ---------------------------------------------------------------------------

func TestClassPlainStringErrorFromMethod(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	require.NoError(t, L.DoString(`
		local c = counter.new()
		local ok, msg = pcall(function() c:increment_plain_error(-1) end)
		assert(ok == false, "expected failure")
		assert(type(msg) == "string", "expected string error, got " .. type(msg))
		assert(string.find(msg, "plain:") ~= nil, "expected 'plain:' in: " .. msg)
	`))
}

// ---------------------------------------------------------------------------
// 9. Structured *luareg.Error from a method — raises table; err.kind/message
// ---------------------------------------------------------------------------

func TestClassStructuredErrorFromMethod(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	require.NoError(t, L.DoString(`
		local c = counter.new()
		local ok, err = pcall(function() c:increment_boom(-5) end)
		assert(ok == false, "expected failure")
		assert(type(err) == "table", "expected table error, got " .. type(err))
		assert(err.kind == "NegativeDelta", "expected NegativeDelta, got " .. tostring(err.kind))
		assert(err.message == "delta must be non-negative", "unexpected message: " .. tostring(err.message))
	`))
}

// ---------------------------------------------------------------------------
// 10. Structured error from a Module.Fn (not a method)
// ---------------------------------------------------------------------------

func TestStructuredErrorFromModuleFn(t *testing.T) {
	m := luareg.NewModule("m", "test")
	m.Fn("boom", func(n int) (int, error) {
		if n < 0 {
			return 0, &luareg.Error{Kind: "Negative", Message: "must be non-negative"}
		}
		return n * 2, nil
	}, "")

	L := lua.NewState()
	defer L.Close()
	m.PushTo(L)
	L.SetGlobal("m", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		local ok, err = pcall(m.boom, -1)
		assert(ok == false)
		assert(type(err) == "table", "expected table, got " .. type(err))
		assert(err.kind == "Negative", "kind: " .. tostring(err.kind))
		assert(err.message == "must be non-negative", "msg: " .. tostring(err.message))

		-- Success path still works.
		local v = m.boom(3)
		assert(v == 6)
	`))
}

// ---------------------------------------------------------------------------
// 11. Duplicate method name on a class panics
// ---------------------------------------------------------------------------

func TestClassDuplicateMethodPanics(t *testing.T) {
	c := luareg.NewClass[*Counter]("test.Counter", "test")
	c.Method("get", (*Counter).Get, "first")
	assert.PanicsWithValue(t,
		`luareg: class "test.Counter": duplicate registration of method "get"`,
		func() { c.Method("get", (*Counter).Get, "second") },
	)
}

// ---------------------------------------------------------------------------
// 12. Duplicate class name on a module panics
// ---------------------------------------------------------------------------

func TestModuleDuplicateClassPanics(t *testing.T) {
	m := luareg.NewModule("m", "test")
	c1 := luareg.NewClass[*Counter]("Foo", "first")
	c2 := luareg.NewClass[*Counter]("Foo", "second")
	m.RegisterClass(c1)
	assert.PanicsWithValue(t,
		`luareg: module "m": duplicate class registration "Foo"`,
		func() { m.RegisterClass(c2) },
	)
}

// ---------------------------------------------------------------------------
// 13. Method with *lua.LState as second param (escape hatch)
// ---------------------------------------------------------------------------

func TestClassLStateEscapeHatch(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	require.NoError(t, L.DoString(`
		local c = counter.new()
		c:increment(99)
		local tbl = c:get_with_state()
		assert(type(tbl) == "table", "expected table")
		assert(tbl.value == 99, "expected 99, got " .. tostring(tbl.value))
	`))
}

// ---------------------------------------------------------------------------
// 14. Wrap and Check directly
// ---------------------------------------------------------------------------

func TestClassWrapAndCheck(t *testing.T) {
	c := luareg.NewClass[*Counter]("test.Counter", "test")
	c.Method("get", (*Counter).Get, "get")

	m := luareg.NewModule("m", "test")
	m.RegisterClass(c)

	L := lua.NewState()
	defer L.Close()
	m.PushTo(L)

	instance := &Counter{value: 42}
	ud := c.Wrap(L, instance)
	L.SetGlobal("ud", ud)

	// Check extracts the same pointer.
	require.NoError(t, L.DoString(`assert(type(ud) == "userdata")`))

	// Use Check via the method dispatch path.
	require.NoError(t, L.DoString(`
		local mt = getmetatable(ud)
		local v = mt.__index.get(ud)
		assert(v == 42, "expected 42, got " .. tostring(v))
	`))
}

// ---------------------------------------------------------------------------
// 15. Value-receiver type (NewClass[ValCounter])
// ---------------------------------------------------------------------------

func TestClassValueReceiverType(t *testing.T) {
	c := luareg.NewClass[ValCounter]("test.ValCounter", "value counter")
	c.Method("get_val", (ValCounter).GetVal, "get val")

	m := luareg.NewModule("m", "test")
	m.RegisterClass(c)
	m.Fn("new_val", func() ValCounter { return ValCounter{V: 7} }, "new val counter")

	L := lua.NewState()
	defer L.Close()
	m.PushTo(L)
	L.SetGlobal("m", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		local vc = m.new_val()
		local v = vc:get_val()
		assert(v == 7, "expected 7, got " .. tostring(v))
	`))
}

// ---------------------------------------------------------------------------
// 16. Class registered but module never PushTo'd — metatable doesn't exist
// ---------------------------------------------------------------------------

func TestClassNoMetatableWithoutPushTo(t *testing.T) {
	c := luareg.NewClass[*Counter]("missing.Counter", "test")
	c.Method("get", (*Counter).Get, "get")

	m := luareg.NewModule("m", "test")
	m.RegisterClass(c)
	// Do NOT call m.PushTo(L)

	L := lua.NewState()
	defer L.Close()

	// Without PushTo the metatable for "missing.Counter" doesn't exist in L.
	// Wrap still succeeds (returns a UserData with nil metatable from GetTypeMetatable),
	// but the method dispatch has no __index — calling a method should fail.
	instance := &Counter{value: 5}
	ud := c.Wrap(L, instance)
	L.SetGlobal("ud", ud)

	err := L.DoString(`ud:get()`)
	require.Error(t, err, "expected error when metatable not registered")
}

// ---------------------------------------------------------------------------
// 17. End-to-end: PushTo + method dispatch + auto-wrap
// ---------------------------------------------------------------------------

func TestClassEndToEnd(t *testing.T) {
	m := luareg.NewModule("app", "end-to-end test")

	cc := luareg.NewClass[*Counter]("app.Counter", "counter")
	cc.Method("increment", (*Counter).Increment, "add delta", luareg.Args("delta"))
	cc.Method("get", (*Counter).Get, "get value")
	cc.Method("reset", (*Counter).Reset, "reset")
	m.RegisterClass(cc)

	chc := luareg.NewClass[*Chainer]("app.Chainer", "chainer")
	chc.Method("with_label", (*Chainer).WithLabel, "new label", luareg.Args("label"))
	chc.Method("label", (*Chainer).Label, "get label")
	m.RegisterClass(chc)

	// Factory functions.
	m.Fn("new_counter", NewCounter, "create counter")
	m.Fn("new_chainer", NewChainer, "create chainer", luareg.Args("label"))

	L := lua.NewState()
	defer L.Close()
	n := m.PushTo(L)
	require.Equal(t, 1, n)
	L.SetGlobal("app", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		-- Counter lifecycle.
		local c = app.new_counter()
		c:increment(10)
		c:increment(5)
		assert(c:get() == 15)
		c:reset()
		assert(c:get() == 0)

		-- Two independent counters.
		local a = app.new_counter()
		local b = app.new_counter()
		a:increment(100)
		b:increment(1)
		assert(a:get() == 100)
		assert(b:get() == 1)

		-- Chainer with auto-wrapped return.
		local ch = app.new_chainer("hello")
		local ch2 = ch:with_label("world")
		assert(ch2:label() == "world")
		assert(ch:label() == "hello")  -- original unchanged
	`))
}

// ---------------------------------------------------------------------------
// 18. A1 plain-error compat: errors.New still raises as string
// ---------------------------------------------------------------------------

func TestA1PlainErrorCompat(t *testing.T) {
	// This duplicates the A1 TestErrorRaise test but makes the intent explicit:
	// plain errors.New must NOT become a table after the A2 raiseError changes.
	m := luareg.NewModule("m", "test")
	m.Fn("boom", func(n int) (int, error) {
		if n < 0 {
			return 0, fmt.Errorf("negative number")
		}
		return n * 2, nil
	}, "")

	L := lua.NewState()
	defer L.Close()
	m.PushTo(L)
	L.SetGlobal("m", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		local ok, msg = pcall(m.boom, -1)
		assert(ok == false)
		assert(type(msg) == "string", "A1 compat: expected string, got " .. type(msg))
		assert(string.find(msg, "negative number") ~= nil)
	`))
}

// ---------------------------------------------------------------------------
// 19. Method param docs/opts populate FnMeta correctly
// ---------------------------------------------------------------------------

func TestClassMethodMetadataOpts(t *testing.T) {
	c := luareg.NewClass[*Counter]("test.Counter", "counter")
	c.Method("increment", (*Counter).Increment, "add delta to counter",
		luareg.Args("delta"),
		luareg.ArgDoc("delta", "the amount to add"),
		luareg.ReturnDoc(0, "result", "unused"),
	)

	meta := c.MethodsMeta()
	require.Len(t, meta, 1)
	assert.Equal(t, "increment", meta[0].LuaName)
	assert.Equal(t, "delta", meta[0].ArgName(0))
	assert.Equal(t, "add delta to counter", meta[0].Doc)
}

// ---------------------------------------------------------------------------
// 20. Classes() accessor returns registered classes in order; Doc() accessible
// ---------------------------------------------------------------------------

func TestModuleClassesAccessor(t *testing.T) {
	m := luareg.NewModule("m", "test")
	c1 := luareg.NewClass[*Counter]("Counter", "counter doc")
	c2 := luareg.NewClass[*Chainer]("Chainer", "chainer doc")
	m.RegisterClass(c1)
	m.RegisterClass(c2)

	classes := m.Classes()
	require.Len(t, classes, 2)
	assert.Equal(t, "Counter", classes[0].Name())
	assert.Equal(t, "counter doc", classes[0].Doc())
	assert.Equal(t, "Chainer", classes[1].Name())
	assert.Equal(t, "chainer doc", classes[1].Doc())
}

// ---------------------------------------------------------------------------
// 21. Method with wrong receiver-type param panics at registration
// ---------------------------------------------------------------------------

func TestClassMethodWrongReceiverTypePanics(t *testing.T) {
	c := luareg.NewClass[*Counter]("test.Counter", "counter")
	// This method takes *Chainer, not *Counter — should panic.
	assert.Panics(t, func() {
		c.Method("bad", (*Chainer).Label, "wrong receiver type")
	})
}

// ---------------------------------------------------------------------------
// 22. Method arity mismatch raises an error
// ---------------------------------------------------------------------------

func TestClassMethodArityMismatch(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	// increment expects 1 arg (delta); call with 0.
	err := L.DoString(`
		local c = counter.new()
		c:increment()
	`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected 1 argument(s), got 0")

	// Call with too many args.
	err = L.DoString(`
		local c = counter.new()
		c:increment(1, 2)
	`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected 1 argument(s), got 2")
}

// ---------------------------------------------------------------------------
// 23. Error.Error() returns the message (implements error interface)
// ---------------------------------------------------------------------------

func TestLuaregErrorInterface(t *testing.T) {
	e := &luareg.Error{Kind: "TestKind", Message: "test message"}
	assert.Equal(t, "test message", e.Error())
}

// ---------------------------------------------------------------------------
// 24. Successful increment_boom path (no error raised)
// ---------------------------------------------------------------------------

func TestClassStructuredErrorNoError(t *testing.T) {
	L := newCounterState(t)
	defer L.Close()

	require.NoError(t, L.DoString(`
		local c = counter.new()
		c:increment_boom(5)
		assert(c:get() == 5)
	`))
}

// ---------------------------------------------------------------------------
// 25. *luareg.Error wrapped via fmt.Errorf %w is still detected via errors.As
// ---------------------------------------------------------------------------

// TestClassWrappedStructuredError: a *luareg.Error wrapped with fmt.Errorf %w
// must still raise a structured Lua table. The wrapping context is lost but
// kind/message survive — that's what pcall consumers care about.
func TestClassWrappedStructuredError(t *testing.T) {
	mod := luareg.NewModule("m", "test")
	mod.Fn("boom", func() (int, error) {
		inner := &luareg.Error{Kind: "NotFound", Message: "pod missing"}
		return 0, fmt.Errorf("looking up pod: %w", inner)
	}, "")

	L := lua.NewState()
	defer L.Close()
	mod.PushTo(L)
	L.SetGlobal("m", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		local ok, err = pcall(m.boom)
		assert(ok == false, "expected error")
		assert(type(err) == "table", "expected structured table, got " .. type(err))
		assert(err.kind == "NotFound", "kind: " .. tostring(err.kind))
		assert(err.message == "pod missing", "message: " .. tostring(err.message))
	`))
}

// ---------------------------------------------------------------------------
// 26. Ordering — Fn registered BEFORE RegisterClass still auto-wraps returns
// ---------------------------------------------------------------------------

// TestClassFactoryRegisteredBeforeClass: the classLookup closure captures the
// class map by reference, so a factory function registered before the class
// itself still gets auto-wrap behaviour at call time. This is the contract
// claimed in module.go's RegisterClass comment.
func TestClassFactoryRegisteredBeforeClass(t *testing.T) {
	type Box struct{ v int }
	mod := luareg.NewModule("m", "test")

	// Factory registered FIRST, before the class.
	mod.Fn("new_box", func() *Box { return &Box{v: 42} }, "")

	// Class registered AFTER the Fn. If the closure had captured the lookup
	// table by value, the factory would translate Box via JSON instead of
	// wrapping it as userdata.
	cls := luareg.NewClass[*Box]("m.Box", "a box")
	cls.Method("value", func(b *Box) int { return b.v }, "")
	mod.RegisterClass(cls)

	L := lua.NewState()
	defer L.Close()
	mod.PushTo(L)
	L.SetGlobal("m", L.Get(-1))
	L.Pop(1)

	// If auto-wrap worked, b is a userdata and b:value() dispatches. If not,
	// b is a translated table and :value fails (no metatable method).
	require.NoError(t, L.DoString(`
		local b = m.new_box()
		assert(b:value() == 42, "expected 42 from auto-wrapped box")
	`))
}

// ---------------------------------------------------------------------------
// 27. rawLGFunction escape hatch: func(T, *lua.LState) int
// ---------------------------------------------------------------------------

// TestClassRawLGFunctionDispatch: a method whose exact signature is
// func(*Counter, *lua.LState) int (rawLGFunction pattern) must dispatch
// correctly — the int return is the Lua stack count, not a value to push.
// This locks in the rawLGFunction path so a future refactor doesn't break it.
func TestClassRawLGFunctionDispatch(t *testing.T) {
	// rawStackMethod: pushes the sum of its manually-read args onto the Lua stack
	// and returns 1. Uses the rawLGFunction escape hatch signature.
	rawStackMethod := func(c *Counter, L *lua.LState) int {
		// Stack: pos 1 = self (Counter), pos 2 = addend (number).
		// (self already consumed by Check before this function is called)
		// Actually the full stack is visible here; pos 1 is self.
		addend := L.CheckInt(2)
		L.Push(lua.LNumber(c.value + addend))
		return 1
	}

	cls := luareg.NewClass[*Counter]("raw.Counter", "raw counter")
	cls.Method("get", (*Counter).Get, "get value")
	cls.Method("add_raw", rawStackMethod, "add via raw LGFunction")

	mod := luareg.NewModule("raw", "raw test module")
	mod.RegisterClass(cls)
	mod.Fn("new", NewCounter, "create counter")

	L := lua.NewState()
	defer L.Close()
	mod.PushTo(L)
	L.SetGlobal("raw", L.Get(-1))
	L.Pop(1)

	// Seed the counter with value=10 via Get method on a counter pre-seeded in Go.
	require.NoError(t, L.DoString(`
		local c = raw.new()
		c:get()  -- baseline 0
		-- The add_raw method uses the rawLGFunction path:
		-- it reads arg 2 off the stack itself and returns the computed sum.
		local result = c:add_raw(5)
		assert(type(result) == "number", "expected number, got " .. type(result))
		assert(result == 5, "0 + 5 = 5, got " .. tostring(result))
	`))
}

// ---------------------------------------------------------------------------
// 22. Class method with variadic tail
// ---------------------------------------------------------------------------

// AddMany: a variadic method on *Counter. Each value is added to the receiver.
func (c *Counter) AddMany(values ...int) int {
	for _, v := range values {
		c.value += v
	}
	return c.value
}

// TestClassMethodVariadic: a variadic method (receiver + ...int) collects the
// trailing Lua args into a []int, mutates the receiver, and returns the new
// total. Exercises empty/one/many cases through the same method.
func TestClassMethodVariadic(t *testing.T) {
	cls := luareg.NewClass[*Counter]("var.Counter", "counter with variadic")
	cls.Method("add_many", (*Counter).AddMany, "add many at once")
	cls.Method("get", (*Counter).Get, "get value")

	mod := luareg.NewModule("var", "variadic class module")
	mod.RegisterClass(cls)
	mod.Fn("new", NewCounter, "create counter")

	L := lua.NewState()
	defer L.Close()
	mod.PushTo(L)
	L.SetGlobal("var", L.Get(-1))
	L.Pop(1)

	require.NoError(t, L.DoString(`
		local c = var.new()
		assert(c:add_many() == 0, "empty tail should keep value at 0")
		assert(c:add_many(5) == 5, "single arg should give 5")
		assert(c:add_many(1, 2, 3) == 11, "1+2+3 over 5 should give 11")
	`))
}
