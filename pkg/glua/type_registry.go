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

package glua

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// TypeInfo stores information about a registered Go type for Lua stub generation
type TypeInfo struct {
	Name       string                // The Lua-friendly type name (e.g., "corev1.Pod")
	GoType     reflect.Type          // The original Go type
	Fields     map[string]*FieldInfo // Map of field name to field information
	IsArray    bool                  // Whether this is an array type
	ElementKey string                // For arrays, the type key of the element
}

// FieldInfo stores information about a struct field for Lua stub generation
type FieldInfo struct {
	Name    string // The field name (from JSON tag)
	TypeKey string // The Lua type annotation (e.g., "string", "number", "corev1.Container")
	IsArray bool   // Whether this field is an array
}

// TypeRegistry manages type registration and stub generation for Lua.
// It processes Go types recursively and generates Lua LSP annotations.
type TypeRegistry struct {
	types map[string]*TypeInfo // Map of type key to type information (prevents duplicates)
	queue []interface{}        // Queue of objects to process (for discovering types)
}

// NewTypeRegistry: creates a new TypeRegistry instance
func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		types: make(map[string]*TypeInfo),
		queue: make([]interface{}, 0),
	}
}

// Register: adds a Go type to the registry for stub generation.
// The object is queued for processing to discover all dependent types.
// Returns an error if the object is nil or not a valid type.
func (r *TypeRegistry) Register(obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("cannot register nil object")
	}

	t := reflect.TypeOf(obj)
	if t == nil {
		return fmt.Errorf("cannot determine type of object")
	}

	r.queue = append(r.queue, obj)
	return nil
}

// getTypeName: generates a human-readable type name for a Go type.
// Handles Kubernetes API objects specially (e.g., corev1.Pod instead of v1.Pod).
func (r *TypeRegistry) getTypeName(t reflect.Type) string {
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// For primitive types
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Map:
		return "table"
	case reflect.Slice, reflect.Array:
		return "table"
	}

	// For named types with package path
	pkgPath := t.PkgPath()
	typeName := t.Name()

	if pkgPath == "" {
		return typeName
	}

	// Extract package name from path
	parts := strings.Split(pkgPath, "/")
	pkgName := parts[len(parts)-1]

	// Handle Kubernetes API groups specially
	// e.g., k8s.io/api/core/v1 -> corev1
	if strings.Contains(pkgPath, "k8s.io/api/") {
		apiParts := strings.Split(pkgPath, "/")
		for i, part := range apiParts {
			if part == "api" && i+2 < len(apiParts) {
				group := apiParts[i+1]
				version := apiParts[i+2]
				// Combine group and version (e.g., core + v1 = corev1)
				return group + version + "." + typeName
			}
		}
	}

	return pkgName + "." + typeName
}

// getTypeKey: generates a unique key for a type for internal tracking.
// This prevents duplicate type registration and handles circular dependencies.
func (r *TypeRegistry) getTypeKey(t reflect.Type) string {
	// Handle pointers
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Use full package path + type name for uniqueness
	pkgPath := t.PkgPath()
	typeName := t.Name()

	if pkgPath == "" {
		return typeName
	}

	return pkgPath + "." + typeName
}

// processType: recursively processes a type and registers it.
// Returns the Lua type annotation string for the type.
// Handles circular dependencies by checking if a type is already registered.
func (r *TypeRegistry) processType(t reflect.Type) string {
	t = r.unwrapPointer(t)

	if primType := r.getPrimitiveType(t); primType != "" {
		return primType
	}

	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return r.processArrayType(t)
	}

	if t.Kind() == reflect.Map {
		return r.processMapType(t)
	}

	if t.Kind() == reflect.Struct {
		return r.processStructType(t)
	}

	return "any"
}

// unwrapPointer: unwraps pointer types
func (r *TypeRegistry) unwrapPointer(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}

// getPrimitiveType: returns Lua type name for primitive Go types
func (r *TypeRegistry) getPrimitiveType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Interface:
		return "any"
	}
	return ""
}

// processArrayType: processes slice and array types
func (r *TypeRegistry) processArrayType(t reflect.Type) string {
	elemType := t.Elem()
	elemKey := r.processType(elemType)
	return elemKey + "[]"
}

// processMapType: processes map types
func (r *TypeRegistry) processMapType(t reflect.Type) string {
	valueType := t.Elem()
	valueKey := r.processType(valueType)
	return "table<string, " + valueKey + ">"
}

// processStructType: processes struct types and registers them
func (r *TypeRegistry) processStructType(t reflect.Type) string {
	typeKey := r.getTypeKey(t)

	if _, exists := r.types[typeKey]; exists {
		return r.getTypeName(t)
	}

	typeName := r.getTypeName(t)
	typeInfo := &TypeInfo{
		Name:   typeName,
		GoType: t,
		Fields: make(map[string]*FieldInfo),
	}
	r.types[typeKey] = typeInfo

	r.processStructFields(t, typeInfo)
	return typeName
}

// processStructFields: processes all fields in a struct
func (r *TypeRegistry) processStructFields(t reflect.Type, typeInfo *TypeInfo) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		fieldName := strings.Split(jsonTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		fieldTypeKey := r.processType(field.Type)
		typeInfo.Fields[fieldName] = &FieldInfo{
			Name:    fieldName,
			TypeKey: fieldTypeKey,
		}
	}
}

// Process: processes all registered types and discovers dependencies.
// Call this after registering all root types with Register().
func (r *TypeRegistry) Process() error {
	for len(r.queue) > 0 {
		obj := r.queue[0]
		r.queue = r.queue[1:]

		t := reflect.TypeOf(obj)
		r.processType(t)
	}

	return nil
}

// GenerateStubs: generates Lua annotation stubs for all registered types.
// Returns a string containing ---@class and ---@field annotations.
//
// Example output:
//
//	---@class corev1.Pod
//	---@field metadata corev1.ObjectMeta
//	---@field spec corev1.PodSpec
func (r *TypeRegistry) GenerateStubs() (string, error) {
	var sb strings.Builder

	// Sort type names for consistent output
	var typeKeys []string
	for key := range r.types {
		typeKeys = append(typeKeys, key)
	}
	sort.Strings(typeKeys)

	// Generate class definitions
	for _, key := range typeKeys {
		typeInfo := r.types[key]

		sb.WriteString(fmt.Sprintf("---@class %s\n", typeInfo.Name))

		// Sort field names for consistent output
		var fieldNames []string
		for fieldName := range typeInfo.Fields {
			fieldNames = append(fieldNames, fieldName)
		}
		sort.Strings(fieldNames)

		// Generate field annotations
		for _, fieldName := range fieldNames {
			field := typeInfo.Fields[fieldName]
			sb.WriteString(fmt.Sprintf("---@field %s %s\n", field.Name, field.TypeKey))
		}

		sb.WriteString("\n")
	}

	// Create a types table to export all type names for convenience
	sb.WriteString("-- Export all types for convenience\n")
	sb.WriteString("local types = {\n")
	for _, key := range typeKeys {
		typeInfo := r.types[key]
		sb.WriteString(fmt.Sprintf("  %s = {},\n", typeInfo.Name))
	}
	sb.WriteString("}\n\n")
	sb.WriteString("return types\n")

	return sb.String(), nil
}
