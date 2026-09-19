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

// Package collections provides table utilities for Lua scripts: the
// functional trio (map/filter/reduce), grouping/sorting/shaping helpers, and
// deep structural operations (deep_equal, deep_copy, get_path).
//
// Performance note: this is an ergonomics module, not a performance one.
// collections.map(t, fn) is slower than `for i, v in ipairs(t) do ... end`
// in pure Lua, because every element crosses Lua -> Go -> Lua through
// L.CallByParam, where a plain Lua loop never leaves the VM. Reach for the
// callback-based functions for clarity; drop to a hand-written loop in a hot
// path.
//
// This module deliberately does NOT use pkg/glua.Translator (no JSON
// round-trip). It operates directly on *lua.LTable with lua.LValue elements,
// which preserves number precision, non-string keys, function/userdata
// values, and table identity — all of which the Translator's JSON path would
// lose or distort. The cost is that every function below must say, in its
// own doc comment, exactly which part of the table it reads.
//
// Two iteration modes, used consistently:
//
//   - array mode: `for i := 1; i <= t.MaxN(); i++`. Only the contiguous
//     integer-keyed prefix is visited; any hash-part keys are ignored, not
//     an error. Used by map, filter, reduce, find, any, all, group_by,
//     sort_by, partition, uniq, flatten, reverse, zip, chunk.
//   - map mode: `t.ForEach`, which visits every key (array part and hash
//     part together) in gopher-lua's internal, UNSPECIFIED order. Used by
//     keys, values, merge, pick, omit, get_path. Do not rely on the order
//     values come back in.
//   - deep_equal and deep_copy use map mode (recursively), so they cover
//     both the array and hash parts of every nested table.
//
// Callback conventions, followed by every callback-accepting function:
//
//   - Only a real Lua function (*lua.LFunction) is accepted; a table with a
//     __call metamethod is rejected with an immediate error naming the
//     parameter.
//   - Callbacks that take an element also take its 1-based array position,
//     always in (value, index) order, e.g. `fun(v: any, i: integer): any`.
//     This is consistent across map/filter/reduce/find/any/all/group_by/
//     partition. sort_by is the one exception (`fun(v: any): number|string`
//     — position is not part of the sort key).
//   - A callback that raises (calls error(), indexes nil, etc.) does not
//     panic the host: it is invoked with L.CallByParam's Protect: true, and
//     the Lua error is re-raised as a Go error wrapped with function name and
//     the 1-based index at which it failed, e.g.
//     "collections.filter: callback failed at index 7: <message>".
//
// Empty results are always an empty *lua.LTable, never Lua nil — see the
// nil-slice rule in the working-here skill. Every function that can produce
// an empty result is tested for it.
//
// uniq equality: primitives (string/number/boolean) compare by value; tables,
// functions and userdata compare by identity (pointer). This falls out of Go
// interface equality for free: lua.LValue is an interface, so `==` on two
// LValues does value-equality for value-kind dynamic types (LNumber, LString,
// LBool) and pointer-equality for pointer-kind dynamic types (*lua.LTable,
// *lua.LFunction, *lua.LUserData) — exactly the split the spec calls for.
// This is NOT structural equality; use deep_equal in an O(n^2) loop for that.
//
// sort_by is stable (sort.SliceStable) and requires every element's computed
// key to be entirely numbers or entirely strings — a mixed set raises rather
// than inventing a cross-type ordering. A comparator failure (callback raise
// or type mismatch) is detected before sort.SliceStable ever runs, so a
// raising key function cannot leave a table half-sorted.
//
// get_path splits path on "." (so the empty string is a single empty-string
// segment, looked up as t[""] — normally absent, so get_path(t, "", d)
// returns d). A segment that parses as an integer is looked up as a number
// key first, then as a string key if that misses — so
// get_path(t, "items.1", d) finds an array-style t.items[1] as well as a
// literal string key "1". Traversal stops and returns default the moment a
// key is missing or a non-table is indexed.
//
// deep_equal and deep_copy both guard against cyclic tables: deep_equal
// tracks visited (*lua.LTable, *lua.LTable) pairs and treats a revisited pair
// as equal (so two cyclic-but-otherwise-equal structures compare equal
// instead of looping forever); deep_copy tracks source table -> copy table
// so a cycle in the input is reproduced as the same cycle shape in the copy,
// not an infinite recursion. Neither compares nor copies metatables.
package collections

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// checkCallback: validates that v is a real Lua function, raising a
// descriptive error naming fnName and paramName otherwise. __call tables are
// deliberately not supported (F5).
func checkCallback(v lua.LValue, fnName, paramName string) (*lua.LFunction, error) {
	fn, ok := v.(*lua.LFunction)
	if !ok {
		return nil, fmt.Errorf("%s: %s must be a function, got %s", fnName, paramName, v.Type())
	}
	return fn, nil
}

