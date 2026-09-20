-- Test title: first-rune title casing, not language-aware. Pins the
-- ASCII behaviour and a non-ASCII first rune so the limitation is a
-- documented, tested property rather than an accident.
local strings = require("strings")

assert(strings.title("hello world") == "Hello World", "title should capitalise the first letter of each word")
assert(strings.title("hello") == "Hello", "title should work on a single word")
assert(strings.title("") == "", "title should handle empty string")

-- Leading whitespace is preserved verbatim, not collapsed.
assert(strings.title("  hello world") == "  Hello World", "title should preserve leading whitespace")
assert(strings.title("hello   world") == "Hello   World", "title should preserve internal whitespace width")

-- Already-titled or all-caps input is left alone rune-by-rune except the
-- word-start rune (which is idempotent under ToTitle).
assert(strings.title("HELLO WORLD") == "HELLO WORLD", "title should not lowercase the rest of an already-uppercase word")

-- Non-ASCII first rune: 'é' (U+00E9) title-cases to 'É' (U+00C9). This pins
-- that title() operates on Unicode runes, not just ASCII bytes.
assert(strings.title("école") == "École", "title should title-case a non-ASCII first rune")

-- Digraph case: 'ǳ' (U+01F3, lowercase dz) has a TITLECASE form 'ǲ'
-- (U+01F2) that is distinct from its UPPERCASE form 'Ǳ' (U+01F1). Using
-- unicode.ToTitle (not ToUpper) is what makes this assertion hold, and is
-- the whole reason ToTitle was chosen for this implementation.
assert(strings.title("ǳelo") == "ǲelo", "title should use Unicode titlecase, not uppercase, for digraphs")

return true
