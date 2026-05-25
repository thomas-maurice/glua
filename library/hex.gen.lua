---@meta hex

---@class hex
local hex = {}

--- encodes a string to hexadecimal
---@param s string
---@return string
function hex.encode(s) end

--- decodes a hexadecimal string, raises on invalid input
---@param encoded string
---@return string
function hex.decode(encoded) end

return hex
