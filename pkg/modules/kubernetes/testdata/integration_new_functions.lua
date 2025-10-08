
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

return true
