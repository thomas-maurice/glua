# k8sclient API

The `k8sclient` module provides a Kubernetes dynamic client for Lua, allowing you to interact with Kubernetes resources from Lua scripts.

## Initialization

The module is initialized from Go using a `*rest.Config`:

```go
import (
    "github.com/thomas-maurice/glua/pkg/modules/k8sclient"
    "k8s.io/client-go/tools/clientcmd"
    lua "github.com/yuin/gopher-lua"
)

config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
L := lua.NewState()
L.PreloadModule("k8sclient", k8sclient.Loader(config))
```

## API Functions

All functions use the GVKMatcher type for specifying resource types:

```lua
---@class GVKMatcher
---@field group string The API group (empty string for core resources)
---@field version string The API version (e.g., "v1")
---@field kind string The resource kind (e.g., "Pod", "ConfigMap")
```

### `get(gvk, namespace, name)`

Retrieves a single Kubernetes resource.

**Parameters:**
- `gvk` (GVKMatcher): The resource type to get
- `namespace` (string): The namespace of the resource
- `name` (string): The name of the resource

**Returns:**
- `table|nil`: The Kubernetes object as a Lua table, or nil on error
- `string|nil`: Error message if retrieval failed

**Example:**
```lua
local client = require("k8sclient")
local gvk = {group = "", version = "v1", kind = "Pod"}
local pod, err = client.get(gvk, "default", "my-pod")

if err then
    error("Failed to get pod: " .. err)
end

print("Pod name: " .. pod.metadata.name)
```

### `create(object)`

Creates a new Kubernetes resource.

**Parameters:**
- `object` (table): The Kubernetes object to create (must have `apiVersion`, `kind`, and `metadata` fields)

**Returns:**
- `table|nil`: The created object with server-generated fields, or nil on error
- `string|nil`: Error message if creation failed

**Example:**
```lua
local client = require("k8sclient")

local configmap = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {
        name = "my-config",
        namespace = "default"
    },
    data = {
        key1 = "value1",
        key2 = "value2"
    }
}

local created, err = client.create(configmap)

if err then
    error("Failed to create ConfigMap: " .. err)
end

print("Created ConfigMap with UID: " .. created.metadata.uid)
```

### `update(object)`

Updates an existing Kubernetes resource.

**Parameters:**
- `object` (table): The Kubernetes object to update (must include `metadata.resourceVersion`)

**Returns:**
- `table|nil`: The updated object, or nil on error
- `string|nil`: Error message if update failed

**Example:**
```lua
local client = require("k8sclient")
local gvk = {group = "", version = "v1", kind = "ConfigMap"}

-- Get the resource
local cm, err = client.get(gvk, "default", "my-config")
if err then error(err) end

-- Modify it
cm.data.key3 = "value3"

-- Update it
local updated, err = client.update(cm)
if err then
    error("Failed to update ConfigMap: " .. err)
end
```

### `delete(gvk, namespace, name)`

Deletes a Kubernetes resource.

**Parameters:**
- `gvk` (GVKMatcher): The resource type to delete
- `namespace` (string): The namespace of the resource
- `name` (string): The name of the resource

**Returns:**
- `string|nil`: Error message if deletion failed, or nil on success

**Example:**
```lua
local client = require("k8sclient")
local gvk = {group = "", version = "v1", kind = "ConfigMap"}

local err = client.delete(gvk, "default", "my-config")

if err then
    error("Failed to delete ConfigMap: " .. err)
end

print("ConfigMap deleted successfully")
```

### `list(gvk, namespace)`

Lists all resources of a given type in a namespace.

**Parameters:**
- `gvk` (GVKMatcher): The resource type to list
- `namespace` (string): The namespace to list resources from

**Returns:**
- `table|nil`: Array of Kubernetes objects, or nil on error
- `string|nil`: Error message if listing failed

**Example:**
```lua
local client = require("k8sclient")
local gvk = {group = "", version = "v1", kind = "Pod"}

local pods, err = client.list(gvk, "default")

if err then
    error("Failed to list pods: " .. err)
end

for i, pod in ipairs(pods) do
    print("Pod " .. i .. ": " .. pod.metadata.name)
end
```

## Common Resource Types

### Core Resources (group = "")

```lua
{group = "", version = "v1", kind = "Pod"}
{group = "", version = "v1", kind = "Service"}
{group = "", version = "v1", kind = "ConfigMap"}
{group = "", version = "v1", kind = "Secret"}
{group = "", version = "v1", kind = "Namespace"}
{group = "", version = "v1", kind = "PersistentVolumeClaim"}
```

### Apps Resources

```lua
{group = "apps", version = "v1", kind = "Deployment"}
{group = "apps", version = "v1", kind = "StatefulSet"}
{group = "apps", version = "v1", kind = "DaemonSet"}
{group = "apps", version = "v1", kind = "ReplicaSet"}
```

### Batch Resources

```lua
{group = "batch", version = "v1", kind = "Job"}
{group = "batch", version = "v1", kind = "CronJob"}
```

## Complete Example

```lua
local client = require("k8sclient")

-- Define resource type
local gvk = {group = "", version = "v1", kind = "ConfigMap"}

-- Create
local cm = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {
        name = "app-config",
        namespace = "default"
    },
    data = {
        database_url = "postgres://localhost:5432/mydb"
    }
}

local created, err = client.create(cm)
if err then error("Create failed: " .. err) end
print("Created: " .. created.metadata.name)

-- Get
local fetched, err = client.get(gvk, "default", "app-config")
if err then error("Get failed: " .. err) end
print("Current database_url: " .. fetched.data.database_url)

-- Update
fetched.data.database_url = "postgres://db.example.com:5432/mydb"
fetched.data.api_key = "secret123"

local updated, err = client.update(fetched)
if err then error("Update failed: " .. err) end
print("Updated with " .. #updated.data .. " keys")

-- List
local items, err = client.list(gvk, "default")
if err then error("List failed: " .. err) end
print("Total ConfigMaps: " .. #items)

-- Delete
local err = client.delete(gvk, "default", "app-config")
if err then error("Delete failed: " .. err) end
print("Deleted successfully")
```

## Error Handling

All functions follow Lua's error handling convention: they return `(result, error)`. Always check the error before using the result:

```lua
local result, err = client.get(gvk, namespace, name)
if err then
    -- Handle error
    print("Error: " .. err)
    return
end

-- Use result safely
print(result.metadata.name)
```

## Type Annotations for LSP

The module includes full type annotations for Lua language servers:

```lua
---@class k8sclient
local k8sclient = {}

---@param gvk GVKMatcher
---@param namespace string
---@param name string
---@return table|nil
---@return string|nil
function k8sclient.get(gvk, namespace, name) end

---@param object table
---@return table|nil
---@return string|nil
function k8sclient.create(object) end

---@param object table
---@return table|nil
---@return string|nil
function k8sclient.update(object) end

---@param gvk GVKMatcher
---@param namespace string
---@param name string
---@return string|nil
function k8sclient.delete(gvk, namespace, name) end

---@param gvk GVKMatcher
---@param namespace string
---@return table|nil
---@return string|nil
function k8sclient.list(gvk, namespace) end
```
