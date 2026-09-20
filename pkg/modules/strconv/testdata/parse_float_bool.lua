-- parse_float and every parse_bool accepted spelling.
local strconv = require("strconv")

assert(strconv.parse_float("3.14") == 3.14, "parse_float basic decimal")
assert(strconv.parse_float("-2.5e3") == -2500, "parse_float exponent notation")
-- gopher-lua's math.huge is math.MaxFloat64, not IEEE +Inf, so compare
-- against 1/0 (real IEEE division) instead.
assert(strconv.parse_float("Inf") == 1 / 0, "parse_float accepts the Inf spelling")
assert(strconv.parse_float("-Inf") == -1 / 0, "parse_float accepts the -Inf spelling")
assert(strconv.parse_float("NaN") ~= strconv.parse_float("NaN"), "parse_float accepts the NaN spelling (NaN ~= NaN)")

local bool_true = {"1", "t", "T", "TRUE", "true", "True"}
for _, s in ipairs(bool_true) do
  assert(strconv.parse_bool(s) == true, "parse_bool should accept '" .. s .. "' as true")
end

local bool_false = {"0", "f", "F", "FALSE", "false", "False"}
for _, s in ipairs(bool_false) do
  assert(strconv.parse_bool(s) == false, "parse_bool should accept '" .. s .. "' as false")
end

return true
