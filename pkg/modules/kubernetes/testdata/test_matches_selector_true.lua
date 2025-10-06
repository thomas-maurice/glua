-- Test: matches_selector - positive match
--
-- Verifies that matches_selector returns true when
-- all selector labels are present in the label set.

local k8s = require("kubernetes")

local labels = {
	app = "nginx",
	tier = "frontend",
	environment = "production"
}

local selector = {
	app = "nginx",
	tier = "frontend"
}

local matches = k8s.matches_selector(labels, selector)
if not matches then
	error("Expected selector to match labels")
end

return true
