---@meta strings

---@class strings
local strings = {}

---@param s string The string to check
---@param prefix string The prefix to look for
---@return boolean True if string has prefix
function strings.has_prefix(s, prefix) end

---@param s string The string to check
---@param suffix string The suffix to look for
---@return boolean True if string has suffix
function strings.has_suffix(s, suffix) end

---@param s string The string to trim
---@param cutset string The characters to remove
---@return string The trimmed string
function strings.trim(s, cutset) end

---@param s string The string to trim
---@param cutset string The characters to remove
---@return string The trimmed string
function strings.trim_left(s, cutset) end

---@param s string The string to trim
---@param cutset string The characters to remove
---@return string The trimmed string
function strings.trim_right(s, cutset) end

---@param s string The string to split
---@param sep string The separator
---@return table Array of split parts
function strings.split(s, sep) end

---@param parts table Array of strings to join
---@param sep string The separator
---@return string The joined string
function strings.join(parts, sep) end

---@param s string The string to convert
---@return string The uppercase string
function strings.to_upper(s) end

---@param s string The string to convert
---@return string The lowercase string
function strings.to_lower(s) end

---@param s string The string to search
---@param substr string The substring to find
---@return boolean True if string contains substring
function strings.contains(s, substr) end

---@param s string The string to search
---@param substr string The substring to count
---@return number The count of occurrences
function strings.count(s, substr) end

---@param s string The string to search
---@param old string The substring to replace
---@param new string The replacement string
---@param n number Number of replacements (-1 for all)
---@return string The string with replacements
function strings.replace(s, old, new, n) end

return strings
