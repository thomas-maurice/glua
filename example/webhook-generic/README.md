# Generic Kubernetes Mutating Webhook with Lua

A production-ready Kubernetes mutating admission webhook that uses Lua scripts to mutate both Pods and Nodes.

## Features

- **Generic resource handling**: Supports multiple Kubernetes resource types (Pods, Nodes)
- **Lua-based mutations**: Easy-to-modify mutation logic without recompiling
- **Selective enabling**: Enable/disable mutations per resource type
- **Production-ready**: TLS support, health checks, structured logging
- **Type-safe**: Uses Kubernetes client-go types with proper deserialization

## What It Does

This webhook demonstrates two mutation examples:

### Pod Mutations

Adds a `even-mem: true` label to pods where any container requests an even number of memory bytes.

**Example:**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: app
    image: nginx
    resources:
      requests:
        memory: "128Mi"  # 134217728 bytes (even)
```

After mutation:

```yaml
metadata:
  labels:
    even-mem: "true"  # Added by webhook
```

### Node Mutations

Adds a `hello: ok` label to all nodes.

**Example:**

Before:

```yaml
apiVersion: v1
kind: Node
metadata:
  name: worker-1
```

After mutation:

```yaml
metadata:
  labels:
    hello: "ok"  # Added by webhook
```

## Architecture

```
┌─────────────────────────────────────────┐
│   Kubernetes API Server                 │
│                                          │
│   ┌──────────────────────────────────┐  │
│   │  Admission Controller            │  │
│   └──────────────┬───────────────────┘  │
└──────────────────┼──────────────────────┘
                   │ Admission Review
                   ▼
        ┌──────────────────────┐
        │  Webhook Server      │
        │                      │
        │  ┌────────────────┐  │
        │  │  Generic       │  │
        │  │  Mutator       │  │
        │  └────┬───────────┘  │
        │       │              │
        │  ┌────▼───────┐      │
        │  │ Lua        │      │
        │  │ Scripts    │      │
        │  └────────────┘      │
        └──────────────────────┘
```

## Quick Start

### Local Development

```bash
# Build
make build

# Run locally (without TLS)
make run

# Clean
make clean
```

### Docker

```bash
# Build image
make docker-build

# Push image
make docker-push TAG=v1.0.0
```

### Deploy to Kubernetes

```bash
# Install with Helm
helm install glua-webhook-generic ./charts/glua-webhook-generic

# Check deployment
kubectl get pods -l app=glua-webhook-generic
kubectl logs -l app=glua-webhook-generic
```

## Configuration

### Command-Line Flags

| Flag             | Default                      | Description                          |
|------------------|------------------------------|--------------------------------------|
| `-address`       | `:8443`                      | Address to listen on                 |
| `-cert`          | `/etc/webhook/certs/tls.crt` | Path to TLS certificate              |
| `-key`           | `/etc/webhook/certs/tls.key` | Path to TLS private key              |
| `-scripts`       | `/etc/webhook/scripts`       | Path to Lua scripts directory        |
| `-enable-nodes`  | `true`                       | Enable node mutations                |
| `-enable-pods`   | `true`                       | Enable pod mutations                 |

### Environment Variables (Helm)

Set via `values.yaml`:

```yaml
config:
  enableNodes: true
  enablePods: true
  logLevel: info
```

## Lua Scripts

### Script Structure

Each resource type has its own Lua script:

- `mutate_pod.lua` - Handles Pod mutations
- `mutate_node.lua` - Handles Node mutations

### Global Variables

Scripts have access to:

- `pod` or `node` - The Kubernetes resource object
- `patches` - Empty table to populate with JSON patch operations
- `kubernetes` module - Helper functions for Kubernetes operations

### Available Kubernetes Functions

```lua
local k8s = require("kubernetes")

-- Initialize metadata
pod = k8s.init_defaults(pod)

