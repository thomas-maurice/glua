local hex = require("hex")

-- encode returns a plain string
local encoded = hex.encode("hello")
assert(encoded == "68656c6c6f", "Expected correct hex encoding, got: " .. encoded)

-- decode returns the decoded string (raises on error)
local decoded = hex.decode("68656c6c6f")
assert(decoded == "hello", "Expected 'hello', got: " .. decoded)

-- round-trip
local original = "The quick brown fox"
local rt = hex.decode(hex.encode(original))
assert(rt == original, "Round-trip should match")

return true
