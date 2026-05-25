-- Test: JSON stringify - nested structure

local json = require("json")
local data = {user = {name = "Bob", tags = {"admin", "user"}}, count = 42}

local tbl = json.parse(json.stringify(data))

assert(tbl.user.name == "Bob", "Expected nested name 'Bob' after round-trip")
assert(tbl.user.tags[1] == "admin", "Expected first tag 'admin' after round-trip")
assert(tbl.count == 42, "Expected count 42 after round-trip")
