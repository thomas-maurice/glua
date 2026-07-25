
-- Kubernetes client integration test
-- This script creates, updates, lists, and deletes a ConfigMap

local k8sclient = require("k8sclient")
local client = k8sclient.new_client()

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

local created = client:create(configmap)
print("✓ Created ConfigMap: " .. created.metadata.name)
print("  Data: key1=" .. created.data.key1 .. ", key2=" .. created.data.key2)

-- 2. Get the ConfigMap
print("\n2. Getting ConfigMap...")
local fetched = client:get(gvk, "default", "test-config")
print("✓ Fetched ConfigMap: " .. fetched.metadata.name)
print("  UID: " .. fetched.metadata.uid)

-- 3. Update the ConfigMap
print("\n3. Updating ConfigMap...")
fetched.data.key3 = "value3"
fetched.data.key1 = "updated_value1"

local updated = client:update(fetched)
print("✓ Updated ConfigMap")
print("  Data: key1=" .. updated.data.key1 .. ", key3=" .. updated.data.key3)

-- 4. List ConfigMaps
print("\n4. Listing ConfigMaps in default namespace...")
local items = client:list(gvk, "default")
print("✓ Found " .. #items .. " ConfigMap(s)")
for i, item in ipairs(items) do
	print("  - " .. item.metadata.name)
end

-- 5. Delete the ConfigMap
print("\n5. Deleting ConfigMap...")
client:delete(gvk, "default", "test-config")
print("✓ Deleted ConfigMap")

-- 6. Verify deletion
print("\n6. Verifying deletion...")
local ok, err = pcall(function()
	return client:get(gvk, "default", "test-config")
end)
if ok then
	error("ConfigMap should have been deleted but still exists!")
end
print("✓ ConfigMap successfully deleted (not found)")

print("\n=== All tests passed! ===")
