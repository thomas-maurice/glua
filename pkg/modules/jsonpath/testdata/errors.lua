-- A malformed path/template must raise a catchable Lua error naming the
-- path, never crash the host, on every function that takes one. Wrong
-- arity must also raise (arity is exact in this framework).
local jsonpath = require("jsonpath")

local data = {spec = {replicas = 3}}

local function expect_error(fn, ...)
  local ok, err = pcall(fn, ...)
  assert(not ok, "expected an error but call succeeded")
  assert(type(err) == "string", "expected a string error message")
  return err
end

-- Unterminated bracket.
expect_error(jsonpath.query, data, "{.spec[")
expect_error(jsonpath.first, data, "{.spec[", nil)
expect_error(jsonpath.exists, data, "{.spec[")
expect_error(jsonpath.render, data, "{.spec[")

-- Malformed filter expression.
local err = expect_error(jsonpath.query, data, "{.spec[?(@.x==)]}")
assert(string.find(err, "jsonpath.query"), "error message should be prefixed with the function name")

-- Wrong arity raises. A missing argument always raises, on every function.
expect_error(jsonpath.query, data)
expect_error(jsonpath.first, data, ".spec.replicas")
expect_error(jsonpath.exists, data)
expect_error(jsonpath.render, data)

-- exists and render take no *lua.LState, so their arity is exactly enforced:
-- a surplus argument also raises (see the package doc's "Arity" section for
-- why query/first cannot make the same guarantee).
expect_error(jsonpath.exists, data, ".spec.replicas", "extra")
expect_error(jsonpath.render, data, ".spec.replicas", "extra")

-- A path/template that parses fine but fails at evaluation time (not a
-- missing-key miss) still raises even under AllowMissingKeys(true): an
-- out-of-bounds array index is a real error, not "field not set".
local arrayData = {spec = {containers = {{image = "nginx"}}}}
expect_error(jsonpath.query, arrayData, ".spec.containers[99].image")

-- Passing a value the Translator cannot convert (a function has no JSON
-- representation) raises a catchable error instead of propagating a Go
-- panic, on every function that accepts a data argument.
local unconvertible = function() end
expect_error(jsonpath.query, unconvertible, ".x")
expect_error(jsonpath.first, unconvertible, ".x", nil)
expect_error(jsonpath.exists, unconvertible, ".x")
expect_error(jsonpath.render, unconvertible, ".x")

return true
