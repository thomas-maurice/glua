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

package luareg

import (
	"fmt"
	"reflect"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
)

var (
	errorType      = reflect.TypeOf((*error)(nil)).Elem()
	luaStateType   = reflect.TypeOf((*lua.LState)(nil))
	luaTableType   = reflect.TypeOf((*lua.LTable)(nil))
	luaValueIfType = reflect.TypeOf((*lua.LValue)(nil)).Elem() // the interface type
)

// buildWrapper: reflects on goFn and returns an lua.LGFunction that performs
// argument conversion, calls goFn, converts return values, and handles errors.
//
// classLookup maps a reflect.Type to an AnyClass for auto-wrap of return values.
// Pass nil if no class auto-wrap is needed.
//
// Panics immediately (at registration time) if goFn is not a function or its
// parameter/return types are not in the supported set. Variadic functions are
// supported: trailing Lua args are collected into a slice of the variadic
// element type before the Go call.
func buildWrapper(meta *FnMeta, classLookup func(reflect.Type) AnyClass) lua.LGFunction {
	fn := reflect.ValueOf(meta.GoFn)
	ft := fn.Type()

	// Determine whether the first param is *lua.LState (escape hatch).
	// If so, the Lua arity check ignores that param and we inject L ourselves.
	startsWithLState := ft.NumIn() > 0 && ft.In(0) == luaStateType
	isVariadic := ft.IsVariadic()

	// fixedLuaArgCount: number of Lua args that map to FIXED (non-variadic,
	// non-LState) Go params. The variadic tail (if any) consumes extras.
	fixedLuaArgCount := ft.NumIn()
	if startsWithLState {
		fixedLuaArgCount--
	}
	if isVariadic {
		fixedLuaArgCount--
	}

	// Determine whether the last return type is error.
	numOut := ft.NumOut()
	hasErrorReturn := numOut > 0 && ft.Out(numOut-1).Implements(errorType)
	nonErrReturns := numOut
	if hasErrorReturn {
		nonErrReturns--
	}

	// rawLGFunction: true when the signature is func(*lua.LState) int.
	// The function manages its own Lua stack pushes; the int return IS the
	// Lua return count. Variadic signatures never qualify.
	rawLGFunction := !isVariadic && startsWithLState && fixedLuaArgCount == 0 && numOut == 1 && ft.Out(0).Kind() == reflect.Int

	translator := glua.NewTranslator()

	return func(L *lua.LState) int {
		// Arity check. RaiseError (not ArgError) — arity is a call-level
		// problem, not an "argument #1 is bad" problem; ArgError would
		// produce a misleading "bad argument #1" prefix.
		// With *lua.LState escape hatch OR variadic: enforce minimum
		// (fixedLuaArgCount), allow extras. Without either: enforce exact.
		if startsWithLState || isVariadic {
			if fixedLuaArgCount > 0 && L.GetTop() < fixedLuaArgCount {
				L.RaiseError("expected at least %d argument(s), got %d", fixedLuaArgCount, L.GetTop())
				return 0
			}
		} else {
			if L.GetTop() != fixedLuaArgCount {
				L.RaiseError("expected %d argument(s), got %d", fixedLuaArgCount, L.GetTop())
				return 0
			}
		}

		// Build Go argument slice.
		args := make([]reflect.Value, ft.NumIn())
		argOffset := 0
		if startsWithLState {
			args[0] = reflect.ValueOf(L)
			argOffset = 1
		}

		for i := 0; i < fixedLuaArgCount; i++ {
			paramType := ft.In(i + argOffset)
			luaPos := i + 1 // Lua is 1-indexed
			val, err := luaToGo(L, luaPos, paramType, translator)
			if err != nil {
				L.ArgError(luaPos, err.Error())
				return 0
			}
			args[i+argOffset] = val
		}

		// Collect the variadic tail (if any) into a typed slice.
		if isVariadic {
			sliceType := ft.In(ft.NumIn() - 1) // e.g. []string
			elemType := sliceType.Elem()       // e.g. string
			top := L.GetTop()
			nVariadic := top - fixedLuaArgCount
			if nVariadic < 0 {
				nVariadic = 0
			}
			variadicSlice := reflect.MakeSlice(sliceType, nVariadic, nVariadic)
			for i := 0; i < nVariadic; i++ {
				luaPos := fixedLuaArgCount + i + 1
				val, err := luaToGo(L, luaPos, elemType, translator)
				if err != nil {
					L.ArgError(luaPos, err.Error())
					return 0
				}
				variadicSlice.Index(i).Set(val)
			}
			args[ft.NumIn()-1] = variadicSlice
		}

		// Call the Go function. CallSlice passes args[len-1] as the variadic
		// slice; Call handles the non-variadic case.
		var results []reflect.Value
		if isVariadic {
			results = fn.CallSlice(args)
		} else {
			results = fn.Call(args)
		}

		// Raw lua.LGFunction pattern: the function manages its own stack pushes
		// and the single int return IS the Lua stack count. Return it directly.
		if rawLGFunction {
			return int(results[0].Int())
		}

		// Handle error return (last result if hasErrorReturn).
		if hasErrorReturn {
			errVal := results[len(results)-1]
			if !errVal.IsNil() {
				raiseError(L, errVal.Interface().(error))
				return 0
			}
			results = results[:len(results)-1]
		}

		// Push non-error return values, auto-wrapping class instances.
		for _, rv := range results {
			lv, err := goToLuaWithClasses(L, rv, translator, classLookup)
			if err != nil {
				L.RaiseError("return conversion error: %s", err.Error())
				return 0
			}
			L.Push(lv)
		}

		return nonErrReturns
	}
}

