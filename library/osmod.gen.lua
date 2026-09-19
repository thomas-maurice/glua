---@meta osmod

---@class osmod
local osmod = {}

--- returns the value of an environment variable
---@param name string the environment variable name
---@return string value the variable's value, or empty string if unset
function osmod.getenv(name) end

--- sets an environment variable, raises on error
---@param name string the environment variable name
---@param value string the value to set
---@return boolean ok true on success
function osmod.setenv(name, value) end

--- unsets an environment variable, raises on error
---@param name string the environment variable name to remove
---@return boolean ok true on success
function osmod.unsetenv(name) end

--- returns the system hostname, raises on error
---@return string name the system's hostname as reported by the OS
function osmod.hostname() end

--- returns the default temporary directory path
---@return string path the directory the OS designates for temporary files
function osmod.tmpdir() end

return osmod
