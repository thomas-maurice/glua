-- Every documented error path, asserted via pcall.
local netaddr = require("netaddr")

local function expect_error(fn, ...)
  local ok, err = pcall(fn, ...)
  assert(not ok, "expected an error but call succeeded")
  assert(type(err) == "string", "expected a string error message")
  return err
end

-- parse_ip / normalize_ip / ip_version / is_private / is_loopback / is_global
-- all raise on genuinely unparseable input.
expect_error(netaddr.parse_ip, "not an ip")
expect_error(netaddr.normalize_ip, "not an ip")
expect_error(netaddr.ip_version, "not an ip")
expect_error(netaddr.is_private, "not an ip")
expect_error(netaddr.is_loopback, "not an ip")
expect_error(netaddr.is_global, "not an ip")
expect_error(netaddr.parse_ip, "999.999.999.999")
expect_error(netaddr.parse_ip, "10.0.0.0/8") -- parse_ip rejects CIDR notation

-- parse_cidr raises on garbage and on an out-of-range prefix length.
expect_error(netaddr.parse_cidr, "not a cidr")
local err = expect_error(netaddr.parse_cidr, "10.0.0.0/33")
assert(err:find("between 0 and 32", 1, true) ~= nil, "v4 prefix-length error names the bound: " .. err)
local err6 = expect_error(netaddr.parse_cidr, "2001:db8::/200")
assert(err6:find("between 0 and 128", 1, true) ~= nil, "v6 prefix-length error names the bound: " .. err6)

-- cidr_contains / subnet_of / cidr_overlaps raise on a malformed argument
-- (but NOT on a valid cross-family comparison -- see cidr_relations.lua).
expect_error(netaddr.cidr_contains, "garbage", "10.0.0.1")
expect_error(netaddr.cidr_contains, "10.0.0.0/8", "garbage")
expect_error(netaddr.subnet_of, "garbage", "10.0.0.0/8")
expect_error(netaddr.cidr_overlaps, "garbage", "10.0.0.0/8")

return true
