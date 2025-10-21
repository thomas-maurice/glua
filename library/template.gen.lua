---@meta template

---@class template
local template = {}

---@param tmpl string The template string
---@param data table The data to render with
---@return string result The rendered template, or nil on error
---@return string|nil err Error message if rendering failed
function template.render(tmpl, data) end

---@param path string The path to the template file
---@param data table The data to render with
---@return string result The rendered template, or nil on error
---@return string|nil err Error message if rendering failed
function template.render_file(path, data) end

return template
