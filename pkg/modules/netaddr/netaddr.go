// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package netaddr provides IP address and CIDR arithmetic for Lua scripts as
// pure computation. It deliberately imports only net/netip and math/big, and
// NEVER net — net is the package that carries the DNS resolver, so not
// importing it makes "no name resolution, ever" a property of the import
// list a reviewer can check in one glance, rather than a promise in a
// comment. netaddr_import_test.go enforces this at test time.
//
// IPv6 is a first-class citizen throughout this module, not an afterthought:
//
//   - v4-mapped v6 addresses (::ffff:192.0.2.1) are unwrapped (net/netip's
//     Unmap) before every computation, including inside parse_ip. A
//     v4-mapped v6 address therefore reports version 4 and behaves
//     identically, everywhere in this module, to its plain v4 form. This is
//     the same treatment normalize_ip's spec already mandates; it is applied
//     consistently rather than only there.
//   - is_private/is_global/is_loopback cover the v6 special ranges (ULA
//     fc00::/7, link-local fe80::/10, loopback ::1, unspecified ::) in
//     addition to the v4 ones — see isPrivateAddr/isGlobalAddr below for the
//     exact, auditable definitions.
//   - cidr_contains, subnet_of and cidr_overlaps return false (never raise)
//     when comparing across address families (a v4 host/prefix against a v6
//     one, or vice versa) after unmapping. Only a malformed IP/CIDR string
//     raises.
//   - num_addresses is a decimal string, not a Lua number: a /0 IPv6 prefix
//     holds 2^128 addresses, far past what a float64 can hold exactly. This
//     costs v4 callers a tonumber() call for the common case, in exchange
//     for v6 callers never receiving a silently-wrong count.
package netaddr

import (
	"fmt"
	"math/big"
	"net/netip"
	"strings"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// IPInfo: everything parse_ip reports about a single address. For a
// v4-mapped v6 address, every field reflects the unwrapped (v4) form — see
// the package doc.
type IPInfo struct {
	IP            string `json:"ip"`
	Version       int    `json:"version"`
	IsPrivate     bool   `json:"is_private"`
	IsLoopback    bool   `json:"is_loopback"`
	IsGlobal      bool   `json:"is_global"`
	IsMulticast   bool   `json:"is_multicast"`
	IsUnspecified bool   `json:"is_unspecified"`
	IsLinkLocal   bool   `json:"is_link_local"`
}

// CIDRInfo: everything parse_cidr reports about a prefix. cidr and network
// are both the normalised (masked) form: a non-canonical input like
// "10.0.0.5/8" is normalised to network "10.0.0.0" / cidr "10.0.0.0/8"
// rather than rejected — see parseCIDRInfo.
type CIDRInfo struct {
	CIDR         string `json:"cidr"`
	Network      string `json:"network"`
	PrefixLen    int    `json:"prefix_len"`
	Version      int    `json:"version"`
	First        string `json:"first"`
	Last         string `json:"last"`
	NumAddresses string `json:"num_addresses"`
}

// v4GlobalExclusions: IPv4 ranges that are IsGlobalUnicast() but are
// documentation/testing/CGNAT ranges, not genuinely globally routable.
// Verbatim from SPECS.md S5's is_global definition.
var v4GlobalExclusions = []netip.Prefix{
	netip.MustParsePrefix("100.64.0.0/10"),   // CGNAT (RFC 6598)
	netip.MustParsePrefix("192.0.0.0/24"),    // IETF protocol assignments
	netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1
	netip.MustParsePrefix("198.18.0.0/15"),   // benchmarking
	netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2
	netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3
	netip.MustParsePrefix("240.0.0.0/4"),     // reserved
}

// v6GlobalExclusions: IPv6 documentation/discard ranges excluded from
// is_global for the same reason as v4GlobalExclusions.
var v6GlobalExclusions = []netip.Prefix{
	netip.MustParsePrefix("100::/64"),      // discard-only
	netip.MustParsePrefix("2001:db8::/32"), // documentation
}

// isPrivateAddr: exact definition of is_private, auditable in one place.
// Covers the v4 RFC 1918 set (via netip's own IsPrivate) plus, explicitly,
// the v6 ranges the spec calls out that IsPrivate() alone does not cover:
// link-local (fe80::/10, and v4's 169.254.0.0/16 via IsLinkLocalUnicast),
// loopback (::1 / 127.0.0.0/8) and unspecified (:: / 0.0.0.0). addr must
// already be unmapped.
func isPrivateAddr(addr netip.Addr) bool {
	return addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() || addr.IsUnspecified()
}

// isGlobalAddr: exact definition of is_global, auditable in one place.
// Stricter than IsGlobalUnicast() alone (which admits RFC 1918 space and
// admits documentation/testing ranges): a policy that trusts is_global must
// not be satisfied by an address that will never actually route on the
// public Internet. addr must already be unmapped.
func isGlobalAddr(addr netip.Addr) bool {
	if !addr.IsValid() || !addr.IsGlobalUnicast() {
		return false
	}
	if addr.IsPrivate() || addr.IsLoopback() || addr.IsLinkLocalUnicast() {
		return false
	}
	exclusions := v6GlobalExclusions
	if addr.Is4() {
		exclusions = v4GlobalExclusions
	}
	for _, p := range exclusions {
		if p.Contains(addr) {
			return false
		}
	}
	return true
}

// parseAddr: parses s and unwraps a v4-mapped v6 address, so every caller
// downstream sees the same (v4) form regardless of which syntax was used.
// name is the calling Lua function's name, used to prefix the error.
func parseAddr(name, s string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(s)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("%s: %w", name, err)
	}
	return addr.Unmap(), nil
}

