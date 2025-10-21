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

package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the http module for Lua.
// This function should be registered with L.PreloadModule("http", http.Loader)
//
// @luamodule http
//
// Example usage in Lua:
//
//	local http = require("http")
//	local response = http.get("https://api.example.com/data")
//	print(response.status)
//	print(response.body)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"get":     get,
	"post":    post,
	"put":     put,
	"delete":  delete_,
	"request": request,
}

// get: performs an HTTP GET request.
//
// @luafunc get
// @luaparam url string The URL to request
// @luaparam headers table|nil Optional headers table
// @luareturn table Response table with status, body, headers, or nil on error
// @luareturn string|nil Error message if request failed
//
// Example:
//
//	local resp, err = http.get("https://api.example.com/data", {
//	    ["Authorization"] = "Bearer token"
//	})
//	if err then
//	    print("Error: " .. err)
//	else
//	    print("Status: " .. resp.status)
//	    print("Body: " .. resp.body)
//	end
func get(L *lua.LState) int {
	url := L.CheckString(1)
	headers := L.OptTable(2, nil)
	return doRequest(L, "GET", url, "", headers)
}

// post: performs an HTTP POST request.
//
// @luafunc post
// @luaparam url string The URL to request
// @luaparam body string The request body
// @luaparam headers table|nil Optional headers table
// @luareturn table Response table with status, body, headers, or nil on error
// @luareturn string|nil Error message if request failed
//
// Example:
//
//	local json = require("json")
//	local resp, err = http.post(
//	    "https://api.example.com/data",
//	    json.stringify({key = "value"}),
//	    {["Content-Type"] = "application/json"}
//	)
func post(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)
	headers := L.OptTable(3, nil)
	return doRequest(L, "POST", url, body, headers)
}

// put: performs an HTTP PUT request.
//
// @luafunc put
// @luaparam url string The URL to request
// @luaparam body string The request body
// @luaparam headers table|nil Optional headers table
// @luareturn table Response table with status, body, headers, or nil on error
// @luareturn string|nil Error message if request failed
//
// Example:
//
//	local resp, err = http.put(
//	    "https://api.example.com/resource/123",
//	    json.stringify({key = "new_value"}),
//	    {["Content-Type"] = "application/json"}
//	)
func put(L *lua.LState) int {
	url := L.CheckString(1)
	body := L.CheckString(2)
	headers := L.OptTable(3, nil)
	return doRequest(L, "PUT", url, body, headers)
}

// delete_: performs an HTTP DELETE request.
//
// @luafunc delete
// @luaparam url string The URL to request
// @luaparam headers table|nil Optional headers table
// @luareturn table Response table with status, body, headers, or nil on error
// @luareturn string|nil Error message if request failed
//
// Example:
//
//	local resp, err = http.delete("https://api.example.com/resource/123")
func delete_(L *lua.LState) int {
	url := L.CheckString(1)
	headers := L.OptTable(2, nil)
	return doRequest(L, "DELETE", url, "", headers)
}

// request: performs a generic HTTP request with custom method.
//
// @luafunc request
// @luaparam method string The HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
// @luaparam url string The URL to request
// @luaparam body string|nil Optional request body
// @luaparam headers table|nil Optional headers table
// @luareturn table Response table with status, body, headers, or nil on error
// @luareturn string|nil Error message if request failed
//
// Example:
//
//	local resp, err = http.request("PATCH", "https://api.example.com/resource/123",
//	    json.stringify({key = "value"}),
//	    {["Content-Type"] = "application/json"}
//	)
func request(L *lua.LState) int {
	method := L.CheckString(1)
	url := L.CheckString(2)
	body := L.OptString(3, "")
	headers := L.OptTable(4, nil)
	return doRequest(L, method, url, body, headers)
}

// doRequest: internal function to perform the HTTP request
func doRequest(L *lua.LState, method, url, body string, headers *lua.LTable) int {
	// Create request
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create request: %v", err)))
		return 2
	}

	// Add headers
	if headers != nil {
		headers.ForEach(func(key lua.LValue, val lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				if valStr, ok := val.(lua.LString); ok {
					req.Header.Set(string(keyStr), string(valStr))
				}
			}
		})
	}

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("request failed: %v", err)))
		return 2
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read response body: %v", err)))
		return 2
	}

	// Build response table
	respTable := L.NewTable()
	respTable.RawSetString("status", lua.LNumber(resp.StatusCode))
	respTable.RawSetString("body", lua.LString(string(respBody)))

	// Add response headers
	headersTable := L.NewTable()
	for key, values := range resp.Header {
		if len(values) > 0 {
			headersTable.RawSetString(key, lua.LString(values[0]))
		}
	}
	respTable.RawSetString("headers", headersTable)

	L.Push(respTable)
	L.Push(lua.LNil)
	return 2
}
