// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package x509

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	timemod "github.com/thomas-maurice/glua/pkg/modules/time"
	lua "github.com/yuin/gopher-lua"
)

// Test certificates in this file are generated fresh on every test run with
// crypto/x509.CreateCertificate, never checked in as static PEM fixtures —
// per SPECS.md S12's test plan, a checked-in certificate eventually expires
// and breaks the suite years later for a reason unrelated to the code being
// tested.

// genKey: generates a fresh P-256 key for a test certificate. ECDSA is used
// instead of RSA purely for test speed; key type is irrelevant to what this
// module tests (jwt is where key-type/algorithm-family handling is
// exercised).
func genKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	return key
}

// encodePEM: wraps DER bytes in a CERTIFICATE PEM block.
func encodePEM(der []byte) string {
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

// newCA: creates a self-signed CA certificate.
func newCA(t *testing.T, cn string, serial *big.Int, notBefore, notAfter time.Time) (pemStr string, cert *x509.Certificate, key *ecdsa.PrivateKey) {
	t.Helper()
	key = genKey(t)
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: cn, Organization: []string{"Test CA Org"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            1,
		MaxPathLenZero:        false,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	require.NoError(t, err)
	cert, err = x509.ParseCertificate(der)
	require.NoError(t, err)
	return encodePEM(der), cert, key
}

// newLeaf: creates a leaf certificate signed by parent/parentKey, populated
// with every SAN kind (dns, v4 IP, v6 IP, email, URI) plus several key
// usages, so a single fixture can drive the "parse populates every field"
// test.
func newLeaf(t *testing.T, cn string, parent *x509.Certificate, parentKey *ecdsa.PrivateKey, serial *big.Int, notBefore, notAfter time.Time) (pemStr string, cert *x509.Certificate, key *ecdsa.PrivateKey) {
	t.Helper()
	key = genKey(t)
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: cn, Organization: []string{"Leaf Org"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		DNSNames:              []string{cn, "www." + cn},
		IPAddresses:           []net.IP{net.ParseIP("192.0.2.10"), net.ParseIP("2001:db8::1")},
		EmailAddresses:        []string{"admin@example.com"},
		URIs:                  []*url.URL{{Scheme: "spiffe", Host: "example.com", Path: "/leaf"}},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, parent, &key.PublicKey, parentKey)
	require.NoError(t, err)
	cert, err = x509.ParseCertificate(der)
	require.NoError(t, err)
	return encodePEM(der), cert, key
}

// bigSerial: a serial number well past float64's exact integer range
// (2^53), used to pin that x509.parse renders it as an exact decimal
// string rather than a rounded Lua number.
func bigSerial() *big.Int {
	s, ok := new(big.Int).SetString("123456789012345678901234567890", 10)
	if !ok {
		panic("bad test serial")
	}
	return s
}

// TestPublicKeyBits_KeyTypes: pins publicKeyBits' behaviour for the three
// key types it recognises, plus its documented fallback for anything else.
func TestPublicKeyBits_KeyTypes(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	assert.Equal(t, 2048, publicKeyBits(&rsaKey.PublicKey))

	ecKey := genKey(t)
	assert.Equal(t, 256, publicKeyBits(&ecKey.PublicKey))

	edPub, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	assert.Equal(t, 256, publicKeyBits(edPub))

	assert.Equal(t, 0, publicKeyBits("not a key"), "unknown key types must report 0, not raise")
}

// TestExtKeyUsageStrings_UnknownValue: an ExtKeyUsage outside the known set
// must render as "unknown_<n>" rather than being silently dropped -- an
// operator auditing a certificate's usages must see it, even unlabelled.
func TestExtKeyUsageStrings_UnknownValue(t *testing.T) {
	got := extKeyUsageStrings([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsage(999)})
	assert.Contains(t, got, "server_auth")
	assert.Contains(t, got, "unknown_999")
}

func TestParse_AllFields(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour).Truncate(time.Second)
	notAfter := notBefore.Add(24 * time.Hour)
	_, ca, caKey := newCA(t, "Test Root CA", big.NewInt(1), notBefore, notAfter)
	pemStr, _, _ := newLeaf(t, "leaf.example.com", ca, caKey, bigSerial(), notBefore, notAfter)

	c, err := parseFn(pemStr)
	require.NoError(t, err)

	assert.Equal(t, "leaf.example.com", c.SubjectCN)
	assert.Equal(t, "Test Root CA", c.IssuerCN)
	assert.Equal(t, []string{"Leaf Org"}, c.SubjectOrg)
	assert.Equal(t, "123456789012345678901234567890", c.Serial, "serial must be an exact decimal string, not a rounded float64")
	assert.Equal(t, notBefore.Unix(), c.NotBefore)
	assert.Equal(t, notAfter.Unix(), c.NotAfter)
	assert.ElementsMatch(t, []string{"leaf.example.com", "www.leaf.example.com"}, c.DNSNames)
	assert.Contains(t, c.IPAddresses, "192.0.2.10", "IPv4 SAN must be present")
	assert.Contains(t, c.IPAddresses, "2001:db8::1", "IPv6 SAN must be present and correctly formatted (D5)")
	assert.Equal(t, []string{"admin@example.com"}, c.EmailAddresses)
	assert.Equal(t, []string{"spiffe://example.com/leaf"}, c.URIs)
	assert.False(t, c.IsCA)
	assert.ElementsMatch(t, []string{"digital_signature", "key_encipherment"}, c.KeyUsage)
	assert.ElementsMatch(t, []string{"server_auth", "client_auth"}, c.ExtKeyUsage)
	assert.Equal(t, 256, c.PublicKeyBits, "P-256 key must report 256 bits")
	assert.Len(t, c.FingerprintSHA256, 64, "sha256 fingerprint must be 64 lowercase hex chars")
	assert.Regexp(t, "^[0-9a-f]{64}$", c.FingerprintSHA256)
}

func TestParse_CAFields(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour)
	notAfter := notBefore.Add(24 * time.Hour)
	pemStr, _, _ := newCA(t, "Test Root CA", big.NewInt(1), notBefore, notAfter)

	c, err := parseFn(pemStr)
	require.NoError(t, err)

	assert.True(t, c.IsCA)
	assert.Equal(t, 1, c.MaxPathLen)
	assert.ElementsMatch(t, []string{"digital_signature", "cert_sign", "crl_sign"}, c.KeyUsage)
}

