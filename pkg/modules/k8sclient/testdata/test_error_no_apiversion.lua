-- Test: error handling - object without apiVersion
--
-- Tests that create fails when object is missing apiVersion.

local obj = {
	kind = "ConfigMap",
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
	error("Expected error for missing apiVersion, got nil")
end

return true
