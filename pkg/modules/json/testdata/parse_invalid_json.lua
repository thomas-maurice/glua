-- Test: JSON parse - invalid JSON
--
-- Verifies that json.parse() returns an error
-- when given malformed JSON.

local json = require("json")
local tbl, err = json.parse('{invalid json}')

if tbl ~= nil then
	error("Expected tbl to be nil for invalid JSON")
end

if err == nil then
	error("Expected error for invalid JSON")
end
