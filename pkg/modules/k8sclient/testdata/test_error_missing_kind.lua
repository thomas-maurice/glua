-- Test: error handling - missing kind
--
-- Tests that get fails when GVK is missing kind field.

local gvk = {group = "", version = "v1"}  -- Missing kind
local result, err = client.get(gvk, TEST_NAMESPACE, TEST_CONFIG_NAME)

if not err then
	error("Expected error for missing kind, got nil")
end

return true
