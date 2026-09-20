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

// Package jwt decodes, verifies and signs JSON Web Tokens, designing out the
// two classic JWT footguns at the API level rather than documenting around
// them.
//
// # decode_unverified, not decode
//
// The unverified-decode function is deliberately named decode_unverified,
// with no decode alias. jwt.decode(token) would be the most-misused
// function in this library — naming it decode_unverified means the
// dangerous shortcut can't be reached by the shorter, more inviting name.
// It performs NO signature check; use it only to read kid/iss before
// choosing a verification key.
//
// # Algorithm confusion is structurally impossible, not just discouraged
//
// The classic attack: a server picks the key type from the token's own
// header, so an attacker takes the server's RSA PUBLIC key and signs a
// forged token with HS256 using that public key bytes as the HMAC secret.
// If the server's verification code blindly trusts the header's declared
// algorithm to choose how to interpret its key material, the forged HS256
// token verifies.
//
// This package closes that off two ways:
//
//   - The acceptable algorithm set comes ONLY from the caller's
//     opts.algorithms, never from the token header. The token's declared alg
//     is checked for membership in that set (via jwt/v5's WithValidMethods)
//     and is NEVER used to pick which key type to parse `key` as.
//   - opts.algorithms may not mix families: a caller passing
//     {"HS256", "RS256"} gets a raised error before the token is even
//     touched, because a single `key` argument cannot safely be interpreted
//     two ways. A caller needing both must call verify twice, once per key.
//     This package draws the family line at four buckets -- HMAC, RSA-ish
//     (RS*/PS*, which share the same *rsa.PublicKey/*rsa.PrivateKey key
//     type), EC (ES*), and EdDSA (Ed25519) -- rather than only forbidding
//     the literal HMAC-vs-RSA pairing the classic attack uses: mixing RS256
//     and ES256 would be exactly as ambiguous (which PEM key type is
//     `key`?), even though it isn't the textbook CVE. EdDSA gets its own
//     family rather than folding into EC: Ed25519 keys are neither RSA nor
//     ECDSA keys, and {"EdDSA", "ES256"} is exactly as ambiguous as any
//     other cross-family pairing.
//
// alg: none is refused unconditionally: it can never appear in
// opts.algorithms (rejected at option-validation time, before parsing), so
// WithValidMethods can never admit a token that declares it.
//
// # Accepted key encodings
//
// RS*/PS*/ES*/EdDSA all take `key` as PEM (a PEM public key or certificate
// for verify, a PEM private key for sign) -- never raw key bytes. This
// matters most for EdDSA: an Ed25519 private key can also be represented as
// a raw 32-byte seed or 64-byte expanded key, neither of which is PEM. This
// package does not accept either raw form, deliberately: a raw byte string
// is indistinguishable at the Lua boundary from an HS* secret (both are
// just an arbitrary-length Lua string), so accepting it would reopen a
// smaller version of the exact key/algorithm ambiguity the family
// partitioning exists to prevent. PEM-only keeps one uniform answer to
// "what shape is `key`" for every non-HMAC algorithm.
//
// # Claim validation
//
// exp and nbf are always validated when present (there is no option to
// disable either) via golang-jwt/jwt/v5's own Validator, which is exactly
// the code this dependency is being taken FOR -- exp/nbf semantics with
// leeway are where hand-rolled JWT validation tends to go subtly wrong.
// leeway_seconds is capped at maxLeewaySeconds (300s / 5 minutes): it is the
// one option that can effectively switch OFF exp/nbf enforcement if left
// unbounded (an unvalidated math.huge or 1e18 makes an hours-expired token
// verify successfully), so unlike the other options it gets a hard ceiling,
// not just a sign check. iat
// is never validated (RFC 7519 defines it as informational, not a security
// control, and comparing it invites clock-skew noise for no benefit).
// leeway_seconds applies equally to exp and nbf. issuer/audience/subject are
// checked only when the corresponding option is non-empty.
//
// Every VerifyOptions field is named so its Go zero value is the SAFE
// value: an omitted issuer/audience/subject means "don't check" (the safe
// default for a check that hasn't been configured), an omitted
// leeway_seconds means zero leeway (the strict default), and — the one that
// matters most — allow_missing_exp defaults to false, meaning a token with
// no exp claim is REJECTED unless a caller opts in. There is deliberately
// no require_exp field: naming it that way would make the safe behaviour
// depend on the caller remembering to set a flag, rather than needing to
// deliberately relax it.
package jwt

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// algFamily: the three key-interpretation families this package recognises.
// Every supported algorithm maps to exactly one, and opts.algorithms may
// only ever span one family per call — see the package doc.
type algFamily int

