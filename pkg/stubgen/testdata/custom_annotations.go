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
