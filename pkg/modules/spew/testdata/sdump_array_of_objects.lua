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

-- Test: spew.sdump - array of objects
--
-- Verifies that spew.sdump() handles arrays containing
-- object elements (tables with named fields).

local spew = require("spew")
local users = {
	{name="Alice", age=30},
	{name="Bob", age=25},
	{name="Charlie", age=35}
}

local result = spew.sdump(users)

-- Check all names are present
if not string.find(result, "Alice") then
	error("Expected result to contain 'Alice'")
end

if not string.find(result, "Bob") then
	error("Expected result to contain 'Bob'")
end

if not string.find(result, "Charlie") then
	error("Expected result to contain 'Charlie'")
end

return result
