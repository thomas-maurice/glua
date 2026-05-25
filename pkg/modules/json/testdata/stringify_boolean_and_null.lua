-- Test: JSON stringify - boolean and null values

local json = require("json")
local data = {active = true, disabled = false, value = nil}

local str = json.stringify(data)
local tbl = json.parse(str)

assert(tbl.active == true, "Expected active true after round-trip")
assert(tbl.disabled == false, "Expected disabled false after round-trip")
