---@meta template

---@class template
local template = {}

--- renders a template string with data, raises on error
---@param tmpl string
---@param data gopher-lua.LTable
---@return string
function template.render(tmpl, data) end

--- renders a template file with data, raises on error
---@param path string
---@param data gopher-lua.LTable
---@return string
function template.render_file(path, data) end

return template
