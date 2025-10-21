-- Test object hashing with arrays and edge cases
local hash = require("hash")

-- Test with array
local h1, err1 = hash.sha256_obj({1, 2, 3})
assert(err1 == nil, "Expected no error for array")
assert(type(h1) == "string", "Expected string result")

-- Same array should produce same hash
local h2, err2 = hash.sha256_obj({1, 2, 3})
assert(err2 == nil, "Expected no error")
assert(h1 == h2, "Same array should produce same hash")

-- Different order should produce different hash
local h3, err3 = hash.sha256_obj({3, 2, 1})
assert(err3 == nil, "Expected no error")
assert(h1 ~= h3, "Different order should produce different hash")

-- Test with empty table
local h_empty, err_empty = hash.sha256_obj({})
assert(err_empty == nil, "Expected no error for empty table")
assert(type(h_empty) == "string", "Expected string result")

return true
