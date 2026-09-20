-- cidr_contains / subnet_of / cidr_overlaps, within a family, across
-- families, and the v4-mapped v6 host decision.
local netaddr = require("netaddr")

-- cidr_contains: v4 within v4.
assert(netaddr.cidr_contains("10.0.0.0/8", "10.1.2.3") == true, "v4 in v4")
assert(netaddr.cidr_contains("10.0.0.0/8", "11.1.2.3") == false, "v4 not in v4")

-- cidr_contains: v6 host in v6 cidr.
assert(netaddr.cidr_contains("2001:db8::/32", "2001:db8::1") == true, "v6 in v6")
assert(netaddr.cidr_contains("2001:db8::/32", "2001:db9::1") == false, "v6 not in v6")

-- cidr_contains: cross-family returns false, does not raise.
assert(netaddr.cidr_contains("10.0.0.0/8", "2001:db8::1") == false, "v6 host against v4 cidr is false")
assert(netaddr.cidr_contains("2001:db8::/32", "10.0.0.1") == false, "v4 host against v6 cidr is false")

-- cidr_contains: a v4-mapped v6 host is compared as its unwrapped v4 form.
assert(netaddr.cidr_contains("10.0.0.0/8", "::ffff:10.1.2.3") == true, "v4-mapped v6 host matches v4 cidr")

-- subnet_of: same-prefix and adjacent-prefix edges.
assert(netaddr.subnet_of("10.0.0.0/24", "10.0.0.0/24") == true, "a prefix is a subnet of itself")
assert(netaddr.subnet_of("10.0.0.0/25", "10.0.0.0/24") == true, "more specific prefix is a subnet")
assert(netaddr.subnet_of("10.0.1.0/24", "10.0.0.0/24") == false, "adjacent prefix is not a subnet")
assert(netaddr.subnet_of("10.0.0.0/24", "10.0.0.0/25") == false, "less specific prefix is not a subnet of a more specific one")

-- subnet_of: v6 within v6, and cross-family false.
assert(netaddr.subnet_of("2001:db8:1::/48", "2001:db8::/32") == true, "v6 subnet")
assert(netaddr.subnet_of("10.0.0.0/24", "2001:db8::/32") == false, "cross-family subnet_of is false")

-- cidr_overlaps: same-prefix, adjacent, and containment overlap.
assert(netaddr.cidr_overlaps("10.0.0.0/24", "10.0.0.0/24") == true, "identical prefixes overlap")
assert(netaddr.cidr_overlaps("10.0.0.0/24", "10.0.1.0/24") == false, "adjacent prefixes do not overlap")
assert(netaddr.cidr_overlaps("10.0.0.0/16", "10.0.1.0/24") == true, "containing prefixes overlap")

-- cidr_overlaps: v6 within v6, and cross-family false.
assert(netaddr.cidr_overlaps("2001:db8::/32", "2001:db8:1::/48") == true, "v6 overlap")
assert(netaddr.cidr_overlaps("10.0.0.0/8", "::/0") == false, "cross-family cidr_overlaps is false")

return true
