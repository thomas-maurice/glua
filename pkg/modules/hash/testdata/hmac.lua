-- Test HMAC-SHA256 function
local hash = require("hash")

-- Test basic HMAC
local h = hash.hmac_sha256("message", "secret_key")
assert(type(h) == "string", "Expected string result")
assert(#h == 64, "Expected 64 character hex string for HMAC-SHA256, got length: " .. #h)

-- Same key should produce same hash
local h1 = hash.hmac_sha256("message", "key")
local h2 = hash.hmac_sha256("message", "key")
assert(h1 == h2, "Expected same HMAC for same inputs")

-- Different key should produce different hash
local h3 = hash.hmac_sha256("message", "key1")
local h4 = hash.hmac_sha256("message", "key2")
assert(h3 ~= h4, "Expected different HMAC for different keys")

return true
