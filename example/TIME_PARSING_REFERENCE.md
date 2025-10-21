# Time Parsing Reference

This document shows how to work with Kubernetes timestamp fields in Lua.

## Understanding Kubernetes Time Types

In Kubernetes YAML/JSON, time fields are RFC3339 strings:

```yaml
metadata:
  creationTimestamp: "2024-01-15T10:30:00Z"
```

In Lua, these are represented as **strings** with the type `v1.Time` or `v1.MicroTime`.

## Parsing Timestamps

### Basic Parsing

```lua
local k8s = require("kubernetes")

-- metadata.creationTimestamp is a string in Lua
local timestampStr = pod.metadata.creationTimestamp  -- "2024-01-15T10:30:00Z"

-- Parse to Unix timestamp (seconds since epoch)
local unixTime, err = k8s.parse_time(timestampStr)
if err then
    error("Failed to parse: " .. err)
end

print(unixTime)  -- 1705315800 (Unix timestamp)
```

### Formatting Timestamps

```lua
-- Convert Unix timestamp back to RFC3339 string
local unixTime = 1705315800
local rfc3339, err = k8s.format_time(unixTime)
if err then
    error("Failed to format: " .. err)
end

print(rfc3339)  -- "2024-01-15T10:30:00Z"
```

## Common Use Cases

### Calculate Pod Age

```lua
local k8s = require("kubernetes")

local creationTime, err = k8s.parse_time(pod.metadata.creationTimestamp)
if err then
    error("Parse failed: " .. err)
end

local currentTime = os.time()
local ageSeconds = currentTime - creationTime
local ageMinutes = ageSeconds / 60
local ageHours = ageMinutes / 60
local ageDays = ageHours / 24

print(string.format("Pod age: %.1f days", ageDays))
```

### Check if Resource is Recent

```lua
local k8s = require("kubernetes")

local creationTime, err = k8s.parse_time(pod.metadata.creationTimestamp)
if err then
    error("Parse failed: " .. err)
end

local ONE_HOUR = 60 * 60
if os.time() - creationTime < ONE_HOUR then
    print("Pod was created within the last hour")
end
```

### Compare Two Timestamps

```lua
local k8s = require("kubernetes")

local creation, _ = k8s.parse_time(pod.metadata.creationTimestamp)
local deletion, _ = k8s.parse_time(pod.metadata.deletionTimestamp)

if deletion and deletion > creation then
    local lifetime = deletion - creation
    print(string.format("Pod lived for %d seconds", lifetime))
end
```

### Check Deletion Grace Period

```lua
-- deletionTimestamp is v1.Time (string or nil)
if pod.metadata.deletionTimestamp then
    local deletionTime, err = k8s.parse_time(pod.metadata.deletionTimestamp)
    if not err then
        local gracePeriod = pod.metadata.deletionGracePeriodSeconds or 30
        local finalDeletion = deletionTime + gracePeriod
        local remaining = finalDeletion - os.time()

        if remaining > 0 then
            print(string.format("Pod will be deleted in %d seconds", remaining))
        else
            print("Pod should be deleted now")
        end
    end
end
```

### Format Custom Timestamps

```lua
local k8s = require("kubernetes")

-- Create a timestamp for "1 day ago"
local oneDayAgo = os.time() - (24 * 60 * 60)
local formatted, err = k8s.format_time(oneDayAgo)
if not err then
    print("24 hours ago was: " .. formatted)
end
```

### Working with Lease Times (v1.MicroTime)

```lua
local k8s = require("kubernetes")

-- For coordinationv1.Lease resources
---@type coordinationv1.Lease
local lease = {
    spec = {
        acquireTime = "2024-01-15T10:30:00.123456Z",  -- v1.MicroTime
        renewTime = "2024-01-15T10:31:00.654321Z"     -- v1.MicroTime
    }
}

-- Parse microsecond-precision timestamps (same function works)
local acquireTime, _ = k8s.parse_time(lease.spec.acquireTime)
local renewTime, _ = k8s.parse_time(lease.spec.renewTime)

local renewalInterval = renewTime - acquireTime
print(string.format("Lease renewed after %d seconds", renewalInterval))
```

### Working with Event Times

```lua
local k8s = require("kubernetes")

-- For eventsv1.Event resources
---@type eventsv1.Event
local event = {
    eventTime = "2024-01-15T10:30:00.123456Z",  -- v1.MicroTime
    series = {
        lastObservedTime = "2024-01-15T10:35:00.987654Z"  -- v1.MicroTime
    }
}

local eventTime, _ = k8s.parse_time(event.eventTime)
local lastTime, _ = k8s.parse_time(event.series.lastObservedTime)

local duration = lastTime - eventTime
print(string.format("Event occurred over %d seconds", duration))
```

## Type Information

### v1.Time Fields

Common fields with `v1.Time` type:

- `metadata.creationTimestamp` - When resource was created
- `metadata.deletionTimestamp` - When resource deletion was requested
- `status.conditions[].lastTransitionTime` - When condition last changed
- `status.conditions[].lastUpdateTime` - When condition was last updated

### v1.MicroTime Fields

Common fields with `v1.MicroTime` type (microsecond precision):

- `coordinationv1.Lease.spec.acquireTime` - When lease was acquired
- `coordinationv1.Lease.spec.renewTime` - When lease was last renewed
- `eventsv1.Event.eventTime` - When event occurred
- `eventsv1.EventSeries.lastObservedTime` - Last time event was observed

## Error Handling

Always check for errors when parsing timestamps:

```lua
local k8s = require("kubernetes")

local timestamp, err = k8s.parse_time(pod.metadata.creationTimestamp)
if err then
    -- Handle invalid timestamp
    print("WARNING: Invalid timestamp: " .. err)
    return
end

-- Safe to use timestamp
print("Unix time: " .. timestamp)
```

## See Also

- [05_timestamp_operations.lua](scripts/05_timestamp_operations.lua) - Full example with age calculations
- [kubernetes module API](../README.md#kubernetes-module) - Full kubernetes module documentation
