
---@diagnostic disable-next-line: undefined-global
-- Test: parse_cpu with input from Go
--
-- This test is parameterized from Go. The input value is set
-- via L.SetGlobal() before running this script.

local k8s = require("kubernetes")
local result, err = k8s.parse_cpu(test_input)
return result, err
