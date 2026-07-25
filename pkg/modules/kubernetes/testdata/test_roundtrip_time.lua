
-- Test: format/parse time round-trip
--
-- Verifies that converting timestamp -> string -> timestamp
-- preserves the original value.

local k8s = require("kubernetes")

-- Format timestamp to string (raises on error)
local timestr = k8s.format_time(test_timestamp)

-- Parse string back to timestamp (raises on error)
local timestamp = k8s.parse_time(timestr)

-- Should match original
assert(timestamp == test_timestamp, string.format("Round-trip failed: %d != %d", timestamp, test_timestamp))

return true
