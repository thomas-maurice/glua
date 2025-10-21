
-- Example 6: Multi-Container Pod Analysis
-- Analyzes pods with multiple containers (sidecars, init containers)

local k8s = require("kubernetes")

---@type corev1.Pod
---@diagnostic disable-next-line: undefined-global
local pod = myPod

print("=== Multi-Container Analysis ===")
print("Pod: " .. pod.metadata.name)
print()

-- Count containers
local containerCount = #pod.spec.containers
local initContainerCount = 0
if pod.spec.initContainers then
    initContainerCount = #pod.spec.initContainers
end

print("Container Count: " .. containerCount)
print("Init Container Count: " .. initContainerCount)
print()

-- Identify container patterns
print("=== Container Patterns ===")

local hasNginx = false
local hasEnvoy = false
local hasSidecar = false

for i, container in ipairs(pod.spec.containers) do
    local image = container.image:lower()

    if image:find("nginx") then
        hasNginx = true
        print("  ✓ Nginx container detected: " .. container.name)
    elseif image:find("envoy") then
        hasEnvoy = true
        print("  ✓ Envoy sidecar detected: " .. container.name)
    elseif i > 1 then
        hasSidecar = true
        print("  ✓ Additional container (sidecar): " .. container.name)
    end
end
print()

-- Analyze resource distribution
print("=== Resource Distribution ===")
local resources = {}

for i, container in ipairs(pod.spec.containers) do
    local cpu = 0
    local memory = 0

    if container.resources.limits then
        if container.resources.limits["cpu"] then
            cpu, _ = k8s.parse_cpu(container.resources.limits["cpu"])
        end
        if container.resources.limits["memory"] then
            memory, _ = k8s.parse_memory(container.resources.limits["memory"])
        end
    end

    table.insert(resources, {
        name = container.name,
        cpu = cpu,
        memory = memory
    })
end

-- Calculate totals
local totalCPU = 0
local totalMemory = 0
for _, res in ipairs(resources) do
    totalCPU = totalCPU + res.cpu
    totalMemory = totalMemory + res.memory
end

-- Show percentages
for _, res in ipairs(resources) do
    local cpuPercent = 0
    local memPercent = 0

    if totalCPU > 0 then
        cpuPercent = (res.cpu / totalCPU) * 100
    end
    if totalMemory > 0 then
        memPercent = (res.memory / totalMemory) * 100
    end

    print(string.format("%s:", res.name))
    print(string.format("  CPU: %d millicores (%.1f%%)", res.cpu, cpuPercent))
    print(string.format("  Memory: %.2f MB (%.1f%%)",
        res.memory / (1024 * 1024), memPercent))
end
print()

-- Recommendations
print("=== Recommendations ===")
if containerCount > 3 then
    print("  ⚠️  Consider reducing number of containers (" .. containerCount .. " is high)")
else
    print("  ✓ Container count is reasonable")
end

if hasSidecar and not hasEnvoy then
    print("  ℹ️  Consider using Envoy for sidecar proxy patterns")
end
