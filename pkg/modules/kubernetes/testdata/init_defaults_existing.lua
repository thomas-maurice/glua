-- Test: init_defaults with existing labels and annotations
--
-- This test verifies that init_defaults() preserves existing
-- labels and annotations when they are already present,
-- and still allows adding new entries.

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

-- Initialize defaults (should not overwrite existing)
k8s.init_defaults(obj)

-- Existing values should be preserved
assert(obj.metadata.labels.existing == "label", "existing label should be preserved")
assert(obj.metadata.annotations.existing == "annotation", "existing annotation should be preserved")

-- Should still be able to add new entries
obj.metadata.labels.new = "value"
assert(obj.metadata.labels.new == "value", "should be able to add new label")

return true
