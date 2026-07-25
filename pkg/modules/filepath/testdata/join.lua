local filepath = require("filepath")

-- Test join with multiple parts (pass as table)
local result = filepath.join({"/home", "user", "documents", "file.txt"})
assert(result == "/home/user/documents/file.txt", "join failed: " .. result)

-- Test join with empty table
local result2 = filepath.join({})
assert(result2 == "", "empty join should return empty string, got: " .. result2)

return true
