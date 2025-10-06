-- Test: toleration_matches - Exists operator
--
-- Verifies that a toleration with Exists operator
-- matches any taint with the same key (value doesn't matter).

local k8s = require("kubernetes")

local toleration = {
	key = "node-role",
	operator = "Exists",
	effect = "NoSchedule"
}

local taint = {
	key = "node-role",
	value = "anything",  -- value doesn't matter for Exists
	effect = "NoSchedule"
}

local matches = k8s.toleration_matches(toleration, taint)
if not matches then
	error("Expected toleration with Exists to match taint")
end

return true
