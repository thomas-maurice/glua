local regexp = require("regexp")

-- find_all: returns all matches as a table
local matches = regexp.find_all("[0-9]+", "abc123def456ghi", -1)
assert(#matches == 2, "should find 2 matches, got: " .. #matches)
assert(matches[1] == "123", "first match should be '123', got: " .. matches[1])
assert(matches[2] == "456", "second match should be '456', got: " .. matches[2])

-- find_all with invalid pattern raises
local ok, err = pcall(regexp.find_all, "[", "text", -1)
assert(not ok, "find_all should raise on invalid pattern")

return true
