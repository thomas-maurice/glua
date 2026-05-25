---@meta osmod

---@class osmod
local osmod = {}

--- returns the value of an environment variable
---@param name string
---@return string
function osmod.getenv(name) end

--- sets an environment variable, raises on error
---@param name string
---@param value string
---@return boolean
function osmod.setenv(name, value) end

--- unsets an environment variable, raises on error
---@param name string
---@return boolean
function osmod.unsetenv(name) end

--- returns the system hostname, raises on error
---@return string
function osmod.hostname() end

--- returns the default temporary directory path
---@return string
function osmod.tmpdir() end

return osmod
