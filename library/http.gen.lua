---@meta http

---@class http
local http = {}

---@param url string The URL to request
---@param headers table|nil Optional headers table
---@return table Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.get(url, headers) end

---@param url string The URL to request
---@param body string The request body
---@param headers table|nil Optional headers table
---@return table Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.post(url, body, headers) end

---@param url string The URL to request
---@param body string The request body
---@param headers table|nil Optional headers table
---@return table Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.put(url, body, headers) end

---@param url string The URL to request
---@param headers table|nil Optional headers table
---@return table Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.delete(url, headers) end

---@param method string The HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
---@param url string The URL to request
---@param body string|nil Optional request body
---@param headers table|nil Optional headers table
---@return table Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.request(method, url, body, headers) end

return http
