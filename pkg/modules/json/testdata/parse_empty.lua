-- Test: JSON parse - empty object and array

local json = require("json")

-- Empty object
local obj = json.parse('{}')
assert(type(obj) == "table", "Expected table for empty object")

-- Empty array
local arr = json.parse('[]')
assert(#arr == 0, "Expected empty array length 0, got: " .. #arr)
