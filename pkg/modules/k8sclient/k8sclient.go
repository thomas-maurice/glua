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

package k8sclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/modules/kubernetes"
	lua "github.com/yuin/gopher-lua"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var translator = glua.NewTranslator()

// Client holds the Kubernetes dynamic client
type Client struct {
	dynamic dynamic.Interface
}

// NewClient: creates a new Kubernetes client from a rest.Config
func NewClient(config *rest.Config) (*Client, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Client{
		dynamic: dynamicClient,
	}, nil
}

// Loader: creates and returns the k8sclient module for Lua.
// This function should be called with a rest.Config and then registered:
//
//	loader := k8sclient.Loader(config)
//	L.PreloadModule("k8sclient", loader)
//
// @luamodule k8sclient
//
// Example usage in Lua:
//
//	local k8sclient = require("k8sclient")
//	local client = k8sclient.new_client()
//	local gvk = {group = "", version = "v1", kind = "ConfigMap"}
//	local cm, err = client:get(gvk, "default", "my-config")
func Loader(config *rest.Config) lua.LGFunction {
	return func(L *lua.LState) int {
		// Create module table with new_client factory function
		mod := L.NewTable()
		L.SetField(mod, "new_client", L.NewFunction(func(L *lua.LState) int {
			return newClientLua(L, config)
		}))

		// For backwards compatibility, also export functions at module level
		client, err := NewClient(config)
		if err != nil {
			L.RaiseError("failed to create k8s client: %v", err)
			return 0
		}

		L.SetField(mod, "get", L.NewFunction(client.get))
		L.SetField(mod, "create", L.NewFunction(client.create))
		L.SetField(mod, "update", L.NewFunction(client.update))
		L.SetField(mod, "delete", L.NewFunction(client.delete))
		L.SetField(mod, "list", L.NewFunction(client.list))

		L.Push(mod)
		return 1
	}
}

// newClientLua: creates a new client instance in Lua
//
// @luafunc new_client
// @luareturn table The client instance with methods: get, create, update, delete, list
// @luareturn string|nil Error message if client creation failed
//
// Example:
//
//	local k8sclient = require("k8sclient")
//	local client = k8sclient.new_client()
//	local pod, err = client:get({group="", version="v1", kind="Pod"}, "default", "my-pod")
func newClientLua(L *lua.LState, config *rest.Config) int {
	client, err := NewClient(config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create client: %v", err)))
		return 2
	}

	// Create client table with methods
	clientTable := L.NewTable()
	L.SetField(clientTable, "get", L.NewFunction(client.get))
	L.SetField(clientTable, "create", L.NewFunction(client.create))
	L.SetField(clientTable, "update", L.NewFunction(client.update))
	L.SetField(clientTable, "delete", L.NewFunction(client.delete))
	L.SetField(clientTable, "list", L.NewFunction(client.list))

	L.Push(clientTable)
	L.Push(lua.LNil)
	return 2
}

// get: retrieves a Kubernetes resource by GVK, namespace, and name.
//
// @luafunc get
// @luaparam gvk GVKMatcher The GVK matcher with group, version, and kind
// @luaparam namespace string The namespace of the resource
// @luaparam name string The name of the resource
// @luareturn table|nil The Kubernetes object, or nil on error
// @luareturn string|nil Error message if retrieval failed
//
// Example:
//
//	local gvk = {group = "", version = "v1", kind = "ConfigMap"}
//	local cm, err = client.get(gvk, "default", "my-config")
func (c *Client) get(L *lua.LState) int {
	gvkTable := L.CheckTable(1)
	namespace := L.CheckString(2)
	name := L.CheckString(3)

	// Parse GVK
	var gvk kubernetes.GVKMatcher
	if err := translator.FromLua(L, gvkTable, &gvk); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse GVK: %v", err)))
		return 2
	}

	// Validate GVK
	if gvk.Kind == "" || gvk.Version == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("GVK requires 'kind' and 'version' fields"))
		return 2
	}

	// Build GVR
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}

	// Get resource
	obj, err := c.dynamic.Resource(gvr).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to get resource: %v", err)))
		return 2
	}

	// Convert to Lua
	luaObj, err := translator.ToLua(L, obj.Object)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to convert to Lua: %v", err)))
		return 2
	}

	L.Push(luaObj)
	L.Push(lua.LNil)
	return 2
}

// create: creates a Kubernetes resource from a Lua table.
//
// @luafunc create
// @luaparam obj table The Kubernetes object to create
// @luareturn table|nil The created Kubernetes object, or nil on error
// @luareturn string|nil Error message if creation failed
//
// Example:
//
//	local cm = {
//	  apiVersion = "v1",
//	  kind = "ConfigMap",
//	  metadata = {name = "my-config", namespace = "default"},
//	  data = {key = "value"}
//	}
//	local created, err = client.create(cm)
func (c *Client) create(L *lua.LState) int {
	objTable := L.CheckTable(1)

	// Convert to Go map
	var objMap map[string]interface{}
	if err := translator.FromLua(L, objTable, &objMap); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse object: %v", err)))
		return 2
	}

	// Create unstructured object
	obj := &unstructured.Unstructured{Object: objMap}

	// Extract GVK
	gvk := obj.GroupVersionKind()
	if gvk.Kind == "" || gvk.Version == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("object missing apiVersion or kind"))
		return 2
	}

	// Build GVR
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}

	// Get namespace
	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// Create resource
	created, err := c.dynamic.Resource(gvr).Namespace(namespace).Create(context.Background(), obj, metav1.CreateOptions{})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create resource: %v", err)))
		return 2
	}

	// Convert to Lua
	luaObj, err := translator.ToLua(L, created.Object)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to convert to Lua: %v", err)))
		return 2
	}

	L.Push(luaObj)
	L.Push(lua.LNil)
	return 2
}

