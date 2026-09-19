-- Test text.indent and text.dedent.
local text = require("text")

-- indent prefixes every line.
assert(text.indent("a\nb\nc", "> ") == "> a\n> b\n> c", "indent should prefix every line")
assert(text.indent("solo", "  | ") == "  | solo", "indent should prefix a single line")

-- A trailing empty line (s ends with "\n") is NOT prefixed.
assert(text.indent("a\nb\n", "> ") == "> a\n> b\n", "indent should not prefix the trailing empty line from a final newline")

-- An internal blank line IS prefixed (only the trailing one is special-cased).
assert(text.indent("a\n\nb", "> ") == "> a\n> \n> b", "indent should prefix internal blank lines")

assert(text.indent("", "> ") == "", "indent of an empty string should stay empty")

-- dedent removes the longest common leading whitespace.
local block = "    def f():\n        return 1\n"
assert(text.dedent(block) == "def f():\n    return 1\n", "dedent should remove the common 4-space margin")

-- Mixed indentation: the common prefix is only as long as the shortest one.
local mixed = "    a\n      b\n    c"
assert(text.dedent(mixed) == "a\n  b\nc", "dedent should use the shortest common indentation as the margin")

-- Blank lines do not participate in computing the margin, and are
-- normalized to empty in the output even if they had different whitespace.
local withBlank = "    a\n\t\n    b"
assert(text.dedent(withBlank) == "a\n\nb", "dedent should ignore blank/whitespace-only lines when computing the margin, and normalize them to empty")

-- Tabs and spaces are literal, not expanded: a tab-indented line and a
-- space-indented line share no common prefix, so nothing is removed.
local tabsAndSpaces = "\ta\n    b"
assert(text.dedent(tabsAndSpaces) == tabsAndSpaces, "dedent should not remove anything when tab- and space-indented lines share no literal common prefix")

return true
