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
	"path/filepath"
	"testing"

	lua "github.com/yuin/gopher-lua"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/rest"
)

// TestPluralize tests the pluralization function
func TestPluralize(t *testing.T) {
	tests := []struct {
		kind     string
		expected string
	}{
		{"Pod", "pods"},
		{"ConfigMap", "configmaps"},
		{"Service", "services"},
		{"Deployment", "deployments"},
		{"Ingress", "ingresses"},
		{"Endpoints", "endpoints"},
		{"StatefulSet", "statefulsets"},
		{"NetworkPolicy", "networkpolicies"},
		{"CronJob", "cronjobs"},
		{"PersistentVolume", "persistentvolumes"},
		// Fallback cases
		{"CustomResource", "customresources"},
		{"Entity", "entities"},
		{"Glass", "glasses"},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			result := pluralize(tt.kind)
			if result != tt.expected {
				t.Errorf("pluralize(%q) = %q, want %q", tt.kind, result, tt.expected)
			}
		})
	}
}

// TestLoaderCreation tests that the Loader can be created with valid config
func TestLoaderCreation(t *testing.T) {
	config := &rest.Config{
		Host: "https://localhost:6443",
	}

	loader := Loader(config)
	if loader == nil {
		t.Errorf("Expected non-nil loader function")
	}
}

// setupFakeClient creates a fake dynamic client for testing
func setupFakeClient(objects ...runtime.Object) *Client {
	scheme := runtime.NewScheme()

	fakeClient := fake.NewSimpleDynamicClientWithCustomListKinds(
		scheme,
		map[schema.GroupVersionResource]string{
			{Group: "", Version: "v1", Resource: "configmaps"}: "ConfigMapList",
		},
		objects...,
	)

	return &Client{
		dynamic: fakeClient,
	}
}

// Test constants
const (
	testNamespace        = "default"
	testConfigName       = "test-config"
	testConfigKey        = "key1"
	testConfigValue      = "value1"
	testNewConfigName    = "new-config"
	testDataKey          = "foo"
	testDataValue        = "bar"
	testUpdateKey        = "updated"
	testUpdateValue      = "new-value"
	testDeleteConfigName = "delete-config"
	testDeleteDataKey    = "test"
	testDeleteDataValue  = "data"
	testUpdateConfigName = "update-config"
	testOriginalKey      = "original"
	testOriginalValue    = "value"
	testIntegrationName  = "integration-config"
	testInitialKey       = "initial"
	testInitialValue     = "value"
	testListConfig1      = "list-config-1"
	testListConfig2      = "list-config-2"
	testListKey1         = "key1"
	testListValue1       = "value1"
	testListKey2         = "key2"
	testListValue2       = "value2"
)

// setupLuaWithClient creates a Lua state with the k8sclient module loaded and test constants.
// The client is exposed as a UserData global named "client" with the class metatable.
func setupLuaWithClient(client *Client) *lua.LState {
	L := lua.NewState()

	// Register the class metatable so Wrap works.
	cls := newClientClass()
	tmp := build(nil)
	_ = tmp
	// We need to register the class metatable on L. Do this by pushing and popping a module.
	m := build(nil)
	m.PushTo(L)
	L.Pop(1)

	// Wrap the client as UserData with the class metatable.
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable(cls.Name()))
	L.SetGlobal("client", ud)

	// Set test constants as globals
	L.SetGlobal("TEST_NAMESPACE", lua.LString(testNamespace))
	L.SetGlobal("TEST_CONFIG_NAME", lua.LString(testConfigName))
	L.SetGlobal("TEST_CONFIG_KEY", lua.LString(testConfigKey))
	L.SetGlobal("TEST_CONFIG_VALUE", lua.LString(testConfigValue))
	L.SetGlobal("TEST_NEW_CONFIG_NAME", lua.LString(testNewConfigName))
	L.SetGlobal("TEST_DATA_KEY", lua.LString(testDataKey))
	L.SetGlobal("TEST_DATA_VALUE", lua.LString(testDataValue))
	L.SetGlobal("TEST_UPDATE_KEY", lua.LString(testUpdateKey))
	L.SetGlobal("TEST_UPDATE_VALUE", lua.LString(testUpdateValue))
	L.SetGlobal("TEST_DELETE_CONFIG_NAME", lua.LString(testDeleteConfigName))
	L.SetGlobal("TEST_DELETE_DATA_KEY", lua.LString(testDeleteDataKey))
	L.SetGlobal("TEST_DELETE_DATA_VALUE", lua.LString(testDeleteDataValue))
	L.SetGlobal("TEST_UPDATE_CONFIG_NAME", lua.LString(testUpdateConfigName))
	L.SetGlobal("TEST_ORIGINAL_KEY", lua.LString(testOriginalKey))
	L.SetGlobal("TEST_ORIGINAL_VALUE", lua.LString(testOriginalValue))
	L.SetGlobal("TEST_INTEGRATION_NAME", lua.LString(testIntegrationName))
	L.SetGlobal("TEST_INITIAL_KEY", lua.LString(testInitialKey))
	L.SetGlobal("TEST_INITIAL_VALUE", lua.LString(testInitialValue))
	L.SetGlobal("TEST_LIST_CONFIG_1", lua.LString(testListConfig1))
	L.SetGlobal("TEST_LIST_CONFIG_2", lua.LString(testListConfig2))
	L.SetGlobal("TEST_LIST_KEY_1", lua.LString(testListKey1))
	L.SetGlobal("TEST_LIST_VALUE_1", lua.LString(testListValue1))
	L.SetGlobal("TEST_LIST_KEY_2", lua.LString(testListKey2))
	L.SetGlobal("TEST_LIST_VALUE_2", lua.LString(testListValue2))

	return L
}

