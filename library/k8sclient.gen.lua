---@meta k8sclient

---@class kubernetes.GVKMatcher
---@field group string
---@field kind string
---@field version string

---@class k8sclient.Client
local Client = {}

--- get a resource by GVK, namespace, and name
---@param gvk kubernetes.GVKMatcher
---@param namespace string
---@param name string
---@return table<string, any>
function Client:get(gvk, namespace, name) end

--- create a resource from a Lua table
---@param obj table<string, any>
---@return table<string, any>
function Client:create(obj) end

--- update a resource from a Lua table
---@param obj table<string, any>
---@return table<string, any>
function Client:update(obj) end

--- delete a resource by GVK, namespace, and name
---@param gvk kubernetes.GVKMatcher
---@param namespace string
---@param name string
function Client:delete(gvk, namespace, name) end

--- list resources by GVK and namespace
---@param gvk kubernetes.GVKMatcher
---@param namespace string
---@return table<string, any>[]
function Client:list(gvk, namespace) end

---@class k8sclient
---@field Client k8sclient.Client
---@field POD kubernetes.GVKMatcher Pod GVK constant
---@field NAMESPACE kubernetes.GVKMatcher Namespace GVK constant
---@field NODE kubernetes.GVKMatcher Node GVK constant
---@field CONFIGMAP kubernetes.GVKMatcher ConfigMap GVK constant
---@field SECRET kubernetes.GVKMatcher Secret GVK constant
---@field SERVICE kubernetes.GVKMatcher Service GVK constant
---@field SERVICEACCOUNT kubernetes.GVKMatcher ServiceAccount GVK constant
---@field PERSISTENTVOLUME kubernetes.GVKMatcher PersistentVolume GVK constant
---@field PERSISTENTVOLUMECLAIM kubernetes.GVKMatcher PersistentVolumeClaim GVK constant
---@field DEPLOYMENT kubernetes.GVKMatcher Deployment GVK constant
---@field STATEFULSET kubernetes.GVKMatcher StatefulSet GVK constant
---@field DAEMONSET kubernetes.GVKMatcher DaemonSet GVK constant
---@field REPLICASET kubernetes.GVKMatcher ReplicaSet GVK constant
---@field JOB kubernetes.GVKMatcher Job GVK constant
---@field CRONJOB kubernetes.GVKMatcher CronJob GVK constant
---@field INGRESS kubernetes.GVKMatcher Ingress GVK constant
---@field NETWORKPOLICY kubernetes.GVKMatcher NetworkPolicy GVK constant
---@field ROLE kubernetes.GVKMatcher Role GVK constant
---@field CLUSTERROLE kubernetes.GVKMatcher ClusterRole GVK constant
---@field ROLEBINDING kubernetes.GVKMatcher RoleBinding GVK constant
---@field CLUSTERROLEBINDING kubernetes.GVKMatcher ClusterRoleBinding GVK constant
local k8sclient = {}

--- create a new Kubernetes client
---@return number
function k8sclient.new_client() end

k8sclient.Client = Client

return k8sclient
