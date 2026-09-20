-- Test: JSON stringify - a self-referential table raises a catchable error
-- instead of crashing the process, and the Lua state keeps working
-- afterwards.
--
-- Regression test: json.stringify used to hand-roll its own unguarded
-- recursive Lua->Go conversion (luaToGo), so handing it a cyclic table
-- (`t.self = t`) recursed forever and killed the whole host process with an
-- unrecoverable "fatal error: stack overflow" -- not something pcall could
-- ever catch. It now shares pkg/glua's TableGuard, so the same input raises
-- an ordinary, catchable Lua error.

local json = require("json")

local t = {}
t.self = t

local ok, err = pcall(json.stringify, t)
assert(not ok, "Expected error for a self-referential table")
assert(type(err) == "string", "Error should be a string message")
assert(string.find(err, "cycle") ~= nil, "Expected error to mention 'cycle', got: " .. err)

-- Prove the Lua state is still usable after the caught error.
local result = json.stringify({hello = "world"})
assert(result == '{"hello":"world"}', "Expected normal stringify to still work, got: " .. result)
