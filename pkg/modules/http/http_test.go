// Copyright (c) 2024-2025 Thomas Maurice
// SPDX-License-Identifier: MIT

package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
)

func TestGet(t *testing.T) {
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
		local resp = http.get("` + server.URL + `", nil)
		assert(resp.status == 200, "Expected status 200, got " .. resp.status)
		assert(resp.body == '{"message":"hello"}', "Expected correct body")
	`
	require.NoError(t, L.DoString(code))
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
		local resp = http.post("` + server.URL + `", "test body", {["Content-Type"] = "text/plain"})
		assert(resp.status == 201, "Expected status 201")
	`
	require.NoError(t, L.DoString(code))
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
		local resp = http.get("` + server.URL + `", {["Authorization"] = "Bearer token123"})
		assert(resp.status == 200, "Expected status 200")
	`
	require.NoError(t, L.DoString(code))
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
		local resp = http.request("PATCH", "` + server.URL + `", "patch body", nil)
		assert(resp.status == 200, "Expected status 200")
	`
	require.NoError(t, L.DoString(code))
}

// TestInvalidURL: invalid URL raises a Lua error.
func TestInvalidURL(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("http", Loader)

	code := `
		local http = require("http")
		local ok, err = pcall(http.get, "not-a-valid-url", nil)
		assert(not ok, "Expected error for invalid URL")
		assert(type(err) == "string", "Error should be a string")
	`
	require.NoError(t, L.DoString(code))
}
