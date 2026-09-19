-- decode_unverified must read header/claims WITHOUT checking the signature
-- at all -- that is its entire, deliberately dangerous contract. A tampered
-- signature must not stop it from decoding.
local jwt = require("jwt")

local secret = "shared-secret"
local token = jwt.sign({ sub = "svc-a", exp = 4102444800 }, secret, { algorithm = "HS256" })

local decoded = jwt.decode_unverified(token)
assert(decoded.header.alg == "HS256", "header.alg must be readable")
assert(decoded.claims.sub == "svc-a", "claims must be readable")

-- Tamper with the signature segment (flip the last character).
local headerPart, claimsPart, sigPart = token:match("^([^.]+)%.([^.]+)%.(.+)$")
local lastChar = sigPart:sub(-1, -1)
local replacement = (lastChar == "A") and "B" or "A"
local tamperedSig = sigPart:sub(1, -2) .. replacement
local tampered = headerPart .. "." .. claimsPart .. "." .. tamperedSig

local ok, decodedTampered = pcall(jwt.decode_unverified, tampered)
assert(ok, "decode_unverified must succeed even on a token with an invalid signature")
assert(decodedTampered.claims.sub == "svc-a", "claims must still be readable from a tampered token")

-- But verify MUST reject the same tampered token.
local verifyOk = pcall(jwt.verify, tampered, secret, { algorithms = { "HS256" } })
assert(not verifyOk, "verify must reject a token with a tampered signature")

-- Malformed input (wrong number of segments) must raise.
local malformedOk = pcall(jwt.decode_unverified, "not-a-real-token")
assert(not malformedOk, "decode_unverified must raise on a token that isn't three segments")

return true
