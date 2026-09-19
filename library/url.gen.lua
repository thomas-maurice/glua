---@meta url

---@class url.URL
---@field fragment string
---@field host string
---@field hostname string
---@field opaque string
---@field password string
---@field path string
---@field port string
---@field raw_path string
---@field raw_query string
---@field scheme string
---@field username string

---@class url
local url = {}

--- parses a URL into its component parts, raises on failure
---@param s string the URL string to parse
---@return url.URL u url.URL table: scheme, opaque, username, password, host, hostname, port, path, raw_path, raw_query, fragment
function url.parse(s) end

--- reconstructs a URL string from a url.URL-shaped table, raises if the result is not a valid URL
---@param parts url.URL a table with the same shape parse returns; host is preferred verbatim when present, otherwise hostname+port are combined and hostname is auto-bracketed if it looks like an IPv6 literal
---@return string s the reconstructed URL string
function url.build(parts) end

--- resolves ref against base per RFC 3986 reference resolution, raises if either fails to parse
---@param base string the base URL
---@param ref string the reference to resolve against base; e.g. resolve("https://x/a/b", "c") -> "https://x/a/c", while resolve("https://x/a/b/", "c") -> "https://x/a/b/c"
---@return string s the resolved absolute URL string
function url.resolve(base, ref) end

--- escapes a string for safe inclusion in a URL query component, using + for space
---@param s string the raw string to escape
---@return string out the escaped string
function url.query_escape(s) end

--- the inverse of query_escape, raises on an invalid %-escape
---@param s string the escaped string
---@return string out the decoded raw string
function url.query_unescape(s) end

--- escapes a string for safe inclusion in a URL path segment, using %20 for space
---@param s string the raw string to escape
---@return string out the escaped string
function url.path_escape(s) end

--- the inverse of path_escape, raises on an invalid %-escape
---@param s string the escaped string
---@return string out the decoded raw string
function url.path_unescape(s) end

--- parses a URL query string, raises on malformed input
---@param raw string the raw query string, without a leading '?'
---@return table<string, string[]> values table<string, string[]>: every key maps to an array of values, even a single one; empty (never nil) for an empty query string
function url.parse_query(raw) end

--- encodes a table into a URL query string, sorted by key, raises on an invalid value shape
---@param t table<string, string|string[]> table<string, string|string[]>: a plain string for a single value, or an array of strings for a repeated key
---@return string raw the encoded query string, with keys sorted (net/url.Values.Encode's own behaviour, kept for deterministic output)
function url.build_query(t) end

return url
