
-- Test: JSON parse - simple object
--
-- Verifies that json.parse() can parse a simple JSON object
-- with string, number, and boolean fields.

local json = require("json")
local tbl, err = json.parse('{"name":"John","age":30,"active":true}')

if err then
	error("Parse failed: " .. err)
end

if tbl.name ~= "John" then
	error("Expected name to be 'John', got: " .. tostring(tbl.name))
end

if tbl.age ~= 30 then
	error("Expected age to be 30, got: " .. tostring(tbl.age))
end

if tbl.active ~= true then
	error("Expected active to be true, got: " .. tostring(tbl.active))
end
