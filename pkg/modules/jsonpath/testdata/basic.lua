-- Test every supported syntax form against a fixture table, and the three
-- equivalent path spellings (bare, $-rooted, braced).
local jsonpath = require("jsonpath")

local data = {
  spec = {
    replicas = 3,
    containers = {
      {name = "web", image = "nginx"},
      {name = "cache", image = "redis"},
      {name = "sidecar", image = "envoy"},
    },
  },
}

-- Field access: bare, $-rooted and braced forms are equivalent.
assert(jsonpath.first(data, ".spec.replicas", nil) == 3, "bare path should resolve a field")
assert(jsonpath.first(data, "$.spec.replicas", nil) == 3, "$-rooted path should resolve a field")
assert(jsonpath.first(data, "{.spec.replicas}", nil) == 3, "braced path should resolve a field")

-- Index: a single bracketed integer selects one element.
assert(jsonpath.first(data, ".spec.containers[0].name", nil) == "web", "index [0] should select the first element")
assert(jsonpath.first(data, ".spec.containers[2].name", nil) == "sidecar", "index [2] should select the third element")

-- Slice: [0:2] selects elements 0 and 1.
local sliced = jsonpath.query(data, ".spec.containers[0:2].name")
assert(#sliced == 2 and sliced[1] == "web" and sliced[2] == "cache", "slice [0:2] should select the first two elements")

-- Wildcard: [*] selects every element.
local images = jsonpath.query(data, "{.spec.containers[*].image}")
assert(#images == 3, "wildcard should select every element")
assert(images[1] == "nginx" and images[2] == "redis" and images[3] == "envoy", "wildcard should preserve order")

-- Recursive descent: ..name finds every "name" key at any depth.
local names = jsonpath.query(data, "..name")
assert(#names == 3, "recursive descent should find every matching key at any depth")

-- Filter: ?(@.name=="cache") selects elements matching the comparison.
local filtered = jsonpath.query(data, "{.spec.containers[?(@.name==\"cache\")].image}")
assert(#filtered == 1 and filtered[1] == "redis", "filter should select only the matching element")

return true
