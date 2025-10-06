-- Test: toleration_matches - Equal operator
--
-- Verifies that a toleration with Equal operator
-- matches a taint with the same key/value/effect.

local k8s = require("kubernetes")

local toleration = {
	key = "node-role",
	operator = "Equal",
	value = "master",
	effect = "NoSchedule"
}

local taint = {
	key = "node-role",
	value = "master",
	effect = "NoSchedule"
}

local matches = k8s.toleration_matches(toleration, taint)
if not matches then
	error("Expected toleration to match taint")
end

return true
