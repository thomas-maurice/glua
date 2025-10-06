-- Test: kubernetes module integration test
--
-- This test verifies that all kubernetes module functions work
-- correctly together, including parsing, formatting, and
-- round-trip conversions.

local k8s = require("kubernetes")

-- Test parse_memory
local mem_bytes = k8s.parse_memory("1Gi")
assert(mem_bytes == 1073741824, "Memory parsing failed")

-- Test parse_cpu
local cpu_millis = k8s.parse_cpu("100m")
assert(cpu_millis == 100, "CPU parsing failed")

-- Test parse_time
local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")
assert(timestamp > 0, "Time parsing failed")

-- Test format_time
local timestr = k8s.format_time(1759509540)
assert(timestr == "2025-10-03T16:39:00Z", "Time formatting failed")

-- Test round-trip
local formatted = k8s.format_time(timestamp)
local parsed = k8s.parse_time(formatted)
assert(parsed == timestamp, "Round-trip failed")

-- Test init_defaults
local obj = {}
k8s.init_defaults(obj)
obj.metadata.labels.test = "value"
assert(obj.metadata.labels.test == "value", "init_defaults failed")

return true
