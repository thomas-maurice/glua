package kubernetes

import (
	"fmt"
	"time"

	lua "github.com/yuin/gopher-lua"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Loader: creates and returns the kubernetes module for Lua.
// This function should be registered with L.PreloadModule("kubernetes", kubernetes.Loader)
//
// @luamodule kubernetes
//
// Example usage in Lua:
//   local k8s = require("kubernetes")
//   local bytes = k8s.parse_memory("1024Mi")
//   local millicores = k8s.parse_cpu("100m")
//   local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")
//   local timestr = k8s.format_time(1759509540)
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"parse_memory":  parseMemory,
	"parse_cpu":     parseCPU,
	"parse_time":    parseTime,
	"format_time":   formatTime,
	"init_defaults": initDefaults,
}

// parseMemory: parses a Kubernetes memory quantity (e.g., "1024Mi", "1Gi", "512M") and returns bytes as a number.
// Returns nil and error message on failure.
//
// @luafunc parse_memory
// @luaparam quantity string The memory quantity to parse (e.g., "1024Mi", "1Gi")
// @luareturn number The memory value in bytes, or nil on error
// @luareturn string|nil Error message if parsing failed
//
// Example:
//   local bytes = k8s.parse_memory("1024Mi")  -- returns 1073741824
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
// @luareturn number The CPU value in millicores, or nil on error
// @luareturn string|nil Error message if parsing failed
//
// Example:
//   local millicores = k8s.parse_cpu("100m")  -- returns 100
//   local millicores = k8s.parse_cpu("1")     -- returns 1000
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
// @luareturn number The Unix timestamp, or nil on error
// @luareturn string|nil Error message if parsing failed
//
// Example:
//   local timestamp = k8s.parse_time("2025-10-03T16:39:00Z")  -- returns Unix timestamp
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
// @luareturn string The time in RFC3339 format (e.g., "2025-10-03T16:39:00Z"), or nil on error
// @luareturn string|nil Error message if formatting failed
//
// Example:
//   local timestr = k8s.format_time(1759509540)  -- returns "2025-10-03T16:39:00Z"
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
// @luareturn table The same object with initialized defaults (modified in-place)
//
// Example:
//   local k8s = require("kubernetes")
//   k8s.init_defaults(myPod)
//   myPod.metadata.labels.app = "myapp"  -- safe even if labels was nil before
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
