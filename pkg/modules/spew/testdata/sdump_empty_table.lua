-- Test: spew.sdump - empty table
--
-- Verifies that spew.sdump() handles empty tables correctly.

local spew = require("spew")
local result = spew.sdump({})

if #result == 0 then
	error("Expected non-empty result for empty table")
end

return result
