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

Or manually add to `go.mod` and run `go mod tidy`:

```go
require github.com/thomas-maurice/glua v0.0.12 // or latest version
```

### Installing Lua Stubs for IDE Autocomplete

Download the latest Lua stubs for IDE autocomplete support:

```bash
# Download and extract to your project
VERSION=v0.0.12  # Replace with the latest version
curl -sL https://github.com/thomas-maurice/glua/releases/download/${VERSION}/glua-stubs_${VERSION}.tar.gz | tar xz

# This extracts to library/*.gen.lua
# Configure your IDE to recognize the library/ directory
```

**VS Code Setup** (with Lua extension):

Create or update `.vscode/settings.json`:

```json
{
  "Lua.workspace.library": ["library"]
}
```

Now you'll get autocomplete for all glua modules in your Lua scripts!

### Installing stubgen Binary (Optional)

If you want to generate stubs for your own modules:

```bash
# Linux/macOS
VERSION=v0.0.12  # Replace with the latest version
curl -sL https://github.com/thomas-maurice/glua/releases/download/${VERSION}/stubgen_${VERSION}_$(uname -s | tr '[:upper:]' '[:lower:]')_$(uname -m | sed 's/x86_64/amd64/').tar.gz | tar xz
sudo mv stubgen /usr/local/bin/

# Verify installation
stubgen --help
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

        -- Parse Kubernetes quantities
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

- `bin/stubgen` - Generates Lua LSP stubs for IDE autocomplete
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
local memBytes, err = k8s.parse_memory("256Mi")  -- 268435456
local memBytes2 = k8s.parse_memory("1Gi")         -- 1073741824

-- Parse CPU quantities (returns millicores)
local cpuMillis, err = k8s.parse_cpu("100m")     -- 100
local cpuMillis2 = k8s.parse_cpu("1.5")          -- 1500

-- Parse timestamps (returns Unix timestamp)
local timestamp, err = k8s.parse_time("2025-10-03T16:39:00Z")  -- 1759509540

-- Format timestamps (Unix timestamp → RFC3339 string)
local timeStr, err = k8s.format_time(1759509540)  -- "2025-10-03T16:39:00Z"

-- All functions return (value, error) tuple
local bytes, err = k8s.parse_memory("invalid")
if err then
    print("Parse error: " .. err)
end
```

Example use case - process pod resource limits/requests:

```lua
local k8s = require("kubernetes")
local pod = myPod

for i, container in ipairs(pod.spec.containers) do
    print("Container: " .. container.name)

    -- Parse memory limit
    if container.resources.limits["memory"] then
        local memBytes = k8s.parse_memory(container.resources.limits["memory"])
        print(string.format("  Memory limit: %.2f MB", memBytes / (1024 * 1024)))
    end

    -- Parse CPU limit
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
local client = require("k8sclient")

-- Define GVK (Group/Version/Kind)
local pod_gvk = {group = "", version = "v1", kind = "Pod"}

-- Create a Pod
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
local created, err = client.create(pod)

-- Get a resource
local fetched, err = client.get(pod_gvk, "default", "nginx")

-- Update a resource
fetched.metadata.labels = {app = "web"}
local updated, err = client.update(fetched)

-- List resources
local pods, err = client.list(pod_gvk, "default")
for i, pod in ipairs(pods) do
    print(pod.metadata.name)
end

-- Delete a resource
local err = client.delete(pod_gvk, "default", "nginx")
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

local bytes, err = k8s.parse_memory("256Mi")
if err then
    print("Error: " .. err)
    return
end

print("Parsed successfully: " .. bytes)
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

### Step 1: Create Module

You can create two types of Lua modules: **function-based modules** (simple) and **UserData-based modules** (for stateful objects).

#### Option A: Simple Function-Based Module

Create `pkg/modules/mymodule/mymodule.go`:

```go
package mymodule

import (
    lua "github.com/yuin/gopher-lua"
)

// Loader: creates the mymodule Lua module
//
// @luamodule mymodule
func Loader(L *lua.LState) int {
    mod := L.SetFuncs(L.NewTable(), exports)
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "greet": greet,
    "add":   add,
}

// greet: returns a personalized greeting
//
// @luafunc greet
// @luaparam name string The name to greet
// @luareturn string The greeting message
func greet(L *lua.LState) int {
    name := L.CheckString(1)
    L.Push(lua.LString("Hello, " + name + "!"))
    return 1
}