// invoke: calls fn with args via a protected call, so a Lua error raised
// inside the callback (error(), a nil index, etc.) cannot panic the host. On
// failure the error is wrapped with fnName and the 1-based idx at which the
// callback was invoked, e.g. "collections.map: callback failed at index 3:
// <message>". idx is caller-supplied context for the error message only; it
// need not be one of args.
func invoke(L *lua.LState, fnName string, fn *lua.LFunction, idx int, args ...lua.LValue) (lua.LValue, error) {
	if err := L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}, args...); err != nil {
		return nil, fmt.Errorf("%s: callback failed at index %d: %s", fnName, idx, err.Error())
	}
	ret := L.Get(-1)
	L.Pop(1)
	return ret, nil
}

// mapFn: array mode. Applies fn(v, i) to every element of t and returns a new
// array of the results, same length as t (a callback returning nil produces
// a trailing gap per normal Lua array semantics, not an error).
func mapFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (*lua.LTable, error) {
	fn, err := checkCallback(fnArg, "collections.map", "fn")
	if err != nil {
		return nil, err
	}
	n := t.MaxN()
	result := L.NewTable()
	for i := 1; i <= n; i++ {
		ret, err := invoke(L, "collections.map", fn, i, t.RawGetInt(i), lua.LNumber(i))
		if err != nil {
			return nil, err
		}
		result.RawSetInt(i, ret)
	}
	return result, nil
}

// filterFn: array mode. Returns a new, contiguously re-indexed array of the
// elements for which fn(v, i) is truthy.
func filterFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (*lua.LTable, error) {
	fn, err := checkCallback(fnArg, "collections.filter", "fn")
	if err != nil {
		return nil, err
	}
	n := t.MaxN()
	result := L.NewTable()
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		ret, err := invoke(L, "collections.filter", fn, i, v, lua.LNumber(i))
		if err != nil {
			return nil, err
		}
		if lua.LVAsBool(ret) {
			result.Append(v)
		}
	}
	return result, nil
}

// reduceFn: array mode. Folds fn(acc, v, i) left-to-right over t, starting
// from the required init value, and returns the final accumulator.
func reduceFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue, initial lua.LValue) (lua.LValue, error) {
	fn, err := checkCallback(fnArg, "collections.reduce", "fn")
	if err != nil {
		return nil, err
	}
	acc := initial
	n := t.MaxN()
	for i := 1; i <= n; i++ {
		ret, err := invoke(L, "collections.reduce", fn, i, acc, t.RawGetInt(i), lua.LNumber(i))
		if err != nil {
			return nil, err
		}
		acc = ret
	}
	return acc, nil
}

