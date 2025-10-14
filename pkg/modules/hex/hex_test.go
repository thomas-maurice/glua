// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hex

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestEncode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local encoded = hex.encode("hello")
		assert(encoded == "68656c6c6f", "Expected correct hex encoding, got: " .. encoded)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDecode(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local decoded, err = hex.decode("68656c6c6f")
		assert(err == nil, "Expected no error")
		assert(decoded == "hello", "Expected 'hello', got: " .. decoded)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRoundTrip(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local original = "The quick brown fox"
		local encoded = hex.encode(original)
		local decoded, err = hex.decode(encoded)
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

	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local decoded, err = hex.decode("not valid hex!")
		assert(decoded == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestDecodeOddLength(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hex", Loader)

	code := `
		local hex = require("hex")
		local decoded, err = hex.decode("123")  -- odd length
		assert(decoded == nil, "Expected nil result")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
