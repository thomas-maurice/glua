-- Test: format_duration with input from Go
--
-- This test is parameterized from Go. The seconds value is set
-- via L.SetGlobal() before running this script.

local k8s = require("kubernetes")
local result, err = k8s.format_duration(test_seconds)
return result, err
