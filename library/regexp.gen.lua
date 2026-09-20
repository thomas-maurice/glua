---@meta regexp

---@class regexp
local regexp = {}

--- reports whether pattern matches text, raises on invalid pattern
---@param pattern string an RE2 regular expression
---@param text string the string to test
---@return boolean ok true if pattern matches anywhere in text
function regexp.match(pattern, text) end

--- returns the first match of pattern in text, raises on invalid pattern
---@param pattern string an RE2 regular expression
---@param text string the string to search
---@return string match the first matching substring, or empty string if no match
function regexp.find(pattern, text) end

--- returns all matches of pattern in text up to n, raises on invalid pattern
---@param pattern string an RE2 regular expression
---@param text string the string to search
---@param n number maximum number of matches to return; a negative value returns all matches
---@return string[] matches table (array) of matching substrings, in order of appearance
function regexp.find_all(pattern, text, n) end

--- replaces the first match of pattern with replacement, raises on invalid pattern
---@param pattern string an RE2 regular expression
---@param text string the string to search
---@param replacement string the literal string to substitute for the first match
---@return string result text with its first match (if any) replaced
function regexp.replace(pattern, text, replacement) end

--- replaces all matches of pattern with replacement, raises on invalid pattern
---@param pattern string an RE2 regular expression
---@param text string the string to search
---@param replacement string the literal string to substitute for every match
---@return string result text with all matches replaced
function regexp.replace_all(pattern, text, replacement) end

--- splits text by pattern into at most n parts, raises on invalid pattern
---@param pattern string an RE2 regular expression used as the separator
---@param text string the string to split
---@param n number maximum number of substrings to return; a negative value returns all substrings
---@return string[] parts table (array) of substrings between matches of pattern
function regexp.split(pattern, text, n) end

return regexp
