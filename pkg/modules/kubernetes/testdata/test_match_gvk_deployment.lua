-- Test: match_gvk - apps/v1 Deployment
--
-- Verifies that match_gvk correctly identifies a Deployment object.

local k8s = require("kubernetes")

local deployment = {
	apiVersion = "apps/v1",
	kind = "Deployment",
}

local matcher = {group = "apps", version = "v1", kind = "Deployment"}
local matches = k8s.match_gvk(deployment, matcher)
if not matches then
	error("Expected Deployment to match GVK apps/v1/Deployment")
end

return true
