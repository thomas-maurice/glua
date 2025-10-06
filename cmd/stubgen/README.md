# stubgen - Lua Module Stub Generator

A code analyzer that scans Go source files for Lua module definitions and generates Lua LSP annotation stubs.

## Usage

```bash
go run ./cmd/stubgen -dir <source_directory> -output <output_file>
```

### Options

- `-dir` - Directory to scan for Go module files (default: ".")
- `-output` - Output file for combined stubs (default: "module_stubs.gen.lua")
- `-output-dir` - Output directory for per-module stub files (RECOMMENDED for Neovim/LSP)

### Example

```bash
# Generate per-module files (RECOMMENDED for LSP autocomplete)
go run ./cmd/stubgen -dir pkg/modules -output-dir library
# Creates: library/kubernetes.lua, library/mymodule.lua, etc.

# Generate combined stub file (alternative)
go run ./cmd/stubgen -dir pkg/modules -output stubs.lua

# Generate stubs for a specific module
go run ./cmd/stubgen -dir pkg/modules/kubernetes -output-dir library
```

**Important for Neovim users**: Use `-output-dir library` to generate per-module files. This allows Lua LSP to properly recognize `require("kubernetes")` statements.

## Annotation Format

To make your Go Lua modules discoverable by stubgen, add structured comments:

### Module Declaration

```go
// @luamodule <module_name>
func Loader(L *lua.LState) int {
    // ...
}
```

### Function Declaration

```go
// @luafunc <function_name>
// @luaparam <param_name> <type> <description>
// @luareturn <type> <description>
func myFunc(L *lua.LState) int {
    // ...
}
```

### Complete Example

```go
package mymodule

import lua "github.com/yuin/gopher-lua"

// Loader: creates the mymodule module
//
// @luamodule mymodule
func Loader(L *lua.LState) int {
    mod := L.SetFuncs(L.NewTable(), exports)
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "add": add,
}

// add: adds two numbers
//
// @luafunc add
// @luaparam a number First number
// @luaparam b number Second number
// @luareturn number The sum
func add(L *lua.LState) int {
    a := L.CheckNumber(1)
    b := L.CheckNumber(2)
    L.Push(lua.LNumber(a + b))
    return 1
}
```

This will generate:

```lua
---@meta

--- mymodule module
---@class mymodule
local mymodule = {}

--- add: adds two numbers
---@param a number First number
---@param b number Second number
---@return number The sum
function mymodule.add(a, b) end

return mymodule
```

**Note**: The `---@meta` annotation at the top is crucial - it tells Lua LSP that this is a definition file for type checking and autocomplete. Without it, the LSP won't recognize `require("mymodule")` properly in Neovim.

## Supported Types

- `string` - Lua string
- `number` - Lua number
- `boolean` - Lua boolean
- `table` - Lua table
- `any` - Any type
- `string|nil` - Union types (e.g., string or nil)
- Custom types (any valid Lua LSP type annotation)
