# Example Lua Scripts - Complete Guide

This directory contains 7 comprehensive example scripts demonstrating all features of the glua library.

## 🎯 Quick Start

### 1. Generate Stubs

From the repository root:
```bash
make stubgen
```

This generates:
- `library/kubernetes.lua` - Kubernetes module autocomplete
- `example/annotations.gen.lua` - Type definitions for Pod and related types

### 2. Run an Example

```bash
cd example
go run run_script.go scripts/01_basic_pod_info.lua
```

### 3. Open Scripts in Your IDE

The scripts are located in `example/scripts/`. Open any `.lua` file in VSCode or Neovim and you'll get:
- ✨ **Full autocomplete** for Pod types
- ✨ **Autocomplete** for kubernetes module functions
- ✨ **Type checking** and error detection
- ✨ **Hover documentation** for all functions

## 📝 Example Scripts

### 01_basic_pod_info.lua
**What it does**: Displays basic pod information including name, namespace, labels, and containers.

**Features demonstrated**:
- Accessing pod metadata
- Iterating over Lua tables
- Working with Pod labels
- Iterating over containers

**Autocomplete features**:
- `pod.metadata.` → shows name, namespace, labels, etc.
- `pod.spec.` → shows containers, volumes, etc.
- `container.` → shows name, image, resources, etc.

---

### 02_resource_limits.lua
**What it does**: Analyzes and calculates resource limits/requests for all containers.

**Features demonstrated**:
- Using `kubernetes.parse_cpu()` and `kubernetes.parse_memory()`
- Error handling with Lua tuples `(value, error)`
- Calculating totals across containers
- String formatting with `string.format()`

**Autocomplete features**:
- `k8s.parse_` → shows parse_memory, parse_cpu, parse_time, format_time
- Hover over functions shows parameter types and return values
- `container.resources.limits` → autocomplete for CPU/memory

---

### 03_policy_validation.lua
**What it does**: Validates pod against organizational policies and reports violations.

**Features demonstrated**:
- Policy enforcement in Lua
- Building violation/warning lists
- Complex validation logic
- Using Lua's `error()` to fail validation

**Policies checked**:
1. All containers must have resource limits
2. Memory limit must not exceed 2GB
3. CPU limit should not exceed 2 cores (warning)
4. Required labels: app, team, environment

**Autocomplete features**:
- Full navigation through pod structure
- Type-safe access to resources
- Label autocomplete

---

### 04_environment_vars.lua
**What it does**: Analyzes environment variables and adds new ones to containers.

**Features demonstrated**:
- Reading environment variables from containers
- Modifying pod data in Lua
- Adding new fields to Lua tables
- Exporting modified data back to Go

**Key technique**: Shows how to modify pod data and export it:
```lua
modifiedPod = pod  -- Export to Go
```

**Autocomplete features**:
- `container.env` → autocomplete for environment variable array
- `envVar.` → shows name, value, valueFrom

---

### 05_timestamp_operations.lua
**What it does**: Parses Kubernetes timestamps, calculates pod age, formats timestamps.

**Features demonstrated**:
- `kubernetes.parse_time()` - RFC3339 → Unix timestamp
- `kubernetes.format_time()` - Unix timestamp → RFC3339
- Using `os.time()` for current time
- Date/time calculations in Lua
- Time formatting and age calculations

**Calculations shown**:
- Pod age in seconds, minutes, hours, days
- Formatting timestamps to different formats
- Creating custom timestamps

**Autocomplete features**:
- `k8s.parse_time` → shows parameters and return type
- `k8s.format_time` → shows conversion direction

---

### 06_multi_container_analysis.lua
**What it does**: Analyzes pods with multiple containers (sidecars, init containers).

**Features demonstrated**:
- Detecting container patterns (nginx, envoy, sidecars)
- Resource distribution calculations
- Percentage calculations
- Providing recommendations based on analysis

**Analysis performed**:
- Container count (regular + init containers)
- Pattern detection (web servers, proxies)
- Resource distribution per container
- Recommendations for optimization

**Autocomplete features**:
- `pod.spec.initContainers` → autocomplete for init containers
- Full navigation through container arrays
- Access to all container properties

---

### 07_json_export.lua
**What it does**: Transforms pod data into a custom report structure for export.

**Features demonstrated**:
- Building complex nested Lua tables
- Data transformation and aggregation
- Creating reports for export to Go
- Collecting data from multiple sources

**Export format**:
```lua
{
  podName = "...",
  namespace = "...",
  containerCount = 1,
  containers = { ... },
  totalResources = {
    cpuMillicores = 100,
    memoryBytes = 104857600
  }
}
```

**Autocomplete features**:
- Full type safety while building custom structures
- Autocomplete for all pod fields being read
- Type hints for kubernetes module functions

