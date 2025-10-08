
-- Test: init_defaults full workflow
--
-- This test demonstrates a complete workflow of using init_defaults()
-- with a pod object, verifying initial state, initialization,
-- and adding labels/annotations.

local k8s = require("kubernetes")

-- Test case 1: nil labels and annotations
local pod = {
	metadata = {
		name = "test-pod",
		namespace = "default"
	}
}

-- Before init_defaults, labels and annotations are nil
assert(pod.metadata.labels == nil, "labels should be nil initially")
assert(pod.metadata.annotations == nil, "annotations should be nil initially")

k8s.init_defaults(pod)

-- After init_defaults, they should be tables
assert(type(pod.metadata.labels) == "table", "labels should be table")
assert(type(pod.metadata.annotations) == "table", "annotations should be table")

-- Add some labels and annotations
pod.metadata.labels.app = "myapp"
pod.metadata.labels.tier = "backend"
pod.metadata.annotations.version = "1.0.0"

-- Verify they were added successfully
assert(pod.metadata.labels.app == "myapp", "label app should be set")
assert(pod.metadata.labels.tier == "backend", "label tier should be set")
assert(pod.metadata.annotations.version == "1.0.0", "annotation version should be set")

return true
