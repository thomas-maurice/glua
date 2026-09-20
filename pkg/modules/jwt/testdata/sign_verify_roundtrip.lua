-- HS256 sign -> verify round trip, and the claims that come back must match
-- what was signed.
local jwt = require("jwt")

local secret = "correct-horse-battery-staple"
-- exp is a fixed far-future Unix timestamp (year 2100) rather than a
-- computed "now + n" value, since this script has no access to a clock
-- module -- verify rejects a token with no exp by default, so one must be
-- present.
local claims = { sub = "svc-a", role = "admin", exp = 4102444800 }

local token = jwt.sign(claims, secret, { algorithm = "HS256" })
assert(type(token) == "string", "sign must return a string")

local verified = jwt.verify(token, secret, { algorithms = { "HS256" } })
assert(verified.sub == "svc-a", "sub claim must round-trip")
assert(verified.role == "admin", "role claim must round-trip")

-- Wrong secret must be rejected.
local ok = pcall(jwt.verify, token, "wrong-secret", { algorithms = { "HS256" } })
assert(not ok, "verify must reject a token signed with a different secret")

return true
