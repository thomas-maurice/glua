-- max_bytes < 1 must raise, naming the limit, with no "0 means unlimited"
-- escape hatch.
local compress = require("compress")

local blob = compress.gzip_compress("hello world", compress.DEFAULT_COMPRESSION)

local ok, err = pcall(compress.gzip_decompress, blob, 0)
assert(not ok, "expected max_bytes = 0 to raise")
assert(string.find(err, "max_bytes"), "error should mention max_bytes, got: " .. tostring(err))

local ok2, err2 = pcall(compress.gzip_decompress, blob, -1)
assert(not ok2, "expected negative max_bytes to raise")
assert(string.find(err2, "max_bytes"), "error should mention max_bytes, got: " .. tostring(err2))

return true
