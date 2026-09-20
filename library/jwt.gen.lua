---@meta jwt

---@class jwt.Decoded
---@field claims table<string, any> the token's claims segment, decoded
---@field header table<string, any> the token's header segment, decoded

---@class jwt.SignOptions
---@field algorithm string REQUIRED; never "none"
---@field kid string "" = omit the header field
---@field typ string "" = "JWT" (jwt/v5's own default)

---@class jwt.VerifyOptions
---@field algorithms string[] REQUIRED, non-empty, single family, never "none"
---@field allow_missing_exp boolean false (default): a token with no exp is REJECTED
---@field audience string "" = do not check
---@field issuer string "" = do not check
---@field leeway_seconds number 0 = no leeway; applied to both exp and nbf
---@field subject string "" = do not check

---@class jwt
local jwt = {}

--- decodes a JWT's header and claims WITHOUT verifying its signature -- DOES NOT verify anything; use only to read kid/iss before choosing a key
---@param token string the compact JWT string (header.claims.signature)
---@return jwt.Decoded decoded the unverified header and claims
function jwt.decode_unverified(token) end

--- verifies a JWT's signature and claims; raises on ANY failure: bad signature, alg not in opts.algorithms, alg:none, expired, not yet valid, issuer/audience/subject mismatch, or an unparseable key
---@param token string the compact JWT string
---@param key string HS*: the raw shared secret; RS*/PS*/ES*/EdDSA: a PEM public key or certificate (never a raw Ed25519 seed/expanded key)
---@param opts jwt.VerifyOptions required options: algorithms (non-empty, single family, never "none"), issuer, audience, subject, leeway_seconds (0..300, seconds), allow_missing_exp
---@return table<string, any> claims the token's verified claims
function jwt.verify(token, key, opts) end

--- signs claims into a compact JWT; raises on an unknown/forbidden algorithm, a key that does not match the algorithm, or claims that cannot be JSON-encoded
---@param claims table<string, any> the claims to encode
---@param key string HS*: the raw shared secret; RS*/PS*/ES*/EdDSA: a PEM private key (never a raw Ed25519 seed/expanded key)
---@param opts jwt.SignOptions required options: algorithm (never "none"), kid, typ
---@return string token the compact JWT string
function jwt.sign(claims, key, opts) end

return jwt
