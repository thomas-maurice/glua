-- Test: error handling - object without kind
--
-- Tests that create fails when object is missing kind.

local obj = {
	apiVersion = "v1",
	metadata = {
		name = "test-config",
		namespace = "default"
	},
	data = {
		key = "value"
	}
}

local result, err = client.create(obj)

if not err then
	error("Expected error for missing kind, got nil")
end

return true
