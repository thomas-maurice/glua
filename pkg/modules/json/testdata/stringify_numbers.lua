-- Test: JSON stringify - numbers
--
-- Verifies that json.stringify() correctly handles
-- various number types (integers, negatives, floats, zero).

local json = require("json")
local data = {
	integer = 42,
	negative = -10,
	float = 3.14,
	zero = 0
}

local str, err = json.stringify(data)
if err then
	error("Stringify failed: " .. err)
end

-- Parse it back to verify
local tbl, parseErr = json.parse(str)
if parseErr then
	error("Failed to parse stringified JSON: " .. parseErr)
end

if tbl.integer ~= 42 then
	error("Expected integer to be 42 after round-trip")
end

if tbl.negative ~= -10 then
	error("Expected negative to be -10 after round-trip")
end

if math.abs(tbl.float - 3.14) > 0.01 then
	error("Expected float to be ~3.14 after round-trip")
end

if tbl.zero ~= 0 then
	error("Expected zero to be 0 after round-trip")
end
