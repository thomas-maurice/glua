# glua-webhook: Lua-Powered Kubernetes Mutating Webhook

A Kubernetes mutating admission webhook that uses Lua scripts to dynamically modify pod specifications. This example demonstrates how to build a production-ready webhook using the glua library.

## Features

- **Lua-based mutation logic**: Write custom mutation rules in Lua instead of Go
- **Dynamic pod modification**: Add annotations, labels, or modify any pod specification
- **Production-ready**: Includes proper TLS certificate management via cert-manager
- **Helm chart**: Easy deployment to any Kubernetes cluster
- **Label-based targeting**: Only mutates pods with `thomas.maurice/mutate=true` label

## What This Example Does

This webhook demonstrates a **streamlined approach** to Kubernetes mutations using the `kubernetes` module's helper functions. It adds two annotations to any pod that has the label `thomas.maurice/mutate=true`:

1. `coucou.lil: hello` - A custom annotation demonstrating the mutation capability
2. `glua.mutated-at: <timestamp>` - Records when the mutation occurred

**Key improvement**: Uses `kubernetes` module helpers (`init_defaults()`, `add_annotation()`) instead of manual JSON patch manipulation, making the Lua code much simpler and more maintainable.

## Architecture

```
┌─────────────────┐
│  Kubernetes API │
│     Server      │
└────────┬────────┘
         │ Pod CREATE/UPDATE
         ▼
┌─────────────────────┐
│ MutatingWebhook     │
│ Configuration       │
└────────┬────────────┘
         │ HTTPS Request
         ▼
┌─────────────────────┐
│  glua-webhook       │
│  (Go Server)        │
│  ┌───────────────┐  │
│  │ Lua Runtime   │  │
│  │ mutate.lua    │  │
│  └───────────────┘  │
└────────┬────────────┘
         │ JSON Patches
         ▼
┌─────────────────┐
│  Modified Pod   │
└─────────────────┘
```

## Prerequisites

- Kubernetes cluster (local or remote)
- `kubectl` configured to access the cluster
- `helm` (v3+) installed
- `docker` for building the image
- `kind` for local testing (optional)

## Quick Start

### 1. Create a Kind Cluster (Local Testing)

```bash
make kind
```

This creates a local Kubernetes cluster named `glua-webhook-test`.

### 2. Install cert-manager

The webhook requires TLS certificates managed by cert-manager:

```bash
make cert-manager
```

This installs cert-manager and waits for it to be ready.

### 3. Build and Deploy the Webhook

```bash
make install
```

This will:

- Build the Docker image
- Load it into the Kind cluster
- Deploy the webhook using Helm
- Configure the mutating webhook

### 4. Test the Mutation

```bash
make test
```

This creates a test pod with the mutation label and displays the annotations that were added.

Expected output:

```json
{
  "coucou.lil": "hello",
  "glua.mutated-at": "2025-10-07T00:00:00Z"
}
```

## Manual Testing

Create a pod with the mutation label:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: my-test-pod
  labels:
    thomas.maurice/mutate: "true"
spec:
  containers:
  - name: nginx
    image: nginx:alpine
EOF
```

Verify the annotations were added:

```bash
kubectl get pod my-test-pod -o jsonpath='{.metadata.annotations}' | jq .
```

## Customizing the Mutation Logic

The mutation logic is defined in Lua and stored in a ConfigMap. To customize:

1. Edit [`scripts/mutate.lua`](scripts/mutate.lua) or the ConfigMap in [`charts/glua-webhook/templates/configmap.yaml`](charts/glua-webhook/templates/configmap.yaml)

2. Upgrade the Helm release:

```bash
helm upgrade glua-webhook ./charts/glua-webhook --namespace glua-webhook
```

### Lua Script API

The Lua script has access to:

- `pod`: The Kubernetes Pod object as a Lua table
- `patches`: An empty table to populate with JSON patch operations
- `kubernetes`: The kubernetes module with helper functions

**Streamlined approach using kubernetes module:**

```lua
local k8s = require("kubernetes")

-- Ensure metadata structures exist
pod = k8s.init_defaults(pod)

-- Add annotations easily - no manual path handling needed!
pod = k8s.add_annotation(pod, "my-key", "my-value")
pod = k8s.add_label(pod, "my-label", "my-value")

-- Generate patches for the annotations
for key, value in pairs(pod.metadata.annotations) do
    if key == "my-key" then  -- Only patch what we added
        local escaped_key = key:gsub("~", "~0"):gsub("/", "~1")
        table.insert(patches, {
            op = "add",
            path = "/metadata/annotations/" .. escaped_key,
            value = value
        })
    end
