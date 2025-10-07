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
