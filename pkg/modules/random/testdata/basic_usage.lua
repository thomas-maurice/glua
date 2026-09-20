-- Exercises the documented example composition end-to-end.
local random = require("random")

local apiKey = random.token(32)
assert(type(apiKey) == "string" and #apiKey > 0, "token should be a non-empty string")

local salt = random.hex(16)
assert(#salt == 32, "hex(16) should be 32 hex characters")

local pin = random.string(6, "0123456789")
assert(#pin == 6, "string(6, ...) should be 6 characters long")
for i = 1, #pin do
  local c = pin:sub(i, i)
  assert(string.find("0123456789", c, 1, true), "pin contains a character outside the charset: " .. c)
end

local dieRoll = random.int(1, 6)
assert(dieRoll >= 1 and dieRoll <= 6, "int(1,6) out of range: " .. tostring(dieRoll))

-- min == max is a legal edge case: always returns that single value.
assert(random.int(42, 42) == 42, "int(42, 42) must always return 42")

return true
