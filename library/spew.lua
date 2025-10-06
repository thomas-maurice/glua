---@meta

---@class spew
local spew = {}

--- dump: prints the contents of a Lua value to stdout with detailed formatting. This is useful for debugging and inspecting complex table structures.
---@param value any The Lua value to dump (table, string, number, etc.)
function spew.dump(value) end

--- sdump: returns a string representation of a Lua value with detailed formatting. Unlike dump, this returns the string instead of printing to stdout.
---@param value any The Lua value to dump (table, string, number, etc.)
---@return string A detailed string representation of the value
function spew.sdump(value) end

return spew
