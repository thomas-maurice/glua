-- Copyright (c) 2024-2025 Thomas Maurice
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
-- SOFTWARE.

-- k8sclient example: Complete CRUD operations demonstration
--
-- This script demonstrates the full lifecycle of Kubernetes resource management from Lua:
-- - Dynamic client operations (create, get, update, delete, list)
-- - Working with multiple resource types (Pods, ConfigMaps, Services)
-- - Managing resource relationships (ConfigMap → Pod, Service → Pod)
-- - Label selectors and resource queries
-- - Practical nginx deployment with configuration
--
-- Use cases:
-- - Admission controllers that need to query cluster state
-- - Policy engines validating resource relationships
-- - Custom operators written in Lua
-- - GitOps reconciliation scripts
-- - Testing and development automation

local k8sclient = require("k8sclient")
local client = k8sclient.new_client()

print("╔═══════════════════════════════════════════════════════════════╗")
print("║   K8s Client Example: Full Resource Lifecycle Management     ║")
print("╚═══════════════════════════════════════════════════════════════╝")
print("")
print("This example demonstrates:")
print("  • Creating ConfigMaps for application configuration")
print("  • Deploying Pods with volume mounts and labels")
print("  • Creating Services to expose applications")
print("  • Querying resources with list operations")
print("  • Updating resources with annotations")
print("  • Deleting resources with cascading cleanup")
print("")

-- Define GVKs (Group/Version/Kind) for different resource types
local pod_gvk = {group = "", version = "v1", kind = "Pod"}
local configmap_gvk = {group = "", version = "v1", kind = "ConfigMap"}
local service_gvk = {group = "", version = "v1", kind = "Service"}

-- Step 1: Create a ConfigMap with nginx config
print("1. Creating ConfigMap with nginx config...")

---@type corev1.ConfigMap
local configmap = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "nginx-config",
		namespace = "default"
	},
	data = {
		["nginx.conf"] = [[
events {}
http {
    server {
        listen 8080;
        location / {
            return 200 'Hello from Lua-managed nginx!\n';
            add_header Content-Type text/plain;
        }
    }
}
]]
	}
}

local created_cm, err = client.create(configmap)
if err then
	error("Failed to create ConfigMap: " .. err)
end

print("   ✓ ConfigMap created: " .. created_cm.metadata.name)
print("   ✓ UID: " .. created_cm.metadata.uid)

-- Step 2: Create an nginx Pod
print("\n2. Creating nginx Pod...")

---@type corev1.Pod Pod definition with nginx configuration
local pod = {
	apiVersion = "v1",
	kind = "Pod",
	metadata = {
		name = "nginx-example",
		namespace = "default",
		labels = {
			app = "nginx",
			managed_by = "lua"
		}
	},
	---@type PodSpec
	spec = {
		containers = {
			{
				name = "nginx",
				image = "nginx:alpine",
				ports = {
					{
						containerPort = 8080,
						name = "http"
					}
				},
				volumeMounts = {
					{
						name = "config",
						mountPath = "/etc/nginx/nginx.conf",
						subPath = "nginx.conf"
					}
				}
			}
		},
		volumes = {
			{
				name = "config",
				configMap = {
					name = "nginx-config"
				}
			}
		}
	}
}

local created_pod, err = client.create(pod)
if err then
	error("Failed to create Pod: " .. err)
end

print("   ✓ Pod created: " .. created_pod.metadata.name)
print("   ✓ UID: " .. created_pod.metadata.uid)

-- Step 2b: Create a Service to expose the Pod
print("\n2b. Creating Service to expose nginx...")

---@type corev1.Service
local service = {
	apiVersion = "v1",
	kind = "Service",
	metadata = {
		name = "nginx-service",
		namespace = "default",
		labels = {
			app = "nginx"
		}
	},
	spec = {
		selector = {
			app = "nginx"
		},
		ports = {
			{
				protocol = "TCP",
				port = 80,
				targetPort = 8080,
				name = "http"
			}
		},
		type = "ClusterIP"
	}
}

local created_svc, err = client.create(service)
if err then
	error("Failed to create Service: " .. err)
end

print("   ✓ Service created: " .. created_svc.metadata.name)
print("   ✓ Type: " .. created_svc.spec.type)
print("   ✓ Cluster IP: " .. (created_svc.spec.clusterIP or "pending"))
print("   ✓ Selects pods with: app=" .. created_svc.spec.selector.app)

-- Step 3: Get the Pod
print("\n3. Retrieving Pod...")

local fetched_pod, err = client.get(pod_gvk, "default", "nginx-example")
if err then
	error("Failed to get Pod: " .. err)
end

print("   ✓ Pod retrieved: " .. fetched_pod.metadata.name)
print("   ✓ Image: " .. fetched_pod.spec.containers[1].image)
print("   ✓ Labels: app=" .. (fetched_pod.metadata.labels.app or "none"))

-- Step 4: Update the Pod (add annotation)
print("\n4. Updating Pod (adding annotation)...")

-- Wait a bit for pod to stabilize
os.execute("sleep 0.5")

-- Refetch the pod to get the latest version (avoid conflicts)
local latest_pod, err = client.get(pod_gvk, "default", "nginx-example")
if err then
	error("Failed to refetch Pod: " .. err)
end

if not latest_pod.metadata.annotations then
	latest_pod.metadata.annotations = {}
