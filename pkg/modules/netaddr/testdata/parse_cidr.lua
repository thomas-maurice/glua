-- parse_cidr: v4 and v6 prefixes, first/last, num_addresses as an exact
-- string, and the non-canonical-input normalization decision.
local netaddr = require("netaddr")

-- v4 /24.
local c24 = netaddr.parse_cidr("192.168.1.0/24")
assert(c24.cidr == "192.168.1.0/24", "cidr")
assert(c24.network == "192.168.1.0", "network")
assert(c24.prefix_len == 24, "prefix_len")
assert(c24.version == 4, "version")
assert(c24.first == "192.168.1.0", "first")
assert(c24.last == "192.168.1.255", "last")
assert(c24.num_addresses == "256", "num_addresses is a string: got " .. type(c24.num_addresses))
assert(type(c24.num_addresses) == "string", "num_addresses must be a string, not a number")

-- v4 /32: host route, first == last, 1 address.
local c32 = netaddr.parse_cidr("10.1.2.3/32")
assert(c32.first == "10.1.2.3", "host route first")
assert(c32.last == "10.1.2.3", "host route last")
assert(c32.num_addresses == "1", "host route has exactly 1 address")

-- Non-canonical input: host bits set. Decision: normalize, don't raise.
local nc = netaddr.parse_cidr("10.0.0.5/8")
assert(nc.cidr == "10.0.0.0/8", "non-canonical cidr is normalized: got " .. nc.cidr)
assert(nc.network == "10.0.0.0", "non-canonical network is normalized")

-- v6 /128: host route.
local c128 = netaddr.parse_cidr("2001:db8::1/128")
assert(c128.first == "2001:db8::1", "v6 host route first")
assert(c128.last == "2001:db8::1", "v6 host route last")
assert(c128.num_addresses == "1", "v6 host route has exactly 1 address")

-- v6 /64: 2^64 addresses -- past float64's exact integer range (2^53).
local c64 = netaddr.parse_cidr("2001:db8::/64")
assert(c64.num_addresses == "18446744073709551616", "exact 2^64: got " .. c64.num_addresses)
assert(c64.first == "2001:db8::", "v6 /64 first")
assert(c64.last == "2001:db8::ffff:ffff:ffff:ffff", "v6 /64 last")

-- v6 /0: the entire address space, 2^128 addresses -- a 39-digit string.
local c0 = netaddr.parse_cidr("::/0")
assert(#c0.num_addresses == 39, "2^128 is a 39-digit decimal string: got " .. #c0.num_addresses .. " digits")
assert(c0.num_addresses == "340282366920938463463374607431768211456", "exact 2^128")
assert(c0.first == "::", "v6 /0 first is the unspecified address")
assert(c0.version == 6, "v6 /0 version")

return true
