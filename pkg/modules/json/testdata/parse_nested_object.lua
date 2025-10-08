
-- Test: JSON parse - nested object
--
-- Verifies that json.parse() can handle nested objects
-- with multiple levels of depth.

local json = require("json")
local tbl, err = json.parse('{"person":{"name":"Jane","address":{"city":"NYC"}}}')

if err then
	error("Parse failed: " .. err)
end

if tbl.person.name ~= "Jane" then
	error("Expected nested name to be 'Jane', got: " .. tostring(tbl.person.name))
end

if tbl.person.address.city ~= "NYC" then
	error("Expected nested city to be 'NYC', got: " .. tostring(tbl.person.address.city))
end