// luaToGo: converts the Lua value at stack position pos to a Go reflect.Value
// of the requested type. Returns an error with a user-readable message on
// type mismatch so callers can use L.ArgError.
func luaToGo(L *lua.LState, pos int, t reflect.Type, tr *glua.Translator) (reflect.Value, error) {
	// lua.LValue interface pass-through: return the raw Lua value as-is.
	// Functions can declare a lua.LValue param to accept any Lua value.
	if t == luaValueIfType {
		lv := L.Get(pos)
		return reflect.ValueOf(&lv).Elem(), nil
	}

	switch t.Kind() {
	case reflect.String:
		// gopher-lua's CheckString coerces numbers; enforce actual type.
		if lv := L.Get(pos); lv.Type() != lua.LTString {
			return reflect.Value{}, fmt.Errorf("arg %d: expected string, got %s", pos, lv.Type())
		}
		v := L.CheckString(pos)
		return reflect.ValueOf(v), nil

	case reflect.Bool:
		v := L.CheckBool(pos)
		return reflect.ValueOf(v), nil

	case reflect.Int:
		v := L.CheckInt(pos)
		return reflect.ValueOf(v), nil

	case reflect.Int8:
		v := L.CheckInt(pos)
		return reflect.ValueOf(int8(v)), nil

	case reflect.Int16:
		v := L.CheckInt(pos)
		return reflect.ValueOf(int16(v)), nil

	case reflect.Int32:
		v := L.CheckInt(pos)
		return reflect.ValueOf(int32(v)), nil

	case reflect.Int64:
		v := L.CheckInt64(pos)
		return reflect.ValueOf(v), nil

	// uint types: gopher-lua has no CheckUint; read as int and convert.
	// Values that don't fit (negative or > max) will wrap silently — this
	// matches the Lua numeric model where all numbers are float64.
	case reflect.Uint:
		v := L.CheckInt(pos)
		return reflect.ValueOf(uint(v)), nil //nolint:gosec

	case reflect.Uint8:
		v := L.CheckInt(pos)
		return reflect.ValueOf(uint8(v)), nil //nolint:gosec

	case reflect.Uint16:
		v := L.CheckInt(pos)
		return reflect.ValueOf(uint16(v)), nil //nolint:gosec

	case reflect.Uint32:
		v := L.CheckInt(pos)
		return reflect.ValueOf(uint32(v)), nil //nolint:gosec

	case reflect.Uint64:
		// L.CheckInt returns int; use CheckNumber for uint64 to preserve range.
		v := L.CheckNumber(pos)
		return reflect.ValueOf(uint64(v)), nil //nolint:gosec

	case reflect.Float32:
		v := L.CheckNumber(pos)
		return reflect.ValueOf(float32(v)), nil

	case reflect.Float64:
		v := L.CheckNumber(pos)
		return reflect.ValueOf(float64(v)), nil

	case reflect.Slice:
		lv := L.CheckTable(pos)
		return luaTableToSlice(L, lv, t, tr, pos)

	case reflect.Map:
		lv := L.CheckTable(pos)
		ptr := reflect.New(t)
		if err := tr.FromLua(L, lv, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: cannot convert table to map: %s", pos, err.Error())
		}
		return ptr.Elem(), nil

	case reflect.Struct:
		lv := L.CheckTable(pos)
		ptr := reflect.New(t)
		if err := tr.FromLua(L, lv, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: cannot convert table to struct: %s", pos, err.Error())
		}
		return ptr.Elem(), nil

	case reflect.Pointer:
		// *lua.LTable is a pass-through: if the Lua value is nil, return a nil
		// *lua.LTable (allowing optional table params). Otherwise check it is a
		// table and return the pointer directly without going through the
		// Translator. This mirrors the *lua.LState escape hatch.
		if t == luaTableType {
			lv := L.Get(pos)
			if lv == lua.LNil || lv.Type() == lua.LTNil {
				return reflect.Zero(t), nil
			}
			tbl := L.CheckTable(pos)
			return reflect.ValueOf(tbl), nil
		}

		// For other pointer-to-struct types, convert via Translator.
		if t.Elem().Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("arg %d: unsupported pointer type %s", pos, t)
		}
		lv := L.CheckTable(pos)
		ptr := reflect.New(t.Elem())
		if err := tr.FromLua(L, lv, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: cannot convert table to %s: %s", pos, t, err.Error())
		}
		return ptr, nil

	default:
		return reflect.Value{}, fmt.Errorf("arg %d: unsupported type %s", pos, t)
	}
}

