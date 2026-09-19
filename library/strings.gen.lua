---@meta strings

---@class strings
local strings = {}

--- checks if string has prefix
---@param s string the string to test
---@param prefix string the prefix to look for
---@return boolean ok true if s starts with prefix
function strings.has_prefix(s, prefix) end

--- checks if string has suffix
---@param s string the string to test
---@param suffix string the suffix to look for
---@return boolean ok true if s ends with suffix
function strings.has_suffix(s, suffix) end

--- removes cutset characters from both ends of a string
---@param s string the string to trim
---@param cutset string the set of characters to remove from both ends
---@return string out s with leading and trailing cutset characters removed
function strings.trim(s, cutset) end

--- removes cutset characters from the left end of a string
---@param s string the string to trim
---@param cutset string the set of characters to remove from the left
---@return string out s with leading cutset characters removed
function strings.trim_left(s, cutset) end

--- removes cutset characters from the right end of a string
---@param s string the string to trim
---@param cutset string the set of characters to remove from the right
---@return string out s with trailing cutset characters removed
function strings.trim_right(s, cutset) end

--- splits a string by separator into a table
---@param s string the string to split
---@param sep string the separator; if empty, splits after every UTF-8 character
---@return string[] parts table (array) of substrings between occurrences of sep
function strings.split(s, sep) end

--- joins a table of strings with a separator, coercing values to strings
---@param parts gopher-lua.LTable table (array) of values to join; non-string values are coerced via tostring
---@param sep string the separator to place between elements
---@return string out the joined string
function strings.join(parts, sep) end

--- converts a string to uppercase
---@param s string the string to convert
---@return string out s with all letters mapped to upper case
function strings.to_upper(s) end

--- converts a string to lowercase
---@param s string the string to convert
---@return string out s with all letters mapped to lower case
function strings.to_lower(s) end

--- checks if a string contains a substring
---@param s string the string to search
---@param substr string the substring to look for
---@return boolean ok true if substr appears anywhere in s
function strings.contains(s, substr) end

--- counts occurrences of substr in s
---@param s string the string to search
---@param substr string the non-overlapping substring to count; empty string counts UTF-8 characters plus one
---@return number n the number of non-overlapping instances of substr in s
function strings.count(s, substr) end

--- replaces occurrences of old with new in s
---@param s string the string to search
---@param old string the substring to replace
---@param new string the replacement substring
---@param n number maximum number of replacements; a negative value replaces all occurrences
---@return string out s with up to n occurrences of old replaced by new
function strings.replace(s, old, new, n) end

return strings
