-- Test: integration test for new kubernetes functions
--
-- Verifies that all new helper functions work correctly together.

local k8s = require("kubernetes")

-- Test parse_duration
local seconds, err = k8s.parse_duration("5m")
if err then error("parse_duration failed: " .. err) end
if seconds ~= 300 then error("Expected 5m to be 300 seconds") end

-- Test format_duration
local duration_str, err2 = k8s.format_duration(300)
if err2 then error("format_duration failed: " .. err2) end
if duration_str ~= "5m0s" then error("Expected 300 seconds to be 5m0s") end

-- Test parse_int_or_string with number
local val1, is_str1 = k8s.parse_int_or_string(8080)
if is_str1 then error("Expected number to not be string") end
if val1 ~= 8080 then error("Expected value to be 8080") end

-- Test parse_int_or_string with string
local val2, is_str2 = k8s.parse_int_or_string("http")
if not is_str2 then error("Expected string to be string") end
if val2 ~= "http" then error("Expected value to be 'http'") end

-- Test matches_selector
local matches1 = k8s.matches_selector(
	{app="nginx", tier="frontend"},
	{app="nginx"}
)
if not matches1 then error("Expected selector to match") end

-- Test toleration_matches
local matches2 = k8s.toleration_matches(
	{key="node-role", operator="Equal", value="master", effect="NoSchedule"},
	{key="node-role", value="master", effect="NoSchedule"}
)
if not matches2 then error("Expected toleration to match taint") end

return true
