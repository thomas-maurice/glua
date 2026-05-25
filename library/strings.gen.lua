---@meta strings

---@class strings
local strings = {}

--- checks if string has prefix
---@param s string
---@param prefix string
---@return boolean
function strings.has_prefix(s, prefix) end

--- checks if string has suffix
---@param s string
---@param suffix string
---@return boolean
function strings.has_suffix(s, suffix) end

--- removes cutset characters from both ends of a string
---@param s string
---@param cutset string
---@return string
function strings.trim(s, cutset) end

--- removes cutset characters from the left end of a string
---@param s string
---@param cutset string
---@return string
function strings.trim_left(s, cutset) end

--- removes cutset characters from the right end of a string
---@param s string
---@param cutset string
---@return string
function strings.trim_right(s, cutset) end

--- splits a string by separator into a table
---@param s string
---@param sep string
---@return string[]
function strings.split(s, sep) end

--- joins a table of strings with a separator, coercing values to strings
---@param parts gopher-lua.LTable
---@param sep string
---@return string
function strings.join(parts, sep) end

--- converts a string to uppercase
---@param s string
---@return string
function strings.to_upper(s) end

--- converts a string to lowercase
---@param s string
---@return string
function strings.to_lower(s) end

--- checks if a string contains a substring
---@param s string
---@param substr string
---@return boolean
function strings.contains(s, substr) end

--- counts occurrences of substr in s
---@param s string
---@param substr string
---@return number
function strings.count(s, substr) end

--- replaces occurrences of old with new in s
---@param s string
---@param old string
---@param new string
---@param n number
---@return string
function strings.replace(s, old, new, n) end

return strings
