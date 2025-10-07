---@meta

---@class k8sclient
local k8sclient = {}

--- newClientLua: creates a new client instance in Lua
---@return table The client instance with methods: get, create, update, delete, list
---@return string|nil Error message if client creation failed
function k8sclient.new_client() end

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
