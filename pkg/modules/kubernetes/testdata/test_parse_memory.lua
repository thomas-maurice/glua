
---@diagnostic disable-next-line: undefined-global
-- Test: parse_memory with input from Go
--
-- This test is parameterized from Go. The input value is set
-- via L.SetGlobal() before running this script.
-- With luareg auto-raise, errors are raised as Lua errors (no 2nd return value).

local k8s = require("kubernetes")
local result = k8s.parse_memory(test_input)
return result
