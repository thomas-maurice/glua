local osmod = require("osmod")

-- Test setenv
local err = osmod.setenv("TEST_GLUA_VAR", "hello")
assert(err == nil, "setenv should not return error")

-- Verify it was set
local value = osmod.getenv("TEST_GLUA_VAR")
assert(value == "hello", "Environment variable should be set to 'hello', got: " .. value)

-- Test unsetenv
err = osmod.unsetenv("TEST_GLUA_VAR")
assert(err == nil, "unsetenv should not return error")

-- Verify it was unset
value = osmod.getenv("TEST_GLUA_VAR")
assert(value == "", "Environment variable should be empty after unset")

return true
