// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package jwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

// jwtNewWithClaimsHS256: signs claims with HS256 using golang-jwt directly,
// bypassing this package's sign() entirely. Used to simulate an attacker
// forging a token, not a legitimate caller of this module.
func jwtNewWithClaimsHS256(claims map[string]interface{}, secret []byte) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return tok.SignedString(secret)
}

// rawUnsignedNoneToken: hand-builds a compact JWT whose header declares
// alg: none with an empty signature segment -- the classic bypass payload,
// constructed without going through this package's sign() (which can never
// produce such a token, by design).
func rawUnsignedNoneToken(t *testing.T, claims map[string]interface{}) string {
	t.Helper()
	header := map[string]interface{}{"alg": "none", "typ": "JWT"}
	hb, err := json.Marshal(header)
	require.NoError(t, err)
	cb, err := json.Marshal(claims)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(hb) + "." + base64.RawURLEncoding.EncodeToString(cb) + "."
}

// Keys in this file are generated fresh per test run, never checked in --
// same reasoning as pkg/modules/x509's test fixtures.

func genRSAKeyPEMs(t *testing.T) (privPEM, pubPEM string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privDER := x509.MarshalPKCS1PrivateKey(key)
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privDER}))

	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return privPEM, pubPEM
}

func genECKeyPEMs(t *testing.T) (privPEM, pubPEM string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	privDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privDER}))

	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	require.NoError(t, err)
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return privPEM, pubPEM
}

// genEd25519KeyPEMs: PKCS#8 (private) / PKIX (public) PEM, the only
// encodings this package accepts for EdDSA -- see the package doc's
// "Accepted key encodings" section.
func genEd25519KeyPEMs(t *testing.T) (privPEM, pubPEM string) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	privDER, err := x509.MarshalPKCS8PrivateKey(priv)
	require.NoError(t, err)
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privDER}))

	pubDER, err := x509.MarshalPKIXPublicKey(pub)
	require.NoError(t, err)
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return privPEM, pubPEM
}

// fakeEd448PKCS8PEM: hand-builds a syntactically valid PKCS#8 structure
// declaring the Ed448 OID (1.3.101.113, RFC 8410) as its algorithm. Go's
// standard library implements no Ed448 support at all (crypto/ed25519 is
// Ed25519-only, and crypto/x509 has no Ed448 case), so an Ed448 key cannot
// be generated with stdlib -- this builds only the ASN.1 shape needed to
// prove parseSignKey/parseVerifyKey reject the OID cleanly, not that
// signing/verification with Ed448 works (which stdlib cannot do at all).
func fakeEd448PKCS8PEM(t *testing.T) string {
	t.Helper()
	// Mirrors the private crypto/x509 "pkcs8" struct shape (version, algo,
	// opaque private key octet string) closely enough for asn1.Marshal to
	// produce something ParsePKCS8PrivateKey's outer unmarshal accepts; the
	// actual key bytes are never read because the OID switch rejects the
	// algorithm before getting that far.
	type pkcs8 struct {
		Version    int
		Algo       pkix.AlgorithmIdentifier
		PrivateKey []byte
	}
	der, err := asn1.Marshal(pkcs8{
		Version: 0,
		Algo:    pkix.AlgorithmIdentifier{Algorithm: asn1.ObjectIdentifier{1, 3, 101, 113}},
		PrivateKey: func() []byte {
			b, err := asn1.Marshal(make([]byte, 57)) // Ed448 seed size, irrelevant to the outcome
			require.NoError(t, err)
			return b
		}(),
	})
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

func TestSignVerify_RS256_RoundTrip(t *testing.T) {
	privPEM, pubPEM := genRSAKeyPEMs(t)

	token, err := signFn(map[string]interface{}{"sub": "svc-a", "exp": float64(time.Now().Add(time.Hour).Unix())}, privPEM, SignOptions{Algorithm: "RS256"})
	require.NoError(t, err)

	claims, err := verifyFn(token, pubPEM, VerifyOptions{Algorithms: []string{"RS256"}})
	require.NoError(t, err)
	assert.Equal(t, "svc-a", claims["sub"])
}

func TestSignVerify_ES256_RoundTrip(t *testing.T) {
	privPEM, pubPEM := genECKeyPEMs(t)

	token, err := signFn(map[string]interface{}{"sub": "svc-b", "exp": float64(time.Now().Add(time.Hour).Unix())}, privPEM, SignOptions{Algorithm: "ES256"})
	require.NoError(t, err)

	claims, err := verifyFn(token, pubPEM, VerifyOptions{Algorithms: []string{"ES256"}})
	require.NoError(t, err)
	assert.Equal(t, "svc-b", claims["sub"])
}

func TestSignVerify_HS256_RoundTrip(t *testing.T) {
	token, err := signFn(map[string]interface{}{"sub": "svc-c", "exp": float64(time.Now().Add(time.Hour).Unix())}, "secret", SignOptions{Algorithm: "HS256"})
	require.NoError(t, err)

	claims, err := verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}})
	require.NoError(t, err)
	assert.Equal(t, "svc-c", claims["sub"])
}

