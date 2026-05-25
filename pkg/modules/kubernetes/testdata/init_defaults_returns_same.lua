
-- Test: init_defaults returns the updated object
--
-- init_defaults() returns the (possibly new) object with metadata initialized.
-- The caller must assign the return value. The returned object has the same
-- content as the input (it is not a different object conceptually, but table
-- identity is not preserved due to JSON round-trip).

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test"
	}
}

local result = k8s.init_defaults(obj)

-- The returned object should have metadata initialized
assert(result ~= nil, "should return a non-nil object")
assert(type(result.metadata) == "table", "metadata should be a table")
assert(result.metadata.name == "test", "name should be preserved")
assert(type(result.metadata.labels) == "table", "labels should be initialized")
assert(type(result.metadata.annotations) == "table", "annotations should be initialized")

return true
