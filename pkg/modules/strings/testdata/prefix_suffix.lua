-- Test has_prefix and has_suffix functions
local strings = require("strings")

-- Test has_prefix
local text = "hello world"
assert(strings.has_prefix(text, "hello") == true, "has_prefix should return true for matching prefix")
assert(strings.has_prefix(text, "world") == false, "has_prefix should return false for non-matching prefix")
assert(strings.has_prefix(text, "") == true, "has_prefix should return true for empty prefix")
assert(strings.has_prefix("", "hello") == false, "has_prefix should return false for empty string")

-- Test has_suffix
assert(strings.has_suffix(text, "world") == true, "has_suffix should return true for matching suffix")
assert(strings.has_suffix(text, "hello") == false, "has_suffix should return false for non-matching suffix")
assert(strings.has_suffix(text, "") == true, "has_suffix should return true for empty suffix")
assert(strings.has_suffix("", "world") == false, "has_suffix should return false for empty string")

return true
