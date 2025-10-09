
-- Test: spew module loading
--
-- Verifies that the spew module loads correctly
-- and exposes the expected functions.

local spew = require("spew")

-- Verify module has expected functions
if type(spew.dump) ~= "function" then
	error("Expected spew.dump to be a function")
end

if type(spew.sdump) ~= "function" then
	error("Expected spew.sdump to be a function")
end
