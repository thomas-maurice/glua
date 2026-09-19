-- Proves the byte path survives both modules: gzip-compress, base64-encode,
-- then reverse it, and the original string must come back unchanged.
local compress = require("compress")
local base64 = require("base64")

local payload = "the quick brown fox jumps over the lazy dog\0\1\2\255"

local blob = base64.encode(compress.gzip_compress(payload, compress.BEST_COMPRESSION))
local roundtripped = compress.gzip_decompress(base64.decode(blob), compress.MAX_BYTES_DEFAULT)

assert(roundtripped == payload, "round-trip mismatch")

return true
