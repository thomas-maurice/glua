-- Test to_upper and to_lower functions
local strings = require("strings")

-- Test to_upper
local text = "hello world"
assert(strings.to_upper(text) == "HELLO WORLD", "to_upper should convert to uppercase")
assert(strings.to_upper("ALREADY UPPER") == "ALREADY UPPER", "to_upper should not change uppercase text")
assert(strings.to_upper("MiXeD CaSe") == "MIXED CASE", "to_upper should convert mixed case")
assert(strings.to_upper("") == "", "to_upper should handle empty string")
assert(strings.to_upper("123abc") == "123ABC", "to_upper should handle alphanumeric")

-- Test to_lower
assert(strings.to_lower("HELLO WORLD") == "hello world", "to_lower should convert to lowercase")
assert(strings.to_lower("already lower") == "already lower", "to_lower should not change lowercase text")
assert(strings.to_lower("MiXeD CaSe") == "mixed case", "to_lower should convert mixed case")
assert(strings.to_lower("") == "", "to_lower should handle empty string")
assert(strings.to_lower("123ABC") == "123abc", "to_lower should handle alphanumeric")

-- Test round-trip conversion
local original = "Hello World 123"
local upper = strings.to_upper(original)
local lower = strings.to_lower(original)
assert(strings.to_lower(upper) == strings.to_lower(original), "round-trip conversion should work")
assert(strings.to_upper(lower) == strings.to_upper(original), "reverse round-trip conversion should work")

return true
