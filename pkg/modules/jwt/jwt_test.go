// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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
