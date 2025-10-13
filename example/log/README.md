# Log Module Example

This example demonstrates how to use the `log` module to add structured logging to your Lua scripts using the beautiful [charmbracelet/log](https://github.com/charmbracelet/log) library.

## Features

- **Zero configuration required** - just load the module and start logging
- **Two usage patterns**: Module-level functions (`log.info()`) or logger objects (`logger:info()`)
- Multiple log levels: debug, info, warn, error, fatal
- Structured logging with key-value pairs
- Logger objects with `logger:with()` for adding context
- Optional context injection from Go (request IDs, user info, etc.)
- Automatic format detection via `LOG_FORMAT` env var (text or json)
- Beautiful colored output (text mode)

## Running the Example

```bash
# Text format (default - pretty colored output)
go run ./example/log

# JSON format (structured logging)
LOG_FORMAT=json go run ./example/log
```

## Quick Start

The simplest way to use the log module:

```go
package main

import (
    logmodule "github.com/thomas-maurice/glua/pkg/modules/log"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    // Just preload the module - that's it!
    L.PreloadModule("log", logmodule.Loader)

    L.DoString(`
        local log = require("log")
        log.info("Hello, world!")
    `)
}
```

## Usage Patterns

### Pattern 1: Module-Level Functions (Simplest)

Use module-level functions for simple logging:

```lua
local log = require("log")

log.debug("Detailed debug information")
log.info("General information")
log.warn("Warning message")
log.error("Error occurred")
log.fatal("Fatal error - exits program")
```

### Pattern 2: Logger Objects

Get a logger object for more control:

```lua
local log = require("log")
local logger = log.logger()

logger:info("Using logger object")
logger:warn("This is a warning")
```

### Adding Context with logger:with()

Create loggers with persistent context fields:

```lua
local log = require("log")
local logger = log.logger()

-- Create a new logger with context
local contextLogger = logger:with("request_id", "abc-123", "session", "xyz-789")

-- Context fields are included in all logs from this logger
contextLogger:info("Processing started")
contextLogger:info("Processing completed")

-- Original logger is unchanged
logger:info("System message")  -- No context fields
```

### Chaining with() Calls

Build up context incrementally:

```lua
local log = require("log")
local logger = log.logger()

local reqLogger = logger:with("request_id", "abc-123")
local userLogger = reqLogger:with("user_id", 456)

reqLogger:info("Request received")      -- has request_id
userLogger:info("User authenticated")   -- has request_id AND user_id
logger:info("System status")            -- no extra fields
```

### Structured Logging

Add key-value pairs to any log call:

```lua
local log = require("log")

log.info("User logged in",
    "user_id", 12345,
    "email", "user@example.com",
    "ip", "192.168.1.1"
)

log.error("Database query failed",
    "query", "SELECT * FROM users",
    "error", "connection timeout",
    "duration_ms", 5000
)
```

## Advanced: Injecting Context from Go

For advanced use cases where you need to inject context from Go (like request IDs, user info, etc.):

```go
package main

import (
    logmodule "github.com/thomas-maurice/glua/pkg/modules/log"
    lua "github.com/yuin/gopher-lua"
)

func processRequest(requestID string, userID int) {
    L := lua.NewState()
    defer L.Close()

    // Get the default logger and add fields from Go
    logger := logmodule.GetDefaultLogger().With(
        "request_id", requestID,
        "user_id", userID,
    )
    logmodule.InjectLogger(L, logger)

    L.PreloadModule("log", logmodule.Loader)

    // All Lua logs will now include request_id and user_id
    L.DoFile("script.lua")
}
```

## Output Formats

### Text Format (Default)

Pretty, colored output with timestamps and caller information:

```
2025-10-13T14:30:45Z INF Application started app=log-example version=1.0.0
```

### JSON Format

Structured JSON logs for log aggregation systems (set `LOG_FORMAT=json`):

```json
{"level":"info","msg":"Application started","time":"2025-10-13T14:30:45Z","caller":"log/log.go:376","app":"log-example","version":"1.0.0"}
```

## Usage Comparison

### Module-Level vs Logger Object

```lua
local log = require("log")

-- Module-level (simpler, uses default logger)
log.info("Simple message")

-- Logger object (more control, can create multiple loggers)
local logger = log.logger()
logger:info("Message from logger object")

-- Create context loggers
local requestLogger = logger:with("request_id", "123")
requestLogger:info("Processing")  -- includes request_id
```

Both patterns work with the same logger under the hood, so you can mix and match:

```lua
local log = require("log")

log.info("Module-level log")

local logger = log.logger()  -- Gets the same default logger
logger:info("Object log")     -- Same output style

-- Create contextual logger
local ctxLogger = logger:with("key", "value")
ctxLogger:info("With context")
```

## Use Cases

1. **Simple Logging**: Use module-level functions for straightforward logging
2. **Contextual Logging**: Use logger objects with `with()` to add request/session context
3. **Request Tracking**: Inject request IDs from Go and track them through Lua execution
4. **Multi-tenant Apps**: Use `logger:with()` to add tenant/user context
5. **Debugging**: Use debug level logs with detailed information
6. **Production Monitoring**: JSON format with log aggregation systems

## Complete Example

```go
package main

import (
    logmodule "github.com/thomas-maurice/glua/pkg/modules/log"
    lua "github.com/yuin/gopher-lua"
)

func main() {
    L := lua.NewState()
    defer L.Close()

    L.PreloadModule("log", logmodule.Loader)

    L.DoString(`
        local log = require("log")

        -- Simple module-level logging
        log.info("Application started")

        -- Get logger object
        local logger = log.logger()
        logger:info("Using logger object")

        -- Create contextual logger
        local requestLogger = logger:with("request_id", "req-123")
        requestLogger:info("Processing request")

        -- Add more context
        local userLogger = requestLogger:with("user_id", 456)
        userLogger:info("User action")  -- has both request_id and user_id

        -- Original logger unchanged
        logger:info("System message")  -- no extra context
    `)
}
```
