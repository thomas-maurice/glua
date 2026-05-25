-- Test: JSON parse - simple object

local json = require("json")
local tbl = json.parse('{"name":"John","age":30,"active":true}')

assert(tbl.name == "John", "Expected name 'John', got: " .. tostring(tbl.name))
assert(tbl.age == 30, "Expected age 30, got: " .. tostring(tbl.age))
assert(tbl.active == true, "Expected active true, got: " .. tostring(tbl.active))
