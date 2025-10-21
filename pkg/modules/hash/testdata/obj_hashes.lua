-- Test object hashing functions
local hash = require("hash")

-- Test MD5 object hashing
local h, err = hash.md5_obj({name="John", age=30})
assert(err == nil, "Expected no error, got: " .. tostring(err))
assert(type(h) == "string", "Expected string result")
assert(#h == 32, "Expected 32 character hex string for MD5")

-- Same object should produce same hash
local h2, err2 = hash.md5_obj({name="John", age=30})
assert(err2 == nil, "Expected no error on second call")
assert(h == h2, "Same object should produce same hash")

-- Different object should produce different hash
local h3, err3 = hash.md5_obj({name="Jane", age=30})
assert(err3 == nil, "Expected no error on third call")
assert(h ~= h3, "Different object should produce different hash")

-- Test SHA1 object hashing
local sha1_h, sha1_err = hash.sha1_obj({foo="bar", num=123})
assert(sha1_err == nil, "Expected no error, got: " .. tostring(sha1_err))
assert(type(sha1_h) == "string", "Expected string result")
assert(#sha1_h == 40, "Expected 40 character hex string for SHA1")

-- Test SHA256 object hashing
local sha256_h, sha256_err = hash.sha256_obj({key="value", nested={deep="data"}})
assert(sha256_err == nil, "Expected no error, got: " .. tostring(sha256_err))
assert(type(sha256_h) == "string", "Expected string result")
assert(#sha256_h == 64, "Expected 64 character hex string for SHA256")

-- Test SHA512 object hashing
local sha512_h, sha512_err = hash.sha512_obj({data="test"})
assert(sha512_err == nil, "Expected no error, got: " .. tostring(sha512_err))
assert(type(sha512_h) == "string", "Expected string result")
assert(#sha512_h == 128, "Expected 128 character hex string for SHA512")

return true
