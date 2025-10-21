---@meta hex

---@class hex
local hex = {}

---@param string str The string to encode
---@return string encoded The hex encoded string
function hex.encode(str) end

---@param string encoded The hex encoded string
---@return string decoded The decoded string, or nil on error
---@return string|nil err Error message if decoding failed
function hex.decode(encoded) end

return hex
