local osmod = require("osmod")

-- setenv returns true on success (raises on error)
local ok = osmod.setenv("TEST_GLUA_VAR", "hello")
assert(ok == true, "setenv should return true on success, got: " .. tostring(ok))

-- Verify it was set
local value = osmod.getenv("TEST_GLUA_VAR")
assert(value == "hello", "Environment variable should be set to 'hello', got: " .. value)

-- unsetenv returns true on success
ok = osmod.unsetenv("TEST_GLUA_VAR")
assert(ok == true, "unsetenv should return true on success, got: " .. tostring(ok))

-- Verify it was unset
value = osmod.getenv("TEST_GLUA_VAR")
assert(value == "", "Environment variable should be empty after unset")

-- setenv with empty name raises — catch with pcall
local raised, err = pcall(osmod.setenv, "", "anything")
assert(not raised, "setenv with empty name should raise an error")
assert(type(err) == "string", "Error should be a string")

return true
