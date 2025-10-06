# k8sclient Example

This example demonstrates how to use the `k8sclient` module to interact with Kubernetes clusters from Lua scripts.

## Overview

The example creates an nginx Pod with a custom configuration served via a ConfigMap, then performs full CRUD operations:

1. **Create** a ConfigMap with nginx configuration
2. **Create** an nginx Pod using the ConfigMap
3. **Read** (Get) the Pod to verify creation
4. **Update** the Pod with annotations
5. **List** Pods and ConfigMaps
6. **Delete** all resources
7. **Verify** cleanup

## Prerequisites

- Go 1.24+
- [Kind](https://kind.sigs.k8s.io/) (Kubernetes in Docker)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

### Installing Kind

```bash
# Linux/macOS
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# macOS (via Homebrew)
brew install kind
```

### Installing kubectl

```bash
# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/kubectl

# macOS (via Homebrew)
brew install kubectl
```

## Usage

### Run with Kind Cluster (Automated)

The easiest way to run the example is using the included shell script:

```bash
./run-test.sh
```

This will:
1. Create a temporary Kind cluster
2. Build the example program
3. Run the Lua test script
4. Clean up the Kind cluster

### Run with Existing Cluster

If you have an existing Kubernetes cluster configured in `~/.kube/config`:

```bash
go build -o k8sclient-example .
./k8sclient-example -script test.lua
```

Or specify a custom kubeconfig:

```bash
./k8sclient-example -kubeconfig /path/to/kubeconfig -script test.lua
```

### Run from Project Root

Using Make:

```bash
# From project root
make test-k8sclient
```

## Example Output

```
=== K8s Client Example: nginx Pod CRUD Operations ===

1. Creating ConfigMap with nginx config...
   ✓ ConfigMap created: nginx-config
   ✓ UID: b5ac3a4d-a0cd-4019-a4d1-cb4141f246b5

2. Creating nginx Pod...
   ✓ Pod created: nginx-example
   ✓ UID: 315b4807-8cba-42b3-bfe9-3da7776781a6

3. Retrieving Pod...
   ✓ Pod retrieved: nginx-example
   ✓ Image: nginx:alpine
   ✓ Labels: app=nginx

4. Updating Pod (adding annotation)...
   ✓ Pod updated with annotations:
   ✓ updated_by: lua-script
   ✓ update_time: 2025-10-07 01:42:43

5. Listing Pods in default namespace...
   ✓ Found 1 Pod(s):
   1. nginx-example (labels: app=nginx managed_by=lua )

6. Listing ConfigMaps in default namespace...
   ✓ Found 2 ConfigMap(s):
   1. kube-root-ca.crt (data keys: ca.crt )
   2. nginx-config (data keys: nginx.conf )

7. Deleting Pod...
   ✓ Pod deleted successfully

8. Verifying Pod deletion...
   ✓ Pod is terminating (deletionTimestamp: 2025-10-06T23:43:13Z)

9. Deleting ConfigMap...
   ✓ ConfigMap deleted successfully

10. Final verification - listing all resources...
   ℹ Pod 'nginx-example' is terminating
   ✓ User Pods remaining: 0
   ✓ User ConfigMaps remaining: 0
   ✓ All user resources cleaned up successfully!

=== All CRUD operations completed successfully! ===
```

## Type Annotations

The example includes comprehensive type annotations for LSP support. When using a Lua language server (like lua-language-server), you'll get:

- Autocompletion for Kubernetes resource structures
- Type checking for Pod specs, containers, volumes, etc.
- Inline documentation

Example type definitions used:

```lua
---@class PodSpec
---@field containers Container[] List of containers in the pod
---@field volumes Volume[] List of volumes that can be mounted by containers

---@class Container
---@field name string Name of the container
---@field image string Docker image name
---@field ports ContainerPort[] List of ports to expose from the container
---@field volumeMounts VolumeMount[] Pod volumes to mount into the container
```

## Files

- [`main.go`](main.go) - Go wrapper that loads kubeconfig and executes Lua scripts
- [`test.lua`](test.lua) - Comprehensive CRUD test demonstrating all k8sclient operations
- [`run-test.sh`](run-test.sh) - Automated test script with Kind cluster management
- [`.luarc.json`](.luarc.json) - Lua language server configuration

## Learning Resources

The test script demonstrates:

1. **GVK (Group/Version/Kind) matching** - How to specify Kubernetes resource types
2. **Creating resources** - Building ConfigMaps and Pods from Lua tables
3. **Reading resources** - Fetching individual resources by name
4. **Updating resources** - Modifying existing resources (handling resourceVersion conflicts)
5. **Listing resources** - Querying multiple resources with filtering
6. **Deleting resources** - Cleaning up with proper verification
7. **Error handling** - Lua-style error handling with (result, error) pattern

## CI/CD Integration

This example is automatically tested in GitHub Actions as part of the project's CI pipeline. See [`.github/workflows/test.yml`](../../.github/workflows/test.yml) for the workflow configuration.
