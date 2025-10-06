-- Test: spew.sdump - primitive types
--
-- Verifies that spew.sdump() handles all Lua primitive types:
-- strings, numbers, booleans, and nil.

local spew = require("spew")

-- Test string
local str = spew.sdump("hello")
if #str == 0 then
	error("Expected non-empty result for string")
end

-- Test number
local num = spew.sdump(42)
if #num == 0 then
	error("Expected non-empty result for number")
end

-- Test boolean
local bool = spew.sdump(true)
if #bool == 0 then
	error("Expected non-empty result for boolean")
end

-- Test nil
local nil_val = spew.sdump(nil)
if #nil_val == 0 then
	error("Expected non-empty result for nil")
end
