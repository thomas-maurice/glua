-- Copyright (c) 2024-2025 Thomas Maurice
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
-- SOFTWARE.

-- mutate.lua: Lua mutation script for Kubernetes pods
--
-- This script adds the annotation "coucou.lil: hello" to any pod
-- that has the label "thomas.maurice/mutate=true"
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
			value = {}
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
	if c == "~" then return "~0" end
	if c == "/" then return "~1" end
	return c
end)

if annotations_created then
	-- Annotations were just created, so we need to replace the empty object
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations",
		value = {
			[annotation_key] = annotation_value
		}
	})
else
	-- Annotations exist, just add the new field
	table.insert(patches, {
		op = "add",
		path = "/metadata/annotations/" .. escaped_key,
		value = annotation_value
	})
end

-- Also add a timestamp annotation to track when the mutation occurred
local timestamp = os.date("%Y-%m-%dT%H:%M:%SZ")
table.insert(patches, {
	op = "add",
	path = "/metadata/annotations/glua.mutated-at",
	value = timestamp
})

print(string.format("Generated %d patches", #patches))
