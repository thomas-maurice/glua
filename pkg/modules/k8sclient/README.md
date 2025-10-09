## Kubernetes Client Module

The `k8sclient` module provides a dynamic Kubernetes client for Lua scripts, allowing you to interact with Kubernetes clusters using unstructured objects.

### Features

- **Dynamic Client**: Works with any Kubernetes resource type
- **CRUD Operations**: Get, Create, Update, Delete, and List resources
- **Type-Safe**: Uses `GVKMatcher` type for LSP support
- **Unstructured**: Returns Lua tables directly from Kubernetes objects

### Usage

#### Setup

The module must be initialized with a `*rest.Config` from the Go side:

```go
import (
    "github.com/thomas-maurice/glua/pkg/modules/k8sclient"
    lua "github.com/yuin/gopher-lua"
    "k8s.io/client-go/tools/clientcmd"
)

// Load kubeconfig
config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
if err != nil {
    log.Fatal(err)
}

// Create Lua state and register module
L := lua.NewState()
defer L.Close()

L.PreloadModule("k8sclient", k8sclient.Loader(config))
```

#### Functions

##### `get(gvk, namespace, name)`

Retrieves a Kubernetes resource.

```lua
local client = require("k8sclient")

local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local cm, err = client.get(gvk, "default", "my-config")

if err then
    error("Failed to get ConfigMap: " .. err)
end

print("ConfigMap:", cm.metadata.name)
```

##### `create(object)`

Creates a new Kubernetes resource.

```lua
local configmap = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {
        name = "test-config",
        namespace = "default"
    },
    data = {
        key1 = "value1",
        key2 = "value2"
    }
}

local created, err = client.create(configmap)
if err then
    error("Failed to create: " .. err)
end

print("Created:", created.metadata.uid)
```

##### `update(object)`

Updates an existing Kubernetes resource.

```lua
-- Modify the object
cm.data.newkey = "newvalue"

local updated, err = client.update(cm)
if err then
    error("Failed to update: " .. err)
end

print("Updated:", updated.metadata.resourceVersion)
```

##### `delete(gvk, namespace, name)`

Deletes a Kubernetes resource.

```lua
local gvk = {group = "", version = "v1", kind = "ConfigMap"}
local err = client.delete(gvk, "default", "my-config")

if err then
    error("Failed to delete: " .. err)
end

print("Deleted successfully")
```

##### `list(gvk, namespace)`

Lists Kubernetes resources in a namespace.

```lua
local gvk = {group = "", version = "v1", kind = "Pod"}
local items, err = client.list(gvk, "default")

if err then
    error("Failed to list: " .. err)
end

print("Found", #items, "pods")
for i, pod in ipairs(items) do
    print("-", pod.metadata.name)
end
```

### GVKMatcher

The `GVKMatcher` type is used to specify Group/Version/Kind:

```lua
-- Core resources (empty group)
local pod_gvk = {group = "", version = "v1", kind = "Pod"}
local svc_gvk = {group = "", version = "v1", kind = "Service"}

-- Other API groups
local deploy_gvk = {group = "apps", version = "v1", kind = "Deployment"}
local ing_gvk = {group = "networking.k8s.io", version = "v1", kind = "Ingress"}
```

### Complete Example

```lua
local client = require("k8sclient")

-- ConfigMap GVK
local gvk = {group = "", version = "v1", kind = "ConfigMap"}

-- Create
local cm = {
    apiVersion = "v1",
    kind = "ConfigMap",
    metadata = {name = "app-config", namespace = "default"},
    data = {env = "prod"}
}

local created, err = client.create(cm)
if err then error(err) end

-- Get
local fetched, err = client.get(gvk, "default", "app-config")
if err then error(err) end

-- Update
fetched.data.version = "1.0"
local updated, err = client.update(fetched)
if err then error(err) end

-- List
local items, err = client.list(gvk, "default")
if err then error(err) end
print("Total ConfigMaps:", #items)

-- Delete
local err = client.delete(gvk, "default", "app-config")
if err then error(err) end
```

### Testing

The module includes integration tests that run against a Kind cluster:

```bash
cd example/kubernetes
./run-test.sh
```

This will:

1. Create a Kind cluster
2. Run the Lua integration test
3. Clean up the cluster

### Error Handling

All functions return errors as the last return value (Lua idiom):

```lua
local result, err = client.get(gvk, ns, name)
if err then
    -- Handle error
    error("Operation failed: " .. err)
end

-- Use result
print(result.metadata.name)
```

### LSP Support

The module is fully annotated for Lua LSP autocomplete. After running `stubgen`, you'll get proper type hints for all functions and the `GVKMatcher` type.
