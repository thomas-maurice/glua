---@meta bit32

---@class bit32
local bit32 = {}

--- bitwise AND of x and every value in rest, reduced modulo 2^32
---@param x number an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)
---@param ... number additional operands, same constraints as x
---@return number result the AND of all operands, in [0, 2^32)
function bit32.band(x, ...) end

--- bitwise OR of x and every value in rest, reduced modulo 2^32
---@param x number an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)
---@param ... number additional operands, same constraints as x
---@return number result the OR of all operands, in [0, 2^32)
function bit32.bor(x, ...) end

--- bitwise XOR of x and every value in rest, reduced modulo 2^32
---@param x number an exact integer in [-2^53, 2^53]; two's-complement negatives are honoured (-1 acts as 0xFFFFFFFF)
---@param ... number additional operands, same constraints as x
---@return number result the XOR of all operands, in [0, 2^32)
function bit32.bxor(x, ...) end

--- bitwise NOT of x, reduced modulo 2^32
---@param x number an exact integer in [-2^53, 2^53]
---@return number result the bitwise complement of x, in [0, 2^32); bnot(0) is 4294967295, never -1
function bit32.bnot(x) end

--- logical left shift of x by n bits
---@param x number an exact integer in [-2^53, 2^53]
---@param n number shift amount; must be a non-negative integer. n >= 32 yields 0
---@return number result x shifted left by n bits, in [0, 2^32)
function bit32.lshift(x, n) end

--- logical (zero-filling) right shift of x by n bits
---@param x number an exact integer in [-2^53, 2^53]
---@param n number shift amount; must be a non-negative integer. n >= 32 yields 0
---@return number result x shifted right by n bits with zero fill, in [0, 2^32)
function bit32.rshift(x, n) end

--- arithmetic (sign-propagating) right shift of x by n bits
---@param x number an exact integer in [-2^53, 2^53]
---@param n number shift amount; must be a non-negative integer. n >= 32 yields 0xFFFFFFFF if bit 31 of x was set, else 0
---@return number result x shifted right by n bits with bit-31 (sign) fill, in [0, 2^32)
function bit32.arshift(x, n) end

--- reports whether bit n of x is set
---@param x number an exact integer in [-2^53, 2^53]
---@param n number bit position, least significant bit is 0; must be in [0, 31]
---@return boolean set true if bit n of x is 1
function bit32.test(x, n) end

--- returns x with bit n set
---@param x number an exact integer in [-2^53, 2^53]
---@param n number bit position, least significant bit is 0; must be in [0, 31]
---@return number result x with bit n forced to 1, in [0, 2^32)
function bit32.set(x, n) end

--- returns x with bit n cleared
---@param x number an exact integer in [-2^53, 2^53]
---@param n number bit position, least significant bit is 0; must be in [0, 31]
---@return number result x with bit n forced to 0, in [0, 2^32)
function bit32.clear(x, n) end

return bit32
