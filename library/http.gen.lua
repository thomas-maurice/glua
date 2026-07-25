---@meta http

---@class http
local http = {}

--- performs an HTTP GET request, raises on network error
---@param url string
---@param headers gopher-lua.LTable
---@return any
function http.get(url, headers) end

--- performs an HTTP POST request, raises on network error
---@param url string
---@param body string
---@param headers gopher-lua.LTable
---@return any
function http.post(url, body, headers) end

--- performs an HTTP PUT request, raises on network error
---@param url string
---@param body string
---@param headers gopher-lua.LTable
---@return any
function http.put(url, body, headers) end

--- performs an HTTP DELETE request, raises on network error
---@param url string
---@param headers gopher-lua.LTable
---@return any
function http.delete(url, headers) end

--- performs an HTTP request with a custom method, raises on network error
---@param method string
---@param url string
---@param body string
---@param headers gopher-lua.LTable
---@return any
function http.request(method, url, body, headers) end

return http
