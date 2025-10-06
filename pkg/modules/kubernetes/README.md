# Kubernetes Module

The `kubernetes` module provides utilities for working with Kubernetes resource quantities and timestamps in Lua scripts.

## Functions

### `kubernetes.parse_memory(quantity)`

Parses a Kubernetes memory quantity string and returns the value in bytes.

**Parameters:**
- `quantity` (string): The memory quantity to parse (e.g., "1024Mi", "1Gi", "512M")

**Returns:**
- `number`: The memory value in bytes, or nil on error
- `string|nil`: Error message if parsing failed

**Example:**
```lua
local k8s = require("kubernetes")
local bytes, err = k8s.parse_memory("1024Mi")

if err then
    print("Error: " .. err)
else
    print(bytes)  -- prints 1073741824
end
```

### `kubernetes.parse_cpu(quantity)`

Parses a Kubernetes CPU quantity string and returns the value in millicores.

**Parameters:**
- `quantity` (string): The CPU quantity to parse (e.g., "100m", "1", "2000m")

**Returns:**
- `number`: The CPU value in millicores, or nil on error
- `string|nil`: Error message if parsing failed

**Example:**
```lua
local k8s = require("kubernetes")
local millicores, err = k8s.parse_cpu("100m")

if err then
    print("Error: " .. err)
else
    print(millicores)  -- prints 100
end

-- Whole CPUs are converted to millicores
local mc2 = k8s.parse_cpu("1")  -- returns 1000
```

### `kubernetes.parse_time(timestr)`

Parses a Kubernetes time string (RFC3339 format) and returns a Unix timestamp.

**Parameters:**
- `timestr` (string): The time string in RFC3339 format (e.g., "2025-10-03T16:39:00Z")

**Returns:**
- `number`: The Unix timestamp, or nil on error
- `string|nil`: Error message if parsing failed

**Example:**
```lua
local k8s = require("kubernetes")
local timestamp, err = k8s.parse_time("2025-10-03T16:39:00Z")

if err then
    print("Error: " .. err)
else
    print(timestamp)  -- prints 1759509540
end
```

### `kubernetes.format_time(timestamp)`

Converts a Unix timestamp to a Kubernetes time string in RFC3339 format.

**Parameters:**
- `timestamp` (number): The Unix timestamp to convert

**Returns:**
- `string`: The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
- `string|nil`: Error message if formatting failed

**Example:**
```lua
local k8s = require("kubernetes")
local timestr, err = k8s.format_time(1759509540)

if err then
    print("Error: " .. err)
else
    print(timestr)  -- prints "2025-10-03T16:39:00Z"
end
```

### `kubernetes.init_defaults(obj)`

Initializes default empty tables for `metadata.labels` and `metadata.annotations` if they are nil. This is useful for ensuring these fields are tables instead of nil, making it easier to add labels/annotations in Lua without checking for nil first.

**Parameters:**
- `obj` (table): The Kubernetes object (must have a metadata field)

**Returns:**
- `table`: The same object with initialized defaults (modified in-place)

**Example:**
```lua
local k8s = require("kubernetes")

local pod = {
    metadata = {
        name = "my-pod",
        namespace = "default"
    }
}

-- Before init_defaults, labels and annotations are nil
-- This would error: pod.metadata.labels.app = "myapp"

-- Initialize defaults
k8s.init_defaults(pod)

-- Now safe to add labels and annotations
pod.metadata.labels.app = "myapp"
pod.metadata.labels.tier = "backend"
pod.metadata.annotations["version"] = "1.0.0"

print(pod.metadata.labels.app)  -- prints "myapp"
```

**Notes:**
- The function modifies the object in-place and also returns it
- If metadata doesn't exist, it will be created
- Existing labels and annotations are preserved
- Only nil values are replaced with empty tables

### `kubernetes.parse_duration(duration_str)`

Parses a duration string (e.g., "5m", "1h30m") and returns the value in seconds.

**Parameters:**
- `duration_str` (string): The duration string to parse (e.g., "5m", "1h", "1h30m45s")

**Returns:**
- `number`: The duration in seconds, or nil on error
- `string|nil`: Error message if parsing failed

**Example:**
```lua
local k8s = require("kubernetes")
local seconds, err = k8s.parse_duration("5m")

if err then
    print("Error: " .. err)
else
    print(seconds)  -- prints 300
end

-- Complex durations
local s2 = k8s.parse_duration("1h30m")  -- returns 5400
```

### `kubernetes.format_duration(seconds)`

Converts a duration in seconds to a duration string.

**Parameters:**
- `seconds` (number): The duration in seconds

**Returns:**
- `string`: The formatted duration string (e.g., "5m0s"), or nil on error
- `string|nil`: Error message if formatting failed

**Example:**
```lua
local k8s = require("kubernetes")
local duration_str, err = k8s.format_duration(300)

if err then
    print("Error: " .. err)
else
    print(duration_str)  -- prints "5m0s"
end
```