// TestParse_EmptyIPAddresses_ReturnsEmptyTable: a certificate with no IP SANs
// must report an empty (non-nil) table, not Lua nil -- see the skill's
// nil-slice rule.
func TestParse_EmptyIPAddresses_ReturnsEmptyTable(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour)
	notAfter := notBefore.Add(24 * time.Hour)
	pemStr, _, _ := newCA(t, "Test Root CA", big.NewInt(1), notBefore, notAfter)

	c, err := parseFn(pemStr)
	require.NoError(t, err)
	assert.NotNil(t, c.IPAddresses)
	assert.Empty(t, c.IPAddresses)
	assert.NotNil(t, c.DNSNames)
	assert.NotNil(t, c.EmailAddresses)
	assert.NotNil(t, c.URIs)
}

func TestParse_EmptyInput_Raises(t *testing.T) {
	_, err := parseFn("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no PEM block found")
}

func TestParse_PrivateKeyBlock_Raises(t *testing.T) {
	key := genKey(t)
	der, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	pemStr := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}))

	_, err = parseFn(pemStr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not CERTIFICATE")
}

func TestParse_TruncatedDER_Raises(t *testing.T) {
	pemStr := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("not valid der")}))
	_, err := parseFn(pemStr)
	require.Error(t, err)
}

