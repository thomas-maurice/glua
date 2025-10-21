local regexp = require("regexp")

-- Test find_all
local matches, err = regexp.find_all("[0-9]+", "abc123def456ghi", -1)
assert(err == nil, "find_all should not error")
assert(#matches == 2, "should find 2 matches, got: " .. #matches)
assert(matches[1] == "123", "first match should be '123', got: " .. matches[1])
assert(matches[2] == "456", "second match should be '456', got: " .. matches[2])

return true
