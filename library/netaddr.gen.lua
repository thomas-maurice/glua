---@meta netaddr

---@class netaddr.CIDRInfo
---@field cidr string
---@field first string
---@field last string
---@field network string
---@field num_addresses string
---@field prefix_len number
---@field version number

---@class netaddr.IPInfo
---@field ip string
---@field is_global boolean
---@field is_link_local boolean
---@field is_loopback boolean
---@field is_multicast boolean
---@field is_private boolean
---@field is_unspecified boolean
---@field version number

---@class netaddr
local netaddr = {}

--- parses an IP address and reports version and policy-relevant predicates, raises on invalid input
---@param s string the address to parse; a v4-mapped v6 address (::ffff:192.0.2.1) is unwrapped to its plain v4 form
---@return netaddr.IPInfo info IPInfo table: ip, version, is_private, is_loopback, is_global, is_multicast, is_unspecified, is_link_local
function netaddr.parse_ip(s) end

--- reports whether s parses as a valid IP address; never raises
---@param s string the string to check
---@return boolean ok true if parse_ip(s) would succeed
function netaddr.is_ip(s) end

--- returns the canonical string form of an address, raises on invalid input
---@param s string the address to normalize
---@return string out compressed IPv6 / zone stripped / v4-mapped v6 unwrapped to plain v4
function netaddr.normalize_ip(s) end

--- returns the IP version, raises on invalid input
---@param s string the address to check; a v4-mapped v6 address reports 4
---@return number version 4 or 6
function netaddr.ip_version(s) end

--- reports whether an address is private, raises on invalid input
---@param s string the address to check
---@return boolean ok true for RFC 1918 (v4), ULA fc00::/7, link-local, loopback or unspecified
function netaddr.is_private(s) end

--- reports whether an address is a loopback address, raises on invalid input
---@param s string the address to check
---@return boolean ok true for 127.0.0.0/8 or ::1
function netaddr.is_loopback(s) end

--- reports whether an address is globally routable, raises on invalid input
---@param s string the address to check
---@return boolean ok true if global unicast and not private/loopback/link-local/documentation/CGNAT
function netaddr.is_global(s) end

--- parses a CIDR prefix, raises on invalid input
---@param s string the prefix to parse, e.g. "10.0.0.0/8" or "2001:db8::/32"; a non-canonical input like "10.0.0.5/8" is normalized to its network address
---@return netaddr.CIDRInfo info CIDRInfo table: cidr, network, prefix_len, version, first, last, num_addresses (a decimal string, exact even for a v6 /0)
function netaddr.parse_cidr(s) end

--- reports whether an IP falls inside a CIDR prefix, raises on a malformed argument
---@param cidr string the containing prefix
---@param ip string the address to test; a v4-mapped v6 address is compared as its plain v4 form
---@return boolean ok false (not an error) if cidr and ip are valid but of different address families
function netaddr.cidr_contains(cidr, ip) end

--- reports whether inner is fully contained within outer, raises on a malformed argument
---@param inner string the candidate subnet
---@param outer string the candidate supernet; a prefix is a subnet of itself
---@return boolean ok false (not an error) if inner and outer are valid but of different address families
function netaddr.subnet_of(inner, outer) end

--- reports whether two CIDR prefixes share any address, raises on a malformed argument
---@param a string the first prefix
---@param b string the second prefix
---@return boolean ok false (not an error) if a and b are valid but of different address families
function netaddr.cidr_overlaps(a, b) end

return netaddr
