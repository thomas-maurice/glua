-- format_float: verbs and shortest round-trip precision.
local strconv = require("strconv")

assert(strconv.format_float(1.5, "f", 2) == "1.50", "format_float fixed precision")
assert(strconv.format_float(1234.5, "e", 2) == "1.23e+03", "format_float exponent verb")

-- prec == -1 gives the shortest representation that round-trips exactly.
local third = 1 / 3
local s = strconv.format_float(third, "g", -1)
assert(strconv.parse_float(s) == third, "shortest round-trip formatting must parse back to the same float")

return true
