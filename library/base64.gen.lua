---@meta base64

---@class base64
local base64 = {}

---@param str string The string to encode
---@return encoded string The base64 encoded string
function base64.encode(str) end

---@param encoded string The base64 encoded string
---@return decoded string The decoded string, or nil on error
---@return err string|nil Error message if decoding failed
function base64.decode(encoded) end

---@param str string The string to encode
---@return encoded string The URL-safe base64 encoded string
function base64.encode_url(str) end

---@param encoded string The URL-safe base64 encoded string
---@return decoded string The decoded string, or nil on error
---@return err string|nil Error message if decoding failed
function base64.decode_url(encoded) end

return base64
