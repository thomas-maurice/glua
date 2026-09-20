-- normalize_ip and ip_version: compression, zone stripping, v4-mapped unwrap.
local netaddr = require("netaddr")

assert(netaddr.normalize_ip("2001:0db8:0000::1") == "2001:db8::1", "compresses v6")
assert(netaddr.normalize_ip("::ffff:1.2.3.4") == "1.2.3.4", "unwraps v4-mapped v6")
assert(netaddr.normalize_ip("fe80::1%eth0") == "fe80::1", "strips zone")

assert(netaddr.ip_version("192.168.1.1") == 4, "v4 version")
assert(netaddr.ip_version("2001:db8::1") == 6, "v6 version")
assert(netaddr.ip_version("::ffff:8.8.8.8") == 4, "v4-mapped v6 reports version 4")
assert(netaddr.ip_version("::") == 6, ":: is v6")

-- is_ip / parse_ip pairing: is_ip must return false exactly where parse_ip raises.
assert(netaddr.is_ip("192.168.1.1") == true, "valid v4")
assert(netaddr.is_ip("2001:db8::1") == true, "valid v6")
assert(netaddr.is_ip("not an ip") == false, "garbage")
local ok = pcall(netaddr.parse_ip, "not an ip")
assert(ok == false, "parse_ip must raise on the same input is_ip rejects")

return true
