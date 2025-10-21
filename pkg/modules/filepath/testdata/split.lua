local filepath = require("filepath")

-- Test split
local dir, file = filepath.split("/home/user/file.txt")
assert(dir == "/home/user/", "dir should be '/home/user/', got: " .. dir)
assert(file == "file.txt", "file should be 'file.txt', got: " .. file)

return true
