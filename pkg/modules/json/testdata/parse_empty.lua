-- Test: JSON parse - empty object and array
--
-- Verifies that json.parse() correctly handles
-- empty objects and empty arrays.

local json = require("json")

-- Empty object
local obj, err1 = json.parse('{}')
if err1 then
	error("Parse empty object failed: " .. err1)
end

-- Empty array
local arr, err2 = json.parse('[]')
if err2 then
	error("Parse empty array failed: " .. err2)
end

if #arr ~= 0 then
	error("Expected empty array length to be 0, got: " .. #arr)
end
