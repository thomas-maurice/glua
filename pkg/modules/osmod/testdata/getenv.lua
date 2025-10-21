local osmod = require("osmod")

-- Test getenv with existing variable
local path = osmod.getenv("PATH")
assert(path ~= "", "PATH should not be empty")

-- Test getenv with non-existing variable
local nonexist = osmod.getenv("NONEXISTENT_VAR_XYZ")
assert(nonexist == "", "Non-existent variable should return empty string")

return true