func TestParseChain_PreservesOrder(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour)
	notAfter := notBefore.Add(24 * time.Hour)
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), notBefore, notAfter)
	interPEM, inter, interKey := newCA(t, "Intermediate", big.NewInt(2), notBefore, notAfter)
	_ = root
	leafPEM, _, _ := newLeaf(t, "leaf.example.com", inter, interKey, big.NewInt(3), notBefore, notAfter)

	bundle := rootPEM + interPEM + leafPEM
	certs, err := parseChainFn(bundle)
	require.NoError(t, err)
	require.Len(t, certs, 3)
	assert.Equal(t, "Root", certs[0].SubjectCN)
	assert.Equal(t, "Intermediate", certs[1].SubjectCN)
	assert.Equal(t, "leaf.example.com", certs[2].SubjectCN)
	_ = rootKey
}

func TestParseChain_EmptyInput_Raises(t *testing.T) {
	_, err := parseChainFn("")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no CERTIFICATE blocks found")
}

func TestParseChain_BadBlock_Raises(t *testing.T) {
	notBefore := time.Now().Add(-time.Hour)
	notAfter := notBefore.Add(24 * time.Hour)
	good, _, _ := newCA(t, "Root", big.NewInt(1), notBefore, notAfter)
	bad := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte("garbage")}))

	_, err := parseChainFn(good + bad)
	require.Error(t, err)
}

func TestExpiresInDays_PositiveNegativeFractional(t *testing.T) {
	now := time.Now().Truncate(time.Second)

	// Not yet expired: 10 days remaining.
	pemStr, _, _ := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(10*24*time.Hour))
	days, err := expiresInDaysFn(pemStr, now.Unix())
	require.NoError(t, err)
	assert.InDelta(t, 10.0, days, 0.01)

	// Already expired: negative.
	pemStr2, _, _ := newCA(t, "Root2", big.NewInt(2), now.Add(-48*time.Hour), now.Add(-24*time.Hour))
	days2, err := expiresInDaysFn(pemStr2, now.Unix())
	require.NoError(t, err)
	assert.Less(t, days2, 0.0)

	// Fractional: 12 hours remaining -> 0.5 days.
	pemStr3, _, _ := newCA(t, "Root3", big.NewInt(3), now.Add(-time.Hour), now.Add(12*time.Hour))
	days3, err := expiresInDaysFn(pemStr3, now.Unix())
	require.NoError(t, err)
	assert.InDelta(t, 0.5, days3, 0.01)
}

func TestExpiresInDays_BadPEM_Raises(t *testing.T) {
	_, err := expiresInDaysFn("garbage", time.Now().Unix())
	require.Error(t, err)
}

func TestIsValidAt_Boundaries(t *testing.T) {
	notBefore := time.Now().Truncate(time.Second)
	notAfter := notBefore.Add(time.Hour)
	pemStr, _, _ := newCA(t, "Root", big.NewInt(1), notBefore, notAfter)

	okAtStart, err := isValidAtFn(pemStr, notBefore.Unix())
	require.NoError(t, err)
	assert.True(t, okAtStart, "exactly not_before must be valid (inclusive)")

	okAtEnd, err := isValidAtFn(pemStr, notAfter.Unix())
	require.NoError(t, err)
	assert.True(t, okAtEnd, "exactly not_after must be valid (inclusive)")

	okBefore, err := isValidAtFn(pemStr, notBefore.Add(-time.Second).Unix())
	require.NoError(t, err)
	assert.False(t, okBefore)

	okAfter, err := isValidAtFn(pemStr, notAfter.Add(time.Second).Unix())
	require.NoError(t, err)
	assert.False(t, okAfter)
}

func TestIsValidAt_BadPEM_Raises(t *testing.T) {
	_, err := isValidAtFn("garbage", time.Now().Unix())
	require.Error(t, err)
}

func TestVerifyChain_Success(t *testing.T) {
	now := time.Now()
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))
	leafPEM, _, _ := newLeaf(t, "api.example.com", root, rootKey, big.NewInt(2), now.Add(-time.Hour), now.Add(24*time.Hour))

	ok, reason, err := verifyChainFn(leafPEM, "", rootPEM, VerifyChainOptions{
		DNSName:   "api.example.com",
		AtTime:    now.Unix(),
		KeyUsages: []string{"server_auth"},
	})
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Empty(t, reason)
}

