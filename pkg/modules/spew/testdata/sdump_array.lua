
-- Test: spew.sdump - array
--
-- Verifies that spew.sdump() correctly handles Lua arrays.

local spew = require("spew")
local result = spew.sdump({1, 2, 3, 4, 5})

-- Check that result is non-empty
if #result == 0 then
	error("Expected non-empty result")
end

return result
