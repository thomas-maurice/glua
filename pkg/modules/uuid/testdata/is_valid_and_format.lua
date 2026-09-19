-- is_valid never raises, even on garbage; format raises on a bad style;
-- format round-trips through all four accepted spellings.
local uuid = require("uuid")

assert(uuid.is_valid(uuid.v4()) == true, "a freshly generated v4 must be valid")
assert(uuid.is_valid("not-a-uuid") == false, "garbage must not be valid")
assert(uuid.is_valid("") == false, "empty string must not be valid")

local ok, err = pcall(uuid.parse, "not-a-uuid")
assert(not ok, "parse should raise on garbage input")

local id = uuid.v4()
local ok2, err2 = pcall(uuid.format, id, "nope")
assert(not ok2, "format should raise on an unknown style")
assert(string.find(err2, "style") or string.find(err2, "nope"),
  "error should mention the bad style, got: " .. tostring(err2))

-- format round-trips: parse(format(id, style)) still reports the same
-- canonical uuid for every style.
for _, style in ipairs({"canonical", "plain", "urn", "braced"}) do
  local formatted = uuid.format(id, style)
  local info = uuid.parse(formatted)
  assert(info.uuid == id, "round-trip through style " .. style .. " changed the uuid")
end

return true
