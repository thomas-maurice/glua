---@meta k8sclient

---@class GVKMatcher
---@field group string The API group (empty string for core resources)
---@field version string The API version (e.g., "v1", "v1beta1")
---@field kind string The resource kind (e.g., "Pod", "Deployment")

---@class k8sclient
local k8sclient = {}

---@return table client The client instance with methods: get, create, update, delete, list
---@return string|nil err Error message if client creation failed
function k8sclient.new_client() end

---@param GVKMatcher gvk The GVK matcher with group, version, and kind
---@param string namespace The namespace of the resource
---@param string name The name of the resource
---@return table|nil obj The Kubernetes object, or nil on error
---@return string|nil err Error message if retrieval failed
function k8sclient.get(gvk, namespace, name) end

---@param table obj The Kubernetes object to create
---@return table|nil obj The created Kubernetes object, or nil on error
---@return string|nil err Error message if creation failed
function k8sclient.create(obj) end

---@param table obj The Kubernetes object to update
---@return table|nil obj The updated Kubernetes object, or nil on error
---@return string|nil err Error message if update failed
function k8sclient.update(obj) end

---@param GVKMatcher gvk The GVK matcher with group, version, and kind
---@param string namespace The namespace of the resource
---@param string name The name of the resource
---@return string|nil err Error message if deletion failed, nil on success
function k8sclient.delete(gvk, namespace, name) end

---@param GVKMatcher gvk The GVK matcher with group, version, and kind
---@param string namespace The namespace to list from
---@return table[]|nil objects Array of Kubernetes objects, or nil on error
---@return string|nil err Error message if listing failed
function k8sclient.list(gvk, namespace) end

---@type table Pod GVK constant {group="", version="v1", kind="Pod"}
k8sclient.POD = nil

---@type table Namespace GVK constant {group="", version="v1", kind="Namespace"}
k8sclient.NAMESPACE = nil

---@type table Node GVK constant {group="", version="v1", kind="Node"}
k8sclient.NODE = nil

---@type table ConfigMap GVK constant {group="", version="v1", kind="ConfigMap"}
k8sclient.CONFIGMAP = nil

---@type table Secret GVK constant {group="", version="v1", kind="Secret"}
k8sclient.SECRET = nil

---@type table Service GVK constant {group="", version="v1", kind="Service"}
k8sclient.SERVICE = nil

---@type table ServiceAccount GVK constant {group="", version="v1", kind="ServiceAccount"}
k8sclient.SERVICEACCOUNT = nil

---@type table PersistentVolume GVK constant {group="", version="v1", kind="PersistentVolume"}
k8sclient.PERSISTENTVOLUME = nil

---@type table PersistentVolumeClaim GVK constant {group="", version="v1", kind="PersistentVolumeClaim"}
k8sclient.PERSISTENTVOLUMECLAIM = nil

---@type table Deployment GVK constant {group="apps", version="v1", kind="Deployment"}
k8sclient.DEPLOYMENT = nil

---@type table StatefulSet GVK constant {group="apps", version="v1", kind="StatefulSet"}
k8sclient.STATEFULSET = nil

---@type table DaemonSet GVK constant {group="apps", version="v1", kind="DaemonSet"}
k8sclient.DAEMONSET = nil

---@type table ReplicaSet GVK constant {group="apps", version="v1", kind="ReplicaSet"}
k8sclient.REPLICASET = nil

---@type table Job GVK constant {group="batch", version="v1", kind="Job"}
k8sclient.JOB = nil

---@type table CronJob GVK constant {group="batch", version="v1", kind="CronJob"}
k8sclient.CRONJOB = nil

---@type table Ingress GVK constant {group="networking.k8s.io", version="v1", kind="Ingress"}
k8sclient.INGRESS = nil

---@type table NetworkPolicy GVK constant {group="networking.k8s.io", version="v1", kind="NetworkPolicy"}
k8sclient.NETWORKPOLICY = nil

---@type table Role GVK constant {group="rbac.authorization.k8s.io", version="v1", kind="Role"}
k8sclient.ROLE = nil

---@type table ClusterRole GVK constant {group="rbac.authorization.k8s.io", version="v1", kind="ClusterRole"}
k8sclient.CLUSTERROLE = nil

---@type table RoleBinding GVK constant {group="rbac.authorization.k8s.io", version="v1", kind="RoleBinding"}
k8sclient.ROLEBINDING = nil

---@type table ClusterRoleBinding GVK constant {group="rbac.authorization.k8s.io", version="v1", kind="ClusterRoleBinding"}
k8sclient.CLUSTERROLEBINDING = nil

return k8sclient
