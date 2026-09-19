# glua - Go to Lua Bridge with type safety

[![Tests](https://github.com/thomas-maurice/glua/actions/workflows/test.yml/badge.svg)](https://github.com/thomas-maurice/glua/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/thomas-maurice/glua)](https://goreportcard.com/report/github.com/thomas-maurice/glua)
[![GoDoc](https://godoc.org/github.com/thomas-maurice/glua?status.svg)](https://godoc.org/github.com/thomas-maurice/glua)

A comprehensive toolkit for embedding Lua in Go applications with full type safety and IDE autocomplete support. Designed specifically for Kubernetes API types but works with any Go structs.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage Guide](#usage-guide)
- [Creating Custom Lua Modules](#creating-custom-lua-modules)
- [IDE Setup](#ide-setup)
- [API Reference](#api-reference)
- [Features](#features)
- [Use Cases](#use-cases)
- [Performance](#performance)
- [Testing](#testing)
- [Contributing](#contributing)

## Installation

### Using in Your Project

Add glua to your `go.mod`:

```bash
go get github.com/thomas-maurice/glua@latest
```

Or import the package in your code and run `go mod tidy`, which will add the
latest version to `go.mod` for you.

### Installing Lua Stubs for IDE Autocomplete

Download the latest Lua stubs for IDE autocomplete support:

```bash
# Resolve the newest release's stub tarball and extract it to your project.
# (The asset name is versioned, e.g. glua-stubs_0.0.12.tar.gz, so we ask the
# GitHub API for whatever the latest release actually published instead of
# hardcoding a tag.)
ASSET_URL=$(curl -sL https://api.github.com/repos/thomas-maurice/glua/releases/latest \
  | grep browser_download_url \
  | grep glua-stubs \
  | cut -d '"' -f4)
curl -sL "$ASSET_URL" | tar xz

# This extracts to library/*.gen.lua
# Configure your IDE to recognize the library/ directory
```

**VS Code Setup** (with Lua extension):

Create or update `.vscode/settings.json`:

```json
{
  "Lua.workspace.library": ["library"],
  "Lua.diagnostics.disable": ["duplicate-doc-field"]
}
```

Or create a `.luarc.json` in your project root:

```json
{
  "runtime": {
    "version": "Lua 5.1"
  },
  "workspace": {
    "library": ["library"],
    "checkThirdParty": false
  },
  "diagnostics": {
    "disable": ["duplicate-doc-field"]
  }
}
```

**Why disable `duplicate-doc-field`?** The Kubernetes stubs contain many classes with the same field names (e.g., different volume types each have a `volumeID` field). This is valid, but LuaLS shows false positive warnings. Disabling this diagnostic suppresses these harmless warnings.

Now you'll get autocomplete for all glua modules in your Lua scripts!

### Installing glua-gen (Optional)

`glua-gen` regenerates `library/*.gen.lua` for glua's own built-in modules.
Most users embedding glua will instead write a tiny `tools/stubgen/main.go`
in their own project (see [Embedding glua](#embedding-glua-in-your-project) below) — that
covers both glua's modules and any custom modules you add. `glua-gen` is only
useful if you want glua's stubs alone without writing any Go boilerplate.

```bash
go install github.com/thomas-maurice/glua/cmd/glua-gen@latest
glua-gen -out ./library
```

### Cloning the Repository

To work on glua or run the examples:

```bash
git clone https://github.com/thomas-maurice/glua.git
cd glua
go mod download
```

### Requirements

- Go 1.24 or later
- For Kubernetes support: `k8s.io/client-go`, `k8s.io/api`, `k8s.io/apimachinery`
- For integration tests: `kind` and `kubectl`

## Quick Start

Here's a complete example showing the main features:

```go
package main

import (
    "github.com/thomas-maurice/glua/pkg/glua"
    "github.com/thomas-maurice/glua/pkg/modules/kubernetes"
    lua "github.com/yuin/gopher-lua"
    corev1 "k8s.io/api/core/v1"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // 1. Load Kubernetes module
    L.PreloadModule("kubernetes", kubernetes.Loader)

    // 2. Create translator for Go ↔ Lua conversion
    translator := glua.NewTranslator()

    // 3. Convert Go struct to Lua table
    pod := &corev1.Pod{ /* ... */ }
    luaTable, _ := translator.ToLua(L, pod)
    L.SetGlobal("myPod", luaTable)

    // 4. Execute Lua script that uses kubernetes module
    L.DoString(`
        local k8s = require("kubernetes")
        local pod = myPod

        -- Parse Kubernetes quantities (raise on invalid input)
        local memBytes = k8s.parse_memory(pod.spec.containers[1].resources.limits["memory"])
        local cpuMillis = k8s.parse_cpu(pod.spec.containers[1].resources.limits["cpu"])
        local timestamp = k8s.parse_time(pod.metadata.creationTimestamp)

        print(string.format("Memory: %.2f MB", memBytes / (1024 * 1024)))
        print(string.format("CPU: %d millicores", cpuMillis))

        -- Modify and return
        modifiedPod = pod
    `)

    // 5. Convert modified Lua table back to Go
    modifiedTable := L.GetGlobal("modifiedPod")
    var reconstructedPod corev1.Pod
    translator.FromLua(L, modifiedTable, &reconstructedPod)

    // Round-trip complete! Data integrity preserved
}
```

## Usage Guide

### Building

The project includes a Makefile for common tasks:

```bash
# Run all tests and build binaries (default)
make

# Run tests only
make test

# Build binaries only
make build

# Clean build artifacts
make clean

# Show help
make help
```

Built binaries:

- `bin/glua-gen` - Generates Lua LSP stubs for IDE autocomplete
- `bin/example` - Complete working example with all features

### Go to Lua Conversion

Convert any Go struct to a Lua table with full type preservation:

```go
translator := glua.NewTranslator()

// Works with any Go type
pod := &corev1.Pod{
    ObjectMeta: metav1.ObjectMeta{
        Name: "my-pod",
        CreationTimestamp: metav1.Time{Time: time.Now()},
        Labels: map[string]string{"app": "demo"},
    },
    Spec: corev1.PodSpec{
        Containers: []corev1.Container{{
            Name:  "nginx",
            Image: "nginx:latest",
        }},
    },
}

luaTable, err := translator.ToLua(L, pod)
if err != nil {
    panic(err)
}

L.SetGlobal("myPod", luaTable)
```

Features:

- Preserves timestamps (RFC3339 strings)
- Preserves resource quantities (CPU/memory strings like "100m", "256Mi")
- Handles nested structures, arrays, and maps
- Converts via JSON for robustness
- A `[]byte` field nested inside a struct or map (e.g. a Kubernetes
  `Secret.Data`) is JSON's own base64 encoding of `[]byte`, so it round-trips
  as a base64 Lua string — matching what `kubectl get -o yaml` already shows
  for the same field. This is unrelated to the raw-string `[]byte` handling
  used by `luareg.Fn`/`Class.Method` top-level arguments and returns (see
  "Register with `luareg`" below) — the Translator's JSON path is not
  involved there.

### Lua to Go Conversion

Convert Lua tables back to Go structs with type safety:

```go
// After Lua modifies the table
modifiedTable := L.GetGlobal("modifiedPod")

var reconstructedPod corev1.Pod
err := translator.FromLua(L, modifiedTable, &reconstructedPod)
if err != nil {
    panic(err)
}

// Full round-trip integrity - data is identical to original
```

Features:

- Full round-trip integrity (original == reconstructed)
- Automatic type coercion
- Preserves all complex Kubernetes types
- Handles any `LValue` (LTable, LString, LNumber, etc.)

### Type Registry and LSP Stub Generation

Generate Lua LSP annotations for IDE autocomplete:

```go
registry := glua.NewTypeRegistry()

// Register your types
registry.Register(&corev1.Pod{})
registry.Register(&corev1.Service{})
registry.Register(&corev1.ConfigMap{})

// Process and generate stubs
registry.Process()
stubs, _ := registry.GenerateStubs()

// Write to file for IDE consumption
os.WriteFile("annotations.gen.lua", []byte(stubs), 0644)
```

This generates complete type definitions:

```lua
---@meta

---@class corev1.Pod
---@field kind string
---@field apiVersion string
---@field metadata v1.ObjectMeta
---@field spec corev1.PodSpec
---@field status corev1.PodStatus

---@class corev1.PodSpec
---@field containers corev1.Container[]
---@field volumes corev1.Volume[]
---@field nodeName string
-- ... all fields with correct types

---@class corev1.Container
---@field name string
---@field image string
---@field resources corev1.ResourceRequirements
-- ... complete definitions
```

Field descriptions come from the field's Go doc comment by default — reflection
alone can't see it, so `stubgen.Generator` (used by `RegisterType`,
`GenerateModule`, `GenerateFromRegistry`) resolves it from source via
`golang.org/x/tools/go/packages`, falling back to the field's trailing inline
comment if it has no doc comment above it:

```go
type Message struct {
    // Body is the plaintext message body.
    Body string `json:"body"`
}
```

```lua
---@class mymod.Message
---@field body string Body is the plaintext message body.
```

Reading source is a dev-time-only, best-effort step: if the package can't be
resolved (stripped module cache, vendored-only build, etc.) the field simply
gets no description — stub generation never fails because of it. To force a
specific description regardless of the doc comment (or when there isn't one),
add a `luadoc` struct tag; it always wins:

```go
type Message struct {
    Body string `json:"body" luadoc:"plaintext message body"`
}
```

Precedence is `luadoc` tag > Go doc/inline comment > no description.

Note this comment resolution only happens through `stubgen.Generator`. Calling
`glua.TypeRegistry.GenerateStubs()` directly, as in the example above, only
picks up `luadoc` tags unless you also call `SetFieldDocFunc` yourself.

Now in your Lua scripts, you get full autocomplete:

```lua
---@type corev1.Pod
local pod = myPod

-- IDE shows autocomplete for all fields!
print(pod.metadata.name)
print(pod.spec.containers[1].image)
```

### Kubernetes Module

The built-in Kubernetes module provides utility functions for parsing K8s resource quantities and timestamps.

Load in Go:

```go
L.PreloadModule("kubernetes", kubernetes.Loader)
```

Use in Lua:

```lua
local k8s = require("kubernetes")

-- Parse memory quantities (returns bytes)
local memBytes = k8s.parse_memory("256Mi")   -- 268435456
local memBytes2 = k8s.parse_memory("1Gi")    -- 1073741824

-- Parse CPU quantities (returns millicores)
local cpuMillis = k8s.parse_cpu("100m")      -- 100
local cpuMillis2 = k8s.parse_cpu("1.5")      -- 1500

-- Parse timestamps (returns Unix timestamp)
local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")  -- 1759509540

-- Format timestamps (Unix timestamp → RFC3339 string)
local timeStr = k8s.format_time(1759509540)  -- "2025-10-03T16:39:00Z"

-- All functions raise on error; use pcall to handle errors gracefully
local ok, err = pcall(k8s.parse_memory, "invalid")
if not ok then
    print("Parse error: " .. tostring(err))
end
```

Example use case - process pod resource limits/requests:

```lua
local k8s = require("kubernetes")
local pod = myPod

for i, container in ipairs(pod.spec.containers) do
    print("Container: " .. container.name)

    -- Parse memory limit (raises on invalid input)
    if container.resources.limits["memory"] then
        local memBytes = k8s.parse_memory(container.resources.limits["memory"])
        print(string.format("  Memory limit: %.2f MB", memBytes / (1024 * 1024)))
    end

    -- Parse CPU limit (raises on invalid input)
    if container.resources.limits["cpu"] then
        local cpuMillis = k8s.parse_cpu(container.resources.limits["cpu"])
        print(string.format("  CPU limit: %d millicores", cpuMillis))
    end
end
```

### K8s Client Module

The `k8sclient` module provides a dynamic Kubernetes client for Lua, allowing full CRUD operations on any Kubernetes resource directly from Lua scripts.

Load in Go:

```go
import "github.com/thomas-maurice/glua/pkg/modules/k8sclient"

config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
L.PreloadModule("k8sclient", k8sclient.Loader(config))
```

Use in Lua:

```lua
local k8sclient = require("k8sclient")
local client = k8sclient.new_client()

-- Define GVK (or use predefined constants like k8sclient.POD)
local pod_gvk = {group = "", version = "v1", kind = "Pod"}

-- Create a Pod (raises on error)
local pod = {
    apiVersion = "v1",
    kind = "Pod",
    metadata = {name = "nginx", namespace = "default"},
    spec = {
        containers = {{
            name = "nginx",
            image = "nginx:alpine"
        }}
    }
}
local created = client:create(pod)

-- Get a resource (raises on error)
local fetched = client:get(pod_gvk, "default", "nginx")

-- Update a resource (raises on error)
fetched.metadata.labels = {app = "web"}
local updated = client:update(fetched)

-- List resources (raises on error)
local pods = client:list(pod_gvk, "default")
for i, pod in ipairs(pods) do
    print(pod.metadata.name)
end

-- Delete a resource (raises on error)
client:delete(pod_gvk, "default", "nginx")
```

**Complete Example:** See [example/k8sclient/](./example/k8sclient) for a full working example with nginx Pod, ConfigMaps, and Kind cluster integration.

**Run the example:**

```bash
make test-k8sclient  # Runs with temporary Kind cluster
```

### Error Handling

Proper error handling throughout:

```go
// Go side
luaTable, err := translator.ToLua(L, pod)
if err != nil {
    log.Fatalf("Conversion failed: %v", err)
}

err = translator.FromLua(L, luaTable, &reconstructedPod)
if err != nil {
    log.Fatalf("Reconstruction failed: %v", err)
}
```

```lua
-- Lua side
local k8s = require("kubernetes")

-- Functions raise on error; use pcall to handle gracefully
local ok, result = pcall(k8s.parse_memory, "256Mi")
if not ok then
    print("Error: " .. tostring(result))
    return
end

print("Parsed successfully: " .. result)
```

### Round-Trip Integrity

glua ensures perfect round-trip conversion:

```go
// Original Go struct
originalPod := sample.GetPod()

// Convert to Lua
luaTable, _ := translator.ToLua(L, originalPod)
L.SetGlobal("pod", luaTable)

// Process in Lua (even if unchanged)
L.DoString("modifiedPod = pod")

// Convert back to Go
modifiedTable := L.GetGlobal("modifiedPod")
var reconstructedPod corev1.Pod
translator.FromLua(L, modifiedTable, &reconstructedPod)

// Verify integrity
originalJSON, _ := json.Marshal(originalPod)
reconstructedJSON, _ := json.Marshal(reconstructedPod)

if string(originalJSON) == string(reconstructedJSON) {
    fmt.Println("Perfect round-trip!")
}
```

This works for:

- Timestamps (RFC3339 strings preserved)
- Resource quantities ("100m", "256Mi" preserved as strings)
- Maps and labels
- Nested arrays and structures
- Complex Kubernetes objects

## Creating Custom Lua Modules

A glua module is a Go package that exposes idiomatic Go functions and types to
Lua. You write plain Go (no `lua.LState` plumbing), `luareg` handles argument
conversion, error propagation, and UserData wrapping via reflection. The same
metadata that wires runtime calls also drives stub generation — the LSP gets
proper type-annotated autocomplete for free.

This walk-through builds one complete `widget` module that exercises every
mechanism: plain functions, structs that cross the boundary, a stateful class
with methods, structured errors, and module-level constants. The final form
runs as-is; copy and adapt.

### 1. Plain Go

Start with idiomatic Go. The module has zero awareness of Lua at this point.

```go
package widget

import "errors"

// Order is a value that crosses the boundary as a Lua table. The json tags
// become Lua field names AND the field names in the generated stub.
type Order struct {
    ID    string  `json:"id"`
    Total float64 `json:"total"`
    Items int     `json:"items"`
}

// Receipt: what the script gets back from CalculateReceipt.
type Receipt struct {
    OrderID string  `json:"order_id"`
    Tax     float64 `json:"tax"`
    Grand   float64 `json:"grand_total"`
}

// CalculateReceipt: function-style. (Receipt, error) → Lua returns the
// Receipt on success, raises on non-nil error.
func CalculateReceipt(o Order, taxRate float64) (Receipt, error) {
    if o.Total < 0 {
        return Receipt{}, errors.New("negative total")
    }
    tax := o.Total * taxRate
    return Receipt{OrderID: o.ID, Tax: tax, Grand: o.Total + tax}, nil
}

// Cart: a stateful object exposed as a Lua UserData class.
type Cart struct{ items []string }

func NewCart() *Cart                    { return &Cart{} }
func (c *Cart) Add(item string)         { c.items = append(c.items, item) }
func (c *Cart) Items() []string         { return c.items }
func (c *Cart) Size() int               { return len(c.items) }
func (c *Cart) Pop() (string, error) {
    if len(c.items) == 0 {
        // Structured error: Lua receives a {message, kind} table rather
        // than a string. Useful so pcall callers can dispatch on err.kind.
        return "", &luareg.Error{Kind: "Empty", Message: "cart is empty"}
    }
    n := len(c.items) - 1
    last, c.items = c.items[n], c.items[:n]
    return last, nil
}
```

The `Cart` is intentionally a normal Go type — no `lua.LState` in any
signature. The Lua bridge is added in step 2.

### 2. Register with `luareg`

```go
import (
    "github.com/thomas-maurice/glua/pkg/luareg"
    lua "github.com/yuin/gopher-lua"
)

// build: constructs the module definition once. Loader and Register both
// reuse it — Loader to push into a Lua state, Register to record metadata
// for stub generation.
func build() *luareg.Module {
    m := luareg.NewModule("widget", "widget business logic")

    // Module-level function. Reflection inspects CalculateReceipt's signature
    // and generates the Lua wrapper. Args/ArgDoc/ReturnDoc enrich the stubs.
    m.Fn("calculate_receipt", CalculateReceipt, "compute tax and grand total",
        luareg.Args("order", "tax_rate"),
        luareg.ArgDoc("order", "the order to bill"),
        luareg.ArgDoc("tax_rate", "tax rate in [0,1]"),
        luareg.ReturnDoc(0, "receipt", "the computed receipt"),
    )

    // Class registration. The [*Cart] type parameter is the receiver type;
    // no instance is required. Methods are passed as Go method expressions.
    cart := luareg.NewClass[*Cart]("widget.Cart", "a mutable shopping cart")
    cart.Method("add",   (*Cart).Add,   "add an item",          luareg.Args("item"))
    cart.Method("items", (*Cart).Items, "return all items")
    cart.Method("size",  (*Cart).Size,  "return the number of items")
    cart.Method("pop",   (*Cart).Pop,   "remove and return the last item; raises Empty if no items")
    m.RegisterClass(cart)

    // Factory. Returning *Cart auto-wraps the value as Lua UserData bound
    // to the widget.Cart metatable.
    m.Fn("new_cart", NewCart, "create a new empty cart")

    // Module-level constant. Set on the Lua module table at PushTo time;
    // also emitted as a ---@field annotation in the stub.
    m.Const("DEFAULT_TAX_RATE", 0.2, "number", "default sales-tax rate")

    return m
}

// Loader: gopher-lua module entry point.
//   L.PreloadModule("widget", widget.Loader)
func Loader(L *lua.LState) int { return build().PushTo(L) }

// Register: stub-generation entry point. Called by the host's tools/stubgen.
func Register(reg *luareg.Registry) { build().Register(reg) }
```

Notes:

- `build()` is the single source of truth — both runtime and stub generation
  use it, so they cannot drift.
- For methods that genuinely need access to the Lua stack (e.g. true variadic
  args), use `*lua.LState` as the second parameter after the receiver. It is
  invisible to argument counting and stubs. See `pkg/modules/log/log.go` for
  the `Logger:with` pattern.
- Module-level constants are set on the module table at `PushTo` time, and
  also appear as `---@field NAME type doc` annotations in the stub.
- When a parameter's Go type reflects to something unhelpful for LuaLS (e.g.
  a `lua.LValue` used as a callback), override the annotation with
  `luareg.ArgType(name, luaType)`:

  ```go
  m.Fn("on_message", handler, "register a handler for room message events",
      luareg.Args("handler"),
      luareg.ArgType("handler", "fun(evt: core.Event)"),
  )
  ```

  generates `---@param handler fun(evt: core.Event)` instead of
  `---@param handler any`. It has no runtime effect — stub generation only.
- The same escape hatch exists for return values: `luareg.ReturnType(index,
  luaType)` overrides the Nth return's inferred type (0-indexed, excluding a
  trailing `error` return). Use it when a return type's Go type carries no
  useful Lua type on its own — e.g. a `*lua.LTable` built by hand and returned
  as an escape hatch, where a more specific type (`any[]`, `table<string,
  string>`, ...) is more useful to callers than the generic fallback:

  ```go
  m.Fn("keys", collectionsKeys, "returns a table's keys as a new table",
      luareg.ReturnDoc(0, "keys", "the table's keys"),
      luareg.ReturnType(0, "any[]"),
  )
  ```

  Without an override, an un-annotated `*lua.LTable` return stubs as `table`
  and a `lua.LValue` return stubs as `any` — never a bogus class reference.
  It has no runtime effect — stub generation only.
- A top-level `[]byte` argument or return value (directly in a function or
  method signature, not nested inside a struct/map field) maps to a raw Lua
  string, not base64 and not a table of numbers — Lua strings are 8-bit
  clean, so this is a lossless, zero-copy representation, and stubs emit
  `string` accordingly. Arguments also still accept a table of numbers for
  backwards compatibility. This is the opposite of `[]byte` nested inside a
  struct field, which goes through `Translator.ToLua`/`FromLua` (JSON) and is
  base64-encoded — see "Go to Lua Conversion" above.

### 3. Use from Lua

Once `widget.Loader` is preloaded into a `*lua.LState`, scripts can use it:

```lua
local widget = require("widget")

-- Function-style call: struct arg in, struct result out.
local order = {id = "A-100", total = 50.0, items = 3}
local receipt = widget.calculate_receipt(order, widget.DEFAULT_TAX_RATE)
print(receipt.grand_total)   -- 60.0

-- UserData class with : method syntax.
local cart = widget.new_cart()
cart:add("apple")
cart:add("bread")
print(cart:size())           -- 2
print(cart:pop())            -- bread

-- Structured error: pcall returns a table the caller can dispatch on.
cart:pop()                   -- ok
local ok, err = pcall(function() return cart:pop() end)
if not ok and type(err) == "table" and err.kind == "Empty" then
    print("nothing left in the cart")
end
```

### 4. Testing

Test your module the same way glua tests its own: spin up a real `lua.LState`,
preload the module, run a Lua snippet, assert outcomes. No mocks.

```go
package widget_test

import (
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/thomas-maurice/glua/example/widget"
    lua "github.com/yuin/gopher-lua"
)

func TestCalculateReceipt(t *testing.T) {
    L := lua.NewState()
    defer L.Close()
    L.PreloadModule("widget", widget.Loader)

    require.NoError(t, L.DoString(`
        local widget = require("widget")
        local r = widget.calculate_receipt({id = "X", total = 100, items = 1}, 0.2)
        assert(r.grand_total == 120, "got " .. r.grand_total)
    `))
}

func TestCartEmptyPopRaisesStructured(t *testing.T) {
    L := lua.NewState()
    defer L.Close()
    L.PreloadModule("widget", widget.Loader)

    require.NoError(t, L.DoString(`
        local widget = require("widget")
        local cart = widget.new_cart()
        local ok, err = pcall(function() return cart:pop() end)
        assert(not ok)
        assert(type(err) == "table" and err.kind == "Empty", "expected Empty error")
    `))
}
```

For a real-world example with file-glob-loaded testdata, see
`pkg/modules/strings/strings_test.go` and `pkg/modules/strings/testdata/*.lua`.

### Error handling recap

Go `error` returns auto-raise on the Lua side. There are two flavours:

- **Plain `errors.New(...)` / `fmt.Errorf(...)`** raises a Lua string error.
  Catch with `pcall`; `err` is a string.
- **`*luareg.Error{Kind, Message}`** raises a Lua table `{kind=..., message=...}`.
  Catch with `pcall`; branch on `err.kind`. Use this when callers genuinely
  need to dispatch on the failure category (NotFound, Expired, …) rather
  than string-match.

The `(value, err)` tuple-return pattern from the pre-luareg era is gone.

### Generating stubs

```bash
# From the glua repo root:
make gen-stubs
# or:
go run ./cmd/glua-gen -out library
```

This writes one `library/<module>.gen.lua` file per registered module. The generated files are what IDEs consume for autocomplete.

### Embedding glua in your project

When you embed glua in your own application as a downstream consumer and add custom modules, ship a `tools/stubgen/main.go` that merges glua's built-in modules with yours:

**Directory layout:**

```
myproject/
  main.go
  modules/
    widget/
      widget.go      # your module with Loader + Register
  tools/
    stubgen/
      main.go        # stub generator
  library/           # generated stubs land here
  .luarc.json
```

**`tools/stubgen/main.go` template:**

```go
package main

import (
    "log"

    "github.com/thomas-maurice/glua/pkg/luareg"
    "github.com/thomas-maurice/glua/pkg/modules"
    "github.com/thomas-maurice/glua/pkg/stubgen"
    "github.com/myorg/myproject/modules/widget"
)

func main() {
    reg := luareg.NewRegistry()
    modules.RegisterAll(reg)   // glua's 20 built-in modules
    widget.Register(reg)       // your module

    gen := stubgen.NewGenerator()
    files, err := gen.GenerateFromRegistry(reg, "./library")
    if err != nil {
        log.Fatal(err)
    }
    for _, f := range files {
        log.Println("wrote", f)
    }
}
```

Run it:

```bash
go run ./tools/stubgen
```

Both glua's stubs and your custom module stubs land in `./library/`.

**`.luarc.json` (LuaLS config):**

```json
{
  "runtime": {
    "version": "Lua 5.1"
  },
  "workspace": {
    "library": ["library"],
    "checkThirdParty": false
  },
  "diagnostics": {
    "disable": ["duplicate-doc-field"]
  }
}
```

**Loading your module at runtime:**

```go
L.PreloadModule("widget", widget.Loader)
```

### Generated stub shape

For the `widget` module from the walk-through above, `go run ./tools/stubgen`
produces `library/widget.gen.lua`. Note that `Order` and `Receipt` get
`---@class` blocks automatically — they were discovered as struct types in the
function signatures, with field names taken from JSON tags. **You never wrote
an annotation comment.**

```lua
---@meta widget

---@class widget.Order
---@field id string
---@field items number
---@field total number

---@class widget.Receipt
---@field grand_total number
---@field order_id string
---@field tax number

---@class widget.Cart
local Cart = {}

--- add an item
---@param item string
function Cart:add(item) end

--- return all items
---@return string[]
function Cart:items() end

--- return the number of items
---@return number
function Cart:size() end

--- remove and return the last item; raises Empty if no items
---@return string
function Cart:pop() end

---@class widget
---@field Cart widget.Cart
---@field DEFAULT_TAX_RATE number default sales-tax rate
local widget = {}

--- compute tax and grand total
---@param order widget.Order the order to bill
---@param tax_rate number tax rate in [0,1]
---@return widget.Receipt receipt the computed receipt
function widget.calculate_receipt(order, tax_rate) end

--- create a new empty cart
---@return widget.Cart
function widget.new_cart() end

widget.Cart = Cart

return widget
```

## IDE Setup

This section explains how to enable autocomplete for your Lua scripts.

### VSCode Setup

1. Install Lua Language Server extension: [Lua](https://marketplace.visualstudio.com/items?itemName=sumneko.lua)

2. Generate stubs:

```bash
# Generate module stubs (for kubernetes, custom modules)
make gen-stubs  # Creates library/kubernetes.gen.lua, etc.

# Generate type stubs (run your app that uses TypeRegistry)
go run .  # Creates annotations.gen.lua
```

3. Create `.vscode/settings.json`:

```json
{
  "Lua.workspace.library": [
    "${workspaceFolder}/library",
    "${workspaceFolder}/annotations.gen.lua"
  ],
  "Lua.runtime.version": "Lua 5.1",
  "Lua.diagnostics.globals": ["myPod", "originalPod"]
}
```

4. Reload VSCode and enjoy autocomplete.

### Neovim Setup

1. Install lua-language-server:

```vim
:MasonInstall lua-language-server
```

2. Generate stubs:

```bash
make gen-stubs
go run .  # If using TypeRegistry
```

3. Create `.luarc.json` in project root:

```json
{
  "runtime": { "version": "Lua 5.1" },
  "workspace": {
    "library": [".", "library", "annotations.gen.lua"],
    "checkThirdParty": false
  },
  "diagnostics": {
    "globals": ["myPod", "originalPod", "modifiedPod"]
  }
}
```

4. Restart LSP: `:LspRestart`

### What You Get

```lua
---@type corev1.Pod
local pod = myPod

-- Full autocomplete with Ctrl+Space
pod.metadata.name
pod.spec.containers[1].image

local k8s = require("kubernetes")
k8s.parse_memory("256Mi")  -- Shows parameters and return types
```

## API Reference

This section is split into two parts: the **Go API** (for embedding Lua in your Go application) and the **Lua Standard Modules** (available to your Lua scripts).

### Go API

These are the Go packages and types you use in your Go application to interact with Lua.

#### Translator

Bidirectional converter between Go structs and Lua tables.

```go
type Translator struct{}

// NewTranslator: creates a new bidirectional Go ↔ Lua translator
func NewTranslator() *Translator

// ToLua: converts a Go value to a Lua value
// Supports structs, maps, slices, primitives
// Preserves timestamps and resource quantities
func (t *Translator) ToLua(L *lua.LState, o interface{}) (lua.LValue, error)

// FromLua: converts a Lua value to a Go value
// Accepts any LValue (LTable, LString, LNumber, etc.)
// Requires pointer to output variable
func (t *Translator) FromLua(L *lua.LState, lv lua.LValue, output interface{}) error
```

**Usage:**

```go
translator := glua.NewTranslator()

// Go → Lua
pod := &corev1.Pod{...}
luaTable, err := translator.ToLua(L, pod)
L.SetGlobal("myPod", luaTable)

// Lua → Go
modifiedTable := L.GetGlobal("myPod")
var reconstructedPod corev1.Pod
err = translator.FromLua(L, modifiedTable, &reconstructedPod)
```

#### TypeRegistry

Generates Lua LSP annotations for IDE autocomplete from Go types.

```go
type TypeRegistry struct{}

// NewTypeRegistry: creates a new type registry for stub generation
func NewTypeRegistry() *TypeRegistry

// Register: registers a Go type for Lua stub generation
func (r *TypeRegistry) Register(obj interface{}) error

// Process: processes all registered types and their dependencies
func (r *TypeRegistry) Process() error

// GenerateStubs: generates Lua LSP annotation code
func (r *TypeRegistry) GenerateStubs() (string, error)
```

**Usage:**

```go
registry := glua.NewTypeRegistry()
registry.Register(&corev1.Pod{})
registry.Register(&corev1.Service{})
registry.Process()
stubs, _ := registry.GenerateStubs()
os.WriteFile("annotations.gen.lua", []byte(stubs), 0644)
```

### Lua Standard Modules

These are the modules available to your Lua scripts via `require()`. Load them in Go with `L.PreloadModule()`.

#### kubernetes

Utility functions for parsing and formatting Kubernetes resource quantities and timestamps.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/kubernetes"

L.PreloadModule("kubernetes", kubernetes.Loader)
```

**Lua API:**

```lua
local k8s = require("kubernetes")

-- Parse memory: "256Mi" → 268435456 (bytes)
bytes = k8s.parse_memory(quantity)

-- Parse CPU: "100m" → 100 (millicores)
millis = k8s.parse_cpu(quantity)

-- Parse duration: "5m" → 300 (seconds)
seconds = k8s.parse_duration(duration)

-- Parse time: "2025-10-03T16:39:00Z" → 1759509540 (Unix timestamp)
timestamp = k8s.parse_time(timestr)

-- Format time: 1759509540 → "2025-10-03T16:39:00Z"
timestr = k8s.format_time(timestamp)

-- Format duration: 300 → "5m0s"
duration = k8s.format_duration(seconds)

-- Initialize defaults: ensures metadata.labels and metadata.annotations exist
obj = k8s.init_defaults(obj)

-- Add/manipulate labels and annotations
obj = k8s.add_label(obj, "app", "nginx")
obj = k8s.add_annotation(obj, "version", "1.0")
has = k8s.has_label(obj, "app")
value = k8s.get_label(obj, "app")
obj = k8s.remove_label(obj, "app")

-- Match GVK (Group/Version/Kind)
matches = k8s.match_gvk(obj, {group="apps", version="v1", kind="Deployment"})
```

All functions raise a Lua error on invalid input. Use `pcall` to handle errors gracefully.

#### k8sclient

Dynamic Kubernetes client for CRUD operations on any Kubernetes resource from Lua.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/k8sclient"

config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
L.PreloadModule("k8sclient", k8sclient.Loader(config))
```

**Lua API:**

```lua
local k8sclient = require("k8sclient")

-- GVK constants (predefined)
k8sclient.POD          -- {group="", version="v1", kind="Pod"}
k8sclient.DEPLOYMENT   -- {group="apps", version="v1", kind="Deployment"}
k8sclient.SERVICE      -- {group="", version="v1", kind="Service"}
k8sclient.CONFIGMAP    -- {group="", version="v1", kind="ConfigMap"}
-- ... and many more

-- Create a client (needs a config passed from Go via Loader)
local client = k8sclient.new_client()

-- Create a resource (raises on error)
created = client:create(resource_table)

-- Get a resource (raises on error)
resource = client:get(gvk, namespace, name)

-- Update a resource (raises on error)
updated = client:update(resource_table)

-- List resources (raises on error)
resources = client:list(gvk, namespace)

-- Delete a resource (raises on error)
client:delete(gvk, namespace, name)
```

**Example:**

```lua
local k8sclient = require("k8sclient")
local client = k8sclient.new_client()

-- Create a ConfigMap
local cm = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {name = "my-config", namespace = "default"},
    data = {key = "value"}
}
local created = client:create(cm)

-- Get it back
local fetched = client:get(k8sclient.CONFIGMAP, "default", "my-config")

-- Update it
fetched.data.newkey = "newvalue"
local updated = client:update(fetched)

-- Delete it
client:delete(k8sclient.CONFIGMAP, "default", "my-config")
```

#### json

JSON encoding and decoding.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/json"

L.PreloadModule("json", json.Loader)
```

**Lua API:**

```lua
local json = require("json")

-- Parse JSON string to Lua table (raises on invalid JSON)
table = json.parse('{"name":"John","age":30}')

-- Stringify Lua table to JSON (raises on error)
jsonstr = json.stringify({name="John", age=30})
```

#### yaml

YAML encoding and decoding.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/yaml"

L.PreloadModule("yaml", yaml.Loader)
```

**Lua API:**

```lua
local yaml = require("yaml")

-- Parse YAML string to Lua table (raises on invalid YAML)
table = yaml.parse("name: John\nage: 30")

-- Stringify Lua table to YAML (raises on error)
yamlstr = yaml.stringify({name="John", age=30})
```

#### spew

Pretty-printing for debugging (like Go's spew package).

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/spew"

L.PreloadModule("spew", spew.Loader)
```

**Lua API:**

```lua
local spew = require("spew")

-- Dump to string (returns formatted string)
str = spew.sdump({name="John", nested={deep={value=42}}})

-- Dump to stdout (prints directly)
spew.dump({name="John", age=30})
```

#### http

HTTP client for making requests. Every request (get/post/put/delete/request)
is bounded by a 30-second timeout so a slow or unresponsive endpoint cannot
block the calling goroutine forever; a timed-out request raises a Lua error
like any other network failure.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/http"

L.PreloadModule("http", http.Loader)
```

**Lua API:**

```lua
local http = require("http")

-- GET request (raises on network/HTTP error)
response = http.get("https://api.example.com/data")
-- response = {status=200, body="...", headers={...}}

-- POST request (raises on network/HTTP error)
-- Arguments are positional: (url, body, headers)
response = http.post(
    "https://api.example.com/data",
    '{"key":"value"}',
    {["Content-Type"] = "application/json"}
)

-- Body is required (pass "" for none); headers may be nil
response = http.post("https://api.example.com/ping", "")

-- Other methods: put, patch, delete
```

#### template

Go template rendering.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/template"

L.PreloadModule("template", template.Loader)
```

**Lua API:**

```lua
local template = require("template")

-- Render template with data (raises on template parse/execute error)
result = template.render("Hello {{.name}}, you are {{.age}} years old",
    {name="John", age=30})
-- result = "Hello John, you are 30 years old"
```

#### fs

Filesystem operations.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/fs"

L.PreloadModule("fs", fs.Loader)
```

**Lua API:**

```lua
local fs = require("fs")

-- Read file (raises on error)
content = fs.read_file("/path/to/file.txt")

-- Write file (raises on error)
fs.write_file("/path/to/file.txt", "content")

-- Check existence
exists = fs.exists("/path/to/file")

-- Create directory (raises on error)
fs.mkdir("/path/to/dir")
fs.mkdir_all("/path/to/nested/dir")

-- Remove (raises on error)
fs.remove("/path/to/file")
fs.remove_all("/path/to/dir")

-- List directory (raises on error)
files = fs.list("/path/to/dir")

-- Get file info (raises on error)
info = fs.stat("/path/to/file")
-- info = {size=1234, mode=420, mod_time=1234567890, is_dir=false}
-- mode is the raw Go FileMode as a number (420 == 0644), not an octal string
```

#### time

Time manipulation and formatting.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/time"

L.PreloadModule("time", time.Loader)
```

**Lua API:**

```lua
local time = require("time")

-- Current Unix timestamp
now = time.now()

-- Parse date string (raises on parse error)
-- Arguments are (timestr, layout) - the value first, the Go layout second
timestamp = time.parse("2025-10-21 14:30:00", "2006-01-02 15:04:05")

-- Format timestamp
datestr = time.format(timestamp, "2006-01-02 15:04:05")

-- Sleep
time.sleep(2)  -- sleep for 2 seconds
```

#### base64, hex, hash

Encoding and hashing utilities.

**Load in Go:**

```go
import (
    "github.com/thomas-maurice/glua/pkg/modules/base64"
    "github.com/thomas-maurice/glua/pkg/modules/hex"
    "github.com/thomas-maurice/glua/pkg/modules/hash"
)

L.PreloadModule("base64", base64.Loader)
L.PreloadModule("hex", hex.Loader)
L.PreloadModule("hash", hash.Loader)
```

**Lua API:**

```lua
local base64 = require("base64")
local hex = require("hex")
local hash = require("hash")

-- Base64
encoded = base64.encode("hello")
decoded = base64.decode(encoded)  -- raises on invalid base64

-- Hex
encoded = hex.encode("hello")
decoded = hex.decode(encoded)  -- raises on invalid hex

-- Hash strings
md5 = hash.md5("hello")
sha1 = hash.sha1("hello")
sha256 = hash.sha256("hello")
sha512 = hash.sha512("hello")

-- Hash Lua tables (converted to JSON; raises on conversion error)
hash_val = hash.md5_obj({name="John", age=30})
hash_val = hash.sha1_obj({key="value"})
hash_val = hash.sha256_obj({foo="bar", nested={data=123}})
hash_val = hash.sha512_obj({items={1, 2, 3}})
```

**Object Hashing:**

The `*_obj` functions convert Lua tables to JSON before hashing, making them useful for:

- Content-based resource identifiers
- Detecting configuration changes
- Caching keys for complex data structures
- Checksums for nested objects

```lua
-- Example: Detect if a Pod spec has changed
local hash = require("hash")
local pod_hash = hash.sha256_obj(pod.spec)
-- Store/compare this hash to detect changes
```

#### log

Structured logging with fields (similar to logrus).

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/log"

L.PreloadModule("log", log.Loader)

// Level and output format are configured on the Go side: build a
// *charmbracelet/log.Logger the way you want it and inject it, and every
// Lua-side call logs through it. There is no Lua-facing level/format setter.
log.InjectLogger(L, myLogger)
```

**Lua API:**

```lua
local log = require("log")

-- Simple logging
log.info("Application started")
log.warn("Deprecated feature used")
log.error("Connection failed")
log.debug("Debug information")

-- Structured logging with fields
log.info("User logged in", {user_id=123, ip="1.2.3.4"})

-- Logger with preset fields: get the default logger, then derive a child
logger = log.logger():with({component="api", version="1.0"})
logger:info("Request received", {path="/api/users"})
-- Output includes: component=api version=1.0 path=/api/users

-- Logger objects expose the same levels as the module
logger:debug("Debug information")
logger:warn("Deprecated feature used")
```

#### osmod

Operating system utilities for environment variables, hostname, and temp directories.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/osmod"

L.PreloadModule("osmod", osmod.Loader)
```

**Lua API:**

```lua
local osmod = require("osmod")

-- Environment variables
value = osmod.getenv("PATH")
osmod.setenv("MY_VAR", "my_value")
osmod.unsetenv("MY_VAR")

-- System information
hostname = osmod.hostname()

-- Temporary directory
tmpdir = osmod.tmpdir()
```

#### filepath

Path manipulation utilities for file and directory paths.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/filepath"

L.PreloadModule("filepath", filepath.Loader)
```

**Lua API:**

```lua
local filepath = require("filepath")

-- Join path components
path = filepath.join("/usr", "local", "bin")  -- "/usr/local/bin"

-- Split path into directory and file
dir, file = filepath.split("/usr/local/bin/tool")  -- "/usr/local/bin", "tool"

-- Get absolute path (raises on error)
abspath = filepath.abs("../relative/path")

-- Get file extension
ext = filepath.ext("/path/to/file.txt")  -- ".txt"

-- Get base name
base = filepath.base("/path/to/file.txt")  -- "file.txt"

-- Get directory
dir = filepath.dir("/path/to/file.txt")  -- "/path/to"

-- Clean path (simplify)
clean = filepath.clean("/path//to/../file")  -- "/path/file"
```

#### regexp

Regular expression matching and manipulation.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/regexp"

L.PreloadModule("regexp", regexp.Loader)
```

**Lua API:**

```lua
local regexp = require("regexp")

-- Match pattern (boolean; raises on invalid pattern)
matches = regexp.match("^[a-z]+$", "hello")  -- true

-- Find first match (raises on invalid pattern)
match = regexp.find("([0-9]+)", "version 123 build 456")  -- "123"

-- Find all matches (raises on invalid pattern)
matches = regexp.find_all("([0-9]+)", "version 123 build 456", -1)
-- matches = {"123", "456"}

-- Replace occurrences (raises on invalid pattern)
result = regexp.replace("([0-9]+)", "version 123", "999")
-- result = "version 999"

result = regexp.replace_all("([0-9]+)", "version 123 build 456", "X")
-- result = "version X build X"

-- Split by pattern (raises on invalid pattern)
parts = regexp.split("\\s+", "one  two   three", -1)
-- parts = {"one", "two", "three"}
```

#### strings

String manipulation utilities.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/strings"

L.PreloadModule("strings", strings.Loader)
```

**Lua API:**

```lua
local strings = require("strings")

-- Prefix/suffix checking
has = strings.has_prefix("hello world", "hello")  -- true
has = strings.has_suffix("hello world", "world")  -- true

-- Trimming
trimmed = strings.trim("  hello  ", " ")  -- "hello"
trimmed = strings.trim_left("  hello  ", " ")  -- "hello  "
trimmed = strings.trim_right("  hello  ", " ")  -- "  hello"

-- Split and join
parts = strings.split("a,b,c", ",")  -- {"a", "b", "c"}
joined = strings.join({"a", "b", "c"}, ",")  -- "a,b,c"

-- Case conversion
upper = strings.to_upper("hello")  -- "HELLO"
lower = strings.to_lower("WORLD")  -- "world"

-- Search and count
has = strings.contains("hello world", "world")  -- true
count = strings.count("banana", "a")  -- 3

-- Replace
result = strings.replace("hello world", "world", "there", -1)  -- "hello there"

-- Whitespace trimming (the single most-missed function before this existed)
trimmed = strings.trim_space("  hello  ")  -- "hello"

-- Prefix/suffix removal (not just checking)
s = strings.trim_prefix("hello world", "hello ")  -- "world"
s = strings.trim_suffix("app.tar.gz", ".gz")      -- "app.tar"

-- Split on whitespace runs, with no empty entries
parts = strings.fields("  the quick  brown fox  ")  -- {"the", "quick", "brown", "fox"}

-- Repeat. Named rep, not repeat: "repeat" is a reserved word in Lua, so
-- strings.repeat(s, n) would be a syntax error at the call site.
s = strings.rep("ab", 3)  -- "ababab"

-- index/last_index are 1-based, returning 0 when absent -- NOT Go's
-- 0-based/-1 convention. This composes with string.sub directly.
pos = strings.index("key=value", "=")        -- 4
missing = strings.index("hello", "xyz")      -- 0
pos = strings.last_index("banana", "an")     -- 4

-- Case-insensitive equality (simple Unicode case-folding)
eq = strings.equal_fold("Hello", "HELLO")  -- true

-- Title-case the first rune of each whitespace-separated word. First-rune
-- only, not language-aware -- does not implement locale casing exceptions
-- (e.g. "of"/"the" staying lowercase in real title case).
s = strings.title("hello world")  -- "Hello World"

-- Cut: split around the first occurrence of a separator
before, after, found = strings.cut("app=nginx", "=")  -- "app", "nginx", true

-- Split with a limit (the remainder stays unsplit in the last element)
parts = strings.split_n("a,b,c,d", ",", 2)  -- {"a", "b,c,d"}
```

#### bit32

32-bit unsigned bitwise operations. gopher-lua implements Lua 5.1, which has
no bitwise operators and no `bit`/`bit32` library at all, so this fills a real
gap rather than duplicating something Lua already has.

The module is named `bit32`, not `bit`, and it is deliberate: Lua numbers are
float64, which can only represent integers exactly up to 2^53. A 64-bit
bitwise result would routinely exceed that and silently lose bits — worse
than no bitwise support at all. Restricting every operand and result to
`[0, 2^32)` keeps everything exactly representable. Results are always
unsigned, so `bit32.bnot(0)` is `4294967295`, never `-1`. **LuaJIT users
reaching for `require("bit")` should note the different name and the
unsigned result convention** — this is not that library.

Negative operands are accepted as two's complement (`-1` behaves as
`0xFFFFFFFF`); non-integer floats, `NaN`/`Inf`, and operands outside
`[-2^53, 2^53]` all raise. Shift counts must be non-negative (negative
raises); `n >= 32` clamps to `0` (or `0xFFFFFFFF` for `arshift` on a
negative value) rather than raising. Bit positions for `test`/`set`/`clear`
must be in `[0, 31]` — out of range raises, unlike shift counts.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/bit32"

L.PreloadModule("bit32", bit32.Loader)
```

**Lua API:**

```lua
local bit32 = require("bit32")

-- band/bor/bxor take a variadic tail
local rw = bit32.bor(0x4, 0x2)              -- 6
local masked = bit32.band(0xFF, 0x0F, 0x03) -- 3

-- Two's complement: -1 behaves as 0xFFFFFFFF
print(bit32.band(-1, 0xFF))                 -- 255
print(bit32.bnot(0))                        -- 4294967295, not -1

-- arshift sign-extends bit 31; rshift never does
print(bit32.rshift(0x80000000, 4))          -- 0x08000000
print(bit32.arshift(0x80000000, 4))         -- 0xF8000000

-- test/set/clear operate on individual bit positions [0, 31]
local mode = 0x1A4                          -- 0644 octal
if bit32.test(mode, 8) then
  print("owner can read")
end
```

#### strconv

Numeric parsing/formatting and Go-syntax string quoting, with errors that
actually say what went wrong.

Lua 5.1 already has `tonumber(s)` / `tonumber(s, base)`, so `parse_int` and
`parse_float` are admittedly thin wrappers over them — the value this module
adds is narrow but real:

- `tonumber` returns `nil` with no explanation. `strconv.parse_int` and
  `strconv.parse_float` **raise** with Go's own message
  (`strconv.ParseInt: parsing "12a": invalid syntax`), matching this
  library's fail-loud convention.
- `tonumber` silently hands back a float for an integer that doesn't fit.
  `parse_int`/`format_int` instead **raise** when a value's magnitude
  exceeds 2^53 — the largest integer a Lua number (a float64) can represent
  exactly — rather than silently rounding it.
- `format_int(n, base)` and `format_float(f, fmt, prec)` (with shortest
  round-trip formatting via `prec = -1`) have no Lua 5.1 equivalent at all.
- `quote`/`unquote` (Go-syntax string literals with escapes) have no
  equivalent in Lua 5.1.

There is no `itoa` — it is exactly `format_int(n, 10)`.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/strconv"

L.PreloadModule("strconv", strconv.Loader)
```

**Lua API:**

```lua
local strconv = require("strconv")

local mode = strconv.parse_int("644", 8)     -- 420
print(strconv.format_int(mode, 8))           -- "644"
print(strconv.format_float(1/3, "g", -1))    -- "0.3333333333333333"

-- tonumber("12a") would silently return nil; strconv raises with a reason.
local ok, err = pcall(strconv.atoi, "12a")
print(ok, err) -- false, ".../strconv.ParseInt: parsing "12a": invalid syntax"

-- parse_bool accepts Go's spelling set
print(strconv.parse_bool("TRUE"))            -- true

-- quote/unquote round-trip Go-syntax string literals
local q = strconv.quote("line one\nline two")
print(strconv.unquote(q) == "line one\nline two") -- true
```

#### text

Presentation/layout string operations that Go's standard library does not
provide: word wrap, indent/dedent, truncation with an ellipsis, and rune
padding.

The split from `strings` is deliberate: `strings` binds Go's `strings`
package one-to-one (if Go has it, it goes there); `text` is for operations Go
does not have at all.

**Every function here is rune-oriented, not display-width-oriented.**
"Width" and "length" always mean a count of Unicode code points, never bytes
and never terminal display cells. A CJK or emoji string will not visually
align in a monospace terminal even though these functions agree it is "N
runes wide" — getting real display width right needs a dependency
(`mattn/go-runewidth`) that is deliberately not taken for this module. This
is a stated limitation, not a bug.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/text"

L.PreloadModule("text", text.Loader)
```

**Lua API:**

```lua
local text = require("text")

-- Greedy word wrap. Existing "\n" are hard paragraph breaks. A word longer
-- than width is not split -- it overflows its own line. Raises if width < 1.
local wrapped = text.wrap("the quick brown fox jumps over the lazy dog", 20)

-- Prefix every line. A trailing empty line (s ends with "\n") is not
-- prefixed, so a trailing newline survives indenting unchanged.
print(text.indent(wrapped, "  | "))

-- Remove the common leading-whitespace prefix across all non-blank lines.
-- Tabs and spaces are compared literally, never expanded; blank lines are
-- normalized to empty rather than participating in the margin calculation.
local code = "    def f():\n        return 1\n"
print(text.dedent(code))  -- "def f():\n    return 1\n"

-- Truncate to a maximum rune width, appending an ellipsis when shortened.
print(text.truncate("sha256:0123456789abcdef", 12, "…"))  -- "sha256:012…"

-- Pad with a single rune. Never truncates -- an already-wider string is
-- returned unchanged.
print(text.pad_right("NAME", 20, " ") .. "STATUS")
print(text.pad_left("42", 5, "0"))  -- "00042"
```

`text.table` (an ASCII table formatter) was considered and deliberately cut
from this module — its design surface (column alignment, cell wrapping,
border styles) is a caller concern, not a stdlib concern.

## Features

- **Bidirectional Conversion**: Seamlessly convert Go structs to Lua tables and vice versa with full round-trip integrity
- **Automatic Stub Generation**: Generate Lua LSP annotations from Go types for IDE autocomplete
- **Lua Module System**: Create type-safe Lua modules with Go functions
- **IDE Support**: Full autocomplete and type checking in VSCode, Neovim, and other editors
- **Kubernetes Ready**: Built-in support for K8s API types and resource quantities
- **K8s Dynamic Client**: Full CRUD operations on any Kubernetes resource from Lua scripts
- **Type Safety**: Preserve complex types like timestamps, quantities, maps, and nested structures
- **Well Tested**: 79%+ code coverage with comprehensive unit and integration tests

## Use Cases

### Kubernetes Admission Controllers

Process and validate K8s resources in Lua scripts:

```go
// Load validation script
L.DoString(validationScript)

// Convert admission request to Lua
luaRequest, _ := translator.ToLua(L, admissionRequest)
L.SetGlobal("request", luaRequest)

// Run validation logic in Lua
L.DoString(`
    local pod = request.object
    if pod.spec.containers[1].resources.limits["memory"] == nil then
        reject("Memory limit required")
    end
`)
```

### Policy Engines

Define policies in Lua, enforce in Go:

```lua
local k8s = require("kubernetes")

function validate_pod(pod)
    -- Policy: Memory must be under 2GB
    local memBytes = k8s.parse_memory(pod.spec.containers[1].resources.limits["memory"])
    if memBytes > 2 * 1024 * 1024 * 1024 then
        return false, "Memory limit exceeds 2GB"
    end

    return true, nil
end
```

### Configuration Processing

Process complex config files with Lua logic:

```go
config := LoadConfig()
luaConfig, _ := translator.ToLua(L, config)
L.SetGlobal("config", luaConfig)

L.DoString(configProcessingScript)

var processedConfig Config
translator.FromLua(L, L.GetGlobal("result"), &processedConfig)
```

## Performance

Performance benchmarks demonstrate glua's efficiency for production use:

```
BenchmarkGoToLuaSimple-16              409915       3074 ns/op     4327 B/op       45 allocs/op
BenchmarkGoToLuaComplex-16              75747      15458 ns/op    23979 B/op      221 allocs/op
BenchmarkGoToLuaPod-16                  29530      40777 ns/op    58909 B/op      468 allocs/op
BenchmarkLuaToGoSimple-16              609123       2044 ns/op     1000 B/op       23 allocs/op
BenchmarkLuaToGoComplex-16             124522       9876 ns/op     4901 B/op      118 allocs/op
BenchmarkRoundTripSimple-16            202911       5640 ns/op     5330 B/op       68 allocs/op
BenchmarkRoundTripPod-16                30727      39612 ns/op    42924 B/op      391 allocs/op
BenchmarkLuaFieldAccess-16              88826      13688 ns/op    33937 B/op      114 allocs/op
BenchmarkLuaNestedFieldAccess-16        49378      24432 ns/op    37745 B/op      271 allocs/op
BenchmarkLuaArrayIteration-16           42882      28862 ns/op    36633 B/op      335 allocs/op
BenchmarkLuaMapIteration-16             72069      17139 ns/op    34833 B/op      125 allocs/op
BenchmarkLuaFieldModification-16        74764      16009 ns/op    34705 B/op      154 allocs/op
BenchmarkLuaComplexOperation-16         19200      76427 ns/op   195660 B/op      458 allocs/op
```

Key performance characteristics:

- **Simple conversions**: ~3µs Go→Lua, ~2µs Lua→Go
- **Kubernetes Pod**: ~40µs full round-trip conversion
- **Field access**: ~14µs for simple Lua operations
- **Production ready**: Suitable for request processing, admission controllers, policy evaluation

Run benchmarks yourself:

```bash
make bench          # View benchmark results
make bench-update   # Update benchmarks/README.md with latest results
```

## Testing

```bash
# Run ALL tests: unit + k8sclient integration (recommended)
make test

# Or just run make (default target runs all tests)
make

# Unit tests only
make test-unit

# K8s integration test only (requires Kind & kubectl)
make test-k8sclient

# Verbose per-package
make test-verbose

# Fast (no race detection)
make test-short
```

Coverage: 79%+ overall with comprehensive unit and integration tests

What's tested:

- Go ↔ Lua conversions in real Lua VMs (not just Go unit tests)
- Kubernetes module functions with actual K8s types
- K8s client CRUD operations with real Kind cluster
- Round-trip integrity (Go → Lua → Go preserves data)
- Stub generation from Go code
- Race detection enabled
- CI/CD across Go 1.21, 1.22, 1.23

## Example Application

The [example/](./example) directory contains a complete working demo showing all features.

```bash
# From repo root
make example && ./bin/example

# Or
cd example && go run .
```

Features demonstrated:

- Go → Lua conversion (Pod struct to Lua table)
- Lua script execution with kubernetes module
- Parsing timestamps, CPU, and memory quantities
- Lua → Go conversion (table back to Pod struct)
- Round-trip integrity verification
- Stub generation for IDE autocomplete
- Error handling

To get autocomplete in the example:

1. Run `make gen-stubs` from repo root
2. Run `go run .` from example/ directory
3. Open `script.lua` in your IDE - autocomplete works

## Troubleshooting

### Autocomplete doesn't work

1. Run `make gen-stubs` to generate module stubs
2. Check `.luarc.json` or `.vscode/settings.json` includes `"library"` directory
3. Verify `library/kubernetes.gen.lua` exists and starts with `---@meta`
4. Restart LSP: `:LspRestart` (Neovim) or reload window (VSCode)

### Module not found error

Ensure `L.PreloadModule("mymodule", mymodule.Loader)` is called before `L.DoString()`

### Round-trip data mismatch

Ensure all struct fields are exported (capitalized) and JSON-marshallable

## Contributing

Contributions welcome! Please:

1. Ensure all tests pass (`go test -cover -race ./...`)
2. Add tests for new functionality
3. Follow existing code style (gofmt)
4. Update documentation
5. Add examples for new features

## License

MIT License - see LICENSE file for details

## Credits

Built on top of [gopher-lua](https://github.com/yuin/gopher-lua) by Yusuke Inuzuka.

Kubernetes API types from [k8s.io/api](https://github.com/kubernetes/api).
