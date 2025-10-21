local regexp = require("regexp")

-- Test replace_all
local result, err = regexp.replace_all("[0-9]+", "abc123def456", "NUM")
assert(err == nil, "replace_all should not error")
assert(result == "abcNUMdefNUM", "should replace all numbers, got: " .. result)

-- Test split
local parts, err2 = regexp.split(",", "a,b,c,d", -1)
assert(err2 == nil, "split should not error")
assert(#parts == 4, "should have 4 parts, got: " .. #parts)
assert(parts[1] == "a" and parts[4] == "d", "parts should be correct")

return true
