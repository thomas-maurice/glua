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

-- Test: init_defaults with existing labels and annotations
--
-- This test verifies that init_defaults() preserves existing
-- labels and annotations when they are already present,
-- and still allows adding new entries.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test-pod",
		labels = {
			existing = "label"
		},
		annotations = {
			existing = "annotation"
		}
	}
}

-- Initialize defaults (should not overwrite existing)
k8s.init_defaults(obj)

-- Existing values should be preserved
assert(obj.metadata.labels.existing == "label", "existing label should be preserved")
assert(obj.metadata.annotations.existing == "annotation", "existing annotation should be preserved")

-- Should still be able to add new entries
obj.metadata.labels.new = "value"
assert(obj.metadata.labels.new == "value", "should be able to add new label")

return true
