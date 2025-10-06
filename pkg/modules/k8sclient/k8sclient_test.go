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

// setupLuaWithClient creates a Lua state with the k8sclient module loaded
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
