---@meta

---@class json
local json = {}

--- parse: parses a JSON string and returns a Lua table. Returns nil and error message on failure.
---@param jsonstr string The JSON string to parse
---@return table The parsed JSON as a Lua table, or nil on error
---@return string|nil Error message if parsing failed
function json.parse(jsonstr) end

--- stringify: converts a Lua table to a JSON string. Returns nil and error message on failure.
---@param tbl table The Lua table to convert to JSON
---@return string The JSON string, or nil on error
---@return string|nil Error message if conversion failed
function json.stringify(tbl) end
