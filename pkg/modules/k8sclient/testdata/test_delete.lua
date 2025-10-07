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

-- Test: delete operation
--
-- Tests deleting a Kubernetes resource.

-- Create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_DELETE_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_DELETE_DATA_KEY] = TEST_DELETE_DATA_VALUE
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

-- Delete it
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local err = client.delete(gvk, TEST_NAMESPACE, TEST_DELETE_CONFIG_NAME)

if err then
	error("Failed to delete: " .. err)
end

-- Try to get it (should fail)
local fetched, err = client.get(gvk, TEST_NAMESPACE, TEST_DELETE_CONFIG_NAME)

if not err then
	error("Expected error getting deleted resource, got nil")
end

return true
