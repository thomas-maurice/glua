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