const (
	familyUnknown algFamily = iota
	familyHMAC
	familyRSA // RS* and PS*: both use *rsa.PublicKey / *rsa.PrivateKey
	familyEC
	familyEdDSA // EdDSA (Ed25519 only -- see algFamilies comment on Ed448)
)

// maxLeewaySeconds: the ceiling accepted for VerifyOptions.LeewaySeconds.
// 300s (5 minutes) is a generous bound for real clock skew between an
// issuer and a verifier (NTP-synced clocks are typically within seconds;
// this covers even badly-drifted hosts) while remaining far below anything
// that could meaningfully weaken exp/nbf enforcement. leeway_seconds is the
// one VerifyOptions field that can switch OFF the check this module exists
// to perform (see MEDIUM 2 in the security review this const references),
// so unlike other options it needs a hard ceiling, not just a sign check --
// math.huge (which gopher-lua maps to math.MaxFloat64, not +Inf) or 1e18
// both make an exp-hours-ago token verify successfully if left unbounded.
const maxLeewaySeconds = 300

// algFamilies: the complete, closed set of algorithms this package accepts.
// "none" and anything not listed here are unsupported. This map is the
// single source of truth both for what verify/sign will accept and for
// family-mixing detection.
//
// EdDSA has exactly one member, unlike the other families: JOSE defines no
// EdDSA256/384/512 variants. In JOSE, the string "EdDSA" is also the alg
// name for Ed448, not just Ed25519 -- the algorithm name alone does not
// disambiguate the curve. This package only ever produces/accepts Ed25519
// keys: Go's crypto/x509 and crypto/ed25519 do not implement Ed448 at all,
// so an Ed448 key fails to parse (see parseVerifyKey/parseSignKey) rather
// than being silently misinterpreted as Ed25519.
var algFamilies = map[string]algFamily{
	"HS256": familyHMAC, "HS384": familyHMAC, "HS512": familyHMAC,
	"RS256": familyRSA, "RS384": familyRSA, "RS512": familyRSA,
	"PS256": familyRSA, "PS384": familyRSA, "PS512": familyRSA,
	"ES256": familyEC, "ES384": familyEC, "ES512": familyEC,
	"EdDSA": familyEdDSA,
}

// Decoded: the result of decode_unverified. Both fields are populated
// WITHOUT any signature check -- see the package doc.
type Decoded struct {
	Header map[string]interface{} `json:"header"` // the token's header segment, decoded
	Claims map[string]interface{} `json:"claims"` // the token's claims segment, decoded
}

// VerifyOptions: required options for verify (F2 -- there are no optional
// arguments in this framework). Every field's zero value is the safe value
// -- see the package doc for why allow_missing_exp is spelled that way
// rather than require_exp.
type VerifyOptions struct {
	Algorithms      []string `json:"algorithms"`        // REQUIRED, non-empty, single family, never "none"
	Issuer          string   `json:"issuer"`            // "" = do not check
	Audience        string   `json:"audience"`          // "" = do not check
	Subject         string   `json:"subject"`           // "" = do not check
	LeewaySeconds   float64  `json:"leeway_seconds"`    // 0 = no leeway; applied to both exp and nbf
	AllowMissingExp bool     `json:"allow_missing_exp"` // false (default): a token with no exp is REJECTED
}