// add: adds two numbers
//
// @luafunc add
// @luaparam a number First number
// @luaparam b number Second number
// @luareturn number Sum of a and b
// @luareturn string|nil Error message if any
func add(L *lua.LState) int {
    a := L.CheckNumber(1)
    b := L.CheckNumber(2)
    L.Push(lua.LNumber(a + b))
    L.Push(lua.LNil)
    return 2
}
```

#### Option B: UserData-Based Module (Stateful Objects)

For modules that need to maintain state or provide object-oriented APIs, use UserData. This example shows how to create a Logger object with methods (like the built-in `log` module):

Create `pkg/modules/counter/counter.go`:

```go
package counter

import (
    lua "github.com/yuin/gopher-lua"
)

const counterTypeName = "counter.Counter"

// Counter: a simple counter object
type Counter struct {
    value int
}

// Loader: creates the counter Lua module
//
// @luamodule counter
func Loader(L *lua.LState) int {
    // Register the Counter type with methods
    registerCounterType(L)

    // Create module table with functions
    mod := L.SetFuncs(L.NewTable(), exports)
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "new": newCounter,
}

// registerCounterType: registers the Counter UserData type with its metatable
func registerCounterType(L *lua.LState) {
    mt := L.NewTypeMetatable(counterTypeName)
    L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), counterMethods))
}

// counterMethods: map of Counter object methods
var counterMethods = map[string]lua.LGFunction{
    "increment": counterIncrement,
    "decrement": counterDecrement,
    "get":       counterGet,
    "reset":     counterReset,
}

// wrapCounter: wraps a Counter in UserData for use in Lua
func wrapCounter(L *lua.LState, c *Counter) *lua.LUserData {
    ud := L.NewUserData()
    ud.Value = c
    L.SetMetatable(ud, L.GetTypeMetatable(counterTypeName))
    return ud
}

// checkCounter: extracts Counter from UserData
func checkCounter(L *lua.LState, n int) *Counter {
    ud := L.CheckUserData(n)
    if c, ok := ud.Value.(*Counter); ok {
        return c
    }
    L.ArgError(n, "Counter expected")
    return nil
}

// newCounter: creates a new Counter object
//
// @luafunc new
// @luaparam initialValue number Optional initial value (default: 0)
// @luareturn counter.Counter counter.Counter A new counter object
func newCounter(L *lua.LState) int {
    initialValue := 0
    if L.GetTop() >= 1 {
        initialValue = int(L.CheckNumber(1))
    }

    counter := &Counter{value: initialValue}
    ud := wrapCounter(L, counter)
    L.Push(ud)
    return 1
}

// counterIncrement: increments the counter
//
// @luamethod counter.Counter increment
// @luaparam self counter.Counter The counter object
// @luaparam amount number Optional amount to add (default: 1)
func counterIncrement(L *lua.LState) int {
    c := checkCounter(L, 1)
    amount := 1
    if L.GetTop() >= 2 {
        amount = int(L.CheckNumber(2))
    }
    c.value += amount
    return 0
}

// counterDecrement: decrements the counter
//
// @luamethod counter.Counter decrement
// @luaparam self counter.Counter The counter object
// @luaparam amount number Optional amount to subtract (default: 1)
func counterDecrement(L *lua.LState) int {
    c := checkCounter(L, 1)
    amount := 1
    if L.GetTop() >= 2 {
        amount = int(L.CheckNumber(2))
    }
    c.value -= amount
    return 0
}

// counterGet: gets the current counter value
//
// @luamethod counter.Counter get
// @luaparam self counter.Counter The counter object
// @luareturn number The current counter value
func counterGet(L *lua.LState) int {
    c := checkCounter(L, 1)
    L.Push(lua.LNumber(c.value))
    return 1
}

// counterReset: resets the counter to zero
//
// @luamethod counter.Counter reset
// @luaparam self counter.Counter The counter object
func counterReset(L *lua.LState) int {
    c := checkCounter(L, 1)
    c.value = 0
    return 0
}
```

**Usage in Lua:**

```lua
local counter = require("counter")

-- Create counter objects
local c1 = counter.new()
local c2 = counter.new(10)

-- Use object methods with : notation
c1:increment()
c1:increment(5)
print(c1:get())  -- 6

c2:decrement()
print(c2:get())  -- 9

