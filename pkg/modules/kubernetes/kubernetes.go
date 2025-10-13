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

var (
	// translator: handles Go-Lua type conversion
	translator = glua.NewTranslator()
	// typeRegistry: manages type registration for stub generation
	typeRegistry = glua.NewTypeRegistry()
)

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

// @luaclass v1.Pod

// @luaclass v1.PodList

// @luaclass v1.Namespace

// @luaclass v1.NamespaceList

// @luaclass v1.Node

// @luaclass v1.NodeList

// @luaclass v1.ConfigMap

// @luaclass v1.ConfigMapList

// @luaclass v1.Secret

// @luaclass v1.SecretList

// @luaclass v1.Service

// @luaclass v1.ServiceList

// @luaclass v1.ServiceAccount

// @luaclass v1.ServiceAccountList

// @luaclass v1.PersistentVolume

// @luaclass v1.PersistentVolumeList

// @luaclass v1.PersistentVolumeClaim

// @luaclass v1.PersistentVolumeClaimList

// @luaclass v1.Deployment

// @luaclass v1.DeploymentList

// @luaclass v1.StatefulSet

// @luaclass v1.StatefulSetList

// @luaclass v1.DaemonSet

// @luaclass v1.DaemonSetList

// @luaclass v1.ReplicaSet

// @luaclass v1.ReplicaSetList

// @luaclass v1.Job

// @luaclass v1.JobList

// @luaclass v1.CronJob

// @luaclass v1.CronJobList

// @luaclass v1.Ingress

// @luaclass v1.IngressList

// @luaclass v1.NetworkPolicy

// @luaclass v1.NetworkPolicyList

// @luaclass v1.Role

// @luaclass v1.RoleList

// @luaclass v1.ClusterRole

// @luaclass v1.ClusterRoleList

// @luaclass v1.RoleBinding

// @luaclass v1.RoleBindingList

// @luaclass v1.ClusterRoleBinding

// @luaclass v1.ClusterRoleBindingList

// @luaclass v1.ObjectMeta

// @luaclass v1.TypeMeta

// @luaclass v1.Time

// @luaclass v1.MicroTime

// @luaclass v1.Duration

// @luaclass v1.Status

// @luaclass v1.StatusDetails

// @luaclass v1.StatusCause

// @luaclass v1.ListMeta

// @luaclass v1.OwnerReference

// @luaclass v1.LabelSelector

// @luaclass v1.LabelSelectorRequirement

// Loader: creates and returns the kubernetes module for Lua.
// This function should be registered with L.PreloadModule("kubernetes", kubernetes.Loader)
//
// @luamodule kubernetes
//
// Example usage in Lua:
//
//	local k8s = require("kubernetes")
//	local bytes = k8s.parse_memory("1024Mi")
//	local millicores = k8s.parse_cpu("100m")
//	local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")
//	local timestr = k8s.format_time(1759509540)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"parse_memory":    parseMemory,
	"parse_cpu":       parseCPU,
	"parse_time":      parseTime,
	"format_time":     formatTime,
	"init_defaults":   initDefaults,
	"parse_duration":  parseDuration,
	"format_duration": formatDuration,
	"match_gvk":       matchGVK,
}

// parseMemory: parses a Kubernetes memory quantity (e.g., "1024Mi", "1Gi", "512M") and returns bytes as a number.
// Returns nil and error message on failure.
//
// @luafunc parse_memory
// @luaparam quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
// @luareturn bytes number The memory value in bytes, or nil on error
// @luareturn err string|nil Error message if parsing failed
//
// Example:
//
//	local bytes = k8s.parse_memory("1024Mi")  -- returns 1073741824
func parseMemory(L *lua.LState) int {
	str := L.CheckString(1)

	quantity, err := resource.ParseQuantity(str)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse memory quantity: %v", err)))
		return 2
	}

	// Get value in bytes
	bytes := quantity.Value()

	L.Push(lua.LNumber(bytes))
	return 1
}

