---@diagnostic disable: undefined-global, lowercase-global
-- mutate_pod.lua: Adds "even-mem: true" label to pods with even memory bytes
--
-- Global variables available:
--   pod: the Kubernetes Pod object as a Lua table
--   patches: empty table to populate with JSON patch operations
--   kubernetes: the kubernetes module with helper functions

local k8s = require("kubernetes")

-- Initialize metadata if needed
pod = k8s.init_defaults(pod)

-- Function to check if a pod has even memory allocation
local function hasEvenMemory(podSpec)
	for _, container in ipairs(podSpec.containers or {}) do
		if container.resources and container.resources.requests and container.resources.requests.memory then
			local memStr = container.resources.requests.memory
			local memBytes, err = k8s.parse_memory(memStr)

			if not err and memBytes then
				-- Check if memory bytes is even
				if memBytes % 2 == 0 then
					return true
				end
			end
		end
	end
	return false
end

-- Check if pod has even memory and add label
if hasEvenMemory(pod.spec) then
	pod = k8s.add_label(pod, "even-mem", "true")

	-- Generate the JSON patch for the label
	local escaped_key = "even-mem"
	table.insert(patches, {
		op = "add",
		path = "/metadata/labels/even-mem",
		value = "true"
	})

	print(string.format("Added even-mem label to pod %s/%s",
		pod.metadata.namespace or "default",
		pod.metadata.name or "unknown"))
end