// luaTableToSlice: converts a Lua table (array portion) into a Go slice of
// the element type inferred from the slice type t.
func luaTableToSlice(L *lua.LState, tbl *lua.LTable, t reflect.Type, tr *glua.Translator, argPos int) (reflect.Value, error) {
	elemType := t.Elem()
	n := tbl.MaxN()
	slice := reflect.MakeSlice(t, n, n)
	for i := 1; i <= n; i++ {
		raw := tbl.RawGetInt(i)
		elem, err := luaToGoValue(L, raw, elemType, tr, argPos, i)
		if err != nil {
			return reflect.Value{}, err
		}
		slice.Index(i - 1).Set(elem)
	}
	return slice, nil
}

// luaToGoValue: converts a raw lua.LValue to a reflect.Value of the given Go
// type. Used for slice element conversion where we already have the LValue.
func luaToGoValue(L *lua.LState, lv lua.LValue, t reflect.Type, tr *glua.Translator, argPos, elemIdx int) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.String:
		s, ok := lv.(lua.LString)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected string, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(string(s)), nil

	case reflect.Bool:
		b, ok := lv.(lua.LBool)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected boolean, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(bool(b)), nil

	case reflect.Int:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(int(n)), nil

	case reflect.Int8:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(int8(n)), nil //nolint:gosec

	case reflect.Int16:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(int16(n)), nil //nolint:gosec

	case reflect.Int32:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(int32(n)), nil //nolint:gosec

	case reflect.Int64:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(int64(n)), nil

	case reflect.Uint:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(uint(n)), nil //nolint:gosec

	case reflect.Uint8:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(uint8(n)), nil //nolint:gosec

	case reflect.Uint16:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(uint16(n)), nil //nolint:gosec

	case reflect.Uint32:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(uint32(n)), nil //nolint:gosec

	case reflect.Uint64:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(uint64(n)), nil //nolint:gosec

	case reflect.Float32:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(float32(n)), nil

	case reflect.Float64:
		n, ok := lv.(lua.LNumber)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected number, got %s", argPos, elemIdx, lv.Type())
		}
		return reflect.ValueOf(float64(n)), nil

	case reflect.Struct:
		// Composite element: convert the inner table via the Translator.
		// Mirrors the top-level Struct case in luaToGo.
		tbl, ok := lv.(*lua.LTable)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected table, got %s", argPos, elemIdx, lv.Type())
		}
		ptr := reflect.New(t)
		if err := tr.FromLua(L, tbl, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: cannot convert table to %s: %s", argPos, elemIdx, t, err.Error())
		}
		return ptr.Elem(), nil

	case reflect.Pointer:
		// Pointer element: only *struct is supported (matches luaToGo).
		// A Lua-nil entry materializes as a typed nil pointer.
		if t.Elem().Kind() != reflect.Struct {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: unsupported pointer type %s", argPos, elemIdx, t)
		}
		if lv == lua.LNil || lv.Type() == lua.LTNil {
			return reflect.Zero(t), nil
		}
		tbl, ok := lv.(*lua.LTable)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected table, got %s", argPos, elemIdx, lv.Type())
		}
		ptr := reflect.New(t.Elem())
		if err := tr.FromLua(L, tbl, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: cannot convert table to %s: %s", argPos, elemIdx, t, err.Error())
		}
		return ptr, nil

	case reflect.Slice, reflect.Map:
		// Nested slice/map element: delegate to the Translator. This makes
		// [][]T, []map[K]V, map[K][]V work on the input path.
		tbl, ok := lv.(*lua.LTable)
		if !ok {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: expected table, got %s", argPos, elemIdx, lv.Type())
		}
		ptr := reflect.New(t)
		if err := tr.FromLua(L, tbl, ptr.Interface()); err != nil {
			return reflect.Value{}, fmt.Errorf("arg %d: table element %d: cannot convert table to %s: %s", argPos, elemIdx, t, err.Error())
		}
		return ptr.Elem(), nil

	default:
		return reflect.Value{}, fmt.Errorf("arg %d: table element %d: unsupported element type %s", argPos, elemIdx, t)
	}
}