c1:reset()
print(c1:get())  -- 0
```

**Key concepts for UserData objects:**

1. **Type Name**: Unique identifier for your UserData type (e.g., `"counter.Counter"`)
2. **Metatable**: Defines methods available on your object via `__index`
3. **Wrapper Function**: `wrapCounter()` creates UserData from Go struct
4. **Checker Function**: `checkCounter()` extracts Go struct from UserData
5. **Method Signature**: Methods use `:` notation in Lua, which implicitly passes `self` as first argument
6. **Annotations**: Use `@luamethod ModuleName.ClassName methodName` for methods vs `@luafunc` for module functions

**See also:** The built-in `log` module (`pkg/modules/log/`) is a complete real-world example of UserData objects with proper stub generation.

### Step 2: Generate Stubs with stubgen

The stubgen tool scans your Go code for special annotations and generates Lua LSP stubs for IDE autocomplete.

**Run stubgen:**

```bash
# Scan pkg/modules directory and generate stubs in library/
make stubgen

# Or manually:
go run ./cmd/stubgen -dir pkg/modules -output-dir library

# For a single combined file:
go run ./cmd/stubgen -dir pkg/modules -output mymodules.gen.lua
```

**What stubgen looks for:**

The tool scans for these comment annotations in your Go code:

### Module and Function Annotations

- **`@luamodule <name>`** - Marks the Loader function (required for each module)
  - Must be directly above the `Loader` function
  - Example: `@luamodule mymodule`

- **`@luafunc <name>`** - Defines a module-level function
  - For functions like `mymodule.greet()`
  - Example: `@luafunc greet`

- **`@luamethod <ClassName> <methodName>`** - Defines a method on a UserData object
  - For object methods like `logger:info()`
  - Class name should be namespaced (e.g., `log.Logger`, `counter.Counter`)
  - Example: `@luamethod log.Logger info`

- **`@luaclass <ClassName>`** - Defines a data structure type with manual field annotations
  - For table-like data structures where you want explicit control over field documentation
  - Use with `@luafield` to define each field with custom descriptions
  - Class name can be simple (`GVKMatcher`) or namespaced (`mymodule.Config`)
  - Example: `@luaclass GVKMatcher`
  - **TIP**: For complex Go structs with many nested fields, Type Registry offers automatic field discovery

- **`@luafield <fieldName> <type> [description]`** - Defines a field in a `@luaclass`
  - Must be used after `@luaclass` annotation
  - `<fieldName>`: Name of the field
  - `<type>`: Lua type of the field (`string`, `number`, `table`, etc.)
  - `[description]`: Optional human-readable description
  - Example: `@luafield group string The API group`

- **`@luaconst <NAME> <type> [description]`** - Defines a module constant
  - For constants like `k8sclient.POD`, `k8sclient.DEPLOYMENT`
  - `<NAME>`: Constant name (typically UPPERCASE)
  - `<type>`: Lua type of the constant (`table`, `string`, `number`, etc.)
  - `[description]`: Optional human-readable description
  - Can be placed anywhere in the file (not just above functions)
  - Example: `@luaconst POD table Pod GVK constant {group="", version="v1", kind="Pod"}`

### Parameter Annotations

- **`@luaparam <name> <type> [description]`** - Defines a function/method parameter
  - `<name>`: Parameter name (use `self` for the implicit object parameter in methods)
  - `<type>`: Lua type (`string`, `number`, `boolean`, `table`, `any`, or custom type like `log.Logger`)
  - `[description]`: Optional human-readable description
  - Examples:
    - `@luaparam msg string The message to log`
    - `@luaparam count number`
    - `@luaparam ... any Optional key-value pairs` (for variadic parameters)
    - `@luaparam self log.Logger The logger object` (for methods)

### Return Value Annotations

- **`@luareturn <type> [description]`** - Defines a return value
  - Can have multiple `@luareturn` annotations for multiple return values
  - `<type>`: Lua type or custom type
  - `[description]`: Optional human-readable description
  - Type can include union types: `string|nil` for optional returns
  - Examples:
    - `@luareturn string The greeting message`
    - `@luareturn log.Logger log.Logger A new logger with additional fields`
    - `@luareturn err string|nil Error message if any`

### Annotation Placement

Annotations must be in **Go-style comments** (`//`) directly above the function:

```go
// functionName: brief description of what it does
//
// @luafunc functionName
// @luaparam param1 string Description of param1
// @luaparam param2 number Description of param2
// @luareturn result string Description of return value
// @luareturn err string|nil Error message if operation failed
func functionName(L *lua.LState) int { ... }
```

