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

-- Test: create operation
--
-- Tests creating a new Kubernetes resource.

local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = TEST_NEW_CONFIG_NAME,
		namespace = TEST_NAMESPACE
	},
	data = {
		[TEST_DATA_KEY] = TEST_DATA_VALUE
	}
}

local created, err = client.create(cm)

if err then
	error("Failed to create: " .. err)
end

if created.metadata.name ~= TEST_NEW_CONFIG_NAME then
	error("Expected name '" .. TEST_NEW_CONFIG_NAME .. "', got " .. tostring(created.metadata.name))
end

if created.data[TEST_DATA_KEY] ~= TEST_DATA_VALUE then
	error("Expected data." .. TEST_DATA_KEY .. " '" .. TEST_DATA_VALUE .. "', got " .. tostring(created.data[TEST_DATA_KEY]))
end

return true
