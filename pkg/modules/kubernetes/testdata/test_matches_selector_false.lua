-- Test: matches_selector - negative match
--
-- Verifies that matches_selector returns false when
-- selector labels don't match.

local k8s = require("kubernetes")

local labels = {
	app = "nginx",
	tier = "frontend"
}

local selector = {
	app = "nginx",
	tier = "backend"  -- different value
}

local matches = k8s.matches_selector(labels, selector)
if matches then
	error("Expected selector to NOT match labels")
end

return true
