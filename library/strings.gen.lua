---@meta strings

---@class strings
local strings = {}

---@param string s The string to check
---@param string prefix The prefix to look for
---@return boolean True if string has prefix
function strings.has_prefix(s, prefix) end

---@param string s The string to check
---@param string suffix The suffix to look for
---@return boolean True if string has suffix
function strings.has_suffix(s, suffix) end

---@param string s The string to trim
---@param string cutset The characters to remove
---@return string The trimmed string
function strings.trim(s, cutset) end

---@param string s The string to trim
---@param string cutset The characters to remove
---@return string The trimmed string
function strings.trim_left(s, cutset) end

---@param string s The string to trim
---@param string cutset The characters to remove
---@return string The trimmed string
function strings.trim_right(s, cutset) end

---@param string s The string to split
---@param string sep The separator
---@return table Array of split parts
function strings.split(s, sep) end

---@param table parts Array of strings to join
---@param string sep The separator
---@return string The joined string
function strings.join(parts, sep) end

---@param string s The string to convert
---@return string The uppercase string
function strings.to_upper(s) end

---@param string s The string to convert
---@return string The lowercase string
function strings.to_lower(s) end

---@param string s The string to search
---@param string substr The substring to find
---@return boolean True if string contains substring
function strings.contains(s, substr) end

---@param string s The string to search
---@param string substr The substring to count
---@return number The count of occurrences
function strings.count(s, substr) end

---@param string s The string to search
---@param string old The substring to replace
---@param string new The replacement string
---@param number n Number of replacements (-1 for all)
---@return string The string with replacements
function strings.replace(s, old, new, n) end

return strings
