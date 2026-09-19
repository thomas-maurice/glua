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

// Package jsonpath queries nested Lua tables with JSONPath for Lua scripts,
// backed by k8s.io/client-go/util/jsonpath — the same engine `kubectl -o
// jsonpath=` uses. Zero new dependencies: client-go is already a direct
// requirement of this module (go.mod), so importing its jsonpath subpackage
// changes nothing in go.mod/go.sum.
//
// # Dialect
//
// This is the "kubectl jsonpath" dialect, NOT Goessner-canonical JSONPath and
// NOT RFC 9535. It has no length()/min()/max(), a narrower filter grammar
// (`?(@.x==y)` style comparisons only), and is as much a text-template
// engine (`{range}...{end}`, literal text between `{...}` segments) as a
// query language. Supported forms include `$`, `.a`, `['a']`, `[0]`,
// `[0:2]`, `..` recursive descent, `[*]` wildcard and `?(@.x==y)` filters.
// If a script works in `kubectl get -o jsonpath=`, it works here — that is
// the point of choosing this engine over a more "standard" one that would be
// a dependency for a dialect nobody in this project's audience actually
// writes. Do not expect RFC 9535 semantics.
//
// # Data path
//
// The input (`data`) is an arbitrary Lua value, normally a table. It is
// converted to Go via pkg/glua.Translator.FromLua before being handed to the
// jsonpath engine, which walks it by reflection — the engine expects
// map[string]interface{}/[]interface{}, which is exactly the Translator's
// JSON-round-trip shape. This means this module inherits the Translator's
// JSON-path caveats: all numbers arrive as float64 (a Lua integer above
// 2^53 will not round-trip exactly — this is a pre-existing glua limitation,
// not something introduced here), an empty Lua table reads back as an empty
// JSON object rather than an array, and array/map discrimination on the way
// in is LTable.MaxN()-based. Results are converted back to Lua via
// Translator.ToLua, so a matched JSON boolean/string/integral-valued number
// keeps its Lua type across the round trip (see the mixed-type test).
//
// # Cardinality
//
// query always returns an array-style table, even for zero or one match —
// consistency over convenience, so calling code never has to branch on
// "did I get a table or a bare value". A miss is an empty table, never nil
// (see the working-here skill's nil-slice rule). first exists specifically
// for the common "give me one value or a default" case.
//
// Because Lua cannot store a genuine nil as a trailing array element (a
// trailing nil is invisible to both `#t` and `ipairs`), a query whose LAST
// match is a JSON null will under-report its length by one. A null matched
// anywhere but the last position is stored correctly. Use jsonpath.exists to
// check presence independently of value when this matters.
//
// # Auto-brace
//
// The kubectl dialect's templates are `{...}`-wrapped. This module accepts
// bare paths too: if the given path/template does not already start with
// "{", it is wrapped automatically. So `.spec.replicas`, `$.spec.replicas`
// and `{.spec.replicas}` are all equivalent inputs to every function below.
//
// # Missing keys: query/first/exists vs render
//
// query, first and exists use AllowMissingKeys(true): a path segment that
// is not present in the data yields no results, not an error, because "is
// this field set?" is the whole point of exists (and a common use of
// query/first). render uses AllowMissingKeys(false): a text template that
// silently renders a gap is a bug in the template, so a missing key raises.
//
// # Arity
//
// exists and render take exactly two arguments; a wrong count always
// raises. query and first need *lua.LState to build their Lua return value
// (a fresh table for query, an arbitrary value for first), and
// pkg/luareg's reflection wrapper only enforces a MINIMUM argument count
// for a function that takes *lua.LState first — this is an existing,
// repo-wide property of the wrapper (every module function that returns a
// constructed table/value has it), not something specific to this module.
// A missing argument to query/first still raises; a surplus one is
// silently ignored.
//
// A malformed path/template (one that fails to parse) always raises,
// regardless of function, with a message naming the offending path. The
// underlying parser is not known to panic on malformed input — this was
// verified against a battery of malformed expressions (unterminated
// brackets, bad filters, out-of-range array indices, non-array indexing)
// during development — but every entry point still recovers from a panic
// defensively and turns it into a normal Go error, since it is a
// third-party parser this module does not control.
package jsonpath

