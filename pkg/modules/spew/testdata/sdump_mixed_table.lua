
-- Test: spew.sdump - mixed table
--
-- Verifies that spew.sdump() handles tables containing
-- mixed types: strings, numbers, booleans, arrays, and nested tables.

local spew = require("spew")
local data = {
	name = "test",
	count = 42,
	active = true,
	tags = {"a", "b", "c"},
	nested = {
		key = "value"
	}
}

local result = spew.sdump(data)

-- Verify various types are present
if not string.find(result, "test") then
	error("Expected result to contain string value")
end

if not string.find(result, "42") then
	error("Expected result to contain number value")
end

return result