// SignOptions: required options for sign.
type SignOptions struct {
	Algorithm string `json:"algorithm"` // REQUIRED; never "none"
	Kid       string `json:"kid"`       // "" = omit the header field
	Typ       string `json:"typ"`       // "" = "JWT" (jwt/v5's own default)
}

// validateAlgorithms: the single choke point for the whole "algorithm
// confusion" defense (see package doc). Used both by verify (opts.algorithms,
// which may list several algorithms of one family) and sign (a single
// algorithm, passed as a one-element slice so the same rejection rules
// apply uniformly). Raises on: an empty list, "none" appearing anywhere, an
// algorithm outside the supported set, or algorithms spanning more than one
// family.
func validateAlgorithms(algs []string, context string) (algFamily, error) {
	if len(algs) == 0 {
		return familyUnknown, fmt.Errorf("%s: algorithms is required and must not be empty", context)
	}
	var family algFamily
	var first string
	for i, alg := range algs {
		if alg == "none" {
			return familyUnknown, fmt.Errorf("%s: algorithm \"none\" is never accepted", context)
		}
		f, ok := algFamilies[alg]
		if !ok {
			return familyUnknown, fmt.Errorf("%s: unsupported algorithm %q", context, alg)
		}
		if i == 0 {
			family = f
			first = alg
			continue
		}
		if f != family {
			return familyUnknown, fmt.Errorf(
				"%s: algorithms mixes algorithm families (%q and %q) -- a single key argument cannot be interpreted as both; call verify separately per key type",
				context, first, alg)
		}
	}
	return family, nil
}

