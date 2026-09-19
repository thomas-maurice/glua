-- Claim validation: exp/nbf are always checked, leeway smooths clock skew,
-- issuer/audience/subject are checked only when opts asks for them, and a
-- missing exp is rejected unless the caller opts in.
local jwt = require("jwt")

local secret = "shared-secret"
local FAR_FUTURE = 4102444800 -- year 2100, used whenever exp itself isn't the thing under test

-- Expired token rejected.
local expiredToken = jwt.sign({ sub = "svc-a", exp = 1000 }, secret, { algorithm = "HS256" })
local okExpired, errExpired = pcall(jwt.verify, expiredToken, secret, { algorithms = { "HS256" } })
assert(not okExpired, "verify must reject an expired token")
assert(string.find(tostring(errExpired), "expired"), "error should say the token expired, got: " .. tostring(errExpired))

-- leeway_seconds making a just-expired token pass needs a controllable
-- clock (this script only has real time), so that case is covered at the
-- Go level: TestVerify_LeewayAllowsJustExpiredToken.

-- nbf in the future rejected.
local nbfToken = jwt.sign({ sub = "svc-a", exp = FAR_FUTURE, nbf = FAR_FUTURE - 1000 }, secret, { algorithm = "HS256" })
-- nbf is also far in the future relative to "now" (real clock), so this must fail.
local okNbf, errNbf = pcall(jwt.verify, nbfToken, secret, { algorithms = { "HS256" } })
assert(not okNbf, "verify must reject a token not valid yet")
assert(string.find(tostring(errNbf), "not valid yet"), "error should say not valid yet, got: " .. tostring(errNbf))

-- Missing exp rejected by default, accepted with allow_missing_exp = true.
local noExpToken = jwt.sign({ sub = "svc-a" }, secret, { algorithm = "HS256" })
local okNoExp, errNoExp = pcall(jwt.verify, noExpToken, secret, { algorithms = { "HS256" } })
assert(not okNoExp, "verify must reject a token with no exp by default")
assert(string.find(tostring(errNoExp), "required claim"), "error should say a required claim is missing, got: " .. tostring(errNoExp))

local okAllowed, claimsAllowed = pcall(jwt.verify, noExpToken, secret, { algorithms = { "HS256" }, allow_missing_exp = true })
assert(okAllowed, "verify must accept a token with no exp when allow_missing_exp is true")
assert(claimsAllowed.sub == "svc-a")

-- Wrong issuer / audience rejected.
local issAudToken = jwt.sign({ sub = "svc-a", exp = FAR_FUTURE, iss = "https://idp.example.com", aud = "api.example.com" }, secret, { algorithm = "HS256" })

local okIss, errIss = pcall(jwt.verify, issAudToken, secret, { algorithms = { "HS256" }, issuer = "https://someone-else.example.com" })
assert(not okIss, "verify must reject a wrong issuer")
assert(string.find(tostring(errIss), "issuer"), "error should mention issuer, got: " .. tostring(errIss))

local okAud, errAud = pcall(jwt.verify, issAudToken, secret, { algorithms = { "HS256" }, audience = "someone-else" })
assert(not okAud, "verify must reject a wrong audience")
assert(string.find(tostring(errAud), "audience"), "error should mention audience, got: " .. tostring(errAud))

-- Correct issuer and audience together must succeed.
local okBoth, claimsBoth = pcall(jwt.verify, issAudToken, secret, {
  algorithms = { "HS256" },
  issuer = "https://idp.example.com",
  audience = "api.example.com",
})
assert(okBoth, "verify must accept matching issuer and audience: " .. tostring(claimsBoth))

return true
