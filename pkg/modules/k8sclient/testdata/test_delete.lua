
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

client:create(cm)

-- Delete it
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
client:delete(gvk, TEST_NAMESPACE, TEST_DELETE_CONFIG_NAME)

-- Try to get it (should fail with an error)
local ok, err = pcall(function()
	return client:get(gvk, TEST_NAMESPACE, TEST_DELETE_CONFIG_NAME)
end)

if ok then
	error("Expected error getting deleted resource, but succeeded")
end

return true
