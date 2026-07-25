---@meta yaml

---@class yaml
local yaml = {}

--- parses a YAML string into a Lua value, raises on invalid YAML
---@param yamlstr string
---@return any
function yaml.parse(yamlstr) end

--- converts a Lua value to a YAML string, raises on error
---@param value any
---@return string
function yaml.stringify(value) end

return yaml
