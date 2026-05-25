
-- Test: error handling - missing version
--
-- Tests that get raises an error when GVK is missing version field.

local gvk = {group = "", kind = "ConfigMap"}  -- Missing version
local ok, err = pcall(function()
	return client:get(gvk, TEST_NAMESPACE, TEST_CONFIG_NAME)
end)

if ok then
	error("Expected error for missing version, but got success")
end

return true