// parseCPU: parses a Kubernetes CPU quantity (e.g., "100m", "1", "2000m") and returns millicores as a number.
// Returns nil and error message on failure.
//
// @luafunc parse_cpu
// @luaparam quantity string The CPU quantity to parse (e.g., "100m", "1", "2000m")
// @luareturn millicores number The CPU value in millicores, or nil on error
// @luareturn err string|nil Error message if parsing failed
//
// Example:
//
//	local millicores = k8s.parse_cpu("100m")  -- returns 100
//	local millicores = k8s.parse_cpu("1")     -- returns 1000
func parseCPU(L *lua.LState) int {
	str := L.CheckString(1)

	quantity, err := resource.ParseQuantity(str)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse CPU quantity: %v", err)))
		return 2
	}

	// Get value in millicores
	millicores := quantity.MilliValue()

	L.Push(lua.LNumber(millicores))
	return 1
}

// parseTime: parses a Kubernetes time string (RFC3339 format like "2025-10-03T16:39:00Z") and returns a Unix timestamp.
// Returns nil and error message on failure.
//
// @luafunc parse_time
// @luaparam timestr string The time string in RFC3339 format (e.g., "2025-10-03T16:39:00Z")
// @luareturn timestamp number The Unix timestamp, or nil on error
// @luareturn err string|nil Error message if parsing failed
//
// Example:
//
//	local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")  -- returns Unix timestamp
func parseTime(L *lua.LState) int {
	str := L.CheckString(1)

	// Parse using Kubernetes Time format
	var k8sTime metav1.Time
	if err := k8sTime.UnmarshalJSON([]byte(fmt.Sprintf(`"%s"`, str))); err != nil {
		// Try standard RFC3339 parsing as fallback
		t, parseErr := time.Parse(time.RFC3339, str)
		if parseErr != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to parse time: %v", err)))
			return 2
		}
		k8sTime = metav1.NewTime(t)
	}

	// Return Unix timestamp
	L.Push(lua.LNumber(k8sTime.Unix()))
	return 1
}

// formatTime: converts a Unix timestamp (int64) to a Kubernetes time string in RFC3339 format.
// Returns nil and error message on failure.
//
// @luafunc format_time
// @luaparam timestamp number The Unix timestamp to convert
// @luareturn timestr string The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
// @luareturn err string|nil Error message if formatting failed
//
// Example:
//
//	local timestr = k8s.format_time(1759509540)  -- returns "2025-10-03T16:39:00Z"
func formatTime(L *lua.LState) int {
	timestamp := L.CheckNumber(1)

	// Convert to time.Time
	t := time.Unix(int64(timestamp), 0).UTC()

	// Format as RFC3339 (Kubernetes standard format)
	formatted := t.Format(time.RFC3339)

	L.Push(lua.LString(formatted))
	L.Push(lua.LNil)
	return 2
}

// initDefaults: initializes default empty tables for metadata.labels and metadata.annotations
// if they are nil. This is useful for ensuring these fields are tables instead of nil,
// making it easier to add labels/annotations in Lua without checking for nil first.
//
// @luafunc init_defaults
// @luaparam obj table The Kubernetes object (must have a metadata field)
// @luareturn obj table The same object with initialized defaults (modified in-place)
//
// Example:
//
//	local k8s = require("kubernetes")
//	k8s.init_defaults(myPod)
//	myPod.metadata.labels.app = "myapp"  -- safe even if labels was nil before
func initDefaults(L *lua.LState) int {
	obj := L.CheckTable(1)

	// Get metadata field
	metadata := L.GetField(obj, "metadata")
	if metadata == lua.LNil {
		// If metadata doesn't exist, create it
		metadata = L.NewTable()
		L.SetField(obj, "metadata", metadata)
	}

	metadataTable, ok := metadata.(*lua.LTable)
	if !ok {
		L.Push(obj)
		return 1
	}

	// Initialize labels if nil
	labels := L.GetField(metadataTable, "labels")
	if labels == lua.LNil {
		L.SetField(metadataTable, "labels", L.NewTable())
	}

	// Initialize annotations if nil
	annotations := L.GetField(metadataTable, "annotations")
	if annotations == lua.LNil {
		L.SetField(metadataTable, "annotations", L.NewTable())
	}

	L.Push(obj)
	return 1
}

