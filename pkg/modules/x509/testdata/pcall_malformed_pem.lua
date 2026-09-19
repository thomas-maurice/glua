-- x509.parse must raise a catchable error on malformed PEM input, never
-- panic and never silently return a half-populated certificate.
local x509 = require("x509")

local ok1, err1 = pcall(x509.parse, "this is not PEM at all")
assert(not ok1, "expected x509.parse to raise on non-PEM input")
assert(type(err1) == "string" or type(err1) == "table", "pcall must return an error value")

local ok2, err2 = pcall(x509.parse, "")
assert(not ok2, "expected x509.parse to raise on empty input")

local ok3, err3 = pcall(x509.parse_chain, "")
assert(not ok3, "expected x509.parse_chain to raise on empty input")

local ok4, err4 = pcall(x509.is_valid_at, "garbage", 0)
assert(not ok4, "expected x509.is_valid_at to raise on malformed PEM")

local ok5, err5 = pcall(x509.expires_in_days, "garbage", 0)
assert(not ok5, "expected x509.expires_in_days to raise on malformed PEM")

return true
