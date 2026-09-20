---
name: working-here
description: How to work on the glua repository - build, test, benchmark, generate Lua stubs, add modules, upgrade dependencies, and release. Load this before making any change to this repo, including doc-only changes.
---

# Working on glua

glua embeds Lua (gopher-lua) in Go and binds Go functions, structs and types into
Lua with as little boilerplate as possible. It ships a reflection-based binding
layer, a stub generator for editor autocomplete, and a standard library of modules.

**Load this skill before touching the repo.** It is the operational source of
truth. If it disagrees with what you find in the tree, the tree wins — and then
you fix this file (see "Maintaining this skill", which is not optional).

## Golden rules

1. **Never `git commit` or `git push` unless explicitly asked in that turn.** The
   owner signs commits with GPG and commits himself. "Fix X", "make this pass"
   and "investigate Y" do NOT imply commit. Stop at the working-tree edit and
   report. Same for `rebase`, `amend`, force-push, tag push, branch deletion.
2. **Never run destructive git.** No `stash`, `pop`, `checkout`, `reset`,
   `clean`, `restore`. Read-only git only: `status`, `diff`, `log`. If you need
   to know whether a change predates your session, ask.
3. **`make` must pass before a task is complete.** Unit tests AND the Kind
   integration test. Not "the tests I touched".
4. **`make bench-update` before any commit.** The pre-commit hook does NOT check
   this. See "Benchmarks".
5. **No emojis in Go code.** (The Makefile has them; that is pre-existing, leave it.)
6. **Go doc comments use `<funcName>: <description>`**, not the stdlib
   `FuncName ...` form. This is deliberate and `revive`'s `exported` rule is
   disabled in `.golangci.yml` because of it. Do not "fix" the style.
7. **Write unit tests for every new function.** Tests must be able to fail when
   the business logic changes.
8. **Remove unused functions.** Dead exports get deleted, not kept "just in case".
9. **Keep the docs in sync** — README.md, this skill. See "Maintaining this
   skill". Note there is **no root-level API.md**; README.md is the main
   document. `pkg/modules/k8sclient/API.md` exists but is module-specific.

## Layout

```
pkg/glua/       Translator (Go <-> Lua value conversion) + TypeRegistry (offline
                reflect walker used only by stub generation; no Lua code in it)
pkg/luareg/     The binding layer. Module/Registry (module.go), Class[T] and
                metatables (class.go), the reflection wrapper that turns a Go
                func into a lua.LGFunction (reflect.go)
pkg/stubgen/    Generates library/*.gen.lua LuaLS annotations from the registry.
                comments.go pulls Go doc comments via x/tools/go/packages
pkg/modules/    The shipped stdlib. all.go::RegisterAll registers every module
cmd/glua-gen/   CLI that runs RegisterAll and writes library/*.gen.lua
cmd/run-script/ Minimal script runner
library/        GENERATED Lua stubs. Never hand-edit; run `make gen-stubs`
example/        Runnable examples; example/k8sclient is the Kind integration test
benchmarks/     Conversion + Lua access benchmarks, README.md is generated
```

## How the Lua <-> Go translation works

Two distinct paths. Knowing which one applies is the single most useful thing
when debugging a conversion bug.

**1. Data path — JSON round-trip** (`pkg/glua/translator.go`):
`ToLua` does `json.Marshal` -> `json.Unmarshal` into `interface{}` -> reflect-walk
into an `*lua.LTable`. `FromLua` reverses it. Consequences:

- **Struct field names come from `json` tags.** There is no `lua:"..."` tag.
  `json:"-"` fields are invisible to Lua; `omitempty` makes zero values vanish.
- Custom `MarshalJSON`/`UnmarshalJSON` are honoured for free.
- All numbers become `lua.LNumber` (float64) — int64 precision is lost above 2^53.
- Arrays vs maps on the way back are discriminated by `LTable.MaxN()`, so an
  empty table is a map and a table with a gap in its integer keys degrades.