For methods:

```go
// methodName: brief description of what it does
//
// @luamethod ClassName methodName
// @luaparam self ClassName The object instance
// @luaparam param1 string Description of param1
// @luareturn result any Description of return value
func methodName(L *lua.LState) int { ... }
```

For standalone classes (data structures):

```go
// Loader: creates the mymodule Lua module
//
// @luamodule mymodule
//
// @luaclass GVKMatcher
// @luafield group string The API group
// @luafield version string The API version
// @luafield kind string The resource kind
func Loader(L *lua.LState) int { ... }
```

For constants (can appear anywhere in the file):

```go
// @luaconst POD table Pod GVK constant {group="", version="v1", kind="Pod"}

// @luaconst DEPLOYMENT table Deployment GVK constant {group="apps", version="v1", kind="Deployment"}

func setupConstants(L *lua.LState, mod *lua.LTable) {
    L.SetField(mod, "POD", createGVKTable(L, "", "v1", "Pod"))
    L.SetField(mod, "DEPLOYMENT", createGVKTable(L, "apps", "v1", "Deployment"))
}
```

**Example from our code above:**

```go
// Module-level function annotation
// @luamodule mymodule    <- Tells stubgen this is a Lua module
func Loader(L *lua.LState) int { ... }

// @luafunc greet         <- Function name in Lua
// @luaparam name string The name to greet    <- Parameter with type and description
// @luareturn string The greeting message     <- Return type and description
func greet(L *lua.LState) int { ... }

// UserData method annotation
// @luamethod counter.Counter increment    <- Method on Counter class
// @luaparam self counter.Counter The counter object
// @luaparam amount number Optional amount to add
func counterIncrement(L *lua.LState) int { ... }
```

**Generated output for simple module** (`library/mymodule.gen.lua`):

```lua
---@meta mymodule

---@class mymodule
local mymodule = {}

---@param name string The name to greet
---@return string The greeting message
function mymodule.greet(name) end

---@param a number First number
---@param b number Second number
---@return number Sum of a and b
---@return string|nil Error message if any
function mymodule.add(a, b) end

return mymodule
```

**Generated output for UserData module** (`library/counter.gen.lua`):

```lua
---@meta counter

---@class counter.Counter
local Counter = {}

---@param amount number Optional amount to add (default: 1)
function Counter:increment(amount) end

---@param amount number Optional amount to subtract (default: 1)
function Counter:decrement(amount) end

---@return number The current counter value
function Counter:get() end

function Counter:reset() end

---@class counter
---@field Counter counter.Counter
local counter = {}

---@param initialValue number Optional initial value (default: 0)
---@return counter.Counter A new counter object
function counter.new(initialValue) end

counter.Counter = Counter

return counter
```

Note how the UserData class (`Counter`) is defined first with its methods, then the module (`counter`) is defined with the class as a field, and finally they're linked together with `counter.Counter = Counter`. This structure enables proper IDE autocomplete.

### Choosing Between @luaclass and Type Registry

Both approaches are fully supported and actively used. Choose based on your needs:

**Use `@luaclass` + `@luafield` annotations when:**

- You want explicit control over documentation and field descriptions
- Defining simple table-like data structures (e.g., configuration objects)
- You need custom field descriptions that differ from Go struct tags
- Working with interface types or non-struct data
- Example: `GVKMatcher` with detailed field descriptions

**Use Type Registry (`typeRegistry.Register()`) when:**

- You want automatic field discovery from Go struct definitions
- Working with complex Go structs with many nested fields
- Using third-party types (e.g., Kubernetes API types)
- You have many similar types that need consistent documentation
- Field names and types from JSON tags are sufficient
- Example: Kubernetes resources, complex configuration structs

**Example comparison:**

```go
// @luaclass approach - simple, manual
// @luaclass GVKMatcher
// @luafield group string API group
// @luafield version string API version
// @luafield kind string Resource kind

// Type Registry approach - automatic, comprehensive
typeRegistry.Register(corev1.Pod{})        // Auto-discovers ALL fields
typeRegistry.Register(corev1.Service{})    // Including nested types
typeRegistry.Register(corev1.ConfigMap{})  // With proper type references
```

