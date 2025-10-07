-- Copyright (c) 2024-2025 Thomas Maurice
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
-- SOFTWARE.

-- Script: JSON Processing and Serialization
-- Description: Demonstrates the json module for parsing and stringifying data

local json = require("json")

print("=== JSON Module Demo ===\n")

-- 1. Parse JSON string
print("[1] Parsing JSON string...")
local jsonInput = '{"cluster":"production","nodes":5,"healthy":true,"services":["api","db","cache"]}'
local data, parseErr = json.parse(jsonInput)

if parseErr then
    error("Failed to parse JSON: " .. parseErr)
end

print("  Cluster: " .. data.cluster)
print("  Nodes: " .. tostring(data.nodes))
print("  Healthy: " .. tostring(data.healthy))
print("  Services: " .. table.concat(data.services, ", "))

-- 2. Create complex Lua table
print("\n[2] Creating and stringifying complex data...")
local podSummary = {
    name = myPod.metadata.name,
    namespace = myPod.metadata.namespace,
    labels = myPod.metadata.labels,
    containerCount = #myPod.spec.containers,
    containers = {}
}

-- Extract container info
for i, container in ipairs(myPod.spec.containers) do
    table.insert(podSummary.containers, {
        name = container.name,
        image = container.image,
        envCount = #(container.env or {})
    })
end

-- Stringify to JSON
local jsonOutput, stringifyErr = json.stringify(podSummary)

if stringifyErr then
    error("Failed to stringify: " .. stringifyErr)
end

print("  JSON output length: " .. #jsonOutput .. " bytes")
print("  JSON preview: " .. jsonOutput:sub(1, 100) .. "...")

-- 3. Round-trip test
print("\n[3] Testing round-trip conversion...")
local roundtrip, rtErr = json.parse(jsonOutput)

if rtErr then
    error("Round-trip failed: " .. rtErr)
end

if roundtrip.name ~= podSummary.name then
    error("Round-trip validation failed: name mismatch")
end

if roundtrip.containerCount ~= podSummary.containerCount then
    error("Round-trip validation failed: container count mismatch")
end

print("  ✓ Name preserved: " .. roundtrip.name)
print("  ✓ Namespace preserved: " .. roundtrip.namespace)
print("  ✓ Container count preserved: " .. roundtrip.containerCount)

-- 4. Parse nested JSON
print("\n[4] Parsing nested JSON structure...")
local nestedJson = [[
{
    "metadata": {
        "version": "1.0",
        "author": "DevOps Team"
    },
    "metrics": {
        "requests": 12345,
        "errors": 23,
        "latency": [100, 150, 200, 175]
    }
}
]]

local nested, nestedErr = json.parse(nestedJson)

if nestedErr then
    error("Failed to parse nested JSON: " .. nestedErr)
end

print("  Metadata version: " .. nested.metadata.version)
print("  Metrics requests: " .. tostring(nested.metrics.requests))
print("  Metrics errors: " .. tostring(nested.metrics.errors))
print("  Latency samples: " .. #nested.metrics.latency)

-- 5. Export pod summary for use in Go
print("\n[5] Exporting pod summary...")
exportedPodData = jsonOutput
print("  ✓ Exported as 'exportedPodData' global variable")
print("  ✓ Can be retrieved in Go code")

print("\n=== JSON Processing Complete ===")
