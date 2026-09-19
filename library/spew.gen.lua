---@meta spew

---@class spew
local spew = {}

--- prints a Lua value to stdout as colored indented JSON
---@param value any the Lua value to dump; tables are walked recursively
function spew.dump(value) end

--- returns a JSON string representation of a Lua value
---@param value any the Lua value to dump; tables are walked recursively
---@return string s the indented JSON representation of value
function spew.sdump(value) end

return spew
