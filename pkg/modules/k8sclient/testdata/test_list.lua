
-- Test: list operation
--
-- Tests listing Kubernetes resources.

-- Create multiple ConfigMaps
local cm1 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_LIST_CONFIG_1,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_LIST_KEY_1] = TEST_LIST_VALUE_1
	}
}

local cm2 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_LIST_CONFIG_2,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_LIST_KEY_2] = TEST_LIST_VALUE_2
	}
}

client:create(cm1)
client:create(cm2)

-- List ConfigMaps
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local items = client:list(gvk, TEST_NAMESPACE)

-- Verify we have at least 2 items
local count = 0
for i, item in ipairs(items) do
	count = count + 1
end

if count < 2 then
	error("Expected at least 2 ConfigMaps, got " .. count)
end

-- Verify items have correct kind
for i, item in ipairs(items) do
	if item.kind ~= "ConfigMap" then
		error("Item " .. i .. " has wrong kind: " .. (item.kind or "nil"))
	end
end

return true
