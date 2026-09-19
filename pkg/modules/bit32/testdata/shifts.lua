-- Shift behaviour: lshift/rshift/arshift, including n >= 32 clamping and the
-- logical-vs-arithmetic distinction on a value with the high bit set.
local bit32 = require("bit32")

assert(bit32.lshift(1, 4) == 16, "lshift should shift left")
assert(bit32.lshift(1, 32) == 0, "lshift by >= 32 should yield 0")
assert(bit32.lshift(1, 40) == 0, "lshift by > 32 should yield 0")

assert(bit32.rshift(16, 4) == 1, "rshift should shift right")
assert(bit32.rshift(0xFFFFFFFF, 32) == 0, "rshift by >= 32 should yield 0")

-- The whole point of arshift: it must genuinely differ from rshift on a
-- value with the high bit set.
local logical = bit32.rshift(0x80000000, 4)
local arithmetic = bit32.arshift(0x80000000, 4)
assert(logical == 0x08000000, "rshift should zero-fill, got " .. logical)
assert(arithmetic == 0xF8000000, "arshift should sign-fill, got " .. arithmetic)
assert(logical ~= arithmetic, "rshift and arshift must differ for a high-bit-set value")

-- n >= 32 for arshift depends on the sign bit of x.
assert(bit32.arshift(0x80000000, 32) == 0xFFFFFFFF, "arshift by >= 32 with sign bit set should yield all-ones")
assert(bit32.arshift(0x7FFFFFFF, 32) == 0, "arshift by >= 32 without sign bit should yield 0")

-- No high bit set: arshift behaves like rshift.
assert(bit32.arshift(0x0F, 2) == bit32.rshift(0x0F, 2), "arshift should match rshift when the sign bit is clear")

return true
