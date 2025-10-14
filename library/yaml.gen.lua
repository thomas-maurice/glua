---@meta yaml

---@class yaml
local yaml = {}

---@param yamlstr string The YAML string to parse
---@return tbl table The parsed YAML as a Lua table, or nil on error
---@return err string|nil Error message if parsing failed
function yaml.parse(yamlstr) end

---@param tbl table The Lua table to convert to YAML
---@return str string The YAML string, or nil on error
---@return err string|nil Error message if conversion failed
function yaml.stringify(tbl) end

return yaml
