-- random.int(min, max) requires min <= max; random.string requires a
-- non-empty charset. Both must raise a catchable error, not silently
-- misbehave.
local random = require("random")

local ok1, err1 = pcall(random.int, 5, 1)
assert(not ok1, "expected random.int(5, 1) to raise")
assert(string.find(err1, "min"), "error should mention min, got: " .. tostring(err1))

local ok2, err2 = pcall(random.string, 4, "")
assert(not ok2, "expected random.string(4, \"\") to raise")
assert(string.find(err2, "charset"), "error should mention charset, got: " .. tostring(err2))

local ok3 = pcall(random.bytes, 0)
assert(not ok3, "expected random.bytes(0) to raise")

local ok4 = pcall(random.hex, -1)
assert(not ok4, "expected random.hex(-1) to raise")

return true
