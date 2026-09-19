-- Test keys, values, merge, pick and omit -- all map-mode (array part and
-- hash part together), order-unspecified where noted.
local c = require("collections")

-- keys/values on a pure array: numeric indices are still keys.
local akeys = c.keys({10, 20, 30})
table.sort(akeys) -- order is unspecified, so sort before comparing
assert(#akeys == 3 and akeys[1] == 1 and akeys[2] == 2 and akeys[3] == 3, "keys of an array should be its 1..n indices")
local avalues = c.values({10, 20, 30})
table.sort(avalues)
assert(#avalues == 3 and avalues[1] == 10 and avalues[2] == 20 and avalues[3] == 30, "values of an array should be its elements")

-- keys/values on a mixed table cover both the array and hash parts.
local mixed = {1, 2, name = "x", active = true}
local mkeys = c.keys(mixed)
assert(#mkeys == 4, "keys must cover both the array part and the hash part")

local emptyKeys = c.keys({})
assert(type(emptyKeys) == "table" and #emptyKeys == 0, "keys({}) must return an empty table, not nil")
local emptyValues = c.values({})
assert(type(emptyValues) == "table" and #emptyValues == 0, "values({}) must return an empty table, not nil")

-- merge: shallow, left-to-right, later wins, new table, inputs untouched.
local base = {a = 1, b = 2}
local override = {b = 20, c = 30}
local merged = c.merge(base, override)
assert(merged.a == 1 and merged.b == 20 and merged.c == 30, "merge should apply later tables over earlier ones")
assert(base.b == 2, "merge must not mutate its inputs")

local mergedThree = c.merge({a = 1}, {a = 2}, {a = 3})
assert(mergedThree.a == 3, "merge should apply every extra table left-to-right, last one winning")

local mergedNone = c.merge({a = 1})
assert(mergedNone.a == 1, "merge with no extra tables should copy the base table")

-- merge on array-shaped tables also covers the array part (map mode).
local mergedArrays = c.merge({1, 2, 3}, {10, 20})
assert(mergedArrays[1] == 10 and mergedArrays[2] == 20 and mergedArrays[3] == 3, "merge should overlay array indices too")

-- pick: only the named keys survive; a missing name is silently skipped.
local picked = c.pick({a = 1, b = 2, c = 3}, {"a", "c", "missing"})
assert(picked.a == 1 and picked.c == 3 and picked.b == nil, "pick should keep only the named keys")
assert(picked.missing == nil, "pick should silently skip a name absent from the source table")

local emptyPick = c.pick({a = 1}, {})
assert(type(emptyPick) == "table", "pick with an empty names list must return an empty table, not nil")
for _ in pairs(emptyPick) do error("pick with no names should produce no keys") end

-- omit: every key except the named ones; non-string keys are never omitted.
local omitted = c.omit({a = 1, b = 2, c = 3}, {"b"})
assert(omitted.a == 1 and omitted.c == 3 and omitted.b == nil, "omit should drop only the named keys")

local omittedArray = c.omit({1, 2, 3}, {"b"})
assert(#omittedArray == 3, "omit must never drop numeric (non-string) keys, since names only match string keys")

return true
