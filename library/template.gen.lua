---@meta template

---@class template
local template = {}

---@param tmpl string The template string
---@param data table The data to render with
---@return result string The rendered template, or nil on error
---@return err string|nil Error message if rendering failed
function template.render(tmpl, data) end

---@param path string The path to the template file
---@param data table The data to render with
---@return result string The rendered template, or nil on error
---@return err string|nil Error message if rendering failed
function template.render_file(path, data) end

return template
