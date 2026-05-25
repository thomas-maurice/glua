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
	"fmt"
	"time"

	"github.com/thomas-maurice/glua/pkg/glua"
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GVKMatcher: represents a Kubernetes Group/Version/Kind matcher.
type GVKMatcher struct {
	Group   string `json:"group"`
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

// typeRegistry: manages type registration for stub generation.
var typeRegistry = glua.NewTypeRegistry()

func init() {
	// Register GVKMatcher with the type registry
	if err := typeRegistry.Register(GVKMatcher{}); err != nil {
		panic(fmt.Sprintf("failed to register GVKMatcher: %v", err))
	}

	// Register core Kubernetes types
	types := []interface{}{
		// Core resources
		corev1.Pod{},
		corev1.PodList{},
		corev1.Namespace{},
		corev1.NamespaceList{},
		corev1.Node{},
		corev1.NodeList{},
		corev1.ConfigMap{},
		corev1.ConfigMapList{},
		corev1.Secret{},
		corev1.SecretList{},
		corev1.Service{},
		corev1.ServiceList{},
		corev1.ServiceAccount{},
		corev1.ServiceAccountList{},
		corev1.PersistentVolume{},
		corev1.PersistentVolumeList{},
		corev1.PersistentVolumeClaim{},
		corev1.PersistentVolumeClaimList{},
		// Apps resources
		appsv1.Deployment{},
		appsv1.DeploymentList{},
		appsv1.StatefulSet{},
		appsv1.StatefulSetList{},
		appsv1.DaemonSet{},
		appsv1.DaemonSetList{},
		appsv1.ReplicaSet{},
		appsv1.ReplicaSetList{},
		// Batch resources
		batchv1.Job{},
		batchv1.JobList{},
		batchv1.CronJob{},
		batchv1.CronJobList{},
		// Networking resources
		networkingv1.Ingress{},
		networkingv1.IngressList{},
		networkingv1.NetworkPolicy{},
		networkingv1.NetworkPolicyList{},
		// RBAC resources
		rbacv1.Role{},
		rbacv1.RoleList{},
		rbacv1.ClusterRole{},
		rbacv1.ClusterRoleList{},
		rbacv1.RoleBinding{},
		rbacv1.RoleBindingList{},
		rbacv1.ClusterRoleBinding{},
		rbacv1.ClusterRoleBindingList{},
		// Metav1 types
		metav1.ObjectMeta{},
		metav1.TypeMeta{},
		metav1.Time{},
		metav1.MicroTime{},
		metav1.Duration{},
		metav1.Status{},
		metav1.StatusDetails{},
		metav1.StatusCause{},
		metav1.ListMeta{},
		metav1.OwnerReference{},
		metav1.LabelSelector{},
		metav1.LabelSelectorRequirement{},
	}

	for _, t := range types {
		if err := typeRegistry.Register(t); err != nil {
			panic(fmt.Sprintf("failed to register type %T: %v", t, err))
		}
	}
}

// ParseMemory: parses a Kubernetes memory quantity and returns bytes.
// Raises a Lua error on invalid input.
//
// Example:
//
//	local bytes = k8s.parse_memory("1024Mi")  -- returns 1073741824
func ParseMemory(quantity string) (int64, error) {
	q, err := resource.ParseQuantity(quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to parse memory quantity: %w", err)
	}
	return q.Value(), nil
}

// ParseCPU: parses a Kubernetes CPU quantity and returns millicores.
// Raises a Lua error on invalid input.
//
// Example:
//
//	local millicores = k8s.parse_cpu("100m")  -- returns 100
func ParseCPU(quantity string) (int64, error) {
	q, err := resource.ParseQuantity(quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to parse CPU quantity: %w", err)
	}
	return q.MilliValue(), nil
}

// ParseTime: parses a Kubernetes time string (RFC3339) and returns a Unix timestamp.
// Raises a Lua error on invalid input.
//
// Example:
//
//	local ts = k8s.parse_time("2025-10-03T16:39:00Z")
func ParseTime(timestr string) (int64, error) {
	var k8sTime metav1.Time
	if err := k8sTime.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, timestr))); err != nil {
		t, parseErr := time.Parse(time.RFC3339, timestr)
		if parseErr != nil {
			return 0, fmt.Errorf("failed to parse time: %w", err)
		}
		k8sTime = metav1.NewTime(t)
	}
	return k8sTime.Unix(), nil
}