### `kubernetes.parse_int_or_string(value)`

Handles Kubernetes IntOrString type, determining if a value is a number or string.

**Parameters:**
- `value` (number|string): The value to check

**Returns:**
- `number|string`: The input value
- `boolean`: true if the value is a string, false if it's a number

**Example:**
```lua
local k8s = require("kubernetes")

-- With a number
local val1, is_str1 = k8s.parse_int_or_string(8080)
print(val1)     -- prints 8080
print(is_str1)  -- prints false

-- With a string
local val2, is_str2 = k8s.parse_int_or_string("http")
print(val2)     -- prints "http"
print(is_str2)  -- prints true
```

### `kubernetes.matches_selector(labels, selector)`

Checks if a set of labels matches a label selector.

**Parameters:**
- `labels` (table): The label set to check
- `selector` (table): The label selector (key-value pairs)

**Returns:**
- `boolean`: true if all selector labels match, false otherwise

**Example:**
```lua
local k8s = require("kubernetes")

local pod_labels = {
    app = "nginx",
    tier = "frontend",
    version = "v1"
}

local selector = {
    app = "nginx",
    tier = "frontend"
}

local matches = k8s.matches_selector(pod_labels, selector)
print(matches)  -- prints true
```

### `kubernetes.toleration_matches(toleration, taint)`

Checks if a toleration matches a taint.

**Parameters:**
- `toleration` (table): The toleration with fields: key, operator, value, effect
- `taint` (table): The taint with fields: key, value, effect

**Returns:**
- `boolean`: true if the toleration matches the taint, false otherwise

**Example:**
```lua
local k8s = require("kubernetes")

-- Equal operator
local toleration = {
    key = "node-role",
    operator = "Equal",
    value = "master",
    effect = "NoSchedule"
}

local taint = {
    key = "node-role",
    value = "master",
    effect = "NoSchedule"
}

local matches = k8s.toleration_matches(toleration, taint)
print(matches)  -- prints true

-- Exists operator (value doesn't matter)
local tol2 = {
    key = "node-role",
    operator = "Exists",
    effect = "NoSchedule"
}

local matches2 = k8s.toleration_matches(tol2, taint)
print(matches2)  -- prints true
```

### `kubernetes.match_gvk(obj, matcher)`

Checks if a Kubernetes object matches the specified Group/Version/Kind (GVK) matcher.

**Parameters:**
- `obj` (table): The Kubernetes object to check
- `matcher` (GVKMatcher): A table with `group`, `version`, and `kind` fields

**Returns:**
- `boolean`: true if the object's apiVersion and kind match the matcher's GVK

**Example:**
```lua
local k8s = require("kubernetes")

-- Check if object is a Pod
local pod = {
    apiVersion = "v1",
    kind = "Pod",
}

local pod_matcher = {group = "", version = "v1", kind = "Pod"}
local is_pod = k8s.match_gvk(pod, pod_matcher)
print(is_pod)  -- prints true

-- Check if object is a Deployment
local deployment = {
    apiVersion = "apps/v1",
    kind = "Deployment",
}

local deployment_matcher = {group = "apps", version = "v1", kind = "Deployment"}
local is_deployment = k8s.match_gvk(deployment, deployment_matcher)
print(is_deployment)  -- prints true

-- Check for wrong type
local service_matcher = {group = "", version = "v1", kind = "Service"}
local is_service = k8s.match_gvk(pod, service_matcher)
print(is_service)  -- prints false
```

**Notes:**
- For core API resources (Pod, Service, ConfigMap, etc.), use an empty string for the group field
- For resources in other API groups (Deployment, StatefulSet, etc.), specify the group name (e.g., "apps", "batch")
- The apiVersion field in the object should match the group/version format (e.g., "v1" for core resources, "apps/v1" for apps group)
- The GVKMatcher is a Go type registered with the TypeRegistry, allowing seamless conversion between Go and Lua

## Usage in Go

```go
package main

import (
    "github.com/thomas-maurice/glua/pkg/modules/kubernetes"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // Register the kubernetes module
    L.PreloadModule("kubernetes", kubernetes.Loader)

    // Use in Lua
    L.DoString(`
        local k8s = require("kubernetes")

        -- Parse resource quantities
        local mem_bytes = k8s.parse_memory("512Mi")
        local cpu_millis = k8s.parse_cpu("250m")

        print("Memory: " .. mem_bytes .. " bytes")
        print("CPU: " .. cpu_millis .. " millicores")

        -- Work with timestamps
        local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")
        local formatted = k8s.format_time(timestamp)
        print("Time: " .. formatted)

        -- Initialize Kubernetes object defaults
        local pod = {
            metadata = {
                name = "my-pod"
            }
        }

        k8s.init_defaults(pod)
        pod.metadata.labels.app = "web"
        pod.metadata.annotations.version = "1.0"
    `)
}
```

