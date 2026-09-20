-- Test text.pad_left and text.pad_right.
local text = require("text")

assert(text.pad_left("42", 5, "0") == "00042", "pad_left should left-pad to width")
assert(text.pad_right("42", 5, "0") == "42000", "pad_right should right-pad to width")

-- Never truncates: already-wider input is returned unchanged.
assert(text.pad_left("hello world", 3, " ") == "hello world", "pad_left must not truncate a string already wider than width")
assert(text.pad_right("hello world", 3, " ") == "hello world", "pad_right must not truncate a string already wider than width")

-- Exactly at width: unchanged, no padding added.
assert(text.pad_left("abc", 3, "-") == "abc", "pad_left at exactly width should add no padding")

-- Empty string pads to a string of just the pad rune.
assert(text.pad_left("", 3, "*") == "***", "pad_left of an empty string should be all pad runes")

-- pad must be exactly one rune: multi-rune pad raises.
local ok, err = pcall(text.pad_left, "x", 5, "ab")
assert(ok == false, "pad_left should raise when pad is more than one rune")
assert(string.find(err, "exactly one rune") ~= nil, "pad_left error should explain the constraint, got: " .. tostring(err))

ok, err = pcall(text.pad_right, "x", 5, "")
assert(ok == false, "pad_right should raise when pad is the empty string (zero runes)")

-- width is capped so a typo'd huge width cannot turn into an unbounded
-- allocation (security review LOW 3): a value far past the cap must raise,
-- not silently allocate.
local okHuge, errHuge = pcall(text.pad_left, "x", 1e8, "-")
assert(okHuge == false, "pad_left with a width far past the cap should raise")
assert(string.find(errHuge, "width must be") ~= nil, "pad_left error should name the width constraint, got: " .. tostring(errHuge))

return true
