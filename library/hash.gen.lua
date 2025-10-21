---@meta hash

---@class hash
local hash = {}

---@param string str The string to hash
---@return string hash The hex-encoded MD5 hash
function hash.md5(str) end

---@param string str The string to hash
---@return string hash The hex-encoded SHA1 hash
function hash.sha1(str) end

---@param string str The string to hash
---@return string hash The hex-encoded SHA256 hash
function hash.sha256(str) end

---@param string str The string to hash
---@return string hash The hex-encoded SHA512 hash
function hash.sha512(str) end

---@param string message The message to authenticate
---@param string key The secret key
---@return string hash The hex-encoded HMAC-SHA256
function hash.hmac_sha256(message, key) end

---@param table obj The table to hash
---@return string hash The hex-encoded MD5 hash
---@return string|nil err Error message if conversion fails
function hash.md5_obj(obj) end

---@param table obj The table to hash
---@return string hash The hex-encoded SHA1 hash
---@return string|nil err Error message if conversion fails
function hash.sha1_obj(obj) end

---@param table obj The table to hash
---@return string hash The hex-encoded SHA256 hash
---@return string|nil err Error message if conversion fails
function hash.sha256_obj(obj) end

---@param table obj The table to hash
---@return string hash The hex-encoded SHA512 hash
---@return string|nil err Error message if conversion fails
function hash.sha512_obj(obj) end

return hash