## Common Use Cases

### 1. Parsing Resource Limits

```lua
local k8s = require("kubernetes")

local container = myPod.spec.containers[1]

-- Parse memory limit
local mem_limit = container.resources.limits.memory
local mem_bytes = k8s.parse_memory(mem_limit)

if mem_bytes > 1073741824 then  -- 1Gi
    print("Memory limit exceeds 1Gi")
end

-- Parse CPU request
local cpu_request = container.resources.requests.cpu
local cpu_millis = k8s.parse_cpu(cpu_request)

if cpu_millis < 100 then
    print("Warning: CPU request is very low")
end
```

### 2. Calculating Pod Age

```lua
local k8s = require("kubernetes")

local creation_time = myPod.metadata.creationTimestamp
local created_at = k8s.parse_time(creation_time)
local now = os.time()

local age_seconds = now - created_at
local age_hours = age_seconds / 3600

print(string.format("Pod age: %.1f hours", age_hours))
```

### 3. Adding Labels and Annotations

```lua
local k8s = require("kubernetes")

-- Initialize defaults to avoid nil errors
k8s.init_defaults(myPod)

-- Add labels
myPod.metadata.labels.environment = "production"
myPod.metadata.labels.version = "v2.0"
myPod.metadata.labels.owner = "platform-team"

-- Add annotations
myPod.metadata.annotations["deployment.kubernetes.io/revision"] = "5"
myPod.metadata.annotations["prometheus.io/scrape"] = "true"
myPod.metadata.annotations["prometheus.io/port"] = "8080"
```

### 4. Validating Resource Configurations

```lua
local k8s = require("kubernetes")

function validate_resources(pod)
    for i, container in ipairs(pod.spec.containers) do
        -- Check if resources are defined
        if not container.resources or not container.resources.limits then
            error("Container " .. container.name .. " has no resource limits")
        end

        -- Parse and validate memory
        local mem_limit = k8s.parse_memory(container.resources.limits.memory)
        if mem_limit > 8589934592 then  -- 8Gi
            error("Container " .. container.name .. " exceeds max memory of 8Gi")
        end

        -- Parse and validate CPU
        local cpu_limit = k8s.parse_cpu(container.resources.limits.cpu)
        if cpu_limit > 4000 then  -- 4 CPUs
            error("Container " .. container.name .. " exceeds max CPU of 4 cores")
        end
    end
end

validate_resources(myPod)
```

### 5. Time-based Operations

```lua
local k8s = require("kubernetes")

-- Parse pod creation time
local created = k8s.parse_time(myPod.metadata.creationTimestamp)

-- Calculate time windows
local one_hour_ago = os.time() - 3600
local one_day_ago = os.time() - 86400

if created > one_hour_ago then
    print("Pod created within the last hour")
elseif created > one_day_ago then
    print("Pod created within the last day")
else
    print("Pod is older than one day")
end

-- Format times for logging
local created_str = k8s.format_time(created)
print("Pod created at: " .. created_str)
```

## Data Type Conversions

### Memory Quantities

| Input | Bytes | Notes |
|-------|-------|-------|
| "1Ki" | 1024 | Kibibytes (binary) |
| "1Mi" | 1048576 | Mebibytes (binary) |
| "1Gi" | 1073741824 | Gibibytes (binary) |
| "1Ti" | 1099511627776 | Tebibytes (binary) |
| "1K" | 1000 | Kilobytes (decimal) |
| "1M" | 1000000 | Megabytes (decimal) |
| "1G" | 1000000000 | Gigabytes (decimal) |

### CPU Quantities

| Input | Millicores | Notes |
|-------|------------|-------|
| "100m" | 100 | 100 millicores |
| "250m" | 250 | 250 millicores |
| "1" | 1000 | 1 full CPU = 1000 millicores |
| "2" | 2000 | 2 full CPUs = 2000 millicores |
| "0.5" | 500 | Half a CPU = 500 millicores |

### Time Formats

The module uses RFC3339 format (ISO 8601), which is the standard for Kubernetes:

- Input format: `"2025-10-03T16:39:00Z"` or `"2025-10-03T16:39:00+00:00"`
- Output format: `"2025-10-03T16:39:00Z"` (always UTC)
- Unix timestamps are integers (seconds since epoch)

## Error Handling

All parsing functions return two values: the result and an error message. Always check for errors:

```lua
local k8s = require("kubernetes")

local bytes, err = k8s.parse_memory("invalid")
if err then
    print("Parse error: " .. err)
    return
end

-- Use bytes safely
print("Parsed: " .. bytes)
```

## Testing

Run the test suite:

```bash
go test ./pkg/modules/kubernetes/
```

## Integration

The kubernetes module is automatically included when you run `make stubgen` and will generate IDE autocomplete stubs in `library/kubernetes.gen.lua`.
