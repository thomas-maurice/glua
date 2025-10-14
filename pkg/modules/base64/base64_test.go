// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package base64

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestEncode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local encoded = base64.encode("hello world")
		assert(encoded == "aGVsbG8gd29ybGQ=", "Expected correct base64 encoding, got: " .. encoded)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDecode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local decoded, err = base64.decode("aGVsbG8gd29ybGQ=")
		assert(err == nil, "Expected no error")
		assert(decoded == "hello world", "Expected 'hello world', got: " .. decoded)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local original = "The quick brown fox jumps over the lazy dog"
		local encoded = base64.encode(original)
		local decoded, err = base64.decode(encoded)
		assert(err == nil, "Expected no error")
		assert(decoded == original, "Expected round-trip to match")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDecodeInvalid(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local decoded, err = base64.decode("not!valid!base64!")
		assert(decoded == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestEncodeURL(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local encoded = base64.encode_url("hello world")
		-- URL encoding uses - and _ instead of + and /
		assert(type(encoded) == "string", "Expected string result")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDecodeURL(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("base64", Loader)

	code := `
		local base64 = require("base64")
		local original = "test data"
		local encoded = base64.encode_url(original)
		local decoded, err = base64.decode_url(encoded)
		assert(err == nil, "Expected no error")
		assert(decoded == original, "Expected URL-safe round-trip to match")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
