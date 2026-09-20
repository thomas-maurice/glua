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

// Package uuid provides UUID v4, v5 and v7 generation, parsing and
// formatting for Lua scripts, backed by github.com/google/uuid.
//
// This is deliberately a separate module from random, not
// random.uuid_v4()/random.uuid_v7(): v7 is *time-ordered*, not random —
// consecutive v7 values sort by creation time by design — so presenting it
// next to random.token() would teach the wrong mental model. parse/format/
// is_valid have nothing to do with randomness either.
//
// The google/uuid dependency exists specifically for v7: RFC 9562 requires
// a monotonic counter in the sub-millisecond bits so two UUIDs minted in the
// same millisecond still sort correctly, which needs real synchronisation.
// v4 alone would not have justified the dependency.
//
// parse/is_valid/v5's namespace argument accept every form
// github.com/google/uuid's Parse accepts: canonical hyphenated
// (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx), the raw 32-hex form with no
// hyphens, urn:uuid:..., and the braced {xxxxxxxx-...-xxxxxxxxxxxx} form.
// format converts a parsed UUID between the "canonical", "plain", "urn" and
// "braced" spellings.
package uuid

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// Info: everything parse reports about a UUID string.
type Info struct {
	UUID    string `json:"uuid"`    // canonical hyphenated form
	Version int    `json:"version"` // the UUID version (0 if unrecognised)
	Variant string `json:"variant"` // e.g. "RFC4122"

	// Timestamp: Unix seconds, populated for v1/v6/v7 only; 0 for every
	// other version. Unix seconds (rather than google/uuid's own Time type)
	// matches the currency the time module already uses elsewhere in glua.
	Timestamp float64 `json:"timestamp"`
}

// v4: generates a random (version 4) UUID.
func v4() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("uuid.v4: %w", err)
	}
	return id.String(), nil
}

// v7: generates a time-ordered (version 7) UUID. Values minted within the
// same millisecond are still monotonically increasing, courtesy of
// google/uuid's internal counter — this ordering guarantee is the entire
// reason this module takes google/uuid as a dependency.
func v7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("uuid.v7: %w", err)
	}
	return id.String(), nil
}

// v5: generates a deterministic (version 5, SHA-1 name-based) UUID from a
// namespace UUID and a name. The same (namespace, name) pair always
// produces the same UUID, which is useful for deriving a stable ID from a
// Kubernetes object name without a database.
func v5(namespace, name string) (string, error) {
	ns, err := uuid.Parse(namespace)
	if err != nil {
		return "", fmt.Errorf("uuid.v5: invalid namespace: %w", err)
	}
	return uuid.NewSHA1(ns, []byte(name)).String(), nil
}

// v7TimestampSeconds: reads the leading 48 bits of a v7 UUID directly as
// Unix milliseconds (see RFC 9562 layout: unix_ts_ms occupies uuid[0:6]).
// google/uuid's own UUID.Time() is oriented around v1's Gregorian-epoch
// counter representation; extracting the v7 timestamp directly is simpler
// and more obviously correct than routing through it.
func v7TimestampSeconds(id uuid.UUID) float64 {
	ms := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
		int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
	return float64(ms) / 1000.0
}

// timestampSeconds: returns id's embedded timestamp as Unix seconds, or 0
// if id's version does not carry a timestamp.
func timestampSeconds(id uuid.UUID) float64 {
	switch id.Version() {
	case 7:
		return v7TimestampSeconds(id)
	case 1, 2, 6:
		sec, _ := id.Time().UnixTime()
		return float64(sec)
	default:
		return 0
	}
}

// parse: parses s (any of the canonical, plain, urn or braced forms) and
// reports its version, variant and (for v1/v6/v7) timestamp.
func parse(s string) (Info, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return Info{}, fmt.Errorf("uuid.parse: %w", err)
	}
	return Info{
		UUID:      id.String(),
		Version:   int(id.Version()),
		Variant:   id.Variant().String(),
		Timestamp: timestampSeconds(id),
	}, nil
}

// isValid: reports whether s parses as a valid UUID. Never raises — this is
// the non-raising predicate companion to parse, for callers who want a
// boolean instead of a pcall.
func isValid(s string) bool {
	return uuid.Validate(s) == nil
}

// format: reformats a valid UUID string into the requested style.
func format(s, style string) (string, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", fmt.Errorf("uuid.format: %w", err)
	}
	switch style {
	case "canonical":
		return id.String(), nil
	case "plain":
		return strings.ReplaceAll(id.String(), "-", ""), nil
	case "urn":
		return id.URN(), nil
	case "braced":
		return "{" + id.String() + "}", nil
	default:
		return "", fmt.Errorf("uuid.format: unknown style %q (want canonical, plain, urn, or braced)", style)
	}
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("uuid", "UUID v4/v5/v7 generation, parsing and formatting")

	m.Fn("v4", v4, "generates a random (version 4) UUID",
		luareg.ReturnDoc(0, "id", "canonical lowercase hyphenated UUID"))

	m.Fn("v7", v7, "generates a time-ordered (version 7) UUID; sorts by creation time and is monotonic within a millisecond",
		luareg.ReturnDoc(0, "id", "canonical lowercase hyphenated UUID"))

	m.Fn("v5", v5, "generates a deterministic (version 5) UUID from a namespace UUID and a name",
		luareg.Args("namespace", "name"),
		luareg.ArgDoc("namespace", "a UUID string identifying the namespace; see the NAMESPACE_* constants"),
		luareg.ArgDoc("name", "the name to derive the UUID from"),
		luareg.ReturnDoc(0, "id", "canonical lowercase hyphenated UUID; identical inputs always produce the same output"))

	m.Fn("parse", parse, "parses a UUID string and reports its version, variant and timestamp",
		luareg.Args("s"),
		luareg.ArgDoc("s", "a UUID in canonical, plain (no hyphens), urn:uuid:, or {braced} form"),
		luareg.ReturnDoc(0, "info", "uuid.Info: uuid, version, variant, timestamp (Unix seconds, v1/v6/v7 only, 0 otherwise)"))

	m.Fn("is_valid", isValid, "reports whether s parses as a valid UUID; never raises",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the string to validate"),
		luareg.ReturnDoc(0, "ok", "true if parse(s) would succeed"))

	m.Fn("format", format, "reformats a valid UUID string into the requested style",
		luareg.Args("s", "style"),
		luareg.ArgDoc("s", "a UUID in any form parse accepts"),
		luareg.ArgDoc("style", "one of: canonical, plain, urn, braced"),
		luareg.ReturnDoc(0, "formatted", "s reformatted into the requested style"))

	m.Const("NIL", uuid.Nil.String(), "string", "the nil UUID, 00000000-0000-0000-0000-000000000000")
	m.Const("NAMESPACE_DNS", uuid.NameSpaceDNS.String(), "string", "RFC 4122 DNS namespace, for use with v5")
	m.Const("NAMESPACE_URL", uuid.NameSpaceURL.String(), "string", "RFC 4122 URL namespace, for use with v5")
	m.Const("NAMESPACE_OID", uuid.NameSpaceOID.String(), "string", "RFC 4122 OID namespace, for use with v5")
	m.Const("NAMESPACE_X500", uuid.NameSpaceX500.String(), "string", "RFC 4122 X.500 namespace, for use with v5")

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("uuid", uuid.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
