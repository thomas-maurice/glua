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

package netaddr

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// TestLuaScripts: runs all Lua test scripts in testdata/.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "No Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("netaddr", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("Test script returned %v, expected true", result)
			}
		})
	}
}

// TestParseIPUnwrapsV4MappedV6: pins the module-wide decision that a
// v4-mapped v6 address (::ffff:192.0.2.1) is unwrapped to its plain v4 form
// everywhere, not just in normalize_ip -- version, ip and every predicate
// reflect the v4 address, so it behaves identically to typing "192.0.2.1".
func TestParseIPUnwrapsV4MappedV6(t *testing.T) {
	info, err := parseIP("::ffff:192.0.2.1")
	require.NoError(t, err)
	assert.Equal(t, "192.0.2.1", info.IP)
	assert.Equal(t, 4, info.Version)

	plain, err := parseIP("192.0.2.1")
	require.NoError(t, err)
	assert.Equal(t, plain, info, "a v4-mapped v6 address must parse identically to its plain v4 form")
}

// TestCIDRContainsV4MappedV6AgainstV4CIDR: pins the decision that a
// v4-mapped v6 host address IS considered contained in a v4 CIDR, because it
// is unwrapped to its v4 form before the family check runs.
func TestCIDRContainsV4MappedV6AgainstV4CIDR(t *testing.T) {
	ok, err := cidrContains("10.0.0.0/8", "::ffff:10.1.2.3")
	require.NoError(t, err)
	assert.True(t, ok, "a v4-mapped v6 host must be compared as its unwrapped v4 form")
}

// TestCIDRContainsCrossFamilyReturnsFalse: D5 -- a v4 host against a v6 CIDR
// (and vice versa) must return false, not raise, once both sides parse.
func TestCIDRContainsCrossFamilyReturnsFalse(t *testing.T) {
	ok, err := cidrContains("2001:db8::/32", "10.0.0.1")
	require.NoError(t, err)
	assert.False(t, ok)

	ok, err = cidrContains("10.0.0.0/8", "2001:db8::1")
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestSubnetOfAndOverlapsCrossFamilyReturnFalse: same cross-family contract
// as cidr_contains, for subnet_of and cidr_overlaps.
func TestSubnetOfAndOverlapsCrossFamilyReturnFalse(t *testing.T) {
	ok, err := subnetOf("10.0.0.0/24", "2001:db8::/32")
	require.NoError(t, err)
	assert.False(t, ok)

	ok, err = cidrOverlaps("10.0.0.0/24", "::/0")
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestParseCIDRNonCanonicalIsNormalized: pins the decision that a
// non-canonical CIDR (host bits set) is normalized to its network address
// rather than rejected, matching net.ParseCIDR's own (perhaps surprising)
// precedent of returning the masked network.
func TestParseCIDRNonCanonicalIsNormalized(t *testing.T) {
	info, err := parseCIDR("10.0.0.5/8")
	require.NoError(t, err)
	assert.Equal(t, "10.0.0.0/8", info.CIDR)
	assert.Equal(t, "10.0.0.0", info.Network)
}

// TestParseCIDRV6ZeroSlashZero: "::/0" is the entire IPv6 address space --
// 2^128 addresses, exercised as an exact decimal string (39 digits), and its
// first/last addresses are "::" and the all-ones address.
func TestParseCIDRV6ZeroSlashZero(t *testing.T) {
	info, err := parseCIDR("::/0")
	require.NoError(t, err)
	assert.Equal(t, 6, info.Version)
	assert.Equal(t, "::", info.First)
	assert.Equal(t, "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff", info.Last)
	assert.Equal(t, "340282366920938463463374607431768211456", info.NumAddresses)
	assert.Len(t, info.NumAddresses, 39)
}

// TestParseCIDRV6Slash64: a v6 /64 holds 2^64 addresses -- past float64's
// exact integer range (2^53) -- so num_addresses must be an exact string,
// not a rounded Lua number.
func TestParseCIDRV6Slash64(t *testing.T) {
	info, err := parseCIDR("2001:db8::/64")
	require.NoError(t, err)
	assert.Equal(t, "18446744073709551616", info.NumAddresses)
	assert.Equal(t, "2001:db8::", info.First)
	assert.Equal(t, "2001:db8::ffff:ffff:ffff:ffff", info.Last)
}

// TestParseCIDRV4Slash24: the v4 ergonomic cost of F4's "always a string"
// rule -- a /24 fits comfortably in a float64, but callers still get a
// string and must tonumber() it themselves.
func TestParseCIDRV4Slash24(t *testing.T) {
	info, err := parseCIDR("192.168.1.0/24")
	require.NoError(t, err)
	assert.Equal(t, "256", info.NumAddresses)
	assert.Equal(t, "192.168.1.0", info.First)
	assert.Equal(t, "192.168.1.255", info.Last)
}

// TestPrefixLengthErrorNamesTheBound: prefix validation must name the actual
// bound (0-32 for v4, 0-128 for v6) rather than Go's bare "prefix length out
// of range", per D5's "validated with useful errors".
func TestPrefixLengthErrorNamesTheBound(t *testing.T) {
	_, err := parseCIDR("10.0.0.0/33")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "between 0 and 32")

	_, err = parseCIDR("2001:db8::/129")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "between 0 and 128")
}

