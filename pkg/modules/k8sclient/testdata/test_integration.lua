
-- Test: integration - full CRUD workflow
--
-- Tests a complete create, read, update, delete workflow.

-- Create a ConfigMap
local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_INTEGRATION_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_INITIAL_KEY] = TEST_INITIAL_VALUE
	}
}

local created = client:create(cm)

if created.metadata.name ~= TEST_INTEGRATION_NAME then
	error("Created object has wrong name: " .. (created.metadata.name or "nil"))
end

-- Get the ConfigMap
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local fetched = client:get(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)

if fetched.data[TEST_INITIAL_KEY] ~= TEST_INITIAL_VALUE then
	error("Fetched object has wrong data: " .. (fetched.data[TEST_INITIAL_KEY] or "nil"))
end

-- Update the ConfigMap
fetched.data[TEST_UPDATE_KEY] = TEST_UPDATE_VALUE
local updated = client:update(fetched)

if updated.data[TEST_UPDATE_KEY] ~= TEST_UPDATE_VALUE then
	error("Updated object doesn't have new data")
end

if updated.data[TEST_INITIAL_KEY] ~= TEST_INITIAL_VALUE then
	error("Updated object lost original data")
end

-- List to verify it exists
local items = client:list(gvk, TEST_NAMESPACE)

local found = false
for i, item in ipairs(items) do
	if item.metadata.name == TEST_INTEGRATION_NAME then
		found = true
		break
	end
end

if not found then
	error("ConfigMap not found in list")
end

-- Delete the ConfigMap
client:delete(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)

-- Verify it's gone
local ok, err = pcall(function()
	return client:get(gvk, TEST_NAMESPACE, TEST_INTEGRATION_NAME)
end)

if ok then
	error("Expected error getting deleted resource, but succeeded")
end

return true
