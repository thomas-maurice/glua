-- Test split and join functions
local strings = require("strings")

-- Test split
local text = "hello,world,foo,bar"
local parts = strings.split(text, ",")
assert(type(parts) == "table", "split should return a table")
assert(#parts == 4, "split should return 4 parts")
assert(parts[1] == "hello", "first part should be 'hello'")
assert(parts[2] == "world", "second part should be 'world'")
assert(parts[3] == "foo", "third part should be 'foo'")
assert(parts[4] == "bar", "fourth part should be 'bar'")

-- Test split with multi-character separator
local text2 = "hello::world::foo"
local parts2 = strings.split(text2, "::")
assert(#parts2 == 3, "split should handle multi-char separator")
assert(parts2[1] == "hello", "first part should be 'hello'")
assert(parts2[2] == "world", "second part should be 'world'")
assert(parts2[3] == "foo", "third part should be 'foo'")

-- Test split with empty string
local parts3 = strings.split("", ",")
assert(#parts3 == 1, "split of empty string should return table with one empty string")
assert(parts3[1] == "", "split of empty string should contain empty string")

-- Test join
local joined = strings.join(parts, ",")
assert(joined == text, "join should reconstruct the original string")

-- Test join with different separator
local joined2 = strings.join(parts, " - ")
assert(joined2 == "hello - world - foo - bar", "join should use provided separator")

-- Test join with empty table
local empty_table = {}
local joined3 = strings.join(empty_table, ",")
assert(joined3 == "", "join of empty table should return empty string")

-- Test join with single element
local single = {"hello"}
local joined4 = strings.join(single, ",")
assert(joined4 == "hello", "join with single element should return that element")

return true
