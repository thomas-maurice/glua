---@meta filepath

---@class filepath
local filepath = {}

--- joins a table of path elements into a single path
---@param elem string[] table (array) of path segments to join with the OS path separator
---@return string path the joined and Clean-ed path
function filepath.join(elem) end

--- splits a path into directory and file components
---@param path string the path to split, e.g. "/a/b/c.txt"
---@return string dir everything up to and including the final separator, e.g. "/a/b/"
---@return string file everything after the final separator, e.g. "c.txt"
function filepath.split(path) end

--- returns the absolute form of the path, raises on error
---@param path string relative or absolute path to resolve against the process working directory
---@return string path the absolute, Clean-ed form of path
function filepath.abs(path) end

--- returns the file extension including the dot
---@param path string the path whose extension to extract
---@return string ext the file extension including the leading dot, or empty string if none
function filepath.ext(path) end

--- returns the last element of the path
---@param path string the path whose last element to extract
---@return string name the last path element, with trailing separators removed
function filepath.base(path) end

--- returns all but the last element of the path
---@param path string the path whose parent directory to extract
---@return string dir all but the last element of path
function filepath.dir(path) end

--- returns the shortest path equivalent to path
---@param path string the path to simplify, e.g. containing ".." or repeated separators
---@return string path the shortest path lexically equivalent to path
function filepath.clean(path) end

return filepath
