---@meta jsonpath

---@class jsonpath
local jsonpath = {}

--- evaluates path against data and returns every match as a new array table, in match order; empty when there is no match
---@param data any the value to query, normally a table
---@param path string a kubectl-jsonpath path; bare (".spec.replicas"), $-rooted ("$.spec.replicas") and braced ("{.spec.replicas}") forms are all accepted
---@return any[] matches an array of every matched value, in match order; always a table, never nil, even for zero or one match
function jsonpath.query(data, path) end

--- evaluates path against data and returns its first match, or default if there is no match
---@param data any the value to query, normally a table
---@param path string a kubectl-jsonpath path; bare, $-rooted and braced forms are all accepted
---@param default any returned as-is when path has no match; pass nil explicitly for no default
---@return any value the first matched value, or default
function jsonpath.first(data, path, default) end

--- reports whether path has at least one match in data
---@param data any the value to query, normally a table
---@param path string a kubectl-jsonpath path; bare, $-rooted and braced forms are all accepted
---@return boolean found true if path matched at least one value in data
function jsonpath.exists(data, path) end

--- renders template against data using the native kubectl jsonpath text-template mode, concatenating matched values and literal text
---@param data any the value to render from, normally a table
---@param template string a kubectl-jsonpath template, e.g. "{.metadata.name}{'\t'}{.status.phase}"; bare and braced forms are both accepted
---@return string text the rendered text
function jsonpath.render(data, template) end

return jsonpath