import (
	"bytes"
	"fmt"
	"strings"

	glua "github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
	k8sjsonpath "k8s.io/client-go/util/jsonpath"
)

// translator: reusable translator for Lua <-> Go conversions.
var translator = glua.NewTranslator()

// normalisePath: wraps path in "{...}" if it is not already brace-wrapped,
// so bare paths (".spec.replicas"), $-rooted paths ("$.spec.replicas") and
// already-braced paths ("{.spec.replicas}") are all accepted identically.
func normalisePath(path string) string {
	if strings.HasPrefix(path, "{") {
		return path
	}
	return "{" + path + "}"
}

// toGoData: converts a Lua value to a Go value suitable for the jsonpath
// engine's reflection walk, riding pkg/glua.Translator's JSON data path.
// Translator.FromLua does not use its *lua.LState parameter (the conversion
// is pure data), so nil is passed deliberately rather than threading one
// through every caller.
func toGoData(data lua.LValue) (interface{}, error) {
	var out interface{}
	if err := translator.FromLua(nil, data, &out); err != nil {
		return nil, fmt.Errorf("jsonpath: failed to convert input data: %w", err)
	}
	return out, nil
}

// evalPath: parses path (auto-braced) and evaluates it against data,
// returning every matched value flattened into a single slice, in match
// order. allowMissingKeys controls whether an absent field is a miss (true)
// or a parse-time-adjacent error (false). fnName is used to prefix error
// messages. Recovers from any panic in the underlying parser/evaluator and
// converts it to a Go error, since that engine is a third-party dependency
// this module does not control.
func evalPath(fnName string, data interface{}, path string, allowMissingKeys bool) (results []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			results, err = nil, fmt.Errorf("%s: internal jsonpath error: %v", fnName, r)
		}
	}()

	jp := k8sjsonpath.New(fnName)
	jp.AllowMissingKeys(allowMissingKeys)
	if perr := jp.Parse(normalisePath(path)); perr != nil {
		return nil, fmt.Errorf("%s: invalid path %q: %w", fnName, path, perr)
	}
	groups, ferr := jp.FindResults(data)
	if ferr != nil {
		return nil, fmt.Errorf("%s: %w", fnName, ferr)
	}
	out := make([]interface{}, 0, len(groups))
	for _, group := range groups {
		for _, v := range group {
			out = append(out, v.Interface())
		}
	}
	return out, nil
}

// queryFn: evaluates path against data and returns every match as a new Lua
// array table, in match order. Always a table, even for zero or one match.
func queryFn(L *lua.LState, data lua.LValue, path string) (*lua.LTable, error) {
	goData, err := toGoData(data)
	if err != nil {
		return nil, err
	}
	matches, err := evalPath("jsonpath.query", goData, path, true)
	if err != nil {
		return nil, err
	}
	result := L.NewTable()
	for i, v := range matches {
		lv, err := translator.ToLua(L, v)
		if err != nil {
			return nil, fmt.Errorf("jsonpath.query: failed to convert result %d: %w", i, err)
		}
		// RawSetInt (not Append) so a JSON null match physically occupies
		// its slot instead of being silently dropped -- see the package doc
		// for the one case (a null as the LAST match) this cannot fix.
		result.RawSetInt(i+1, lv)
	}
	return result, nil
}

// firstFn: evaluates path against data and returns its first match, or def
// if there is no match. def is required (no optional-argument mechanism in
// this framework), which makes the miss case visible at the call site.
func firstFn(L *lua.LState, data lua.LValue, path string, def lua.LValue) (lua.LValue, error) {
	goData, err := toGoData(data)
	if err != nil {
		return nil, err
	}
	matches, err := evalPath("jsonpath.first", goData, path, true)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return def, nil
	}
	lv, err := translator.ToLua(L, matches[0])
	if err != nil {
		return nil, fmt.Errorf("jsonpath.first: failed to convert result: %w", err)
	}
	return lv, nil
}

