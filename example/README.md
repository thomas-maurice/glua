# glua Example Application

This example demonstrates **all core features** of the glua library in a single comprehensive demo.

## What It Demonstrates

The example application showcases:

1. **Type Registry & LSP Stub Generation** - Generate Lua type annotations from Go types for IDE autocomplete
2. **Go → Lua Conversion** - Convert Kubernetes Pod structs to Lua tables seamlessly
3. **Lua Script Execution** - Run comprehensive Lua scripts with full Pod access
4. **Running Example Scripts** - Execute 3 example scripts showcasing real-world use cases:
   - Resource analysis with kubernetes module
   - Policy validation and enforcement
   - Data modification in Lua
5. **Lua → Go Conversion** - Convert modified Lua tables back to Go structs
6. **Round-Trip Verification** - Verify perfect data integrity preservation
7. **Feature Summary** - Complete overview of all library capabilities

## Quick Start

### Option 1: Run the Main Demo

```bash
# From the example directory
cd example
go run main.go
```

Or from repository root:
```bash
# Build and run
make example
cd example && ../bin/example
```

### Option 2: Run Individual Example Scripts

The script runner is now a standalone command:

```bash
# From repository root
go run ./cmd/run-script scripts/01_basic_pod_info.lua
go run ./cmd/run-script scripts/02_resource_limits.lua
go run ./cmd/run-script scripts/03_policy_validation.lua
```

See [scripts/README.md](scripts/README.md) for all 7 available scripts.

## Example Output

```
╔════════════════════════════════════════════════════════════╗
║         glua - Go ↔ Lua Translator Demo                   ║
╚════════════════════════════════════════════════════════════╝

[1/7] Generating LSP type annotations...
    ✓ Generated annotations.gen.lua for IDE autocomplete

[2/7] Converting Go Pod struct to Lua table...
    ✓ Converted Go struct to Lua table

[3/7] Running main demonstration script...
    (Full Lua script output showing pod processing)
    ✓ Lua script executed successfully
    ✓ Lua parsed CPU: 100 millicores
    ✓ Lua parsed Memory: 104857600 bytes (100.00 MB)
    ✓ Lua parsed Timestamp: 1759509540 (Unix time)

[4/7] Running example scripts to showcase features...
    [1/3] Parse CPU/memory with kubernetes module...
        ✓ Success
    [2/3] Enforce organizational policies...
        ✓ Success
    [3/3] Modify pod data in Lua...
        ✓ Success

[5/7] Converting modified Lua table back to Go...
    ✓ Converted Lua table back to Go Pod struct

[6/7] Verifying data integrity (round-trip test)...
    ✓ Timestamp preserved: 2025-10-03 16:39:00 UTC
    ✓ CPU limit preserved: 100m (100 millicores)
    ✓ Memory limit preserved: 100Mi (104857600 bytes)
    ✓ Full JSON round-trip verified (3722 bytes)

[7/7] Summary of glua capabilities demonstrated:
    ✓ Type Registry - Generate LSP stubs for IDE autocomplete
    ✓ Go → Lua - Convert any Go struct to Lua table
    ✓ Lua Modules - kubernetes module functions
    ✓ Lua Execution - Run complex scripts with full Pod access
    ✓ Lua → Go - Convert Lua tables back to Go structs
    ✓ Round-trip Integrity - Perfect data preservation
    ✓ Example Scripts - 7 scripts showcasing real-world use cases

╔════════════════════════════════════════════════════════════╗
║                  ALL TESTS PASSED ✓                        ║
╚════════════════════════════════════════════════════════════╝

Next steps:
  • Explore example scripts in scripts/ directory
  • Run individual scripts: go run ./cmd/run-script scripts/01_basic_pod_info.lua
  • Open scripts in your IDE for full autocomplete support
  • See EXAMPLES.md for detailed documentation

⏱️  Total execution time: ~4ms
```

## Project Structure

```
glua/
├── example/
│   ├── main.go              # Main demo application (runs all features)
│   ├── script.lua           # Main demonstration script
│   ├── scripts/             # 7 example scripts showing different use cases
│   │   ├── 01_basic_pod_info.lua
│   │   ├── 02_resource_limits.lua
│   │   ├── 03_policy_validation.lua
│   │   ├── 04_environment_vars.lua
│   │   ├── 05_timestamp_operations.lua
│   │   ├── 06_multi_container_analysis.lua
│   │   └── 07_json_export.lua
│   ├── sample/              # Sample data (realistic Pod object)
│   ├── .luarc.json          # Lua LSP configuration for IDE
│   └── EXAMPLES.md          # Detailed guide for all example scripts
├── cmd/
│   ├── stubgen/             # Stub generator tool
│   └── run-script/          # Script runner for individual examples
├── pkg/
│   ├── glua/                # Core library
│   ├── modules/             # Lua modules (kubernetes, etc.)
│   └── stubgen/             # Stub generation logic
└── library/                 # Generated module stubs (kubernetes.lua)
```

