-- parse_int / atoi: base handling, prefix inference, and round-trip with
-- format_int across representative bases.
local strconv = require("strconv")

assert(strconv.parse_int("644", 8) == 420, "parse_int with explicit base 8")
assert(strconv.parse_int("2A", 16) == 42, "parse_int with explicit base 16")
assert(strconv.atoi("123") == 123, "atoi is base-10 shorthand")
assert(strconv.atoi("-45") == -45, "atoi handles negative numbers")

-- base == 0 infers from prefix.
assert(strconv.parse_int("0x1A", 0) == 26, "parse_int base 0 infers hex from 0x prefix")
assert(strconv.parse_int("0b101", 0) == 5, "parse_int base 0 infers binary from 0b prefix")
assert(strconv.parse_int("0o17", 0) == 15, "parse_int base 0 infers octal from 0o prefix")
assert(strconv.parse_int("42", 0) == 42, "parse_int base 0 defaults to decimal with no prefix")

-- Underscore digit separators are only accepted when base == 0.
assert(strconv.parse_int("1_000", 0) == 1000, "digit separators allowed when base is 0")

-- Round-trip format_int . parse_int over bases 2, 8, 10, 16, 36.
local bases = {2, 8, 10, 16, 36}
local values = {0, 1, 42, 420, 1000000}
for _, base in ipairs(bases) do
  for _, v in ipairs(values) do
    local s = strconv.format_int(v, base)
    local back = strconv.parse_int(s, base)
    assert(back == v, "round-trip failed for value " .. v .. " base " .. base .. " (got " .. s .. " -> " .. back .. ")")
  end
end

return true
