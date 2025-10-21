-- Test basic string hashing functions
local hash = require("hash")

-- Test MD5
local md5_result = hash.md5("hello world")
assert(md5_result == "5eb63bbbe01eeed093cb22bb8f5acdc3", "Expected correct MD5 hash, got: " .. md5_result)

-- Test SHA1
local sha1_result = hash.sha1("hello world")
assert(sha1_result == "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", "Expected correct SHA1 hash, got: " .. sha1_result)

-- Test SHA256
local sha256_result = hash.sha256("hello world")
assert(sha256_result == "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", "Expected correct SHA256 hash, got: " .. sha256_result)

-- Test SHA512
local sha512_result = hash.sha512("hello world")
assert(type(sha512_result) == "string", "Expected string result")
assert(#sha512_result == 128, "Expected 128 character hex string for SHA512, got length: " .. #sha512_result)

-- Test empty string
local empty_sha256 = hash.sha256("")
assert(empty_sha256 == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "Expected correct SHA256 of empty string")

return true
