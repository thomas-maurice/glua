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

package fs

import (
	"fmt"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

// Loader: creates and returns the fs module for Lua.
// This function should be registered with L.PreloadModule("fs", fs.Loader)
//
// @luamodule fs
//
// Example usage in Lua:
//
//	local fs = require("fs")
//	local content, err = fs.read_file("/path/to/file.txt")
//	fs.write_file("/path/to/output.txt", "content")
func Loader(L *lua.LState) int {
	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations
var exports = map[string]lua.LGFunction{
	"read_file":  readFile,
	"write_file": writeFile,
	"exists":     exists,
	"mkdir":      mkdir,
	"mkdir_all":  mkdirAll,
	"remove":     remove,
	"remove_all": removeAll,
	"list":       list,
	"stat":       stat,
}

// readFile: reads the entire contents of a file.
//
// @luafunc read_file
// @luaparam path string The path to the file to read
// @luareturn string The file contents, or nil on error
// @luareturn string|nil Error message if reading failed
//
// Example:
//
//	local content, err = fs.read_file("/etc/config.yaml")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print(content)
//	end
func readFile(L *lua.LState) int {
	path := L.CheckString(1)

	content, err := os.ReadFile(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read file: %v", err)))
		return 2
	}

	L.Push(lua.LString(string(content)))
	L.Push(lua.LNil)
	return 2
}

// writeFile: writes content to a file, creating it if it doesn't exist.
//
// @luafunc write_file
// @luaparam path string The path to the file to write
// @luaparam content string The content to write
// @luareturn string|nil Error message if writing failed, nil on success
//
// Example:
//
//	local err = fs.write_file("/tmp/output.txt", "Hello World")
//	if err then
//	    print("Error: " .. err)
//	end
func writeFile(L *lua.LState) int {
	path := L.CheckString(1)
	content := L.CheckString(2)

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// exists: checks if a file or directory exists.
//
// @luafunc exists
// @luaparam path string The path to check
// @luareturn boolean True if the path exists, false otherwise
//
// Example:
//
//	if fs.exists("/etc/config.yaml") then
//	    print("Config file exists")
//	end
func exists(L *lua.LState) int {
	path := L.CheckString(1)

	_, err := os.Stat(path)
	L.Push(lua.LBool(err == nil))
	return 1
}

// mkdir: creates a directory.
//
// @luafunc mkdir
// @luaparam path string The directory path to create
// @luareturn string|nil Error message if creation failed, nil on success
//
// Example:
//
//	local err = fs.mkdir("/tmp/mydir")
//	if err then
//	    print("Error: " .. err)
//	end
func mkdir(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.Mkdir(path, 0755)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// mkdirAll: creates a directory and all necessary parent directories.
//
// @luafunc mkdir_all
// @luaparam path string The directory path to create
// @luareturn string|nil Error message if creation failed, nil on success
//
// Example:
//
//	local err = fs.mkdir_all("/tmp/path/to/nested/dir")
//	if err then
//	    print("Error: " .. err)
//	end
func mkdirAll(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.MkdirAll(path, 0755)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to create directories: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// remove: removes a file or empty directory.
//
// @luafunc remove
// @luaparam path string The path to remove
// @luareturn string|nil Error message if removal failed, nil on success
//
// Example:
//
//	local err = fs.remove("/tmp/file.txt")
//	if err then
//	    print("Error: " .. err)
//	end
func remove(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.Remove(path)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to remove: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// removeAll: removes a path and all its contents recursively.
//
// @luafunc remove_all
// @luaparam path string The path to remove recursively
// @luareturn string|nil Error message if removal failed, nil on success
//
// Example:
//
//	local err = fs.remove_all("/tmp/mydir")
//	if err then
//	    print("Error: " .. err)
//	end
func removeAll(L *lua.LState) int {
	path := L.CheckString(1)

	err := os.RemoveAll(path)
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("failed to remove recursively: %v", err)))
		return 1
	}

	L.Push(lua.LNil)
	return 1
}

// list: lists all entries in a directory.
//
// @luafunc list
// @luaparam path string The directory path to list
// @luareturn table Array of entry names, or nil on error
// @luareturn string|nil Error message if listing failed
//
// Example:
//
//	local entries, err = fs.list("/tmp")
//	if err then
//	    print("Error: " .. err)
//	else
//	    for i, entry in ipairs(entries) do
//	        print(entry)
//	    end
//	end
func list(L *lua.LState) int {
	path := L.CheckString(1)

	entries, err := os.ReadDir(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to list directory: %v", err)))
		return 2
	}

	tbl := L.NewTable()
	for i, entry := range entries {
		tbl.RawSetInt(i+1, lua.LString(entry.Name()))
	}

	L.Push(tbl)
	L.Push(lua.LNil)
	return 2
}

// stat: gets information about a file or directory.
//
// @luafunc stat
// @luaparam path string The path to stat
// @luareturn table File info with name, size, is_dir, mode, mod_time, or nil on error
// @luareturn string|nil Error message if stat failed
//
// Example:
//
//	local info, err = fs.stat("/etc/config.yaml")
//	if err then
//	    print("Error: " .. err)
//	else
//	    print("Name: " .. info.name)
//	    print("Size: " .. info.size .. " bytes")
//	    print("Is directory: " .. tostring(info.is_dir))
//	end
func stat(L *lua.LState) int {
	path := L.CheckString(1)

	info, err := os.Stat(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stat: %v", err)))
		return 2
	}

	tbl := L.NewTable()
	tbl.RawSetString("name", lua.LString(filepath.Base(path)))
	tbl.RawSetString("size", lua.LNumber(info.Size()))
	tbl.RawSetString("is_dir", lua.LBool(info.IsDir()))
	tbl.RawSetString("mode", lua.LNumber(info.Mode()))
	tbl.RawSetString("mod_time", lua.LNumber(info.ModTime().Unix()))

	L.Push(tbl)
	L.Push(lua.LNil)
	return 2
}