// findFn: array mode. Returns the first element (and its 1-based index) for
// which fn(v, i) is truthy. Returns nil, 0 if no element matches — a
// legitimate negative answer, not an error (F1).
func findFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (lua.LValue, int, error) {
	fn, err := checkCallback(fnArg, "collections.find", "fn")
	if err != nil {
		return nil, 0, err
	}
	n := t.MaxN()
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		ret, err := invoke(L, "collections.find", fn, i, v, lua.LNumber(i))
		if err != nil {
			return nil, 0, err
		}
		if lua.LVAsBool(ret) {
			return v, i, nil
		}
	}
	return lua.LNil, 0, nil
}

// anyFn: array mode. Reports whether fn(v, i) is truthy for at least one
// element, short-circuiting on the first match. An empty table is not any.
func anyFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (bool, error) {
	fn, err := checkCallback(fnArg, "collections.any", "fn")
	if err != nil {
		return false, err
	}
	n := t.MaxN()
	for i := 1; i <= n; i++ {
		ret, err := invoke(L, "collections.any", fn, i, t.RawGetInt(i), lua.LNumber(i))
		if err != nil {
			return false, err
		}
		if lua.LVAsBool(ret) {
			return true, nil
		}
	}
	return false, nil
}

// allFn: array mode. Reports whether fn(v, i) is truthy for every element,
// short-circuiting on the first failure. An empty table is vacuously all.
func allFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (bool, error) {
	fn, err := checkCallback(fnArg, "collections.all", "fn")
	if err != nil {
		return false, err
	}
	n := t.MaxN()
	for i := 1; i <= n; i++ {
		ret, err := invoke(L, "collections.all", fn, i, t.RawGetInt(i), lua.LNumber(i))
		if err != nil {
			return false, err
		}
		if !lua.LVAsBool(ret) {
			return false, nil
		}
	}
	return true, nil
}

// groupByFn: array mode. Groups elements of t by the key returned from
// fn(v, i), preserving each group's encounter order. The result maps each
// distinct key (any Lua value, used as a raw table key) to a new array of
// the elements that produced it.
func groupByFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (*lua.LTable, error) {
	fn, err := checkCallback(fnArg, "collections.group_by", "fn")
	if err != nil {
		return nil, err
	}
	n := t.MaxN()
	result := L.NewTable()
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		key, err := invoke(L, "collections.group_by", fn, i, v, lua.LNumber(i))
		if err != nil {
			return nil, err
		}
		group, ok := result.RawGet(key).(*lua.LTable)
		if !ok {
			group = L.NewTable()
			result.RawSet(key, group)
		}
		group.Append(v)
	}
	return result, nil
}

// sortByFn: array mode. Returns a new array sorted by the key returned from
// fn(v) (position is not part of the key), using sort.SliceStable so elements
// with equal keys keep their relative input order. Every key must be a
// number or a string, and all keys must share the same type — a mixed set
// raises "key type mismatch" naming the offending index, rather than
// inventing a cross-type ordering.
func sortByFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (*lua.LTable, error) {
	fn, err := checkCallback(fnArg, "collections.sort_by", "fn")
	if err != nil {
		return nil, err
	}
	n := t.MaxN()
	type item struct {
		value lua.LValue
		key   lua.LValue
	}
	items := make([]item, n)
	var keyKind lua.LValueType
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		key, err := invoke(L, "collections.sort_by", fn, i, v)
		if err != nil {
			return nil, err
		}
		switch key.Type() {
		case lua.LTNumber, lua.LTString:
		default:
			return nil, fmt.Errorf("collections.sort_by: key must be a number or string, got %s at index %d", key.Type(), i)
		}
		if i == 1 {
			keyKind = key.Type()
		} else if key.Type() != keyKind {
			return nil, fmt.Errorf("collections.sort_by: key type mismatch at index %d: got %s, expected %s", i, key.Type(), keyKind)
		}
		items[i-1] = item{value: v, key: key}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if keyKind == lua.LTNumber {
			return items[i].key.(lua.LNumber) < items[j].key.(lua.LNumber)
		}
		return items[i].key.(lua.LString) < items[j].key.(lua.LString)
	})
	result := L.NewTable()
	for _, it := range items {
		result.Append(it.value)
	}
	return result, nil
}

