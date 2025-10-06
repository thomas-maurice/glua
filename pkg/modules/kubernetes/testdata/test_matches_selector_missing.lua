-- Test: matches_selector - missing label
--
-- Verifies that matches_selector returns false when
-- a selector label is missing from the label set.

local k8s = require("kubernetes")

local labels = {
	app = "nginx"
}

local selector = {
	app = "nginx",
	tier = "frontend"  -- missing from labels
}

local matches = k8s.matches_selector(labels, selector)
if matches then
	error("Expected selector to NOT match labels (missing label)")
end

return true