func TestSignVerify_EdDSA_RoundTrip(t *testing.T) {
	privPEM, pubPEM := genEd25519KeyPEMs(t)

	token, err := signFn(map[string]interface{}{"sub": "svc-d", "exp": float64(time.Now().Add(time.Hour).Unix())}, privPEM, SignOptions{Algorithm: "EdDSA"})
	require.NoError(t, err)

	claims, err := verifyFn(token, pubPEM, VerifyOptions{Algorithms: []string{"EdDSA"}})
	require.NoError(t, err)
	assert.Equal(t, "svc-d", claims["sub"])
}

// TestAlgorithmConfusion_ForgedHS256UsingRSAPublicKey is the single most
// important test in this chunk (SPECS.md S13): an attacker who only has the
// server's RSA PUBLIC key must not be able to forge a token by signing it
// with HS256 using that public key's PEM bytes as the HMAC secret.
func TestAlgorithmConfusion_ForgedHS256UsingRSAPublicKey(t *testing.T) {
	_, pubPEM := genRSAKeyPEMs(t)

	// The attacker signs an HS256 token using the RSA public key PEM as the
	// HMAC secret -- this is standard library usage, deliberately bypassing
	// this package's sign() to simulate an attacker who doesn't need our
	// API to forge a token.
	forged, err := jwtNewWithClaimsHS256(map[string]interface{}{
		"sub": "admin",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}, []byte(pubPEM))
	require.NoError(t, err)

	// The server only ever expects RS256 and hands verify the same public
	// key it always uses for RS256 verification.
	_, err = verifyFn(forged, pubPEM, VerifyOptions{Algorithms: []string{"RS256"}})
	require.Error(t, err, "forged HS256 token must be rejected when the server only accepts RS256")

	// Assert on the SPECIFIC error jwt/v5's WithValidMethods produces
	// (parser.go: newError(fmt.Sprintf("signing method %v is invalid", alg),
	// ErrTokenSignatureInvalid)), not just "an error occurred". Both the
	// protected path (WithValidMethods rejects the header's declared alg
	// before any key is used) and jwt/v5's later HMAC-key-type-mismatch path
	// wrap the SAME ErrTokenSignatureInvalid sentinel, so errors.Is alone
	// cannot tell them apart -- confirmed by reproducing both call paths
	// directly against golang-jwt/jwt/v5 v5.3.1:
	//   with WithValidMethods:    "token signature is invalid: signing method HS256 is invalid"
	//   without WithValidMethods: "token signature is invalid: key is of invalid type: HMAC verify expects []byte"
	// Only the first message proves the algorithm-family defence fired; the
	// second would still occur (with a different message) if
	// WithValidMethods were ever removed from verifyFn, which is exactly
	// the regression this test must catch (repo Rule 9: the previous
	// assertion, require.Error only, passed whether or not the defence
	// existed, because the RSA-public-key-as-HMAC-secret type mismatch
	// alone was enough to fail verification).
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	assert.Contains(t, err.Error(), "signing method HS256 is invalid",
		"error must show WithValidMethods rejected the token's declared alg, not merely that some later step failed")
}

