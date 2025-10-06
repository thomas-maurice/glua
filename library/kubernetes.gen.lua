---@meta

---@class kubernetes
local kubernetes = {}

--- parseMemory: parses a Kubernetes memory quantity (e.g., "1024Mi", "1Gi", "512M") and returns bytes as a number. Returns nil and error message on failure.
---@param quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
---@return number The memory value in bytes, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_memory(quantity) end

--- parseCPU: parses a Kubernetes CPU quantity (e.g., "100m", "1", "2000m") and returns millicores as a number. Returns nil and error message on failure.
---@param quantity string The CPU quantity to parse (e.g., "100m", "1", "2000m")
---@return number The CPU value in millicores, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_cpu(quantity) end

--- parseTime: parses a Kubernetes time string (RFC3339 format like "2025-10-03T16:39:00Z") and returns a Unix timestamp. Returns nil and error message on failure.
---@param timestr string The time string in RFC3339 format (e.g., "2025-10-03T16:39:00Z")
---@return number The Unix timestamp, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_time(timestr) end

--- formatTime: converts a Unix timestamp (int64) to a Kubernetes time string in RFC3339 format. Returns nil and error message on failure.
---@param timestamp number The Unix timestamp to convert
---@return string The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_time(timestamp) end

return kubernetes
