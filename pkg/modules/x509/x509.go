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

// Package x509 parses PEM certificates FROM A STRING and answers questions
// about them. There is no x509.parse_file — a script that needs to read a
// certificate off disk uses the existing fs module and passes the resulting
// string in.
//
// Two hard invariants, both enforced by construction rather than by
// documentation alone:
//
//  1. No file reads. Every function takes PEM text as a Lua string.
//  2. No network, and no system trust store. crypto/x509.Verify never does
//     AIA/OCSP/CRL fetching, so "no network" is free. But
//     crypto/x509.SystemCertPool reads the host filesystem and makes
//     verification results depend on whichever machine the script happens to
//     run on — this package NEVER calls it, anywhere, and never will. The
//     root and intermediate pools passed to verify_chain are built ONLY from
//     the PEM text the caller supplies; an empty `roots` argument can never
//     verify anything (see TestVerifyChain_EmptyRootsNeverSucceeds).
//
// Timestamps (not_before, not_after) are Unix-second numbers, not RFC3339
// strings, so they compose with the time module (time.format, time.diff)
// the same way the time and uuid modules do. This deliberately differs from
// the kubernetes module, which speaks RFC3339 strings — see the README's
// "Time representation" note for which module speaks which.
//
// serial is rendered as an exact decimal string rather than a Lua number:
// serial numbers can be up to 20 bytes, far past float64's exact integer
// range, and a Lua number would silently round it (see also
// netaddr.num_addresses and x509's own public_key_bits, which stays a
// regular number because RSA/EC key sizes never approach that range).
package x509

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// Certificate: everything x509.parse/x509.parse_chain report about a single
// parsed certificate.
type Certificate struct {
	Subject            string   `json:"subject"`              // RFC 2253-style distinguished name
	Issuer             string   `json:"issuer"`               // RFC 2253-style distinguished name
	SubjectCN          string   `json:"subject_cn"`           // subject's CommonName, "" if absent
	IssuerCN           string   `json:"issuer_cn"`            // issuer's CommonName, "" if absent
	SubjectOrg         []string `json:"subject_org"`          // subject's Organization values
	IssuerOrg          []string `json:"issuer_org"`           // issuer's Organization values
	Serial             string   `json:"serial"`               // exact decimal string; serials can exceed 2^53
	Version            int      `json:"version"`              // X.509 version (3 for a v3 certificate)
	NotBefore          int64    `json:"not_before"`           // Unix seconds
	NotAfter           int64    `json:"not_after"`            // Unix seconds
	DNSNames           []string `json:"dns_names"`            // SAN dNSName entries
	IPAddresses        []string `json:"ip_addresses"`         // SAN iPAddress entries, v4 and v6, in net.IP.String() form
	EmailAddresses     []string `json:"email_addresses"`      // SAN rfc822Name entries
	URIs               []string `json:"uris"`                 // SAN uniformResourceIdentifier entries
	IsCA               bool     `json:"is_ca"`                // basic constraints CA flag
	MaxPathLen         int      `json:"max_path_len"`         // basic constraints path length constraint, 0 if unset
	KeyUsage           []string `json:"key_usage"`            // e.g. {"digital_signature", "key_encipherment"}
	ExtKeyUsage        []string `json:"ext_key_usage"`        // e.g. {"server_auth", "client_auth"}
	SignatureAlgorithm string   `json:"signature_algorithm"`  // e.g. "SHA256-RSA"
	PublicKeyAlgorithm string   `json:"public_key_algorithm"` // e.g. "RSA", "ECDSA", "Ed25519"
	PublicKeyBits      int      `json:"public_key_bits"`      // e.g. 2048 for a 2048-bit RSA key
	FingerprintSHA256  string   `json:"fingerprint_sha256"`   // lowercase hex SHA-256 of the raw DER
}

// VerifyChainOptions: required options for verify_chain (F2 — there are no
// optional arguments in this framework, so this whole table argument is
// mandatory).
//
// at_time has no safe zero value: silently verifying "as of 1970-01-01"
// would pass or fail every certificate for the wrong reason, so at_time ==
// 0 raises rather than defaulting to any particular time.
type VerifyChainOptions struct {
	DNSName   string   `json:"dns_name"`   // hostname to check the leaf against; "" skips the hostname check
	AtTime    int64    `json:"at_time"`    // Unix seconds; REQUIRED, must not be 0 — pass time.now()
	KeyUsages []string `json:"key_usages"` // required extended key usages; empty means {"server_auth"}
}

// keyUsageBit: one bit of x509.KeyUsage and its Lua-facing name.
type keyUsageBit struct {
	bit  x509.KeyUsage
	name string
}

