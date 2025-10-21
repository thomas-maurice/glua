local filepath = require("filepath")

-- Test ext
local ext = filepath.ext("/home/user/file.txt")
assert(ext == ".txt", "ext should be '.txt', got: " .. ext)

-- Test base
local base = filepath.base("/home/user/file.txt")
assert(base == "file.txt", "base should be 'file.txt', got: " .. base)

-- Test dir
local dir = filepath.dir("/home/user/file.txt")
assert(dir == "/home/user", "dir should be '/home/user', got: " .. dir)

return true
