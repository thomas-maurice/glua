
-- Test: create operation
--
-- Tests creating a new Kubernetes resource.

local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_NEW_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_DATA_KEY] = TEST_DATA_VALUE
	}
}

local created = client:create(cm)

if created.metadata.name ~= TEST_NEW_CONFIG_NAME then
	error("Expected name '" .. TEST_NEW_CONFIG_NAME .. "', got " .. tostring(created.metadata.name))
end

if created.data[TEST_DATA_KEY] ~= TEST_DATA_VALUE then
	error("Expected data." .. TEST_DATA_KEY .. " '" .. TEST_DATA_VALUE .. "', got " .. tostring(created.data[TEST_DATA_KEY]))
end

return true
