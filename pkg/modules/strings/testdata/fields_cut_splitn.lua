-- Test fields, cut and split_n.
local strings = require("strings")

-- fields: splits around runs of whitespace, no empty entries.
local f = strings.fields("  the   quick brown\tfox  ")
assert(#f == 4, "fields should produce 4 words, got " .. tostring(#f))
assert(f[1] == "the" and f[2] == "quick" and f[3] == "brown" and f[4] == "fox", "fields should split on whitespace runs")
assert(#strings.fields("") == 0, "fields on empty string should produce an empty table")
assert(#strings.fields("   ") == 0, "fields on an all-whitespace string should produce an empty table")

-- cut: before, after, found.
local before, after, found = strings.cut("app=nginx", "=")
assert(before == "app" and after == "nginx" and found == true, "cut should split on the separator and report found")

before, after, found = strings.cut("no-separator-here", "=")
assert(before == "no-separator-here" and after == "" and found == false, "cut should return the whole string and found=false when sep is absent")

before, after, found = strings.cut("a=b=c", "=")
assert(before == "a" and after == "b=c" and found == true, "cut should only split on the FIRST occurrence")

-- split_n: split with a limit.
local parts = strings.split_n("a,b,c,d", ",", 2)
assert(#parts == 2, "split_n with n=2 should produce 2 elements")
assert(parts[1] == "a" and parts[2] == "b,c,d", "split_n should leave the remainder unsplit in the last element")

parts = strings.split_n("a,b,c", ",", 0)
assert(#parts == 0, "split_n with n=0 should produce an empty table")

parts = strings.split_n("a,b,c", ",", -1)
assert(#parts == 3, "split_n with n=-1 should split all occurrences, like split")

return true