// partitionFn: array mode. Splits t into two new arrays: elements for which
// fn(v, i) is truthy (matching), and the rest, each contiguously re-indexed
// and preserving relative order.
func partitionFn(L *lua.LState, t *lua.LTable, fnArg lua.LValue) (*lua.LTable, *lua.LTable, error) {
	fn, err := checkCallback(fnArg, "collections.partition", "fn")
	if err != nil {
		return nil, nil, err
	}
	n := t.MaxN()
	matching := L.NewTable()
	rest := L.NewTable()
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		ret, err := invoke(L, "collections.partition", fn, i, v, lua.LNumber(i))
		if err != nil {
			return nil, nil, err
		}
		if lua.LVAsBool(ret) {
			matching.Append(v)
		} else {
			rest.Append(v)
		}
	}
	return matching, rest, nil
}

// uniqFn: array mode. Returns a new array with duplicate elements removed,
// keeping the first occurrence. Primitives (string/number/boolean) are
// deduplicated by value; tables, functions and userdata by identity — this
// is NOT deep_equal-style structural equality.
func uniqFn(L *lua.LState, t *lua.LTable) *lua.LTable {
	n := t.MaxN()
	seen := make(map[lua.LValue]bool, n)
	result := L.NewTable()
	for i := 1; i <= n; i++ {
		v := t.RawGetInt(i)
		if !seen[v] {
			seen[v] = true
			result.Append(v)
		}
	}
	return result
}

// flattenFn: array mode, recursively. Nested tables (found via array-mode
// iteration, i.e. their own array part) are inlined up to depth levels; depth
// -1 flattens fully. depth must be >= -1.
func flattenFn(L *lua.LState, t *lua.LTable, depth int) (*lua.LTable, error) {
	if depth < -1 {
		return nil, fmt.Errorf("collections.flatten: depth must be >= -1 (-1 means fully), got %d", depth)
	}
	result := L.NewTable()
	var walk func(tbl *lua.LTable, level int)
	walk = func(tbl *lua.LTable, level int) {
		n := tbl.MaxN()
		for i := 1; i <= n; i++ {
			v := tbl.RawGetInt(i)
			if sub, ok := v.(*lua.LTable); ok && (depth == -1 || level < depth) {
				walk(sub, level+1)
				continue
			}
			result.Append(v)
		}
	}
	walk(t, 0)
	return result, nil
}

// reverseFn: array mode. Returns a new array with t's elements in reverse
// order.
func reverseFn(L *lua.LState, t *lua.LTable) *lua.LTable {
	n := t.MaxN()
	result := L.NewTable()
	for i := n; i >= 1; i-- {
		result.Append(t.RawGetInt(i))
	}
	return result
}

// zipFn: array mode on both a and b. Returns a new array of 2-element arrays
// {a[i], b[i]}, of length min(#a, #b) — extra elements in the longer table
// are dropped, not padded with nil.
func zipFn(L *lua.LState, a, b *lua.LTable) *lua.LTable {
	n := min(a.MaxN(), b.MaxN())
	result := L.NewTable()
	for i := 1; i <= n; i++ {
		pair := L.NewTable()
		pair.Append(a.RawGetInt(i))
		pair.Append(b.RawGetInt(i))
		result.Append(pair)
	}
	return result
}

// chunkFn: array mode. Splits t into new arrays of at most size elements
// each, in order; the last chunk may be shorter than size. size must be >= 1.
func chunkFn(L *lua.LState, t *lua.LTable, size int) (*lua.LTable, error) {
	if size < 1 {
		return nil, fmt.Errorf("collections.chunk: size must be >= 1, got %d", size)
	}
	n := t.MaxN()
	result := L.NewTable()
	var cur *lua.LTable
	for i := 1; i <= n; i++ {
		if (i-1)%size == 0 {
			cur = L.NewTable()
			result.Append(cur)
		}
		cur.Append(t.RawGetInt(i))
	}
	return result, nil
}

