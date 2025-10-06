-- Test: get operation
--
-- Tests retrieving a Kubernetes resource by GVK, namespace, and name.

-- First create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "test-config",
		namespace = "default"
	},
	data = {
		key1 = "value1"
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

-- Now get it
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local fetched, err = client.get(gvk, "default", "test-config")

if err then
	error("Failed to get: " .. err)
end

if fetched.metadata.name ~= "test-config" then
	error("Expected name 'test-config', got " .. tostring(fetched.metadata.name))
end

if fetched.data.key1 ~= "value1" then
	error("Expected data.key1 'value1', got " .. tostring(fetched.data.key1))
end

return true