func TestVerifyChain_DefaultKeyUsageIsServerAuth(t *testing.T) {
	now := time.Now()
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))
	leafPEM, _, _ := newLeaf(t, "api.example.com", root, rootKey, big.NewInt(2), now.Add(-time.Hour), now.Add(24*time.Hour))

	ok, reason, err := verifyChainFn(leafPEM, "", rootPEM, VerifyChainOptions{
		DNSName: "api.example.com",
		AtTime:  now.Unix(),
		// KeyUsages intentionally omitted: must default to {"server_auth"}.
	})
	require.NoError(t, err)
	assert.True(t, ok, "empty key_usages must default to server_auth: %s", reason)
}

func TestVerifyChain_ExpiredFails(t *testing.T) {
	now := time.Now()
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), now.Add(-48*time.Hour), now.Add(48*time.Hour))
	leafPEM, _, _ := newLeaf(t, "api.example.com", root, rootKey, big.NewInt(2), now.Add(-48*time.Hour), now.Add(-24*time.Hour))

	ok, reason, err := verifyChainFn(leafPEM, "", rootPEM, VerifyChainOptions{
		DNSName:   "api.example.com",
		AtTime:    now.Unix(), // past leaf's not_after
		KeyUsages: []string{"server_auth"},
	})
	require.NoError(t, err)
	assert.False(t, ok)
	assert.NotEmpty(t, reason, "reason must explain why verification failed")
}

func TestVerifyChain_WrongDNSNameFails(t *testing.T) {
	now := time.Now()
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))
	leafPEM, _, _ := newLeaf(t, "api.example.com", root, rootKey, big.NewInt(2), now.Add(-time.Hour), now.Add(24*time.Hour))

	ok, reason, err := verifyChainFn(leafPEM, "", rootPEM, VerifyChainOptions{
		DNSName:   "not-the-right-host.example.com",
		AtTime:    now.Unix(),
		KeyUsages: []string{"server_auth"},
	})
	require.NoError(t, err)
	assert.False(t, ok)
	assert.NotEmpty(t, reason)
}

func TestVerifyChain_MissingIntermediateFails(t *testing.T) {
	now := time.Now()
	rootPEM, root, rootKey := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))
	_, inter, interKey := newCA(t, "Intermediate", big.NewInt(2), now.Add(-time.Hour), now.Add(24*time.Hour))
	// inter is signed by root, but we never present rootPEM as its parent
	// authority in a chain -- we sign the leaf under inter and only supply
	// root as the trust anchor, withholding the intermediate PEM.
	interPEM, interCert, _ := newLeaf(t, "intermediate-cert", root, rootKey, big.NewInt(3), now.Add(-time.Hour), now.Add(24*time.Hour))
	_ = interPEM
	_ = interCert
	leafPEM, _, _ := newLeaf(t, "api.example.com", inter, interKey, big.NewInt(4), now.Add(-time.Hour), now.Add(24*time.Hour))

	ok, reason, err := verifyChainFn(leafPEM, "", rootPEM, VerifyChainOptions{
		DNSName:   "api.example.com",
		AtTime:    now.Unix(),
		KeyUsages: []string{"server_auth"},
	})
	require.NoError(t, err)
	assert.False(t, ok, "verification must fail without the intermediate that signed the leaf")
	assert.NotEmpty(t, reason)
}

// TestVerifyChain_EmptyRootsNeverSucceeds is the invariant test that matters
// most for this module: even a self-signed certificate presented as its own
// leaf must NEVER verify against an empty roots argument. If this test ever
// fails, x509.SystemCertPool has been reintroduced somewhere in
// verifyChainFn, silently making verification depend on the host's trust
// store instead of only the caller-supplied roots.
func TestVerifyChain_EmptyRootsNeverSucceeds(t *testing.T) {
	now := time.Now()
	// A self-signed cert that is also a CA -- the most favourable case for
	// an accidental SystemCertPool fallback to "succeed" wrongly, since a
	// self-signed cert would trivially verify against itself if it were
	// ever (incorrectly) added as its own root.
	pemStr, _, _ := newCA(t, "Self Signed", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))

	ok, reason, err := verifyChainFn(pemStr, "", "", VerifyChainOptions{
		AtTime:    now.Unix(),
		KeyUsages: []string{"server_auth"},
	})
	require.NoError(t, err)
	assert.False(t, ok, "verify_chain must never succeed with an empty roots argument")
	assert.NotEmpty(t, reason)
}

