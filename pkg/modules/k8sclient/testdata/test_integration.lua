-- Test: integration - full CRUD workflow
--
-- Tests a complete create, read, update, delete workflow.

-- Create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "integration-config",
		namespace = "default"
	},
	data = {
		initial = "value"
	}
}

local created, err = client.create(cm)
if err then
	error("Failed to create: " .. err)
end

if created.metadata.name ~= "integration-config" then
	error("Created object has wrong name: " .. (created.metadata.name or "nil"))
end

-- Get the ConfigMap
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local fetched, err = client.get(gvk, "default", "integration-config")

if err then
	error("Failed to get: " .. err)
end

if fetched.data.initial ~= "value" then
	error("Fetched object has wrong data: " .. (fetched.data.initial or "nil"))
end

-- Update the ConfigMap
fetched.data.updated = "new-value"
local updated, err = client.update(fetched)

if err then
	error("Failed to update: " .. err)
end

if updated.data.updated ~= "new-value" then
	error("Updated object doesn't have new data")
end

if updated.data.initial ~= "value" then
	error("Updated object lost original data")
end

-- List to verify it exists
local items, err = client.list(gvk, "default")

if err then
	error("Failed to list: " .. err)
end

local found = false
for i, item in ipairs(items) do
	if item.metadata.name == "integration-config" then
		found = true
		break
	end
end

if not found then
	error("ConfigMap not found in list")
end

-- Delete the ConfigMap
local err = client.delete(gvk, "default", "integration-config")

if err then
	error("Failed to delete: " .. err)
end

-- Verify it's gone
local deleted, err = client.get(gvk, "default", "integration-config")

if not err then
	error("Expected error getting deleted resource, got nil")
end

return true
