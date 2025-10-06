-- Test: match_gvk - core v1 Pod
--
-- Verifies that match_gvk correctly identifies a Pod object.

local k8s = require("kubernetes")

local pod = {
	apiVersion = "v1",
	kind = "Pod",
}

local matcher = {group = "", version = "v1", kind = "Pod"}
local matches = k8s.match_gvk(pod, matcher)
if not matches then
	error("Expected Pod to match GVK v1/Pod")
end

return true