- Every hand-written Lua→Go table walk in this codebase is guarded against a
  self-referential or pathologically deep Lua table via `pkg/glua.TableGuard`
  (`table_guard.go`): create one with `glua.NewTableGuard()` per top-level
  conversion call, then at every `*lua.LTable` node call `guard.Enter(tbl)`
  before recursing into its children and `defer guard.Leave(tbl)` right
  after checking the error. `Enter` tracks a *path-scoped* ancestor set (so a
  true cycle errors immediately, but a non-cyclic diamond reference — the
  same table reachable via two different fields — still converts, matching
  `encoding/json`) plus a hard `maxTableDepth` (100, defined once in
  `table_guard.go`) cap as defense-in-depth against non-cyclic runaway
  nesting. `Translator.FromLua` (`fromLuaValue`/`fromLuaValueGuarded`) uses
  this internally; so do `json.stringify`, `yaml.stringify`,
  `template.render`/`render_file`, and `spew.dump`/`sdump` — each of those
  hand-rolls its own Lua→Go conversion (different array/map-detection
  details, different marshal target) but threads a shared `*glua.TableGuard`
  through its own recursion rather than reimplementing the bookkeeping.
  Without this, a script doing `t.self = t` and handing `t` to ANY
  table-walking module function recursed forever and killed the whole
  process with an unrecoverable `fatal error: stack overflow` — not a
  catchable panic, gopher-lua's `pcall` cannot catch it and neither can
  `recover()`.
  **Enter/Leave, not a closure-returning API:** `Enter` returns only an
  `error`; it deliberately does NOT return a `func()` to call on exit. A
  closure captured per table node and returned up the stack has to escape to
  the heap, which showed up as a real, measured allocation regression on
  this hot path (recursion depth × extra allocs) during benchmarking. Do not
  "simplify" this back into a closure-returning form.
  **Do not hand-roll a new unguarded `*lua.LTable` recursion.** If you add a
  module function that walks a Lua table itself instead of relying on
  `Translator.FromLua`/`ToLua`, it MUST thread a `*glua.TableGuard` through
  that recursion the same way. Copying an old conversion loop as a starting
  point for a new module is exactly how this bug shipped in four places at
  once — grep for `TableGuard` usage in `pkg/modules/{json,yaml,template,spew}`
  for the pattern to copy instead.

**2. Handle path — userdata + metatables** (`pkg/luareg/class.go`):
a registered `Class[T]` stores the Go value intact in an `*lua.LUserData` with a
type metatable whose `__index` is a methods table. Nothing is serialised, the Go
pointer survives, and Lua sees only the registered methods. This is how
`k8sclient.Client` and `log.Logger` cross the boundary.

`pkg/luareg/reflect.go` is the glue: per argument it picks handle (registered
class -> `Class.Check`) or data (everything else -> `Translator.FromLua`).

**`[]byte` is a deliberate special case.** Top-level args and returns map to a
raw `lua.LString` (Lua strings are 8-bit clean). `[]byte` *nested inside a
struct* stays base64, because it goes through the JSON path — that is
intentional, not a bug: a k8s `Secret.Data` renders as base64 in
`kubectl get -o yaml`, which is the mental model module users already have.
Do not "fix" the nested case without a deliberate decision.

**Arity is exact — EXCEPT when the function takes `*lua.LState`.** The reflection
wrapper normally rejects both too few and too many arguments. But a function
whose first Go parameter is `*lua.LState` gets a *minimum* check only: too few
raises, surplus arguments are silently ignored.

This is deliberate, not a bug. The LState is the escape hatch for reading extra
stack slots yourself, which is exactly how `log.info("msg", {fields})` works —
`moduleLuaInfo` takes `(L, msg)` and pulls the fields from stack position 2 via
`extractFields`. Enforcing exact arity would break it.

The consequence to know: any function taking `*lua.LState` — including every one
that returns a constructed table, since building one needs the state — opts out
of surplus-argument checking. `collections.map(t, f, "junk")` is accepted.
If a function does NOT need the state, leave it out and get the stricter check.

**A nil Go slice becomes Lua `nil`, not an empty table.** It travels the JSON
path, and `json.Marshal` renders a nil slice as `null`. An empty non-nil slice
renders as `[]` and arrives as an empty table. So a function returning a nil
slice hands Lua something that breaks `#result` and `ipairs(result)`.

Several stdlib functions return nil rather than an empty slice on "no results" —
`regexp.FindAllString` and `strings.SplitN(s, sep, 0)` both do. **Always
normalise before returning:**

```go
if matches == nil {
    return []string{}, nil
}
```

`pkg/modules/regexp/regexp.go:60` is the precedent. Any new function returning a
slice needs a test for the empty case that asserts `type(x) == "table"`.

**`[]byte` has the same trap with a different symptom.** It does not take the
JSON path — it maps to a raw `lua.LString` — but a nil `[]byte` still arrives as
Lua `nil` rather than `""`. `bytes.Buffer.Bytes()` returns nil for an empty
buffer, so decompressing to an empty string, or any other empty byte result,
silently yields nil unless you normalise:

