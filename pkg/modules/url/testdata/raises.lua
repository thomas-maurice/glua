-- Every documented error path, asserted via pcall.
local url = require("url")

local function expect_error(fn, ...)
  local ok, err = pcall(fn, ...)
  assert(not ok, "expected an error but call succeeded")
  assert(type(err) == "string", "expected a string error message")
  return err
end

expect_error(url.parse, "://")
expect_error(url.resolve, "://", "c")
expect_error(url.resolve, "https://example.com/a", "://")
expect_error(url.query_unescape, "%zz")
expect_error(url.path_unescape, "%zz")
expect_error(url.parse_query, "%zz=1")

-- build_query: bad key type and bad value shape.
local err = expect_error(url.build_query, {bad = 1})
assert(err:find("bad", 1, true) ~= nil, "build_query error names the offending key: " .. err)

return true
