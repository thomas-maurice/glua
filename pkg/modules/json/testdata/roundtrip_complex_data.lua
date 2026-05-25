-- Test: JSON round-trip - complex data

local json = require("json")

local original = '{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}],"active":true,"count":2}'

local tbl = json.parse(original)
local stringified = json.stringify(tbl)
local final = json.parse(stringified)

assert(final.users[1].name == "Alice", "Round-trip failed: name mismatch")
assert(final.users[2].age == 25, "Round-trip failed: age mismatch")
assert(final.active == true, "Round-trip failed: active mismatch")
assert(final.count == 2, "Round-trip failed: count mismatch")