// FormatTime: converts a Unix timestamp to a Kubernetes time string in RFC3339 format.
//
// Example:
//
//	local timestr = k8s.format_time(1759509540)  -- returns "2025-10-03T16:39:00Z"
func FormatTime(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(time.RFC3339)
}

// ParseDuration: parses a duration string and returns seconds.
// Raises a Lua error on invalid input.
//
// Example:
//
//	local seconds = k8s.parse_duration("5m")  -- returns 300
func ParseDuration(durationStr string) (float64, error) {
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}
	return d.Seconds(), nil
}

// FormatDuration: converts seconds to a duration string.
//
// Example:
//
//	local s = k8s.format_duration(300)  -- returns "5m0s"
func FormatDuration(seconds float64) string {
	return (time.Duration(seconds) * time.Second).String()
}

// MatchGVK: checks if a Kubernetes object (as a map) matches a GVKMatcher.
// The obj map must have "apiVersion" and "kind" keys.
func MatchGVK(obj map[string]any, matcher GVKMatcher) bool {
	if matcher.Kind == "" || matcher.Version == "" {
		return false
	}

	apiVersion, _ := obj["apiVersion"].(string)
	objKind, _ := obj["kind"].(string)

	if objKind != matcher.Kind {
		return false
	}

	var expectedAPIVersion string
	if matcher.Group == "" {
		expectedAPIVersion = matcher.Version
	} else {
		expectedAPIVersion = matcher.Group + "/" + matcher.Version
	}

	return apiVersion == expectedAPIVersion
}

// ensureMetadata: ensures metadata, labels, and annotations keys exist in a map.
// Returns the updated map.
func ensureMetadata(obj map[string]any) map[string]any {
	if obj == nil {
		obj = make(map[string]any)
	}
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		meta = make(map[string]any)
		obj["metadata"] = meta
	}
	if _, ok := meta["labels"]; !ok {
		meta["labels"] = make(map[string]any)
	}
	if _, ok := meta["annotations"]; !ok {
		meta["annotations"] = make(map[string]any)
	}
	return obj
}

// EnsureMetadata: ensures metadata.labels and metadata.annotations exist.
// Returns the updated obj. Callers must assign the return value:
//
//	obj = k8s.ensure_metadata(obj)
func EnsureMetadata(obj map[string]any) map[string]any {
	return ensureMetadata(obj)
}

// InitDefaults: alias for EnsureMetadata. Returns the updated obj.
//
//	obj = k8s.init_defaults(obj)
func InitDefaults(obj map[string]any) map[string]any {
	return ensureMetadata(obj)
}

// AddLabel: adds a label to obj.metadata.labels and returns the updated obj.
//
//	obj = k8s.add_label(obj, "app", "nginx")
func AddLabel(obj map[string]any, key, value string) map[string]any {
	obj = ensureMetadata(obj)
	meta := obj["metadata"].(map[string]any)
	labels := meta["labels"].(map[string]any)
	labels[key] = value
	return obj
}

// AddLabels: adds multiple labels from a table and returns the updated obj.
//
//	obj = k8s.add_labels(obj, {app="nginx", env="prod"})
func AddLabels(obj map[string]any, toAdd map[string]any) map[string]any {
	obj = ensureMetadata(obj)
	meta := obj["metadata"].(map[string]any)
	labels := meta["labels"].(map[string]any)
	for k, v := range toAdd {
		labels[k] = v
	}
	return obj
}

// RemoveLabel: removes a label from obj.metadata.labels and returns the updated obj.
//
//	obj = k8s.remove_label(obj, "app")
func RemoveLabel(obj map[string]any, key string) map[string]any {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return obj
	}
	labels, _ := meta["labels"].(map[string]any)
	if labels == nil {
		return obj
	}
	delete(labels, key)
	return obj
}

// HasLabel: returns true if obj.metadata.labels contains key.
func HasLabel(obj map[string]any, key string) bool {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return false
	}
	labels, _ := meta["labels"].(map[string]any)
	if labels == nil {
		return false
	}
	_, ok := labels[key]
	return ok
}

// GetLabel: returns the value of label key, or "" if absent.
func GetLabel(obj map[string]any, key string) string {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return ""
	}
	labels, _ := meta["labels"].(map[string]any)
	if labels == nil {
		return ""
	}
	v, _ := labels[key].(string)
	return v
}

// AddAnnotation: adds an annotation to obj.metadata.annotations and returns the updated obj.
//
//	obj = k8s.add_annotation(obj, "kubectl.kubernetes.io/last-applied-configuration", "...")
func AddAnnotation(obj map[string]any, key, value string) map[string]any {
	obj = ensureMetadata(obj)
	meta := obj["metadata"].(map[string]any)
	annotations := meta["annotations"].(map[string]any)
	annotations[key] = value
	return obj
}

