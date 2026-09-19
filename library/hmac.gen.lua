---@meta hmac

---@class hmac
local hmac = {}

--- computes the hex-encoded HMAC-SHA1 tag of a message with a key
---@param message string the message to authenticate
---@param key string the shared secret key
---@return string tag the lowercase hex-encoded HMAC-SHA1 tag
function hmac.sha1(message, key) end

--- computes the hex-encoded HMAC-SHA256 tag of a message with a key
---@param message string the message to authenticate
---@param key string the shared secret key
---@return string tag the lowercase hex-encoded HMAC-SHA256 tag
function hmac.sha256(message, key) end

--- computes the hex-encoded HMAC-SHA512 tag of a message with a key
---@param message string the message to authenticate
---@param key string the shared secret key
---@return string tag the lowercase hex-encoded HMAC-SHA512 tag
function hmac.sha512(message, key) end

--- verifies an HMAC-SHA1 tag in constant time; never raises, a malformed tag simply returns false
---@param message string the message that was authenticated
---@param key string the shared secret key
---@param tag string the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false
---@return boolean ok true if tag is the correct HMAC-SHA1 tag for message under key
function hmac.verify_sha1(message, key, tag) end

--- verifies an HMAC-SHA256 tag in constant time; never raises, a malformed tag simply returns false
---@param message string the message that was authenticated
---@param key string the shared secret key
---@param tag string the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false
---@return boolean ok true if tag is the correct HMAC-SHA256 tag for message under key
function hmac.verify_sha256(message, key, tag) end

--- verifies an HMAC-SHA512 tag in constant time; never raises, a malformed tag simply returns false
---@param message string the message that was authenticated
---@param key string the shared secret key
---@param tag string the hex-encoded tag to verify (case-insensitive); a non-hex or wrong-length tag returns false
---@return boolean ok true if tag is the correct HMAC-SHA512 tag for message under key
function hmac.verify_sha512(message, key, tag) end

return hmac
