
-- Test: init_defaults with existing labels and annotations
--
-- This test verifies that init_defaults() preserves existing
-- labels and annotations when they are already present,
-- and still allows adding new entries.
-- Note: init_defaults returns the updated object; the caller must assign.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test-pod",
		labels = {
			existing = "label"
		},
		annotations = {
			existing = "annotation"
		}
	}
}

-- Initialize defaults (must assign return value)
obj = k8s.init_defaults(obj)

-- Existing values should be preserved
assert(k8s.get_label(obj, "existing") == "label", "existing label should be preserved")
assert(k8s.get_annotation(obj, "existing") == "annotation", "existing annotation should be preserved")

-- Should still be able to add new entries
obj = k8s.add_label(obj, "new", "value")
assert(k8s.get_label(obj, "new") == "value", "should be able to add new label")

return true
