-- Test text.wrap: greedy word wrap, hard newlines preserved, long words
-- overflow rather than being split, and width < 1 raises.
local text = require("text")

-- Basic greedy wrap.
local wrapped = text.wrap("the quick brown fox jumps", 10)

-- Rejoining wrapped lines with spaces must reproduce the original words in
-- order: wrap must not drop or reorder anything, only insert line breaks.
local rebuilt = string.gsub(wrapped, "\n", " ")
assert(rebuilt == "the quick brown fox jumps", "wrap must not drop or reorder words, got: " .. wrapped)

for line in string.gmatch(wrapped, "[^\n]+") do
  assert(#line <= 10, "each wrapped line should respect the width, got line of length " .. #line .. ": " .. line)
end

-- A single word longer than width is NOT split; it overflows its own line.
local long = text.wrap("supercalifragilisticexpialidocious", 5)
assert(long == "supercalifragilisticexpialidocious", "a word longer than width must not be broken")

local mixed = text.wrap("a supercalifragilisticexpialidocious b", 5)
local mixedLines = {}
for line in string.gmatch(mixed, "[^\n]+") do
  table.insert(mixedLines, line)
end
assert(mixedLines[2] == "supercalifragilisticexpialidocious", "the overflow word should occupy its own line unsplit, got: " .. mixed)

-- Existing "\n" are hard breaks: paragraphs wrap independently and a blank
-- line survives.
local para = text.wrap("first paragraph here\n\nsecond one", 8)
local paraLines = {}
for line in string.gmatch(para .. "\n", "(.-)\n") do
  table.insert(paraLines, line)
end
local blankFound = false
for _, l in ipairs(paraLines) do
  if l == "" then blankFound = true end
end
assert(blankFound, "a blank line in the input must survive as a blank line in the output, got: " .. para)

-- Empty string wraps to empty string.
assert(text.wrap("", 10) == "", "wrap of an empty string should be empty")

-- width < 1 raises.
local ok, err = pcall(text.wrap, "hello", 0)
assert(ok == false, "wrap with width 0 should raise")
assert(string.find(err, "width must be >= 1") ~= nil, "wrap error should explain the constraint, got: " .. tostring(err))

ok, err = pcall(text.wrap, "hello", -5)
assert(ok == false, "wrap with a negative width should raise")

return true
