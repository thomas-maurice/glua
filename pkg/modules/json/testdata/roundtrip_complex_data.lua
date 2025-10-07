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

-- Test: JSON round-trip - complex data
--
-- Verifies full round-trip conversion (JSON -> Lua -> JSON -> Lua)
-- with a complex data structure including nested arrays and objects.

local json = require("json")

-- Original JSON
local original = '{"users":[{"name":"Alice","age":30},{"name":"Bob","age":25}],"active":true,"count":2}'

-- Parse to table
local tbl, err1 = json.parse(original)
if err1 then
	error("Parse failed: " .. err1)
end

-- Stringify back to JSON
local stringified, err2 = json.stringify(tbl)
if err2 then
	error("Stringify failed: " .. err2)
end

-- Parse again to verify
local final, err3 = json.parse(stringified)
if err3 then
	error("Second parse failed: " .. err3)
end

-- Verify structure
if final.users[1].name ~= "Alice" then
	error("Round-trip failed: name mismatch")
end

if final.users[2].age ~= 25 then
	error("Round-trip failed: age mismatch")
end

if final.active ~= true then
	error("Round-trip failed: active mismatch")
end

if final.count ~= 2 then
	error("Round-trip failed: count mismatch")
end
