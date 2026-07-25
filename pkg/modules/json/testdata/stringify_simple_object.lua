-- Test: JSON stringify - simple object

local json = require("json")
local tbl = json.parse(json.stringify({name="Alice", age=25}))

assert(tbl.name == "Alice", "Expected name 'Alice' after round-trip")
assert(tbl.age == 25, "Expected age 25 after round-trip")
