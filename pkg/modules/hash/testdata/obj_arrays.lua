-- Test object hashing with arrays and edge cases (raise-on-error shape)
local hash = require("hash")

-- Test with array
local h1 = hash.sha256_obj({1, 2, 3})
assert(type(h1) == "string", "Expected string result")

-- Same array should produce same hash
local h2 = hash.sha256_obj({1, 2, 3})
assert(h1 == h2, "Same array should produce same hash")

-- Different order should produce different hash
local h3 = hash.sha256_obj({3, 2, 1})
assert(h1 ~= h3, "Different order should produce different hash")

-- Test with empty table
local h_empty = hash.sha256_obj({})
assert(type(h_empty) == "string", "Expected string result for empty table")

return true
