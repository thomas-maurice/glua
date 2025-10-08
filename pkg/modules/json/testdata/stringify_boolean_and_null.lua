
-- Test: JSON stringify - boolean and null values
--
-- Verifies that json.stringify() correctly handles
-- boolean values and nil (which becomes null in JSON).

local json = require("json")
local data = {
	active = true,
	disabled = false,
	value = nil
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

if tbl.active ~= true then
	error("Expected active to be true after round-trip")
end

if tbl.disabled ~= false then
	error("Expected disabled to be false after round-trip")
end
