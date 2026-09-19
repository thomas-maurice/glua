-- Two's-complement negative handling and modulo-2^32 reduction of oversized
-- operands.
local bit32 = require("bit32")

-- -1 is two's-complement 0xFFFFFFFF, so masking with 0xFF yields 255.
assert(bit32.band(-1, 0xFF) == 255, "band(-1, 0xFF) should be 255 via two's complement")

-- An operand of 0x1FFFFFFFF (2^33 - 1) reduces mod 2^32 to 0xFFFFFFFF before
-- the AND is applied.
assert(bit32.band(0x1FFFFFFFF, 0xFF) == 0xFF, "band should reduce oversized operand mod 2^32")

return true
