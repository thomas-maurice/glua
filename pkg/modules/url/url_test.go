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

package url

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
			L.PreloadModule("url", Loader)
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

// TestParseIPv6LiteralHostSplitsHostnameAndPort: D5's headline url
// requirement -- "http://[::1]:8080/path" must split into hostname "::1"
// (unbracketed) and port "8080" separately from host (the raw, bracketed
// authority), not force callers to strip brackets themselves.
func TestParseIPv6LiteralHostSplitsHostnameAndPort(t *testing.T) {
	u, err := parseURL("http://[::1]:8080/path")
	require.NoError(t, err)
	assert.Equal(t, "::1", u.Hostname, "hostname must be unbracketed")
	assert.Equal(t, "8080", u.Port)
	assert.Equal(t, "[::1]:8080", u.Host, "host mirrors net/url.URL.Host verbatim, brackets included")
	assert.Equal(t, "/path", u.Path)
}

// TestBuildAddsBracketsAutomatically: a caller who builds a URL table from
// scratch, knowing only hostname+port (no host field), must not have to
// know the IPv6 bracket rule themselves -- build must add it.
func TestBuildAddsBracketsAutomatically(t *testing.T) {
	s, err := buildURL(URL{Scheme: "http", Hostname: "::1", Port: "8080", Path: "/path"})
	require.NoError(t, err)
	assert.Equal(t, "http://[::1]:8080/path", s)
}

// TestParseBuildRoundTripIPv6: build(parse(u)) must be lossless for an IPv6
// literal host URL.
func TestParseBuildRoundTripIPv6(t *testing.T) {
	const original = "http://[::1]:8080/path"
	parts, err := parseURL(original)
	require.NoError(t, err)
	rebuilt, err := buildURL(parts)
	require.NoError(t, err)
	assert.Equal(t, original, rebuilt)
}

// TestParseBuildRoundTripIPv6Zone: a zone id on a link-local IPv6 literal
// (RFC 6874, percent-encoded as %25 in the URL text) must survive a
// parse/build round-trip without corruption.
func TestParseBuildRoundTripIPv6Zone(t *testing.T) {
	const original = "http://[fe80::1%25eth0]/"
	parts, err := parseURL(original)
	require.NoError(t, err)
	// Hostname() is unbracketed AND percent-decodes the zone marker, so the
	// Go string holds a literal '%' rather than "%25" -- net/url's own
	// String() re-escapes it back to %25 when rendering from the Host
	// field, which is what makes the round-trip below lossless.
	assert.Equal(t, "fe80::1%eth0", parts.Hostname)
	rebuilt, err := buildURL(parts)
	require.NoError(t, err)
	assert.Equal(t, original, rebuilt, "zone id must round-trip without corruption")
}

// TestParseBuildRoundTripVariousURLs: parse . build must be lossless for a
// representative set of URLs, v4 and v6 alike.
func TestParseBuildRoundTripVariousURLs(t *testing.T) {
	cases := []string{
		"https://user:pw@example.com:8443/a/b?x=1&x=2#frag",
		"https://example.com/a/b",
		"http://example.com/",
		"http://[::1]:8080/path",
		"http://[2001:db8::1]/x?q=1",
	}
	for _, original := range cases {
		parts, err := parseURL(original)
		require.NoError(t, err, original)
		rebuilt, err := buildURL(parts)
		require.NoError(t, err, original)
		assert.Equal(t, original, rebuilt, "round-trip for %s", original)
	}
}

// TestParseOpaqueURL: a mailto: URL has no host/path, only an opaque part.
func TestParseOpaqueURL(t *testing.T) {
	u, err := parseURL("mailto:x@y.com")
	require.NoError(t, err)
	assert.Equal(t, "mailto", u.Scheme)
	assert.Equal(t, "x@y.com", u.Opaque)
	assert.Empty(t, u.Host)
}

// TestParseSchemeRelativeURL: a scheme-relative URL ("//host/path") has an
// empty scheme but a populated host.
func TestParseSchemeRelativeURL(t *testing.T) {
	u, err := parseURL("//example.com/a")
	require.NoError(t, err)
	assert.Empty(t, u.Scheme)
	assert.Equal(t, "example.com", u.Host)
	assert.Equal(t, "/a", u.Path)
}

