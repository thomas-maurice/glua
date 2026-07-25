local regexp = require("regexp")

-- match: raises on invalid pattern, returns bool on valid pattern
local matched = regexp.match("^hello", "hello world")
assert(matched == true, "should match")

local matched2 = regexp.match("^world", "hello world")
assert(matched2 == false, "should not match")

-- find: raises on invalid pattern, returns match string on valid pattern
local result = regexp.find("world", "hello world")
assert(result == "world", "should find 'world', got: " .. result)

-- find on no match returns empty string
local empty = regexp.find("xyz", "hello world")
assert(empty == "", "no match should return empty string, got: " .. empty)

return true
