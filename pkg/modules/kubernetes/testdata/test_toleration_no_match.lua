-- Test: toleration_matches - no match
--
-- Verifies that tolerations don't match when
-- keys or effects differ.

local k8s = require("kubernetes")

local toleration = {
	key = "node-role",
	operator = "Equal",
	value = "master",
	effect = "NoSchedule"
}

local taint = {
	key = "different-key",
	value = "master",
	effect = "NoSchedule"
}

local matches = k8s.toleration_matches(toleration, taint)
if matches then
	error("Expected toleration to NOT match taint (different key)")
end

return true
