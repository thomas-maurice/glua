-- The intended usage: hash once, verify on login, rehash if the stored cost
-- is stale. Cost 4 is used throughout purely to keep this test fast; real
-- callers should pass password.DEFAULT_COST.
local password = require("password")

local TEST_COST = 4

local stored = password.hash("hunter2", TEST_COST)
assert(type(stored) == "string", "expected string hash")
assert(stored:sub(1, 7) == "$2a$04$", "expected $2a$04$ prefix, got " .. stored:sub(1, 7))

-- Correct password verifies.
assert(password.verify("hunter2", stored), "expected verify to succeed for the correct password")

-- Wrong password returns false, not an error.
assert(not password.verify("wrong-password", stored), "expected verify to fail for the wrong password")

-- cost() returns what we hashed with, which is what rehash-on-login compares against.
assert(password.cost(stored) == TEST_COST, "expected cost to round-trip")

return true
