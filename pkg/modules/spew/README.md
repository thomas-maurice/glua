# Spew Module

The `spew` module provides deep pretty-printing and debugging functionality for Lua values using [go-spew](https://github.com/davecgh/go-spew). It's invaluable for debugging complex table structures and understanding data at runtime.

## Functions

### `spew.dump(value)`

Prints a detailed representation of a Lua value to stdout with syntax highlighting and type information.

**Parameters:**
- `value` (any): The Lua value to dump (table, string, number, boolean, nil, etc.)

**Returns:**
- None (prints directly to stdout)

**Example:**
```lua
local spew = require("spew")

local data = {
    name = "John",
    age = 30,
    tags = {"developer", "golang"}
}

spew.dump(data)
-- Output:
-- (map[string]interface {}) (len=3) {
--  (string) (len=4) "name": (string) (len=4) "John",
--  (string) (len=3) "age": (float64) 30,
--  (string) (len=4) "tags": ([]interface {}) (len=2) {
--   (string) (len=9) "developer",
--   (string) (len=6) "golang"
--  }
-- }
```

### `spew.sdump(value)`

Returns a detailed string representation of a Lua value. Unlike `dump()`, this returns the string instead of printing to stdout.

**Parameters:**
- `value` (any): The Lua value to dump

**Returns:**
- `string`: A detailed string representation of the value with type information

**Example:**
```lua
local spew = require("spew")

local data = {
    status = "active",
    count = 42
}

local result = spew.sdump(data)
print(result)
-- Or save it to a file, log it, etc.
```

## Features

### Deep Inspection

Spew recursively inspects nested structures:

```lua
local spew = require("spew")

local complex = {
    server = {
        host = "localhost",
        port = 8080,
        endpoints = {
            api = "/api/v1",
            health = "/health"
        }
    },
    clients = {
        {name = "client1", active = true},
        {name = "client2", active = false}
    }
}

spew.dump(complex)
-- Shows full nested structure with types and lengths
```

### Type Information

Every value includes its Go type:

```lua
local spew = require("spew")

spew.dump({
    str = "hello",        -- (string) (len=5) "hello"
    num = 42,             -- (float64) 42
    bool = true,          -- (bool) true
    list = {1, 2, 3}      -- ([]interface {}) (len=3) {...}
})
```

### Array vs Map Detection

Automatically distinguishes between arrays and maps:

```lua
local spew = require("spew")

-- Array (consecutive integer keys starting at 1)
spew.dump({1, 2, 3})
-- ([]interface {}) (len=3) {
--  (float64) 1,
--  (float64) 2,
--  (float64) 3
-- }

-- Map (string keys or non-consecutive integers)
spew.dump({name = "test", value = 123})
-- (map[string]interface {}) (len=2) {
--  (string) (len=4) "name": (string) (len=4) "test",
--  (string) (len=5) "value": (float64) 123
-- }
```

## Usage in Go

```go
package main

import (
    "github.com/thomas-maurice/glua/pkg/modules/spew"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // Register the spew module
    L.PreloadModule("spew", spew.Loader)

    // Use in Lua
    L.DoString(`
        local spew = require("spew")

        -- Debug a complex data structure
        local config = {
            database = {
                host = "db.example.com",
                port = 5432,
                credentials = {
                    user = "admin",
                    -- password redacted
                }
            },
            cache = {
                enabled = true,
                ttl = 3600
            }
        }

        print("=== Configuration Debug ===")
        spew.dump(config)

        -- Get string representation for logging
        local configStr = spew.sdump(config)
        -- Now you can log configStr, save to file, etc.
    `)
}
```

## Common Use Cases

### 1. Debugging Pod Structures

```lua
local spew = require("spew")

-- Inspect a Kubernetes pod structure
print("=== Pod Structure ===")
spew.dump(myPod)

-- See exactly what fields are available
-- and their types
```

### 2. Comparing Data

```lua
local spew = require("spew")

local before = spew.sdump(originalData)
local after = spew.sdump(modifiedData)

if before ~= after then
    print("Data changed!")
    print("Before:", before)
    print("After:", after)
end
```

### 3. Logging Complex Errors

```lua
local spew = require("spew")

local function processRequest(request)
    if not validateRequest(request) then
        local dump = spew.sdump(request)
        error("Invalid request structure:\n" .. dump)
    end
end
```

### 4. Exploring Unknown Data

```lua
local spew = require("spew")
local json = require("json")

-- Parse unknown JSON and explore its structure
local data, err = json.parse(unknownJsonString)
if not err then
    print("=== Exploring Unknown Data ===")
    spew.dump(data)
    -- Now you can see exactly what fields and types exist
end
```

## Differences from `dump()` vs `sdump()`

| Feature | `dump()` | `sdump()` |
|---------|----------|-----------|
| Output | Prints to stdout | Returns string |
| Use case | Quick debugging | Logging, string manipulation |
| Return value | None | String |
| Performance | Slightly faster | Slight overhead for string building |

## Performance Considerations

- Spew is intended for debugging and development
- For production logging, consider using structured logging instead
- Large data structures can produce verbose output
- Use `sdump()` when you need to control where output goes

## Testing

Run the test suite:

```bash
go test ./pkg/modules/spew/
```

## Integration

The spew module is automatically included when you run `make stubgen` and will generate IDE autocomplete stubs in `library/spew.lua`.

## Credits

Built on top of [go-spew](https://github.com/davecgh/go-spew) by Dave Collins.
