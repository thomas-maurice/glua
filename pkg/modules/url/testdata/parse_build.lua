-- parse: every component populated, plus build round-trips.
local url = require("url")

local u = url.parse("https://user:pw@example.com:8443/a/b?x=1&x=2#frag")
assert(u.scheme == "https", "scheme")
assert(u.username == "user", "username")
assert(u.password == "pw", "password")
assert(u.hostname == "example.com", "hostname")
assert(u.port == "8443", "port")
assert(u.path == "/a/b", "path")
assert(u.raw_query == "x=1&x=2", "raw_query")
assert(u.fragment == "frag", "fragment")

local rebuilt = url.build(u)
assert(rebuilt == "https://user:pw@example.com:8443/a/b?x=1&x=2#frag", "round-trip: got " .. rebuilt)

-- IPv6 literal host: the D5 trap. hostname must be unbracketed, port separate.
local v6 = url.parse("http://[::1]:8080/path")
assert(v6.hostname == "::1", "v6 hostname is unbracketed: got " .. v6.hostname)
assert(v6.port == "8080", "v6 port")
assert(url.build(v6) == "http://[::1]:8080/path", "v6 round-trip via build")

-- build must bracket a bare v6 literal automatically, even without a "host" field.
local manual = url.build({scheme = "http", hostname = "::1", port = "8080", path = "/path"})
assert(manual == "http://[::1]:8080/path", "build auto-brackets a bare v6 hostname: got " .. manual)

-- Zone id: must survive a round-trip without corruption.
local zoned = url.parse("http://[fe80::1%25eth0]/")
assert(url.build(zoned) == "http://[fe80::1%25eth0]/", "zone id round-trips")

-- Opaque URL (mailto:) has no host.
local mail = url.parse("mailto:x@y.com")
assert(mail.scheme == "mailto", "mailto scheme")
assert(mail.opaque == "x@y.com", "mailto opaque")
assert(mail.host == "", "mailto has no host")

-- Scheme-relative URL.
local rel = url.parse("//example.com/a")
assert(rel.scheme == "", "scheme-relative has empty scheme")
assert(rel.host == "example.com", "scheme-relative host")

return true