The Type Registry automatically processes nested types, creates proper type hierarchies, and handles complex struct relationships. See the [Type Registry section](#type-registry-and-lsp-stub-generation) for details.

**How stubgen works:**

1. Stubgen scans all `.go` files in the specified directory
2. Finds functions with `@luamodule` annotation (these are module Loaders)
3. Finds functions with `@luafunc` annotation (these are exported Lua functions)
4. Extracts `@luaparam` and `@luareturn` annotations for each function
5. Generates EmmyLua-compatible annotation files that LSP servers understand
6. Outputs one `.gen.lua` file per module (or a single combined file)

**Verification:**

```bash
# Check generated files
ls library/
# Output: json.gen.lua  kubernetes.gen.lua  mymodule.gen.lua  spew.gen.lua

# View generated stub
cat library/mymodule.gen.lua
```

### Step 3: Register and Use

```go
package main

import (
    "your-project/pkg/modules/mymodule"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // Register module
    L.PreloadModule("mymodule", mymodule.Loader)

    // Use in Lua
    L.DoString(`
        local m = require("mymodule")
        print(m.greet("World"))  -- Hello, World!
        print(m.add(5, 3))       -- 8
    `)
}
```

### Stubgen Tool

The `stubgen` command generates Lua LSP stubs from your Go module code.

Usage:

```bash
# Generate stubs for all modules (recommended)
make stubgen

# Or manually:
go run ./cmd/stubgen -dir pkg/modules -output-dir library

# Single combined file:
go run ./cmd/stubgen -dir pkg/modules -output stubs.lua
```

Options:

- `-dir`: Directory to scan for Go modules (default: ".")
- `-output-dir`: Generate per-module files in this directory (recommended for LSP)
- `-output`: Generate single combined file (default: "module_stubs.gen.lua")

What it does:

Scans Go files for these annotations:

- `@luamodule <name>` - Marks module Loader function
- `@luafunc <name>` - Module-level function
- `@luamethod <ClassName> <methodName>` - Method on UserData object (e.g., `@luamethod log.Logger info`)
- `@luaconst <NAME> <type> [description]` - Module constant (e.g., `@luaconst POD table`)
  - Can appear anywhere in the file
  - Generates `---@type` annotations
- `@luaparam <name> <type> [description]` - Function/method parameter
  - Supports variadic: `@luaparam ... any`
  - For methods, include: `@luaparam self ClassName`
- `@luareturn <type> [description]` - Return value
  - Supports multiple returns
  - Supports union types: `string|nil`

Generates Lua LSP annotation files that IDEs use for autocomplete. For UserData objects, it automatically generates the proper class structure with methods and exports the class as a module field.

**See detailed annotation reference in the "Creating Custom Lua Modules" section above.**

## IDE Setup

This section explains how to enable autocomplete for your Lua scripts.

### VSCode Setup

1. Install Lua Language Server extension: [Lua](https://marketplace.visualstudio.com/items?itemName=sumneko.lua)

2. Generate stubs:

```bash
# Generate module stubs (for kubernetes, custom modules)
make stubgen  # Creates library/kubernetes.gen.lua, etc.

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
make stubgen
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
bytes, err = k8s.parse_memory(quantity)

-- Parse CPU: "100m" → 100 (millicores)
millis, err = k8s.parse_cpu(quantity)

-- Parse duration: "5m" → 300 (seconds)
seconds, err = k8s.parse_duration(duration)

-- Parse time: "2025-10-03T16:39:00Z" → 1759509540 (Unix timestamp)
timestamp, err = k8s.parse_time(timestr)

-- Format time: 1759509540 → "2025-10-03T16:39:00Z"
timestr, err = k8s.format_time(timestamp)

-- Format duration: 300 → "5m0s"
duration, err = k8s.format_duration(seconds)

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

All functions return `(value, error)` tuples (except boolean helpers).

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
local client = require("k8sclient")

-- GVK constants (predefined)
client.POD          -- {group="", version="v1", kind="Pod"}
client.DEPLOYMENT   -- {group="apps", version="v1", kind="Deployment"}
client.SERVICE      -- {group="", version="v1", kind="Service"}
client.CONFIGMAP    -- {group="", version="v1", kind="ConfigMap"}
-- ... and many more

-- Create a resource
created, err = client.create(resource_table)

-- Get a resource
resource, err = client.get(gvk, namespace, name)

-- Update a resource
updated, err = client.update(resource_table)

-- List resources
resources, err = client.list(gvk, namespace)

-- Delete a resource
err = client.delete(gvk, namespace, name)
```

**Example:**

```lua
local client = require("k8sclient")

-- Create a ConfigMap
local cm = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {name = "my-config", namespace = "default"},
    data = {key = "value"}
}
local created, err = client.create(cm)

-- Get it back
local fetched, err = client.get(client.CONFIGMAP, "default", "my-config")

-- Update it
fetched.data.newkey = "newvalue"
local updated, err = client.update(fetched)

-- Delete it
local err = client.delete(client.CONFIGMAP, "default", "my-config")
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

-- Parse JSON string to Lua table
table, err = json.parse('{"name":"John","age":30}')

-- Stringify Lua table to JSON
jsonstr, err = json.stringify({name="John", age=30})
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

-- Parse YAML string to Lua table
table, err = yaml.parse("name: John\nage: 30")

-- Stringify Lua table to YAML
yamlstr, err = yaml.stringify({name="John", age=30})
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

HTTP client for making requests.

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/http"

L.PreloadModule("http", http.Loader)
```

**Lua API:**

```lua
local http = require("http")

-- GET request
response, err = http.get("https://api.example.com/data")
-- response = {status=200, body="...", headers={...}}

-- POST request
response, err = http.post("https://api.example.com/data", {
    body = '{"key":"value"}',
    headers = {["Content-Type"] = "application/json"}
})

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

-- Render template with data
result, err = template.render("Hello {{.name}}, you are {{.age}} years old",
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

-- Read file
content, err = fs.read_file("/path/to/file.txt")

-- Write file
err = fs.write_file("/path/to/file.txt", "content")

-- Check existence
exists = fs.exists("/path/to/file")

-- Create directory
err = fs.mkdir("/path/to/dir")
err = fs.mkdir_all("/path/to/nested/dir")

-- Remove
err = fs.remove("/path/to/file")
err = fs.remove_all("/path/to/dir")

-- List directory
files, err = fs.list("/path/to/dir")

-- Get file info
info, err = fs.stat("/path/to/file")
-- info = {size=1234, mode="0644", mod_time=1234567890, is_dir=false}
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

-- Parse date string
timestamp, err = time.parse("2006-01-02 15:04:05", "2025-10-21 14:30:00")

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
decoded, err = base64.decode(encoded)

-- Hex
encoded = hex.encode("hello")
decoded, err = hex.decode(encoded)

-- Hash
md5 = hash.md5("hello")
sha1 = hash.sha1("hello")
sha256 = hash.sha256("hello")
sha512 = hash.sha512("hello")
```

#### log

Structured logging with fields (similar to logrus).

**Load in Go:**

```go
import "github.com/thomas-maurice/glua/pkg/modules/log"

L.PreloadModule("log", log.Loader)
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

-- Logger with preset fields
logger = log.with_fields({component="api", version="1.0"})
logger:info("Request received", {path="/api/users"})
-- Output includes: component=api version=1.0 path=/api/users

-- Set log level
log.set_level("debug")  -- "debug", "info", "warn", "error"

-- Set output format
log.set_format("json")  -- "json" or "text"
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

-- Get absolute path
abspath, err = filepath.abs("../relative/path")

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

-- Match pattern (boolean)
matches = regexp.match("^[a-z]+$", "hello")  -- true

-- Find first match
match, err = regexp.find("([0-9]+)", "version 123 build 456")  -- "123"

-- Find all matches
matches, err = regexp.find_all("([0-9]+)", "version 123 build 456", -1)
-- matches = {"123", "456"}

-- Replace first occurrence
result, err = regexp.replace("([0-9]+)", "version 123", "999", 1)
-- result = "version 999"

-- Replace all occurrences
result, err = regexp.replace_all("([0-9]+)", "version 123 build 456", "X")
-- result = "version X build X"

-- Split by pattern
parts, err = regexp.split("\\s+", "one  two   three", -1)
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
```

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

1. Run `make stubgen` from repo root
2. Run `go run .` from example/ directory
3. Open `script.lua` in your IDE - autocomplete works

## Troubleshooting

### Autocomplete doesn't work

1. Run `make stubgen` to generate module stubs
2. Check `.luarc.json` or `.vscode/settings.json` includes `"library"` directory
3. Verify `library/kubernetes.gen.lua` exists and starts with `---@meta`
4. Restart LSP: `:LspRestart` (Neovim) or reload window (VSCode)

### Module not found error

Ensure `L.PreloadModule("mymodule", mymodule.Loader)` is called before `L.DoString()`

### Stubgen finds no modules

Check `@luamodule` annotation is directly above Loader function in Go code

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
