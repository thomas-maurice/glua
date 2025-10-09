
-- Test: spew.dump - outputs to stdout
--
-- Verifies that spew.dump() executes without errors.
-- Note: dump() prints to stdout and returns nothing, so we
-- just verify it doesn't crash.

local spew = require("spew")
spew.dump({name="Bob", age=25})
