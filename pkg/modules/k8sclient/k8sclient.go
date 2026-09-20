// Copyright (c) 2024-2025 Thomas Maurice
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package k8sclient provides a Kubernetes dynamic client for Lua scripts:
// get, create, update, delete and list arbitrary resources by GVK, plus GVK
// constants for common built-in resource types.
package k8sclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/thomas-maurice/glua/pkg/luareg"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	lua "github.com/yuin/gopher-lua"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// defaultCallTimeout: per-call timeout applied to every apiserver request.
const defaultCallTimeout = 30 * time.Second

// callContext: returns a context with the default per-call timeout.
func callContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultCallTimeout)
}

// Client: holds the Kubernetes dynamic client
type Client struct {
	dynamic dynamic.Interface
}

// NewClient: creates a new Kubernetes client from a rest.Config
func NewClient(config *rest.Config) (*Client, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}
	return &Client{dynamic: dynamicClient}, nil
}

// Get: retrieves a Kubernetes resource by GVK, namespace, and name. Raises on error.
//
// Example:
//
//	local obj = client:get(k8sclient.POD, "default", "my-pod")
func (c *Client) Get(gvk kubernetes.GVKMatcher, namespace, name string) (map[string]any, error) {
	if gvk.Kind == "" || gvk.Version == "" {
		return nil, fmt.Errorf("GVK requires 'kind' and 'version' fields")
	}
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}
	ctx, cancel := callContext()
	defer cancel()
	obj, err := c.dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}
	return obj.Object, nil
}

// Create: creates a Kubernetes resource from a map. Raises on error.
//
// Example:
//
//	local created = client:create(obj)
func (c *Client) Create(obj map[string]any) (map[string]any, error) {
	unstrObj := &unstructured.Unstructured{Object: obj}
	gvk := unstrObj.GroupVersionKind()
	if gvk.Kind == "" || gvk.Version == "" {
		return nil, fmt.Errorf("object missing apiVersion or kind")
	}
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}
	namespace := unstrObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	ctx, cancel := callContext()
	defer cancel()
	created, err := c.dynamic.Resource(gvr).Namespace(namespace).Create(ctx, unstrObj, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	return created.Object, nil
}

// Update: updates a Kubernetes resource. Raises on error.
//
// Example:
//
//	local updated = client:update(obj)
func (c *Client) Update(obj map[string]any) (map[string]any, error) {
	unstrObj := &unstructured.Unstructured{Object: obj}
	gvk := unstrObj.GroupVersionKind()
	if gvk.Kind == "" || gvk.Version == "" {
		return nil, fmt.Errorf("object missing apiVersion or kind")
	}
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}
	namespace := unstrObj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}
	ctx, cancel := callContext()
	defer cancel()
	updated, err := c.dynamic.Resource(gvr).Namespace(namespace).Update(ctx, unstrObj, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update resource: %w", err)
	}
	return updated.Object, nil
}

// Delete: deletes a Kubernetes resource by GVK, namespace, and name. Raises on error.
//
// Example:
//
//	client:delete(k8sclient.POD, "default", "my-pod")
func (c *Client) Delete(gvk kubernetes.GVKMatcher, namespace, name string) error {
	if gvk.Kind == "" || gvk.Version == "" {
		return fmt.Errorf("GVK requires 'kind' and 'version' fields")
	}
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}
	ctx, cancel := callContext()
	defer cancel()
	if err := c.dynamic.Resource(gvr).Namespace(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}
	return nil
}

// List: lists Kubernetes resources by GVK and namespace. Raises on error.
//
// Example:
//
//	local items = client:list(k8sclient.POD, "default")
func (c *Client) List(gvk kubernetes.GVKMatcher, namespace string) ([]map[string]any, error) {
	if gvk.Kind == "" || gvk.Version == "" {
		return nil, fmt.Errorf("GVK requires 'kind' and 'version' fields")
	}
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}
	ctx, cancel := callContext()
	defer cancel()
	lst, err := c.dynamic.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	items := make([]map[string]any, len(lst.Items))
	for i, item := range lst.Items {
		items[i] = item.Object
	}
	return items, nil
}

