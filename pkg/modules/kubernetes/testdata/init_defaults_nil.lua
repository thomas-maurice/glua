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

-- Test: init_defaults with nil labels and annotations
--
-- This test verifies that init_defaults() properly initializes
-- metadata.labels and metadata.annotations as empty tables when
-- they are nil, and allows adding new entries afterwards.

local k8s = require("kubernetes")

local obj = {
	metadata = {
		name = "test-pod"
	}
}

-- Before init_defaults, labels and annotations are nil
assert(obj.metadata.labels == nil, "labels should be nil initially")
assert(obj.metadata.annotations == nil, "annotations should be nil initially")

-- Initialize defaults
k8s.init_defaults(obj)

-- After init_defaults, they should be empty tables
assert(type(obj.metadata.labels) == "table", "labels should be table")
assert(type(obj.metadata.annotations) == "table", "annotations should be table")

-- Should be able to add entries
obj.metadata.labels.app = "myapp"
obj.metadata.annotations["version"] = "1.0"

assert(obj.metadata.labels.app == "myapp", "should be able to add label")
assert(obj.metadata.annotations["version"] == "1.0", "should be able to add annotation")

return true
