# JSON Module

The `json` module provides JSON parsing and serialization functionality for Lua scripts.

## Functions

### `json.parse(jsonstr)`

Parses a JSON string and returns a Lua table.

**Parameters:**
- `jsonstr` (string): The JSON string to parse

**Returns:**
- `table`: The parsed JSON as a Lua table (or nil on error)
- `string|nil`: Error message if parsing failed

**Example:**
```lua
local json = require("json")
local tbl, err = json.parse('{"name":"John","age":30}')

if err then
    print("Error: " .. err)
else
    print(tbl.name)  -- prints "John"
    print(tbl.age)   -- prints 30
end
```

### `json.stringify(tbl)`

Converts a Lua table to a JSON string.

**Parameters:**
- `tbl` (table): The Lua table to convert to JSON

**Returns:**
- `string`: The JSON string (or nil on error)
- `string|nil`: Error message if conversion failed

**Example:**
```lua
local json = require("json")
local str, err = json.stringify({name="Jane", age=25})

if err then
    print("Error: " .. err)
else
    print(str)  -- prints '{"age":25,"name":"Jane"}'
end
```

## Data Type Mapping

### JSON to Lua (parse)

| JSON Type | Lua Type |
|-----------|----------|
| object    | table (string keys) |
| array     | table (1-indexed numeric keys) |
| string    | string |
| number    | number |
| boolean   | boolean |
| null      | nil |

### Lua to JSON (stringify)

| Lua Type | JSON Type |
|----------|-----------|
| table with consecutive integer keys (1-indexed) | array |
| table with string keys | object |
| string | string |
| number | number |
| boolean | boolean |
| nil | null |

## Usage in Go

```go
package main

import (
    "github.com/thomas-maurice/glua/pkg/modules/json"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // Register the json module
    L.PreloadModule("json", json.Loader)

    // Use in Lua
    L.DoString(`
        local json = require("json")

        -- Parse JSON
        local data, err = json.parse('{"users":[{"name":"Alice"},{"name":"Bob"}]}')
        if not err then
            for i, user in ipairs(data.users) do
                print(user.name)
            end
        end

        -- Stringify table
        local jsonStr, err2 = json.stringify({
            message = "Hello",
            count = 42,
            items = {1, 2, 3}
        })
        if not err2 then
            print(jsonStr)
        end
    `)
}
```

## Array vs Object Detection

The module automatically detects whether a Lua table should be serialized as a JSON array or object:

- **Arrays**: Tables with consecutive integer keys starting from 1
- **Objects**: Tables with string keys or non-consecutive integer keys

**Examples:**

```lua
local json = require("json")

-- These become JSON arrays
json.stringify({1, 2, 3})           -- [1,2,3]
json.stringify({"a", "b", "c"})     -- ["a","b","c"]

-- These become JSON objects
json.stringify({name="John"})       -- {"name":"John"}
json.stringify({[1]="a", [5]="b"})  -- {"1":"a","5":"b"}  (non-consecutive)
```

## Error Handling

Both functions return two values: the result and an error message. Always check for errors:

```lua
local json = require("json")

local result, err = json.parse(jsonString)
if err then
    print("Parse error: " .. err)
    return
end

-- Use result safely
print(result.field)
```

## Testing

Run the test suite:

```bash
go test ./pkg/modules/json/
```

## Integration

The json module is automatically included when you run `make stubgen` and will generate IDE autocomplete stubs in `library/json.gen.lua`.
