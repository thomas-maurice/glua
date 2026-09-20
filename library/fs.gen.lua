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
---@param path string path to the file to read, absolute or relative to the process working directory
---@return string content the full file contents
function fs.read_file(path) end

--- writes content to a file, raises on error
---@param path string path to the file to write; created if missing, truncated if it exists
---@param content string the data to write, overwriting any existing content
function fs.write_file(path, content) end

--- reports whether a path exists
---@param path string the path to check
---@return boolean ok true if a file or directory exists at path
function fs.exists(path) end

--- creates a directory, raises on error
---@param path string path of the directory to create; the parent directory must already exist
function fs.mkdir(path) end

--- creates a directory and all parents, raises on error
---@param path string path of the directory to create, along with any missing parent directories
function fs.mkdir_all(path) end

--- removes a file or empty directory, raises on error
---@param path string path to the file or empty directory to remove
function fs.remove(path) end

--- removes a path and all its contents, raises on error
---@param path string path to remove recursively, including all files and subdirectories
function fs.remove_all(path) end

--- lists all entries in a directory, raises on error
---@param path string path of the directory to list
---@return string[] names table (array) of entry names directly inside path, not recursive
function fs.list(path) end

--- returns file information, raises on error
---@param path string path of the file or directory to stat
---@return fs.FileInfo info FileInfo table describing name, size, is_dir, mode and mod_time
function fs.stat(path) end

return fs