// newClientClass: creates a fresh luareg class for *Client.
// Each call returns a new instance to avoid duplicate-registration panics when
// build() is called more than once.
func newClientClass() *luareg.Class[*Client] {
	cls := luareg.NewClass[*Client]("k8sclient.Client", "Kubernetes dynamic client")
	cls.Method("get", (*Client).Get, "get a resource by GVK, namespace, and name",
		luareg.Args("gvk", "namespace", "name"),
		luareg.ArgDoc("gvk", "GVK matcher table with group, version and kind fields, e.g. k8sclient.POD"),
		luareg.ArgDoc("namespace", "namespace to look in; pass an empty string for cluster-scoped resources such as nodes or cluster roles"),
		luareg.ArgDoc("name", "name of the resource to fetch"),
		luareg.ReturnDoc(0, "obj", "the resource as a table, in the same shape as `kubectl get -o json`"))
	cls.Method("create", (*Client).Create, "create a resource from a Lua table",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the resource to create; must have apiVersion and kind set. If metadata.namespace is unset, it defaults to \"default\""),
		luareg.ReturnDoc(0, "created", "the created resource as returned by the apiserver, including server-set fields"))
	cls.Method("update", (*Client).Update, "update a resource from a Lua table",
		luareg.Args("obj"),
		luareg.ArgDoc("obj", "the resource to update; must have apiVersion and kind set. If metadata.namespace is unset, it defaults to \"default\""),
		luareg.ReturnDoc(0, "updated", "the updated resource as returned by the apiserver"))
	cls.Method("delete", (*Client).Delete, "delete a resource by GVK, namespace, and name",
		luareg.Args("gvk", "namespace", "name"),
		luareg.ArgDoc("gvk", "GVK matcher table with group, version and kind fields, e.g. k8sclient.POD"),
		luareg.ArgDoc("namespace", "namespace to delete from; pass an empty string for cluster-scoped resources such as nodes or cluster roles"),
		luareg.ArgDoc("name", "name of the resource to delete"))
	cls.Method("list", (*Client).List, "list resources by GVK and namespace",
		luareg.Args("gvk", "namespace"),
		luareg.ArgDoc("gvk", "GVK matcher table with group, version and kind fields, e.g. k8sclient.POD"),
		luareg.ArgDoc("namespace", "namespace to list in; pass an empty string for cluster-scoped resources, or to list across all namespaces"),
		luareg.ReturnDoc(0, "items", "table (array) of matching resources, each in the same shape as `kubectl get -o json`"))
	return cls
}

// newClientLua: module-level factory that creates a Client and wraps it in UserData.
// Raises a Lua error if the dynamic client cannot be created.
func newClientLua(L *lua.LState, config *rest.Config) int {
	client, err := NewClient(config)
	if err != nil {
		L.RaiseError("failed to create client: %v", err)
		return 0
	}

	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("k8sclient.Client"))
	L.Push(ud)
	return 1
}

// gvkConst: convenience wrapper to create a kubernetes.GVKMatcher constant.
func gvkConst(group, version, kind string) kubernetes.GVKMatcher {
	return kubernetes.GVKMatcher{Group: group, Version: version, Kind: kind}
}

