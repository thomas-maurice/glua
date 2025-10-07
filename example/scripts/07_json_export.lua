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

-- Example 7: JSON Export / Data Transformation
-- Transforms pod data into custom JSON structure

local k8s = require("kubernetes")

---@type corev1.Pod
local pod = myPod

print("=== Data Transformation Example ===")
print("Transforming pod to custom format...")
print()

-- Create custom data structure
local report = {
    podName = pod.metadata.name,
    namespace = pod.metadata.namespace,
    labels = pod.metadata.labels or {},
    containerCount = #pod.spec.containers,
    containers = {},
    totalResources = {
        cpuMillicores = 0,
        memoryBytes = 0
    },
    metadata = {
        creationTimestamp = pod.metadata.creationTimestamp,
        uid = pod.metadata.uid or "unknown"
    }
}

-- Process each container
for i, container in ipairs(pod.spec.containers) do
    local containerData = {
        name = container.name,
        image = container.image,
        resources = {
            cpu = 0,
            memory = 0
        },
        environmentVariables = {}
    }

    -- Parse resources
    if container.resources.limits then
        if container.resources.limits["cpu"] then
            local cpu, err = k8s.parse_cpu(container.resources.limits["cpu"])
            if not err then
                containerData.resources.cpu = cpu
                report.totalResources.cpuMillicores = report.totalResources.cpuMillicores + cpu
            end
        end

        if container.resources.limits["memory"] then
            local mem, err = k8s.parse_memory(container.resources.limits["memory"])
            if not err then
                containerData.resources.memory = mem
                report.totalResources.memoryBytes = report.totalResources.memoryBytes + mem
            end
        end
    end

    -- Collect environment variables
    if container.env then
        for j, envVar in ipairs(container.env) do
            table.insert(containerData.environmentVariables, {
                name = envVar.name,
                value = envVar.value or ""
            })
        end
    end

    table.insert(report.containers, containerData)
end

-- Print summary
print("Report Generated:")
print("  Pod: " .. report.podName)
print("  Namespace: " .. report.namespace)
print("  Containers: " .. report.containerCount)
print("  Total CPU: " .. report.totalResources.cpuMillicores .. " millicores")
print(string.format("  Total Memory: %.2f MB",
    report.totalResources.memoryBytes / (1024 * 1024)))
print()

-- Export the report (Go can access this)
exportedReport = report
print("✓ Report exported as 'exportedReport' table")
