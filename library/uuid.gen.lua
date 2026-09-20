---@meta uuid

---@class uuid.Info
---@field timestamp number Timestamp: Unix seconds, populated for v1/v6/v7 only; 0 for every other version. Unix seconds (rather than google/uuid's own Time type) matches the currency the time module already uses elsewhere in glua.
---@field uuid string canonical hyphenated form
---@field variant string e.g. "RFC4122"
---@field version number the UUID version (0 if unrecognised)

---@class uuid
---@field NIL string the nil UUID, 00000000-0000-0000-0000-000000000000
---@field NAMESPACE_DNS string RFC 4122 DNS namespace, for use with v5
---@field NAMESPACE_URL string RFC 4122 URL namespace, for use with v5
---@field NAMESPACE_OID string RFC 4122 OID namespace, for use with v5
---@field NAMESPACE_X500 string RFC 4122 X.500 namespace, for use with v5
local uuid = {}

--- generates a random (version 4) UUID
---@return string id canonical lowercase hyphenated UUID
function uuid.v4() end

--- generates a time-ordered (version 7) UUID; sorts by creation time and is monotonic within a millisecond
---@return string id canonical lowercase hyphenated UUID
function uuid.v7() end

--- generates a deterministic (version 5) UUID from a namespace UUID and a name
---@param namespace string a UUID string identifying the namespace; see the NAMESPACE_* constants
---@param name string the name to derive the UUID from
---@return string id canonical lowercase hyphenated UUID; identical inputs always produce the same output
function uuid.v5(namespace, name) end

--- parses a UUID string and reports its version, variant and timestamp
---@param s string a UUID in canonical, plain (no hyphens), urn:uuid:, or {braced} form
---@return uuid.Info info uuid.Info: uuid, version, variant, timestamp (Unix seconds, v1/v6/v7 only, 0 otherwise)
function uuid.parse(s) end

--- reports whether s parses as a valid UUID; never raises
---@param s string the string to validate
---@return boolean ok true if parse(s) would succeed
function uuid.is_valid(s) end

--- reformats a valid UUID string into the requested style
---@param s string a UUID in any form parse accepts
---@param style string one of: canonical, plain, urn, braced
---@return string formatted s reformatted into the requested style
function uuid.format(s, style) end

return uuid
