-- test/set/clear on individual bit positions.
local bit32 = require("bit32")

assert(bit32.test(0x4, 2) == true, "bit 2 of 0x4 should be set")
assert(bit32.test(0x4, 0) == false, "bit 0 of 0x4 should be clear")
assert(bit32.test(0, 31) == false, "bit 31 of 0 should be clear")
assert(bit32.test(0x80000000, 31) == true, "bit 31 of 0x80000000 should be set")

assert(bit32.set(0, 0) == 1, "set(0, 0) should set the low bit")
assert(bit32.set(0, 31) == 0x80000000, "set(0, 31) should set the high bit")
assert(bit32.set(0xFF, 0) == 0xFF, "setting an already-set bit is a no-op")

assert(bit32.clear(0xFF, 0) == 0xFE, "clear(0xFF, 0) should clear the low bit")
assert(bit32.clear(0, 3) == 0, "clearing an already-clear bit is a no-op")

-- A well-known octal permission bit: owner-read is bit 8 (0o400).
local mode = 0x1A4 -- 0644 octal = 0x1A4
assert(bit32.test(mode, 8) == true, "0644 should have the owner-read bit set")
assert(bit32.test(mode, 1) == false, "0644 should not have the owner-write-execute mixup bit set")

return true