// parsePrefix: parses s as a CIDR prefix and normalises it to its network
// address (Masked()), so a non-canonical input like "10.0.0.5/8" behaves
// identically to "10.0.0.0/8" rather than being rejected. On a prefix-length
// error, re-parses just the address part (if possible) to name the actual
// bound (0-32 for v4, 0-128 for v6) in the error, since Go's own message
// ("prefix length out of range") does not state it.
func parsePrefix(name, s string) (netip.Prefix, error) {
	p, err := netip.ParsePrefix(s)
	if err != nil {
		if idx := strings.LastIndex(s, "/"); idx >= 0 {
			if addr, addrErr := netip.ParseAddr(s[:idx]); addrErr == nil {
				maxBits := 128
				if addr.Is4() {
					maxBits = 32
				}
				return netip.Prefix{}, fmt.Errorf("%s: %q: prefix length must be between 0 and %d: %w", name, s, maxBits, err)
			}
		}
		return netip.Prefix{}, fmt.Errorf("%s: %w", name, err)
	}
	return p.Masked(), nil
}

// addrToBigInt: renders addr as an unsigned big.Int over its network-order bytes.
func addrToBigInt(addr netip.Addr) *big.Int {
	return new(big.Int).SetBytes(addr.AsSlice())
}

// bigIntToAddr: the inverse of addrToBigInt; is4 selects a 4-byte or 16-byte rendering.
func bigIntToAddr(i *big.Int, is4 bool) netip.Addr {
	size := 16
	if is4 {
		size = 4
	}
	buf := make([]byte, size)
	b := i.Bytes()
	copy(buf[size-len(b):], b)
	if is4 {
		var a [4]byte
		copy(a[:], buf)
		return netip.AddrFrom4(a)
	}
	var a [16]byte
	copy(a[:], buf)
	return netip.AddrFrom16(a)
}

// parseIP: parses s and reports version, family and policy-relevant
// predicates. Raises on anything netip cannot parse as a bare address (no
// CIDR notation, no DNS name).
func parseIP(s string) (IPInfo, error) {
	addr, err := parseAddr("netaddr.parse_ip", s)
	if err != nil {
		return IPInfo{}, err
	}
	version := 6
	if addr.Is4() {
		version = 4
	}
	return IPInfo{
		IP:            addr.String(),
		Version:       version,
		IsPrivate:     isPrivateAddr(addr),
		IsLoopback:    addr.IsLoopback(),
		IsGlobal:      isGlobalAddr(addr),
		IsMulticast:   addr.IsMulticast(),
		IsUnspecified: addr.IsUnspecified(),
		IsLinkLocal:   addr.IsLinkLocalUnicast(),
	}, nil
}

// isIP: non-raising predicate companion to parse_ip, for callers who want a
// boolean instead of a pcall.
func isIP(s string) bool {
	_, err := netip.ParseAddr(s)
	return err == nil
}