// update: updates a Kubernetes resource.
//
// @luafunc update
// @luaparam obj table The Kubernetes object to update
// @luareturn table|nil The updated Kubernetes object, or nil on error
// @luareturn string|nil Error message if update failed
//
// Example:
//
//	cm.data.newkey = "newvalue"
//	local updated, err = client.update(cm)
func (c *Client) update(L *lua.LState) int {
	objTable := L.CheckTable(1)

	// Convert to Go map
	var objMap map[string]interface{}
	if err := translator.FromLua(L, objTable, &objMap); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse object: %v", err)))
		return 2
	}

	// Create unstructured object
	obj := &unstructured.Unstructured{Object: objMap}

	// Extract GVK
	gvk := obj.GroupVersionKind()
	if gvk.Kind == "" || gvk.Version == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("object missing apiVersion or kind"))
		return 2
	}

	// Build GVR
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}

	// Get namespace
	namespace := obj.GetNamespace()
	if namespace == "" {
		namespace = "default"
	}

	// Update resource
	updated, err := c.dynamic.Resource(gvr).Namespace(namespace).Update(context.Background(), obj, metav1.UpdateOptions{})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to update resource: %v", err)))
		return 2
	}

	// Convert to Lua
	luaObj, err := translator.ToLua(L, updated.Object)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to convert to Lua: %v", err)))
		return 2
	}

	L.Push(luaObj)
	L.Push(lua.LNil)
	return 2
}

// delete: deletes a Kubernetes resource by GVK, namespace, and name.
//
// @luafunc delete
// @luaparam gvk GVKMatcher The GVK matcher with group, version, and kind
// @luaparam namespace string The namespace of the resource
// @luaparam name string The name of the resource
// @luareturn string|nil Error message if deletion failed, nil on success
//
// Example:
//
//	local gvk = {group = "", version = "v1", kind = "ConfigMap"}
//	local err = client.delete(gvk, "default", "my-config")
func (c *Client) delete(L *lua.LState) int {
	gvkTable := L.CheckTable(1)
	namespace := L.CheckString(2)
	name := L.CheckString(3)

	// Parse GVK
	var gvk kubernetes.GVKMatcher
	if err := translator.FromLua(L, gvkTable, &gvk); err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to parse GVK: %v", err)))
		return 1
	}

	// Validate GVK
	if gvk.Kind == "" || gvk.Version == "" {
		L.Push(lua.LString("GVK requires 'kind' and 'version' fields"))
		return 1
	}

	// Build GVR
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}

	// Delete resource
	err := c.dynamic.Resource(gvr).Namespace(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to delete resource: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// list: lists Kubernetes resources by GVK and namespace.
//
// @luafunc list
// @luaparam gvk GVKMatcher The GVK matcher with group, version, and kind
// @luaparam namespace string The namespace to list from
// @luareturn table[]|nil Array of Kubernetes objects, or nil on error
// @luareturn string|nil Error message if listing failed
//
// Example:
//
//	local gvk = {group = "", version = "v1", kind = "ConfigMap"}
//	local items, err = client.list(gvk, "default")
func (c *Client) list(L *lua.LState) int {
	gvkTable := L.CheckTable(1)
	namespace := L.CheckString(2)

	// Parse GVK
	var gvk kubernetes.GVKMatcher
	if err := translator.FromLua(L, gvkTable, &gvk); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse GVK: %v", err)))
		return 2
	}

	// Validate GVK
	if gvk.Kind == "" || gvk.Version == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("GVK requires 'kind' and 'version' fields"))
		return 2
	}

	// Build GVR
	gvr := schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: pluralize(gvk.Kind),
	}

	// List resources
	list, err := c.dynamic.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list resources: %v", err)))
		return 2
	}

	// Convert items to Lua array
	items := L.CreateTable(len(list.Items), 0)
	for i, item := range list.Items {
		luaObj, err := translator.ToLua(L, item.Object)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to convert item to Lua: %v", err)))
			return 2
		}
		items.RawSetInt(i+1, luaObj)
	}

	L.Push(items)
	L.Push(lua.LNil)
	return 2
}

// pluralize: simple pluralization for resource names.
// This is a basic implementation - Kubernetes has more complex rules.
// Kubernetes resources are all lowercase.
func pluralize(kind string) string {
	// Convert to lowercase
	lower := strings.ToLower(kind)

	// Special cases
	switch lower {
	case "endpoints":
		return "endpoints"
	case "ingress":
		return "ingresses"
	}

	// Simple rule: add 's'
	if strings.HasSuffix(lower, "s") {
		return lower + "es"
	}
	return lower + "s"
}