// validateLeeway: rejects a LeewaySeconds value before it is ever converted
// to a time.Duration. This MUST validate the raw float64, not a value
// already converted to an integer type: Go's float64->int64 conversion for
// an out-of-range value is architecture-defined (it saturates to MaxInt64
// on arm64 and produces MinInt64 on amd64), so checking a post-conversion
// value would behave differently per platform and CI (which runs amd64)
// would not catch an arm64-only bypass. Rejects NaN, +/-Inf, negative
// values and anything above maxLeewaySeconds.
func validateLeeway(seconds float64) (time.Duration, error) {
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return 0, fmt.Errorf("jwt.verify: leeway_seconds must be a finite number, got %v", seconds)
	}
	if seconds < 0 {
		return 0, fmt.Errorf("jwt.verify: leeway_seconds must be >= 0, got %v", seconds)
	}
	if seconds > maxLeewaySeconds {
		return 0, fmt.Errorf("jwt.verify: leeway_seconds must be <= %d, got %v", maxLeewaySeconds, seconds)
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

// parseVerifyKey: interprets key according to family, NEVER according to
// anything read from the token. HMAC families use the raw secret bytes;
// RSA/EC/EdDSA families parse key as a PEM public key (jwt/v5's
// ParseRSAPublicKeyFromPEM/ParseECPublicKeyFromPEM accept a certificate
// too; ParseEdPublicKeyFromPEM only accepts a PKIX public key -- see the
// package doc's "accepted key encodings" note). A key encoding other than
// PEM (raw seed/expanded bytes for EdDSA) is deliberately rejected: see
// parseSignKey.
func parseVerifyKey(family algFamily, key string) (interface{}, error) {
	switch family {
	case familyHMAC:
		return []byte(key), nil
	case familyRSA:
		return jwt.ParseRSAPublicKeyFromPEM([]byte(key))
	case familyEC:
		return jwt.ParseECPublicKeyFromPEM([]byte(key))
	case familyEdDSA:
		return jwt.ParseEdPublicKeyFromPEM([]byte(key))
	default:
		return nil, errors.New("unreachable: unvalidated algorithm family")
	}
}

// parseSignKey: the signing-side mirror of parseVerifyKey. RSA/EC/EdDSA
// families parse key as a PEM PRIVATE key; a key of the wrong type (e.g. an
// EC PEM key for an RSA algorithm, or an RSA key for EdDSA) fails here with
// jwt/v5's own ErrNotRSAPrivateKey/ErrNotECPrivateKey/ErrNotEdPrivateKey,
// which is the "reject a mismatched key with a clear error" requirement.
//
// EdDSA private keys are accepted as PEM (PKCS#8) only, matching RS*/ES* --
// not as a raw 32-byte seed or 64-byte expanded key. A raw byte string is
// indistinguishable from an HMAC secret at the Lua boundary (both arrive as
// a plain Lua string of arbitrary length); accepting one alongside PEM
// would let a caller pass the wrong kind of secret for the wrong algorithm
// and have it silently "work" as something other than intended. PEM keeps
// the "what shape must `key` be" answer uniform across every non-HMAC
// family: jwt.ParseEdPrivateKeyFromPEM (jwt/v5's own helper, mirroring
// ParseRSAPrivateKeyFromPEM/ParseECPrivateKeyFromPEM) rejects anything else.
//
// An Ed448 key is not a supported input: Go's crypto/x509 has no Ed448 OID
// case in ParsePKCS8PrivateKey/ParsePKIXPublicKey, so it fails to parse at
// all (surfaced here, at parse time, as "x509: PKCS#8 wrapping contained
// private key with unknown algorithm: 1.3.101.113" or the PKIX equivalent)
// rather than being silently accepted and misinterpreted as Ed25519 -- see
// TestParseKey_Ed448Rejected_AtParseTime.
func parseSignKey(family algFamily, key string) (interface{}, error) {
	switch family {
	case familyHMAC:
		return []byte(key), nil
	case familyRSA:
		return jwt.ParseRSAPrivateKeyFromPEM([]byte(key))
	case familyEC:
		return jwt.ParseECPrivateKeyFromPEM([]byte(key))
	case familyEdDSA:
		return jwt.ParseEdPrivateKeyFromPEM([]byte(key))
	default:
		return nil, errors.New("unreachable: unvalidated algorithm family")
	}
}

// describeVerifyError: rewraps a jwt/v5 parse/validation error with a short
// prefix naming WHICH check failed, so an operator's pcall message doesn't
// read as a generic "invalid token" -- see SPECS.md S13's test plan.
// The original error remains wrapped (%w) so errors.Is still works for a
// caller that wants to branch on the exact jwt/v5 sentinel.
func describeVerifyError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return fmt.Errorf("token expired: %w", err)
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fmt.Errorf("token not valid yet (nbf): %w", err)
	case errors.Is(err, jwt.ErrTokenInvalidIssuer):
		return fmt.Errorf("issuer mismatch: %w", err)
	case errors.Is(err, jwt.ErrTokenInvalidAudience):
		return fmt.Errorf("audience mismatch: %w", err)
	case errors.Is(err, jwt.ErrTokenInvalidSubject):
		return fmt.Errorf("subject mismatch: %w", err)
	case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
		return fmt.Errorf("required claim missing (exp): %w", err)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return fmt.Errorf("signature invalid or algorithm not permitted: %w", err)
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf("malformed token: %w", err)
	default:
		return err
	}
}

// decodeUnverifiedFn: implements jwt.decode_unverified. Performs NO
// signature check -- see the package doc.
func decodeUnverifiedFn(token string) (Decoded, error) {
	parsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return Decoded{}, fmt.Errorf("jwt.decode_unverified: %w", err)
	}
	claims, _ := parsed.Claims.(jwt.MapClaims)
	header := parsed.Header
	if header == nil {
		header = map[string]interface{}{}
	}
	return Decoded{
		Header: header,
		Claims: map[string]interface{}(claims),
	}, nil
}

