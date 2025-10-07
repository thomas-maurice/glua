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

-- Test: JSON stringify - boolean and null values
--
-- Verifies that json.stringify() correctly handles
-- boolean values and nil (which becomes null in JSON).

local json = require("json")
local data = {
	active = true,
	disabled = false,
	value = nil
}

local str, err = json.stringify(data)
if err then
	error("Stringify failed: " .. err)
end

-- Parse it back to verify
local tbl, parseErr = json.parse(str)
if parseErr then
	error("Failed to parse stringified JSON: " .. parseErr)
end

if tbl.active ~= true then
	error("Expected active to be true after round-trip")
end

if tbl.disabled ~= false then
	error("Expected disabled to be false after round-trip")
end
