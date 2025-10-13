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

package kubernetes

import (
	"testing"
	"time"

	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestKubernetesModuleLoading: tests that kubernetes module loads correctly
func TestKubernetesModuleLoading(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")
		assert(k8s ~= nil, "kubernetes module is nil")
		assert(type(k8s.parse_memory) == "function", "parse_memory is not a function")
		assert(type(k8s.parse_cpu) == "function", "parse_cpu is not a function")
		assert(type(k8s.parse_time) == "function", "parse_time is not a function")
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("Kubernetes module loading failed: %v", err)
	}
}

// TestKubernetesParseMemory: tests memory parsing in Lua context
func TestKubernetesParseMemory(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"parse 1Ki", "1Ki", 1024, false},
		{"parse 1Mi", "1Mi", 1024 * 1024, false},
		{"parse 1Gi", "1Gi", 1024 * 1024 * 1024, false},
		{"parse 100Mi", "100Mi", 100 * 1024 * 1024, false},
		{"parse 512Mi", "512Mi", 512 * 1024 * 1024, false},
		{"parse 2Gi", "2Gi", 2 * 1024 * 1024 * 1024, false},
		{"parse invalid", "invalid", 0, true},
		{"parse empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("input", lua.LString(tt.input))
			L.SetGlobal("expected", lua.LNumber(tt.expected))

			var luaCode string
			if tt.wantErr {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_memory(input)
					assert(err ~= nil, "expected error but got nil")
					assert(result == nil, "expected nil result on error")
				`
			} else {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_memory(input)
					assert(err == nil, "unexpected error: " .. tostring(err))
					assert(result == expected, string.format("expected %d, got %d", expected, result))
				`
			}

			if err := L.DoString(luaCode); err != nil {
				t.Errorf("Lua test failed: %v", err)
			}
		})
	}
}

// TestKubernetesParseCPU: tests CPU parsing in Lua context
func TestKubernetesParseCPU(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"parse 1m", "1m", 1, false},
		{"parse 100m", "100m", 100, false},
		{"parse 500m", "500m", 500, false},
		{"parse 1", "1", 1000, false},
		{"parse 2", "2", 2000, false},
		{"parse 0.5", "0.5", 500, false},
		{"parse invalid", "invalid", 0, true},
		{"parse empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("input", lua.LString(tt.input))
			L.SetGlobal("expected", lua.LNumber(tt.expected))

			var luaCode string
			if tt.wantErr {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_cpu(input)
					assert(err ~= nil, "expected error but got nil")
					assert(result == nil, "expected nil result on error")
				`
			} else {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_cpu(input)
					assert(err == nil, "unexpected error: " .. tostring(err))
					assert(result == expected, string.format("expected %d, got %d", expected, result))
				`
			}

			if err := L.DoString(luaCode); err != nil {
				t.Errorf("Lua test failed: %v", err)
			}
		})
	}
}

