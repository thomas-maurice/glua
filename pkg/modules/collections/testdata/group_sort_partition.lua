-- Test group_by, sort_by (including stability and the mixed-key-type raise)
-- and partition.
local c = require("collections")

-- group_by: groups preserve encounter order within each group.
local words = {"apple", "banana", "avocado", "blueberry", "cherry"}
local byFirst = c.group_by(words, function(w) return w:sub(1, 1) end)
assert(#byFirst.a == 2 and byFirst.a[1] == "apple" and byFirst.a[2] == "avocado", "group_by should preserve order within a group")
assert(#byFirst.b == 2 and byFirst.b[1] == "banana" and byFirst.b[2] == "blueberry", "group_by should collect all matching elements")
assert(#byFirst.c == 1 and byFirst.c[1] == "cherry", "group_by should handle a singleton group")

local emptyGroups = c.group_by({}, function(w) return w end)
assert(type(emptyGroups) == "table", "group_by({}) must return an empty table, not nil")

-- sort_by: ascending by key, stable on ties (equal-key elements keep their
-- original relative order -- this is the whole point of using SliceStable).
local people = {
  {name = "b", age = 30},
  {name = "a", age = 30},
  {name = "c", age = 20},
}
local byAge = c.sort_by(people, function(p) return p.age end)
assert(byAge[1].name == "c", "sort_by should sort ascending by key")
assert(byAge[2].name == "b" and byAge[3].name == "a", "sort_by must be stable: equal keys (age 30) keep b-before-a input order")

-- sort_by on strings works the same way.
local byName = c.sort_by({"banana", "apple", "cherry"}, function(s) return s end)
assert(byName[1] == "apple" and byName[2] == "banana" and byName[3] == "cherry", "sort_by should support string keys")

-- Mixed key types raise rather than inventing a cross-type ordering.
local ok, err = pcall(c.sort_by, {1, "two", 3}, function(v) return v end)
assert(ok == false, "sort_by must raise when keys mix numbers and strings")
assert(string.find(err, "key type mismatch") ~= nil, "sort_by error should say 'key type mismatch', got: " .. tostring(err))
assert(string.find(err, "index 2") ~= nil, "sort_by error should name the offending index, got: " .. tostring(err))

-- A key that is neither number nor string also raises.
local ok2, err2 = pcall(c.sort_by, {{1}, {2}}, function(v) return v end)
assert(ok2 == false, "sort_by must raise when the key itself is not a number or string")

local emptySort = c.sort_by({}, function(v) return v end)
assert(type(emptySort) == "table" and #emptySort == 0, "sort_by({}) must return an empty table, not nil")

-- partition: matching + rest, both re-indexed, order preserved.
local matching, rest = c.partition({1, 2, 3, 4, 5}, function(v) return v % 2 == 0 end)
assert(#matching == 2 and matching[1] == 2 and matching[2] == 4, "partition should collect matching elements in order")
assert(#rest == 3 and rest[1] == 1 and rest[2] == 3 and rest[3] == 5, "partition should collect the remainder in order")

local m0, r0 = c.partition({}, function() return true end)
assert(type(m0) == "table" and #m0 == 0 and type(r0) == "table" and #r0 == 0, "partition({}) must return two empty tables, not nil")

return true