// keyUsageBits: every x509.KeyUsage bit this package decodes, in the order
// they are tested. Verbatim from SPECS.md S12's key_usage example naming.
var keyUsageBits = []keyUsageBit{
	{x509.KeyUsageDigitalSignature, "digital_signature"},
	{x509.KeyUsageContentCommitment, "content_commitment"},
	{x509.KeyUsageKeyEncipherment, "key_encipherment"},
	{x509.KeyUsageDataEncipherment, "data_encipherment"},
	{x509.KeyUsageKeyAgreement, "key_agreement"},
	{x509.KeyUsageCertSign, "cert_sign"},
	{x509.KeyUsageCRLSign, "crl_sign"},
	{x509.KeyUsageEncipherOnly, "encipher_only"},
	{x509.KeyUsageDecipherOnly, "decipher_only"},
}

// extKeyUsageByName: the extended key usages this package understands, both
// for decoding a parsed certificate's ext_key_usage and for accepting
// verify_chain's opts.key_usages by name. Deliberately a small, well-known
// set (the ones SPECS.md's example names) rather than every constant
// crypto/x509 defines.
var extKeyUsageByName = map[string]x509.ExtKeyUsage{
	"any":              x509.ExtKeyUsageAny,
	"server_auth":      x509.ExtKeyUsageServerAuth,
	"client_auth":      x509.ExtKeyUsageClientAuth,
	"code_signing":     x509.ExtKeyUsageCodeSigning,
	"email_protection": x509.ExtKeyUsageEmailProtection,
	"ipsec_end_system": x509.ExtKeyUsageIPSECEndSystem,
	"ipsec_tunnel":     x509.ExtKeyUsageIPSECTunnel,
	"ipsec_user":       x509.ExtKeyUsageIPSECUser,
	"time_stamping":    x509.ExtKeyUsageTimeStamping,
	"ocsp_signing":     x509.ExtKeyUsageOCSPSigning,
}

// extKeyUsageByValue: the reverse of extKeyUsageByName, built once at
// package init for decoding a parsed certificate's ExtKeyUsage list.
var extKeyUsageByValue = func() map[x509.ExtKeyUsage]string {
	m := make(map[x509.ExtKeyUsage]string, len(extKeyUsageByName))
	for name, v := range extKeyUsageByName {
		m[v] = name
	}
	return m
}()

// keyUsageStrings: decodes a KeyUsage bitmask into its set bit names, e.g.
// {"digital_signature", "key_encipherment", "cert_sign"}. Always returns a
// non-nil slice (empty when no bits are set) — see the skill's nil-slice
// rule.
func keyUsageStrings(u x509.KeyUsage) []string {
	out := []string{}
	for _, ku := range keyUsageBits {
		if u&ku.bit != 0 {
			out = append(out, ku.name)
		}
	}
	return out
}

// extKeyUsageStrings: decodes a certificate's ExtKeyUsage list into names,
// e.g. {"server_auth", "client_auth"}. An extended key usage outside the
// known set renders as "unknown_<n>" rather than being silently dropped, so
// a caller auditing a certificate's usages is never shown an incomplete
// list. Always returns a non-nil slice.
func extKeyUsageStrings(us []x509.ExtKeyUsage) []string {
	out := []string{}
	for _, u := range us {
		if name, ok := extKeyUsageByValue[u]; ok {
			out = append(out, name)
		} else {
			out = append(out, fmt.Sprintf("unknown_%d", int(u)))
		}
	}
	return out
}

// extKeyUsagesFromNames: the inverse of extKeyUsageStrings, used to turn
// verify_chain's opts.key_usages into the crypto/x509 constants
// VerifyOptions.KeyUsages expects. Raises on any name outside the known
// set — an unrecognised policy name is a caller mistake, not a legitimate
// "no match" outcome.
func extKeyUsagesFromNames(names []string) ([]x509.ExtKeyUsage, error) {
	out := make([]x509.ExtKeyUsage, 0, len(names))
	for _, n := range names {
		u, ok := extKeyUsageByName[n]
		if !ok {
			return nil, fmt.Errorf("x509.verify_chain: unknown key usage %q", n)
		}
		out = append(out, u)
	}
	return out, nil
}

// publicKeyBits: the bit size of a certificate's public key, for the three
// key types crypto/x509 can produce from a certificate. Unknown key types
// (e.g. DSA, which crypto/x509 can still parse) return 0 rather than
// raising — public_key_bits is informational, not something a policy should
// hard-fail on just because the key type is unusual.
func publicKeyBits(pub any) int {
	switch k := pub.(type) {
	case *rsa.PublicKey:
		return k.N.BitLen()
	case *ecdsa.PublicKey:
		return k.Curve.Params().BitSize
	case ed25519.PublicKey:
		return len(k) * 8
	default:
		return 0
	}
}

