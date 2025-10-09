
-- Test: match_gvk validation - missing version
--
-- Verifies that match_gvk errors when version field is missing.

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

-- Missing 'version' field should error
local matcher = {group = "", kind = "Pod"}
local status, err = pcall(function()
	k8s.match_gvk(pod, matcher)
end)

if status then
	error("Expected error for missing 'version' field, but got success")
end

if not string.find(err, "requires 'version' field") then
	error("Expected error message about 'version' field, got: " .. tostring(err))
end

return true
