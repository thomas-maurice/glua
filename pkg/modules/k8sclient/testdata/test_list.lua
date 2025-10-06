-- Test: list operation
--
-- Tests listing Kubernetes resources.

-- Create multiple ConfigMaps
local cm1 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "list-config-1",
		namespace = "default"
	},
	data = {
		key1 = "value1"
	}
}

local cm2 = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "list-config-2",
		namespace = "default"
	},
	data = {
		key2 = "value2"
	}
}

local created1, err = client.create(cm1)
if err then
	error("Failed to create first ConfigMap: " .. err)
end

local created2, err = client.create(cm2)
if err then
	error("Failed to create second ConfigMap: " .. err)
end

-- List ConfigMaps
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local items, err = client.list(gvk, "default")

if err then
	error("Failed to list: " .. err)
end

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
