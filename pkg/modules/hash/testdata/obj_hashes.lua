-- Test object hashing functions (raise-on-error shape)
local hash = require("hash")

-- Test MD5 object hashing
local h = hash.md5_obj({name="John", age=30})
assert(type(h) == "string", "Expected string result")
assert(#h == 32, "Expected 32 character hex string for MD5")

-- Same object should produce same hash
local h2 = hash.md5_obj({name="John", age=30})
assert(h == h2, "Same object should produce same hash")

-- Different object should produce different hash
local h3 = hash.md5_obj({name="Jane", age=30})
assert(h ~= h3, "Different object should produce different hash")

-- Test SHA1 object hashing
local sha1_h = hash.sha1_obj({foo="bar", num=123})
assert(type(sha1_h) == "string", "Expected string result")
assert(#sha1_h == 40, "Expected 40 character hex string for SHA1")

-- Test SHA256 object hashing
local sha256_h = hash.sha256_obj({key="value", nested={deep="data"}})
assert(type(sha256_h) == "string", "Expected string result")
assert(#sha256_h == 64, "Expected 64 character hex string for SHA256")

-- Test SHA512 object hashing
local sha512_h = hash.sha512_obj({data="test"})
assert(type(sha512_h) == "string", "Expected string result")
assert(#sha512_h == 128, "Expected 128 character hex string for SHA512")

return true
