-- Test: JSON parse - array
--
-- Verifies that json.parse() can parse a JSON array
-- and correctly maps it to a Lua table with numeric indices.

local json = require("json")
local tbl, err = json.parse('[1,2,3,4,5]')

if err then
	error("Parse failed: " .. err)
end

if #tbl ~= 5 then
	error("Expected array length 5, got: " .. #tbl)
end

if tbl[1] ~= 1 or tbl[5] ~= 5 then
	error("Array values incorrect")
end
