---@meta filepath

---@class filepath
local filepath = {}

---@param string ... Path elements to join
---@return string The joined path
function filepath.join(...) end

---@param string path The path to split
---@return string The directory part
---@return string The file part
function filepath.split(path) end

---@param string path The path to make absolute
---@return string The absolute path
---@return string|nil Error message if operation failed
function filepath.abs(path) end

---@param string path The file path
---@return string The file extension (including the dot)
function filepath.ext(path) end

---@param string path The file path
---@return string The base name
function filepath.base(path) end

---@param string path The file path
---@return string The directory path
function filepath.dir(path) end

---@param string path The path to clean
---@return string The cleaned path
function filepath.clean(path) end

return filepath
