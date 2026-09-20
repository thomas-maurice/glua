---@meta yaml

---@class yaml
local yaml = {}

--- parses a YAML string into a Lua value, raises on invalid YAML
---@param yamlstr string a YAML-encoded string
---@return any value the decoded value: table, string, number, boolean or nil
function yaml.parse(yamlstr) end

--- converts a Lua value to a YAML string, raises on error
---@param value any the Lua value to encode; tables become YAML mappings or sequences
---@return string yamlstr the YAML-encoded string
function yaml.stringify(value) end

return yaml
