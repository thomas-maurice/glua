-- Test: spew.sdump - array of objects
--
-- Verifies that spew.sdump() handles arrays containing
-- object elements (tables with named fields).

local spew = require("spew")
local users = {
	{name="Alice", age=30},
	{name="Bob", age=25},
	{name="Charlie", age=35}
}

local result = spew.sdump(users)

-- Check all names are present
if not string.find(result, "Alice") then
	error("Expected result to contain 'Alice'")
end

if not string.find(result, "Bob") then
	error("Expected result to contain 'Bob'")
end

if not string.find(result, "Charlie") then
	error("Expected result to contain 'Charlie'")
end

return result
