-- Test contains, count, and replace functions
local strings = require("strings")

-- Test contains
local text = "hello world"
assert(strings.contains(text, "world") == true, "contains should find substring")
assert(strings.contains(text, "hello") == true, "contains should find substring at start")
assert(strings.contains(text, "missing") == false, "contains should return false for missing substring")
assert(strings.contains(text, "") == true, "contains should return true for empty substring")
assert(strings.contains("", "hello") == false, "contains should return false for empty string")
assert(strings.contains(text, "WORLD") == false, "contains should be case-sensitive")

-- Test count
assert(strings.count("hello", "l") == 2, "count should count occurrences")
assert(strings.count("banana", "a") == 3, "count should count multiple occurrences")
assert(strings.count("hello world", "o") == 2, "count should count in longer string")
assert(strings.count("hello", "x") == 0, "count should return 0 for missing substring")
assert(strings.count("", "a") == 0, "count should return 0 for empty string")
assert(strings.count("aaa", "aa") == 1, "count should count non-overlapping matches")

-- Test replace
local replaced = strings.replace("hello world", "world", "there", -1)
assert(replaced == "hello there", "replace should replace substring")

replaced = strings.replace("foo bar foo", "foo", "baz", 1)
assert(replaced == "baz bar foo", "replace with n=1 should replace only first occurrence")

replaced = strings.replace("foo bar foo", "foo", "baz", 2)
assert(replaced == "baz bar baz", "replace with n=2 should replace two occurrences")

replaced = strings.replace("foo bar foo", "foo", "baz", -1)
assert(replaced == "baz bar baz", "replace with n=-1 should replace all occurrences")

replaced = strings.replace("hello", "x", "y", -1)
assert(replaced == "hello", "replace should not change string without substring")

replaced = strings.replace("", "a", "b", -1)
assert(replaced == "", "replace should handle empty string")

-- Test replace with empty substring
replaced = strings.replace("abc", "", "x", -1)
assert(replaced == "xaxbxcx", "replace empty substring should insert between each char")

return true
