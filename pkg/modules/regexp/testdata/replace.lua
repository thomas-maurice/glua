local regexp = require("regexp")

-- replace_all: replaces all matches
local result = regexp.replace_all("[0-9]+", "abc123def456", "NUM")
assert(result == "abcNUMdefNUM", "should replace all numbers, got: " .. result)

-- replace: replaces only the first match
local first = regexp.replace("[0-9]+", "abc123def456", "NUM")
assert(first == "abcNUMdef456", "should replace only first match, got: " .. first)

-- replace with no match returns text unchanged
local same = regexp.replace("[0-9]+", "no digits here", "NUM")
assert(same == "no digits here", "should return input unchanged, got: " .. same)

-- replace with invalid pattern raises — catch with pcall
local ok, err = pcall(regexp.replace, "[", "text", "X")
assert(not ok, "replace should surface compile errors")

-- split: splits text by pattern
local parts = regexp.split(",", "a,b,c,d", -1)
assert(#parts == 4, "should have 4 parts, got: " .. #parts)
assert(parts[1] == "a" and parts[4] == "d", "parts should be correct")

return true
