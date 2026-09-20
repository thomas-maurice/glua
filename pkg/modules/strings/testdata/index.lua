-- Test index and last_index: per SPECS.md decision D2 these are 1-based,
-- returning 0 when absent -- NOT Go's 0-based/-1 convention. This test pins
-- that deviation and its composability with string.sub.
local strings = require("strings")

-- Basic found case: "world" starts at Lua-index 7 in "hello world".
assert(strings.index("hello world", "world") == 7, "index should be 1-based")
assert(strings.index("hello world", "hello") == 1, "index at the very start should be 1, not 0")

-- Composes with string.sub without an off-by-one adjustment.
local s = "key=value"
local pos = strings.index(s, "=")
assert(pos == 4, "index should find the separator at position 4")
assert(string.sub(s, 1, pos - 1) == "key", "index result should compose with string.sub for the key")
assert(string.sub(s, pos + 1) == "value", "index result should compose with string.sub for the value")

-- Absent case returns 0, not -1 and not nil.
assert(strings.index("hello world", "xyz") == 0, "index should return 0 when substr is absent")
assert(strings.index("", "x") == 0, "index on empty string should return 0")

-- last_index: same 1-based/0-absent convention, but for the last occurrence.
assert(strings.last_index("banana", "an") == 4, "last_index should find the last occurrence, 1-based")
assert(strings.last_index("banana", "xyz") == 0, "last_index should return 0 when substr is absent")
assert(strings.last_index("banana", "b") == 1, "last_index at the very start should be 1")

return true
