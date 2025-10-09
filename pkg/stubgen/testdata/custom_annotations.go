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

package testdata

import lua "github.com/yuin/gopher-lua"

// LoadCustomAnnotationsModule: loads the custom_annotations module
//
// @luamodule custom_annotations
// @luaannotation @alias ID string|number
// @luaannotation @alias Handler fun(id: ID): boolean
func LoadCustomAnnotationsModule(L *lua.LState) int {
	mod := L.NewTable()

	L.SetField(mod, "process_id", L.NewFunction(processID))
	L.SetField(mod, "process_typed_id", L.NewFunction(processTypedID))

	L.Push(mod)
	return 1
}

// processID: processes an ID value
//
// @luafunc process_id
// @luaparam id ID the identifier to process
// @luareturn boolean true if valid
// @luaannotation @deprecated Use process_typed_id instead
// @luaannotation @nodiscard
func processID(L *lua.LState) int {
	L.Push(lua.LBool(true))
	return 1
}

// processTypedID: processes a typed ID value with generics
//
// @luafunc process_typed_id
// @luaparam id any the identifier to process
// @luareturn boolean true if valid
// @luaannotation @generic T
// @luaannotation @param id `T`
// @luaannotation @return boolean
func processTypedID(L *lua.LState) int {
	L.Push(lua.LBool(true))
	return 1
}
