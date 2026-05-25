local hex = require("hex")

-- decode raises on invalid input; catch with pcall
local ok, err = pcall(hex.decode, "not valid hex!")
assert(not ok, "Expected error for invalid hex")
assert(type(err) == "string", "Error should be a string")

-- odd-length string is also invalid
local ok2, err2 = pcall(hex.decode, "123")
assert(not ok2, "Expected error for odd-length hex")
assert(type(err2) == "string", "Error should be a string")

return true
