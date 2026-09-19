---@meta random

---@class random
---@field ALPHANUM string uppercase+lowercase letters and digits, for use as a random.string charset
local random = {}

--- returns n cryptographically random bytes as a raw 8-bit-clean string
---@param n number number of bytes to generate; must be in [1, 1048576]
---@return string data n random bytes
function random.bytes(n) end

--- returns n random bytes, hex-encoded
---@param n number number of underlying random bytes; must be in [1, 1048576]
---@return string hex lowercase hex encoding, 2n characters long
function random.hex(n) end

--- returns n random bytes as an unpadded, URL-safe base64 token
---@param n number number of underlying random bytes; must be in [1, 1048576]
---@return string token base64url encoding without padding, safe to use directly in a URL or header
function random.token(n) end

--- returns a string of n characters drawn uniformly from charset
---@param n number number of characters to generate; must be in [1, 1048576]
---@param charset string the runes to draw from, as a UTF-8 string; must not be empty. Duplicate runes are not de-duplicated and are therefore weighted
---@return string s n runes drawn from charset
function random.string(n, charset) end

--- returns a uniform, bias-free random integer in [min, max], inclusive of both ends (like Lua's math.random, NOT Go's half-open convention)
---@param min number inclusive lower bound; must be an integral value
---@param max number inclusive upper bound; must be an integral value, >= min, with max - min < 2^53
---@return number n a random integer in [min, max]
function random.int(min, max) end

return random
