local base64 = require("base64")

-- encode returns a plain string
local encoded = base64.encode("hello world")
assert(encoded == "aGVsbG8gd29ybGQ=", "Expected correct base64 encoding, got: " .. encoded)

-- decode returns the decoded string (raises on error)
local decoded = base64.decode("aGVsbG8gd29ybGQ=")
assert(decoded == "hello world", "Expected 'hello world', got: " .. decoded)

-- round-trip
local original = "The quick brown fox"
local rt = base64.decode(base64.encode(original))
assert(rt == original, "Round-trip should match")

-- URL-safe round-trip
local urlrt = base64.decode_url(base64.encode_url(original))
assert(urlrt == original, "URL-safe round-trip should match")

return true
