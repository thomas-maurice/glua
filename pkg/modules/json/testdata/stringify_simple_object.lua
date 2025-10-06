-- Test: JSON stringify - simple object
--
-- Verifies that json.stringify() can convert a Lua table
-- to JSON and that it can be parsed back correctly.

local json = require("json")
local str, err = json.stringify({name="Alice", age=25})

if err then
	error("Stringify failed: " .. err)
end

-- Parse it back to verify
local tbl, parseErr = json.parse(str)
if parseErr then
	error("Failed to parse stringified JSON: " .. parseErr)
end

if tbl.name ~= "Alice" then
	error("Expected name to be 'Alice' after round-trip")
end

if tbl.age ~= 25 then
	error("Expected age to be 25 after round-trip")
end
