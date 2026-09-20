---@meta json

---@class json
local json = {}

--- parses a JSON string into a Lua value, raises on invalid JSON
---@param jsonstr string a JSON-encoded string
---@return any value the decoded value: table, string, number, boolean or nil
function json.parse(jsonstr) end

--- converts a Lua value to a JSON string, raises on error
---@param value any the Lua value to encode; tables become JSON objects or arrays
---@return string jsonstr the JSON-encoded string
function json.stringify(value) end

return json
