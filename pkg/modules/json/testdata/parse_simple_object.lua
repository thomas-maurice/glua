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

-- Test: JSON parse - simple object
--
-- Verifies that json.parse() can parse a simple JSON object
-- with string, number, and boolean fields.

local json = require("json")
local tbl, err = json.parse('{"name":"John","age":30,"active":true}')

if err then
	error("Parse failed: " .. err)
end

if tbl.name ~= "John" then
	error("Expected name to be 'John', got: " .. tostring(tbl.name))
end

if tbl.age ~= 30 then
	error("Expected age to be 30, got: " .. tostring(tbl.age))
end

if tbl.active ~= true then
	error("Expected active to be true, got: " .. tostring(tbl.active))
end
