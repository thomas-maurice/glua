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
- [Testing](#testing)
- [Contributing](#contributing)

## Installation

```bash
go get github.com/thomas-maurice/glua
```

Requirements:
- Go 1.21 or later
- For Kubernetes support: `k8s.io/api` and `k8s.io/apimachinery`

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

- `@luamodule <name>` - Marks the Loader function (required)
- `@luafunc <name>` - Exported Lua function name
- `@luaparam <name> <type> <description>` - Function parameter
- `@luareturn <type> <description>` - Return value

**Example from our code above:**

```go
// @luamodule mymodule    <- Tells stubgen this is a Lua module
func Loader(L *lua.LState) int { ... }

// @luafunc greet         <- Function name in Lua
// @luaparam name string The name to greet    <- Parameter with type and description
// @luareturn string The greeting message     <- Return type and description
func greet(L *lua.LState) int { ... }
```

**Generated output** (`library/mymodule.gen.lua`):

```lua
---@meta

---@class mymodule
local mymodule = {}

--- greet: returns a personalized greeting
---@param name string The name to greet
---@return string The greeting message
function mymodule.greet(name) end

--- add: adds two numbers
---@param a number First number
---@param b number Second number
---@return number Sum of a and b
---@return string|nil Error message if any
function mymodule.add(a, b) end

return mymodule
```

**How it works:**

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
- `@luafunc <name>` - Exported function
- `@luaparam <name> <type> <description>` - Function parameter
- `@luareturn <type> <description>` - Return value

Generates Lua LSP annotation files that IDEs use for autocomplete.

## IDE Setup

This section explains how to enable autocomplete for your Lua scripts.

### VSCode Setup

1. Install Lua Language Server extension: [Lua](https://marketplace.visualstudio.com/items?itemName=sumneko.lua)

2. Generate stubs:
```bash
# Generate module stubs (for kubernetes, custom modules)
make stubgen  # Creates library/kubernetes.gen.lua, etc.

# Generate type stubs (run your app that uses TypeRegistry)
go run main.go  # Creates annotations.gen.lua
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
go run main.go  # If using TypeRegistry
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

### Translator

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

### TypeRegistry

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

### Kubernetes Module (Lua API)

```lua
local k8s = require("kubernetes")

-- Parse memory: "256Mi" → 268435456 (bytes)
bytes, err = k8s.parse_memory(quantity)

-- Parse CPU: "100m" → 100 (millicores)
millis, err = k8s.parse_cpu(quantity)

-- Parse time: "2025-10-03T16:39:00Z" → 1759509540 (Unix timestamp)
timestamp, err = k8s.parse_time(timestr)

-- Format time: 1759509540 → "2025-10-03T16:39:00Z"
timestr, err = k8s.format_time(timestamp)
```

All functions return `(value, error)` tuples.

## Features

- **Bidirectional Conversion**: Seamlessly convert Go structs to Lua tables and vice versa with full round-trip integrity
- **Automatic Stub Generation**: Generate Lua LSP annotations from Go types for IDE autocomplete
- **Lua Module System**: Create type-safe Lua modules with Go functions
- **IDE Support**: Full autocomplete and type checking in VSCode, Neovim, and other editors
- **Kubernetes Ready**: Built-in support for K8s API types and resource quantities
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

## Testing

```bash
# Run all tests (recommended)
make test

# Verbose per-package
make test-verbose

# Fast (no race detection)
make test-short
```

Coverage: 79%+ overall with comprehensive unit and integration tests

What's tested:
- Go ↔ Lua conversions in real Lua VMs (not just Go unit tests)
- Kubernetes module functions with actual K8s types
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
cd example && go run main.go
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
2. Run `go run main.go` from example/ directory
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
