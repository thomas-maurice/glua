---@meta json

---@class json
local json = {}

---@param string jsonstr The JSON string to parse
---@return table tbl The parsed JSON as a Lua table, or nil on error
---@return string|nil err Error message if parsing failed
function json.parse(jsonstr) end

---@param table tbl The Lua table to convert to JSON
---@return string str The JSON string, or nil on error
---@return string|nil err Error message if conversion failed
function json.stringify(tbl) end

return json
