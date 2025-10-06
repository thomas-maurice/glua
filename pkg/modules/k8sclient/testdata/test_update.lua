-- Test: update operation
--
-- Tests updating an existing Kubernetes resource.

-- Create initial ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "update-config",
		namespace = "default"
	},
	data = {
		original = "value"
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

-- Update it
created.data.updated = "new-value"
created.data.original = "changed"

local updated, err = client.update(created)
if err then
	error("Failed to update: " .. err)
end

if updated.data.updated ~= "new-value" then
	error("Expected data.updated 'new-value', got " .. tostring(updated.data.updated))
end

if updated.data.original ~= "changed" then
	error("Expected data.original 'changed', got " .. tostring(updated.data.original))
end

return true
