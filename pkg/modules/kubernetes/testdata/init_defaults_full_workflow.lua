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

-- Test: init_defaults full workflow
--
-- This test demonstrates a complete workflow of using init_defaults()
-- with a pod object, verifying initial state, initialization,
-- and adding labels/annotations.

local k8s = require("kubernetes")

-- Test case 1: nil labels and annotations
local pod = {
	metadata = {
		name = "test-pod",
		namespace = "default"
	}
}

-- Before init_defaults, labels and annotations are nil
assert(pod.metadata.labels == nil, "labels should be nil initially")
assert(pod.metadata.annotations == nil, "annotations should be nil initially")

k8s.init_defaults(pod)

-- After init_defaults, they should be tables
assert(type(pod.metadata.labels) == "table", "labels should be table")
assert(type(pod.metadata.annotations) == "table", "annotations should be table")

-- Add some labels and annotations
pod.metadata.labels.app = "myapp"
pod.metadata.labels.tier = "backend"
pod.metadata.annotations.version = "1.0.0"

-- Verify they were added successfully
assert(pod.metadata.labels.app == "myapp", "label app should be set")
assert(pod.metadata.labels.tier == "backend", "label tier should be set")
assert(pod.metadata.annotations.version == "1.0.0", "annotation version should be set")

return true
