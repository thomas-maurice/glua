---@meta regexp

---@class regexp
local regexp = {}

---@param string pattern The regular expression pattern
---@param string text The text to match against
---@return boolean True if pattern matches, false otherwise
---@return string|nil Error message if pattern is invalid
function regexp.match(pattern, text) end

---@param string pattern The regular expression pattern
---@param string text The text to search
---@return string The first match (empty string if no match)
---@return string|nil Error message if pattern is invalid
function regexp.find(pattern, text) end

---@param string pattern The regular expression pattern
---@param string text The text to search
---@param number limit Maximum number of matches (-1 for all)
---@return table Array of matches
---@return string|nil Error message if pattern is invalid
function regexp.find_all(pattern, text, limit) end

---@param string pattern The regular expression pattern
---@param string text The text to search
---@param string replacement The replacement string
---@return string The text with first match replaced
---@return string|nil Error message if pattern is invalid
function regexp.replace(pattern, text, replacement) end

---@param string pattern The regular expression pattern
---@param string text The text to search
---@param string replacement The replacement string
---@return string The text with all matches replaced
---@return string|nil Error message if pattern is invalid
function regexp.replace_all(pattern, text, replacement) end

---@param string pattern The regular expression pattern
---@param string text The text to split
---@param number limit Maximum number of splits (-1 for all)
---@return table Array of split parts
---@return string|nil Error message if pattern is invalid
function regexp.split(pattern, text, limit) end

return regexp
