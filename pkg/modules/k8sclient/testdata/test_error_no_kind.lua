-- Test: error handling - object without kind
--
-- Tests that create fails when object is missing kind.

local obj = {
	apiVersion = "v1",
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
	error("Expected error for missing kind, got nil")
end

return true
