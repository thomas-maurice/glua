---@diagnostic disable: undefined-global, lowercase-global
-- mutate_node.lua: Adds "hello: ok" label to all nodes
--
-- Global variables available:
--   node: the Kubernetes Node object as a Lua table
--   patches: empty table to populate with JSON patch operations
--   kubernetes: the kubernetes module with helper functions

local k8s = require("kubernetes")

-- Initialize metadata if needed
node = k8s.init_defaults(node)

-- Add the "hello: ok" label to the node
node = k8s.add_label(node, "hello", "ok")

-- Generate the JSON patch for the label
table.insert(patches, {
	op = "add",
	path = "/metadata/labels/hello",
	value = "ok"
})

print(string.format("Added hello=ok label to node %s", node.metadata.name or "unknown"))
