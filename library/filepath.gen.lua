---@meta filepath

---@class filepath
local filepath = {}

--- joins a table of path elements into a single path
---@param elem string[]
---@return string
function filepath.join(elem) end

--- splits a path into directory and file components
---@param path string
---@return string
---@return string
function filepath.split(path) end

--- returns the absolute form of the path, raises on error
---@param path string
---@return string
function filepath.abs(path) end

--- returns the file extension including the dot
---@param path string
---@return string
function filepath.ext(path) end

--- returns the last element of the path
---@param path string
---@return string
function filepath.base(path) end

--- returns all but the last element of the path
---@param path string
---@return string
function filepath.dir(path) end

--- returns the shortest path equivalent to path
---@param path string
---@return string
function filepath.clean(path) end

return filepath
