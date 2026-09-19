---@meta compress

---@class compress
---@field NO_COMPRESSION number no compression, framing only
---@field BEST_SPEED number fastest compression level
---@field BEST_COMPRESSION number smallest output, slowest
---@field DEFAULT_COMPRESSION number the codec's own default level
---@field HUFFMAN_ONLY number Huffman-only compression, very fast, weak ratio
---@field MAX_BYTES_DEFAULT number a generous default max_bytes ceiling (64 MiB)
local compress = {}

--- compresses data as a gzip stream
---@param data string the raw bytes to compress
---@param level number compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants
---@return string compressed the gzip stream
function compress.gzip_compress(data, level) end

--- decompresses a gzip stream, raising if the output would exceed max_bytes
---@param data string the gzip stream to decompress
---@param max_bytes number the maximum number of decompressed bytes to allow; must be >= 1, no unlimited option
---@return string data the decompressed bytes
function compress.gzip_decompress(data, max_bytes) end

--- compresses data as a zlib stream
---@param data string the raw bytes to compress
---@param level number compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants
---@return string compressed the zlib stream
function compress.zlib_compress(data, level) end

--- decompresses a zlib stream, raising if the output would exceed max_bytes
---@param data string the zlib stream to decompress
---@param max_bytes number the maximum number of decompressed bytes to allow; must be >= 1, no unlimited option
---@return string data the decompressed bytes
function compress.zlib_decompress(data, max_bytes) end

--- compresses data as a raw DEFLATE stream (no gzip/zlib framing)
---@param data string the raw bytes to compress
---@param level number compression level in [-2, 9]; see the NO_COMPRESSION..HUFFMAN_ONLY constants
---@return string compressed the raw DEFLATE stream
function compress.flate_compress(data, level) end

--- decompresses a raw DEFLATE stream, raising if the output would exceed max_bytes
---@param data string the raw DEFLATE stream to decompress
---@param max_bytes number the maximum number of decompressed bytes to allow; must be >= 1, no unlimited option
---@return string data the decompressed bytes
function compress.flate_decompress(data, max_bytes) end

return compress