// goToLua: converts a reflect.Value produced by a Go function into an
// lua.LValue suitable for pushing onto the Lua stack.
func goToLua(L *lua.LState, v reflect.Value, tr *glua.Translator) (lua.LValue, error) {
	// Fast path: if the value already implements lua.LValue, use it directly.
	// This handles both nilable LValue kinds (*lua.LTable, *lua.LUserData) and
	// value kinds (lua.LString, lua.LNumber, lua.LBool). We only call IsNil
	// on kinds where it is defined; otherwise the value is non-nil by
	// construction.
	if v.Type().Implements(luaValueIfType) {
		switch v.Kind() {
		case reflect.Pointer, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
			if v.IsNil() {
				return lua.LNil, nil
			}
		}
		return v.Interface().(lua.LValue), nil
	}

	// Dereference non-LValue pointers (e.g. *MyStruct).
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return lua.LNil, nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return lua.LString(v.String()), nil

	case reflect.Bool:
		return lua.LBool(v.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(v.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(float64(v.Uint())), nil

	case reflect.Float32, reflect.Float64:
		return lua.LNumber(v.Float()), nil

	case reflect.Slice, reflect.Map, reflect.Struct:
		lv, err := tr.ToLua(L, v.Interface())
		if err != nil {
			return nil, fmt.Errorf("cannot convert %s to Lua: %w", v.Type(), err)
		}
		return lv, nil

	default:
		return nil, fmt.Errorf("unsupported return type %s", v.Type())
	}
}

// validateGoFn: panics with a clear message if goFn is nil or not a function.
// Variadic functions are accepted — see buildWrapper / buildMethodWrapper for
// the runtime handling of the variadic tail.
func validateGoFn(luaName string, goFn any) {
	if goFn == nil {
		panic(fmt.Sprintf("luareg: Fn %q: goFn must not be nil", luaName))
	}
	t := reflect.TypeOf(goFn)
	if t.Kind() != reflect.Func {
		panic(fmt.Sprintf("luareg: Fn %q: goFn must be a function, got %s", luaName, t.Kind()))
	}
}
