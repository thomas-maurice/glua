---@meta text

---@class text
local text = {}

--- greedily word-wraps a string to a maximum rune width per line
---@param s string the string to wrap; existing newlines are preserved as hard paragraph breaks and each paragraph wraps independently
---@param width number maximum runes per output line, must be >= 1; a single word longer than width is not split and overflows its line
---@return string out s re-flowed to width runes per line, using rune counts, not display width
function text.wrap(s, width) end

--- prefixes every line of a string
---@param s string the string to indent
---@param prefix string the string prepended to every line; a trailing empty line produced by a final newline in s is not prefixed
---@return string out s with prefix prepended to each line
function text.indent(s, prefix) end

--- removes the common leading-whitespace prefix shared by every non-blank line
---@param s string the string to dedent; leading whitespace is compared byte-literally (tabs and spaces are distinct, never expanded), and blank/whitespace-only lines do not participate in computing the common prefix and are normalized to empty in the output
---@return string out s with the longest common leading whitespace removed from every non-blank line
function text.dedent(s) end

--- truncates a string to a maximum rune width, appending an ellipsis if it was shortened
---@param s string the string to truncate; rune-counted, not display-width-counted
---@param width number maximum runes of the result, including the ellipsis; must be >= the rune length of ellipsis
---@param ellipsis string the marker appended when s is shortened, e.g. "..." or the single rune "…"; required, there is no default
---@return string out s unchanged if it already fits in width runes, otherwise a prefix of s plus ellipsis, exactly width runes long
function text.truncate(s, width, ellipsis) end

--- left-pads a string with a single rune to a minimum rune width
---@param s string the string to pad
---@param width number minimum runes of the result; s is never truncated if it is already this wide or wider; must be <= 1048576
---@param pad string the single rune to pad with; exactly one rune is required
---@return string out s left-padded with pad to width runes, or s unchanged if already >= width runes
function text.pad_left(s, width, pad) end

--- right-pads a string with a single rune to a minimum rune width
---@param s string the string to pad
---@param width number minimum runes of the result; s is never truncated if it is already this wide or wider; must be <= 1048576
---@param pad string the single rune to pad with; exactly one rune is required
---@return string out s right-padded with pad to width runes, or s unchanged if already >= width runes
function text.pad_right(s, width, pad) end

return text
