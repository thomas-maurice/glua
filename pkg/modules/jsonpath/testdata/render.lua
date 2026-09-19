-- render is the native kubectl-jsonpath text-template mode: it differs from
-- query by producing a single string (concatenating matches and any literal
-- text between {...} segments) and by raising on a missing key instead of
-- silently producing a gap (the query/first/exists family never raises on
-- a miss).
local jsonpath = require("jsonpath")

local data = {
  metadata = {name = "web-1"},
  status = {phase = "Running"},
  spec = {
    containers = {
      {image = "nginx"},
      {image = "redis"},
    },
  },
}

-- Multiple {...} segments with literal text between them are concatenated.
local out = jsonpath.render(data, "{.metadata.name}: {.status.phase}")
assert(out == "web-1: Running", "render should concatenate matches and literal text between segments")

-- A single template with many matches space-joins them (native PrintResults behaviour).
local images = jsonpath.render(data, "{.spec.containers[*].image}")
assert(images == "nginx redis", "render should space-join multiple matches from one segment")

-- Bare and $-rooted forms are auto-braced, same as query/first/exists.
assert(jsonpath.render(data, ".metadata.name") == "web-1", "render should accept a bare path")
assert(jsonpath.render(data, "$.metadata.name") == "web-1", "render should accept a $-rooted path")

-- render raises on a missing key -- the asymmetry with query/exists is deliberate:
-- a text template that silently renders a gap is a template bug.
local ok, err = pcall(jsonpath.render, data, "{.metadata.missing}")
assert(not ok, "render must raise on a missing key")
assert(type(err) == "string", "render's error must be a catchable string")

-- ...while query and exists do NOT raise on the very same missing key.
assert(jsonpath.exists(data, ".metadata.missing") == false, "exists must not raise on the key that makes render raise")
local q = jsonpath.query(data, ".metadata.missing")
assert(type(q) == "table" and #q == 0, "query must not raise on the key that makes render raise")

return true
