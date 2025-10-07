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
