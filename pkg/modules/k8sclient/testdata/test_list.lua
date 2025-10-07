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

-- Test: list operation
--
-- Tests listing Kubernetes resources.

-- Create multiple ConfigMaps
local cm1 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_LIST_CONFIG_1,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_LIST_KEY_1] = TEST_LIST_VALUE_1
	}
}

local cm2 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_LIST_CONFIG_2,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_LIST_KEY_2] = TEST_LIST_VALUE_2
	}
}

local created1, err = client.create(cm1)
if err then
	error("Failed to create first ConfigMap: " .. err)
end

local created2, err = client.create(cm2)
if err then
	error("Failed to create second ConfigMap: " .. err)
end

-- List ConfigMaps
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local items, err = client.list(gvk, TEST_NAMESPACE)

if err then
	error("Failed to list: " .. err)
end

-- Verify we have at least 2 items
local count = 0
for i, item in ipairs(items) do
	count = count + 1
end

if count < 2 then
	error("Expected at least 2 ConfigMaps, got " .. count)
end

-- Verify items have correct kind
for i, item in ipairs(items) do
	if item.kind ~= "ConfigMap" then
		error("Item " .. i .. " has wrong kind: " .. (item.kind or "nil"))
	end
end

return true
