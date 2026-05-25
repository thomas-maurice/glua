---@meta base64

---@class base64
local base64 = {}

--- encodes a string to standard base64
---@param s string
---@return string
function base64.encode(s) end

--- decodes a standard base64 string, raises on invalid input
---@param encoded string
---@return string
function base64.decode(encoded) end

--- encodes a string to URL-safe base64
---@param s string
---@return string
function base64.encode_url(s) end

--- decodes a URL-safe base64 string, raises on invalid input
---@param encoded string
---@return string
function base64.decode_url(encoded) end

return base64
