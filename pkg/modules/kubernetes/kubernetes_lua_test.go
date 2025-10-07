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
