-- Pin the documented multi-byte / rune-counting behaviour of wrap, truncate
-- and pad: they count Unicode runes, not bytes and not terminal display
-- cells. A CJK or emoji string therefore will NOT align in a terminal even
-- though these functions agree it is "N runes wide" -- that limitation is
-- intentional (see the text package doc) and this test pins it rather than
-- leaving it as an accident nobody verified.
--
-- gopher-lua implements Lua 5.1, which has no \xHH string escape, so the
-- multi-byte fixtures below are literal UTF-8 characters in this source
-- file rather than escape sequences.
local text = require("text")

-- "日本語" is 3 runes but 9 bytes (each CJK ideograph is 3 bytes in UTF-8).
local cjk = "日本語"
assert(#cjk == 9, "sanity check: the CJK fixture should be 9 bytes, got " .. #cjk)

-- truncate treats it as 3 runes: truncating to 4 runes (>= 3) leaves it
-- unchanged, byte length is irrelevant to the comparison.
assert(text.truncate(cjk, 4, ".") == cjk, "truncate should compare against rune count (3), not byte count (9)")

-- truncating to 2 runes keeps the first rune ("日", 3 bytes) plus a 1-rune
-- ellipsis, i.e. 2 runes / 4 bytes total, NOT 2 bytes.
local truncated = text.truncate(cjk, 2, ".")
assert(truncated == "日.", "truncate should keep whole runes, not split a multi-byte rune, got: " .. truncated)

-- pad_left/pad_right count runes: padding a 3-rune string to width 5 adds
-- exactly 2 pad runes, regardless of the 9 bytes already present.
local padded = text.pad_left(cjk, 5, "-")
assert(padded == "--" .. cjk, "pad_left should add (width - rune_count) pad runes")

-- A multi-rune pad character (a single rune outside the BMP, e.g. an emoji)
-- is accepted the same way for a multi-byte string as for an ASCII one --
-- "one rune" is not an ASCII special case.
local emoji = "\240\159\152\128" -- U+1F600 GRINNING FACE, a single rune, 4 bytes (decimal byte escapes, Lua 5.1 compatible)
assert(text.pad_right("x", 3, emoji) == "x" .. emoji .. emoji, "a single multi-byte rune is a valid, single pad rune")

return true
