local osmod = require("osmod")

-- hostname raises on error; on success returns a string
local hostname = osmod.hostname()
assert(hostname ~= "", "hostname should not be empty")
assert(type(hostname) == "string", "hostname should be a string")

return true
