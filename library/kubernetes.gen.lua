---@meta kubernetes

---@class kubernetes
local kubernetes = {}

---@param quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
---@return number The memory value in bytes, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_memory(quantity) end

---@param quantity string The CPU quantity to parse (e.g., "100m", "1", "2000m")
---@return number The CPU value in millicores, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_cpu(quantity) end

---@param timestr string The time string in RFC3339 format (e.g., "2025-10-03T16:39:00Z")
---@return number The Unix timestamp, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_time(timestr) end

---@param timestamp number The Unix timestamp to convert
---@return string The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_time(timestamp) end

---@param obj table The Kubernetes object (must have a metadata field)
---@return table The same object with initialized defaults (modified in-place)
function kubernetes.init_defaults(obj) end

---@param duration string The duration string to parse (e.g., "5s", "10m", "2h")
---@return number The duration value in seconds, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_duration(duration) end

---@param seconds number The duration in seconds to convert
---@return string The duration string (e.g., "5m0s", "1h30m0s"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_duration(seconds) end

---@param value any The IntOrString value (number or string)
---@return any The parsed value (preserves type)
---@return boolean true if string, false if number
function kubernetes.parse_int_or_string(value) end

---@param labels table The labels to check (e.g., pod.metadata.labels)
---@param selector table The selector with required labels
---@return boolean true if all selector labels match
function kubernetes.matches_selector(labels, selector) end

---@param obj table The Kubernetes object to check
---@param matcher GVKMatcher The GVK matcher with group, version, and kind fields
---@return boolean true if the GVK matches
function kubernetes.match_gvk(obj, matcher) end

return kubernetes
