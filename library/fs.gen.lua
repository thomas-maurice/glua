---@meta fs

---@class fs
local fs = {}

---@param path string The path to the file to read
---@return content string The file contents, or nil on error
---@return err string|nil Error message if reading failed
function fs.read_file(path) end

---@param path string The path to the file to write
---@param content string The content to write
---@return err string|nil Error message if writing failed, nil on success
function fs.write_file(path, content) end

---@param path string The path to check
---@return exists boolean True if the path exists, false otherwise
function fs.exists(path) end

---@param path string The directory path to create
---@return err string|nil Error message if creation failed, nil on success
function fs.mkdir(path) end

---@param path string The directory path to create
---@return err string|nil Error message if creation failed, nil on success
function fs.mkdir_all(path) end

---@param path string The path to remove
---@return err string|nil Error message if removal failed, nil on success
function fs.remove(path) end

---@param path string The path to remove recursively
---@return err string|nil Error message if removal failed, nil on success
function fs.remove_all(path) end

---@param path string The directory path to list
---@return entries table Array of entry names, or nil on error
---@return err string|nil Error message if listing failed
function fs.list(path) end

---@param path string The path to stat
---@return info table File info with name, size, is_dir, mode, mod_time, or nil on error
---@return err string|nil Error message if stat failed
function fs.stat(path) end

return fs
