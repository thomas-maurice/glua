-- Test: match_gvk - ConfigMap
--
-- Verifies that match_gvk works with ConfigMap (core v1 resource).

local k8s = require("kubernetes")

local configmap = {
	apiVersion = "v1",
	kind = "ConfigMap",
}

local matcher = {group = "", version = "v1", kind = "ConfigMap"}
local matches = k8s.match_gvk(configmap, matcher)
if not matches then
	error("Expected ConfigMap to match GVK v1/ConfigMap")
end

return true