// toCertificate: converts a parsed *x509.Certificate into the Lua-facing
// Certificate struct. ip_addresses uses net.IP.String(), which renders both
// v4 and v6 addresses correctly (D5) — a v6 SAN comes back as e.g.
// "2001:db8::1", not garbled or dropped.
func toCertificate(cert *x509.Certificate) Certificate {
	ips := make([]string, len(cert.IPAddresses))
	for i, ip := range cert.IPAddresses {
		ips[i] = ip.String()
	}
	uris := make([]string, len(cert.URIs))
	for i, u := range cert.URIs {
		uris[i] = u.String()
	}

	fingerprint := sha256.Sum256(cert.Raw)

	return Certificate{
		Subject:            cert.Subject.String(),
		Issuer:             cert.Issuer.String(),
		SubjectCN:          cert.Subject.CommonName,
		IssuerCN:           cert.Issuer.CommonName,
		SubjectOrg:         append([]string{}, cert.Subject.Organization...),
		IssuerOrg:          append([]string{}, cert.Issuer.Organization...),
		Serial:             cert.SerialNumber.String(),
		Version:            cert.Version,
		NotBefore:          cert.NotBefore.Unix(),
		NotAfter:           cert.NotAfter.Unix(),
		DNSNames:           append([]string{}, cert.DNSNames...),
		IPAddresses:        ips,
		EmailAddresses:     append([]string{}, cert.EmailAddresses...),
		URIs:               uris,
		IsCA:               cert.IsCA,
		MaxPathLen:         cert.MaxPathLen,
		KeyUsage:           keyUsageStrings(cert.KeyUsage),
		ExtKeyUsage:        extKeyUsageStrings(cert.ExtKeyUsage),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		PublicKeyBits:      publicKeyBits(cert.PublicKey),
		FingerprintSHA256:  hex.EncodeToString(fingerprint[:]),
	}
}

// parseOne: decodes exactly one PEM block from pemStr and parses it as a
// certificate. Raises if there is no PEM block at all, the block is not a
// CERTIFICATE (e.g. a PRIVATE KEY block), or the DER payload fails to parse.
func parseOne(pemStr string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("x509.parse: no PEM block found")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("x509.parse: PEM block type %q is not CERTIFICATE", block.Type)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("x509.parse: %w", err)
	}
	return cert, nil
}

// parseChainCerts: decodes every PEM block in pemStr and parses the
// CERTIFICATE-typed ones, in file order, silently skipping any other block
// type (so a bundle that happens to also carry a PRIVATE KEY block still
// parses). Raises if a CERTIFICATE block fails to parse, or if the input
// contains no CERTIFICATE blocks at all.
func parseChainCerts(pemStr string) ([]*x509.Certificate, error) {
	rest := []byte(pemStr)
	var certs []*x509.Certificate
	for {
		var block *pem.Block
		block, rest = pem.Decode(rest)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" {
			continue
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("x509.parse_chain: %w", err)
		}
		certs = append(certs, cert)
	}
	if len(certs) == 0 {
		return nil, fmt.Errorf("x509.parse_chain: no CERTIFICATE blocks found")
	}
	return certs, nil
}

// parseFn: implements x509.parse.
func parseFn(pemStr string) (Certificate, error) {
	cert, err := parseOne(pemStr)
	if err != nil {
		return Certificate{}, err
	}
	return toCertificate(cert), nil
}

// parseChainFn: implements x509.parse_chain.
func parseChainFn(pemStr string) ([]Certificate, error) {
	certs, err := parseChainCerts(pemStr)
	if err != nil {
		return nil, err
	}
	out := make([]Certificate, len(certs))
	for i, c := range certs {
		out[i] = toCertificate(c)
	}
	return out, nil
}

// expiresInDaysFn: implements x509.expires_in_days. The result may be
// negative (already expired) or fractional; now is a required parameter
// rather than a hidden time.Now() call, so the result is reproducible and
// testable — pass time.now().
func expiresInDaysFn(pemStr string, now int64) (float64, error) {
	cert, err := parseOne(pemStr)
	if err != nil {
		return 0, err
	}
	return cert.NotAfter.Sub(time.Unix(now, 0)).Hours() / 24, nil
}

// isValidAtFn: implements x509.is_valid_at. Bounds are inclusive: exactly
// not_before or exactly not_after both count as valid.
func isValidAtFn(pemStr string, when int64) (bool, error) {
	cert, err := parseOne(pemStr)
	if err != nil {
		return false, err
	}
	t := time.Unix(when, 0)
	return !t.Before(cert.NotBefore) && !t.After(cert.NotAfter), nil
}

