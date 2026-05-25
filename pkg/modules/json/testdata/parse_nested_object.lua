-- Test: JSON parse - nested object

local json = require("json")
local tbl = json.parse('{"person":{"name":"Jane","address":{"city":"NYC"}}}')

assert(tbl.person.name == "Jane", "Expected nested name 'Jane', got: " .. tostring(tbl.person.name))
assert(tbl.person.address.city == "NYC", "Expected nested city 'NYC', got: " .. tostring(tbl.person.address.city))
