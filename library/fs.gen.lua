---@meta fs

---@class fs
local fs = {}

---@param path string The path to the file to read
---@return string content The file contents, or nil on error
---@return string|nil err Error message if reading failed
function fs.read_file(path) end

---@param path string The path to the file to write
---@param content string The content to write
---@return string|nil err Error message if writing failed, nil on success
function fs.write_file(path, content) end

---@param path string The path to check
---@return boolean exists True if the path exists, false otherwise
function fs.exists(path) end

---@param path string The directory path to create
---@return string|nil err Error message if creation failed, nil on success
function fs.mkdir(path) end

---@param path string The directory path to create
---@return string|nil err Error message if creation failed, nil on success
function fs.mkdir_all(path) end

---@param path string The path to remove
---@return string|nil err Error message if removal failed, nil on success
function fs.remove(path) end

---@param path string The path to remove recursively
---@return string|nil err Error message if removal failed, nil on success
function fs.remove_all(path) end

---@param path string The directory path to list
---@return table entries Array of entry names, or nil on error
---@return string|nil err Error message if listing failed
function fs.list(path) end

---@param path string The path to stat
---@return table info File info with name, size, is_dir, mode, mod_time, or nil on error
---@return string|nil err Error message if stat failed
function fs.stat(path) end

return fs
