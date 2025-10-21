---@meta spew

---@class spew
local spew = {}

---@param any value The Lua value to dump (table, string, number, etc.)
function spew.dump(value) end

---@param any value The Lua value to dump (table, string, number, etc.)
---@return string str A JSON string representation of the value
function spew.sdump(value) end

return spew