```go
if b == nil {
    b = []byte{}
}
```

`pkg/modules/compress` does this. Test the empty case asserting
`type(x) == "string"`.

**gopher-lua's `math.huge` is `math.MaxFloat64`, not `+Inf`.** Stock Lua 5.1 sets
it to `HUGE_VAL`, i.e. infinity. So `math.huge == 1/0` is **false** here, and a
test that expects `math.huge` to behave as infinity will fail confusingly. Real
IEEE infinities do exist — produce one with `1/0` (and `NaN` with `0/0`). Verified
against gopher-lua v1.1.2.

**gopher-lua diverges from stock Lua 5.1 on `setmetatable`.** Stock Lua requires
a table as arg 1, so userdata is protected for free. gopher-lua's
`baseSetMetatable` only type-checks arg 2, so a script CAN reassign a metatable
onto userdata. Classes therefore set `__metatable` to block it. Keep that in
mind before assuming Lua-reference semantics hold here.

## Build and test

```
make                 # DEFAULT. ALL tests (unit + Kind integration) + build everything
make test            # unit + k8sclient integration
make test-unit       # unit only, with -race -cover
make test-short      # unit without race (faster)
make test-verbose    # verbose, but only 3 packages
make test-k8sclient  # Kind integration test (needs kind + kubectl installed)
make build           # builds bin/glua-gen and bin/example (runs `go fmt` first)
make gen-stubs       # regenerate library/*.gen.lua
make fmt             # go fmt ./...
make clean           # rm -rf bin/
```

During development prefer `go test ./pkg/...`, but **always run `make` before
declaring a task done**. Use `-count=1` when re-running: Go caches test results
and a cached pass looks identical to a fresh one.

`make test-k8sclient` creates and destroys a Kind cluster named
`glua-k8sclient-test`. It deletes any pre-existing cluster of that name first.
Do not run it if you cannot afford a cluster churn; do not skip it when
finishing a task.

**Subagents:** do not run the full `make` from a subagent working on a shared
tree — the Kind cluster is a global side effect. Run `go build ./...`,
`go vet ./...`, `go test -race -cover ./pkg/... -count=1` and leave `make` to
the coordinating session at the end.

## Lint

`.golangci.yml` (v2 schema) is authoritative. Run `golangci-lint run ./...`.

Two deliberate deviations, do not undo them without a reason:

- `revive`'s `exported` rule is **disabled** — it demands `FuncName ...` doc
  comments, which contradicts this project's `<funcName>: <description>`
  convention (golden rule 6). Specifying any `rules:` list replaces revive's
  entire default set, which is why the defaults are enumerated explicitly.
- `gocritic`'s `captLocal` is **disabled** — it flags the `L *lua.LState`
  parameter name used throughout, which mirrors the Lua C API's `lua_State *L`.
  `exitAfterDefer` is disabled for the same kind of reason (`os.Exit` after
  `defer Close()` in `main`).

CI pins `golangci-lint-action@v9` + `version: v2.13`. Keep the pinned version and
the locally installed one in the same major line, or findings will differ
between your machine and CI.

## Benchmarks

`make bench` runs them; `make bench-update` runs them **and rewrites
`benchmarks/README.md`**. Run `bench-update` before any commit — the pre-commit
hook does not. After it runs, re-read the "Key Takeaways" section of
`benchmarks/README.md` and update the prose if the numbers moved materially;
the target only replaces the results block, not the analysis.

## Lua stubs

`library/*.gen.lua` are LuaLS annotation files that give editors autocomplete
for the modules. They are **generated** — never hand-edit.

Regenerate with `make gen-stubs` after ANY change to:

- a module's registered functions, constants, classes or their signatures,
- the doc text attached via `luareg.ArgDoc` / `ArgType` / `ReturnDoc`,
- Go doc comments on structs whose fields end up in stubs (`pkg/stubgen/comments.go`
  reads real source comments via `x/tools/go/packages`),
- the type mapping in `pkg/stubgen/generator.go`.

Then commit the regenerated files along with the change. A stub that disagrees
with the runtime is worse than no stub.

Stub doc text comes from FnOpts in `pkg/luareg/module.go`:
`luareg.ArgDoc(name, doc)` for a parameter description, `luareg.ReturnDoc(index,
name, doc)` for a return value, `luareg.ArgType(name, luaType)` to override the
inferred LuaLS type for a parameter (needed for callbacks — `lua.LValue`
carries no useful type), and `luareg.ReturnType(index, luaType)` — the same
override, indexed like `ReturnDoc`, for a return value (needed for
`*lua.LTable`/`lua.LValue` escape-hatch returns, which otherwise stub as the
generic fallback `table`/`any`). **`ArgType` has an external consumer
(matrixbot) — do not delete it.**

