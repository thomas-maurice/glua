---@meta json

---@class json
local json = {}

--- parses a JSON string into a Lua value, raises on invalid JSON
---@param jsonstr string
---@return any
function json.parse(jsonstr) end

--- converts a Lua value to a JSON string, raises on error
---@param value any
---@return string
function json.stringify(value) end

return json
