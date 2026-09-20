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

package stubgen

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomas-maurice/glua/pkg/luareg"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
)

// --- Small helper types used throughout the tests ---

// testPoint: a small struct used to test struct param/return ---@class generation.
type testPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// testRect: a small struct with a nested struct field.
type testRect struct {
	Origin testPoint `json:"origin"`
	Width  float64   `json:"width"`
	Height float64   `json:"height"`
}

// docFixture: exercises every source of a ---@field description the
// resolver in comments.go supports. It must be a package-level declaration
// (not a type declared inside a test function) because loadPackageFieldDocs
// only walks top-level GenDecls in the package's AST.
type docFixture struct {
	// Name is the fixture's display name, read from a doc comment above the
	// field.
	Name string `json:"name"`

	Count int `json:"count"` // Count is read from this inline comment.

	// Overridden has its own doc comment, but the "luadoc" tag takes
	// precedence over it.
	Overridden string `json:"overridden" luadoc:"tag wins over doc comment"`

	Silent string `json:"silent"`
}

// embeddingFixture: exercises extractFieldDocs' embedded-field handling —
// the embedded field's doc comment must be keyed under its declared
// (promoted) name, "docFixture". Package-level for the same AST-walking
// reason as docFixture.
type embeddingFixture struct {
	// Base is embedded directly in embeddingFixture.
	docFixture
}

// --- Pure Go functions registered in tests ---

func addFn(a, b int) int                      { return a + b }
func greetFn(name string) string              { return "hello " + name }
func boolFn(v bool) bool                      { return v }
func errFn(x int) (int, error)                { return x, nil }
func lstateFn(L *lua.LState, s string) string { return s }
func pointFn(p testPoint) testPoint           { return p }
func sliceFn(parts []string) string           { return strings.Join(parts, ",") }
func mapFn(m map[string]string) string        { return "" }
func multiRetFn(a, b string) (string, string) { return a, b }
func ltableReturnFn() *lua.LTable             { return nil }
func lstateReturnFn() *lua.LState             { return nil }

// TestNewGenerator: ensures NewGenerator returns a valid, non-nil generator.
func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	require.NotNil(t, gen)
	require.NotNil(t, gen.typeRegistry)
}

// TestGenerateModule_Empty: an empty module produces a valid, minimal file.
func TestGenerateModule_Empty(t *testing.T) {
	m := luareg.NewModule("empty", "empty module")
	gen := NewGenerator()

	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.True(t, strings.HasPrefix(out, "---@meta empty\n"), "must start with ---@meta empty")
	assert.Contains(t, out, "---@class empty")
	assert.Contains(t, out, "local empty = {}")
	assert.Contains(t, out, "return empty")
}

