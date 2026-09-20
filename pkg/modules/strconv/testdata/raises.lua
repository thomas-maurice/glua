-- Every documented error path, asserted via pcall.
local strconv = require("strconv")

local function expect_error(fn, ...)
  local ok, err = pcall(fn, ...)
  assert(not ok, "expected an error but call succeeded")
  assert(type(err) == "string", "expected a string error message")
  return err
end

-- parse_int: syntax error and the 2^53 range guard.
expect_error(strconv.parse_int, "12a", 10)
expect_error(strconv.parse_int, "9007199254740993", 10) -- 2^53 + 1

-- atoi shares parse_int's error paths.
expect_error(strconv.atoi, "not a number")

-- parse_float: syntax error.
expect_error(strconv.parse_float, "not a float")

-- parse_bool: anything outside the accepted spelling set.
expect_error(strconv.parse_bool, "yes")
expect_error(strconv.parse_bool, "")

-- format_int: non-integral n, oversized n, and bad base.
-- 2^54 is used (not 2^53 + 1) because a literal that close to 2^53 rounds to
-- 2^53 itself at parse time — float64 cannot even represent 2^53 + 1.
expect_error(strconv.format_int, 3.5, 10)
expect_error(strconv.format_int, 2 ^ 54, 10)
expect_error(strconv.format_int, 42, 1)
expect_error(strconv.format_int, 42, 37)

-- format_float: unknown fmt verb and prec < -1.
expect_error(strconv.format_float, 1.5, "z", -1)
expect_error(strconv.format_float, 1.5, "f", -2)

-- unquote: invalid literal.
expect_error(strconv.unquote, "not quoted")

return true