// TestLuaScripts: runs all Lua test scripts in testdata/ directory
func TestLuaScripts(t *testing.T) {
	files, err := filepath.Glob("testdata/test_*.lua")
	if err != nil {
		t.Fatalf("Failed to glob testdata: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("No Lua test files found in testdata/")
	}

	for _, file := range files {
		testName := filepath.Base(file)
		t.Run(testName, func(t *testing.T) {
			client := setupFakeClient()
			L := setupLuaWithClient(client)
			defer L.Close()

			if err := L.DoFile(file); err != nil {
				t.Fatalf("Test script %s failed: %v", testName, err)
			}

			result := L.Get(-1)
			if result != lua.LTrue {
				t.Errorf("Test script %s returned %v, expected true", testName, result)
			}
		})
	}
}

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	config := &rest.Config{
		Host: "https://localhost:6443",
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() failed: %v", err)
	}

	if client == nil {
		t.Fatal("NewClient() returned nil client")
	}

	if client.dynamic == nil {
		t.Error("NewClient() returned client with nil dynamic interface")
	}
}

// TestGVKConstants verifies all 20 GVK constants registered via m.Const in
// build() are actually reachable from Lua with the expected group/version/kind.
// This is the only coverage of the live registration path: createGVKTable and
// addGVKConstants (an earlier, unused approach to the same problem) were
// removed as dead code, so this test replaces the coverage they used to
// provide for a different, no-longer-used code path.
func TestGVKConstants(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	config := &rest.Config{Host: "https://localhost:6443"}
	L.PreloadModule("k8sclient", func(L *lua.LState) int {
		return build(config).PushTo(L)
	})

	constants := []struct {
		name    string
		group   string
		version string
		kind    string
	}{
		{"POD", "", "v1", "Pod"},
		{"NAMESPACE", "", "v1", "Namespace"},
		{"NODE", "", "v1", "Node"},
		{"CONFIGMAP", "", "v1", "ConfigMap"},
		{"SECRET", "", "v1", "Secret"},
		{"SERVICE", "", "v1", "Service"},
		{"SERVICEACCOUNT", "", "v1", "ServiceAccount"},
		{"PERSISTENTVOLUME", "", "v1", "PersistentVolume"},
		{"PERSISTENTVOLUMECLAIM", "", "v1", "PersistentVolumeClaim"},
		{"DEPLOYMENT", "apps", "v1", "Deployment"},
		{"STATEFULSET", "apps", "v1", "StatefulSet"},
		{"DAEMONSET", "apps", "v1", "DaemonSet"},
		{"REPLICASET", "apps", "v1", "ReplicaSet"},
		{"JOB", "batch", "v1", "Job"},
		{"CRONJOB", "batch", "v1", "CronJob"},
		{"INGRESS", "networking.k8s.io", "v1", "Ingress"},
		{"NETWORKPOLICY", "networking.k8s.io", "v1", "NetworkPolicy"},
		{"ROLE", "rbac.authorization.k8s.io", "v1", "Role"},
		{"CLUSTERROLE", "rbac.authorization.k8s.io", "v1", "ClusterRole"},
		{"ROLEBINDING", "rbac.authorization.k8s.io", "v1", "RoleBinding"},
		{"CLUSTERROLEBINDING", "rbac.authorization.k8s.io", "v1", "ClusterRoleBinding"},
	}

	for _, tc := range constants {
		t.Run(tc.name, func(t *testing.T) {
			code := `
				local k8sclient = require("k8sclient")
				local gvk = k8sclient.` + tc.name + `
				return gvk.group, gvk.version, gvk.kind
			`
			if err := L.DoString(code); err != nil {
				t.Fatalf("DoString failed: %v", err)
			}
			if L.GetTop() != 3 {
				t.Fatalf("Expected 3 return values, got %d", L.GetTop())
			}
			group := L.Get(-3).String()
			version := L.Get(-2).String()
			kind := L.Get(-1).String()
			L.Pop(3)

			if group != tc.group {
				t.Errorf("Expected group %q, got %q", tc.group, group)
			}
			if version != tc.version {
				t.Errorf("Expected version %q, got %q", tc.version, version)
			}
			if kind != tc.kind {
				t.Errorf("Expected kind %q, got %q", tc.kind, kind)
			}
		})
	}
}

// TestNewClientLua tests the newClientLua function
func TestNewClientLua(t *testing.T) {
	config := &rest.Config{
		Host: "https://localhost:6443",
	}

	L := lua.NewState()
	defer L.Close()

	// Register the class metatable first so newClientLua can wrap the client.
	m := build(config)
	m.PushTo(L)
	L.Pop(1)

	n := newClientLua(L, config)
	if n != 1 {
		t.Errorf("Expected newClientLua to return 1 value, got %d", n)
	}

	// Check that a UserData was pushed
	ud := L.Get(-1)
	if ud.Type() != lua.LTUserData {
		t.Errorf("Expected UserData, got %v", ud.Type())
	}

	// Verify it wraps a *Client
	udVal := ud.(*lua.LUserData)
	if _, ok := udVal.Value.(*Client); !ok {
		t.Errorf("Expected UserData to wrap *Client, got %T", udVal.Value)
	}
}
