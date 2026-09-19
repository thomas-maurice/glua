-- Test uniq, flatten, reverse, zip and chunk.
local c = require("collections")

-- uniq: primitives dedup by value, tables dedup by identity (not structure).
local dedupPrimitives = c.uniq({1, 2, 2, 3, 1, 3, 3})
assert(#dedupPrimitives == 3, "uniq should dedup primitives by value")
assert(dedupPrimitives[1] == 1 and dedupPrimitives[2] == 2 and dedupPrimitives[3] == 3, "uniq should keep first occurrence order")

local same = {}
local other = {} -- structurally identical to `same` but a different table
local dedupTables = c.uniq({same, same, other})
assert(#dedupTables == 2, "uniq must dedup tables by identity: same,same collapses but other (a distinct, structurally-equal table) survives")

local emptyUniq = c.uniq({})
assert(type(emptyUniq) == "table" and #emptyUniq == 0, "uniq({}) must return an empty table, not nil")

-- flatten: depth-limited and fully (-1).
local nested = {1, {2, 3, {4, 5}}, 6}
local flatOne = c.flatten(nested, 1)
assert(#flatOne == 5, "flatten(nested, 1) should inline one level: {1,2,3,{4,5},6}")
assert(type(flatOne[4]) == "table", "flatten(nested, 1) should leave the doubly-nested {4,5} intact at depth 1")

local flatFull = c.flatten(nested, -1)
assert(#flatFull == 6, "flatten(nested, -1) should inline every level")
for idx = 1, 6 do
  assert(flatFull[idx] == idx, "flatten(nested, -1) should produce 1..6 in order")
end

local flatZero = c.flatten(nested, 0)
assert(#flatZero == 3, "flatten(t, 0) should not descend at all, returning t's own elements unchanged")

local emptyFlatten = c.flatten({}, -1)
assert(type(emptyFlatten) == "table" and #emptyFlatten == 0, "flatten({}, -1) must return an empty table, not nil")

local okDepth, errDepth = pcall(c.flatten, {1}, -2)
assert(okDepth == false, "flatten must raise for a depth less than -1")
assert(string.find(errDepth, "collections.flatten") ~= nil, "flatten error should be prefixed, got: " .. tostring(errDepth))

-- reverse.
local reversed = c.reverse({1, 2, 3})
assert(reversed[1] == 3 and reversed[2] == 2 and reversed[3] == 1, "reverse should reverse element order")
local emptyReverse = c.reverse({})
assert(type(emptyReverse) == "table" and #emptyReverse == 0, "reverse({}) must return an empty table, not nil")

-- zip: length is min(#a, #b), extra elements dropped.
local zipped = c.zip({1, 2, 3}, {"a", "b"})
assert(#zipped == 2, "zip should truncate to the shorter table's length")
assert(zipped[1][1] == 1 and zipped[1][2] == "a", "zip should pair elements positionally")
assert(zipped[2][1] == 2 and zipped[2][2] == "b", "zip should pair the second elements")

local emptyZip = c.zip({}, {1, 2})
assert(type(emptyZip) == "table" and #emptyZip == 0, "zip with an empty input must return an empty table, not nil")

-- chunk: fixed-size groups, last chunk may be short; size < 1 raises.
local chunks = c.chunk({1, 2, 3, 4, 5}, 2)
assert(#chunks == 3, "chunk of 5 elements by 2 should produce 3 chunks")
assert(#chunks[1] == 2 and #chunks[2] == 2 and #chunks[3] == 1, "chunk's last group may be shorter than size")
assert(chunks[3][1] == 5, "chunk should preserve element order within each chunk")

local emptyChunk = c.chunk({}, 3)
assert(type(emptyChunk) == "table" and #emptyChunk == 0, "chunk({}, n) must return an empty table, not nil")

local okSize, errSize = pcall(c.chunk, {1, 2}, 0)
assert(okSize == false, "chunk must raise when size < 1")
assert(string.find(errSize, "collections.chunk") ~= nil, "chunk error should be prefixed, got: " .. tostring(errSize))

return true
