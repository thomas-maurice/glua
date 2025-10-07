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

-- Test: get operation
--
-- Tests retrieving a Kubernetes resource by GVK, namespace, and name.

-- First create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_CONFIG_KEY] = TEST_CONFIG_VALUE
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

-- Now get it
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local fetched, err = client.get(gvk, TEST_NAMESPACE, TEST_CONFIG_NAME)

if err then
	error("Failed to get: " .. err)
end

if fetched.metadata.name ~= TEST_CONFIG_NAME then
	error("Expected name '" .. TEST_CONFIG_NAME .. "', got " .. tostring(fetched.metadata.name))
end

if fetched.data[TEST_CONFIG_KEY] ~= TEST_CONFIG_VALUE then
	error("Expected data." .. TEST_CONFIG_KEY .. " '" .. TEST_CONFIG_VALUE .. "', got " .. tostring(fetched.data[TEST_CONFIG_KEY]))
end

return true
