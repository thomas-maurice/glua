---@meta hash

---@class hash
local hash = {}

---@param str string The string to hash
---@return hash string The hex-encoded MD5 hash
function hash.md5(str) end

---@param str string The string to hash
---@return hash string The hex-encoded SHA1 hash
function hash.sha1(str) end

---@param str string The string to hash
---@return hash string The hex-encoded SHA256 hash
function hash.sha256(str) end

---@param str string The string to hash
---@return hash string The hex-encoded SHA512 hash
function hash.sha512(str) end

---@param message string The message to authenticate
---@param key string The secret key
---@return hash string The hex-encoded HMAC-SHA256
function hash.hmac_sha256(message, key) end

return hash
