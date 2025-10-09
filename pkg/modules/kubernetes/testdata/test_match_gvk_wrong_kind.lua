
-- Test: match_gvk - wrong kind
--
-- Verifies that match_gvk returns false when the kind doesn't match.

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

-- Try to match as a Service (wrong kind)
local matcher = {group = "", version = "v1", kind = "Service"}
local matches = k8s.match_gvk(pod, matcher)
if matches then
	error("Expected Pod NOT to match GVK v1/Service")
end

return true
