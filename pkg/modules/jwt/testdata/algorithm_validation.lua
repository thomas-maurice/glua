-- opts.algorithms is validated BEFORE the token is even looked at -- these
-- assertions all pass a syntactically invalid token, so a pass here proves
-- the rejection came from option validation, not from token parsing.
local jwt = require("jwt")

local garbageToken = "not-a-real-token"
local secret = "some-secret"

-- Missing algorithms entirely.
local ok1, err1 = pcall(jwt.verify, garbageToken, secret, {})
assert(not ok1, "verify must raise when opts.algorithms is absent")
assert(string.find(tostring(err1), "algorithms"), "error should mention algorithms, got: " .. tostring(err1))

-- Empty algorithms list.
local ok2 = pcall(jwt.verify, garbageToken, secret, { algorithms = {} })
assert(not ok2, "verify must raise when opts.algorithms is empty")

-- "none" in the list.
local ok3, err3 = pcall(jwt.verify, garbageToken, secret, { algorithms = { "none" } })
assert(not ok3, "verify must raise when opts.algorithms contains none")
assert(string.find(tostring(err3), "none"), "error should mention none, got: " .. tostring(err3))

-- Mixed families: HMAC + RSA-ish in the same call.
local ok4, err4 = pcall(jwt.verify, garbageToken, secret, { algorithms = { "HS256", "RS256" } })
assert(not ok4, "verify must raise when opts.algorithms mixes families")
assert(string.find(tostring(err4), "mixes"), "error should mention the mixed-family reason, got: " .. tostring(err4))

-- Unsupported algorithm name.
local ok5 = pcall(jwt.verify, garbageToken, secret, { algorithms = { "HS1" } })
assert(not ok5, "verify must raise on an unsupported algorithm name")

-- sign side: "none" is refused at option-validation time too.
local ok6, err6 = pcall(jwt.sign, { sub = "x" }, secret, { algorithm = "none" })
assert(not ok6, "sign must raise when algorithm is none")
assert(string.find(tostring(err6), "none"), "error should mention none, got: " .. tostring(err6))

local ok7 = pcall(jwt.sign, { sub = "x" }, secret, { algorithm = "not-a-real-alg" })
assert(not ok7, "sign must raise on an unsupported algorithm name")

return true
