-- Test: parse_int_or_string with input from Go
--
-- This test is parameterized from Go. The input value is set
-- via L.SetGlobal() before running this script.

local k8s = require("kubernetes")
local value, is_string = k8s.parse_int_or_string(test_input)
return value, is_string
