-- Test get_path: dotted traversal, the numeric-key-then-string-key ambiguity
-- rule, the required default argument, and stopping on a missing key or a
-- non-table.
local c = require("collections")

local obj = {
  spec = {
    containers = {
      {image = "nginx:1"},
      {image = "redis:7"},
    },
  },
}

-- Plain nested lookup.
assert(c.get_path(obj, "spec.containers", nil) ~= nil, "get_path should resolve a simple nested path")

-- A numeric segment indexes the array by position (1-based, matching Lua).
assert(c.get_path(obj, "spec.containers.1.image", "<none>") == "nginx:1", "get_path should index arrays with a numeric segment")
assert(c.get_path(obj, "spec.containers.2.image", "<none>") == "redis:7", "get_path should walk multiple numeric segments")

-- Missing key at any depth returns default, not an error.
assert(c.get_path(obj, "spec.missing.thing", "fallback") == "fallback", "get_path should return default on a missing key")
assert(c.get_path(obj, "spec.containers.99.image", "fallback") == "fallback", "get_path should return default on an out-of-range numeric segment")

-- Indexing through a non-table stops and returns default.
assert(c.get_path(obj, "spec.containers.1.image.nope", "fallback") == "fallback", "get_path must stop and return default when a non-table is indexed")

-- default is REQUIRED (D4): a falsy stored value (false) must be returned as-is,
-- not replaced by default the way `x or default` would incorrectly do.
local withFalse = {enabled = false}
assert(c.get_path(withFalse, "enabled", "fallback") == false, "get_path must return a stored `false` verbatim, not fall back")

-- Passing nil explicitly as default is the documented way to say "no default".
assert(c.get_path(obj, "spec.missing", nil) == nil, "get_path with an explicit nil default should return nil on a miss")

-- The empty string is a single empty-segment path: t[""], normally absent.
assert(c.get_path(obj, "", "fallback") == "fallback", "get_path with an empty path should look up the empty string key and fall back to default")
local withEmptyKey = {}
withEmptyKey[""] = "empty-key-value"
assert(c.get_path(withEmptyKey, "", "fallback") == "empty-key-value", "get_path with an empty path should still find an explicit empty-string key")

-- A numeric-looking segment is tried as a number key first, then a string key.
local ambiguous = {}
ambiguous["1"] = "string-key-value" -- only the string key "1" is set, no array slot 1
assert(c.get_path(ambiguous, "1", "fallback") == "string-key-value", "get_path should fall back to a literal string key when no numeric key exists")

local numericFirst = {}
numericFirst[1] = "number-key-value"
numericFirst["1"] = "string-key-value"
assert(c.get_path(numericFirst, "1", "fallback") == "number-key-value", "get_path should prefer the number key over the string key when both exist")

return true
