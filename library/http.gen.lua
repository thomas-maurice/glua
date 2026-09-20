---@meta http

---@class http
local http = {}

--- performs an HTTP GET request, raises on network error or timeout
---@param url string the absolute URL to request
---@param headers gopher-lua.LTable optional table of header name to value; pass nil for none
---@return any response table with status (number), body (string) and headers (table)
function http.get(url, headers) end

--- performs an HTTP POST request, raises on network error or timeout
---@param url string the absolute URL to request
---@param body string the request body; pass an empty string for no body
---@param headers gopher-lua.LTable optional table of header name to value; pass nil for none
---@return any response table with status (number), body (string) and headers (table)
function http.post(url, body, headers) end

--- performs an HTTP PUT request, raises on network error or timeout
---@param url string the absolute URL to request
---@param body string the request body; pass an empty string for no body
---@param headers gopher-lua.LTable optional table of header name to value; pass nil for none
---@return any response table with status (number), body (string) and headers (table)
function http.put(url, body, headers) end

--- performs an HTTP DELETE request, raises on network error or timeout
---@param url string the absolute URL to request
---@param headers gopher-lua.LTable optional table of header name to value; pass nil for none
---@return any response table with status (number), body (string) and headers (table)
function http.delete(url, headers) end

--- performs an HTTP request with a custom method, raises on network error or timeout
---@param method string the HTTP method, e.g. "PATCH" or "HEAD"
---@param url string the absolute URL to request
---@param body string the request body; pass an empty string for no body
---@param headers gopher-lua.LTable optional table of header name to value; pass nil for none
---@return any response table with status (number), body (string) and headers (table)
function http.request(method, url, body, headers) end

return http
