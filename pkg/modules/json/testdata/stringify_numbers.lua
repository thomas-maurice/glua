-- Test: JSON stringify - numbers

local json = require("json")
local data = {integer = 42, negative = -10, float = 3.14, zero = 0}

local tbl = json.parse(json.stringify(data))

assert(tbl.integer == 42, "Expected integer 42 after round-trip")
assert(tbl.negative == -10, "Expected negative -10 after round-trip")
assert(math.abs(tbl.float - 3.14) < 0.01, "Expected float ~3.14 after round-trip")
assert(tbl.zero == 0, "Expected zero 0 after round-trip")