// TestGenerateModule_FunctionsOnly: a module with only plain functions.
// Verifies param annotations, return annotations, function signature, ordering.
func TestGenerateModule_FunctionsOnly(t *testing.T) {
	m := luareg.NewModule("math", "math utilities")
	m.Fn(
		"add", addFn, "adds two numbers",
		luareg.Args("a", "b"),
		luareg.ArgDoc("a", "first operand"),
		luareg.ArgDoc("b", "second operand"),
		luareg.ReturnDoc(0, "result", "the sum"),
	)
	m.Fn(
		"greet", greetFn, "greets someone",
		luareg.Args("name"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// Header.
	assert.True(t, strings.HasPrefix(out, "---@meta math\n"))

	// Module class declaration.
	assert.Contains(t, out, "---@class math\n")
	assert.Contains(t, out, "local math = {}")

	// add function annotations.
	assert.Contains(t, out, "---@param a number first operand")
	assert.Contains(t, out, "---@param b number second operand")
	assert.Contains(t, out, "---@return number result the sum")
	assert.Contains(t, out, "function math.add(a, b) end")

	// greet function annotations.
	assert.Contains(t, out, "---@param name string")
	assert.Contains(t, out, "---@return string")
	assert.Contains(t, out, "function math.greet(name) end")

	// Return statement.
	assert.Contains(t, out, "return math")

	// Functions appear in registration order (add before greet).
	addPos := strings.Index(out, "function math.add")
	greetPos := strings.Index(out, "function math.greet")
	assert.True(t, addPos < greetPos, "add must appear before greet (registration order)")
}

// TestGenerateModule_BoolReturn: a function returning bool produces "boolean" type.
func TestGenerateModule_BoolReturn(t *testing.T) {
	m := luareg.NewModule("util", "util module")
	m.Fn(
		"identity", boolFn, "identity bool",
		luareg.Args("v"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param v boolean")
	assert.Contains(t, out, "---@return boolean")
}

// TestGenerateModule_LStateSkipped: a function whose first param is *lua.LState
// has that param omitted from the Lua signature.
func TestGenerateModule_LStateSkipped(t *testing.T) {
	m := luareg.NewModule("escape", "escape hatch module")
	m.Fn(
		"passthrough", lstateFn, "passes a string through",
		luareg.Args("s"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// "s" should appear as the sole param.
	assert.Contains(t, out, "---@param s string")
	assert.Contains(t, out, "function escape.passthrough(s) end")
	// Ensure *lua.LState did not leak into the param list.
	assert.NotContains(t, out, "LState")
}

// TestGenerateModule_ErrorReturnOmitted: a function returning (T, error)
// must show T in the return annotation but must not emit anything for the
// error return.
func TestGenerateModule_ErrorReturnOmitted(t *testing.T) {
	m := luareg.NewModule("safe", "safe module")
	m.Fn(
		"safe_add", errFn, "safe add",
		luareg.Args("x"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// number return should be present.
	assert.Contains(t, out, "---@return number")
	// error must NOT appear as a return annotation.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "---@return") {
			assert.NotContains(t, line, "error", "error return must be omitted from annotations")
		}
	}
}

// TestGenerateModule_SliceParam: []string becomes string[].
func TestGenerateModule_SliceParam(t *testing.T) {
	m := luareg.NewModule("arr", "array module")
	m.Fn(
		"join", sliceFn, "joins strings",
		luareg.Args("parts"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param parts string[]")
}

// TestGenerateModule_MapParam: map[string]string becomes table<string, string>.
func TestGenerateModule_MapParam(t *testing.T) {
	m := luareg.NewModule("maps", "map module")
	m.Fn(
		"first", mapFn, "gets first",
		luareg.Args("m"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param m table<string, string>")
}

// TestGenerateModule_MultiReturn: multiple non-error returns all appear.
func TestGenerateModule_MultiReturn(t *testing.T) {
	m := luareg.NewModule("duo", "duo module")
	m.Fn(
		"swap", multiRetFn, "swaps two strings",
		luareg.Args("a", "b"),
		luareg.ReturnDoc(0, "first", "first return"),
		luareg.ReturnDoc(1, "second", "second return"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return string first first return")
	assert.Contains(t, out, "---@return string second second return")
}

// TestGenerateModule_StructParamGeneratesClass: a struct parameter causes a
// ---@class block to appear at the top of the file.
func TestGenerateModule_StructParamGeneratesClass(t *testing.T) {
	m := luareg.NewModule("geo", "geometry module")
	m.Fn(
		"mirror", pointFn, "mirrors a point",
		luareg.Args("p"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// A ---@class block for testPoint must appear.
	assert.Contains(t, out, "---@class stubgen.testPoint", "struct class block must be generated")
	// The param type must reference the class.
	assert.Contains(t, out, "---@param p stubgen.testPoint")
	// The class block must precede the module declaration.
	classPos := strings.Index(out, "---@class stubgen.testPoint")
	moduleClassPos := strings.Index(out, "---@class geo\n")
	assert.True(
		t,
		classPos < moduleClassPos,
		"struct @class block must precede the module @class block",
	)
}

// TestGenerateModule_OneClass: a module with a single registered class produces
// the class block, method stubs, the ---@field on the module, and the assignment line.
func TestGenerateModule_OneClass(t *testing.T) {
	// Build a simple Logger class with two methods.
	type Logger struct{}
	logMethod := func(l *Logger, msg string) {}
	withMethod := func(l *Logger, key, val string) *Logger { return l }

	cls := luareg.NewClass[*Logger]("log.Logger", "structured logger")
	cls.Method(
		"info", logMethod, "logs info",
		luareg.Args("msg"),
		luareg.ArgDoc("msg", "the message"),
	)
	cls.Method(
		"with", withMethod, "returns child logger",
		luareg.Args("key", "val"),
		luareg.ReturnDoc(0, "logger", "child logger"),
	)

	m := luareg.NewModule("log", "logging module")
	m.RegisterClass(cls)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// Class block.
	assert.Contains(t, out, "---@class log.Logger\n")
	assert.Contains(t, out, "local Logger = {}")

	// Method stubs.
	assert.Contains(t, out, "---@param msg string the message")
	assert.Contains(t, out, "function Logger:info(msg) end")
	assert.Contains(t, out, "function Logger:with(key, val) end")

	// Return type for "with" uses the class name.
	assert.Contains(t, out, "---@return log.Logger logger child logger")

	// Module declaration has the class as a field.
	assert.Contains(t, out, "---@class log\n")
	assert.Contains(t, out, "---@field Logger log.Logger")
	assert.Contains(t, out, "local log = {}")

	// Assignment line at the end of the class section.
	assert.Contains(t, out, "log.Logger = Logger")

	// Return statement.
	assert.Contains(t, out, "return log")

	// Class block appears before the module class declaration.
	classPos := strings.Index(out, "---@class log.Logger")
	modClassPos := strings.Index(out, "---@class log\n")
	assert.True(t, classPos < modClassPos)
}

// TestGenerateModule_MultipleClasses: multiple classes are emitted in
// registration order, each with their own block and assignment.
func TestGenerateModule_MultipleClasses(t *testing.T) {
	type Alpha struct{}
	type Beta struct{}

	alphaMethod := func(a *Alpha) string { return "" }
	betaMethod := func(b *Beta) int { return 0 }

	clsA := luareg.NewClass[*Alpha]("mod.Alpha", "alpha class")
	clsA.Method("name", alphaMethod, "returns name")

	clsB := luareg.NewClass[*Beta]("mod.Beta", "beta class")
	clsB.Method("count", betaMethod, "returns count")

	m := luareg.NewModule("mod", "multi-class module")
	m.RegisterClass(clsA)
	m.RegisterClass(clsB)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// Both class blocks present.
	assert.Contains(t, out, "---@class mod.Alpha")
	assert.Contains(t, out, "---@class mod.Beta")

	// Both field annotations on module.
	assert.Contains(t, out, "---@field Alpha mod.Alpha")
	assert.Contains(t, out, "---@field Beta mod.Beta")

	// Both assignment lines.
	assert.Contains(t, out, "mod.Alpha = Alpha")
	assert.Contains(t, out, "mod.Beta = Beta")

	// Registration order preserved: Alpha before Beta.
	alphaPos := strings.Index(out, "---@class mod.Alpha")
	betaPos := strings.Index(out, "---@class mod.Beta")
	assert.True(t, alphaPos < betaPos, "Alpha must appear before Beta (registration order)")
}

// TestGenerateModule_ClassReturnsDifferentClass: a method that returns a
// different registered class uses that class's full Lua name.
func TestGenerateModule_ClassReturnsDifferentClass(t *testing.T) {
	type Builder struct{}
	type Result struct{}

	buildMethod := func(b *Builder) *Result { return &Result{} }

	clsResult := luareg.NewClass[*Result]("factory.Result", "result class")
	clsBuilder := luareg.NewClass[*Builder]("factory.Builder", "builder class")
	clsBuilder.Method("build", buildMethod, "builds a Result")

	m := luareg.NewModule("factory", "factory module")
	m.RegisterClass(clsResult)
	m.RegisterClass(clsBuilder)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// build() should return factory.Result, not factory.Builder.
	assert.Contains(t, out, "---@return factory.Result")
}

// TestGenerateModule_MethodSkipsReceiver: the receiver (first param of a
// method expression) must not appear in the Lua param list.
func TestGenerateModule_MethodSkipsReceiver(t *testing.T) {
	type Counter struct{ n int }
	incMethod := func(c *Counter, delta int) int { return c.n + delta }

	cls := luareg.NewClass[*Counter]("cnt.Counter", "a counter")
	cls.Method(
		"inc", incMethod, "increments counter",
		luareg.Args("delta"),
	)

	m := luareg.NewModule("cnt", "counter module")
	m.RegisterClass(cls)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// Only "delta" should appear as a param, not the receiver.
	assert.Contains(t, out, "---@param delta number")
	assert.Contains(t, out, "function Counter:inc(delta) end")
	// Ensure Counter pointer type is not in the param list.
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "---@param") {
			assert.NotContains(t, line, "Counter", "receiver must not appear as a param")
		}
	}
}

// TestRegisterType_GenerateTypeStubs: RegisterType queues types for the
// annotations.gen.lua file, separate from any module.
func TestRegisterType_GenerateTypeStubs(t *testing.T) {
	type Config struct {
		Name    string  `json:"name"`
		Timeout float64 `json:"timeout"`
	}

	gen := NewGenerator()
	require.NoError(t, gen.RegisterType(Config{}))

	stubs, err := gen.GenerateTypeStubs()
	require.NoError(t, err)
	require.NotEmpty(t, stubs)

	assert.Contains(t, stubs, "---@meta annotations")
	assert.Contains(t, stubs, "---@class stubgen.Config")
	assert.Contains(t, stubs, "---@field name string")
	assert.Contains(t, stubs, "---@field timeout number")
}

// TestRegisterType_LuadocTagDescribesField: a struct field tagged with
// `luadoc:"..."` gets that text appended to its ---@field annotation, since
// reflection cannot see Go doc comments. Untagged fields render unchanged.
func TestRegisterType_LuadocTagDescribesField(t *testing.T) {
	type Message struct {
		Body   string `json:"body"   luadoc:"plaintext message body"`
		Sender string `json:"sender"`
	}

	gen := NewGenerator()
	require.NoError(t, gen.RegisterType(Message{}))

	stubs, err := gen.GenerateTypeStubs()
	require.NoError(t, err)

	assert.Contains(t, stubs, "---@field body string plaintext message body")
	assert.Contains(t, stubs, "---@field sender string\n")
}

// TestRegisterType_FieldDocFromGoComments: field descriptions come from Go
// source comments (resolved via golang.org/x/tools/go/packages), not just
// the "luadoc" tag. Exercises all four precedence cases on docFixture: a
// leading doc comment, a trailing inline comment, a luadoc tag override that
// wins over its own doc comment, and a field with no description at all.
func TestRegisterType_FieldDocFromGoComments(t *testing.T) {
	gen := NewGenerator()
	require.NoError(t, gen.RegisterType(docFixture{}))

	stubs, err := gen.GenerateTypeStubs()
	require.NoError(t, err)

	// Leading doc comment.
	assert.Contains(
		t,
		stubs,
		"---@field name string Name is the fixture's display name, read from a doc comment above the field.",
	)
	// Trailing inline comment.
	assert.Contains(t, stubs, "---@field count number Count is read from this inline comment.")
	// luadoc tag wins over the field's own doc comment.
	assert.Contains(t, stubs, "---@field overridden string tag wins over doc comment")
	assert.NotContains(t, stubs, "own doc comment, but the")
	// No comment, no tag -> no description.
	assert.Contains(t, stubs, "---@field silent string\n")
}

// TestGenerateModule_KubernetesFieldDocsFromSource: the kubernetes module
// registers real k8s.io/api struct shapes (corev1.Pod and friends) via
// RegisterStubType, so its generated stub is the real-world exercise of the
// Go-doc-comment resolver against a large, module-cache-resolved dependency
// graph — not just our own small package. Must complete without error and
// without hanging; a well-known k8s doc comment must render. If package
// resolution ever becomes too slow/flaky here, comment extraction should be
// gated to a package allowlist rather than reverting this test.
func TestGenerateModule_KubernetesFieldDocsFromSource(t *testing.T) {
	reg := luareg.NewRegistry()
	kubernetes.Register(reg)
	require.Len(t, reg.Modules(), 1)

	gen := NewGenerator()
	start := time.Now()
	out, err := gen.GenerateModule(reg.Modules()[0])
	elapsed := time.Since(start)
	require.NoError(t, err)

	assert.Less(t, elapsed, 30*time.Second, "kubernetes module stub generation must not hang")

	// corev1.ObjectMeta.Name carries a real, stable Go doc comment; its
	// presence proves the resolver reached into the k8s.io/api module cache.
	assert.Contains(t, out, "---@field name string Name must be unique within a namespace.")
}

// TestRegisterType_Empty: when no types are registered, GenerateTypeStubs
// returns an empty string (no file should be written).
func TestRegisterType_Empty(t *testing.T) {
	gen := NewGenerator()

	stubs, err := gen.GenerateTypeStubs()
	require.NoError(t, err)
	assert.Empty(t, stubs, "no types registered → empty output, no file to write")
}

// TestGenerateFromRegistry_Empty: an empty registry produces no files and no error.
func TestGenerateFromRegistry_Empty(t *testing.T) {
	reg := luareg.NewRegistry()
	gen := NewGenerator()
	tmpDir := t.TempDir()

	files, err := gen.GenerateFromRegistry(reg, tmpDir)
	require.NoError(t, err)
	assert.Empty(t, files)
}

// TestGenerateFromRegistry_MultipleModules: each module produces one file;
// annotations.gen.lua is not written when no types are registered via RegisterType.
func TestGenerateFromRegistry_MultipleModules(t *testing.T) {
	reg := luareg.NewRegistry()

	mA := luareg.NewModule("alpha", "alpha module")
	mA.Fn("double", func(x int) int { return x * 2 }, "doubles x", luareg.Args("x"))
	mA.Register(reg)

	mB := luareg.NewModule("beta", "beta module")
	mB.Fn("triple", func(x int) int { return x * 3 }, "triples x", luareg.Args("x"))
	mB.Register(reg)

	gen := NewGenerator()
	tmpDir := t.TempDir()

	files, err := gen.GenerateFromRegistry(reg, tmpDir)
	require.NoError(t, err)
	require.Len(t, files, 2, "two modules → two files")

	// Each file name matches the module.
	fileNames := make(map[string]bool)
	for _, f := range files {
		fileNames[filepath.Base(f)] = true
	}
	assert.True(t, fileNames["alpha.gen.lua"])
	assert.True(t, fileNames["beta.gen.lua"])

	// No annotations.gen.lua since no types were registered.
	assert.False(t, fileNames["annotations.gen.lua"])

	// Files must exist and be non-empty.
	for _, f := range files {
		info, err := os.Stat(f)
		require.NoError(t, err)
		assert.Greater(t, info.Size(), int64(0))
	}
}

// TestGenerateFromRegistry_WithAnnotations: when types are registered via
// RegisterType, annotations.gen.lua is also written.
func TestGenerateFromRegistry_WithAnnotations(t *testing.T) {
	type Tag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	reg := luareg.NewRegistry()
	mA := luareg.NewModule("tags", "tags module")
	mA.Register(reg)

	gen := NewGenerator()
	require.NoError(t, gen.RegisterType(Tag{}))

	tmpDir := t.TempDir()
	files, err := gen.GenerateFromRegistry(reg, tmpDir)
	require.NoError(t, err)

	fileNames := make(map[string]bool)
	for _, f := range files {
		fileNames[filepath.Base(f)] = true
	}
	assert.True(
		t,
		fileNames["annotations.gen.lua"],
		"annotations.gen.lua must be written when types are registered",
	)

	// Verify its content.
	content, err := os.ReadFile(filepath.Join(tmpDir, "annotations.gen.lua"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "---@class stubgen.Tag")
}

// TestGenerateFromRegistry_Idempotent: calling GenerateFromRegistry twice with
// the same inputs produces byte-identical files.
func TestGenerateFromRegistry_Idempotent(t *testing.T) {
	reg := luareg.NewRegistry()
	mA := luareg.NewModule("idempotent", "idempotent module")
	mA.Fn("work", greetFn, "greets someone", luareg.Args("name"))
	mA.Register(reg)

	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	gen1 := NewGenerator()
	files1, err := gen1.GenerateFromRegistry(reg, tmpDir1)
	require.NoError(t, err)

	gen2 := NewGenerator()
	files2, err := gen2.GenerateFromRegistry(reg, tmpDir2)
	require.NoError(t, err)

	require.Equal(t, len(files1), len(files2))

	for i := range files1 {
		c1, err := os.ReadFile(files1[i])
		require.NoError(t, err)
		c2, err := os.ReadFile(files2[i])
		require.NoError(t, err)
		assert.Equal(
			t,
			string(c1),
			string(c2),
			"file %s must be byte-identical on repeated generation",
			filepath.Base(files1[i]),
		)
	}
}

// TestGenerateModule_ArgDocWithoutName: if ArgDocs has an entry whose name
// doesn't match any arg, it is simply not emitted (no crash).
func TestGenerateModule_ArgDocWithoutName(t *testing.T) {
	m := luareg.NewModule("util", "util module")
	m.Fn(
		"noop", func(x string) string { return x }, "identity",
		luareg.Args("x"),
		luareg.ArgDoc("y", "nonexistent param"), // "y" not in ArgNames
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// x gets no doc because "x" has no matching ArgDoc entry.
	assert.Contains(t, out, "---@param x string\n")
}

// TestGenerateModule_ArgTypeOverride: a param with luareg.ArgType renders the
// override annotation instead of the reflected type. This is the escape hatch
// for params like lua.LValue callbacks, whose Go type reflects to "any" and
// gives LuaLS nothing useful to check against.
func TestGenerateModule_ArgTypeOverride(t *testing.T) {
	m := luareg.NewModule("bot", "bot module")
	m.Fn(
		"on_message", func(handler lua.LValue) {}, "register a handler for room message events",
		luareg.Args("handler"),
		luareg.ArgType("handler", "fun(evt: core.Event)"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param handler fun(evt: core.Event)\n")
	assert.NotContains(t, out, "---@param handler any")
}

// TestGenerateModule_ArgTypeOverrideWithDoc: ArgType and ArgDoc compose — the
// override type is used alongside the doc string.
func TestGenerateModule_ArgTypeOverrideWithDoc(t *testing.T) {
	m := luareg.NewModule("bot", "bot module")
	m.Fn(
		"on_message", func(handler lua.LValue) {}, "register a handler",
		luareg.Args("handler"),
		luareg.ArgType("handler", "fun(evt: core.Event)"),
		luareg.ArgDoc("handler", "called for each message event"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(
		t,
		out,
		"---@param handler fun(evt: core.Event) called for each message event\n",
	)
}

// TestGenerateModule_ArgTypeFallback: a param without an ArgType override
// still renders its reflected Go type, unchanged from before this feature.
func TestGenerateModule_ArgTypeFallback(t *testing.T) {
	m := luareg.NewModule("math", "math utilities")
	m.Fn("add", addFn, "adds two numbers", luareg.Args("a", "b"))

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param a number\n")
	assert.Contains(t, out, "---@param b number\n")
}

// TestGenerateModule_ArgTypeDuplicate: duplicate ArgType calls for the same
// name are allowed (mirrors ArgDoc); the first registered entry wins because
// findArgType returns on first match, same as findArgDoc.
func TestGenerateModule_ArgTypeDuplicate(t *testing.T) {
	m := luareg.NewModule("bot", "bot module")
	m.Fn(
		"on_message", func(handler lua.LValue) {}, "register a handler",
		luareg.Args("handler"),
		luareg.ArgType("handler", "fun(evt: core.Event)"),
		luareg.ArgType("handler", "fun(evt: core.OtherEvent)"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param handler fun(evt: core.Event)\n")
	assert.NotContains(t, out, "core.OtherEvent")
}

// TestGenerateModule_ArgTypeOnMethod: ArgType also applies to class methods,
// since Method shares the same FnMeta/FnOpt machinery as Fn.
func TestGenerateModule_ArgTypeOnMethod(t *testing.T) {
	type Bot struct{}
	onMessage := func(b *Bot, handler lua.LValue) {}

	cls := luareg.NewClass[*Bot]("bot.Bot", "bot instance")
	cls.Method(
		"on_message", onMessage, "register a handler",
		luareg.Args("handler"),
		luareg.ArgType("handler", "fun(evt: core.Event)"),
	)

	m := luareg.NewModule("bot", "bot module")
	m.RegisterClass(cls)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@param handler fun(evt: core.Event)\n")
}

// TestGenerateModule_FunctionWithDoc: function-level Doc string is emitted as
// a leading --- comment before the annotations.
func TestGenerateModule_FunctionWithDoc(t *testing.T) {
	m := luareg.NewModule("doc", "doc module")
	m.Fn(
		"greet", greetFn, "greets someone by name",
		luareg.Args("name"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "--- greets someone by name\n")
}

// TestGenerateModule_StructReturnGeneratesClass: a struct return type also
// causes a ---@class block to be generated, same as a struct param.
func TestGenerateModule_StructReturnGeneratesClass(t *testing.T) {
	makePoint := func(x, y float64) testPoint { return testPoint{x, y} }

	m := luareg.NewModule("geo2", "geo2 module")
	m.Fn(
		"make_point", makePoint, "creates a point",
		luareg.Args("x", "y"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@class stubgen.testPoint")
	assert.Contains(t, out, "---@return stubgen.testPoint")
}

// TestGenerateFromRegistry_OutputDirCreated: GenerateFromRegistry creates the
// output directory if it does not exist.
func TestGenerateFromRegistry_OutputDirCreated(t *testing.T) {
	reg := luareg.NewRegistry()
	mA := luareg.NewModule("mkdir", "mkdir module")
	mA.Register(reg)

	tmpDir := t.TempDir()
	nested := filepath.Join(tmpDir, "deep", "nested", "dir")

	gen := NewGenerator()
	_, err := gen.GenerateFromRegistry(reg, nested)
	require.NoError(t, err)

	_, statErr := os.Stat(nested)
	assert.NoError(t, statErr, "output directory must be created")
}

// TestGoTypeToLua_Primitives: exercise the type mapping for all primitive kinds.
func TestGoTypeToLua_Primitives(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)

	cases := []struct {
		in  any
		out string
	}{
		{"", "string"},
		{false, "boolean"},
		{0, "number"},
		{int8(0), "number"},
		{int16(0), "number"},
		{int32(0), "number"},
		{int64(0), "number"},
		{uint(0), "number"},
		{uint8(0), "number"},
		{uint16(0), "number"},
		{uint32(0), "number"},
		{uint64(0), "number"},
		{float32(0), "number"},
		{float64(0), "number"},
		{[]string{}, "string[]"},
	}

	for _, tc := range cases {
		t.Run(tc.out, func(t *testing.T) {
			got := goTypeToLua(reflect.TypeOf(tc.in), classLookup)
			assert.Equal(t, tc.out, got)
		})
	}
}

// TestGoTypeToLua_ByteSlice: []byte must emit "string", not "number[]" —
// pkg/luareg maps []byte to a raw Lua string at the top level (see
// pkg/luareg's luaToGo/goToLua []byte special case), so a stub claiming
// number[] would mislead the LSP and callers.
func TestGoTypeToLua_ByteSlice(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)
	assert.Equal(t, "string", goTypeToLua(reflect.TypeOf([]byte(nil)), classLookup))
}

// TestGoTypeToLua_Compound: maps, interfaces, and pointer types.
func TestGoTypeToLua_Compound(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)

	// map[string]int → table<string, number>
	mapType := reflect.TypeOf(map[string]int{})
	assert.Equal(t, "table<string, number>", goTypeToLua(mapType, classLookup))

	// *testPoint → stubgen.testPoint (struct pointer, non-class)
	ptrType := reflect.TypeOf((*testPoint)(nil))
	assert.Equal(t, "stubgen.testPoint", goTypeToLua(ptrType, classLookup))

	// struct value → stubgen.testPoint
	structType := reflect.TypeOf(testPoint{})
	assert.Equal(t, "stubgen.testPoint", goTypeToLua(structType, classLookup))

	// []testPoint → stubgen.testPoint[]
	sliceStructType := reflect.TypeOf([]testPoint{})
	assert.Equal(t, "stubgen.testPoint[]", goTypeToLua(sliceStructType, classLookup))
}

// TestGoTypeToLua_ClassPointerLookup: when a *T type is itself in the class
// lookup, goTypeToLua returns the class name directly without dereferencing.
func TestGoTypeToLua_ClassPointerLookup(t *testing.T) {
	type Widget struct{}
	cls := luareg.NewClass[*Widget]("ui.Widget", "widget class")

	classLookup := map[reflect.Type]luareg.AnyClass{
		cls.GoType(): cls,
	}

	// *Widget is the class's GoType — direct hit.
	assert.Equal(t, "ui.Widget", goTypeToLua(cls.GoType(), classLookup))
}

// TestClassLocalName_WithoutPrefix: a class name that doesn't start with
// "modname." is returned unchanged.
func TestClassLocalName_WithoutPrefix(t *testing.T) {
	assert.Equal(t, "MyClass", classLocalName("MyClass", "mymod"))
	assert.Equal(t, "other.Class", classLocalName("other.Class", "mymod"))
}

// TestGenerateModule_PtrStructParam: a *struct parameter generates a
// ---@class block and the param type references it correctly.
func TestGenerateModule_PtrStructParam(t *testing.T) {
	withPointer := func(p *testPoint) *testPoint { return p }

	m := luareg.NewModule("geo3", "geo3 module")
	m.Fn(
		"copy", withPointer, "copies a point",
		luareg.Args("p"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@class stubgen.testPoint")
	assert.Contains(t, out, "---@param p stubgen.testPoint")
	assert.Contains(t, out, "---@return stubgen.testPoint")
}

// TestGenerateModule_MethodNoGoFn: a FnMeta with nil GoFn (pathological case)
// must not panic and produces a stub with empty params.
func TestGenerateModule_MethodNoGoFn(t *testing.T) {
	// We can't register a nil GoFn via the normal API (it panics), but we can
	// verify writeFuncStub handles the nil GoFn guard by exercising a module
	// with a real function, then confirm nil is handled by calling the internal
	// helpers directly.
	fn := &luareg.FnMeta{LuaName: "noop", GoFn: nil}
	var sb strings.Builder
	classLookup := make(map[reflect.Type]luareg.AnyClass)

	// Should not panic.
	writeFuncStub(&sb, fn, "mymod", classLookup)
	out := sb.String()
	assert.Contains(t, out, "function mymod.noop() end")
}

// TestGenerateModule_ReturnDocByIndex: ReturnDoc index ordering must be
// respected even when the order differs from return position 0, 1, ...
func TestGenerateModule_ReturnDocByIndex(t *testing.T) {
	m := luareg.NewModule("retdoc", "retdoc module")
	m.Fn(
		"pair", multiRetFn, "returns a pair",
		luareg.ReturnDoc(0, "first", "the first value"),
		luareg.ReturnDoc(1, "second", "the second value"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return string first the first value")
	assert.Contains(t, out, "---@return string second the second value")
}

// TestGenerateModule_MethodNilGoFn: writeMethodStub with a nil GoFn must not
// panic and emits a minimal stub.
func TestGenerateModule_MethodNilGoFn(t *testing.T) {
	fn := &luareg.FnMeta{LuaName: "noop", GoFn: nil}
	var sb strings.Builder
	classLookup := make(map[reflect.Type]luareg.AnyClass)

	// Should not panic.
	writeMethodStub(&sb, fn, "MyClass", classLookup)
	out := sb.String()
	assert.Contains(t, out, "function MyClass:noop() end")
}

// TestGoTypeToLua_Interface: an interface{}/any type maps to "any".
func TestGoTypeToLua_Interface(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)

	// The only way to get an interface type via reflect is through a pointer
	// to interface, then Elem(). Use []interface{} → element is interface{}.
	sliceAny := reflect.TypeOf([]interface{}{})
	elemType := sliceAny.Elem()
	assert.Equal(t, "any", goTypeToLua(elemType, classLookup))
}

// TestStructLuaName: exercises the normal package-prefix branch.
func TestStructLuaName(t *testing.T) {
	// Normal in-package struct: package path ends with "stubgen", so prefix is "stubgen".
	assert.Equal(t, "stubgen.testPoint", structLuaName(reflect.TypeOf(testPoint{})))
	assert.Equal(t, "stubgen.testRect", structLuaName(reflect.TypeOf(testRect{})))
}

// TestStructLuaName_K8s: the kubernetes special-case path (k8s.io/api/group/version)
// produces groupversion.TypeName. Uses k8s.io/api/core/v1 from go.mod.
func TestStructLuaName_K8s(t *testing.T) {
	podType := reflect.TypeOf(corev1.Pod{})
	got := structLuaName(podType)
	assert.Equal(t, "corev1.Pod", got, "k8s API type should get group+version prefix")
}

// TestGenerateModule_ReturnDocNameOnly: a ReturnDoc with a name but no doc
// emits "---@return type name" without a trailing space.
func TestGenerateModule_ReturnDocNameOnly(t *testing.T) {
	m := luareg.NewModule("nd", "name only doc module")
	m.Fn(
		"get", func() string { return "" }, "gets a value",
		luareg.ReturnDoc(0, "value", ""),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return string value\n")
}

// TestGenerateModule_ReturnDocOutOfOrder: ReturnDoc entries registered out of
// position-order must still apply to the correct return. Regression test for
// the bug where buildReturnList indexed by registration position instead of
// by the .Index field.
func TestGenerateModule_ReturnDocOutOfOrder(t *testing.T) {
	m := luareg.NewModule("oo", "out of order returns")
	m.Fn(
		"two", func() (string, int) { return "", 0 }, "two returns",
		// Register index 1 BEFORE index 0 to exercise the bug path.
		luareg.ReturnDoc(1, "n", "an integer"),
		luareg.ReturnDoc(0, "s", "a string"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return string s a string\n")
	assert.Contains(t, out, "---@return number n an integer\n")
}

// TestGenerateModule_ReturnDocSparse: registering only ReturnDoc(1, ...) for a
// two-return function attaches the doc to return 1, not return 0.
func TestGenerateModule_ReturnDocSparse(t *testing.T) {
	m := luareg.NewModule("sp", "sparse return docs")
	m.Fn(
		"two", func() (string, int) { return "", 0 }, "two returns",
		luareg.ReturnDoc(1, "n", "an integer"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return string\n")
	assert.Contains(t, out, "---@return number n an integer\n")
}

// TestGenerateModule_LTableReturnFallback: a *lua.LTable return with no
// ReturnType override maps to "table", not a bogus "lua.LTable" class
// reference. Regression test for the bug where goTypeToLua dereferenced
// *lua.LTable, found an unregistered struct, and emitted structLuaName on it.
func TestGenerateModule_LTableReturnFallback(t *testing.T) {
	m := luareg.NewModule("coll", "collections module")
	m.Fn(
		"as_table", ltableReturnFn, "returns a table",
		luareg.ReturnDoc(0, "t", "the table"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return table t the table\n")
	assert.NotContains(t, out, "lua.LTable")
	assert.NotContains(t, out, "---@class lua.LTable")
}

// TestGenerateModule_ReturnTypeOverride: ReturnType overrides the reflected
// type for a return value, mirroring ArgType's precedence for parameters.
// Without this, every *lua.LTable-returning collections function would stub
// as the fallback "table" even where a more specific LuaLS type (e.g.
// "any[]") is more useful to callers.
func TestGenerateModule_ReturnTypeOverride(t *testing.T) {
	m := luareg.NewModule("coll", "collections module")
	m.Fn(
		"keys", ltableReturnFn, "returns a table's keys",
		luareg.ReturnDoc(0, "keys", "the table's keys"),
		luareg.ReturnType(0, "any[]"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return any[] keys the table's keys\n")
	assert.NotContains(t, out, "---@return table")
}

// TestGenerateModule_ReturnTypeWithDoc: ReturnType and ReturnDoc compose on
// the same return index — both the overridden type and the name/doc from
// ReturnDoc must appear together in the annotation.
func TestGenerateModule_ReturnTypeWithDoc(t *testing.T) {
	m := luareg.NewModule("coll", "collections module")
	m.Fn(
		"values", func() *lua.LTable { return nil }, "returns a table's values",
		luareg.ReturnType(0, "any[]"),
		luareg.ReturnDoc(0, "values", "the table's values in iteration order"),
	)

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@return any[] values the table's values in iteration order\n")
}

// TestGenerateModule_LStateReturnNeverLeaks: defensive test for a function
// that (incorrectly) returns *lua.LState. The runtime wrapper cannot convert
// it to a lua.LValue either, but the stub generator must not compound that by
// emitting an "LState" reference into the signature.
func TestGenerateModule_LStateReturnNeverLeaks(t *testing.T) {
	m := luareg.NewModule("escape", "escape hatch module")
	m.Fn("get_state", lstateReturnFn, "returns the Lua state (should never be stubbed)")

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.NotContains(t, out, "LState")
	assert.NotContains(t, out, "---@return")
	assert.Contains(t, out, "function escape.get_state() end")
}

// TestGoTypeToLua_LuaTable: *lua.LTable maps directly to "table" via the
// explicit type-identity check, not by falling through to the generic
// struct-pointer branch (which would emit "lua.LTable").
func TestGoTypeToLua_LuaTable(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)
	got := goTypeToLua(reflect.TypeOf((*lua.LTable)(nil)), classLookup)
	assert.Equal(t, "table", got)
}

// TestGoTypeToLua_LuaValue: lua.LValue (the interface) maps to "any" via the
// explicit type-identity check. This already worked incidentally through the
// generic reflect.Interface case, but is asserted explicitly here since it is
// now also special-cased in goTypeToLua.
func TestGoTypeToLua_LuaValue(t *testing.T) {
	classLookup := make(map[reflect.Type]luareg.AnyClass)
	got := goTypeToLua(reflect.TypeOf((*lua.LValue)(nil)).Elem(), classLookup)
	assert.Equal(t, "any", got)
}

// TestGenerateModule_PureErrorReturn: a function whose only return is `error`
// produces no ---@return lines.
func TestGenerateModule_PureErrorReturn(t *testing.T) {
	m := luareg.NewModule("e", "pure error")
	m.Fn("fail", func() error { return nil }, "may fail")

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.NotContains(t, out, "---@return")
	assert.Contains(t, out, "function e.fail() end")
}

// TestGenerateModule_Consts_PrimitiveAndStruct: a module with a primitive and a
// struct const emits ---@field lines in the ---@class block and generates a
// ---@class block for the struct type (Fix 3).
func TestGenerateModule_Consts_PrimitiveAndStruct(t *testing.T) {
	type MyTag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	m := luareg.NewModule("tags", "tags module")
	m.Const("MAX_TAGS", 16, "number", "maximum number of tags")
	m.Const("DEFAULT_TAG", MyTag{Key: "env", Value: "prod"}, "", "default tag")

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	// The module ---@class block must contain both @field lines.
	assert.Contains(t, out, "---@field MAX_TAGS number maximum number of tags\n")
	assert.Contains(t, out, "---@field DEFAULT_TAG stubgen.MyTag default tag\n")

	// The struct const type must have its own ---@class block.
	assert.Contains(t, out, "---@class stubgen.MyTag")
	assert.Contains(t, out, "---@field key string")
	assert.Contains(t, out, "---@field value string")
}

// TestGenerateModule_Consts_NoDoc: a const without a doc string emits
// "---@field name type" without a trailing space.
func TestGenerateModule_Consts_NoDoc(t *testing.T) {
	m := luareg.NewModule("nd2", "no-doc const module")
	m.Const("LIMIT", 100, "number", "")

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Contains(t, out, "---@field LIMIT number\n")
	assert.NotContains(t, out, "---@field LIMIT number \n") // no trailing space
}

// TestGenerateModule_DoubleErrorReturn: a function returning (error, error)
// has the trailing error stripped; only the first error produces a return.
func TestGenerateModule_DoubleErrorReturn(t *testing.T) {
	m := luareg.NewModule("dd", "double error")
	m.Fn("weird", func() (error, error) { return nil, nil }, "weird")

	gen := NewGenerator()
	out, err := gen.GenerateModule(m)
	require.NoError(t, err)

	assert.Equal(t, 1, strings.Count(out, "---@return"),
		"only the leading error should produce a return annotation, got:\n%s", out)
}