// parseDuration: parses a Kubernetes duration string (e.g., "5s", "10m", "2h") and returns seconds as a number.
// Returns nil and error message on failure.
//
// @luafunc parse_duration
// @luaparam duration string The duration string to parse (e.g., "5s", "10m", "2h")
// @luareturn seconds number The duration value in seconds, or nil on error
// @luareturn err string|nil Error message if parsing failed
//
// Example:
//
//	local seconds = k8s.parse_duration("5m")  -- returns 300
//	local seconds = k8s.parse_duration("1h30m")  -- returns 5400
func parseDuration(L *lua.LState) int {
	str := L.CheckString(1)

	duration, err := time.ParseDuration(str)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse duration: %v", err)))
		return 2
	}

	// Return duration in seconds
	L.Push(lua.LNumber(duration.Seconds()))
	L.Push(lua.LNil)
	return 2
}

// formatDuration: converts seconds to a Kubernetes duration string.
// Returns nil and error message on failure.
//
// @luafunc format_duration
// @luaparam seconds number The duration in seconds to convert
// @luareturn duration string The duration string (e.g., "5m0s", "1h30m0s"), or nil on error
// @luareturn err string|nil Error message if formatting failed
//
// Example:
//
//	local duration_str = k8s.format_duration(300)  -- returns "5m0s"
//	local duration_str = k8s.format_duration(5400)  -- returns "1h30m0s"
func formatDuration(L *lua.LState) int {
	seconds := L.CheckNumber(1)

	duration := time.Duration(seconds) * time.Second
	formatted := duration.String()

	L.Push(lua.LString(formatted))
	L.Push(lua.LNil)
	return 2
}

// matchGVK: checks if a Kubernetes object matches the specified Group/Version/Kind matcher.
// Returns true if the object's apiVersion and kind match the matcher's values.
//
// @luafunc match_gvk
// @luaparam obj table The Kubernetes object to check
// @luaparam matcher kubernetes.GVKMatcher The GVK matcher with group, version, and kind fields
// @luareturn matches boolean true if the GVK matches
//
// Example:
//
//	local matcher = {group = "", version = "v1", kind = "Pod"}
//	local matches = k8s.match_gvk(pod, matcher)  -- returns true for a Pod
func matchGVK(L *lua.LState) int {
	obj := L.CheckTable(1)
	matcherTable := L.CheckTable(2)

	// Convert Lua table to GVKMatcher
	var matcher GVKMatcher
	if err := translator.FromLua(L, matcherTable, &matcher); err != nil {
		L.RaiseError("failed to parse GVKMatcher: %v", err)
		return 0
	}

	// Validate required fields
	if matcher.Kind == "" {
		L.RaiseError("GVKMatcher requires 'kind' field")
		return 0
	}
	if matcher.Version == "" {
		L.RaiseError("GVKMatcher requires 'version' field")
		return 0
	}

	// Get apiVersion and kind from the object
	apiVersion := L.GetField(obj, "apiVersion").String()
	objKind := L.GetField(obj, "kind").String()

	// Check kind first
	if objKind != matcher.Kind {
		L.Push(lua.LFalse)
		return 1
	}

	// Build expected apiVersion
	var expectedAPIVersion string
	if matcher.Group == "" {
		expectedAPIVersion = matcher.Version
	} else {
		expectedAPIVersion = matcher.Group + "/" + matcher.Version
	}

	// Check apiVersion
	if apiVersion == expectedAPIVersion {
		L.Push(lua.LTrue)
	} else {
		L.Push(lua.LFalse)
	}

	return 1
}
