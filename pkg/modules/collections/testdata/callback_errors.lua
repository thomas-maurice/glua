-- Test that a callback raising a Lua error propagates as a normal Lua error
-- (catchable with pcall), never a Go panic that would corrupt the host, and
-- that the wrapped message names the function and the failing index.
local c = require("collections")

local function boom() error("boom") end

local ok, err = pcall(c.map, {1, 2, 3}, boom)
assert(ok == false, "a raising callback must surface as a normal Lua error, not succeed")
assert(string.find(err, "collections.map") ~= nil, "error should name the failing function, got: " .. tostring(err))
assert(string.find(err, "index 1") ~= nil, "error should name the 1-based index the callback failed at, got: " .. tostring(err))
assert(string.find(err, "boom") ~= nil, "error should include the original callback message, got: " .. tostring(err))

-- The host must still be usable after a callback raised -- no corrupted state.
local ok2, sum = pcall(c.reduce, {1, 2, 3}, function(acc, v) return acc + v end, 0)
assert(ok2 == true and sum == 6, "the Lua state must still work normally after a previous callback raised")

-- Every callback-accepting function must protect the same way. Spot-check a
-- few beyond map.
local fns = {
  function() return c.filter({1}, boom) end,
  function() return c.find({1}, boom) end,
  function() return c.any({1}, boom) end,
  function() return c.all({1}, boom) end,
  function() return c.group_by({1}, boom) end,
  function() return c.sort_by({1}, boom) end,
  function() return c.partition({1}, boom) end,
}
for idx, f in ipairs(fns) do
  local fok, ferr = pcall(f)
  assert(fok == false, "callback-accepting function #" .. idx .. " must propagate a raising callback as an error")
  assert(string.find(ferr, "boom") ~= nil, "function #" .. idx .. " error should carry the original message, got: " .. tostring(ferr))
end

-- A non-function callback raises immediately, naming the parameter.
local nok, nerr = pcall(c.map, {1, 2}, "not a function")
assert(nok == false, "passing a non-function as the callback must raise")
assert(string.find(nerr, "fn") ~= nil, "error should name the offending parameter, got: " .. tostring(nerr))

return true
