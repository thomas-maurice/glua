---@meta template

---@class template
local template = {}

---@param string tmpl The template string
---@param table data The data to render with
---@return string result The rendered template, or nil on error
---@return string|nil err Error message if rendering failed
function template.render(tmpl, data) end

---@param string path The path to the template file
---@param table data The data to render with
---@return string result The rendered template, or nil on error
---@return string|nil err Error message if rendering failed
function template.render_file(path, data) end

return template
