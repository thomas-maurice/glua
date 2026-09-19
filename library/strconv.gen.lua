---@meta strconv

---@class strconv
local strconv = {}

--- parses a string as an integer, raising on syntax/range errors instead of returning nil
---@param s string the string to parse
---@param base number numeric base 2-36, or 0 to infer from s's prefix (0x/0o/0b/leading 0); digit-separator underscores are only accepted when base is 0
---@return number n the parsed integer; raises if it cannot fit exactly in a float64 (magnitude > 2^53)
function strconv.parse_int(s, base) end

--- parses a string as a base-10 integer; shorthand for parse_int(s, 10)
---@param s string the string to parse
---@return number n the parsed integer; raises if it cannot fit exactly in a float64 (magnitude > 2^53)
function strconv.atoi(s) end

--- parses a string as a floating-point number, raising on syntax errors instead of returning nil
---@param s string the string to parse; accepts decimal, exponent, hex-float, and Inf/NaN spellings
---@return number f the parsed float
function strconv.parse_float(s) end

--- parses a string as a boolean, accepting Go's spelling set
---@param s string one of 1 t T TRUE true True 0 f F FALSE false False
---@return boolean b the parsed boolean
function strconv.parse_bool(s) end

--- formats an integer in the given base; has no Lua 5.1 equivalent
---@param n number the value to format; must be an exact integer with magnitude <= 2^53
---@param base number numeric base, 2-36; digits above 9 are lowercase letters
---@return string s n rendered in the given base
function strconv.format_int(n, base) end

--- formats a float with an explicit format and precision; has no Lua 5.1 equivalent
---@param f number the value to format
---@param fmt string one of b e E f g G x X, matching Go's strconv.FormatFloat verbs
---@param prec number digits of precision; -1 selects the shortest representation that round-trips exactly through parse_float
---@return string s f rendered per fmt and prec
function strconv.format_float(f, fmt, prec) end

--- returns a Go-syntax double-quoted string literal for s, with non-printables escaped
---@param s string the raw string to quote
---@return string quoted a double-quoted Go string literal representing s
function strconv.quote(s) end

--- parses a Go-syntax string, raw string, or rune literal back into its raw value
---@param s string a "..." double-quoted, `...` raw, or 'c' rune literal
---@return string s the decoded raw string
function strconv.unquote(s) end

return strconv
