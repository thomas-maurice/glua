-- Basic band/bor/bxor/bnot behaviour, including the variadic tail.
local bit32 = require("bit32")

assert(bit32.band(0x4, 0x2) == 0, "band of disjoint bits should be 0")
assert(bit32.band(0xFF, 0x0F) == 0x0F, "band should mask")
assert(bit32.band(0xFF, 0x0F, 0x03) == 0x03, "band should reduce over variadic tail")

assert(bit32.bor(0x4, 0x2) == 6, "bor should combine bits")
assert(bit32.bor(0x1, 0x2, 0x4) == 7, "bor should reduce over variadic tail")

assert(bit32.bxor(0xFF, 0x0F) == 0xF0, "bxor should flip masked bits")
assert(bit32.bxor(0x1, 0x1) == 0, "bxor of equal values is 0")

assert(bit32.bnot(0) == 4294967295, "bnot(0) should be 4294967295, not -1")
assert(bit32.bnot(0xFFFFFFFF) == 0, "bnot(0xFFFFFFFF) should be 0")

-- Round-trip: bnot(bnot(x)) == x
local samples = {0, 1, 0xFFFFFFFF, 0x80000000, 12345}
for _, x in ipairs(samples) do
  assert(bit32.bnot(bit32.bnot(x)) == x, "bnot round-trip failed for " .. x)
end

return true