func TestVerifyChain_AtTimeZero_Raises(t *testing.T) {
	now := time.Now()
	pemStr, _, _ := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))

	_, _, err := verifyChainFn(pemStr, "", pemStr, VerifyChainOptions{AtTime: 0})
	require.Error(t, err, "at_time == 0 must raise rather than silently defaulting to 1970")
	assert.Contains(t, err.Error(), "at_time")
}

func TestVerifyChain_BadLeafPEM_Raises(t *testing.T) {
	_, _, err := verifyChainFn("garbage", "", "", VerifyChainOptions{AtTime: time.Now().Unix()})
	require.Error(t, err)
}

func TestVerifyChain_UnknownKeyUsage_Raises(t *testing.T) {
	now := time.Now()
	pemStr, _, _ := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))

	_, _, err := verifyChainFn(pemStr, "", pemStr, VerifyChainOptions{
		AtTime:    now.Unix(),
		KeyUsages: []string{"not_a_real_usage"},
	})
	require.Error(t, err)
}

// TestNoSystemCertPool: pins the package's central invariant at the source
// level, mirroring netaddr's "no net import" test. crypto/x509.SystemCertPool
// is a specific, greppable identifier, and a text search over the non-test
// source files is a simple, robust way to make "never calls it" a property
// a reviewer (and CI) can check mechanically rather than trusting the
// package doc alone.
func TestNoSystemCertPool(t *testing.T) {
	files, err := filepath.Glob("*.go")
	require.NoError(t, err)
	require.NotEmpty(t, files)

	checked := 0
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		checked++
		data, err := os.ReadFile(file)
		require.NoError(t, err)
		// Match the call form specifically (not the package doc's prose
		// mentions of the identifier) so this test doesn't trip on its own
		// documentation.
		require.NotContains(t, string(data), "SystemCertPool(",
			"%s must never call x509.SystemCertPool -- verification must only trust caller-supplied roots", file)
	}
	require.Positive(t, checked, "no non-test .go files were found to check")
}

// TestLuaScripts: runs every Lua test script in testdata/. These scripts
// deliberately do NOT embed a generated certificate (a static fixture would
// eventually expire); they cover pcall behaviour on malformed PEM only. The
// composition test with the time module (proving the Unix-second choice
// actually composes) lives in TestComposesWithTimeModule below, where a
// fresh certificate can be generated per run.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("x509", Loader)
			if err := L.DoFile(file); err != nil {
				t.Fatalf("Lua script failed: %v", err)
			}
			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("test script returned %v, expected true", result)
			}
		})
	}
}

// TestComposesWithTimeModule: proves x509's Unix-second timestamps compose
// with the time module, as the package doc promises. The certificate is
// generated fresh for this run and injected into the script as a Lua
// long-bracket string literal, so nothing here can go stale.
func TestComposesWithTimeModule(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	pemStr, _, _ := newCA(t, "Root", big.NewInt(1), now.Add(-time.Hour), now.Add(24*time.Hour))

	script := fmt.Sprintf(`
local x509 = require("x509")
local time = require("time")

local c = x509.parse([==[%s]==])
assert(type(c.not_after) == "number", "not_after must be a number")

-- Composing with the time module: format() takes the same Unix-second value.
local formatted = time.format(c.not_after, "2006-01-02")
assert(type(formatted) == "string", "time.format must accept x509's not_after")

local days = x509.expires_in_days([==[%s]==], time.now())
assert(type(days) == "number", "expires_in_days must accept time.now()")
assert(days > 0, "certificate should not be expired yet")

return true
`, pemStr, pemStr)

	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("x509", Loader)
	L.PreloadModule("time", timemod.Loader)
	require.NoError(t, L.DoString(script))
	assert.Equal(t, lua.LTrue, L.Get(-1))
}
