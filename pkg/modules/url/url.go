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

// Package url provides URL parsing, construction, escaping and RFC 3986
// reference resolution for Lua scripts, built entirely on net/url (which
// does not touch the network or the resolver, so the minimal-IO constraint
// holds).
//
// IPv6 literal hosts are the trap this module is designed around:
//
//   - parse("http://[::1]:8080/path") reports hostname "::1" (unbracketed,
//     matching net/url.URL.Hostname()) and port "8080" separately from host
//     "[::1]:8080" (the raw authority, matching net/url.URL.Host verbatim).
//   - build re-adds brackets around a bare IPv6 literal automatically when a
//     caller constructs a URL table from scratch with only hostname/port
//     set — see buildHost. When the host field from a prior parse is
//     present, build uses it directly, which is what makes
//     build(parse(u)) == u for every case below, including a zone id: the
//     hostname field holds the zone's '%' decoded (net/url.Hostname()'s own
//     behaviour), but host and build's fallback path both round-trip
//     through net/url.URL.String(), which re-escapes '%' back to '%25' when
//     rendering a bracketed IPv6 host.
//
// The url.URL table shape returned by parse and accepted by build is:
// scheme, opaque, username, password, host, hostname, port, path, raw_path,
// raw_query, fragment. username/password are flattened out of
// net/url.Userinfo into two plain strings, empty when absent. raw_path is
// net/url.URL.EscapedPath() (the percent-encoded path as it appears in the
// URL) — build always derives the encoded path from path (the decoded
// form), so a raw_path that used non-default escaping is not preserved
// byte-for-byte across a build round-trip; this is a documented limitation,
// not a bug, and does not affect the IPv6/zone guarantees above.
package url

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// URL: the table shape both parse returns and build accepts. See the
// package doc for the host vs hostname/port split and the raw_path caveat.
type URL struct {
	Scheme   string `json:"scheme"`
	Opaque   string `json:"opaque"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Path     string `json:"path"`
	RawPath  string `json:"raw_path"`
	RawQuery string `json:"raw_query"`
	Fragment string `json:"fragment"`
}

// parseURL: parses s into its component parts. Raises on anything net/url
// cannot parse.
func parseURL(s string) (URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return URL{}, fmt.Errorf("url.parse: %w", err)
	}
	var username, password string
	if u.User != nil {
		username = u.User.Username()
		password, _ = u.User.Password()
	}
	return URL{
		Scheme:   u.Scheme,
		Opaque:   u.Opaque,
		Username: username,
		Password: password,
		Host:     u.Host,
		Hostname: u.Hostname(),
		Port:     u.Port(),
		Path:     u.Path,
		RawPath:  u.EscapedPath(),
		RawQuery: u.RawQuery,
		Fragment: u.Fragment,
	}, nil
}

// buildHost: derives the Host field for the reconstructed URL. Prefers the
// caller-supplied host verbatim (the common case: a table came from parse,
// so host is already correctly bracketed/escaped) and falls back to
// composing hostname+port, auto-bracketing hostname if it looks like an
// IPv6 literal, so a caller who only knows about hostname/port never has to
// learn the bracket rule.
func buildHost(parts URL) string {
	if parts.Host != "" {
		return parts.Host
	}
	if parts.Hostname == "" {
		return ""
	}
	host := parts.Hostname
	if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
		host = "[" + host + "]"
	}
	if parts.Port != "" {
		host += ":" + parts.Port
	}
	return host
}

// buildURL: reconstructs a URL string from parts, the same table shape
// parse returns. Raises if the resulting string does not itself parse back
// as a valid URL.
func buildURL(parts URL) (string, error) {
	u := &url.URL{
		Scheme:   parts.Scheme,
		Opaque:   parts.Opaque,
		Host:     buildHost(parts),
		Path:     parts.Path,
		RawQuery: parts.RawQuery,
		Fragment: parts.Fragment,
	}
	if parts.Username != "" || parts.Password != "" {
		if parts.Password != "" {
			u.User = url.UserPassword(parts.Username, parts.Password)
		} else {
			u.User = url.User(parts.Username)
		}
	}
	result := u.String()
	if _, err := url.Parse(result); err != nil {
		return "", fmt.Errorf("url.build: constructed URL %q is not valid: %w", result, err)
	}
	return result, nil
}

// resolve: resolves ref against base per RFC 3986 reference resolution
// (net/url.URL.ResolveReference). Raises if either argument fails to parse.
func resolve(base, ref string) (string, error) {
	b, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("url.resolve: invalid base: %w", err)
	}
	r, err := url.Parse(ref)
	if err != nil {
		return "", fmt.Errorf("url.resolve: invalid ref: %w", err)
	}
	return b.ResolveReference(r).String(), nil
}

// queryEscape: escapes s for safe inclusion in a URL query component,
// using '+' for space. Never raises.
func queryEscape(s string) string {
	return url.QueryEscape(s)
}

// queryUnescape: the inverse of query_escape. Raises on an invalid '%' escape.
func queryUnescape(s string) (string, error) {
	out, err := url.QueryUnescape(s)
	if err != nil {
		return "", fmt.Errorf("url.query_unescape: %w", err)
	}
	return out, nil
}

// pathEscape: escapes s for safe inclusion in a URL path segment, using
// "%20" for space. Never raises.
func pathEscape(s string) string {
	return url.PathEscape(s)
}

// pathUnescape: the inverse of path_escape. Raises on an invalid '%' escape.
func pathUnescape(s string) (string, error) {
	out, err := url.PathUnescape(s)
	if err != nil {
		return "", fmt.Errorf("url.path_unescape: %w", err)
	}
	return out, nil
}

// parseQuery: parses raw as a URL query string. Every key maps to an array
// of values, even a key with a single value, because a repeated key
// ("?a=1&a=2") is legal and a table<string, string> shape would silently
// drop a value. A bare key with no "=" ("?flag") maps to a one-element
// array containing the empty string, matching net/url.ParseQuery. Returns
// an empty (non-nil) table for an empty query string. Raises on a malformed
// query.
func parseQuery(raw string) (map[string][]string, error) {
	v, err := url.ParseQuery(raw)
	if err != nil {
		return nil, fmt.Errorf("url.parse_query: %w", err)
	}
	return map[string][]string(v), nil
}

// buildQuery: encodes t (a table<string, string|string[]> — a plain string
// for a single value, or an array of strings for a repeated key) into a URL
// query string. The output is sorted by key, because net/url.Values.Encode
// sorts — this makes the output deterministic, which callers rely on to
// hash or compare it, but means insertion order in t is not preserved.
// Raises if a value is neither a string nor an array of strings, or if a
// key is not a string.
func buildQuery(tbl *lua.LTable) (string, error) {
	values := url.Values{}
	var iterErr error
	tbl.ForEach(func(k, v lua.LValue) {
		if iterErr != nil {
			return
		}
		key, ok := k.(lua.LString)
		if !ok {
			iterErr = fmt.Errorf("url.build_query: table keys must be strings, got %s", k.Type())
			return
		}
		switch vv := v.(type) {
		case lua.LString:
			values.Add(string(key), string(vv))
		case *lua.LTable:
			n := vv.MaxN()
			for i := 1; i <= n; i++ {
				elem := vv.RawGetInt(i)
				s, ok := elem.(lua.LString)
				if !ok {
					iterErr = fmt.Errorf("url.build_query: key %q: array values must be strings, got %s at index %d", string(key), elem.Type(), i)
					return
				}
				values.Add(string(key), string(s))
			}
		default:
			iterErr = fmt.Errorf("url.build_query: key %q: value must be a string or an array of strings, got %s", string(key), v.Type())
		}
	})
	if iterErr != nil {
		return "", iterErr
	}
	return values.Encode(), nil
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("url", "URL parsing, construction, escaping and RFC 3986 reference resolution")

	m.Fn("parse", parseURL, "parses a URL into its component parts, raises on failure",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the URL string to parse"),
		luareg.ReturnDoc(0, "u", "url.URL table: scheme, opaque, username, password, host, hostname, port, path, raw_path, raw_query, fragment"))
	m.Fn("build", buildURL, "reconstructs a URL string from a url.URL-shaped table, raises if the result is not a valid URL",
		luareg.Args("parts"),
		luareg.ArgDoc("parts", "a table with the same shape parse returns; host is preferred verbatim when present, otherwise hostname+port are combined and hostname is auto-bracketed if it looks like an IPv6 literal"),
		luareg.ReturnDoc(0, "s", "the reconstructed URL string"))
	m.Fn("resolve", resolve, "resolves ref against base per RFC 3986 reference resolution, raises if either fails to parse",
		luareg.Args("base", "ref"),
		luareg.ArgDoc("base", "the base URL"),
		luareg.ArgDoc("ref", "the reference to resolve against base; e.g. resolve(\"https://x/a/b\", \"c\") -> \"https://x/a/c\", while resolve(\"https://x/a/b/\", \"c\") -> \"https://x/a/b/c\""),
		luareg.ReturnDoc(0, "s", "the resolved absolute URL string"))
	m.Fn("query_escape", queryEscape, "escapes a string for safe inclusion in a URL query component, using + for space",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to escape"),
		luareg.ReturnDoc(0, "out", "the escaped string"))
	m.Fn("query_unescape", queryUnescape, "the inverse of query_escape, raises on an invalid %-escape",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the escaped string"),
		luareg.ReturnDoc(0, "out", "the decoded raw string"))
	m.Fn("path_escape", pathEscape, "escapes a string for safe inclusion in a URL path segment, using %20 for space",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the raw string to escape"),
		luareg.ReturnDoc(0, "out", "the escaped string"))
	m.Fn("path_unescape", pathUnescape, "the inverse of path_escape, raises on an invalid %-escape",
		luareg.Args("s"),
		luareg.ArgDoc("s", "the escaped string"),
		luareg.ReturnDoc(0, "out", "the decoded raw string"))
	m.Fn("parse_query", parseQuery, "parses a URL query string, raises on malformed input",
		luareg.Args("raw"),
		luareg.ArgDoc("raw", "the raw query string, without a leading '?'"),
		luareg.ReturnDoc(0, "values", "table<string, string[]>: every key maps to an array of values, even a single one; empty (never nil) for an empty query string"),
		luareg.ReturnType(0, "table<string, string[]>"))
	m.Fn("build_query", buildQuery, "encodes a table into a URL query string, sorted by key, raises on an invalid value shape",
		luareg.Args("t"),
		luareg.ArgDoc("t", "table<string, string|string[]>: a plain string for a single value, or an array of strings for a repeated key"),
		luareg.ArgType("t", "table<string, string|string[]>"),
		luareg.ReturnDoc(0, "raw", "the encoded query string, with keys sorted (net/url.Values.Encode's own behaviour, kept for deterministic output)"))

	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("url", url.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
