---@meta spew

---@class spew
local spew = {}

---@param value any The Lua value to dump (table, string, number, etc.)
function spew.dump(value) end

---@param value any The Lua value to dump (table, string, number, etc.)
---@return str string A JSON string representation of the value
function spew.sdump(value) end

return spew
