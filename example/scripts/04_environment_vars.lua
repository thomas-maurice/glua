
-- Example 4: Environment Variable Management
-- Analyzes and modifies environment variables

---@type corev1.Pod
local pod = myPod

print("=== Environment Variables Analysis ===")
print("Pod: " .. pod.metadata.name)
print()

-- Display all environment variables
for i, container in ipairs(pod.spec.containers) do
    print("Container: " .. container.name)

    if container.env and #container.env > 0 then
        for j, envVar in ipairs(container.env) do
            print(string.format("  %s=%s", envVar.name, envVar.value or ""))
        end
    else
        print("  (no environment variables)")
    end
    print()
end

-- Add new environment variable to all containers
print("=== Adding PROCESSED_BY_LUA environment variable ===")
for i, container in ipairs(pod.spec.containers) do
    if not container.env then
        container.env = {}
    end

    -- Check if variable already exists
    local exists = false
    for j, envVar in ipairs(container.env) do
        if envVar.name == "PROCESSED_BY_LUA" then
            exists = true
            break
        end
    end

    if not exists then
        table.insert(container.env, {
            name = "PROCESSED_BY_LUA",
            value = "true"
        })
        print("  ✓ Added to " .. container.name)
    else
        print("  - Already exists in " .. container.name)
    end
end

-- Export modified pod
modifiedPod = pod
print("\n✓ Pod modified and ready for export")
