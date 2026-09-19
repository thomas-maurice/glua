-- Every documented error path, asserted via pcall.
local bit32 = require("bit32")

local function expect_error(fn, ...)
  local ok, err = pcall(fn, ...)
  assert(not ok, "expected an error but call succeeded")
  assert(type(err) == "string", "expected a string error message")
  return err
end

-- Non-integer float operand.
expect_error(bit32.band, 3.5, 0xFF)
expect_error(bit32.bnot, 3.5)

-- Operand beyond float64's exact integer range (2^60 > 2^53).
expect_error(bit32.band, 2 ^ 60, 0xFF)

-- Negative shift counts always raise.
expect_error(bit32.lshift, 1, -1)
expect_error(bit32.rshift, 1, -1)
expect_error(bit32.arshift, 1, -1)

-- test/set/clear bit index outside [0, 31] raises (unlike shift counts,
-- there is no clamping behaviour for a bit position).
expect_error(bit32.test, 1, 32)
expect_error(bit32.set, 1, 32)
expect_error(bit32.clear, 1, 32)
expect_error(bit32.test, 1, -1)

-- Wrong arity raises (arity is exact in this framework).
expect_error(bit32.bnot)
expect_error(bit32.lshift, 1)

return true
