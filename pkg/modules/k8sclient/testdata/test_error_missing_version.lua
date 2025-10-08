
-- Test: error handling - missing version
--
-- Tests that get fails when GVK is missing version field.

local gvk = {group = "", kind = "ConfigMap"}  -- Missing version
local result, err = client.get(gvk, TEST_NAMESPACE, TEST_CONFIG_NAME)

if not err then
	error("Expected error for missing version, got nil")
end

return true