// TestIsPrivateCoversV6Ranges: D5 -- is_private must cover ULA, link-local,
// loopback and unspecified for v6, not only the v4 RFC 1918 set.
func TestIsPrivateCoversV6Ranges(t *testing.T) {
	cases := map[string]bool{
		"10.1.2.3":    true,  // RFC 1918
		"172.16.0.1":  true,  // RFC 1918
		"192.168.1.1": true,  // RFC 1918
		"8.8.8.8":     false, // public v4
		"fc00::1":     true,  // ULA fc00::/7
		"fdff::1":     true,  // ULA fc00::/7 (upper half)
		"fe80::1":     true,  // link-local
		"::1":         true,  // loopback
		"::":          true,  // unspecified
		"2001:db8::1": false, // documentation range, not private
	}
	for addr, want := range cases {
		got, err := isPrivate(addr)
		require.NoError(t, err, addr)
		assert.Equal(t, want, got, addr)
	}
}

// TestIsGlobalCoversV6Ranges: D5 -- is_global must exclude v6 documentation
// space in addition to the v4 exclusions, and admit a genuine public v6
// address.
func TestIsGlobalCoversV6Ranges(t *testing.T) {
	cases := map[string]bool{
		"8.8.8.8":      true,  // public v4
		"10.0.0.1":     false, // private v4
		"100.64.0.1":   false, // CGNAT
		"192.0.2.1":    false, // TEST-NET-1
		"2001:4860::1": true,  // public v6 (Google range, not a documented exclusion)
		"2001:db8::1":  false, // v6 documentation range
		"100::1":       false, // v6 discard-only range
		"fc00::1":      false, // ULA
		"fe80::1":      false, // link-local
		"::1":          false, // loopback
	}
	for addr, want := range cases {
		got, err := isGlobal(addr)
		require.NoError(t, err, addr)
		assert.Equal(t, want, got, addr)
	}
}

// TestParseIPZoneBearing: netip accepts a zone identifier on a link-local
// address; parse_ip must not reject it, and the zone survives in the ip
// field (normalize_ip is the one that strips it -- see the Lua test suite).
func TestParseIPZoneBearing(t *testing.T) {
	info, err := parseIP("fe80::1%eth0")
	require.NoError(t, err)
	assert.Equal(t, "fe80::1%eth0", info.IP)
	assert.True(t, info.IsLinkLocal)
}

// TestSubnetOfSamePrefixIsTrue: a prefix is a subnet of itself.
func TestSubnetOfSamePrefixIsTrue(t *testing.T) {
	ok, err := subnetOf("10.0.0.0/24", "10.0.0.0/24")
	require.NoError(t, err)
	assert.True(t, ok)
}

// TestSubnetOfAdjacentPrefixIsFalse: adjacent, same-size prefixes are not
// subnets of one another.
func TestSubnetOfAdjacentPrefixIsFalse(t *testing.T) {
	ok, err := subnetOf("10.0.1.0/24", "10.0.0.0/24")
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestCIDROverlapsAdjacentIsFalse: adjacent, non-overlapping prefixes do not overlap.
func TestCIDROverlapsAdjacentIsFalse(t *testing.T) {
	ok, err := cidrOverlaps("10.0.0.0/24", "10.0.1.0/24")
	require.NoError(t, err)
	assert.False(t, ok)
}

// TestIsIPNeverRaises: is_ip is the non-raising companion to parse_ip.
func TestIsIPNeverRaises(t *testing.T) {
	assert.True(t, isIP("192.168.1.1"))
	assert.True(t, isIP("::1"))
	assert.False(t, isIP("not an ip"))
	assert.False(t, isIP("999.999.999.999"))
}
