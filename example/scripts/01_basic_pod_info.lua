
-- Example 1: Basic Pod Information
-- This script extracts and displays basic pod information

---@type corev1.Pod
local pod = myPod

print("=== Basic Pod Information ===")
print("Name: " .. pod.metadata.name)
print("Namespace: " .. pod.metadata.namespace)
print("Kind: " .. pod.kind)
print("API Version: " .. pod.apiVersion)

-- Check if pod has labels
if pod.metadata.labels then
    print("\nLabels:")
    for key, value in pairs(pod.metadata.labels) do
        print("  " .. key .. ": " .. value)
    end
end

-- Container information
print("\nContainers:")
for i, container in ipairs(pod.spec.containers) do
    print("  " .. i .. ". " .. container.name .. " (" .. container.image .. ")")
end
