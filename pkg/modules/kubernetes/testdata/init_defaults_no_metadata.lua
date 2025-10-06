-- Test: init_defaults with no metadata field
--
-- This test verifies that init_defaults() creates the metadata
-- field if it doesn't exist, along with labels and annotations.

local k8s = require("kubernetes")

local obj = {}

-- Initialize defaults (should create metadata)
k8s.init_defaults(obj)

-- Should have created metadata with labels and annotations
assert(type(obj.metadata) == "table", "metadata should be created")
assert(type(obj.metadata.labels) == "table", "labels should be table")
assert(type(obj.metadata.annotations) == "table", "annotations should be table")

return true
