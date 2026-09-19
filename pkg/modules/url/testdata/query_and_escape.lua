-- parse_query / build_query / escape / unescape / resolve.
local url = require("url")

-- parse_query: repeated keys, bare key, and the empty-string case.
local q = url.parse_query("a=1&a=2&flag")
assert(type(q) == "table", "parse_query returns a table")
assert(#q.a == 2, "repeated key keeps both values")
assert(q.a[1] == "1" and q.a[2] == "2", "repeated key values in order")
assert(#q.flag == 1 and q.flag[1] == "", "bare key maps to a one-element array of empty string")

local empty = url.parse_query("")
assert(type(empty) == "table", "parse_query never returns nil for an empty string")
assert(next(empty) == nil, "parse_query(\"\") is an empty table")

-- build_query: multi-value keys as an array, sorted output, determinism.
local built = url.build_query({ns = "default", label = {"a", "b"}})
assert(built == "label=a&label=b&ns=default", "build_query sorts by key: got " .. built)
assert(url.build_query({ns = "default", label = {"a", "b"}}) == built, "build_query is deterministic")

-- query_escape / query_unescape round-trip, space as '+'.
local escaped = url.query_escape("a b/c")
assert(escaped == "a+b%2Fc", "query_escape uses + for space: got " .. escaped)
assert(url.query_unescape(escaped) == "a b/c", "query_unescape inverts query_escape")

-- path_escape / path_unescape round-trip, space as '%20'.
local pescaped = url.path_escape("a b/c")
assert(pescaped:find("%%20") ~= nil, "path_escape uses %%20 for space: got " .. pescaped)
assert(url.path_unescape(pescaped) == "a b/c", "path_unescape inverts path_escape")

-- resolve: the trailing-slash rule, explicitly.
assert(url.resolve("https://example.com/a/b", "c") == "https://example.com/a/c", "no trailing slash drops last segment")
assert(url.resolve("https://example.com/a/b/", "c") == "https://example.com/a/b/c", "trailing slash keeps it")
assert(url.resolve("https://example.com/a/b", "../c") == "https://example.com/c", "absolute-ish ref")

return true
