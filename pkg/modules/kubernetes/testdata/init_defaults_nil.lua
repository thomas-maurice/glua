-- Test: init_defaults with nil labels and annotations
--
-- This test verifies that init_defaults() properly initializes
-- metadata.labels and metadata.annotations as empty tables when
-- they are nil, and allows adding new entries afterwards.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test-pod"
	}
}

-- Before init_defaults, labels and annotations are nil
assert(obj.metadata.labels == nil, "labels should be nil initially")
assert(obj.metadata.annotations == nil, "annotations should be nil initially")

-- Initialize defaults
k8s.init_defaults(obj)

-- After init_defaults, they should be empty tables
assert(type(obj.metadata.labels) == "table", "labels should be table")
assert(type(obj.metadata.annotations) == "table", "annotations should be table")

-- Should be able to add entries
obj.metadata.labels.app = "myapp"
obj.metadata.annotations["version"] = "1.0"

assert(obj.metadata.labels.app == "myapp", "should be able to add label")
assert(obj.metadata.annotations["version"] == "1.0", "should be able to add annotation")

return true
