-- parse_ip: v4, v6, v4-mapped v6, and every IPInfo predicate field.
local netaddr = require("netaddr")

-- Plain v4.
local v4 = netaddr.parse_ip("192.168.1.1")
assert(v4.ip == "192.168.1.1", "v4 ip")
assert(v4.version == 4, "v4 version")
assert(v4.is_private == true, "192.168.1.1 is private")
assert(v4.is_global == false, "192.168.1.1 is not global")
assert(v4.is_loopback == false, "192.168.1.1 is not loopback")

-- Public v4.
local pub = netaddr.parse_ip("8.8.8.8")
assert(pub.is_private == false, "8.8.8.8 is not private")
assert(pub.is_global == true, "8.8.8.8 is global")

-- Plain v6, compressed on output.
local v6 = netaddr.parse_ip("2001:0db8:0000::1")
assert(v6.ip == "2001:db8::1", "v6 ip is compressed: got " .. v6.ip)
assert(v6.version == 6, "v6 version")

-- v6 loopback and unspecified.
local loop6 = netaddr.parse_ip("::1")
assert(loop6.is_loopback == true, "::1 is loopback")
assert(loop6.is_private == true, "::1 counts as private (D5)")

local unspec6 = netaddr.parse_ip("::")
assert(unspec6.is_unspecified == true, ":: is unspecified")
assert(unspec6.is_private == true, ":: counts as private (D5)")

-- v6 ULA and link-local.
local ula = netaddr.parse_ip("fc00::1")
assert(ula.is_private == true, "fc00::1 (ULA) is private")

local ll = netaddr.parse_ip("fe80::1")
assert(ll.is_link_local == true, "fe80::1 is link-local")
assert(ll.is_private == true, "fe80::1 counts as private (D5)")

-- v4-mapped v6 is unwrapped to its plain v4 form everywhere.
local mapped = netaddr.parse_ip("::ffff:192.0.2.1")
assert(mapped.ip == "192.0.2.1", "v4-mapped v6 is unwrapped: got " .. mapped.ip)
assert(mapped.version == 4, "v4-mapped v6 reports version 4")

-- Zone-bearing link-local address: netip accepts it, must not raise.
local zoned = netaddr.parse_ip("fe80::1%eth0")
assert(zoned.ip == "fe80::1%eth0", "zone survives in parse_ip's ip field: got " .. zoned.ip)
assert(zoned.is_link_local == true, "zoned address is still link-local")

-- Multicast.
local mcast = netaddr.parse_ip("ff02::1")
assert(mcast.is_multicast == true, "ff02::1 is multicast")

return true
