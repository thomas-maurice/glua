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

-- Test: integration - full CRUD workflow
--
-- Tests a complete create, read, update, delete workflow.

-- Create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_INTEGRATION_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_INITIAL_KEY] = TEST_INITIAL_VALUE
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

if created.metadata.name ~= TEST_INTEGRATION_NAME then
	error("Created object has wrong name: " .. (created.metadata.name or "nil"))
end

-- Get the ConfigMap
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local fetched, err = client.get(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)

if err then
	error("Failed to get: " .. err)
end

if fetched.data[TEST_INITIAL_KEY] ~= TEST_INITIAL_VALUE then
	error("Fetched object has wrong data: " .. (fetched.data[TEST_INITIAL_KEY] or "nil"))
end

-- Update the ConfigMap
fetched.data[TEST_UPDATE_KEY] = TEST_UPDATE_VALUE
local updated, err = client.update(fetched)

if err then
	error("Failed to update: " .. err)
end

if updated.data[TEST_UPDATE_KEY] ~= TEST_UPDATE_VALUE then
	error("Updated object doesn't have new data")
end

if updated.data[TEST_INITIAL_KEY] ~= TEST_INITIAL_VALUE then
	error("Updated object lost original data")
end

-- List to verify it exists
local items, err = client.list(gvk, TEST_NAMESPACE)

if err then
	error("Failed to list: " .. err)
end

local found = false
for i, item in ipairs(items) do
	if item.metadata.name == TEST_INTEGRATION_NAME then
		found = true
		break
	end
end

if not found then
	error("ConfigMap not found in list")
end

-- Delete the ConfigMap
local err = client.delete(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)

if err then
	error("Failed to delete: " .. err)
end

-- Verify it's gone
local deleted, err = client.get(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)

if not err then
	error("Expected error getting deleted resource, got nil")
end

return true