// normalizeIP: returns the canonical string form of s — compressed IPv6,
// zone stripped, v4-mapped v6 unwrapped to plain v4. Raises on invalid input.
func normalizeIP(s string) (string, error) {
	addr, err := parseAddr("netaddr.normalize_ip", s)
	if err != nil {
		return "", err
	}
	return addr.WithZone("").String(), nil
}

// ipVersion: returns 4 or 6. A v4-mapped v6 address reports 4 (it is
// unwrapped first — see the package doc). Raises on invalid input.
func ipVersion(s string) (float64, error) {
	addr, err := parseAddr("netaddr.ip_version", s)
	if err != nil {
		return 0, err
	}
	if addr.Is4() {
		return 4, nil
	}
	return 6, nil
}

// isPrivate: sugar for parse_ip(s).is_private. Raises on invalid input.
func isPrivate(s string) (bool, error) {
	addr, err := parseAddr("netaddr.is_private", s)
	if err != nil {
		return false, err
	}
	return isPrivateAddr(addr), nil
}

// isLoopback: sugar for parse_ip(s).is_loopback. Raises on invalid input.
func isLoopback(s string) (bool, error) {
	addr, err := parseAddr("netaddr.is_loopback", s)
	if err != nil {
		return false, err
	}
	return addr.IsLoopback(), nil
}

// isGlobal: sugar for parse_ip(s).is_global. Raises on invalid input.
func isGlobal(s string) (bool, error) {
	addr, err := parseAddr("netaddr.is_global", s)
	if err != nil {
		return false, err
	}
	return isGlobalAddr(addr), nil
}

// parseCIDR: parses s as a CIDR prefix, normalising non-canonical input
// (host bits set) to its network address, and computes the first/last
// address and the exact address count as a decimal string (see the package
// doc for why num_addresses is a string). Raises on an unparseable prefix or
// an out-of-range prefix length.
func parseCIDR(s string) (CIDRInfo, error) {
	p, err := parsePrefix("netaddr.parse_cidr", s)
	if err != nil {
		return CIDRInfo{}, err
	}
	addr := p.Addr()
	is4 := addr.Is4()
	version, totalBits := 6, 128
	if is4 {
		version, totalBits = 4, 32
	}
	hostBits := totalBits - p.Bits()
	numAddresses := new(big.Int).Lsh(big.NewInt(1), uint(hostBits)) //nolint:gosec
	last := addr
	if hostBits > 0 {
		mask := new(big.Int).Sub(numAddresses, big.NewInt(1))
		lastInt := new(big.Int).Or(addrToBigInt(addr), mask)
		last = bigIntToAddr(lastInt, is4)
	}
	return CIDRInfo{
		CIDR:         p.String(),
		Network:      addr.String(),
		PrefixLen:    p.Bits(),
		Version:      version,
		First:        addr.String(),
		Last:         last.String(),
		NumAddresses: numAddresses.String(),
	}, nil
}

// cidrContains: reports whether ip falls inside cidr. A v4-mapped v6 ip is
// unwrapped first, so it is compared against a v4 cidr as its plain v4 form.
// After unwrapping, a genuine cross-family comparison (v4 host against a v6
// cidr, or vice versa) returns false rather than raising — only a malformed
// cidr or ip string raises.
func cidrContains(cidr, ip string) (bool, error) {
	p, err := parsePrefix("netaddr.cidr_contains", cidr)
	if err != nil {
		return false, err
	}
	addr, err := parseAddr("netaddr.cidr_contains", ip)
	if err != nil {
		return false, err
	}
	if addr.Is4() != p.Addr().Is4() {
		return false, nil
	}
	return p.Contains(addr), nil
}

// subnetOf: reports whether inner is fully contained within outer — i.e.
// outer's prefix is no more specific than inner's, and outer contains
// inner's network address. A prefix is a subnet of itself. Cross-family
// comparisons return false rather than raising; only a malformed cidr raises.
func subnetOf(inner, outer string) (bool, error) {
	innerP, err := parsePrefix("netaddr.subnet_of", inner)
	if err != nil {
		return false, err
	}
	outerP, err := parsePrefix("netaddr.subnet_of", outer)
	if err != nil {
		return false, err
	}
	if innerP.Addr().Is4() != outerP.Addr().Is4() {
		return false, nil
	}
	if outerP.Bits() > innerP.Bits() {
		return false, nil
	}
	return outerP.Contains(innerP.Addr()), nil
}