// TestKubernetesParseTime: tests timestamp parsing in Lua context
func TestKubernetesParseTime(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	tests := []struct {
		name     string
		input    string
		expected int64
		wantErr  bool
	}{
		{"parse RFC3339", "2025-10-03T16:39:00Z", 1759509540, false},
		{"parse with timezone", "2025-01-01T00:00:00Z", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC).Unix(), false},
		{"parse invalid", "not-a-date", 0, true},
		{"parse empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L.SetGlobal("input", lua.LString(tt.input))
			L.SetGlobal("expected", lua.LNumber(tt.expected))

			var luaCode string
			if tt.wantErr {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_time(input)
					assert(err ~= nil, "expected error but got nil")
					assert(result == nil, "expected nil result on error")
				`
			} else {
				luaCode = `
					local k8s = require("kubernetes")
					local result, err = k8s.parse_time(input)
					assert(err == nil, "unexpected error: " .. tostring(err))
					assert(result == expected, string.format("expected %d, got %d", expected, result))
				`
			}

			if err := L.DoString(luaCode); err != nil {
				t.Errorf("Lua test failed: %v", err)
			}
		})
	}
}

// TestKubernetesModuleWithPodData: tests kubernetes module with real Pod data
func TestKubernetesModuleWithPodData(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)
	translator := glua.NewTranslator()

	// Create a realistic Pod
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "test-app",
			Namespace:         "production",
			CreationTimestamp: metav1.Time{Time: time.Date(2025, 10, 3, 16, 39, 0, 0, time.UTC)},
			Labels: map[string]string{
				"app": "test-app",
				"env": "production",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "app-container",
					Image: "myapp:v1.0.0",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("1"),
							corev1.ResourceMemory: resource.MustParse("1Gi"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
				},
			},
		},
	}

	// Convert to Lua
	luaPod, err := translator.ToLua(L, pod)
	if err != nil {
		t.Fatalf("ToLua() error = %v", err)
	}

	L.SetGlobal("pod", luaPod)

	// Test parsing in Lua
	luaCode := `
		local k8s = require("kubernetes")

		-- Parse CPU limit (should be 1000m)
		local cpu_limit, err = k8s.parse_cpu(pod.spec.containers[1].resources.limits.cpu)
		assert(err == nil, "failed to parse CPU limit: " .. tostring(err))
		assert(cpu_limit == 1000, string.format("CPU limit should be 1000m, got %d", cpu_limit))

		-- Parse CPU request (should be 250m)
		local cpu_request, err = k8s.parse_cpu(pod.spec.containers[1].resources.requests.cpu)
		assert(err == nil, "failed to parse CPU request: " .. tostring(err))
		assert(cpu_request == 250, string.format("CPU request should be 250m, got %d", cpu_request))

		-- Parse memory limit (should be 1Gi = 1073741824 bytes)
		local mem_limit, err = k8s.parse_memory(pod.spec.containers[1].resources.limits.memory)
		assert(err == nil, "failed to parse memory limit: " .. tostring(err))
		assert(mem_limit == 1073741824, string.format("Memory limit should be 1073741824, got %d", mem_limit))

		-- Parse memory request (should be 256Mi = 268435456 bytes)
		local mem_request, err = k8s.parse_memory(pod.spec.containers[1].resources.requests.memory)
		assert(err == nil, "failed to parse memory request: " .. tostring(err))
		assert(mem_request == 268435456, string.format("Memory request should be 268435456, got %d", mem_request))

		-- Parse timestamp
		local timestamp, err = k8s.parse_time(pod.metadata.creationTimestamp)
		assert(err == nil, "failed to parse timestamp: " .. tostring(err))
		assert(timestamp == 1759509540, string.format("Timestamp should be 1759509540, got %d", timestamp))

		-- Calculate some metrics
		local cpu_ratio = cpu_limit / cpu_request
		assert(cpu_ratio == 4, string.format("CPU ratio should be 4, got %f", cpu_ratio))

		local mem_ratio = mem_limit / mem_request
		assert(mem_ratio == 4, string.format("Memory ratio should be 4, got %f", mem_ratio))

		-- Return success
		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("Lua integration test failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if result != lua.LTrue {
		t.Error("Expected Lua script to return true")
	}
}

// TestComplexLuaOperations: tests complex operations in Lua with converted data
func TestComplexLuaOperations(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)
	translator := glua.NewTranslator()

	// Create multiple pods
	pods := []corev1.Pod{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pod-1"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "container-1",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pod-2"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "container-2",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("200Mi"),
							},
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "pod-3"},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "container-3",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("150Mi"),
							},
						},
					},
				},
			},
		},
	}

	luaPods, err := translator.ToLua(L, pods)
	if err != nil {
		t.Fatalf("ToLua() error = %v", err)
	}

	L.SetGlobal("pods", luaPods)

	// Calculate total memory across all pods in Lua
	luaCode := `
		local k8s = require("kubernetes")

		local total_memory = 0
		for i, pod in ipairs(pods) do
			for j, container in ipairs(pod.spec.containers) do
				local mem_str = container.resources.requests.memory
				if mem_str then
					local mem_bytes, err = k8s.parse_memory(mem_str)
					if err then
						error("Failed to parse memory: " .. err)
					end
					total_memory = total_memory + mem_bytes
				end
			end
		end

		-- Total should be 100Mi + 200Mi + 150Mi = 450Mi = 471859200 bytes
		local expected = 471859200
		assert(total_memory == expected, string.format("Total memory should be %d, got %d", expected, total_memory))

		return total_memory
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("Lua calculation failed: %v", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	if num, ok := result.(lua.LNumber); !ok || int64(num) != 471859200 {
		t.Errorf("Expected total memory 471859200, got %v", result)
	}
}

// TestEnsureMetadata: tests the ensure_metadata function
func TestEnsureMetadata(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		-- Create empty pod
		local pod = {kind = "Pod"}

		-- Ensure metadata
		k8s.ensure_metadata(pod)

		assert(pod.metadata ~= nil, "metadata should not be nil")
		assert(pod.metadata.labels ~= nil, "labels should not be nil")
		assert(pod.metadata.annotations ~= nil, "annotations should not be nil")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("ensure_metadata test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected ensure_metadata to succeed")
	}
}

// TestAddLabel: tests the add_label function
func TestAddLabel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add single label
		k8s.add_label(pod, "app", "nginx")
		k8s.add_label(pod, "version", "1.0")

		assert(pod.metadata.labels.app == "nginx", "app label should be nginx")
		assert(pod.metadata.labels.version == "1.0", "version label should be 1.0")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("add_label test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected add_label to succeed")
	}
}

// TestAddLabels: tests the add_labels function
func TestAddLabels(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add multiple labels
		k8s.add_labels(pod, {
			app = "nginx",
			version = "1.0",
			tier = "frontend"
		})

		assert(pod.metadata.labels.app == "nginx", "app label should be nginx")
		assert(pod.metadata.labels.version == "1.0", "version label should be 1.0")
		assert(pod.metadata.labels.tier == "frontend", "tier label should be frontend")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("add_labels test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected add_labels to succeed")
	}
}

// TestRemoveLabel: tests the remove_label function
func TestRemoveLabel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add labels
		k8s.add_labels(pod, {app = "nginx", version = "1.0"})

		-- Remove one label
		k8s.remove_label(pod, "version")

		assert(pod.metadata.labels.app == "nginx", "app label should still exist")
		assert(pod.metadata.labels.version == nil, "version label should be removed")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("remove_label test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected remove_label to succeed")
	}
}

// TestHasLabel: tests the has_label function
func TestHasLabel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}
		k8s.add_label(pod, "app", "nginx")

		assert(k8s.has_label(pod, "app") == true, "should have app label")
		assert(k8s.has_label(pod, "missing") == false, "should not have missing label")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("has_label test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected has_label to succeed")
	}
}

// TestGetLabel: tests the get_label function
func TestGetLabel(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}
		k8s.add_label(pod, "app", "nginx")

		local app = k8s.get_label(pod, "app")
		assert(app == "nginx", "app label value should be nginx")

		local missing = k8s.get_label(pod, "missing")
		assert(missing == nil, "missing label should return nil")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("get_label test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected get_label to succeed")
	}
}

// TestAddAnnotation: tests the add_annotation function
func TestAddAnnotation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add single annotation
		k8s.add_annotation(pod, "description", "My nginx pod")
		k8s.add_annotation(pod, "owner", "team-backend")

		assert(pod.metadata.annotations.description == "My nginx pod", "description annotation should match")
		assert(pod.metadata.annotations.owner == "team-backend", "owner annotation should match")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("add_annotation test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected add_annotation to succeed")
	}
}

// TestAddAnnotations: tests the add_annotations function
func TestAddAnnotations(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add multiple annotations
		k8s.add_annotations(pod, {
			description = "My nginx pod",
			owner = "team-backend",
			version = "1.2.3"
		})

		assert(pod.metadata.annotations.description == "My nginx pod", "description annotation should match")
		assert(pod.metadata.annotations.owner == "team-backend", "owner annotation should match")
		assert(pod.metadata.annotations.version == "1.2.3", "version annotation should match")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("add_annotations test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected add_annotations to succeed")
	}
}

// TestRemoveAnnotation: tests the remove_annotation function
func TestRemoveAnnotation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Add annotations
		k8s.add_annotations(pod, {description = "My pod", owner = "team"})

		-- Remove one annotation
		k8s.remove_annotation(pod, "owner")

		assert(pod.metadata.annotations.description == "My pod", "description should still exist")
		assert(pod.metadata.annotations.owner == nil, "owner annotation should be removed")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("remove_annotation test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected remove_annotation to succeed")
	}
}

// TestHasAnnotation: tests the has_annotation function
func TestHasAnnotation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}
		k8s.add_annotation(pod, "description", "My pod")

		assert(k8s.has_annotation(pod, "description") == true, "should have description annotation")
		assert(k8s.has_annotation(pod, "missing") == false, "should not have missing annotation")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("has_annotation test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected has_annotation to succeed")
	}
}

// TestGetAnnotation: tests the get_annotation function
func TestGetAnnotation(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}
		k8s.add_annotation(pod, "description", "My nginx pod")

		local desc = k8s.get_annotation(pod, "description")
		assert(desc == "My nginx pod", "description value should match")

		local missing = k8s.get_annotation(pod, "missing")
		assert(missing == nil, "missing annotation should return nil")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("get_annotation test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected get_annotation to succeed")
	}
}

// TestLabelAndAnnotationChaining: tests that functions can be chained
func TestLabelAndAnnotationChaining(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("kubernetes", Loader)

	luaCode := `
		local k8s = require("kubernetes")

		local pod = {kind = "Pod"}

		-- Test chaining (though Lua doesn't use return values for this pattern)
		k8s.add_label(pod, "app", "nginx")
		k8s.add_annotation(pod, "description", "Web server")

		assert(k8s.has_label(pod, "app"), "should have app label")
		assert(k8s.has_annotation(pod, "description"), "should have description annotation")

		return true
	`

	if err := L.DoString(luaCode); err != nil {
		t.Fatalf("chaining test failed: %v", err)
	}

	result := L.Get(-1)
	if result != lua.LTrue {
		t.Error("Expected chaining to succeed")
	}
}
