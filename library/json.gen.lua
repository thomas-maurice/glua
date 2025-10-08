---@meta json

---@class json
local json = {}

---@param jsonstr string The JSON string to parse
---@return table The parsed JSON as a Lua table, or nil on error
---@return string|nil Error message if parsing failed
function json.parse(jsonstr) end

---@param tbl table The Lua table to convert to JSON
---@return string The JSON string, or nil on error
---@return string|nil Error message if conversion failed
function json.stringify(tbl) end

return json