// verifyChainFn: implements x509.verify_chain. Returns (ok, reason) as
// normal values rather than raising on a failed verification — "this chain
// is not trusted" is an expected, catchable-by-default outcome, and raising
// would force every caller to pcall the common case. Only a malformed PEM
// argument or an invalid opts.key_usages name raises.
//
// rootPool and interPool are built EXCLUSIVELY from roots/intermediates via
// x509.NewCertPool + AppendCertsFromPEM. x509.SystemCertPool is never
// called, anywhere in this package — see the package doc. AppendCertsFromPEM
// itself skips any non-CERTIFICATE PEM block rather than erroring, so roots
// and intermediates may contain multiple certificates or be entirely empty.
func verifyChainFn(leaf, intermediates, roots string, opts VerifyChainOptions) (bool, string, error) {
	if opts.AtTime == 0 {
		return false, "", fmt.Errorf("x509.verify_chain: opts.at_time is required and must not be 0 (pass time.now())")
	}

	leafCert, err := parseOne(leaf)
	if err != nil {
		return false, "", err
	}

	rootPool := x509.NewCertPool()
	rootPool.AppendCertsFromPEM([]byte(roots))

	interPool := x509.NewCertPool()
	interPool.AppendCertsFromPEM([]byte(intermediates))

	usageNames := opts.KeyUsages
	if len(usageNames) == 0 {
		usageNames = []string{"server_auth"}
	}
	keyUsages, err := extKeyUsagesFromNames(usageNames)
	if err != nil {
		return false, "", err
	}

	_, err = leafCert.Verify(x509.VerifyOptions{
		DNSName:       opts.DNSName,
		Intermediates: interPool,
		Roots:         rootPool,
		CurrentTime:   time.Unix(opts.AtTime, 0),
		KeyUsages:     keyUsages,
	})
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("x509", "PEM certificate parsing and chain verification -- no file reads, no network, no system trust store")

	m.Fn("parse", parseFn,
		"parses the first CERTIFICATE PEM block in a string and returns its fields; raises on missing/wrong-type PEM block or bad DER",
		luareg.Args("pem"),
		luareg.ArgDoc("pem", "PEM-encoded certificate text"),
		luareg.ReturnDoc(0, "cert", "the parsed certificate"))

	m.Fn("parse_chain", parseChainFn,
		"parses every CERTIFICATE PEM block in a string, in file order; raises if any CERTIFICATE block fails to parse or none are found",
		luareg.Args("pem"),
		luareg.ArgDoc("pem", "PEM text containing one or more CERTIFICATE blocks (other block types are ignored)"),
		luareg.ReturnDoc(0, "certs", "the parsed certificates, in the order they appear"))

	m.Fn("expires_in_days", expiresInDaysFn,
		"returns the number of days between now and the certificate's not_after; may be negative (already expired) or fractional",
		luareg.Args("pem", "now"),
		luareg.ArgDoc("pem", "PEM-encoded certificate text"),
		luareg.ArgDoc("now", "the current time as Unix seconds, e.g. time.now()"),
		luareg.ReturnDoc(0, "days", "fractional days until expiry; negative if already expired"))

	m.Fn("is_valid_at", isValidAtFn,
		"reports whether when falls within [not_before, not_after], inclusive of both bounds",
		luareg.Args("pem", "when"),
		luareg.ArgDoc("pem", "PEM-encoded certificate text"),
		luareg.ArgDoc("when", "the time to check, as Unix seconds"),
		luareg.ReturnDoc(0, "ok", "true if not_before <= when <= not_after"))

	m.Fn("verify_chain", verifyChainFn,
		"verifies leaf against intermediates and roots; NEVER consults the system trust store, so an empty roots argument can never succeed",
		luareg.Args("leaf", "intermediates", "roots", "opts"),
		luareg.ArgDoc("leaf", "PEM-encoded leaf certificate to verify"),
		luareg.ArgDoc("intermediates", "PEM text with zero or more intermediate CA certificates"),
		luareg.ArgDoc("roots", "PEM text with zero or more trusted root certificates; empty means nothing is trusted"),
		luareg.ArgDoc("opts", "required options: dns_name, at_time (Unix seconds, required, must not be 0), key_usages (empty defaults to {\"server_auth\"})"),
		luareg.ReturnDoc(0, "ok", "true if the chain verifies against roots"),
		luareg.ReturnDoc(1, "reason", "empty string on success; the underlying x509 verification error text on failure"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("x509", x509.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
