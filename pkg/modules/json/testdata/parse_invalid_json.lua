-- Test: JSON parse - invalid JSON raises
--
-- Verifies that json.parse() raises a Lua error when given malformed JSON.

local json = require("json")

local ok, err = pcall(json.parse, '{invalid json}')
assert(not ok, "Expected error for invalid JSON")
assert(type(err) == "string", "Error should be a string message")