// AddAnnotations: adds multiple annotations from a table and returns the updated obj.
//
//	obj = k8s.add_annotations(obj, {["app.io/version"]="1.0"})
func AddAnnotations(obj map[string]any, toAdd map[string]any) map[string]any {
	obj = ensureMetadata(obj)
	meta := obj["metadata"].(map[string]any)
	annotations := meta["annotations"].(map[string]any)
	for k, v := range toAdd {
		annotations[k] = v
	}
	return obj
}

// RemoveAnnotation: removes an annotation and returns the updated obj.
//
//	obj = k8s.remove_annotation(obj, "kubectl.kubernetes.io/last-applied-configuration")
func RemoveAnnotation(obj map[string]any, key string) map[string]any {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return obj
	}
	annotations, _ := meta["annotations"].(map[string]any)
	if annotations == nil {
		return obj
	}
	delete(annotations, key)
	return obj
}

// HasAnnotation: returns true if obj.metadata.annotations contains key.
func HasAnnotation(obj map[string]any, key string) bool {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return false
	}
	annotations, _ := meta["annotations"].(map[string]any)
	if annotations == nil {
		return false
	}
	_, ok := annotations[key]
	return ok
}

// GetAnnotation: returns the value of annotation key, or "" if absent.
func GetAnnotation(obj map[string]any, key string) string {
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return ""
	}
	annotations, _ := meta["annotations"].(map[string]any)
	if annotations == nil {
		return ""
	}
	v, _ := annotations[key].(string)
	return v
}

// build: constructs the module definition. Reused by Loader and Register.
func build() *luareg.Module {
	m := luareg.NewModule("kubernetes", "Kubernetes utility functions")
	m.Fn("parse_memory", ParseMemory, "parse a Kubernetes memory quantity, returns bytes",
		luareg.Args("quantity"))
	m.Fn("parse_cpu", ParseCPU, "parse a Kubernetes CPU quantity, returns millicores",
		luareg.Args("quantity"))
	m.Fn("parse_time", ParseTime, "parse an RFC3339 time string, returns Unix timestamp",
		luareg.Args("timestr"))
	m.Fn("format_time", FormatTime, "convert a Unix timestamp to RFC3339 string",
		luareg.Args("timestamp"))
	m.Fn("parse_duration", ParseDuration, "parse a duration string, returns seconds",
		luareg.Args("duration"))
	m.Fn("format_duration", FormatDuration, "convert seconds to a duration string",
		luareg.Args("seconds"))
	m.Fn("match_gvk", MatchGVK, "check if a Kubernetes object matches a GVK matcher",
		luareg.Args("obj", "matcher"))
	m.Fn("ensure_metadata", EnsureMetadata, "ensure metadata.labels and annotations exist, returns updated obj",
		luareg.Args("obj"))
	m.Fn("init_defaults", InitDefaults, "ensure metadata.labels and annotations exist, returns updated obj",
		luareg.Args("obj"))
	m.Fn("add_label", AddLabel, "add a label and return the updated obj",
		luareg.Args("obj", "key", "value"))
	m.Fn("add_labels", AddLabels, "add multiple labels and return the updated obj",
		luareg.Args("obj", "labels"))
	m.Fn("remove_label", RemoveLabel, "remove a label and return the updated obj",
		luareg.Args("obj", "key"))
	m.Fn("has_label", HasLabel, "return true if the label exists",
		luareg.Args("obj", "key"))
	m.Fn("get_label", GetLabel, "return the value of a label, or empty string if absent",
		luareg.Args("obj", "key"))
	m.Fn("add_annotation", AddAnnotation, "add an annotation and return the updated obj",
		luareg.Args("obj", "key", "value"))
	m.Fn("add_annotations", AddAnnotations, "add multiple annotations and return the updated obj",
		luareg.Args("obj", "annotations"))
	m.Fn("remove_annotation", RemoveAnnotation, "remove an annotation and return the updated obj",
		luareg.Args("obj", "key"))
	m.Fn("has_annotation", HasAnnotation, "return true if the annotation exists",
		luareg.Args("obj", "key"))
	m.Fn("get_annotation", GetAnnotation, "return the value of an annotation, or empty string if absent",
		luareg.Args("obj", "key"))
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("kubernetes", kubernetes.Loader).
func Loader(L *lua.LState) int {
	return build().PushTo(L)
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
