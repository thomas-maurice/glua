---@meta hash

---@class hash
local hash = {}

--- computes the hex-encoded MD5 hash of a string
---@param s string
---@return string
function hash.md5(s) end

--- computes the hex-encoded SHA1 hash of a string
---@param s string
---@return string
function hash.sha1(s) end

--- computes the hex-encoded SHA256 hash of a string
---@param s string
---@return string
function hash.sha256(s) end

--- computes the hex-encoded SHA512 hash of a string
---@param s string
---@return string
function hash.sha512(s) end

--- computes the hex-encoded HMAC-SHA256 of a message with a key
---@param message string
---@param key string
---@return string
function hash.hmac_sha256(message, key) end

--- computes the MD5 hash of a Lua value serialised to JSON
---@param obj any
---@return string
function hash.md5_obj(obj) end

--- computes the SHA1 hash of a Lua value serialised to JSON
---@param obj any
---@return string
function hash.sha1_obj(obj) end

--- computes the SHA256 hash of a Lua value serialised to JSON
---@param obj any
---@return string
function hash.sha256_obj(obj) end

--- computes the SHA512 hash of a Lua value serialised to JSON
---@param obj any
---@return string
function hash.sha512_obj(obj) end

return hash
