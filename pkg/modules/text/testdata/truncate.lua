-- Test text.truncate.
local text = require("text")

-- Fits already: unchanged.
assert(text.truncate("hello", 10, "...") == "hello", "truncate should not touch a string that already fits")

-- Exactly at width: unchanged (boundary case, no ellipsis appended).
assert(text.truncate("hello", 5, "...") == "hello", "truncate at exactly width should leave the string unchanged")

-- Shortened: prefix + ellipsis, total length == width.
local out = text.truncate("hello world", 8, "...")
assert(out == "hello...", "truncate should keep (width - len(ellipsis)) runes then append ellipsis, got: " .. out)
assert(#out == 8, "truncated result should be exactly width runes long")

-- Single-character ellipsis works the same way. gopher-lua is Lua 5.1, which
-- has no \xHH string escape, so the "…" character is written literally
-- (this file is UTF-8) rather than as an escape sequence.
out = text.truncate("hello world", 6, "…")
assert(out == "hello…", "truncate should support a single-rune ellipsis, got: " .. out)

-- ellipsis longer than width raises.
local ok, err = pcall(text.truncate, "hello world", 2, "...")
assert(ok == false, "truncate should raise when width is smaller than the ellipsis's rune length")
assert(string.find(err, "must be >=") ~= nil, "truncate error should explain the constraint, got: " .. tostring(err))

-- width exactly equal to the ellipsis length is allowed (result is just the ellipsis).
out = text.truncate("hello world", 3, "...")
assert(out == "...", "truncate with width == len(ellipsis) should return just the ellipsis")

return true
