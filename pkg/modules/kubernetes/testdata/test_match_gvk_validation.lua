
-- Test: match_gvk validation - missing kind
--
-- Verifies that match_gvk returns false when kind field is missing
-- (an incomplete matcher cannot match any resource).

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

-- Missing 'kind' field: incomplete matcher should return false
local matcher = {group = "", version = "v1"}
local result = k8s.match_gvk(pod, matcher)

if result then
	error("Expected false for incomplete matcher (missing 'kind' field)")
end

return true
