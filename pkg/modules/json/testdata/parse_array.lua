-- Test: JSON parse - array

local json = require("json")
local tbl = json.parse('[1,2,3,4,5]')

assert(#tbl == 5, "Expected array length 5, got: " .. #tbl)
assert(tbl[1] == 1 and tbl[5] == 5, "Array values incorrect")
