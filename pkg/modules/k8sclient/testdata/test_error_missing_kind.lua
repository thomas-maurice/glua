
-- Test: error handling - missing kind
--
-- Tests that get raises an error when GVK is missing kind field.

local gvk = {group = "", version = "v1"}  -- Missing kind
local ok, err = pcall(function()
	return client:get(gvk, TEST_NAMESPACE, TEST_CONFIG_NAME)
end)

if ok then
	error("Expected error for missing kind, but got success")
end

return true
