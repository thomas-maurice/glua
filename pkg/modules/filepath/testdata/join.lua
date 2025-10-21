local filepath = require("filepath")

-- Test join with multiple parts
local result = filepath.join("/home", "user", "documents", "file.txt")
assert(result == "/home/user/documents/file.txt", "join failed: " .. result)

-- Test join with empty
local result2 = filepath.join()
assert(result2 == "", "empty join should return empty string")

return true
