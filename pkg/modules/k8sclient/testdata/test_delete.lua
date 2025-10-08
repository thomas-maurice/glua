
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
