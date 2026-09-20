---@meta hex

---@class hex
local hex = {}

--- encodes a string to hexadecimal
---@param s string the raw string to encode
---@return string encoded the lowercase hexadecimal encoding of s
function hex.encode(s) end

--- decodes a hexadecimal string, raises on invalid input
---@param encoded string a hexadecimal string with an even number of digits
---@return string s the decoded raw string
function hex.decode(encoded) end

return hex
