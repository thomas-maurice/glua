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

	lua "github.com/yuin/gopher-lua"
)

// maxTableDepth: upper bound on *lua.LTable nesting that TableGuard will
// allow before giving up with an error, instead of letting a hand-written
// recursive walk continue until the Go runtime kills the process with a
// fatal, unrecoverable stack overflow (see TableGuard's doc comment for why
// that matters more than an ordinary panic).
//
// Calibrated by measuring a deliberately elaborate corev1.Pod — multi-container
// spec with node/pod affinity, a projected volume, a lifecycle hook and an
// exec probe — which JSON-marshals to a nesting depth of 10. 100 leaves 10x
// headroom over that and comfortably covers deeper CRD shapes in the wild
// (Argo Workflow templates, Istio VirtualServices commonly nest 15-20
// levels), while remaining many orders of magnitude below the recursion
// depth actually required to exhaust a Go goroutine's stack (the reported
// crash needed a 1,000,000,000-byte stack; each recursive frame is on the
// order of a few hundred bytes, so overflow requires on the order of
// millions of nested levels, not hundreds).
//
// This is a single, package-level constant deliberately shared by every
// caller of TableGuard, so the whole library enforces one consistent limit
// rather than each table walker picking its own.
const maxTableDepth = 100

// TableGuard: shared recursion guard against a cyclic or pathologically deep
// Lua table, for use by anything that hand-walks a *lua.LTable rather than
// going through Translator.FromLua (which already uses this internally).
//
// A script doing `t = {}; t.self = t` and handing t to a table-walking
// function — json.stringify, yaml.stringify, template.render, spew.dump, a
// module function taking a struct/map, ... — makes an unguarded recursive
// walk recurse forever. That crashes the ENTIRE process with "fatal error:
// stack overflow": not an ordinary Go panic, but a fatal runtime error that
// neither gopher-lua's pcall nor Go's recover() can catch. In this
// project's stated deployments (admission webhooks, policy engines running
// semi-trusted Lua), one line of script permanently kills a long-running
// server.
//
// Any new Lua-to-Go table walker MUST thread a *TableGuard through its
// recursion and call Enter at every *lua.LTable it visits, instead of
// hand-rolling its own recursion with no bound. Copying the old (unguarded)
// pattern from before this type existed reintroduces the crash.
type TableGuard struct {
	// ancestors holds the *lua.LTable pointers currently open on the CURRENT
	// recursion path only — added by Enter, removed by the returned leave
	// func. This is deliberately NOT a global "already visited" set: a
	// global set would reject a legitimate diamond reference (the same
	// table reachable twice via different fields, e.g. `t.a = shared; t.b =
	// shared`), which is not a cycle and which encoding/json marshals
	// without complaint — rejecting it here would be a regression for
	// callers relying on that shape. Path-scoped tracking rejects true
	// self-reference while still allowing a table to be visited multiple
	// times as long as it never appears as its own ancestor.
	ancestors map[*lua.LTable]struct{}
	depth     int
}

// NewTableGuard: creates a guard ready to walk a fresh table graph from the
// root (depth 0, no ancestors recorded yet). Create one per top-level
// conversion call (e.g. once per json.stringify invocation) and thread the
// same instance through the whole recursive walk — do not create a new one
// per recursive call, or ancestor tracking and depth counting will not work.
func NewTableGuard() *TableGuard {
	return &TableGuard{ancestors: make(map[*lua.LTable]struct{})}
}

// Enter: call before recursing into tbl's children. Returns an error naming
// whether tbl is a cycle (already open as an ancestor on this path) or
// nesting has passed maxTableDepth, in which case the caller must not
// recurse into tbl at all and must not call Leave.
//
// On success, the caller MUST call Leave(tbl) (typically via
// `defer guard.Leave(tbl)`) once it is done processing tbl's children and is
// about to return, to pop tbl back off the ancestor path. Enter/Leave are a
// pair rather than Enter returning a closure specifically so this stays
// allocation-free on this hot conversion path: a closure returned from Enter
// would have to capture tbl and escape to the heap on every table node
// visited, which showed up as a measurable regression in benchmarking.
func (g *TableGuard) Enter(tbl *lua.LTable) error {
	if _, isAncestor := g.ancestors[tbl]; isAncestor {
		return fmt.Errorf("cannot convert Lua table: self-referential table detected (cycle) at nesting depth %d", g.depth)
	}
	if g.depth >= maxTableDepth {
		return fmt.Errorf("cannot convert Lua table: exceeded maximum nesting depth of %d; the table is either a cycle or is nested implausibly deep for any real payload", maxTableDepth)
	}

	g.ancestors[tbl] = struct{}{}
	g.depth++
	return nil
}

// Leave: pops tbl back off the ancestor path after a successful Enter(tbl).
// Must be called exactly once per successful Enter, typically via
// `defer guard.Leave(tbl)` right after checking Enter's error.
func (g *TableGuard) Leave(tbl *lua.LTable) {
	g.depth--
	delete(g.ancestors, tbl)
}
