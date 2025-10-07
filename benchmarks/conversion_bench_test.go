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
	"time"

	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BenchmarkGoToLuaSimple: benchmarks simple struct conversion to Lua
func BenchmarkGoToLuaSimple(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := translator.ToLua(L, obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGoToLuaComplex: benchmarks complex nested struct conversion
func BenchmarkGoToLuaComplex(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Nested struct {
		Value string   `json:"value"`
		Items []string `json:"items"`
	}

	type Complex struct {
		StringField string            `json:"stringField"`
		IntField    int               `json:"intField"`
		BoolField   bool              `json:"boolField"`
		SliceField  []string          `json:"sliceField"`
		MapField    map[string]int    `json:"mapField"`
		Nested      *Nested           `json:"nested"`
		Tags        map[string]string `json:"tags"`
	}

	obj := Complex{
		StringField: "test",
		IntField:    42,
		BoolField:   true,
		SliceField:  []string{"a", "b", "c", "d", "e"},
		MapField:    map[string]int{"x": 1, "y": 2, "z": 3},
		Nested: &Nested{
			Value: "nested",
			Items: []string{"item1", "item2", "item3"},
		},
		Tags: map[string]string{
			"env":  "prod",
			"team": "platform",
			"app":  "api",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := translator.ToLua(L, obj)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGoToLuaPod: benchmarks Kubernetes Pod conversion
func BenchmarkGoToLuaPod(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "benchmark-pod",
			Namespace:         "default",
			CreationTimestamp: metav1.Time{Time: time.Now()},
			Labels: map[string]string{
				"app":     "test",
				"env":     "prod",
				"version": "1.0.0",
			},
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/port":   "8080",
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
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("100m"),
							corev1.ResourceMemory: resource.MustParse("128Mi"),
						},
					},
					Env: []corev1.EnvVar{
						{Name: "ENV1", Value: "value1"},
						{Name: "ENV2", Value: "value2"},
						{Name: "ENV3", Value: "value3"},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := translator.ToLua(L, pod)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaToGoSimple: benchmarks simple Lua to Go conversion
func BenchmarkLuaToGoSimple(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Simple struct {
		Name  string `json:"name"`
		Age   int    `json:"age"`
		Email string `json:"email"`
	}

	luaCode := `
		return {
			name = "John Doe",
			age = 42,
			email = "john@example.com"
		}
	`

	if err := L.DoString(luaCode); err != nil {
		b.Fatal(err)
	}
	lv := L.Get(-1)
	L.Pop(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Simple
		if err := translator.FromLua(L, lv, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLuaToGoComplex: benchmarks complex Lua to Go conversion
func BenchmarkLuaToGoComplex(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	type Nested struct {
		Value string   `json:"value"`
		Items []string `json:"items"`
	}

	type Complex struct {
		StringField string            `json:"stringField"`
		IntField    int               `json:"intField"`
		BoolField   bool              `json:"boolField"`
		SliceField  []string          `json:"sliceField"`
		MapField    map[string]int    `json:"mapField"`
		Nested      *Nested           `json:"nested"`
		Tags        map[string]string `json:"tags"`
	}

	luaCode := `
		return {
			stringField = "test",
			intField = 42,
			boolField = true,
			sliceField = {"a", "b", "c", "d", "e"},
			mapField = {x = 1, y = 2, z = 3},
			nested = {
				value = "nested",
				items = {"item1", "item2", "item3"}
			},
			tags = {
				env = "prod",
				team = "platform",
				app = "api"
			}
		}
	`

	if err := L.DoString(luaCode); err != nil {
		b.Fatal(err)
	}
	lv := L.Get(-1)
	L.Pop(1)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result Complex
		if err := translator.FromLua(L, lv, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRoundTripSimple: benchmarks full round-trip conversion
func BenchmarkRoundTripSimple(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lv, err := translator.ToLua(L, obj)
		if err != nil {
			b.Fatal(err)
		}

		var result Simple
		if err := translator.FromLua(L, lv, &result); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkRoundTripPod: benchmarks full Kubernetes Pod round-trip
func BenchmarkRoundTripPod(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	translator := glua.NewTranslator()

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:              "benchmark-pod",
			Namespace:         "default",
			CreationTimestamp: metav1.Time{Time: time.Now()},
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lv, err := translator.ToLua(L, pod)
		if err != nil {
			b.Fatal(err)
		}

		var result corev1.Pod
		if err := translator.FromLua(L, lv, &result); err != nil {
			b.Fatal(err)
		}
	}
}
