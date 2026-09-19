-- Test rep. Named rep, not repeat: "repeat" is a reserved word in Lua, so
-- strings.repeat(s, n) would be a SYNTAX error, not a runtime one -- this
-- test file existing and being parseable at all is part of the proof, and
-- the explicit strings.rep(...) call below proves it is callable.
local strings = require("strings")

assert(strings.rep("ab", 3) == "ababab", "rep should repeat the string n times")
assert(strings.rep("x", 0) == "", "rep with count 0 should return empty string")
assert(strings.rep("", 5) == "", "rep of an empty string should return empty string")
assert(strings.rep("ab", 1) == "ab", "rep with count 1 should return s unchanged")

-- Negative count is a programmer error: Go's strings.Repeat panics on this,
-- which this module guards against and turns into a catchable Lua error.
local ok, err = pcall(strings.rep, "ab", -1)
assert(ok == false, "rep with a negative count should raise")
assert(string.find(err, "count must be >= 0") ~= nil, "rep error message should explain the constraint, got: " .. tostring(err))

return true
