-- EdDSA is its own family (Ed25519 keys are neither RSA nor EC keys): mixing
-- it with any other family must raise at option-validation time, before the
-- token is even looked at -- same shape as algorithm_validation.lua's
-- existing family-mixing assertions, extended to the new family. A
-- syntactically invalid token is used throughout so a pass here proves the
-- rejection came from option validation, not from token/key parsing.
local jwt = require("jwt")

local garbageToken = "not-a-real-token"
local secret = "some-secret"

-- EdDSA + EC: neither is RSA/HMAC, but Ed25519 and ECDSA keys are still
-- different key types -- exactly as ambiguous as any other cross-family pair.
local okEC, errEC = pcall(jwt.verify, garbageToken, secret, { algorithms = { "EdDSA", "ES256" } })
assert(not okEC, "verify must raise when opts.algorithms mixes EdDSA and EC")
assert(string.find(tostring(errEC), "mixes"), "error should mention the mixed-family reason, got: " .. tostring(errEC))

-- EdDSA + HMAC: the textbook-shaped confusion (a public key reinterpreted as
-- an HMAC secret) but for the new family.
local okHMAC, errHMAC = pcall(jwt.verify, garbageToken, secret, { algorithms = { "EdDSA", "HS256" } })
assert(not okHMAC, "verify must raise when opts.algorithms mixes EdDSA and HMAC")
assert(string.find(tostring(errHMAC), "mixes"), "error should mention the mixed-family reason, got: " .. tostring(errHMAC))

-- EdDSA alone is a valid, single-family algorithms list -- sign() must get
-- past option validation and fail only once it tries to parse `key` as PEM
-- (this module has no way to generate a real Ed25519 PEM key from Lua; the
-- real sign/verify round trip is covered at the Go level by
-- TestSignVerify_EdDSA_RoundTrip). The failure must be about the key, not
-- about family mixing.
local okSign, errSign = pcall(jwt.sign, { sub = "x" }, secret, { algorithm = "EdDSA" })
assert(not okSign, "sign must raise when key is not a PEM-encoded Ed25519 private key")
assert(not string.find(tostring(errSign), "mixes"), "a single-algorithm EdDSA request must not fail as a family-mixing error, got: " .. tostring(errSign))

return true
