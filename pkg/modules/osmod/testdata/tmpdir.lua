local osmod = require("osmod")

-- Test tmpdir
local tmpdir = osmod.tmpdir()
assert(tmpdir ~= "", "tmpdir should not be empty")
assert(type(tmpdir) == "string", "tmpdir should be a string")

-- Should typically be /tmp on Unix-like systems or similar
print("Temporary directory: " .. tmpdir)

return true