// existsFn: reports whether path has at least one match in data. Takes no
// *lua.LState: the framework's arity check is exact only for functions that
// do not need the state escape hatch, and this one does not (both
// conversion steps here are pure data).
func existsFn(data lua.LValue, path string) (bool, error) {
	goData, err := toGoData(data)
	if err != nil {
		return false, err
	}
	matches, err := evalPath("jsonpath.exists", goData, path, true)
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

// renderFn: renders template against data using the jsonpath engine's
// native text-template mode (the same one `kubectl -o jsonpath=` uses),
// concatenating every matched value's text form and any literal text
// between `{...}` segments. Unlike query/first/exists, a missing key
// raises rather than silently producing a gap. Takes no *lua.LState, for
// the same exact-arity reason as existsFn.
func renderFn(data lua.LValue, template string) (result string, err error) {
	goData, cerr := toGoData(data)
	if cerr != nil {
		return "", cerr
	}

	defer func() {
		if r := recover(); r != nil {
			result, err = "", fmt.Errorf("jsonpath.render: internal jsonpath error: %v", r)
		}
	}()

	jp := k8sjsonpath.New("jsonpath.render")
	jp.AllowMissingKeys(false)
	if perr := jp.Parse(normalisePath(template)); perr != nil {
		return "", fmt.Errorf("jsonpath.render: invalid template %q: %w", template, perr)
	}
	var buf bytes.Buffer
	if eerr := jp.Execute(&buf, goData); eerr != nil {
		return "", fmt.Errorf("jsonpath.render: %w", eerr)
	}
	return buf.String(), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("jsonpath", "JSONPath queries over Lua tables using the kubectl jsonpath dialect -- see the package doc for cardinality, auto-brace and missing-key rules")

	m.Fn("query", queryFn, "evaluates path against data and returns every match as a new array table, in match order; empty when there is no match",
		luareg.Args("data", "path"),
		luareg.ArgDoc("data", "the value to query, normally a table"),
		luareg.ArgDoc("path", "a kubectl-jsonpath path; bare (\".spec.replicas\"), $-rooted (\"$.spec.replicas\") and braced (\"{.spec.replicas}\") forms are all accepted"),
		luareg.ReturnDoc(0, "matches", "an array of every matched value, in match order; always a table, never nil, even for zero or one match"),
		luareg.ReturnType(0, "any[]"))

	m.Fn("first", firstFn, "evaluates path against data and returns its first match, or default if there is no match",
		luareg.Args("data", "path", "default"),
		luareg.ArgDoc("data", "the value to query, normally a table"),
		luareg.ArgDoc("path", "a kubectl-jsonpath path; bare, $-rooted and braced forms are all accepted"),
		luareg.ArgDoc("default", "returned as-is when path has no match; pass nil explicitly for no default"),
		luareg.ReturnDoc(0, "value", "the first matched value, or default"))

	m.Fn("exists", existsFn, "reports whether path has at least one match in data",
		luareg.Args("data", "path"),
		luareg.ArgDoc("data", "the value to query, normally a table"),
		luareg.ArgDoc("path", "a kubectl-jsonpath path; bare, $-rooted and braced forms are all accepted"),
		luareg.ReturnDoc(0, "ok", "true if path matched at least one value in data"))

	m.Fn("render", renderFn, "renders template against data using the native kubectl jsonpath text-template mode, concatenating matched values and literal text",
		luareg.Args("data", "template"),
		luareg.ArgDoc("data", "the value to render from, normally a table"),
		luareg.ArgDoc("template", "a kubectl-jsonpath template, e.g. \"{.metadata.name}{'\\t'}{.status.phase}\"; bare and braced forms are both accepted"),
		luareg.ReturnDoc(0, "text", "the rendered text"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("jsonpath", jsonpath.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
