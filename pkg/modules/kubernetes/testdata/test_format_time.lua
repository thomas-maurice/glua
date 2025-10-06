-- Test: format_time with input from Go
--
-- This test is parameterized from Go. The timestamp value is set
-- via L.SetGlobal() before running this script.

local k8s = require("kubernetes")
local result, err = k8s.format_time(test_timestamp)
return result, err