-- Label operations
pod = k8s.add_label(pod, "key", "value")
pod = k8s.add_labels(pod, {key1 = "value1", key2 = "value2"})
pod = k8s.remove_label(pod, "key")
local exists = k8s.has_label(pod, "key")
local value = k8s.get_label(pod, "key")

-- Annotation operations
pod = k8s.add_annotation(pod, "key", "value")
pod = k8s.remove_annotation(pod, "key")

-- Resource parsing
local bytes, err = k8s.parse_memory("128Mi")
local millicores, err = k8s.parse_cpu("500m")
local timestamp, err = k8s.parse_time("2024-01-15T10:30:00Z")
```

### Example: Custom Pod Mutation

```lua
---@diagnostic disable: undefined-global, lowercase-global
local k8s = require("kubernetes")

pod = k8s.init_defaults(pod)

-- Add annotation with pod's total memory request
local totalMemory = 0
for _, container in ipairs(pod.spec.containers or {}) do
    if container.resources and container.resources.requests then
        local mem = container.resources.requests.memory
        if mem then
            local bytes, err = k8s.parse_memory(mem)
            if not err then
                totalMemory = totalMemory + bytes
            end
        end
    end
end

-- Add annotation
pod = k8s.add_annotation(pod, "total-memory-bytes", tostring(totalMemory))

table.insert(patches, {
    op = "add",
    path = "/metadata/annotations/total-memory-bytes",
    value = tostring(totalMemory)
})
```

## JSON Patch Format

Patches are JSON Patch operations (RFC 6902):

```lua
-- Add a label
table.insert(patches, {
    op = "add",
    path = "/metadata/labels/my-label",
    value = "my-value"
})

-- Replace a value
table.insert(patches, {
    op = "replace",
    path = "/spec/replicas",
    value = 3
})

-- Remove a field
table.insert(patches, {
    op = "remove",
    path = "/metadata/labels/old-label"
})
```

## Testing

### Test Pod Mutation

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: test-even-mem
spec:
  containers:
  - name: nginx
    image: nginx:alpine
    resources:
      requests:
        memory: "128Mi"  # Even number of bytes
EOF

# Check if label was added
kubectl get pod test-even-mem -o jsonpath='{.metadata.labels.even-mem}'
# Output: true
```

### Test Node Mutation

Nodes are typically mutated during creation/updates. To test:

```bash
# Label a node (triggers webhook)
kubectl label node <node-name> test=value

# Check if hello=ok label was added
kubectl get node <node-name> -o jsonpath='{.metadata.labels.hello}'
# Output: ok
```

### View Webhook Logs

```bash
kubectl logs -l app=glua-webhook-generic -f
```

## Troubleshooting

### Webhook Not Called

1. Check webhook configuration:

   ```bash
   kubectl get mutatingwebhookconfigurations glua-webhook-generic -o yaml
   ```

2. Verify certificate and service:

   ```bash
   kubectl get certificate -n glua-webhook-generic
   kubectl get service -n glua-webhook-generic
   ```

### Lua Script Errors

Check webhook logs for script execution errors:

```bash
kubectl logs -l app=glua-webhook-generic | grep "lua mutation failed"
```

### Admission Denied

If admissions are denied, check:

1. TLS certificate validity
2. Service endpoint reachability
3. Webhook timeout settings

## Production Considerations

### Security

- Always use TLS in production
- Run as non-root user (handled by distroless image)
- Use RBAC to limit webhook permissions
- Validate Lua scripts before deployment

### Performance

- Lua scripts execute synchronously - keep them fast
- Consider caching parsed resources if expensive
- Monitor webhook latency metrics

### High Availability

- Run multiple webhook replicas (configured in Helm)
- Set appropriate resource requests/limits
- Configure pod disruption budgets

## Helm Chart

See `charts/glua-webhook-generic/` for deployment configuration.

Key values:

```yaml
replicaCount: 2

config:
  enableNodes: true
  enablePods: true

resources:
  limits:
    cpu: 200m
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

## License

MIT License - See root repository LICENSE file
