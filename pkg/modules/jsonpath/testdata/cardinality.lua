-- query must always return a table -- empty, one element, or many -- never
-- nil (the nil-slice rule); first must fall back to its required default;
-- exists must answer true/false without raising on a miss.
local jsonpath = require("jsonpath")

local data = {
  spec = {
    containers = {
      {image = "nginx"},
      {image = "redis"},
    },
  },
}

-- Zero matches: query returns an empty table, not nil.
local none = jsonpath.query(data, ".spec.missing")
assert(type(none) == "table", "query with zero matches must return a table, not nil")
assert(#none == 0, "query with zero matches must return an empty table")

-- One match: query still returns a table (not the bare value).
local one = jsonpath.query(data, ".spec.containers[0].image")
assert(type(one) == "table" and #one == 1 and one[1] == "nginx",
  "query with one match must still return a 1-element table, not the bare value")

-- Many matches: query returns every match in order.
local many = jsonpath.query(data, ".spec.containers[*].image")
assert(#many == 2 and many[1] == "nginx" and many[2] == "redis",
  "query with many matches must return every match in order")

-- first falls back to the required default when there is no match.
assert(jsonpath.first(data, ".spec.missing", "fallback") == "fallback",
  "first must return default when path has no match")
assert(jsonpath.first(data, ".spec.containers[0].image", "fallback") == "nginx",
  "first must return the first match when one exists")

-- exists answers true/false, never raises on a miss (F1: legitimate negative answer).
assert(jsonpath.exists(data, ".spec.containers[0].image") == true, "exists should be true for a present field")
assert(jsonpath.exists(data, ".spec.missing") == false, "exists should be false for an absent field, not raise")

return true
