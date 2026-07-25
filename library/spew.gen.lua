---@meta spew

---@class spew
local spew = {}

--- prints a Lua value to stdout as colored indented JSON
---@param value any
function spew.dump(value) end

--- returns a JSON string representation of a Lua value
---@param value any
---@return string
function spew.sdump(value) end

return spew
