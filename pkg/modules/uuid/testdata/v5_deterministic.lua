-- v5 is deterministic: the same (namespace, name) always yields the same
-- UUID, and it must raise on an invalid namespace.
local uuid = require("uuid")

local a = uuid.v5(uuid.NAMESPACE_DNS, "example.com")
local b = uuid.v5(uuid.NAMESPACE_DNS, "example.com")
assert(a == b, "v5 must be deterministic for identical inputs")

local c = uuid.v5(uuid.NAMESPACE_DNS, "other.example.com")
assert(a ~= c, "v5 must differ for a different name")

local ok = pcall(uuid.v5, "not-a-namespace", "example.com")
assert(not ok, "v5 should raise on an invalid namespace")

assert(uuid.parse(a).version == 5, "v5 output must report version 5")

return true
