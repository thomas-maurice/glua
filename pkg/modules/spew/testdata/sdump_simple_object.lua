
-- Test: spew.sdump - simple object
--
-- Verifies that spew.sdump() returns a string representation
-- of a simple table containing basic key-value pairs.

local spew = require("spew")
local result = spew.sdump({name="John", age=30})

-- Check that result contains expected values
if not string.find(result, "name") then
	error("Expected result to contain 'name'")
end

if not string.find(result, "John") then
	error("Expected result to contain 'John'")
end

if not string.find(result, "age") then
	error("Expected result to contain 'age'")
end

if not string.find(result, "30") then
	error("Expected result to contain '30'")
end

return result
