-- Test: spew.sdump - a self-referential table raises a catchable error
-- instead of crashing the process, and the Lua state keeps working
-- afterwards.
--
-- Regression test: spew's dump/sdump hand-roll their own unguarded recursive
-- Lua->Go conversion (luaToGo), so handing them a cyclic table (`t.self =
-- t`) recursed forever and killed the whole host process with an
-- unrecoverable "fatal error: stack overflow" -- not something pcall could
-- ever catch. They now share pkg/glua's TableGuard, so the same input
-- raises an ordinary, catchable Lua error.

local spew = require("spew")

local t = {}
t.self = t

local ok, err = pcall(spew.sdump, t)
assert(not ok, "Expected error for a self-referential table")
assert(type(err) == "string", "Error should be a string message")
assert(string.find(err, "cycle") ~= nil, "Expected error to mention 'cycle', got: " .. err)

local ok2, err2 = pcall(spew.dump, t)
assert(not ok2, "Expected error from dump for a self-referential table")
assert(string.find(err2, "cycle") ~= nil, "Expected dump error to mention 'cycle', got: " .. err2)

-- Prove the Lua state is still usable after the caught errors.
local result = spew.sdump({hello = "world"})
assert(string.find(result, "hello") ~= nil, "Expected normal sdump to still work, got: " .. result)
