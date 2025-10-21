local regexp = require("regexp")

-- Test match
local matched, err = regexp.match("^hello", "hello world")
assert(err == nil, "match should not error")
assert(matched == true, "should match")

local matched2, err2 = regexp.match("^world", "hello world")
assert(err2 == nil, "match should not error")
assert(matched2 == false, "should not match")

-- Test find
local result, err3 = regexp.find("world", "hello world")
assert(err3 == nil, "find should not error")
assert(result == "world", "should find 'world', got: " .. result)

return true