end
latest_pod.metadata.annotations.updated_by = "lua-script"
latest_pod.metadata.annotations.update_time = os.date("%Y-%m-%d %H:%M:%S")

local updated_pod, err = client.update(latest_pod)
if err then
	error("Failed to update Pod: " .. err)
end

print("   ✓ Pod updated with annotations:")
print("   ✓ updated_by: " .. (updated_pod.metadata.annotations.updated_by or "none"))
print("   ✓ update_time: " .. (updated_pod.metadata.annotations.update_time or "none"))

-- Step 5: List Pods
print("\n5. Listing Pods in default namespace...")

local pods, err = client.list(pod_gvk, "default")
if err then
	error("Failed to list Pods: " .. err)
end

print("   ✓ Found " .. #pods .. " Pod(s):")
for i, p in ipairs(pods) do
	local labels_str = ""
	if p.metadata.labels then
		for k, v in pairs(p.metadata.labels) do
			labels_str = labels_str .. k .. "=" .. v .. " "
		end
	end
	print("   " .. i .. ". " .. p.metadata.name .. " (labels: " .. labels_str .. ")")
end

-- Step 6: List ConfigMaps
print("\n6. Listing ConfigMaps in default namespace...")

local configmaps, err = client.list(configmap_gvk, "default")
if err then
	error("Failed to list ConfigMaps: " .. err)
end

print("   ✓ Found " .. #configmaps .. " ConfigMap(s):")
for i, cm in ipairs(configmaps) do
	local data_keys = ""
	if cm.data then
		for k, _ in pairs(cm.data) do
			data_keys = data_keys .. k .. " "
		end
	end
	print("   " .. i .. ". " .. cm.metadata.name .. " (data keys: " .. data_keys .. ")")
end

-- Step 7: Delete the Pod
print("\n7. Deleting Pod...")

local err = client.delete(pod_gvk, "default", "nginx-example")
if err then
	error("Failed to delete Pod: " .. err)
end

print("   ✓ Pod deleted successfully")

-- Step 8: Verify Pod deletion (note: pod may be in Terminating state)
print("\n8. Verifying Pod deletion...")

local deleted_pod, err = client.get(pod_gvk, "default", "nginx-example")
if err then
	print("   ✓ Pod confirmed deleted (expected error: " .. err .. ")")
elseif deleted_pod.metadata.deletionTimestamp then
	print("   ✓ Pod is terminating (deletionTimestamp: " .. deleted_pod.metadata.deletionTimestamp .. ")")
else
	print("   ⚠ Pod still exists but deletion was requested")
end

-- Step 9: Delete the Service
print("\n9. Deleting Service...")

local err = client.delete(service_gvk, "default", "nginx-service")
if err then
	error("Failed to delete Service: " .. err)
end

print("   ✓ Service deleted successfully")

-- Step 10: Delete the ConfigMap
print("\n10. Deleting ConfigMap...")

local err = client.delete(configmap_gvk, "default", "nginx-config")
if err then
	error("Failed to delete ConfigMap: " .. err)
end

print("   ✓ ConfigMap deleted successfully")

-- Step 11: Final verification
print("\n11. Final verification - listing all resources...")

local final_pods, err = client.list(pod_gvk, "default")
if err then
	error("Failed to list Pods: " .. err)
end

local final_cms, err = client.list(configmap_gvk, "default")
if err then
	error("Failed to list ConfigMaps: " .. err)
end

-- Count non-system resources
local user_pods = 0
local user_cms = 0

for _, pod in ipairs(final_pods) do
	-- Count terminating pods separately
	if pod.metadata.deletionTimestamp then
		print("   ℹ Pod '" .. pod.metadata.name .. "' is terminating")
	else
		user_pods = user_pods + 1
	end
end

for _, cm in ipairs(final_cms) do
	-- Exclude system ConfigMaps
	if cm.metadata.name ~= "kube-root-ca.crt" then
		user_cms = user_cms + 1
	end
end

print("   ✓ User Pods remaining: " .. user_pods)
print("   ✓ User ConfigMaps remaining: " .. user_cms)

if user_pods == 0 and user_cms == 0 then
	print("   ✓ All user resources cleaned up successfully!")
end

print("")
print("╔═══════════════════════════════════════════════════════════════╗")
print("║            All Operations Completed Successfully!            ║")
print("╚═══════════════════════════════════════════════════════════════╝")
print("")
print("Summary of operations performed:")
print("  ✓ Created ConfigMap with nginx configuration")
print("  ✓ Created Pod with volume mount from ConfigMap")
print("  ✓ Created Service to expose Pod")
print("  ✓ Retrieved resources with get() operations")
print("  ✓ Updated Pod with annotations")
print("  ✓ Listed resources with list() operations")
print("  ✓ Deleted all resources with delete() operations")
print("  ✓ Verified cleanup completion")
print("")
print("Key capabilities demonstrated:")
print("  • Dynamic Kubernetes client from Lua")
print("  • CRUD operations on multiple resource types")
print("  • Resource relationships (ConfigMap → Pod, Service → Pod)")
print("  • Label-based selectors")
print("  • Annotation updates")
print("  • Cascading deletion")
print("")
print("Use cases for k8sclient:")
print("  • Admission controllers querying cluster state")
print("  • Policy engines validating dependencies")
print("  • Custom operators in Lua")
print("  • GitOps reconciliation")
print("  • Testing and automation")
print("")
