---@meta template

---@class template
local template = {}

--- renders a template string with data, raises on error
---@param tmpl string a Go text/template source string
---@param data table table exposed to the template as the root context (dot)
---@return string out the rendered template output
function template.render(tmpl, data) end

--- renders a template file with data, raises on error
---@param path string path to a file containing Go text/template source
---@param data table table exposed to the template as the root context (dot)
---@return string out the rendered template output
function template.render_file(path, data) end

return template
