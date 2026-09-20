-- Test trim_space, trim_prefix and trim_suffix.
--
-- trim_space is the single most-missed function in this module: before it
-- existed, trimming whitespace required strings.trim(s, " \t\n\r") which
-- nobody remembers to spell out correctly.
local strings = require("strings")

assert(strings.trim_space("   hello   ") == "hello", "trim_space should remove surrounding spaces")
assert(strings.trim_space("\t\nhello\r\n") == "hello", "trim_space should remove tabs/newlines/CR")
assert(strings.trim_space("hello") == "hello", "trim_space should be a no-op without whitespace")
assert(strings.trim_space("   ") == "", "trim_space should reduce an all-whitespace string to empty")
assert(strings.trim_space("") == "", "trim_space should handle empty string")

assert(strings.trim_prefix("hello world", "hello ") == "world", "trim_prefix should remove matching prefix")
assert(strings.trim_prefix("hello world", "bye ") == "hello world", "trim_prefix should leave s unchanged if prefix does not match")
assert(strings.trim_prefix("hello", "") == "hello", "trim_prefix with empty prefix is a no-op")

assert(strings.trim_suffix("app.tar.gz", ".gz") == "app.tar", "trim_suffix should remove matching suffix")
assert(strings.trim_suffix("app.tar.gz", ".zip") == "app.tar.gz", "trim_suffix should leave s unchanged if suffix does not match")
assert(strings.trim_suffix("hello", "") == "hello", "trim_suffix with empty suffix is a no-op")

return true
