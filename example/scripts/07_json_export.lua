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
