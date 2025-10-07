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

-- Test: format/parse time round-trip
--
-- Verifies that converting timestamp -> string -> timestamp
-- preserves the original value.

local k8s = require("kubernetes")

-- Format timestamp to string
local timestr, err1 = k8s.format_time(test_timestamp)
assert(err1 == nil, "format_time failed: " .. tostring(err1))

-- Parse string back to timestamp
local timestamp, err2 = k8s.parse_time(timestr)
assert(err2 == nil, "parse_time failed: " .. tostring(err2))

-- Should match original
assert(timestamp == test_timestamp, string.format("Round-trip failed: %d != %d", timestamp, test_timestamp))

return true
