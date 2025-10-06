-- Test: parse_int_or_string validation
--
-- Verifies that parse_int_or_string errors for invalid types.

local k8s = require("kubernetes")

-- Test with boolean (invalid)
local status, err = pcall(function()
	k8s.parse_int_or_string(true)
end)

if status then
	error("Expected error for boolean input, but got success")
end

if not string.find(err, "expects number or string") then
	error("Expected error message about type, got: " .. tostring(err))
end

-- Test with table (invalid)
status, err = pcall(function()
	k8s.parse_int_or_string({})
end)

if status then
	error("Expected error for table input, but got success")
end

if not string.find(err, "expects number or string") then
	error("Expected error message about type, got: " .. tostring(err))
end

return true