### Generated Files (not committed)

These files are generated automatically when you run the example or stubgen:

- **example/annotations.gen.lua** - Type definitions for Pod and related types (~30KB)
- **library/kubernetes.lua** - Kubernetes module stubs for IDE autocomplete (~1.7KB)

## Generating IDE Stubs

glua provides two types of stubs for IDE autocomplete:

### 1. Module Stubs (for `kubernetes` and custom modules)

Generate stubs from Go module code using the `stubgen` tool:

```bash
# From repository root
make stubgen

# Or manually
go run ./cmd/stubgen -dir pkg/modules -output-dir library
```

**Output**:
```
Generated library/kubernetes.lua
Generated Lua stubs for 1 module(s) in library/
```

**What gets generated** (`library/kubernetes.lua`):
```lua
---@meta

---@class kubernetes
local kubernetes = {}

--- parseMemory: parses a Kubernetes memory quantity
---@param quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
---@return number The memory value in bytes, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_memory(quantity) end

--- parseCPU: parses a Kubernetes CPU quantity
---@param quantity string The CPU quantity to parse (e.g., "100m", "1", "2000m")
---@return number The CPU value in millicores, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_cpu(quantity) end

--- parseTime: parses a Kubernetes time string (RFC3339 format)
---@param timestr string The time string in RFC3339 format
---@return number The Unix timestamp, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_time(timestr) end

--- formatTime: converts a Unix timestamp to RFC3339 format
---@param timestamp number The Unix timestamp to convert
---@return string The time in RFC3339 format, or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_time(timestamp) end

return kubernetes
```

**How stubgen works**:
1. Scans Go files in `pkg/modules/` for special annotations
2. Finds `@luamodule` to identify modules
3. Extracts `@luafunc`, `@luaparam`, `@luareturn` annotations
4. Generates EmmyLua-compatible type definitions
5. Creates per-module `.lua` files in `library/` directory

### 2. Type Stubs (for Go structs like `corev1.Pod`)

Generate stubs from Go types using the TypeRegistry:

```bash
cd example
go run main.go
# This generates example/annotations.gen.lua automatically
```

Or programmatically in your code:
```go
registry := glua.NewTypeRegistry()
registry.Register(&corev1.Pod{})
registry.Process()
stubs, _ := registry.GenerateStubs()
os.WriteFile("annotations.gen.lua", []byte(stubs), 0644)
```

## IDE Setup for Autocomplete

### VSCode

1. **Install Lua extension**: [Lua](https://marketplace.visualstudio.com/items?itemName=sumneko.lua)

2. **Generate stubs**:
   ```bash
   # From repo root
   make stubgen          # Generates library/kubernetes.lua
   cd example
   go run main.go        # Generates example/annotations.gen.lua
   ```

3. **Open any Lua file** and enjoy autocomplete

### Neovim

1. **Install lua-language-server**:
   ```vim
   :MasonInstall lua-language-server
   ```

2. **Generate stubs** (same as above)

3. **Open any Lua file** - autocomplete works automatically

### What You Get

```lua
---@type corev1.Pod
local pod = myPod

-- Full autocomplete! Press Ctrl+Space after the dot
pod.metadata.        -- name, namespace, labels, annotations...
pod.spec.            -- containers, volumes, nodeName...
pod.spec.containers[1].  -- name, image, resources...

local k8s = require("kubernetes")
k8s.parse_           -- parse_memory, parse_cpu, parse_time, format_time
```

## Example Scripts

See [scripts/README.md](scripts/README.md) for quick reference, or [EXAMPLES.md](EXAMPLES.md) for comprehensive documentation.

**Quick overview**:
1. **Basic Pod Info** - Display metadata, labels, containers
2. **Resource Limits** - Parse and analyze CPU/memory with kubernetes module
3. **Policy Validation** - Enforce organizational policies (4 policies checked)
4. **Environment Variables** - Analyze and modify env vars, export to Go
5. **Timestamp Operations** - Parse/format timestamps, calculate pod age
6. **Multi-Container Analysis** - Analyze sidecar patterns, resource distribution
7. **JSON Export** - Transform pod data to custom report format

## What Gets Demonstrated

### Core Library Features

- **TypeRegistry** - Register Go types and generate LSP stubs
- **Translator.ToLua()** - Convert Go structs to Lua tables
- **Translator.FromLua()** - Convert Lua tables back to Go
- **Kubernetes Module** - parse_memory, parse_cpu, parse_time, format_time
- **Round-trip Integrity** - Perfect data preservation

### Real-World Use Cases

- **Resource Analysis** - Parse and validate resource limits
- **Policy Enforcement** - Validate pods against organizational rules
- **Data Transformation** - Modify pod data and export to Go
- **Complex Calculations** - Time calculations, resource aggregation
- **IDE Autocomplete** - Full type safety and IntelliSense
