---@meta base64

---@class base64
local base64 = {}

--- encodes a string to standard base64
---@param s string the raw string to encode
---@return string encoded the standard-alphabet base64 encoding of s
function base64.encode(s) end

--- decodes a standard base64 string, raises on invalid input
---@param encoded string a standard-alphabet base64 string to decode
---@return string s the decoded raw string
function base64.decode(encoded) end

--- encodes a string to URL-safe base64
---@param s string the raw string to encode
---@return string encoded the URL-safe base64 encoding of s
function base64.encode_url(s) end

--- decodes a URL-safe base64 string, raises on invalid input
---@param encoded string a URL-safe base64 string to decode
---@return string s the decoded raw string
function base64.decode_url(encoded) end

return base64
