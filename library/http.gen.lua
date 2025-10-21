---@meta http

---@class http
local http = {}

---@param string url The URL to request
---@param table|nil headers Optional headers table
---@return table response Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.get(url, headers) end

---@param string url The URL to request
---@param string body The request body
---@param table|nil headers Optional headers table
---@return table response Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.post(url, body, headers) end

---@param string url The URL to request
---@param string body The request body
---@param table|nil headers Optional headers table
---@return table response Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.put(url, body, headers) end

---@param string url The URL to request
---@param table|nil headers Optional headers table
---@return table response Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.delete(url, headers) end

---@param string method The HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
---@param string url The URL to request
---@param string|nil body Optional request body
---@param table|nil headers Optional headers table
---@return table response Response table with status, body, headers, or nil on error
---@return string|nil Error message if request failed
function http.request(method, url, body, headers) end

return http
