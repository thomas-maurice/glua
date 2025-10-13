# Spew Module Example

This example demonstrates the `spew` module for debugging and pretty-printing Lua data structures with colored JSON output.

## Features

The spew module provides:

- **`spew.dump(value)`** - Pretty-prints any Lua value as colored JSON to stdout
- **`spew.sdump(value)`** - Returns a JSON string representation of any Lua value

## Running the Example

From the example/spew directory:

```bash
cd example/spew
go run .
```

Or from the repository root:

```bash
go run ./example/spew
```

## What the Demo Shows

The demo script (`demo.lua`) demonstrates:

1. **Simple objects** - Basic key-value tables
2. **Nested structures** - Deeply nested tables with multiple levels
3. **Arrays** - Lua tables as arrays
4. **Array of objects** - Tables containing multiple structured objects
5. **Complex nested structures** - Real-world API response-like data
6. **String output** - Using `sdump()` to get JSON as a string

## Use Cases

The spew module is perfect for:

- **Debugging Lua scripts** - Quickly inspect table contents
- **Testing** - Verify data structures in tests
- **Logging** - Output structured data in readable format
- **Development** - Understand complex nested data structures

## Example Output

The output uses colored JSON with syntax highlighting for:

- Keys (field names)
- Strings (green)
- Numbers (cyan)
- Booleans (yellow)
- Null values (gray)

Perfect for terminal-based development and debugging workflows.
