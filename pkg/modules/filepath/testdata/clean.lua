local filepath = require("filepath")

-- Test clean
local result = filepath.clean("/home//user/../user/./file.txt")
assert(result == "/home/user/file.txt", "clean failed: " .. result)

return true
