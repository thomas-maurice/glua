
-- Test: match_gvk validation - missing kind
--
-- Verifies that match_gvk errors when kind field is missing.

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

-- Missing 'kind' field should error
local matcher = {group = "", version = "v1"}
local status, err = pcall(function()
	k8s.match_gvk(pod, matcher)
end)

if status then
	error("Expected error for missing 'kind' field, but got success")
end

if not string.find(err, "requires 'kind' field") then
	error("Expected error message about 'kind' field, got: " .. tostring(err))
end

return true