// cidrOverlaps: reports whether a and b share any address. Cross-family
// comparisons return false rather than raising; only a malformed cidr raises.
func cidrOverlaps(a, b string) (bool, error) {
	aP, err := parsePrefix("netaddr.cidr_overlaps", a)
	if err != nil {
		return false, err
	}
	bP, err := parsePrefix("netaddr.cidr_overlaps", b)
	if err != nil {
		return false, err
	}
	if aP.Addr().Is4() != bP.Addr().Is4() {
		return false, nil
	}
	return aP.Overlaps(bP), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("netaddr", "IP address and CIDR arithmetic, dual-stack (IPv4/IPv6), no name resolution")

	m.Fn("parse_ip", parseIP, "parses an IP address and reports version and policy-relevant predicates, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to parse; a v4-mapped v6 address (::ffff:192.0.2.1) is unwrapped to its plain v4 form"),
		luareg.ReturnDoc(0, "info", "IPInfo table: ip, version, is_private, is_loopback, is_global, is_multicast, is_unspecified, is_link_local"))
	m.Fn("is_ip", isIP, "reports whether s parses as a valid IP address; never raises",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to check"),
		luareg.ReturnDoc(0, "ok", "true if parse_ip(s) would succeed"))
	m.Fn("normalize_ip", normalizeIP, "returns the canonical string form of an address, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to normalize"),
		luareg.ReturnDoc(0, "out", "compressed IPv6 / zone stripped / v4-mapped v6 unwrapped to plain v4"))
	m.Fn("ip_version", ipVersion, "returns the IP version, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to check; a v4-mapped v6 address reports 4"),
		luareg.ReturnDoc(0, "version", "4 or 6"))
	m.Fn("is_private", isPrivate, "reports whether an address is private, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to check"),
		luareg.ReturnDoc(0, "ok", "true for RFC 1918 (v4), ULA fc00::/7, link-local, loopback or unspecified"))
	m.Fn("is_loopback", isLoopback, "reports whether an address is a loopback address, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to check"),
		luareg.ReturnDoc(0, "ok", "true for 127.0.0.0/8 or ::1"))
	m.Fn("is_global", isGlobal, "reports whether an address is globally routable, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the address to check"),
		luareg.ReturnDoc(0, "ok", "true if global unicast and not private/loopback/link-local/documentation/CGNAT"))
	m.Fn("parse_cidr", parseCIDR, "parses a CIDR prefix, raises on invalid input",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the prefix to parse, e.g. \"10.0.0.0/8\" or \"2001:db8::/32\"; a non-canonical input like \"10.0.0.5/8\" is normalized to its network address"),
		luareg.ReturnDoc(0, "info", "CIDRInfo table: cidr, network, prefix_len, version, first, last, num_addresses (a decimal string, exact even for a v6 /0)"))
	m.Fn("cidr_contains", cidrContains, "reports whether an IP falls inside a CIDR prefix, raises on a malformed argument",
		luareg.Args("cidr", "ip"),
		luareg.ArgDoc("cidr", "the containing prefix"),
		luareg.ArgDoc("ip", "the address to test; a v4-mapped v6 address is compared as its plain v4 form"),
		luareg.ReturnDoc(0, "ok", "false (not an error) if cidr and ip are valid but of different address families"))
	m.Fn("subnet_of", subnetOf, "reports whether inner is fully contained within outer, raises on a malformed argument",
		luareg.Args("inner", "outer"),
		luareg.ArgDoc("inner", "the candidate subnet"),
		luareg.ArgDoc("outer", "the candidate supernet; a prefix is a subnet of itself"),
		luareg.ReturnDoc(0, "ok", "false (not an error) if inner and outer are valid but of different address families"))
	m.Fn("cidr_overlaps", cidrOverlaps, "reports whether two CIDR prefixes share any address, raises on a malformed argument",
		luareg.Args("a", "b"),
		luareg.ArgDoc("a", "the first prefix"),
		luareg.ArgDoc("b", "the second prefix"),
		luareg.ReturnDoc(0, "ok", "false (not an error) if a and b are valid but of different address families"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("netaddr", netaddr.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
