-- Test: spew.sdump - number keys
--
-- Verifies that spew.sdump() handles tables with
-- non-consecutive numeric keys (treated as maps).

local spew = require("spew")
-- Non-consecutive numeric keys should be treated as map
local data = {
	[1] = "first",
	[5] = "fifth",
	[10] = "tenth"
}

local result = spew.sdump(data)

if not string.find(result, "first") then
	error("Expected result to contain 'first'")
end

if not string.find(result, "fifth") then
	error("Expected result to contain 'fifth'")
end

return result
