# Example Lua Scripts

This directory contains example Lua scripts demonstrating various use cases with glua.

## Scripts

### 01_basic_pod_info.lua

**Basic Pod Information**

Demonstrates:

- Accessing pod metadata (name, namespace, labels)
- Iterating over containers
- Basic Lua table operations with Go structs

Usage:

```bash
go run ../main.go 01_basic_pod_info.lua
```

### 02_resource_limits.lua

**Resource Limits and Requests**

Demonstrates:

- Using the `kubernetes` module to parse CPU/memory quantities
- Calculating totals across containers
- Error handling with Lua tuples

Usage:

```bash
go run ../main.go 02_resource_limits.lua
```

### 03_policy_validation.lua

**Policy Validation**

Demonstrates:

- Implementing validation rules in Lua
- Collecting violations and warnings
- Using Lua tables for reporting
- Error handling with `error()`

Usage:

```bash
go run ../main.go 03_policy_validation.lua
```

### 04_environment_vars.lua

**Environment Variable Management**

Demonstrates:

- Reading environment variables from containers
- Modifying pod data in Lua
- Exporting modified data back to Go
- Adding new fields to tables

Usage:

```bash
go run ../main.go 04_environment_vars.lua
```

### 05_timestamp_operations.lua

**Timestamp Operations**

Demonstrates:

- Parsing Kubernetes RFC3339 timestamps
- Formatting Unix timestamps back to RFC3339
- Calculating time differences
- Using `os.time()` in Lua

Usage:

```bash
go run ../main.go 05_timestamp_operations.lua
```

### 06_multi_container_analysis.lua

**Multi-Container Pod Analysis**

Demonstrates:

- Analyzing pods with multiple containers
- Pattern detection (sidecars, proxies)
- Resource distribution calculations
- Percentage calculations in Lua

Usage:

```bash
go run ../main.go 06_multi_container_analysis.lua
```

### 07_json_export.lua

**Data Transformation / JSON Export**

Demonstrates:

- Transforming pod data to custom structure
- Building complex nested Lua tables
- Data aggregation and reporting
- Exporting results back to Go

Usage:

```bash
go run ../main.go 07_json_export.lua
```

## IDE Setup

To get autocomplete for these scripts:

1. **Generate stubs**:

   ```bash
   cd .. && make stubgen
   ```

2. **Configure your editor** (see main README.md):
   - VSCode: Create `.vscode/settings.json`
   - Neovim: Create `.luarc.json`

3. **Open any script** and enjoy autocomplete

## Creating Your Own Scripts

All scripts have access to:

- `myPod` (global) - A `corev1.Pod` object
- `kubernetes` (module) - Functions for parsing K8s types
- Standard Lua libraries

Template:

```lua
local k8s = require("kubernetes")

---@type corev1.Pod
local pod = myPod

-- Your code here

-- Export results (optional)
result = someValue
```
