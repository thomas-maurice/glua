-- Test equal_fold: Unicode case-insensitive comparison.
local strings = require("strings")

assert(strings.equal_fold("Hello", "hello") == true, "equal_fold should ignore ASCII case")
assert(strings.equal_fold("HELLO", "hello") == true, "equal_fold should ignore ASCII case fully")
assert(strings.equal_fold("hello", "world") == false, "equal_fold should return false for different strings")
assert(strings.equal_fold("", "") == true, "equal_fold on two empty strings should be true")
assert(strings.equal_fold("Straße", "STRASSE") == false, "equal_fold performs simple case-folding, not full Unicode special-casing (ß != SS)")

return true
