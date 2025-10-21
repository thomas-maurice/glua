---@meta regexp

---@class regexp
local regexp = {}

---@param pattern string The regular expression pattern
---@param text string The text to match against
---@return boolean True if pattern matches, false otherwise
---@return string|nil Error message if pattern is invalid
function regexp.match(pattern, text) end

---@param pattern string The regular expression pattern
---@param text string The text to search
---@return string The first match (empty string if no match)
---@return string|nil Error message if pattern is invalid
function regexp.find(pattern, text) end

---@param pattern string The regular expression pattern
---@param text string The text to search
---@param limit number Maximum number of matches (-1 for all)
---@return table Array of matches
---@return string|nil Error message if pattern is invalid
function regexp.find_all(pattern, text, limit) end

---@param pattern string The regular expression pattern
---@param text string The text to search
---@param replacement string The replacement string
---@return string The text with first match replaced
---@return string|nil Error message if pattern is invalid
function regexp.replace(pattern, text, replacement) end

---@param pattern string The regular expression pattern
---@param text string The text to search
---@param replacement string The replacement string
---@return string The text with all matches replaced
---@return string|nil Error message if pattern is invalid
function regexp.replace_all(pattern, text, replacement) end

---@param pattern string The regular expression pattern
---@param text string The text to split
---@param limit number Maximum number of splits (-1 for all)
---@return table Array of split parts
---@return string|nil Error message if pattern is invalid
function regexp.split(pattern, text, limit) end

return regexp
