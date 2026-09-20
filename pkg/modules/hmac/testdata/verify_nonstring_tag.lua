-- MEDIUM 1 (security review): verify_* must never raise on a non-string
-- tag. The module's own recommended usage is
-- hmac.verify_sha256(payload, secret, request.headers["x-signature"]),
-- which is nil whenever the header is absent -- that is the exact
-- production case this test pins.
local hmac = require("hmac")

local function check(fn, tag, label)
	local ok, result = pcall(fn, "message", "key", tag)
	assert(ok, label .. " must not raise for tag=" .. tostring(tag))
	assert(result == false, label .. " must return false for tag=" .. tostring(tag))
end

for _, fn_name in ipairs({"verify_sha1", "verify_sha256", "verify_sha512"}) do
	local fn = hmac[fn_name]
	check(fn, nil, fn_name)
	check(fn, 42, fn_name)
	check(fn, true, fn_name)
	check(fn, {}, fn_name)
end

return true
