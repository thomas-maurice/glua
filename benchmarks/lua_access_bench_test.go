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

package benchmarks

import (
	"testing"

	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BenchmarkLuaFieldAccess: benchmarks accessing fields from converted objects
func BenchmarkLuaFieldAccess(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Simple struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	obj := Simple{
		Name:  "John Doe",
		Age:   42,
		Email: "john@example.com",
	}

	lv, err := translator.ToLua(L, obj)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("obj", lv)

	luaCode := `
		local name = obj.name
		local age = obj.age
		local email = obj.email
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaNestedFieldAccess: benchmarks accessing nested fields
func BenchmarkLuaNestedFieldAccess(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
				"env": "prod",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.21",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("512Mi"),
						},
					},
				},
			},
		},
	}

	lv, err := translator.ToLua(L, pod)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("pod", lv)

	luaCode := `
		local name = pod.metadata.name
		local namespace = pod.metadata.namespace
		local app_label = pod.metadata.labels.app
		local container_name = pod.spec.containers[1].name
		local cpu_limit = pod.spec.containers[1].resources.limits.cpu
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaArrayIteration: benchmarks iterating over arrays
func BenchmarkLuaArrayIteration(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Container struct {
		Items []string `json:"items"`
	}

	obj := Container{
		Items: make([]string, 100),
	}
	for i := 0; i < 100; i++ {
		obj.Items[i] = "item"
	}

	lv, err := translator.ToLua(L, obj)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("obj", lv)

	luaCode := `
		local count = 0
		for i, item in ipairs(obj.items) do
			count = count + 1
		end
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaMapIteration: benchmarks iterating over maps
func BenchmarkLuaMapIteration(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	obj := map[string]string{
		"key1":  "value1",
		"key2":  "value2",
		"key3":  "value3",
		"key4":  "value4",
		"key5":  "value5",
		"key6":  "value6",
		"key7":  "value7",
		"key8":  "value8",
		"key9":  "value9",
		"key10": "value10",
	}

	lv, err := translator.ToLua(L, obj)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("obj", lv)

	luaCode := `
		local count = 0
		for k, v in pairs(obj) do
			count = count + 1
		end
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaFieldModification: benchmarks modifying fields in Lua
func BenchmarkLuaFieldModification(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Simple struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Count int    `json:"count"`
	}

	obj := Simple{
		Name:  "John Doe",
		Age:   42,
		Count: 0,
	}

	lv, err := translator.ToLua(L, obj)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("obj", lv)

	luaCode := `
		obj.name = "Jane Doe"
		obj.age = obj.age + 1
		obj.count = obj.count + 1
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaComplexOperation: benchmarks complex Lua operations on converted objects
func BenchmarkLuaComplexOperation(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app":     "test",
				"env":     "prod",
				"version": "1.0.0",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.21",
					Env: []corev1.EnvVar{
						{Name: "ENV1", Value: "value1"},
						{Name: "ENV2", Value: "value2"},
						{Name: "ENV3", Value: "value3"},
					},
				},
				{
					Name:  "sidecar",
					Image: "busybox",
					Env: []corev1.EnvVar{
						{Name: "ENV4", Value: "value4"},
					},
				},
			},
		},
	}

	lv, err := translator.ToLua(L, pod)
	if err != nil {
		b.Fatal(err)
	}
	L.SetGlobal("pod", lv)

	luaCode := `
		-- Count labels
		local label_count = 0
		for k, v in pairs(pod.metadata.labels) do
			label_count = label_count + 1
		end

		-- Count total env vars across all containers
		local total_env_vars = 0
		for i, container in ipairs(pod.spec.containers) do
			if container.env then
				total_env_vars = total_env_vars + #container.env
			end
		end

		-- Check if production environment
		local is_prod = pod.metadata.labels.env == "prod"

		-- Modify container images (simulation)
		for i, container in ipairs(pod.spec.containers) do
			container.image = container.image .. ":latest"
		end
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := L.DoString(luaCode); err != nil {
			b.Fatal(err)
		}
	}
}
