-- Test: spew.sdump - nested structure
--
-- Verifies that spew.sdump() handles nested tables correctly
-- and includes all nested values in the output.

local spew = require("spew")
local data = {
	person = {
		name = "Alice",
		address = {
			city = "NYC",
			zip = "10001"
		}
	},
	items = {1, 2, 3}
}

local result = spew.sdump(data)

-- Check nested values are present
if not string.find(result, "Alice") then
	error("Expected result to contain 'Alice'")
end

if not string.find(result, "NYC") then
	error("Expected result to contain 'NYC'")
end

return result
