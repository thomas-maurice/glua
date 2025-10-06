-- Test: matches_selector with nil values
--
-- Verifies that matches_selector skips nil values in selector.

local k8s = require("kubernetes")

local labels = {app = "test", tier = "frontend"}
local selector = {app = nil, tier = "frontend"}

-- Should match because nil values are skipped
local matches = k8s.matches_selector(labels, selector)
if not matches then
	error("Expected matches_selector to skip nil values and return true")
end

return true