## 🎨 IDE Setup for Autocomplete

### VSCode

1. **Install Lua extension**:
   - Install "Lua" by sumneko from the marketplace

2. **Open the example/scripts directory**:
   ```bash
   code /home/thomas/git/glua/example/scripts
   ```

3. **Autocomplete works immediately!** The `.luarc.json` is already configured.

4. **Try it**:
   - Open any `.lua` file
   - Type `pod.` and press Ctrl+Space
   - Type `k8s.` and press Ctrl+Space
   - Hover over any function to see documentation

### Neovim

1. **Install lua-language-server**:
   ```vim
   :MasonInstall lua-language-server
   ```

2. **Open any script**:
   ```bash
   nvim example/scripts/01_basic_pod_info.lua
   ```

3. **Autocomplete works!** The `.luarc.json` is automatically detected.

4. **Try it**:
   - Type `pod.` then `<C-x><C-o>` (or use your completion plugin)
   - Type `k8s.` for module autocomplete
   - Press `K` on any function for hover docs

## 📊 What Gets Autocompleted

### Pod Type (corev1.Pod)

```lua
---@type corev1.Pod
local pod = myPod

pod.                      -- All fields autocomplete!
pod.metadata.             -- name, namespace, labels, annotations...
pod.metadata.labels       -- table<string, string>
pod.spec.                 -- containers, volumes, nodeName...
pod.spec.containers       -- corev1.Container[]
pod.spec.containers[1].   -- name, image, resources...
pod.status.               -- phase, conditions, podIP...
```

### Kubernetes Module

```lua
local k8s = require("kubernetes")

k8s.                      -- All functions autocomplete!
k8s.parse_memory(         -- Shows: (quantity: string) -> number, string|nil
k8s.parse_cpu(            -- Shows: (quantity: string) -> number, string|nil
k8s.parse_time(           -- Shows: (timestr: string) -> number, string|nil
k8s.format_time(          -- Shows: (timestamp: number) -> string, string|nil
```

### Container Type

```lua
local container = pod.spec.containers[1]

container.                -- All fields autocomplete!
container.name            -- string
container.image           -- string
container.resources.      -- limits, requests
container.env             -- corev1.EnvVar[]
container.env[1].         -- name, value, valueFrom
```

## 🔍 How It Works

### Type Annotations

The `annotations.gen.lua` file contains type definitions like:

```lua
---@class corev1.Pod
---@field kind string
---@field apiVersion string
---@field metadata v1.ObjectMeta
---@field spec corev1.PodSpec
---@field status corev1.PodStatus
```

When you write:
```lua
---@type corev1.Pod
local pod = myPod
```

The LSP knows `pod` has the structure of `corev1.Pod` and provides autocomplete!

### Module Stubs

The `library/kubernetes.lua` file contains function signatures:

```lua
---@param quantity string The memory quantity to parse
---@return number The memory value in bytes
---@return string|nil Error message if parsing failed
function kubernetes.parse_memory(quantity) end
```

When you type `k8s.parse_memory(`, the LSP shows you the parameters and return types!

## 🚀 Creating Your Own Scripts

Template:

```lua
-- Require the kubernetes module if you need it
local k8s = require("kubernetes")

-- Type annotation for autocomplete
---@type corev1.Pod
local pod = myPod

-- Your code here
print("Pod: " .. pod.metadata.name)

-- Export results if needed
result = someValue
```

Available globals:
- `myPod` - The sample pod (type: `corev1.Pod`)
- `originalPod` - Alias for `myPod`

Available modules:
- `kubernetes` - CPU/memory/time parsing functions

## 📚 Learn More

- **Main README**: `../README.md` - Full library documentation
- **Scripts README**: `scripts/README.md` - Quick reference for each script
- **API Reference**: See README.md for complete API documentation
- **GoDoc**: https://godoc.org/github.com/thomas-maurice/glua

## 💡 Tips

1. **Always use type annotations**:
   ```lua
   ---@type corev1.Pod
   local pod = myPod
   ```
   This enables autocomplete!

2. **Check for nil before using values**:
   ```lua
   if pod.metadata.labels then
       for k, v in pairs(pod.metadata.labels) do
           print(k, v)
       end
   end
   ```

3. **Handle errors from kubernetes module**:
   ```lua
   local value, err = k8s.parse_memory("256Mi")
   if err then
       error("Failed: " .. err)
   end
   ```

4. **Use string.format for better output**:
   ```lua
   print(string.format("Memory: %.2f MB", bytes / (1024 * 1024)))
   ```

5. **Export data to Go**:
   ```lua
   modifiedPod = pod           -- Go can access this
   exportedReport = myReport   -- Go can access this too
   ```

Happy scripting! ✨
