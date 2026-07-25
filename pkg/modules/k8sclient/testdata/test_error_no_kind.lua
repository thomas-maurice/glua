
-- Test: error handling - object without kind
--
-- Tests that create raises an error when object is missing kind.

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

local ok, err = pcall(function()
	return client:create(obj)
end)

if ok then
	error("Expected error for missing kind, but got success")
end

return true
