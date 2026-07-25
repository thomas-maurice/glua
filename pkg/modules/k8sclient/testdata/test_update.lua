
-- Test: update operation
--
-- Tests updating an existing Kubernetes resource.

-- Create initial ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_UPDATE_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_ORIGINAL_KEY] = TEST_ORIGINAL_VALUE
	}
}

local created = client:create(cm)

-- Update it
created.data[TEST_UPDATE_KEY] = TEST_UPDATE_VALUE
created.data[TEST_ORIGINAL_KEY] = "changed"

local updated = client:update(created)

if updated.data[TEST_UPDATE_KEY] ~= TEST_UPDATE_VALUE then
	error("Expected data." .. TEST_UPDATE_KEY .. " '" .. TEST_UPDATE_VALUE .. "', got " .. tostring(updated.data[TEST_UPDATE_KEY]))
end

if updated.data[TEST_ORIGINAL_KEY] ~= "changed" then
	error("Expected data." .. TEST_ORIGINAL_KEY .. " 'changed', got " .. tostring(updated.data[TEST_ORIGINAL_KEY]))
end

return true
