-- Test the array-mode functional trio plus find/any/all: basic behaviour,
-- the (value, index) callback argument order, and empty-input results.
local c = require("collections")

-- map: transforms every element, callback sees (v, i).
local doubled = c.map({1, 2, 3}, function(v, i) return v * 2 + i end)
assert(doubled[1] == 3 and doubled[2] == 6 and doubled[3] == 9, "map should call fn(v, i) per element")

-- filter: keeps only truthy elements, re-indexed contiguously.
local evens = c.filter({1, 2, 3, 4, 5}, function(v) return v % 2 == 0 end)
assert(#evens == 2 and evens[1] == 2 and evens[2] == 4, "filter should keep only matching elements, re-indexed")

-- reduce: required init, folds left-to-right, sees (acc, v, i).
local sum = c.reduce({1, 2, 3, 4}, function(acc, v) return acc + v end, 0)
assert(sum == 10, "reduce should fold left-to-right from init")
local concatIdx = c.reduce({"a", "b"}, function(acc, v, i) return acc .. i .. v end, "")
assert(concatIdx == "1a2b", "reduce callback should see the index too")

-- find: returns value AND index; nil, 0 when nothing matches.
local v, i = c.find({10, 20, 30}, function(x) return x > 15 end)
assert(v == 20 and i == 2, "find should return the first match and its 1-based index")
local v2, i2 = c.find({1, 2, 3}, function(x) return x > 100 end)
assert(v2 == nil and i2 == 0, "find must return nil, 0 -- not raise -- when nothing matches (F1: legitimate negative answer)")

-- any / all short-circuit correctly, including the vacuous-truth empty case.
assert(c.any({1, 2, 3}, function(x) return x == 2 end) == true, "any should be true when one element matches")
assert(c.any({1, 2, 3}, function(x) return x == 99 end) == false, "any should be false when none match")
assert(c.all({2, 4, 6}, function(x) return x % 2 == 0 end) == true, "all should be true when every element matches")
assert(c.all({2, 3, 4}, function(x) return x % 2 == 0 end) == false, "all should be false when one element fails")
assert(c.any({}, function() return true end) == false, "any of an empty table is false")
assert(c.all({}, function() return false end) == true, "all of an empty table is vacuously true")

-- Empty-table results must be empty TABLES, never nil (nil-slice rule).
local emptyMap = c.map({}, function(v) return v end)
assert(type(emptyMap) == "table" and #emptyMap == 0, "map({}) must return an empty table, not nil")
local emptyFilter = c.filter({}, function() return true end)
assert(type(emptyFilter) == "table" and #emptyFilter == 0, "filter({}) must return an empty table, not nil")

return true