// keysFn: map mode. Returns a new array of all of t's keys (array part and
// hash part together), in gopher-lua's unspecified ForEach order.
func keysFn(L *lua.LState, t *lua.LTable) *lua.LTable {
	result := L.NewTable()
	t.ForEach(func(k, _ lua.LValue) {
		result.Append(k)
	})
	return result
}

// valuesFn: map mode. Returns a new array of all of t's values (array part
// and hash part together), in gopher-lua's unspecified ForEach order.
func valuesFn(L *lua.LState, t *lua.LTable) *lua.LTable {
	result := L.NewTable()
	t.ForEach(func(_, v lua.LValue) {
		result.Append(v)
	})
	return result
}

// mergeFn: map mode, shallow. Returns a new table containing every key from
// first and rest, applied left-to-right so a later table's value for a
// shared key wins. first and rest are never modified.
func mergeFn(L *lua.LState, first *lua.LTable, rest ...*lua.LTable) *lua.LTable {
	result := L.NewTable()
	apply := func(tbl *lua.LTable) {
		tbl.ForEach(func(k, v lua.LValue) {
			result.RawSet(k, v)
		})
	}
	apply(first)
	for _, tbl := range rest {
		apply(tbl)
	}
	return result
}

// pickFn: map mode over names (a string[]), not over all of t's keys. Returns
// a new table containing only the string-keyed fields of t named in names;
// a name with no corresponding entry in t is silently omitted.
func pickFn(L *lua.LState, t *lua.LTable, names []string) *lua.LTable {
	result := L.NewTable()
	for _, name := range names {
		if v := t.RawGetString(name); v != lua.LNil {
			result.RawSetString(name, v)
		}
	}
	return result
}

// omitFn: map mode over all of t's keys. Returns a new table with every
// key/value from t except the string keys named in names; non-string keys
// are never omitted (names can only match string keys).
func omitFn(L *lua.LState, t *lua.LTable, names []string) *lua.LTable {
	drop := make(map[string]bool, len(names))
	for _, name := range names {
		drop[name] = true
	}
	result := L.NewTable()
	t.ForEach(func(k, v lua.LValue) {
		if ks, ok := k.(lua.LString); ok && drop[string(ks)] {
			return
		}
		result.RawSet(k, v)
	})
	return result
}

// lookupPathSegment: map mode. Looks up a single dotted-path segment on tbl.
// A segment that parses as an integer is tried as a number key first, then
// falls back to a string key (so both array-style and literal-digit-string
// keys are reachable). Returns lua.LNil if neither is present.
func lookupPathSegment(tbl *lua.LTable, segment string) lua.LValue {
	if n, err := strconv.Atoi(segment); err == nil {
		if v := tbl.RawGetInt(n); v != lua.LNil {
			return v
		}
	}
	return tbl.RawGetString(segment)
}

// getPathFn: map mode, per segment. Walks t following path's dot-separated
// segments (see lookupPathSegment for the numeric/string ambiguity rule) and
// returns the value found, or def the moment a segment is missing or a
// non-table is indexed. def is a required argument (D4): pass nil explicitly
// for "no default", since `x or default` silently misfires when the stored
// value is boolean false.
func getPathFn(t *lua.LTable, path string, def lua.LValue) lua.LValue {
	var cur lua.LValue = t
	for _, segment := range strings.Split(path, ".") {
		tbl, ok := cur.(*lua.LTable)
		if !ok {
			return def
		}
		cur = lookupPathSegment(tbl, segment)
		if cur == lua.LNil {
			return def
		}
	}
	return cur
}

