-- Test the intended usage: compute a tag, then verify it via verify_sha256.
-- There is deliberately no == comparison on tags anywhere in this script.
local hmac = require("hmac")

local secret = "shared-secret"
local payload = "the quick brown fox"

local tag = hmac.sha256(payload, secret)
assert(type(tag) == "string", "expected string tag")
assert(#tag == 64, "expected 64 hex chars for sha256 tag")

-- Correct tag verifies.
assert(hmac.verify_sha256(payload, secret, tag), "expected verify to succeed for the correct tag")

-- Wrong secret does not verify.
assert(not hmac.verify_sha256(payload, "wrong-secret", tag), "expected verify to fail for the wrong secret")

-- sha1 and sha512 round trip the same way.
local tag1 = hmac.sha1(payload, secret)
assert(hmac.verify_sha1(payload, secret, tag1))
assert(not hmac.verify_sha1(payload, "wrong-secret", tag1))

local tag512 = hmac.sha512(payload, secret)
assert(hmac.verify_sha512(payload, secret, tag512))
assert(not hmac.verify_sha512(payload, "wrong-secret", tag512))

return true
