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
---@param parts table table (array) of values to join; non-string values are coerced via tostring
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

--- removes leading and trailing whitespace from a string
---@param s string the string to trim
---@return string out s with leading and trailing Unicode whitespace removed
function strings.trim_space(s) end

--- removes a leading prefix from a string, if present
---@param s string the string to trim
---@param prefix string the prefix to remove
---@return string out s with prefix removed, or s unchanged if it does not start with prefix
function strings.trim_prefix(s, prefix) end

--- removes a trailing suffix from a string, if present
---@param s string the string to trim
---@param suffix string the suffix to remove
---@return string out s with suffix removed, or s unchanged if it does not end with suffix
function strings.trim_suffix(s, suffix) end

--- splits a string around runs of whitespace
---@param s string the string to split
---@return string[] parts table (array) of substrings between runs of Unicode whitespace; leading/trailing whitespace produces no empty entries
function strings.fields(s) end

--- repeats a string count times
---@param s string the string to repeat
---@param count number number of repetitions; must be in [0, 1048576]
---@return string out s repeated count times
function strings.rep(s, count) end

--- returns the 1-based position of the first occurrence of substr in s, or 0 if absent
---@param s string the string to search
---@param substr string the substring to look for
---@return number pos 1-based index of the first occurrence, or 0 if not found; NOTE this deviates from Go's strings.Index, which is 0-based and returns -1 for not found — Lua's own string.find convention is used instead so the result composes with string.sub
function strings.index(s, substr) end

--- returns the 1-based position of the last occurrence of substr in s, or 0 if absent
---@param s string the string to search
---@param substr string the substring to look for
---@return number pos 1-based index of the last occurrence, or 0 if not found; same 1-based/0-absent convention as index, deviating from Go's strings.LastIndex
function strings.last_index(s, substr) end

--- reports whether two strings are equal under simple Unicode case-folding
---@param a string the first string
---@param b string the second string
---@return boolean ok true if a and b are equal under case-insensitive comparison
function strings.equal_fold(a, b) end

--- title-cases the first rune of every whitespace-separated word
---@param s string the string to title-case
---@return string out s with the first rune of each word title-cased; ASCII/first-rune only, not language-aware (does not implement locale casing exceptions)
function strings.title(s) end

--- slices s around the first instance of sep
---@param s string the string to slice
---@param sep string the separator to find
---@return string before the portion of s before the first occurrence of sep, or all of s if sep is not found
---@return string after the portion of s after the first occurrence of sep, or empty if sep is not found
---@return boolean found true if sep occurs in s
function strings.cut(s, sep) end

--- splits a string by separator into at most n substrings
---@param s string the string to split
---@param sep string the separator; if empty, splits after every UTF-8 character
---@param n number maximum number of substrings: n > 0 limits the result to at most n elements (the last one unsplit); n == 0 returns an empty table; n < 0 splits all occurrences, like split
---@return string[] parts table (array) of at most n substrings
function strings.split_n(s, sep, n) end

return strings
