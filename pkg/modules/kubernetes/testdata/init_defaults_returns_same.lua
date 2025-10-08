
-- Test: init_defaults returns the same object
--
-- This test verifies that init_defaults() modifies the object
-- in-place and returns the same object reference.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test"
	}
}

local result = k8s.init_defaults(obj)

-- Should return the same object (modified in-place)
assert(result == obj, "should return same object")
assert(result.metadata.name == "test", "object should be modified in-place")

return true
