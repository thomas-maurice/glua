-- Test: delete operation
--
-- Tests deleting a Kubernetes resource.

-- Create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "delete-config",
		namespace = "default"
	},
	data = {
		test = "data"
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

-- Delete it
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local err = client.delete(gvk, "default", "delete-config")

if err then
	error("Failed to delete: " .. err)
end

-- Try to get it (should fail)
local fetched, err = client.get(gvk, "default", "delete-config")

if not err then
	error("Expected error getting deleted resource, got nil")
end

return true