## Module API conventions

Rules for any Go function registered via `luareg.Fn`/`Method`, distilled from
the stdlib-module batch (SPECS.md "Part 2", F1–F6). Apply these to every new
module, not just that batch.

- **F1 — errors raise; they do not return `(nil, err)` as a normal value.**
  A trailing `error` return aborts the script (`pkg/luareg/reflect.go`,
  `raiseError`). Classify every failure mode: malformed input / programmer
  error → return a Go `error` (raises); a legitimate negative answer →
  return `false`/`nil` as a normal value. `hmac.verify_sha256(...)` returning
  `false` and `x509.parse("garbage")` raising are both correct, for different
  reasons.
- **F2 — arity is exact; there are no optional arguments.** Unless the Go
  func takes `*lua.LState` first or is variadic, a wrong argument count
  raises, and stubgen always emits `---@param x T` (never `x? T`). Where a
  default is wanted, use one of: an **options struct** (a Go struct with
  `json` tags as a required final table argument — the Translator fills it,
  so missing keys become zero values; this is the preferred mechanism), a
  **named variant** (`strconv.atoi(s)` alongside `strconv.parse_int(s,
  base)`), or a **module constant** passed explicitly
  (`compress.DEFAULT_LEVEL`). Security corollary: name options-struct fields
  so the zero value is the safe value (`allow_missing_exp bool`, never
  `require_exp bool`).