// deepEqualValues: recursively compares a and b. Non-table values compare
// with Lua's own value equality (interface ==, which is by-value for
// value-kind LValues and by-pointer for pointer-kind ones). Tables compare
// key-by-key via map mode; visited guards against infinite recursion on a
// cyclic table by treating a revisited (a, b) pair as equal.
func deepEqualValues(a, b lua.LValue, visited map[[2]*lua.LTable]bool) bool {
	if a == b {
		return true
	}
	ta, aok := a.(*lua.LTable)
	tb, bok := b.(*lua.LTable)
	if !aok || !bok {
		// Different dynamic types, or non-table values that were not == above.
		return false
	}
	pair := [2]*lua.LTable{ta, tb}
	if visited[pair] {
		return true
	}
	visited[pair] = true

	equal := true
	countA, countB := 0, 0
	ta.ForEach(func(k, v lua.LValue) {
		countA++
		if !equal {
			return
		}
		if !deepEqualValues(v, tb.RawGet(k), visited) {
			equal = false
		}
	})
	if !equal {
		return false
	}
	tb.ForEach(func(_, _ lua.LValue) { countB++ })
	return countA == countB
}

// deepEqualFn: general (recursive map mode). Reports whether a and b are
// structurally equal: same keys, and recursively deep_equal values, at every
// level. Metatables are not considered. Cyclic tables do not hang the
// comparison — see deepEqualValues.
func deepEqualFn(a, b lua.LValue) bool {
	return deepEqualValues(a, b, make(map[[2]*lua.LTable]bool))
}

// deepCopyValue: recursively copies v. Non-table values (including functions
// and userdata, which have no deep structure to copy) are returned as-is —
// they are copied "by reference" in the sense that no new Go/Lua value is
// allocated for them. Tables are copied key and value alike via map mode;
// copied tracks source -> new table so a cycle in the input reproduces the
// same cycle shape in the output instead of recursing forever.
func deepCopyValue(L *lua.LState, v lua.LValue, copied map[*lua.LTable]*lua.LTable) lua.LValue {
	tbl, ok := v.(*lua.LTable)
	if !ok {
		return v
	}
	if existing, ok := copied[tbl]; ok {
		return existing
	}
	newTbl := L.NewTable()
	copied[tbl] = newTbl
	tbl.ForEach(func(k, val lua.LValue) {
		newTbl.RawSet(deepCopyValue(L, k, copied), deepCopyValue(L, val, copied))
	})
	return newTbl
}