end
```

**Traditional approach (also supported):**

```lua
-- Manual JSON patch operation
table.insert(patches, {
  op = "add",
  path = "/metadata/annotations/my-key",
  value = "my-value"
})
```

Supported operations: `add`, `remove`, `replace`, `move`, `copy`, `test`

**See also:** The [kubernetes module documentation](../../README.md#kubernetes) for all available helper functions.

## Configuration

### Helm Values

Key configuration options in [`charts/glua-webhook/values.yaml`](charts/glua-webhook/values.yaml):

```yaml
# Image configuration
image:
  repository: glua-webhook
  tag: latest
  pullPolicy: IfNotPresent

# Webhook behavior
webhook:
  failurePolicy: Fail  # Fail or Ignore
  timeoutSeconds: 10
  objectSelector:
    matchLabels:
      thomas.maurice/mutate: "true"

# Resources
resources:
  limits:
    cpu: 200m
    memory: 128Mi
  requests:
    cpu: 100m
    memory: 64Mi
```

### Webhook Scope

The webhook is configured to only mutate pods with the label `thomas.maurice/mutate=true`. To change this:

1. Edit `webhook.objectSelector` in `values.yaml`
2. Upgrade the Helm release

To mutate all pods in specific namespaces:

```yaml
webhook:
  namespaceSelector:
    matchLabels:
      glua-mutate: "enabled"
  objectSelector: {}
```

## Makefile Targets

```bash
make help              # Show available targets
make build             # Build the webhook binary
make docker            # Build Docker image
make kind              # Create Kind cluster
make cert-manager      # Install cert-manager
make install           # Deploy webhook to cluster
make uninstall         # Remove webhook from cluster
make test              # Run mutation test
make logs              # Show webhook logs
make clean             # Clean build artifacts
make all               # Create cluster, install, and test
```

## Development

### Building Locally

```bash
make build
./bin/webhook --help
```

### Running Locally (Outside Kubernetes)

```bash
# Generate self-signed certificates
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes

# Run the webhook
./bin/webhook \
  --address=:8443 \
  --cert=cert.pem \
  --key=key.pem \
  --script=scripts/mutate.lua
```

### Viewing Logs

```bash
make logs
```

Or directly:

```bash
kubectl logs -n glua-webhook -l app=glua-webhook -f
```

## Troubleshooting

### Webhook Not Mutating Pods

1. Check webhook logs: `make logs`
2. Verify the pod has the correct label: `thomas.maurice/mutate=true`
3. Check webhook configuration: `kubectl get mutatingwebhookconfigurations glua-webhook -o yaml`

### Certificate Issues

If you see TLS errors:

```bash
# Check certificate status
kubectl get certificate -n glua-webhook

# Check cert-manager logs
kubectl logs -n cert-manager -l app=cert-manager
```

### Pod Creation Failures

If pods fail to create:

1. Check webhook failure policy (set to `Ignore` for debugging):

```bash
helm upgrade glua-webhook ./charts/glua-webhook \
  --namespace glua-webhook \
  --set webhook.failurePolicy=Ignore
```

2. Check webhook is reachable:

```bash
kubectl get svc -n glua-webhook
kubectl get endpoints -n glua-webhook
```

## Cleanup

```bash
# Uninstall webhook
make uninstall

# Delete Kind cluster
make kind-delete
```

## Project Structure

```
example/webhook/
├── main.go                           # Webhook server implementation
├── Dockerfile                        # Container image definition
├── Makefile                          # Build and deployment automation
├── README.md                         # This file
├── scripts/
│   └── mutate.lua                   # Lua mutation script
└── charts/glua-webhook/             # Helm chart
    ├── Chart.yaml                   # Chart metadata
    ├── values.yaml                  # Default configuration
    └── templates/
        ├── _helpers.tpl             # Template helpers
        ├── serviceaccount.yaml      # Service account
        ├── service.yaml             # Kubernetes service
        ├── deployment.yaml          # Webhook deployment
        ├── configmap.yaml           # Lua script ConfigMap
        ├── certificate.yaml         # TLS certificate
        └── mutatingwebhook.yaml     # Webhook configuration
```

## Security Considerations

- The webhook runs as a non-root user (UID 65532)
- Read-only root filesystem
- Dropped capabilities
- TLS encryption for all webhook traffic
- Certificate rotation via cert-manager

## Learn More

- [Kubernetes Admission Controllers](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
- [cert-manager Documentation](https://cert-manager.io/docs/)
- [glua Library](https://github.com/thomas-maurice/glua)
