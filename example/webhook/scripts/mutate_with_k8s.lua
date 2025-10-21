
-- mutate_with_k8s.lua: Advanced mutation script with Kubernetes client access
--
-- This script demonstrates how to use the k8sclient module within a webhook
-- to query Kubernetes resources during mutation
--
-- Global variables available:
---@diagnostic disable: undefined-global
--   pod: the Kubernetes Pod object being mutated
--   patches: empty table to populate with JSON patch operations
--   k8sclient: Kubernetes dynamic client (if enabled)

-- Helper: ensure annotations exist
local function ensure_annotations_path()
	if not pod.metadata.annotations then
		table.insert(patches, {
			op = "add",
			path = "/metadata/annotations",
			value = {}
		})
		return true
	end
	return false
end

print(string.format("Mutating pod: %s/%s", pod.metadata.namespace or "default", pod.metadata.name or "unknown"))

-- Ensure annotations exist
local annotations_created = ensure_annotations_path()

-- Example 1: Add basic annotation
local annotation_key = "coucou.lil"
local annotation_value = "hello"

local escaped_key = annotation_key:gsub("([%.%/~])", function(c)
	if c == "~" then return "~0" end
	if c == "/" then return "~1" end
	return c
end)

if annotations_created then
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations",
		value = {
			[annotation_key] = annotation_value
		}
	})
else
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations/" .. escaped_key,
		value = annotation_value
	})
end

-- Example 2: Add timestamp
local timestamp = os.date("%Y-%m-%dT%H:%M:%SZ")
table.insert(patches, {
	op = "add",
	path = "/metadata/annotations/glua.mutated-at",
	value = timestamp
})

-- Example 3: Query Kubernetes if client is available
if k8sclient then
	-- Create a client instance
	local client = k8sclient.new_client()

	-- Check if a ConfigMap exists in the same namespace
	local cm_gvk = {
		group = "",
		version = "v1",
		kind = "ConfigMap"
	}

	local namespace = pod.metadata.namespace or "default"
	local config_name = "webhook-config"

	-- Try to get the ConfigMap using the client
	local config, err = client:get(cm_gvk, namespace, config_name)

	if config and not err then
		print(string.format("Found ConfigMap %s in namespace %s", config_name, namespace))

		-- Add annotation indicating config was found
		table.insert(patches, {
			op = "add",
			path = "/metadata/annotations/webhook.config-found",
			value = "true"
		})

		-- If ConfigMap has specific data, use it
		if config.data and config.data["mutation-policy"] then
			local policy = config.data["mutation-policy"]
			table.insert(patches, {
				op = "add",
				path = "/metadata/annotations/webhook.policy",
				value = policy
			})
			print(string.format("Applied policy: %s", policy))
		end
	else
		print(string.format("ConfigMap %s not found or error: %s", config_name, err or "none"))

		-- Add annotation indicating config was not found
		table.insert(patches, {
			op = "add",
			path = "/metadata/annotations/webhook.config-found",
			value = "false"
		})
	end

	-- Example 4: List resources to gather metadata
	local pods, list_err = k8sclient.list(cm_gvk, namespace)
	if pods and not list_err then
		local count = 0
		for i, item in ipairs(pods.items or {}) do
			count = count + 1
		end

		table.insert(patches, {
			op = "add",
			path = "/metadata/annotations/webhook.configmaps-in-namespace",
			value = tostring(count)
		})
		print(string.format("Found %d ConfigMaps in namespace %s", count, namespace))
	end
end

print(string.format("Generated %d patches", #patches))
