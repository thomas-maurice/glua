
-- Test: JSON round-trip - complex data
--
-- Verifies full round-trip conversion (JSON -> Lua -> JSON -> Lua)
-- with a complex data structure including nested arrays and objects.

local json = require("json")

-- Original JSON
local original = '{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}],"active":true,"count":2}'

-- Parse to table
local tbl, err1 = json.parse(original)
if err1 then
	error("Parse failed: " .. err1)
end

-- Stringify back to JSON
local stringified, err2 = json.stringify(tbl)
if err2 then
	error("Stringify failed: " .. err2)
end

-- Parse again to verify
local final, err3 = json.parse(stringified)
if err3 then
	error("Second parse failed: " .. err3)
end

-- Verify structure
if final.users[1].name ~= "Alice" then
	error("Round-trip failed: name mismatch")
end

if final.users[2].age ~= 25 then
	error("Round-trip failed: age mismatch")
end

if final.active ~= true then
	error("Round-trip failed: active mismatch")
end

if final.count ~= 2 then
	error("Round-trip failed: count mismatch")
end
