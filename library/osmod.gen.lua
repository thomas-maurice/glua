---@meta osmod

---@class osmod
local osmod = {}

---@param string name The environment variable name
---@return string The value of the environment variable (empty string if not set)
function osmod.getenv(name) end

---@param string name The environment variable name
---@param string value The value to set
---@return string|nil Error message if operation failed
function osmod.setenv(name, value) end

---@param string name The environment variable name
---@return string|nil Error message if operation failed
function osmod.unsetenv(name) end

---@return string The hostname
---@return string|nil Error message if operation failed
function osmod.hostname() end

---@return string The temporary directory path
function osmod.tmpdir() end

return osmod
