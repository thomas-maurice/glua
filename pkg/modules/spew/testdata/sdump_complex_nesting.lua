
-- Test: spew.sdump - complex nesting
--
-- Verifies that spew.sdump() handles deeply nested structures
-- with multiple levels of nesting.

local spew = require("spew")
local data = {
	level1 = {
		level2 = {
			level3 = {
				level4 = {
					deep_value = "found"
				}
			}
		}
	}
}

local result = spew.sdump(data)

if not string.find(result, "found") then
	error("Expected result to contain deeply nested value")
end

return result
