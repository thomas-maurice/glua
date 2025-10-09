
-- Test: error handling - object without apiVersion
--
-- Tests that create fails when object is missing apiVersion.

local obj = {
	kind = "ConfigMap",
	metadata = {
		name = TEST_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_CONFIG_KEY] = TEST_CONFIG_VALUE
	}
}

local result, err = client.create(obj)

if not err then
	error("Expected error for missing apiVersion, got nil")
end

return true
