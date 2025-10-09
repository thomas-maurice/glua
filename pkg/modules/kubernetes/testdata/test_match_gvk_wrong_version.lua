
-- Test: match_gvk - wrong version
--
-- Verifies that match_gvk returns false when the version doesn't match.

local k8s = require("kubernetes")

local deployment = {
	apiVersion = "apps/v1",
	kind = "Deployment",
}

-- Try to match with wrong version
local matcher = {group = "apps", version = "v1beta1", kind = "Deployment"}
local matches = k8s.match_gvk(deployment, matcher)
if matches then
	error("Expected Deployment NOT to match GVK apps/v1beta1/Deployment")
end

return true
