---@meta kubernetes

---@class kubernetes.GVKMatcher
---@field group string
---@field kind string
---@field version string

---@class kubernetes
local kubernetes = {}

--- parse a Kubernetes memory quantity, returns bytes
---@param quantity string
---@return number
function kubernetes.parse_memory(quantity) end

--- parse a Kubernetes CPU quantity, returns millicores
---@param quantity string
---@return number
function kubernetes.parse_cpu(quantity) end

--- parse an RFC3339 time string, returns Unix timestamp
---@param timestr string
---@return number
function kubernetes.parse_time(timestr) end

--- convert a Unix timestamp to RFC3339 string
---@param timestamp number
---@return string
function kubernetes.format_time(timestamp) end

--- parse a duration string, returns seconds
---@param duration string
---@return number
function kubernetes.parse_duration(duration) end

--- convert seconds to a duration string
---@param seconds number
---@return string
function kubernetes.format_duration(seconds) end

--- check if a Kubernetes object matches a GVK matcher
---@param obj table<string, any>
---@param matcher kubernetes.GVKMatcher
---@return boolean
function kubernetes.match_gvk(obj, matcher) end

--- ensure metadata.labels and annotations exist, returns updated obj
---@param obj table<string, any>
---@return table<string, any>
function kubernetes.ensure_metadata(obj) end

--- ensure metadata.labels and annotations exist, returns updated obj
---@param obj table<string, any>
---@return table<string, any>
function kubernetes.init_defaults(obj) end

--- add a label and return the updated obj
---@param obj table<string, any>
---@param key string
---@param value string
---@return table<string, any>
function kubernetes.add_label(obj, key, value) end

--- add multiple labels and return the updated obj
---@param obj table<string, any>
---@param labels table<string, any>
---@return table<string, any>
function kubernetes.add_labels(obj, labels) end

--- remove a label and return the updated obj
---@param obj table<string, any>
---@param key string
---@return table<string, any>
function kubernetes.remove_label(obj, key) end

--- return true if the label exists
---@param obj table<string, any>
---@param key string
---@return boolean
function kubernetes.has_label(obj, key) end

--- return the value of a label, or empty string if absent
---@param obj table<string, any>
---@param key string
---@return string
function kubernetes.get_label(obj, key) end

--- add an annotation and return the updated obj
---@param obj table<string, any>
---@param key string
---@param value string
---@return table<string, any>
function kubernetes.add_annotation(obj, key, value) end

--- add multiple annotations and return the updated obj
---@param obj table<string, any>
---@param annotations table<string, any>
---@return table<string, any>
function kubernetes.add_annotations(obj, annotations) end

--- remove an annotation and return the updated obj
---@param obj table<string, any>
---@param key string
---@return table<string, any>
function kubernetes.remove_annotation(obj, key) end

--- return true if the annotation exists
---@param obj table<string, any>
---@param key string
---@return boolean
function kubernetes.has_annotation(obj, key) end

--- return the value of an annotation, or empty string if absent
---@param obj table<string, any>
---@param key string
---@return string
function kubernetes.get_annotation(obj, key) end

return kubernetes
