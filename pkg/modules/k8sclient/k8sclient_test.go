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
	// Creating the loader with a valid-looking config should succeed
	// (it only fails when actually making API calls)
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

// setupLuaWithClient creates a Lua state with the k8sclient module loaded and test constants
func setupLuaWithClient(client *Client) *lua.LState {
	L := lua.NewState()

	// Create module table with client methods
	exports := map[string]lua.LGFunction{
		"get":    client.get,
		"create": client.create,
		"update": client.update,
		"delete": client.delete,
		"list":   client.list,
	}

	mod := L.SetFuncs(L.NewTable(), exports)
	L.SetGlobal("client", mod)

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

// runLuaTestFile runs a Lua test file and returns the result
func runLuaTestFile(t *testing.T, filename string) {
	t.Helper()

	client := setupFakeClient()
	L := setupLuaWithClient(client)
	defer L.Close()

	testPath := filepath.Join("testdata", filename)

	if err := L.DoFile(testPath); err != nil {
		t.Fatalf("Test script %s failed: %v", filename, err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Errorf("Test script %s returned %v, expected true", filename, result)
	}
}

// TestGet tests the get operation
func TestGet(t *testing.T) {
	runLuaTestFile(t, "test_get.lua")
}

// TestCreate tests the create operation
func TestCreate(t *testing.T) {
	runLuaTestFile(t, "test_create.lua")
}

// TestUpdate tests the update operation
func TestUpdate(t *testing.T) {
	runLuaTestFile(t, "test_update.lua")
}

// TestDelete tests the delete operation
func TestDelete(t *testing.T) {
	runLuaTestFile(t, "test_delete.lua")
}

// TestList tests the list operation
func TestList(t *testing.T) {
	runLuaTestFile(t, "test_list.lua")
}

// TestIntegration tests a complete CRUD workflow
func TestIntegration(t *testing.T) {
	runLuaTestFile(t, "test_integration.lua")
}

// TestErrorMissingKind tests error handling for missing kind
func TestErrorMissingKind(t *testing.T) {
	runLuaTestFile(t, "test_error_missing_kind.lua")
}

// TestErrorMissingVersion tests error handling for missing version
func TestErrorMissingVersion(t *testing.T) {
	runLuaTestFile(t, "test_error_missing_version.lua")
}

// TestErrorNoAPIVersion tests error handling for missing apiVersion
func TestErrorNoAPIVersion(t *testing.T) {
	runLuaTestFile(t, "test_error_no_apiversion.lua")
}

// TestErrorNoKind tests error handling for missing kind in object
func TestErrorNoKind(t *testing.T) {
	runLuaTestFile(t, "test_error_no_kind.lua")
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

// TestCreateGVKTable tests the createGVKTable helper function
func TestCreateGVKTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	table := createGVKTable(L, "apps", "v1", "Deployment")

	if table == nil {
		t.Fatal("createGVKTable returned nil")
	}

	group := L.GetField(table, "group")
	if group.String() != "apps" {
		t.Errorf("Expected group 'apps', got %q", group.String())
	}

	version := L.GetField(table, "version")
	if version.String() != "v1" {
		t.Errorf("Expected version 'v1', got %q", version.String())
	}

	kind := L.GetField(table, "kind")
	if kind.String() != "Deployment" {
		t.Errorf("Expected kind 'Deployment', got %q", kind.String())
	}
}

// TestAddGVKConstants tests that GVK constants are added correctly
func TestAddGVKConstants(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	mod := L.NewTable()
	addGVKConstants(L, mod)

	// Test a few constants
	constants := []struct {
		name    string
		group   string
		version string
		kind    string
	}{
		{"POD", "", "v1", "Pod"},
		{"DEPLOYMENT", "apps", "v1", "Deployment"},
		{"INGRESS", "networking.k8s.io", "v1", "Ingress"},
		{"ROLE", "rbac.authorization.k8s.io", "v1", "Role"},
	}

	for _, tc := range constants {
		t.Run(tc.name, func(t *testing.T) {
			gvkValue := L.GetField(mod, tc.name)
			if gvkValue.Type() != lua.LTTable {
				t.Fatalf("Expected %s to be a table, got %v", tc.name, gvkValue.Type())
			}

			gvkTable := gvkValue.(*lua.LTable)
			group := L.GetField(gvkTable, "group")
			if group.String() != tc.group {
				t.Errorf("Expected group %q, got %q", tc.group, group.String())
			}

			version := L.GetField(gvkTable, "version")
			if version.String() != tc.version {
				t.Errorf("Expected version %q, got %q", tc.version, version.String())
			}

			kind := L.GetField(gvkTable, "kind")
			if kind.String() != tc.kind {
				t.Errorf("Expected kind %q, got %q", tc.kind, kind.String())
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

	n := newClientLua(L, config)
	if n != 2 {
		t.Errorf("Expected newClientLua to return 2 values, got %d", n)
	}

	// Check that client table was pushed
	client := L.Get(-2)
	if client.Type() != lua.LTTable {
		t.Errorf("Expected client to be a table, got %v", client.Type())
	}

	// Check that error is nil
	err := L.Get(-1)
	if err != lua.LNil {
		t.Errorf("Expected error to be nil, got %v", err)
	}

	// Verify client has methods
	clientTable := client.(*lua.LTable)
	methods := []string{"get", "create", "update", "delete", "list"}
	for _, method := range methods {
		field := L.GetField(clientTable, method)
		if field.Type() != lua.LTFunction {
			t.Errorf("Expected %s to be a function, got %v", method, field.Type())
		}
	}
}