// TestResolveTrailingSlashRule: the exact RFC 3986 behaviour the module
// promises -- a base without a trailing slash discards the last path
// segment, one with a trailing slash keeps it.
func TestResolveTrailingSlashRule(t *testing.T) {
	noSlash, err := resolve("https://example.com/a/b", "c")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/a/c", noSlash)

	withSlash, err := resolve("https://example.com/a/b/", "c")
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/a/b/c", withSlash)
}

// TestParseQueryEmptyStringReturnsEmptyNonNilMap: parse_query must return an
// empty table, never nil, for an empty query string.
func TestParseQueryEmptyStringReturnsEmptyNonNilMap(t *testing.T) {
	q, err := parseQuery("")
	require.NoError(t, err)
	require.NotNil(t, q)
	assert.Empty(t, q)
}

// TestParseQueryRepeatedKeyAndBareKey: a repeated key must keep every value
// (not silently drop one), and a bare key with no '=' maps to a one-element
// array containing the empty string.
func TestParseQueryRepeatedKeyAndBareKey(t *testing.T) {
	q, err := parseQuery("a=1&a=2&flag")
	require.NoError(t, err)
	assert.Equal(t, []string{"1", "2"}, q["a"])
	assert.Equal(t, []string{""}, q["flag"])
}

// TestBuildQueryDeterministic: build_query sorts by key (net/url.Values.Encode's
// own behaviour), so the same input produces the same output every time,
// including across multi-value keys expressed as an array.
func TestBuildQueryDeterministic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	tbl.RawSetString("ns", lua.LString("default"))
	labels := L.NewTable()
	labels.Append(lua.LString("a"))
	labels.Append(lua.LString("b"))
	tbl.RawSetString("label", labels)

	out1, err := buildQuery(tbl)
	require.NoError(t, err)
	out2, err := buildQuery(tbl)
	require.NoError(t, err)
	assert.Equal(t, out1, out2)
	assert.Equal(t, "label=a&label=b&ns=default", out1)
}

// TestBuildQueryRejectsInvalidValueShape: a value that is neither a string
// nor an array of strings raises, naming the offending key.
func TestBuildQueryRejectsInvalidValueShape(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	tbl.RawSetString("bad", lua.LNumber(1))

	_, err := buildQuery(tbl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "bad")
}

// TestBuildQueryRejectsNonStringKey: a non-string table key raises rather
// than being silently coerced or skipped.
func TestBuildQueryRejectsNonStringKey(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	tbl.RawSetInt(1, lua.LString("value-under-a-numeric-key"))

	_, err := buildQuery(tbl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "keys must be strings")
}

// TestBuildQueryRejectsNonStringArrayElement: an array value whose elements
// are not all strings raises, naming the key and the failing index.
func TestBuildQueryRejectsNonStringArrayElement(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tbl := L.NewTable()
	arr := L.NewTable()
	arr.Append(lua.LString("ok"))
	arr.Append(lua.LNumber(5))
	tbl.RawSetString("label", arr)

	_, err := buildQuery(tbl)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "label")
}

// TestBuildHostEmptyWhenNoHostOrHostname: an opaque URL (e.g. mailto:) has
// neither host nor hostname; build must produce an empty Host rather than a
// spurious "[]" or ":".
func TestBuildHostEmptyWhenNoHostOrHostname(t *testing.T) {
	s, err := buildURL(URL{Scheme: "mailto", Opaque: "x@y.com"})
	require.NoError(t, err)
	assert.Equal(t, "mailto:x@y.com", s)
}

