-- Test trim, trim_left, and trim_right functions
local strings = require("strings")

-- Test trim
local text = "  hello world  "
assert(strings.trim(text, " ") == "hello world", "trim should remove spaces from both sides")
assert(strings.trim("xxxhello worldxxx", "x") == "hello world", "trim should remove custom cutset from both sides")
assert(strings.trim("hello", "x") == "hello", "trim should not modify string without cutset chars")
assert(strings.trim("", " ") == "", "trim should handle empty string")

-- Test trim_left
assert(strings.trim_left(text, " ") == "hello world  ", "trim_left should remove spaces from left only")
assert(strings.trim_left("xxxhello worldxxx", "x") == "hello worldxxx", "trim_left should remove custom cutset from left")
assert(strings.trim_left("hello", "x") == "hello", "trim_left should not modify string without cutset chars")

-- Test trim_right
assert(strings.trim_right(text, " ") == "  hello world", "trim_right should remove spaces from right only")
assert(strings.trim_right("xxxhello worldxxx", "x") == "xxxhello world", "trim_right should remove custom cutset from right")
assert(strings.trim_right("hello", "x") == "hello", "trim_right should not modify string without cutset chars")

-- Test with multiple characters in cutset
local multi = "___---hello world---___"
assert(strings.trim(multi, "_-") == "hello world", "trim should handle multiple chars in cutset")
assert(strings.trim_left(multi, "_-") == "hello world---___", "trim_left should handle multiple chars in cutset")
assert(strings.trim_right(multi, "_-") == "___---hello world", "trim_right should handle multiple chars in cutset")

return true
