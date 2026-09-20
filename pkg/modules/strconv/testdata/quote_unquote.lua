-- quote/unquote round-trip, including embedded quotes and newlines.
local strconv = require("strconv")

local cases = {
  "hello",
  "",
  'has "embedded" quotes',
  "has\nnewlines\ttabs",
}
for _, s in ipairs(cases) do
  local q = strconv.quote(s)
  local back = strconv.unquote(q)
  assert(back == s, "unquote(quote(s)) should equal s for: " .. s)
end

-- unquote also accepts raw-string and rune literals.
assert(strconv.unquote("`raw string`") == "raw string", "unquote should accept backtick raw strings")
assert(strconv.unquote("'c'") == "c", "unquote should accept a rune literal")

return true
