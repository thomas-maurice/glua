-- Test: JSON stringify - nested structure
--
-- Verifies that json.stringify() handles complex
-- nested structures with objects and arrays.

local json = require("json")
local data = {
	user = {
		name = "Bob",
		tags = {"admin", "user"}
	},
	count = 42
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

if tbl.user.name ~= "Bob" then
	error("Expected nested name to be 'Bob' after round-trip")
end

if tbl.user.tags[1] ~= "admin" then
	error("Expected first tag to be 'admin' after round-trip")
end

if tbl.count ~= 42 then
	error("Expected count to be 42 after round-trip")
end
