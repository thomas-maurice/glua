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

---@param obj table The table to hash
---@return hash string The hex-encoded MD5 hash
---@return error string|nil Error message if conversion fails
function hash.md5_obj(obj) end

---@param obj table The table to hash
---@return hash string The hex-encoded SHA1 hash
---@return error string|nil Error message if conversion fails
function hash.sha1_obj(obj) end

---@param obj table The table to hash
---@return hash string The hex-encoded SHA256 hash
---@return error string|nil Error message if conversion fails
function hash.sha256_obj(obj) end

---@param obj table The table to hash
---@return hash string The hex-encoded SHA512 hash
---@return error string|nil Error message if conversion fails
function hash.sha512_obj(obj) end

return hash
