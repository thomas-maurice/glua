---@meta hash

---@class hash
local hash = {}

--- computes the hex-encoded MD5 hash of a string
---@param s string the string to hash
---@return string hex the lowercase hex-encoded MD5 digest
function hash.md5(s) end

--- computes the hex-encoded SHA1 hash of a string
---@param s string the string to hash
---@return string hex the lowercase hex-encoded SHA1 digest
function hash.sha1(s) end

--- computes the hex-encoded SHA256 hash of a string
---@param s string the string to hash
---@return string hex the lowercase hex-encoded SHA256 digest
function hash.sha256(s) end

--- computes the hex-encoded SHA512 hash of a string
---@param s string the string to hash
---@return string hex the lowercase hex-encoded SHA512 digest
function hash.sha512(s) end

--- computes the hex-encoded HMAC-SHA256 of a message with a key
---@param message string the message to authenticate
---@param key string the shared secret key
---@return string hex the lowercase hex-encoded HMAC-SHA256 tag
function hash.hmac_sha256(message, key) end

--- computes the MD5 hash of a Lua value serialised to JSON
---@param obj any the Lua value to hash; it is JSON-marshalled before hashing
---@return string hex the lowercase hex-encoded MD5 digest of the JSON encoding
function hash.md5_obj(obj) end

--- computes the SHA1 hash of a Lua value serialised to JSON
---@param obj any the Lua value to hash; it is JSON-marshalled before hashing
---@return string hex the lowercase hex-encoded SHA1 digest of the JSON encoding
function hash.sha1_obj(obj) end

--- computes the SHA256 hash of a Lua value serialised to JSON
---@param obj any the Lua value to hash; it is JSON-marshalled before hashing
---@return string hex the lowercase hex-encoded SHA256 digest of the JSON encoding
function hash.sha256_obj(obj) end

--- computes the SHA512 hash of a Lua value serialised to JSON
---@param obj any the Lua value to hash; it is JSON-marshalled before hashing
---@return string hex the lowercase hex-encoded SHA512 digest of the JSON encoding
function hash.sha512_obj(obj) end

return hash
