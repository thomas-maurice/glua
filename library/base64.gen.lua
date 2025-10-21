---@meta base64

---@class base64
local base64 = {}

---@param string str The string to encode
---@return string encoded The base64 encoded string
function base64.encode(str) end

---@param string encoded The base64 encoded string
---@return string decoded The decoded string, or nil on error
---@return string|nil err Error message if decoding failed
function base64.decode(encoded) end

---@param string str The string to encode
---@return string encoded The URL-safe base64 encoded string
function base64.encode_url(str) end

---@param string encoded The URL-safe base64 encoded string
---@return string decoded The decoded string, or nil on error
---@return string|nil err Error message if decoding failed
function base64.decode_url(encoded) end

return base64
