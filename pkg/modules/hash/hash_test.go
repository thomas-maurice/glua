// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package hash

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestMD5(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.md5("hello world")
		assert(h == "5eb63bbbe01eeed093cb22bb8f5acdc3", "Expected correct MD5 hash, got: " .. h)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestSHA1(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.sha1("hello world")
		assert(h == "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed", "Expected correct SHA1 hash, got: " .. h)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestSHA256(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.sha256("hello world")
		assert(h == "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9", "Expected correct SHA256 hash, got: " .. h)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestSHA512(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.sha512("hello world")
		assert(type(h) == "string", "Expected string result")
		assert(#h == 128, "Expected 128 character hex string for SHA512, got length: " .. #h)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestHMACSHA256(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.hmac_sha256("message", "secret_key")
		assert(type(h) == "string", "Expected string result")
		assert(#h == 64, "Expected 64 character hex string for HMAC-SHA256, got length: " .. #h)
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestHMACSHA256_SameKeyProducesSame(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h1 = hash.hmac_sha256("message", "key")
		local h2 = hash.hmac_sha256("message", "key")
		assert(h1 == h2, "Expected same HMAC for same inputs")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestHMACSHA256_DifferentKeyProducesDifferent(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h1 = hash.hmac_sha256("message", "key1")
		local h2 = hash.hmac_sha256("message", "key2")
		assert(h1 ~= h2, "Expected different HMAC for different keys")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestEmptyString(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("hash", Loader)

	code := `
		local hash = require("hash")
		local h = hash.sha256("")
		assert(h == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "Expected correct SHA256 of empty string")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
