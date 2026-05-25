-- Test: JSON stringify - array

local json = require("json")
local str = json.stringify({1, 2, 3, 4, 5})

local tbl = json.parse(str)
assert(#tbl == 5, "Expected array length 5 after round-trip, got: " .. #tbl)
assert(tbl[1] == 1 and tbl[5] == 5, "Array values incorrect after round-trip")
