---@meta hex

---@class hex
local hex = {}

---@param str string The string to encode
---@return string The hex encoded string
function hex.encode(str) end

---@param encoded string The hex encoded string
---@return string The decoded string, or nil on error
---@return string|nil Error message if decoding failed
function hex.decode(encoded) end

return hex
