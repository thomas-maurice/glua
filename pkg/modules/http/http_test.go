// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestGet(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"hello"}`))
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local resp, err = http.get("` + server.URL + `")
		assert(err == nil, "Expected no error: " .. tostring(err))
		assert(resp.status == 200, "Expected status 200, got " .. resp.status)
		assert(resp.body == '{"message":"hello"}', "Expected correct body")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"created"}`))
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local resp, err = http.post("` + server.URL + `", "test body", {["Content-Type"] = "text/plain"})
		assert(err == nil, "Expected no error")
		assert(resp.status == 201, "Expected status 201")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer token123" {
			t.Errorf("Expected Authorization header, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local resp, err = http.get("` + server.URL + `", {["Authorization"] = "Bearer token123"})
		assert(err == nil, "Expected no error")
		assert(resp.status == 200, "Expected status 200")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local resp, err = http.request("PATCH", "` + server.URL + `", "patch body")
		assert(err == nil, "Expected no error")
		assert(resp.status == 200, "Expected status 200")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}

func TestInvalidURL(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local resp, err = http.get("not-a-valid-url")
		assert(resp == nil, "Expected nil response")
		assert(err ~= nil, "Expected error")
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("Failed to execute Lua code: %v", err)
	}
}
