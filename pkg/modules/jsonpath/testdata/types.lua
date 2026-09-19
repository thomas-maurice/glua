-- Type fidelity through the round trip: a matched integer must not come
-- back with a spurious ".0", a non-integer float must keep its fraction,
-- and a boolean must stay a boolean, not become a string.
local jsonpath = require("jsonpath")

local data = {
  name = "web",
  replicas = 3,
  cpu_limit = 0.5,
  enabled = true,
  disabled = false,
}

local name = jsonpath.first(data, ".name", nil)
assert(type(name) == "string" and name == "web", "string should round-trip as a string")

local replicas = jsonpath.first(data, ".replicas", nil)
assert(type(replicas) == "number" and replicas == 3, "integer should round-trip as a number equal to 3")
assert(tostring(replicas) == "3", "integer should round-trip without a spurious .0")

local cpu = jsonpath.first(data, ".cpu_limit", nil)
assert(type(cpu) == "number" and cpu == 0.5, "non-integer float should keep its fraction")

local enabled = jsonpath.first(data, ".enabled", nil)
assert(type(enabled) == "boolean" and enabled == true, "true should round-trip as a boolean, not a string")

local disabled = jsonpath.first(data, ".disabled", nil)
assert(type(disabled) == "boolean" and disabled == false, "false should round-trip as a boolean, not a string")

-- A mixed-type query result table must preserve each element's own type.
local mixed = jsonpath.query(data, ".*")
local sawNumber, sawString, sawBoolean = false, false, false
for _, v in ipairs(mixed) do
  local t = type(v)
  if t == "number" then sawNumber = true end
  if t == "string" then sawString = true end
  if t == "boolean" then sawBoolean = true end
end
assert(sawNumber and sawString and sawBoolean, "a mixed-type document's matches must each keep their own Lua type")

return true
