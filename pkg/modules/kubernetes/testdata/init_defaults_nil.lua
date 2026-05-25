
-- Test: init_defaults with nil labels and annotations
--
-- This test verifies that init_defaults() properly initializes
-- metadata.labels and metadata.annotations as empty tables when
-- they are nil, and allows adding new entries afterwards.
-- Note: init_defaults returns the updated object; the caller must assign.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test-pod"
	}
}

-- Before init_defaults, labels and annotations are nil
assert(obj.metadata.labels == nil, "labels should be nil initially")
assert(obj.metadata.annotations == nil, "annotations should be nil initially")

-- Initialize defaults (must assign return value)
obj = k8s.init_defaults(obj)

-- After init_defaults, they should be empty tables
assert(type(obj.metadata.labels) == "table", "labels should be table")
assert(type(obj.metadata.annotations) == "table", "annotations should be table")

-- Should be able to add entries via add_label / add_annotation
obj = k8s.add_label(obj, "app", "myapp")
obj = k8s.add_annotation(obj, "version", "1.0")

assert(k8s.get_label(obj, "app") == "myapp", "should be able to add label")
assert(k8s.get_annotation(obj, "version") == "1.0", "should be able to add annotation")

return true
