-- The entire design point of this module: there is no comparison primitive
-- exposed to Lua. If this test ever fails, someone added hmac.equal or an
-- equivalent, which reopens the timing-oracle footgun the module exists to
-- close.
local hmac = require("hmac")

assert(hmac.equal == nil, "hmac.equal must not exist")
assert(hmac.constant_time_compare == nil, "hmac.constant_time_compare must not exist")

-- Also confirm a garbage tag is a clean pcall'd false, not a raised error —
-- this is the shape a caller relies on when the tag comes straight from
-- attacker-controlled input (e.g. a request header) without any
-- pre-validation.
local ok, result = pcall(hmac.verify_sha256, "message", "key", "not hex garbage")
assert(ok, "verify_sha256 must not raise on a garbage tag")
assert(result == false, "verify_sha256 must return false for a garbage tag")

return true