// TestAlgorithmConfusion_ForgedHS256UsingEd25519PublicKey mirrors
// TestAlgorithmConfusion_ForgedHS256UsingRSAPublicKey for the new EdDSA
// family: an attacker who only has the server's Ed25519 PUBLIC key must not
// be able to forge a token by signing it with HS256 using that public key's
// PEM bytes as the HMAC secret.
func TestAlgorithmConfusion_ForgedHS256UsingEd25519PublicKey(t *testing.T) {
	_, pubPEM := genEd25519KeyPEMs(t)

	forged, err := jwtNewWithClaimsHS256(map[string]interface{}{
		"sub": "admin",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}, []byte(pubPEM))
	require.NoError(t, err)

	_, err = verifyFn(forged, pubPEM, VerifyOptions{Algorithms: []string{"EdDSA"}})
	require.Error(t, err, "forged HS256 token must be rejected when the server only accepts EdDSA")

	// As with the RSA analogue: assert on the SPECIFIC WithValidMethods
	// error, not merely "an error occurred" -- confirmed fail-capable by
	// temporarily removing jwt.WithValidMethods(opts.Algorithms) from
	// verifyFn and re-running this test locally: with the defence removed,
	// this assertion fails, because the forged token's declared HS256 alg is
	// dispatched to SigningMethodHS256.Verify with the resolved EdDSA key
	// (an ed25519.PublicKey), which produces jwt/v5's "key is of invalid
	// type: HMAC verify expects []byte" -- the same message and same
	// ErrTokenSignatureInvalid sentinel as the RSA analogue's regression
	// case, and equally NOT proof the algorithm-family defence fired. This
	// test would correctly catch that regression rather than passing either
	// way.
	require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	assert.Contains(t, err.Error(), "signing method HS256 is invalid",
		"error must show WithValidMethods rejected the token's declared alg, not merely that some later step failed")
}

func TestVerify_NoneAlgHeader_Rejected(t *testing.T) {
	// Hand-build a token whose header declares alg: none, with an empty
	// signature -- the classic bypass payload. opts.algorithms is non-empty
	// and deliberately does not (and cannot, since "none" is always
	// refused) include "none".
	token := rawUnsignedNoneToken(t, map[string]interface{}{"sub": "admin", "exp": float64(time.Now().Add(time.Hour).Unix())})

	_, err := verifyFn(token, "any-secret", VerifyOptions{Algorithms: []string{"HS256"}})
	require.Error(t, err, "a token declaring alg:none must never verify, even with an empty signature and a non-empty algorithms list")
}

func TestVerify_LeewayAllowsJustExpiredToken(t *testing.T) {
	token, err := signFn(map[string]interface{}{
		"sub": "svc-a",
		"exp": float64(time.Now().Add(-2 * time.Second).Unix()),
	}, "secret", SignOptions{Algorithm: "HS256"})
	require.NoError(t, err)

	// Without leeway: rejected.
	_, err = verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}})
	require.Error(t, err)

	// With enough leeway: accepted.
	_, err = verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}, LeewaySeconds: 10})
	require.NoError(t, err)
}

// TestVerify_LeewaySeconds_Unbounded_Rejected is the exact regression this
// chunk fixes (security review MEDIUM 2): an unvalidated leeway_seconds lets
// a caller pass math.huge (gopher-lua's math.huge is math.MaxFloat64, not
// +Inf) or another absurdly large value and have an hours-expired token
// verify successfully. Each case here must be rejected before ever reaching
// jwt/v5's exp check.
func TestVerify_LeewaySeconds_Unbounded_Rejected(t *testing.T) {
	// A token that expired 100 hours ago -- the scenario from the security
	// review, confirmed to verify successfully on arm64 before this fix.
	token, err := signFn(map[string]interface{}{
		"sub": "svc-a",
		"exp": float64(time.Now().Add(-100 * time.Hour).Unix()),
	}, "secret", SignOptions{Algorithm: "HS256"})
	require.NoError(t, err)

	cases := map[string]float64{
		"math.MaxFloat64 (gopher-lua's math.huge)": math.MaxFloat64,
		"1e18":                        1e18,
		"NaN":                         math.NaN(),
		"negative":                    -1,
		"just over the ceiling (301)": 301,
	}
	for name, leeway := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}, LeewaySeconds: leeway})
			require.Error(t, err, "leeway_seconds=%v must be rejected, not silently disable exp enforcement", leeway)
		})
	}
}