- **F3 — `[]byte` at the top level is a raw, 8-bit-clean Lua string**, not
  base64 and not a table of numbers (`pkg/luareg/reflect.go`,
  `pkg/stubgen/generator.go`'s `[]byte` special case). `[]byte` nested inside
  a struct/map field still goes through the JSON/Translator path and is
  base64 — see "How the Lua <-> Go translation works" above. Do not put
  binary data in an options struct field.
- **F4 — Lua numbers are float64.** Anything that can exceed 2^53 must not
  silently round: render it as an exact decimal string, or raise. (`bit`
  sidesteps this by choosing 32-bit semantics; `netaddr` host counts and
  `x509` serial numbers use exact strings.)
- **F5 — callbacks are `lua.LValue` + `luareg.ArgType`.** `*lua.LFunction`
  falls into the pointer-to-struct branch and gets `L.CheckTable`'d, so a
  callback parameter must be typed `lua.LValue`, validated manually
  (`v.Type() != lua.LTFunction` → raise), and annotated for the stub with
  `luareg.ArgType("fn", "fun(v: any, i: integer): any")`. Invoke with
  `L.CallByParam(lua.P{Fn: fn, NRet: 1, Protect: true}, args...)` — use
  `Protect: true` and wrap the returned error with context
  (`collections.map: callback failed at index 3: <msg>`) rather than letting
  it propagate bare. Accept only `*lua.LFunction`; do not support `__call`
  tables.
- **F6 — return-type overrides exist: use `luareg.ReturnType`.** Before this
  was added, `*lua.LTable`/`lua.LValue` returns had no override and
  `goTypeToLua` would dereference `*lua.LTable` into an unregistered
  `lua.LTable` struct and emit a bogus class reference. Un-overridden
  `*lua.LTable`/`lua.LValue` returns now stub as `table`/`any`; use
  `luareg.ReturnType(index, luaType)` when a more specific LuaLS type
  (`any[]`, `table<string, string>`, ...) is more useful to callers.

## Adding a module

1. Create `pkg/modules/<name>/<name>.go` with a `Register(reg *luareg.Registry)`.
2. Register it in `pkg/modules/all.go::RegisterAll` (keep the list alphabetical).
3. Add `<name>_test.go` in the same package. Test through real Lua with
   `L.DoString` where behaviour is Lua-visible, not just the Go functions.
4. Give every function `ArgDoc`/`ReturnDoc` text — that is what makes the stubs
   worth shipping.
5. `make gen-stubs`, then `make`.
6. Update README.md's module list.

## Dependency upgrades

`go get -u ./...` then `go mod tidy`, then `go build ./... && go vet ./... &&
go test -race ./pkg/... -count=1`. Beyond that:

### The kube-openapi trap (READ THIS BEFORE BUMPING k8s)

`k8s.io/kube-openapi` has **no semver tags, only pseudo-versions**. So `go get -u`
walks it to the newest commit on master, which may already have migrated to a
`sigs.k8s.io/structured-merge-diff` major that the pinned `apimachinery` does not
use yet. The symptom is a compile error inside apimachinery's own vendored code
mixing `structured-merge-diff/v6` and `/v7` types
(`pkg/util/managedfields/internal/typeconverter.go`).

This is **not** an upstream break to wait out. The procedure:

1. Bump `k8s.io/api`, `k8s.io/apimachinery` and `k8s.io/client-go` **together** to
   the same version.
2. Read that apimachinery version's own go.mod to find the kube-openapi
   pseudo-version it actually declares:
   `curl -s https://proxy.golang.org/k8s.io/apimachinery/@v/vX.Y.Z.mod`
3. Pin `k8s.io/kube-openapi` to **exactly** that pseudo-version.
4. `go mod tidy`, then verify `structured-merge-diff` resolves to the major
   apimachinery declares.

The same reasoning applies to any untagged `k8s.io/*` module: `-u` has no semver
to constrain it, so pin it to what its consumer declares.

### Other dependency notes

- **gopher-lua** is load-bearing for the entire library. After bumping it, run
  the full `pkg` suite and read the upstream changelog for behavioural changes,
  not just API changes.
- **golang.org/x/tools** gates which Go toolchain `pkg/stubgen` can compile
  against. Too old and it fails with `internal error: package "go/ast" without
  types was imported` under a newer Go. If stubgen breaks after a Go upgrade,
  bump x/tools first.
- If an upgrade raises the `go` directive in go.mod, it is a **consumer-visible
  constraint change** — surface it explicitly rather than letting it slide, and
  update the Go version in BOTH workflows (below).

## Go toolchain version — pinned in three places

These must move together:

- `go.mod` (the `go` directive)
- `.github/workflows/test.yml` (`go-version`, and the cache key mentions it)
- `.github/workflows/release.yml` (`go-version`)

Leaving CI below the go.mod floor makes runners silently auto-download a
toolchain on every job, defeating the pin.

## CI and local CI

Workflows: `.github/workflows/test.yml` (jobs: `test`, `lint`, `build`) and
`release.yml` (goreleaser).

Run CI locally with `act` before pushing:

```
make act-check      # is act installed
make act-list       # list workflows/jobs
make act-test-unit  # fast, recommended pre-push check
make act-lint
make act-build
make act-test       # everything incl. K8s, slow
```

## Releases

goreleaser (`.goreleaser.yml`) builds the binaries and the Lua stub tarball. The
stub asset name is **version-stamped** (`glua-stubs_<version>.tar.gz`, no `v`
prefix — goreleaser's `.Version` strips it). That is why README's download
instructions resolve the asset through the GitHub releases API rather than the
`/releases/latest/download/<asset>` path, which would need the version baked
into the filename anyway.

The README must never pin a stale version. If you reference a tag anywhere in
docs or examples, it has to be current.

## Maintaining this skill (MANDATORY)

This skill exists so agents do not re-derive the repo every session. A stale
skill is worse than none — it produces confidently wrong work.

**In the SAME change that alters any of the following, update this file:**

- a Makefile target added, removed or changed in behaviour
- the test story: new test tiers, new external prerequisites (Kind, kubectl,
  act), a change to what `make` runs
- the build/release pipeline: goreleaser config, asset names, workflow jobs,
  the pinned Go or golangci-lint version
- lint policy: a linter enabled/disabled, a new deliberate exclusion
- the binding architecture: how the Translator, Class/metatables or the
  reflect wrapper convert values, or any new special case like `[]byte`
- the stub pipeline: generator behaviour, when regeneration is required
- dependency procedure: a new trap of the kube-openapi kind
- repo layout: a new top-level package or a moved responsibility

**And in the same change, update `README.md`** whenever something of note to a
*user* of the library changes: new or removed modules, changed public API,
changed install or stub-download instructions, new requirements, version
references. `pkg/modules/k8sclient/API.md` too when that module's surface moves.

Checklist before reporting a task done:

- [ ] `make` passes (unit + Kind integration)
- [ ] `golangci-lint run ./...` clean
- [ ] `make bench-update` run if anything could affect performance
- [ ] `make gen-stubs` run if anything stub-visible changed
- [ ] README.md updated if user-visible behaviour changed
- [ ] this skill updated if any of the triggers above fired
- [ ] nothing committed unless explicitly asked