// verifyFn: implements jwt.verify. Raises on ANY failure -- see the package
// doc and SPECS.md S13.
func verifyFn(token, key string, opts VerifyOptions) (map[string]interface{}, error) {
	family, err := validateAlgorithms(opts.Algorithms, "jwt.verify")
	if err != nil {
		return nil, err
	}

	verifyKey, err := parseVerifyKey(family, key)
	if err != nil {
		return nil, fmt.Errorf("jwt.verify: invalid key: %w", err)
	}

	leeway, err := validateLeeway(opts.LeewaySeconds)
	if err != nil {
		return nil, err
	}

	parserOpts := []jwt.ParserOption{
		jwt.WithValidMethods(opts.Algorithms),
		jwt.WithLeeway(leeway),
	}
	if opts.Issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(opts.Issuer))
	}
	if opts.Audience != "" {
		parserOpts = append(parserOpts, jwt.WithAudience(opts.Audience))
	}
	if opts.Subject != "" {
		parserOpts = append(parserOpts, jwt.WithSubject(opts.Subject))
	}
	if !opts.AllowMissingExp {
		parserOpts = append(parserOpts, jwt.WithExpirationRequired())
	}

	claims := jwt.MapClaims{}
	_, err = jwt.NewParser(parserOpts...).ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		// The key is fixed by opts.algorithms' family, resolved above --
		// NEVER by anything read from the token being verified here. See
		// the package doc's "Algorithm confusion" section.
		return verifyKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt.verify: %w", describeVerifyError(err))
	}
	return map[string]interface{}(claims), nil
}

// signFn: implements jwt.sign.
func signFn(claims map[string]interface{}, key string, opts SignOptions) (string, error) {
	family, err := validateAlgorithms([]string{opts.Algorithm}, "jwt.sign")
	if err != nil {
		return "", err
	}

	method := jwt.GetSigningMethod(opts.Algorithm)
	if method == nil {
		return "", fmt.Errorf("jwt.sign: unsupported algorithm %q", opts.Algorithm)
	}

	signKey, err := parseSignKey(family, key)
	if err != nil {
		return "", fmt.Errorf("jwt.sign: key does not match algorithm %q: %w", opts.Algorithm, err)
	}

	token := jwt.NewWithClaims(method, jwt.MapClaims(claims))
	if opts.Kid != "" {
		token.Header["kid"] = opts.Kid
	}
	if opts.Typ != "" {
		token.Header["typ"] = opts.Typ
	}

	signed, err := token.SignedString(signKey)
	if err != nil {
		return "", fmt.Errorf("jwt.sign: claims not encodable: %w", err)
	}
	return signed, nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("jwt", "JWT decode, verify and sign, with alg:none and algorithm confusion designed out")

	m.Fn("decode_unverified", decodeUnverifiedFn,
		"decodes a JWT's header and claims WITHOUT verifying its signature -- DOES NOT verify anything; use only to read kid/iss before choosing a key",
		luareg.Args("token"),
		luareg.ArgDoc("token", "the compact JWT string (header.claims.signature)"),
		luareg.ReturnDoc(0, "decoded", "the unverified header and claims"))

	m.Fn("verify", verifyFn,
		"verifies a JWT's signature and claims; raises on ANY failure: bad signature, alg not in opts.algorithms, alg:none, expired, not yet valid, issuer/audience/subject mismatch, or an unparseable key",
		luareg.Args("token", "key", "opts"),
		luareg.ArgDoc("token", "the compact JWT string"),
		luareg.ArgDoc("key", "HS*: the raw shared secret; RS*/PS*/ES*/EdDSA: a PEM public key or certificate (never a raw Ed25519 seed/expanded key)"),
		luareg.ArgDoc("opts", "required options: algorithms (non-empty, single family, never \"none\"), issuer, audience, subject, leeway_seconds (0..300, seconds), allow_missing_exp"),
		luareg.ReturnDoc(0, "claims", "the token's verified claims"))

	m.Fn("sign", signFn,
		"signs claims into a compact JWT; raises on an unknown/forbidden algorithm, a key that does not match the algorithm, or claims that cannot be JSON-encoded",
		luareg.Args("claims", "key", "opts"),
		luareg.ArgDoc("claims", "the claims to encode"),
		luareg.ArgDoc("key", "HS*: the raw shared secret; RS*/PS*/ES*/EdDSA: a PEM private key (never a raw Ed25519 seed/expanded key)"),
		luareg.ArgDoc("opts", "required options: algorithm (never \"none\"), kid, typ"),
		luareg.ReturnDoc(0, "token", "the compact JWT string"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("jwt", jwt.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