// deepCopyFn: general (recursive map mode). Returns a structural copy of v:
// nested tables are new tables, not aliases. Metatables are not copied.
// Functions and userdata are copied by reference — they are not
// deep-copyable. Cyclic tables are preserved, not expanded infinitely — see
// deepCopyValue.
func deepCopyFn(L *lua.LState, v lua.LValue) lua.LValue {
	return deepCopyValue(L, v, make(map[*lua.LTable]*lua.LTable))
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("collections", "table utilities: functional trio, grouping/sorting/shaping, and deep structural operations — see the package doc for array-vs-map mode and performance notes")

	m.Fn("map", mapFn, "array mode: applies fn to every element, returning a new same-length array of the results",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only (indices 1..#t)"),
		luareg.ArgDoc("fn", "callback invoked as fn(v, i); its return value becomes the new element at i"),
		luareg.ArgType("fn", "fun(v: any, i: integer): any"),
		luareg.ReturnDoc(0, "result", "a new array, same length as t"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("filter", filterFn, "array mode: returns a new, re-indexed array of the elements for which fn is truthy",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "predicate invoked as fn(v, i); the element is kept when this is truthy"),
		luareg.ArgType("fn", "fun(v: any, i: integer): boolean"),
		luareg.ReturnDoc(0, "result", "a new array of the matching elements, in original order"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("reduce", reduceFn, "array mode: folds fn left-to-right over t starting from init, returning the final accumulator",
		luareg.Args("t", "fn", "init"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "callback invoked as fn(acc, v, i); its return value becomes the next acc"),
		luareg.ArgType("fn", "fun(acc: any, v: any, i: integer): any"),
		luareg.ArgDoc("init", "the required initial accumulator; pass nil explicitly if that is the desired seed"),
		luareg.ReturnDoc(0, "result", "the final accumulator after folding every element"))

	m.Fn("find", findFn, "array mode: returns the first element and index for which fn is truthy, or nil, 0 if none",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "predicate invoked as fn(v, i)"),
		luareg.ArgType("fn", "fun(v: any, i: integer): boolean"),
		luareg.ReturnDoc(0, "value", "the first matching element, or nil if none matched"),
		luareg.ReturnDoc(1, "index", "the 1-based index of value, or 0 if none matched"))

	m.Fn("any", anyFn, "array mode: reports whether fn is truthy for at least one element (short-circuits)",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "predicate invoked as fn(v, i)"),
		luareg.ArgType("fn", "fun(v: any, i: integer): boolean"),
		luareg.ReturnDoc(0, "result", "true if any element matched; false for an empty table"))

	m.Fn("all", allFn, "array mode: reports whether fn is truthy for every element (short-circuits, vacuously true when empty)",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "predicate invoked as fn(v, i)"),
		luareg.ArgType("fn", "fun(v: any, i: integer): boolean"),
		luareg.ReturnDoc(0, "result", "true if every element matched, or the table is empty"))

	m.Fn("group_by", groupByFn, "array mode: groups elements by the key fn returns, preserving each group's encounter order",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "callback invoked as fn(v, i); its return value is the group key"),
		luareg.ArgType("fn", "fun(v: any, i: integer): any"),
		luareg.ReturnDoc(0, "groups", "a table mapping each distinct key to a new array of its elements"),
		luareg.ReturnType(0, "table<any, any[]>"))

	m.Fn("sort_by", sortByFn, "array mode: returns a new array stably sorted by the key fn returns (all-number or all-string keys only)",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "callback invoked as fn(v) (no index); its return value is the sort key, and must be a number or string, uniformly across all elements"),
		luareg.ArgType("fn", "fun(v: any): number|string"),
		luareg.ReturnDoc(0, "result", "a new array sorted ascending by key; equal keys keep their original relative order"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("partition", partitionFn, "array mode: splits t into a matching array and a rest array based on fn",
		luareg.Args("t", "fn"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ArgDoc("fn", "predicate invoked as fn(v, i)"),
		luareg.ArgType("fn", "fun(v: any, i: integer): boolean"),
		luareg.ReturnDoc(0, "matching", "a new array of elements for which fn was truthy"),
		luareg.ReturnDoc(1, "rest", "a new array of the remaining elements"),
		luareg.ReturnType(0, "any[]"),
		luareg.ReturnType(1, "any[]"))

	m.Fn("uniq", uniqFn, "array mode: returns a new array with duplicates removed, keeping the first occurrence",
		luareg.Args("t"),
		luareg.ArgDoc("t", "the table to iterate, array part only"),
		luareg.ReturnDoc(0, "result", "a new array; primitives dedup by value, tables/functions/userdata by identity"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("flatten", flattenFn, "array mode, recursive: inlines nested array tables up to depth levels (-1 = fully)",
		luareg.Args("t", "depth"),
		luareg.ArgDoc("t", "the table to flatten, array part only"),
		luareg.ArgDoc("depth", "how many levels of nested tables to inline; -1 flattens fully. Must be >= -1"),
		luareg.ReturnDoc(0, "result", "a new, flattened array"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("reverse", reverseFn, "array mode: returns a new array with t's elements in reverse order",
		luareg.Args("t"),
		luareg.ArgDoc("t", "the table to reverse, array part only"),
		luareg.ReturnDoc(0, "result", "a new array, elements in reverse order"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("zip", zipFn, "array mode on both tables: returns a new array of {a[i], b[i]} pairs, length min(#a, #b)",
		luareg.Args("a", "b"),
		luareg.ArgDoc("a", "the first table, array part only"),
		luareg.ArgDoc("b", "the second table, array part only"),
		luareg.ReturnDoc(0, "result", "a new array of 2-element arrays; extra elements in the longer input are dropped"),
		luareg.ReturnType(0, "any[][]"))

	m.Fn("chunk", chunkFn, "array mode: splits t into new arrays of at most size elements each; the last chunk may be shorter",
		luareg.Args("t", "size"),
		luareg.ArgDoc("t", "the table to chunk, array part only"),
		luareg.ArgDoc("size", "the maximum size of each chunk; must be >= 1"),
		luareg.ReturnDoc(0, "result", "a new array of arrays"),
		luareg.ReturnType(0, "any[][]"))

	m.Fn("keys", keysFn, "map mode: returns a new array of all of t's keys, in unspecified order",
		luareg.Args("t"),
		luareg.ArgDoc("t", "the table to read, array part and hash part together"),
		luareg.ReturnDoc(0, "keys", "a new array of t's keys; order is not defined"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("values", valuesFn, "map mode: returns a new array of all of t's values, in unspecified order",
		luareg.Args("t"),
		luareg.ArgDoc("t", "the table to read, array part and hash part together"),
		luareg.ReturnDoc(0, "values", "a new array of t's values; order is not defined"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("merge", mergeFn, "map mode, shallow: returns a new table with every input table's keys, later tables winning on conflict",
		luareg.Args("t", "..."),
		luareg.ArgDoc("t", "the base table; never modified"),
		luareg.ArgDoc("...", "additional tables applied left-to-right after t; never modified"),
		luareg.ReturnDoc(0, "result", "a new table; a key present in more than one input takes the last input's value"),
		luareg.ReturnType(0, "table<any, any>"))

	m.Fn("pick", pickFn, "map mode over names: returns a new table with only the named string-keyed fields of t",
		luareg.Args("t", "names"),
		luareg.ArgDoc("t", "the table to read from"),
		luareg.ArgDoc("names", "the string keys to keep; a name absent from t is silently skipped"),
		luareg.ReturnDoc(0, "result", "a new table containing only the requested keys"),
		luareg.ReturnType(0, "table<string, any>"))

	m.Fn("omit", omitFn, "map mode: returns a new table with every key of t except the named string keys",
		luareg.Args("t", "names"),
		luareg.ArgDoc("t", "the table to read from, array part and hash part together"),
		luareg.ArgDoc("names", "the string keys to exclude; non-string keys are never excluded"),
		luareg.ReturnDoc(0, "result", "a new table with the named keys removed"),
		luareg.ReturnType(0, "table<any, any>"))

	m.Fn("get_path", getPathFn, "map mode: walks t along path's dot-separated segments, returning default if any segment is missing or a non-table is indexed",
		luareg.Args("t", "path", "default"),
		luareg.ArgDoc("t", "the table to walk"),
		luareg.ArgDoc("path", "dot-separated segments, e.g. \"spec.containers.1.image\"; a segment that parses as an integer is tried as a number key first, then as a string key"),
		luareg.ArgDoc("default", "returned as-is when the path cannot be fully resolved; pass nil explicitly for no default"),
		luareg.ReturnDoc(0, "value", "the value found at path, or default"))

	m.Fn("deep_equal", deepEqualFn, "general: reports whether a and b are structurally equal (same keys, recursively equal values); handles cycles",
		luareg.Args("a", "b"),
		luareg.ArgDoc("a", "the first value"),
		luareg.ArgDoc("b", "the second value"),
		luareg.ReturnDoc(0, "equal", "true if a and b are structurally equal; metatables are not compared"))

	m.Fn("deep_copy", deepCopyFn, "general: returns a structural copy of v; nested tables are new tables, cycles are preserved",
		luareg.Args("v"),
		luareg.ArgDoc("v", "the value to copy"),
		luareg.ReturnDoc(0, "copy", "a deep copy of v; functions and userdata are copied by reference, metatables are not copied"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("collections", collections.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
