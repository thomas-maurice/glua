
-- Kubernetes client integration test
-- This script creates, updates, lists, and deletes a ConfigMap

local client = require("k8sclient")

print("=== Kubernetes Client Integration Test ===\n")

-- Define GVK for ConfigMap
local gvk = {group = "", version = "v1", kind = "ConfigMap"}

-- 1. Create a ConfigMap
print("1. Creating ConfigMap...")
local configmap = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "test-config",
		namespace = "default"
	},
	data = {
		key1 = "value1",
		key2 = "value2"
	}
}

local created, err = client.create(configmap)
if err then
	error("Failed to create ConfigMap: " .. err)
end
print("✓ Created ConfigMap: " .. created.metadata.name)
print("  Data: key1=" .. created.data.key1 .. ", key2=" .. created.data.key2)

-- 2. Get the ConfigMap
print("\n2. Getting ConfigMap...")
local fetched, err = client.get(gvk, "default", "test-config")
if err then
	error("Failed to get ConfigMap: " .. err)
end
print("✓ Fetched ConfigMap: " .. fetched.metadata.name)
print("  UID: " .. fetched.metadata.uid)

-- 3. Update the ConfigMap
print("\n3. Updating ConfigMap...")
fetched.data.key3 = "value3"
fetched.data.key1 = "updated_value1"

local updated, err = client.update(fetched)
if err then
	error("Failed to update ConfigMap: " .. err)
end
print("✓ Updated ConfigMap")
print("  Data: key1=" .. updated.data.key1 .. ", key3=" .. updated.data.key3)

-- 4. List ConfigMaps
print("\n4. Listing ConfigMaps in default namespace...")
local items, err = client.list(gvk, "default")
if err then
	error("Failed to list ConfigMaps: " .. err)
end
print("✓ Found " .. #items .. " ConfigMap(s)")
for i, item in ipairs(items) do
	print("  - " .. item.metadata.name)
end

-- 5. Delete the ConfigMap
print("\n5. Deleting ConfigMap...")
local err = client.delete(gvk, "default", "test-config")
if err then
	error("Failed to delete ConfigMap: " .. err)
end
print("✓ Deleted ConfigMap")

-- 6. Verify deletion
print("\n6. Verifying deletion...")
local fetched, err = client.get(gvk, "default", "test-config")
if not err then
	error("ConfigMap should have been deleted but still exists!")
end
print("✓ ConfigMap successfully deleted (not found)")

print("\n=== All tests passed! ===")
