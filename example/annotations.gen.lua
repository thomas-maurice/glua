---@meta

--- custom_annotations module
---@alias ID string|number
---@alias Handler fun(id: ID): boolean
---@class custom_annotations
local custom_annotations = {}

--- processID: processes an ID value
---@param id ID the identifier to process
---@return boolean true if valid
---@deprecated Use process_typed_id instead
---@nodiscard
function custom_annotations.process_id(id) end

--- processTypedID: processes a typed ID value with generics
---@param id any the identifier to process
---@return boolean true if valid
---@generic T
---@param id `T`
---@return boolean
function custom_annotations.process_typed_id(id) end

--- json module
---@class json
local json = {}

--- parse: parses a JSON string and returns a Lua table. Returns nil and error message on failure.
---@param jsonstr string The JSON string to parse
---@return table The parsed JSON as a Lua table, or nil on error
---@return string|nil Error message if parsing failed
function json.parse(jsonstr) end

--- stringify: converts a Lua table to a JSON string. Returns nil and error message on failure.
---@param tbl table The Lua table to convert to JSON
---@return string The JSON string, or nil on error
---@return string|nil Error message if conversion failed
function json.stringify(tbl) end

--- k8sclient module
---@class k8sclient
local k8sclient = {}

--- get: retrieves a Kubernetes resource by GVK, namespace, and name.
---@param gvk GVKMatcher The GVK matcher with group, version, and kind
---@param namespace string The namespace of the resource
---@param name string The name of the resource
---@return table|nil The Kubernetes object, or nil on error
---@return string|nil Error message if retrieval failed
function k8sclient.get(gvk, namespace, name) end

--- create: creates a Kubernetes resource from a Lua table.
---@param obj table The Kubernetes object to create
---@return table|nil The created Kubernetes object, or nil on error
---@return string|nil Error message if creation failed
function k8sclient.create(obj) end

--- update: updates a Kubernetes resource.
---@param obj table The Kubernetes object to update
---@return table|nil The updated Kubernetes object, or nil on error
---@return string|nil Error message if update failed
function k8sclient.update(obj) end

--- delete: deletes a Kubernetes resource by GVK, namespace, and name.
---@param gvk GVKMatcher The GVK matcher with group, version, and kind
---@param namespace string The namespace of the resource
---@param name string The name of the resource
---@return string|nil Error message if deletion failed, nil on success
function k8sclient.delete(gvk, namespace, name) end

--- list: lists Kubernetes resources by GVK and namespace.
---@param gvk GVKMatcher The GVK matcher with group, version, and kind
---@param namespace string The namespace to list from
---@return table[]|nil Array of Kubernetes objects, or nil on error
---@return string|nil Error message if listing failed
function k8sclient.list(gvk, namespace) end

--- kubernetes module
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

--- initDefaults: initializes default empty tables for metadata.labels and metadata.annotations if they are nil. This is useful for ensuring these fields are tables instead of nil, making it easier to add labels/annotations in Lua without checking for nil first.
---@param obj table The Kubernetes object (must have a metadata field)
---@return table The same object with initialized defaults (modified in-place)
function kubernetes.init_defaults(obj) end

--- parseDuration: parses a Kubernetes duration string (e.g., "5s", "10m", "2h") and returns seconds as a number. Returns nil and error message on failure.
---@param duration string The duration string to parse (e.g., "5s", "10m", "2h")
---@return number The duration value in seconds, or nil on error
---@return string|nil Error message if parsing failed
function kubernetes.parse_duration(duration) end

--- formatDuration: converts seconds to a Kubernetes duration string. Returns nil and error message on failure.
---@param seconds number The duration in seconds to convert
---@return string The duration string (e.g., "5m0s", "1h30m0s"), or nil on error
---@return string|nil Error message if formatting failed
function kubernetes.format_duration(seconds) end

--- parseIntOrString: parses a Kubernetes IntOrString value and returns the value and its type. Returns (number, false) for integers or (string, true) for strings.
---@param value any The IntOrString value (number or string)
---@return any The parsed value (preserves type)
---@return boolean true if string, false if number
function kubernetes.parse_int_or_string(value) end

--- matchesSelector: checks if a set of labels matches a label selector. The selector is a table with key-value pairs that must all match.
---@param labels table The labels to check (e.g., pod.metadata.labels)
---@param selector table The selector with required labels
---@return boolean true if all selector labels match
function kubernetes.matches_selector(labels, selector) end

--- tolerationMatches: checks if a toleration matches a taint. Simplified matching: checks key, operator, value, and effect.
---@param toleration table The toleration object
---@param taint table The taint object
---@return boolean true if the toleration matches the taint
function kubernetes.toleration_matches(toleration, taint) end

--- matchGVK: checks if a Kubernetes object matches the specified Group/Version/Kind matcher. Returns true if the object's apiVersion and kind match the matcher's values.
---@param obj table The Kubernetes object to check
---@param matcher GVKMatcher The GVK matcher with group, version, and kind fields
---@return boolean true if the GVK matches
function kubernetes.match_gvk(obj, matcher) end

--- spew module
---@class spew
local spew = {}

--- dump: prints the contents of a Lua value to stdout with detailed formatting. This is useful for debugging and inspecting complex table structures.
---@param value any The Lua value to dump (table, string, number, etc.)
function spew.dump(value) end

--- sdump: returns a string representation of a Lua value with detailed formatting. Unlike dump, this returns the string instead of printing to stdout.
---@param value any The Lua value to dump (table, string, number, etc.)
---@return string A detailed string representation of the value
function spew.sdump(value) end

return {}
