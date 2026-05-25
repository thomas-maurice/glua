
-- Test: match_gvk validation - missing version
--
-- Verifies that match_gvk returns false when version field is missing
-- (an incomplete matcher cannot match any resource).

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

-- Missing 'version' field: incomplete matcher should return false
local matcher = {group = "", kind = "Pod"}
local result = k8s.match_gvk(pod, matcher)

if result then
	error("Expected false for incomplete matcher (missing 'version' field)")
end

return true
