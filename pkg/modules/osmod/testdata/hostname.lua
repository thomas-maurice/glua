local osmod = require("osmod")

-- Test hostname
local hostname, err = osmod.hostname()
assert(err == nil, "hostname should not return error: " .. tostring(err))
assert(hostname ~= "", "hostname should not be empty")
assert(type(hostname) == "string", "hostname should be a string")

return true