// TestVerify_LeewaySeconds_AtCeiling_Accepted: the ceiling itself
// (maxLeewaySeconds) must still work for a legitimate small-clock-skew use
// -- the fix must reject only what is out of bounds, not the bound itself.
func TestVerify_LeewaySeconds_AtCeiling_Accepted(t *testing.T) {
	token, err := signFn(map[string]interface{}{
		"sub": "svc-a",
		"exp": float64(time.Now().Add(-1 * time.Minute).Unix()),
	}, "secret", SignOptions{Algorithm: "HS256"})
	require.NoError(t, err)

	_, err = verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}, LeewaySeconds: maxLeewaySeconds})
	require.NoError(t, err, "leeway_seconds at the ceiling must still be accepted")
}

func TestVerify_WrongSubject_Rejected(t *testing.T) {
	token, err := signFn(map[string]interface{}{
		"sub": "svc-a",
		"exp": float64(time.Now().Add(time.Hour).Unix()),
	}, "secret", SignOptions{Algorithm: "HS256"})
	require.NoError(t, err)

	_, err = verifyFn(token, "secret", VerifyOptions{Algorithms: []string{"HS256"}, Subject: "svc-b"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "subject")
}

func TestVerify_MalformedToken_Rejected(t *testing.T) {
	_, err := verifyFn("not-a-real-token", "secret", VerifyOptions{Algorithms: []string{"HS256"}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "malformed")
}

func TestSign_KeyDoesNotMatchAlgorithm_Raises(t *testing.T) {
	_, ecPubPEM := genECKeyPEMs(t)
	ecPrivPEM, _ := genECKeyPEMs(t)

	// Signing with RS256 but supplying an EC private key must fail clearly,
	// not silently produce a broken token.
	_, err := signFn(map[string]interface{}{"sub": "x"}, ecPrivPEM, SignOptions{Algorithm: "RS256"})
	require.Error(t, err)

	// Verifying RS256 with an EC public key must fail the same way.
	_, err = verifyFn("irrelevant.token.here", ecPubPEM, VerifyOptions{Algorithms: []string{"RS256"}})
	require.Error(t, err)
}

// TestSign_KeyDoesNotMatchAlgorithm_EdDSA_Raises: the EdDSA-specific
// wrong-key-type cases in both directions -- an RSA key under EdDSA, and an
// Ed25519 key under RS256.
func TestSign_KeyDoesNotMatchAlgorithm_EdDSA_Raises(t *testing.T) {
	rsaPrivPEM, rsaPubPEM := genRSAKeyPEMs(t)
	edPrivPEM, edPubPEM := genEd25519KeyPEMs(t)

	// Signing with EdDSA but supplying an RSA private key must fail clearly.
	_, err := signFn(map[string]interface{}{"sub": "x"}, rsaPrivPEM, SignOptions{Algorithm: "EdDSA"})
	require.Error(t, err, "an RSA private key must not be accepted for EdDSA")

	// Verifying EdDSA with an RSA public key must fail the same way.
	_, err = verifyFn("irrelevant.token.here", rsaPubPEM, VerifyOptions{Algorithms: []string{"EdDSA"}})
	require.Error(t, err, "an RSA public key must not be accepted for EdDSA")

	// Signing with RS256 but supplying an Ed25519 private key must fail
	// clearly -- the reverse direction of the confusion.
	_, err = signFn(map[string]interface{}{"sub": "x"}, edPrivPEM, SignOptions{Algorithm: "RS256"})
	require.Error(t, err, "an Ed25519 private key must not be accepted for RS256")

	// Verifying RS256 with an Ed25519 public key must fail the same way.
	_, err = verifyFn("irrelevant.token.here", edPubPEM, VerifyOptions{Algorithms: []string{"RS256"}})
	require.Error(t, err, "an Ed25519 public key must not be accepted for RS256")
}

// TestParseKey_Ed448Rejected_AtParseTime: Ed448 shares the JOSE alg name
// "EdDSA" with Ed25519 (the alg string alone cannot disambiguate the
// curve), but Go's stdlib crypto/x509 has no Ed448 case at all -- so an
// Ed448 key fails to parse with a clear x509 error rather than being
// silently accepted and misread as Ed25519. Go's stdlib cannot generate a
// real Ed448 key (crypto/ed25519 is Ed25519-only), so this test hand-builds
// only the ASN.1 shape needed to reach the OID switch inside
// x509.ParsePKCS8PrivateKey -- see fakeEd448PKCS8PEM's comment. This proves
// the rejection happens at PEM-parse time (inside parseSignKey, before any
// type switch on the parsed key), which is the "comprehensible error"
// requirement for a curve stdlib cannot even represent.
func TestParseKey_Ed448Rejected_AtParseTime(t *testing.T) {
	fakeEd448PEM := fakeEd448PKCS8PEM(t)

	_, err := signFn(map[string]interface{}{"sub": "x"}, fakeEd448PEM, SignOptions{Algorithm: "EdDSA"})
	require.Error(t, err, "an Ed448 key must be rejected, not silently misread as Ed25519")
	assert.Contains(t, err.Error(), "unknown algorithm",
		"the error must come from x509's OID switch rejecting Ed448 at parse time, not a later, more confusing failure")
}

func TestSign_ClaimsNotEncodable_Raises(t *testing.T) {
	_, err := signFn(map[string]interface{}{"bad": math.NaN()}, "secret", SignOptions{Algorithm: "HS256"})
	require.Error(t, err, "a NaN claim value cannot be JSON-encoded and must raise, not produce a broken token")
}

func TestDecodeUnverified_MalformedInput_Raises(t *testing.T) {
	_, err := decodeUnverifiedFn("only-one-segment")
	require.Error(t, err)
}

func TestValidateAlgorithms_EmptyRaises(t *testing.T) {
	_, err := validateAlgorithms(nil, "jwt.verify")
	require.Error(t, err)
}

func TestValidateAlgorithms_NoneRaises(t *testing.T) {
	_, err := validateAlgorithms([]string{"none"}, "jwt.verify")
	require.Error(t, err)
}

func TestValidateAlgorithms_MixedFamiliesRaises(t *testing.T) {
	_, err := validateAlgorithms([]string{"HS256", "RS256"}, "jwt.verify")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mixes")
}

func TestValidateAlgorithms_MixedRSAAndECRaises(t *testing.T) {
	// Not the textbook HMAC/RSA confusion, but exactly as ambiguous for key
	// parsing -- see the package doc's reasoning for drawing the family
	// line at three buckets rather than only forbidding HMAC+asymmetric.
	_, err := validateAlgorithms([]string{"RS256", "ES256"}, "jwt.verify")
	require.Error(t, err)
}

// TestValidateAlgorithms_MixedEdDSAAndECRaises and
// TestValidateAlgorithms_MixedEdDSAAndHMACRaises: EdDSA is its own family
// (Ed25519 keys are neither RSA nor EC keys), so it must raise on mixing
// with any other family exactly like the pre-existing pairs above -- in the
// same shape as TestAlgorithmConfusion_ForgedHS256UsingRSAPublicKey's
// sibling validation tests.
func TestValidateAlgorithms_MixedEdDSAAndECRaises(t *testing.T) {
	_, err := validateAlgorithms([]string{"EdDSA", "ES256"}, "jwt.verify")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mixes")
}

func TestValidateAlgorithms_MixedEdDSAAndHMACRaises(t *testing.T) {
	_, err := validateAlgorithms([]string{"EdDSA", "HS256"}, "jwt.verify")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mixes")
}

func TestValidateAlgorithms_UnsupportedNameRaises(t *testing.T) {
	_, err := validateAlgorithms([]string{"HS1"}, "jwt.verify")
	require.Error(t, err)
}

func TestValidateAlgorithms_SameFamilyOK(t *testing.T) {
	family, err := validateAlgorithms([]string{"RS256", "PS256"}, "jwt.verify")
	require.NoError(t, err)
	assert.Equal(t, familyRSA, family)
}

// TestLuaScripts: runs every Lua test script in testdata/.
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/*.lua")
	require.NoError(t, err)
	require.NotEmpty(t, files, "no Lua test files found in testdata/")

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()
			L.PreloadModule("jwt", Loader)
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