// build: constructs the module definition. Reused by Loader and Register.
// GVK constants are registered via m.Const so they appear in generated stubs.
func build(config *rest.Config) *luareg.Module {
	m := luareg.NewModule("k8sclient", "Kubernetes dynamic client module")
	cls := newClientClass()
	m.RegisterClass(cls)
	// new_client factory: captures config in closure.
	m.Fn("new_client", func(L *lua.LState) int {
		return newClientLua(L, config)
	}, "create a new Kubernetes client")

	// GVK constants for common Kubernetes resources.
	m.Const("POD", gvkConst("", "v1", "Pod"), "kubernetes.GVKMatcher", "Pod GVK constant")
	m.Const("NAMESPACE", gvkConst("", "v1", "Namespace"), "kubernetes.GVKMatcher", "Namespace GVK constant")
	m.Const("NODE", gvkConst("", "v1", "Node"), "kubernetes.GVKMatcher", "Node GVK constant")
	m.Const("CONFIGMAP", gvkConst("", "v1", "ConfigMap"), "kubernetes.GVKMatcher", "ConfigMap GVK constant")
	m.Const("SECRET", gvkConst("", "v1", "Secret"), "kubernetes.GVKMatcher", "Secret GVK constant")
	m.Const("SERVICE", gvkConst("", "v1", "Service"), "kubernetes.GVKMatcher", "Service GVK constant")
	m.Const("SERVICEACCOUNT", gvkConst("", "v1", "ServiceAccount"), "kubernetes.GVKMatcher", "ServiceAccount GVK constant")
	m.Const("PERSISTENTVOLUME", gvkConst("", "v1", "PersistentVolume"), "kubernetes.GVKMatcher", "PersistentVolume GVK constant")
	m.Const("PERSISTENTVOLUMECLAIM", gvkConst("", "v1", "PersistentVolumeClaim"), "kubernetes.GVKMatcher", "PersistentVolumeClaim GVK constant")
	m.Const("DEPLOYMENT", gvkConst("apps", "v1", "Deployment"), "kubernetes.GVKMatcher", "Deployment GVK constant")
	m.Const("STATEFULSET", gvkConst("apps", "v1", "StatefulSet"), "kubernetes.GVKMatcher", "StatefulSet GVK constant")
	m.Const("DAEMONSET", gvkConst("apps", "v1", "DaemonSet"), "kubernetes.GVKMatcher", "DaemonSet GVK constant")
	m.Const("REPLICASET", gvkConst("apps", "v1", "ReplicaSet"), "kubernetes.GVKMatcher", "ReplicaSet GVK constant")
	m.Const("JOB", gvkConst("batch", "v1", "Job"), "kubernetes.GVKMatcher", "Job GVK constant")
	m.Const("CRONJOB", gvkConst("batch", "v1", "CronJob"), "kubernetes.GVKMatcher", "CronJob GVK constant")
	m.Const("INGRESS", gvkConst("networking.k8s.io", "v1", "Ingress"), "kubernetes.GVKMatcher", "Ingress GVK constant")
	m.Const("NETWORKPOLICY", gvkConst("networking.k8s.io", "v1", "NetworkPolicy"), "kubernetes.GVKMatcher", "NetworkPolicy GVK constant")
	m.Const("ROLE", gvkConst("rbac.authorization.k8s.io", "v1", "Role"), "kubernetes.GVKMatcher", "Role GVK constant")
	m.Const("CLUSTERROLE", gvkConst("rbac.authorization.k8s.io", "v1", "ClusterRole"), "kubernetes.GVKMatcher", "ClusterRole GVK constant")
	m.Const("ROLEBINDING", gvkConst("rbac.authorization.k8s.io", "v1", "RoleBinding"), "kubernetes.GVKMatcher", "RoleBinding GVK constant")
	m.Const("CLUSTERROLEBINDING", gvkConst("rbac.authorization.k8s.io", "v1", "ClusterRoleBinding"), "kubernetes.GVKMatcher", "ClusterRoleBinding GVK constant")

	return m
}

// Loader: creates and returns the k8sclient module for Lua.
//
// Example:
//
//	loader := k8sclient.Loader(config)
//	L.PreloadModule("k8sclient", loader)
func Loader(config *rest.Config) lua.LGFunction {
	return func(L *lua.LState) int {
		return build(config).PushTo(L)
	}
}

// Register: adds this module to reg for stub generation.
// Passes nil config since stub generation does not make API calls.
func Register(reg *luareg.Registry) {
	build(nil).Register(reg)
}

// pluralize: converts a Kubernetes resource kind to its plural form.
func pluralize(kind string) string {
	lower := strings.ToLower(kind)

	pluralMap := map[string]string{
		// Core resources (group="")
		"pod":                   "pods",
		"namespace":             "namespaces",
		"node":                  "nodes",
		"configmap":             "configmaps",
		"secret":                "secrets",
		"service":               "services",
		"serviceaccount":        "serviceaccounts",
		"persistentvolume":      "persistentvolumes",
		"persistentvolumeclaim": "persistentvolumeclaims",
		"endpoints":             "endpoints",
		// Apps resources (group="apps")
		"deployment":  "deployments",
		"statefulset": "statefulsets",
		"daemonset":   "daemonsets",
		"replicaset":  "replicasets",
		// Batch resources (group="batch")
		"job":     "jobs",
		"cronjob": "cronjobs",
		// Networking resources (group="networking.k8s.io")
		"ingress":       "ingresses",
		"networkpolicy": "networkpolicies",
		// RBAC resources (group="rbac.authorization.k8s.io")
		"role":               "roles",
		"clusterrole":        "clusterroles",
		"rolebinding":        "rolebindings",
		"clusterrolebinding": "clusterrolebindings",
	}

	if plural, ok := pluralMap[lower]; ok {
		return plural
	}

	if strings.HasSuffix(lower, "s") {
		return lower + "es"
	}
	if strings.HasSuffix(lower, "y") {
		return strings.TrimSuffix(lower, "y") + "ies"
	}
	return lower + "s"
}
