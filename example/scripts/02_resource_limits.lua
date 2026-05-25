
-- Example 2: Resource Limits and Requests
-- Parse and validate resource limits for all containers

local k8s = require("kubernetes")

---@type corev1.Pod
---@diagnostic disable-next-line: undefined-global
local pod = myPod

print("=== Resource Limits Analysis ===")
print("Pod: " .. pod.metadata.name)
print()

local totalCPU = 0
local totalMemory = 0

for i, container in ipairs(pod.spec.containers) do
    print("Container: " .. container.name)

    -- Parse CPU limits. Q5: parse_cpu raises on a malformed value; wrap in
    -- pcall if you want to keep going past a bad container.
    if container.resources.limits and container.resources.limits["cpu"] then
        local cpuMillis = k8s.parse_cpu(container.resources.limits["cpu"])
        print("  CPU Limit: " .. cpuMillis .. " millicores")
        totalCPU = totalCPU + cpuMillis
    else
        print("  CPU Limit: NOT SET")
    end

    -- Parse memory limits.
    if container.resources.limits and container.resources.limits["memory"] then
        local memBytes = k8s.parse_memory(container.resources.limits["memory"])
        local memMB = memBytes / (1024 * 1024)
        print(string.format("  Memory Limit: %.2f MB", memMB))
        totalMemory = totalMemory + memBytes
    else
        print("  Memory Limit: NOT SET")
    end

    print()
end

print("=== Totals ===")
print("Total CPU: " .. totalCPU .. " millicores (" .. (totalCPU / 1000) .. " cores)")
print(string.format("Total Memory: %.2f MB (%.2f GB)",
    totalMemory / (1024 * 1024),
    totalMemory / (1024 * 1024 * 1024)))
