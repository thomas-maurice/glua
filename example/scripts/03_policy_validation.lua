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

-- Example 3: Policy Validation
-- Validates pod against organizational policies

local k8s = require("kubernetes")

---@type corev1.Pod
local pod = myPod

print("=== Policy Validation ===")
print("Pod: " .. pod.metadata.name)
print()

local violations = {}
local warnings = {}

-- Policy 1: All containers must have resource limits
print("[Policy 1] Resource Limits Required")
for i, container in ipairs(pod.spec.containers) do
    if not container.resources.limits then
        table.insert(violations, "Container '" .. container.name .. "' missing resource limits")
    elseif not container.resources.limits["memory"] then
        table.insert(violations, "Container '" .. container.name .. "' missing memory limit")
    elseif not container.resources.limits["cpu"] then
        table.insert(violations, "Container '" .. container.name .. "' missing CPU limit")
    else
        print("  ✓ " .. container.name .. " has resource limits")
    end
end
print()

-- Policy 2: Memory limit must not exceed 2GB
print("[Policy 2] Memory Limit ≤ 2GB")
for i, container in ipairs(pod.spec.containers) do
    if container.resources.limits and container.resources.limits["memory"] then
        local memBytes, err = k8s.parse_memory(container.resources.limits["memory"])
        if not err then
            local maxBytes = 2 * 1024 * 1024 * 1024
            if memBytes > maxBytes then
                table.insert(violations, "Container '" .. container.name .. "' exceeds 2GB memory limit")
            else
                print("  ✓ " .. container.name .. " within limit")
            end
        end
    end
end
print()

-- Policy 3: CPU limit should not exceed 2 cores
print("[Policy 3] CPU Limit ≤ 2 cores (warning)")
for i, container in ipairs(pod.spec.containers) do
    if container.resources.limits and container.resources.limits["cpu"] then
        local cpuMillis, err = k8s.parse_cpu(container.resources.limits["cpu"])
        if not err then
            if cpuMillis > 2000 then
                table.insert(warnings, "Container '" .. container.name .. "' uses more than 2 CPU cores")
            else
                print("  ✓ " .. container.name .. " within limit")
            end
        end
    end
end
print()

-- Policy 4: Required labels
print("[Policy 4] Required Labels")
local requiredLabels = {"app", "team", "environment"}
for _, label in ipairs(requiredLabels) do
    if not pod.metadata.labels or not pod.metadata.labels[label] then
        table.insert(violations, "Missing required label: " .. label)
    else
        print("  ✓ Label '" .. label .. "' present")
    end
end
print()

-- Summary
print("=== Validation Summary ===")
print("Violations: " .. #violations)
for i, v in ipairs(violations) do
    print("  " .. i .. ". " .. v)
end
print()

print("Warnings: " .. #warnings)
for i, w in ipairs(warnings) do
    print("  " .. i .. ". " .. w)
end
print()

if #violations > 0 then
    error("Policy validation failed with " .. #violations .. " violation(s)")
end

print("✓ All policies passed!")
