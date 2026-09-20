-- Test deep_equal and deep_copy: structural comparison/copying, that
-- metatables are not considered, and that a self-referencing (cyclic) table
-- is handled rather than hanging the script.
local c = require("collections")

-- deep_equal: structurally equal but not the same table.
local a = {x = 1, y = {2, 3}}
local b = {x = 1, y = {2, 3}}
assert(c.deep_equal(a, b) == true, "deep_equal should compare structurally, not by identity")
assert(c.deep_equal(a, a) == true, "deep_equal should be reflexive")

local different = {x = 1, y = {2, 4}}
assert(c.deep_equal(a, different) == false, "deep_equal should detect a nested difference")

local extraKey = {x = 1, y = {2, 3}, z = 9}
assert(c.deep_equal(a, extraKey) == false, "deep_equal should detect an extra key")

assert(c.deep_equal(1, 1) == true, "deep_equal should work on plain primitives")
assert(c.deep_equal(1, "1") == false, "deep_equal must not consider a number and its string form equal")
assert(c.deep_equal(false, nil) == false, "deep_equal must distinguish false from nil")

-- deep_copy: a structural copy, not an alias -- mutating the copy must not
-- affect the original, including through nested tables.
local original = {x = 1, nested = {y = 2}}
local copy = c.deep_copy(original)
assert(c.deep_equal(original, copy) == true, "deep_copy's result should deep_equal the original")
copy.nested.y = 999
assert(original.nested.y == 2, "mutating the copy's nested table must not affect the original")

-- Cyclic table: must not hang, and the cycle must be preserved in the copy
-- (not expanded into an infinite/oversized structure).
local cyclic = {}
cyclic.self = cyclic
assert(c.deep_equal(cyclic, cyclic) == true, "deep_equal on a self-referencing table must terminate and report equal")

local cyclicCopy = c.deep_copy(cyclic)
assert(cyclicCopy ~= cyclic, "deep_copy must produce a new table, not alias the original")
assert(cyclicCopy.self == cyclicCopy, "deep_copy must preserve the cycle: the copy's self-reference should point back to the copy itself")

-- Two independently-built cyclic-but-equal structures still compare equal.
local cyclicB = {}
cyclicB.self = cyclicB
assert(c.deep_equal(cyclic, cyclicB) == true, "two structurally-equal cyclic tables should deep_equal")

return true
