---@meta fs

---@class fs.FileInfo
---@field is_dir boolean
---@field mod_time number
---@field mode number
---@field name string
---@field size number

---@class fs
local fs = {}

--- reads the entire contents of a file, raises on error
---@param path string
---@return string
function fs.read_file(path) end

--- writes content to a file, raises on error
---@param path string
---@param content string
function fs.write_file(path, content) end

--- reports whether a path exists
---@param path string
---@return boolean
function fs.exists(path) end

--- creates a directory, raises on error
---@param path string
function fs.mkdir(path) end

--- creates a directory and all parents, raises on error
---@param path string
function fs.mkdir_all(path) end

--- removes a file or empty directory, raises on error
---@param path string
function fs.remove(path) end

--- removes a path and all its contents, raises on error
---@param path string
function fs.remove_all(path) end

--- lists all entries in a directory, raises on error
---@param path string
---@return string[]
function fs.list(path) end

--- returns file information, raises on error
---@param path string
---@return fs.FileInfo
function fs.stat(path) end

return fs
