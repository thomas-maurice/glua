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
	"encoding/json"
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

// Translator handles conversion between Go values and Lua values
type Translator struct{}

// NewTranslator: creates a new Translator instance
func NewTranslator() *Translator {
	return &Translator{}
}

// ToLua converts an arbitrary Go value to a Lua value.
// It supports primitive types (string, int64, etc.) and complex structs.
// The conversion process:
//  1. Marshal the object to JSON
//  2. Unmarshal into map[string]interface{}
//  3. Use reflect to walk the map and populate LTable
func (t *Translator) ToLua(L *lua.LState, o interface{}) (lua.LValue, error) {
	// Marshal to JSON
	jsonBytes, err := json.Marshal(o)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Unmarshal into map[string]interface{}
	var data interface{}
	if err := json.Unmarshal(jsonBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert to Lua value
	return t.toLuaValue(L, data)
}

// toLuaValue: recursively converts Go values to Lua values using reflection
func (t *Translator) toLuaValue(L *lua.LState, v interface{}) (lua.LValue, error) {
	if v == nil {
		return lua.LNil, nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.String:
		return lua.LString(val.String()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(val.Int()), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return lua.LNumber(val.Uint()), nil

	case reflect.Float32, reflect.Float64:
		return lua.LNumber(val.Float()), nil

	case reflect.Bool:
		return lua.LBool(val.Bool()), nil

	case reflect.Slice, reflect.Array:
		table := L.NewTable()
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			luaVal, err := t.toLuaValue(L, item)
			if err != nil {
				return nil, fmt.Errorf("failed to convert slice element %d: %w", i, err)
			}
			table.Append(luaVal)
		}
		return table, nil

	case reflect.Map:
		table := L.NewTable()
		iter := val.MapRange()
		for iter.Next() {
			key := iter.Key().Interface()
			value := iter.Value().Interface()

			// Convert key to string (Lua tables typically use string keys)
			keyStr := fmt.Sprintf("%v", key)

			luaVal, err := t.toLuaValue(L, value)
			if err != nil {
				return nil, fmt.Errorf("failed to convert map value for key %v: %w", key, err)
			}

			table.RawSetString(keyStr, luaVal)
		}
		return table, nil

	default:
		return nil, fmt.Errorf("unsupported type: %v", val.Kind())
	}
}

// FromLua converts a Lua value to a Go value.
// Works similar to json.Unmarshal - pass a state, Lua value and output object.
// The conversion process:
//  1. Create Go value from Lua value (map[string]interface{} for tables, primitives for others)
//  2. Marshal the value to JSON
//  3. Unmarshal JSON into output object
func (t *Translator) FromLua(L *lua.LState, lv lua.LValue, output interface{}) error {
	// Convert Lua value to Go value
	data, err := t.fromLuaValue(lv)
	if err != nil {
		return fmt.Errorf("failed to convert Lua value: %w", err)
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Unmarshal into output
	if err := json.Unmarshal(jsonBytes, output); err != nil {
		return fmt.Errorf("failed to unmarshal into output: %w", err)
	}

	return nil
}

// fromLuaValue: recursively converts Lua values to Go values
func (t *Translator) fromLuaValue(lv lua.LValue) (interface{}, error) {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil, nil

	case lua.LString:
		return string(v), nil

	case lua.LNumber:
		return float64(v), nil

	case lua.LBool:
		return bool(v), nil

	case *lua.LTable:
		// Check if it's an array or a map
		maxN := v.MaxN()

		// If MaxN > 0, treat it as an array
		if maxN > 0 {
			arr := make([]interface{}, 0, maxN)
			for i := 1; i <= maxN; i++ {
				val := v.RawGetInt(i)
				item, err := t.fromLuaValue(val)
				if err != nil {
					return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
				}
				arr = append(arr, item)
			}
			return arr, nil
		}

		// Otherwise, treat it as a map
		m := make(map[string]interface{})
		v.ForEach(func(key, value lua.LValue) {
			keyStr := key.String()
			val, err := t.fromLuaValue(value)
			if err != nil {
				// Store error for handling outside ForEach
				m["__error__"] = err
				return
			}
			m[keyStr] = val
		})

		// Check if error occurred during ForEach
		if errVal, ok := m["__error__"]; ok {
			return nil, errVal.(error)
		}

		return m, nil

	default:
		return nil, fmt.Errorf("unsupported Lua type: %T", v)
	}
}