// TestBuildRaisesOnInvalidResult: build must raise if the parts produce a
// string that does not itself parse back as a valid URL, rather than
// silently returning garbage.
func TestBuildRaisesOnInvalidResult(t *testing.T) {
	_, err := buildURL(URL{Scheme: "ht tp", Path: "/x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not valid")
}

// TestBuildHost_ConsistentHostAndHostname_RoundTrips: host present together
// with a hostname/port that agree with it must still work -- the fix must
// not break the legitimate case of rebuilding a table produced by parse.
func TestBuildHost_ConsistentHostAndHostname_RoundTrips(t *testing.T) {
	parts, err := parseURL("https://example.com:8443/a")
	require.NoError(t, err)
	require.Equal(t, "example.com:8443", parts.Host)
	require.Equal(t, "example.com", parts.Hostname)
	require.Equal(t, "8443", parts.Port)

	s, err := buildURL(parts)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com:8443/a", s)
}

// TestBuildHost_MutatedHostnameAfterParse_Raises is the exact regression
// this chunk fixes (security review LOW 1): mutating .hostname on a table
// obtained from parse, while leaving .host untouched, must raise -- not
// silently return the original (pre-mutation) host. Before the fix, this
// was a silent no-op: an SSRF allow-list written as
// `u.hostname = "good.com"; url.build(u)` looked like it worked but did
// nothing, because build always preferred host verbatim.
func TestBuildHost_MutatedHostnameAfterParse_Raises(t *testing.T) {
	parts, err := parseURL("https://evil.com/path")
	require.NoError(t, err)
	require.Equal(t, "evil.com", parts.Host)

	parts.Hostname = "good.com" // host is untouched, still "evil.com"

	_, err = buildURL(parts)
	require.Error(t, err, "build must not silently keep the original host when hostname was changed")
	assert.Contains(t, err.Error(), "evil.com")
	assert.Contains(t, err.Error(), "good.com")
}

// TestBuildHost_MutatedPortAfterParse_Raises: same as the hostname case,
// but for a port left inconsistent with host.
func TestBuildHost_MutatedPortAfterParse_Raises(t *testing.T) {
	parts, err := parseURL("https://example.com:443/path")
	require.NoError(t, err)

	parts.Port = "9999" // host is untouched, still "example.com:443"

	_, err = buildURL(parts)
	require.Error(t, err, "build must not silently keep the original host when port was changed")
	assert.Contains(t, err.Error(), "9999")
}

// TestBuildHost_HostAloneWorks: a caller who sets only host (no
// hostname/port at all) must not trip the new consistency check.
func TestBuildHost_HostAloneWorks(t *testing.T) {
	s, err := buildURL(URL{Scheme: "https", Host: "example.com:8080", Path: "/x"})
	require.NoError(t, err)
	assert.Equal(t, "https://example.com:8080/x", s)
}

// TestBuildHost_HostnameAloneWorks: a caller who sets only hostname/port
// (no host field at all) must not trip the new consistency check -- this is
// the pre-existing "build from scratch" path (TestBuildAddsBracketsAutomatically
// covers its IPv6 form).
func TestBuildHost_HostnameAloneWorks(t *testing.T) {
	s, err := buildURL(URL{Scheme: "https", Hostname: "example.com", Port: "8080", Path: "/x"})
	require.NoError(t, err)
	assert.Equal(t, "https://example.com:8080/x", s)
}

// TestBuildHost_ZoneIDRoundTrip_StillWorks: the consistency check must not
// false-positive on the case host-verbatim exists to handle: a zone-id IPv6
// literal, where re-composing hostname+port would not reproduce host
// byte-for-byte (host keeps '%25', Hostname() decodes it to a literal '%').
// This pins that the new check compares DECODED forms, not raw strings.
func TestBuildHost_ZoneIDRoundTrip_StillWorks(t *testing.T) {
	const original = "http://[fe80::1%25eth0]/"
	parts, err := parseURL(original)
	require.NoError(t, err)

	rebuilt, err := buildURL(parts)
	require.NoError(t, err)
	assert.Equal(t, original, rebuilt)
}

// TestBuildRaisesOnInvalidResult_PasswordNotLeaked is the exact regression
// this chunk fixes (security review LOW 2): the error raised when the
// constructed URL fails to re-parse must not embed the userinfo password
// verbatim -- that text is exactly the kind of string that ends up in
// log.error, a 500 body, or a webhook status.
func TestBuildRaisesOnInvalidResult_PasswordNotLeaked(t *testing.T) {
	_, err := buildURL(URL{
		Scheme:   "ht tp", // the space makes String() produce an unparseable result
		Username: "alice",
		Password: "s3cr3t",
		Host:     "example.com",
		Path:     "/x",
	})
	require.Error(t, err)
	assert.NotContains(t, err.Error(), "s3cr3t", "the raw password must never appear in a raised error")
	assert.Contains(t, err.Error(), "xxxxx", "the error should show the redacted form (net/url.URL.Redacted)")
}
