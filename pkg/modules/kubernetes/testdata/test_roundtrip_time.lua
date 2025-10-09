
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
