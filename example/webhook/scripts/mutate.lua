-- mutate.lua: Lua mutation script for Kubernetes pods
--
-- This script adds the annotation "coucou.lil: hello" to any pod
-- that has the label "glua-webhook/mutate=true"
--
-- Global variables available:
--   pod: the Kubernetes Pod object as a Lua table
--   patches: empty table to populate with JSON patch operations
--
-- Example usage:
--   The patches table should contain JSON patch operations like:
--   {
--     op = "add",
--     path = "/metadata/annotations/coucou.lil",
--     value = "hello"
--   }

-- Helper function to check if annotations exist
local function ensure_annotations_path()
	if not pod.metadata.annotations then
		-- If annotations don't exist, we need to create the whole annotations object
		table.insert(patches, {
			op = "add",
			path = "/metadata/annotations",
			value = {},
		})
		return true
	end
	return false
end

-- Main mutation logic
print(string.format("Mutating pod: %s/%s", pod.metadata.namespace or "default", pod.metadata.name or "unknown"))

-- Ensure annotations exist
local annotations_created = ensure_annotations_path()

-- Add the coucou.lil annotation
-- If we just created annotations, the value is already empty and we'll add to it
-- If annotations existed, we add a new field to the existing object
local annotation_key = "coucou.lil"
local annotation_value = "hello"

-- Escape special characters in the annotation key for JSON path
local escaped_key = annotation_key:gsub("([%.%/~])", function(c)
	if c == "~" then
		return "~0"
	end
	if c == "/" then
		return "~1"
	end
	return c
end)

if annotations_created then
	-- Annotations were just created, so we need to replace the empty object
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations",
		value = {
			[annotation_key] = annotation_value,
		},
	})
else
	-- Annotations exist, just add the new field
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations/" .. escaped_key,
		value = annotation_value,
	})
end

-- Also add a timestamp annotation to track when the mutation occurred
local timestamp = os.date("%Y-%m-%dT%H:%M:%SZ")
table.insert(patches, {
	op = "add",
	path = "/metadata/annotations/glua.mutated-at",
	value = timestamp,
})

print(string.format("Generated %d patches", #patches))
