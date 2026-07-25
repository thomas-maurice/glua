---@meta regexp

---@class regexp
local regexp = {}

--- reports whether pattern matches text, raises on invalid pattern
---@param pattern string
---@param text string
---@return boolean
function regexp.match(pattern, text) end

--- returns the first match of pattern in text, raises on invalid pattern
---@param pattern string
---@param text string
---@return string
function regexp.find(pattern, text) end

--- returns all matches of pattern in text up to n, raises on invalid pattern
---@param pattern string
---@param text string
---@param n number
---@return string[]
function regexp.find_all(pattern, text, n) end

--- replaces the first match of pattern with replacement, raises on invalid pattern
---@param pattern string
---@param text string
---@param replacement string
---@return string
function regexp.replace(pattern, text, replacement) end

--- replaces all matches of pattern with replacement, raises on invalid pattern
---@param pattern string
---@param text string
---@param replacement string
---@return string
function regexp.replace_all(pattern, text, replacement) end

--- splits text by pattern into at most n parts, raises on invalid pattern
---@param pattern string
---@param text string
---@param n number
---@return string[]
function regexp.split(pattern, text, n) end

return regexp
