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

// Package password provides bcrypt password hashing and verification for Lua
// scripts.
//
// bcrypt only: argon2id is deliberately deferred. bcrypt's output
// ($2a$<cost>$<salt><hash>) is self-describing — cost and salt travel with
// the hash, so verify needs no extra parameters and a cost rotation needs no
// schema migration. argon2 returns raw bytes; storing it requires
// hand-writing the PHC string format ($argon2id$v=19$m=...,t=...,p=...$...)
// and a parser for it, which is exactly the kind of hand-rolled format
// parsing this library avoids. If argon2id is needed later it is additive
// (argon2id_hash/argon2id_verify), not a replacement.
//
// Mismatch vs malformed hash, and why they are NOT the same outcome: this is
// the deliberate mirror image of hmac.verify_* (pkg/modules/hmac).
// verify's plaintext argument is attacker-controlled (a login attempt), so a
// wrong password is an ordinary, expected outcome and verify returns false
// for it. hash, by contrast, is the application's OWN data — read back from
// a database or a config file — so a malformed one means the storage layer
// or a migration is broken, not that the user mistyped their password. verify
// therefore RAISES on a malformed/non-bcrypt hash rather than returning
// false, so that class of bug surfaces immediately instead of silently
// reading as "wrong password" forever. See hmac's package doc for the
// opposite case (tag is attacker input, so verify_* never raises on it).
//
// Constant-time by construction: bcrypt.CompareHashAndPassword re-derives a
// hash from the candidate plaintext using the salt and cost embedded in the
// stored hash, then compares the two derived outputs with
// subtle.ConstantTimeCompare. The plaintext is never compared directly, and
// — same rule as hmac — no comparison primitive is exposed to Lua; there is
// no password.equal.
package password

import (
	"errors"
	"fmt"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/crypto/bcrypt"
)

// defaultCost: the production-recommended bcrypt cost, exposed as
// password.DEFAULT_COST. bcrypt's own package default (bcrypt.DefaultCost)
// is 10; this module recommends 12, matching current general guidance for
// interactive login hashing as of writing. Tests use a much lower cost (4,
// bcrypt's minimum) purely to keep the suite fast — cost 12 is roughly
// 250ms per call.
const defaultCost = 12

// maxPasswordBytes: bcrypt operates on at most 72 bytes of input and
// silently ignores anything beyond that internally. This module does not
// rely on that silent truncation — see hashPassword below — because two
// distinct passwords sharing the first 72 bytes would otherwise verify
// identically, which is a real footgun for callers who don't know about the
// limit. Passwords over this length are rejected outright.
const maxPasswordBytes = 72

// hashPassword: hashes plaintext with bcrypt at the given cost. Raises if
// cost is outside bcrypt's valid range [4, 31], or if plaintext exceeds 72
// bytes.
//
// The 72-byte check is explicit here rather than left to bcrypt alone:
// golang.org/x/crypto/bcrypt.GenerateFromPassword does return
// ErrPasswordTooLong for input over 72 bytes, so the two checks currently
// agree, but this module does not want its "no silent truncation" guarantee
// to depend on that implementation detail continuing to be enforced
// upstream in exactly this way.
func hashPassword(plaintext string, cost int) (string, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return "", fmt.Errorf("password.hash: cost %d outside allowed range [%d, %d]", cost, bcrypt.MinCost, bcrypt.MaxCost)
	}
	if len(plaintext) > maxPasswordBytes {
		return "", fmt.Errorf("password.hash: password is %d bytes, bcrypt supports at most %d", len(plaintext), maxPasswordBytes)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), cost)
	if err != nil {
		return "", fmt.Errorf("password.hash: %w", err)
	}
	return string(hashed), nil
}

// verifyPassword: reports whether plaintext matches the bcrypt hash. Returns
// false for a wrong plaintext (a normal outcome — plaintext is
// attacker-controlled). Raises if hash is not a well-formed bcrypt hash, since
// hash is the application's own stored data and a malformed one indicates a
// storage or migration bug, not a login failure — see the package doc for
// the full reasoning and its mirror image in hmac.verify_*.
func verifyPassword(plaintext, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	// Any other error (ErrHashTooShort, InvalidHashPrefixError,
	// HashVersionTooNewError, ...) means hash was not a valid bcrypt hash at
	// all — that is our data being malformed, not a login failure, so raise.
	return false, fmt.Errorf("password.verify: malformed hash: %w", err)
}

// costOf: returns the bcrypt cost embedded in hash. Raises on a malformed
// hash, for the same reason verify does.
func costOf(hash string) (int, error) {
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return 0, fmt.Errorf("password.cost: malformed hash: %w", err)
	}
	return cost, nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("password", "bcrypt password hashing and verification")

	m.Fn("hash", hashPassword,
		"hashes a password with bcrypt at the given cost; raises if cost is out of range or the password exceeds 72 bytes",
		luareg.Args("plaintext", "cost"),
		luareg.ArgDoc("plaintext", "the password to hash; must not exceed 72 bytes (bcrypt's own limit)"),
		luareg.ArgDoc("cost", "bcrypt cost factor in [4, 31]; use password.DEFAULT_COST unless you have a specific reason not to"),
		luareg.ReturnDoc(0, "hash", "a self-describing bcrypt hash ($2a$<cost>$...) suitable for storage"))
	m.Fn("verify", verifyPassword,
		"verifies a password against a stored bcrypt hash in constant time; returns false on a wrong password, raises on a malformed hash",
		luareg.Args("plaintext", "hash"),
		luareg.ArgDoc("plaintext", "the password attempt"),
		luareg.ArgDoc("hash", "a bcrypt hash previously produced by password.hash; this is application data, not attacker input, so a malformed value raises rather than returning false"),
		luareg.ReturnDoc(0, "ok", "true if plaintext matches hash"))
	m.Fn("cost", costOf,
		"returns the bcrypt cost embedded in a hash, for deciding whether to rehash on login; raises on a malformed hash",
		luareg.Args("hash"),
		luareg.ArgDoc("hash", "a bcrypt hash previously produced by password.hash"),
		luareg.ReturnDoc(0, "cost", "the cost factor the hash was created with"))

	m.Const("DEFAULT_COST", defaultCost, "number", "recommended bcrypt cost for production use (tests should use a much lower cost, e.g. 4, to stay fast)")

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("password", password.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
