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

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// Response: represents an HTTP response returned to Lua callers.
type Response struct {
	Status  int               `json:"status"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers"`
}

// doRequest: executes an HTTP request and returns a *Response or an error.
// The optional headers table is read directly from the Lua stack via L.
func doRequest(L *lua.LState, method, url, body string, headers *lua.LTable) (*Response, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if headers != nil {
		headers.ForEach(func(key lua.LValue, val lua.LValue) {
			if keyStr, ok := key.(lua.LString); ok {
				if valStr, ok := val.(lua.LString); ok {
					req.Header.Set(string(keyStr), string(valStr))
				}
			}
		})
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	hdrs := make(map[string]string, len(resp.Header))
	for k, vals := range resp.Header {
		if len(vals) > 0 {
			hdrs[k] = vals[0]
		}
	}

	return &Response{
		Status:  resp.StatusCode,
		Body:    string(respBody),
		Headers: hdrs,
	}, nil
}

// responseToLua: converts a Response to a Lua table. Used by all request
// wrappers to produce a consistent response shape.
func responseToLua(L *lua.LState, r *Response) *lua.LTable {
	tbl := L.NewTable()
	tbl.RawSetString("status", lua.LNumber(r.Status))
	tbl.RawSetString("body", lua.LString(r.Body))
	hdrs := L.NewTable()
	for k, v := range r.Headers {
		hdrs.RawSetString(k, lua.LString(v))
	}
	tbl.RawSetString("headers", hdrs)
	return tbl
}

// get: performs an HTTP GET request with optional headers table.
func get(L *lua.LState, url string, headers *lua.LTable) (lua.LValue, error) {
	r, err := doRequest(L, "GET", url, "", headers)
	if err != nil {
		return lua.LNil, err
	}
	return responseToLua(L, r), nil
}

// post: performs an HTTP POST request with optional headers table.
func post(L *lua.LState, url, body string, headers *lua.LTable) (lua.LValue, error) {
	r, err := doRequest(L, "POST", url, body, headers)
	if err != nil {
		return lua.LNil, err
	}
	return responseToLua(L, r), nil
}

// put: performs an HTTP PUT request with optional headers table.
func put(L *lua.LState, url, body string, headers *lua.LTable) (lua.LValue, error) {
	r, err := doRequest(L, "PUT", url, body, headers)
	if err != nil {
		return lua.LNil, err
	}
	return responseToLua(L, r), nil
}

// delete_: performs an HTTP DELETE request with optional headers table.
func delete_(L *lua.LState, url string, headers *lua.LTable) (lua.LValue, error) {
	r, err := doRequest(L, "DELETE", url, "", headers)
	if err != nil {
		return lua.LNil, err
	}
	return responseToLua(L, r), nil
}

// request: performs an HTTP request with a custom method.
func request(L *lua.LState, method, url, body string, headers *lua.LTable) (lua.LValue, error) {
	r, err := doRequest(L, method, url, body, headers)
	if err != nil {
		return lua.LNil, err
	}
	return responseToLua(L, r), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("http", "HTTP client utilities")
	m.Fn("get", get, "performs an HTTP GET request, raises on network error",
		luareg.Args("url", "headers"))
	m.Fn("post", post, "performs an HTTP POST request, raises on network error",
		luareg.Args("url", "body", "headers"))
	m.Fn("put", put, "performs an HTTP PUT request, raises on network error",
		luareg.Args("url", "body", "headers"))
	m.Fn("delete", delete_, "performs an HTTP DELETE request, raises on network error",
		luareg.Args("url", "headers"))
	m.Fn("request", request, "performs an HTTP request with a custom method, raises on network error",
		luareg.Args("method", "url", "body", "headers"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("http", http.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
