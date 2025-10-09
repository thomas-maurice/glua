
-- Test: JSON stringify - array
--
-- Verifies that json.stringify() correctly converts
-- Lua arrays to JSON arrays.

local json = require("json")
local str, err = json.stringify({1, 2, 3, 4, 5})

if err then
	error("Stringify failed: " .. err)
end

-- Parse it back to verify
local tbl, parseErr = json.parse(str)
if parseErr then
	error("Failed to parse stringified JSON: " .. parseErr)
end

if #tbl ~= 5 then
	error("Expected array length 5 after round-trip, got: " .. #tbl)
end

if tbl[1] ~= 1 or tbl[5] ~= 5 then
	error("Array values incorrect after round-trip")
end
