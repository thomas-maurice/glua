---@meta yaml

---@class yaml
local yaml = {}

---@param yamlstr string The YAML string to parse
---@return table tbl The parsed YAML as a Lua table, or nil on error
---@return string|nil err Error message if parsing failed
function yaml.parse(yamlstr) end

---@param tbl table The Lua table to convert to YAML
---@return string str The YAML string, or nil on error
---@return string|nil err Error message if conversion failed
function yaml.stringify(tbl) end

return yaml
