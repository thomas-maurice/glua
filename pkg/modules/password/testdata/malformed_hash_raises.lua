-- password.verify's stored hash is OUR data, unlike hmac's tag which is
-- attacker input (see the package doc). A malformed hash must be a
-- catchable Lua error, not a silent "wrong password" false.
local password = require("password")

local ok, err = pcall(password.verify, "hunter2", "not-a-hash")
assert(not ok, "expected verify to raise on a malformed hash")
assert(type(err) == "string" and #err > 0, "expected a non-empty error message")

local ok2, err2 = pcall(password.cost, "not-a-hash")
assert(not ok2, "expected cost to raise on a malformed hash")
assert(type(err2) == "string" and #err2 > 0, "expected a non-empty error message")

-- Also confirm there is no comparison primitive exposed, mirroring hmac's
-- design (see hmac package doc): password.equal must not exist.
assert(password.equal == nil, "password.equal must not exist")

return true
