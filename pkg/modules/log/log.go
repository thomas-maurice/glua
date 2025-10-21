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

package log

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/thomas-maurice/glua/pkg/glua"
	lua "github.com/yuin/gopher-lua"
)

const (
	loggerTypeName = "Logger"
	loggerKey      = "__logger__"
)

var (
	// defaultLogger: the default logger instance used when no logger is injected
	defaultLogger *log.Logger
	// defaultLoggerOnce: ensures default logger is initialized only once
	defaultLoggerOnce sync.Once
	// translator: reusable translator for Lua to Go conversions
	translator = glua.NewTranslator()
)

// GetDefaultLogger: returns the default logger instance, initializing it if necessary.
// The default logger uses:
//   - os.Stderr for output
//   - RFC3339 timestamp format
//   - Caller reporting enabled
//   - TextFormatter by default, or JSONFormatter if LOG_FORMAT=json env var is set
//
// You can use this to get the default logger and add fields before injecting:
//
//	logger := logmodule.GetDefaultLogger().With("app", "myapp")
//	logmodule.InjectLogger(L, logger)
func GetDefaultLogger() *log.Logger {
	defaultLoggerOnce.Do(func() {
		formatter := log.TextFormatter
		if os.Getenv("LOG_FORMAT") == "json" {
			formatter = log.JSONFormatter
		}
		defaultLogger = log.NewWithOptions(os.Stderr, log.Options{
			TimeFormat:      time.RFC3339,
			ReportTimestamp: true,
			ReportCaller:    true,
			Formatter:       formatter,
		})
	})
	return defaultLogger
}

// getDefaultLoggerUserData: returns the default logger wrapped in UserData
func getDefaultLoggerUserData(L *lua.LState) *lua.LUserData {
	lv := L.GetGlobal(loggerKey)
	if ud, ok := lv.(*lua.LUserData); ok {
		return ud
	}
	// Create and cache default logger UserData
	ud := wrapLogger(L, GetDefaultLogger())
	L.SetGlobal(loggerKey, ud)
	return ud
}

// InjectLogger: injects a pre-configured logger instance into the Lua state.
// This is OPTIONAL - if you don't inject a logger, a default one will be created automatically.
// Use this only when you want to add pre-set fields (like request ID, user ID, etc.) from Go.
//
// Example usage:
//
//	// Get default logger and add fields
//	logger := logmodule.GetDefaultLogger().With("request_id", "abc123", "user", "john")
//	logmodule.InjectLogger(L, logger)
//
// Or create a custom logger:
//
//	logger := logmodule.NewLogger(os.Stderr, true, time.RFC3339, log.TextFormatter)
//	logger = logger.With("app", "myapp", "version", "1.0.0")
//	logmodule.InjectLogger(L, logger)
func InjectLogger(L *lua.LState, logger *log.Logger) {
	ud := wrapLogger(L, logger)
	L.SetGlobal(loggerKey, ud)
}

// wrapLogger: wraps a log.Logger in UserData with the Logger metatable
func wrapLogger(L *lua.LState, logger *log.Logger) *lua.LUserData {
	// Ensure Logger type is registered
	registerLoggerType(L)

	ud := L.NewUserData()
	ud.Value = logger
	L.SetMetatable(ud, L.GetTypeMetatable(loggerTypeName))
	return ud
}

// registerLoggerType: ensures the Logger type metatable is registered
func registerLoggerType(L *lua.LState) {
	mt := L.GetTypeMetatable(loggerTypeName)
	if mt == lua.LNil {
		mt = L.NewTypeMetatable(loggerTypeName)
		L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), getLoggerMethods()))
	}
}

// getLoggerMethods: returns the logger methods map (lazy initialization to avoid init cycle)
func getLoggerMethods() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"debug": loggerDebug,
		"info":  loggerInfo,
		"warn":  loggerWarn,
		"error": loggerError,
		"fatal": loggerFatal,
		"with":  loggerWith,
	}
}

// checkLogger: extracts a log.Logger from UserData
func checkLogger(L *lua.LState, index int) *log.Logger {
	ud := L.CheckUserData(index)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logger
	}
	L.ArgError(index, "Logger expected")
	return nil
}

// Loader: creates and returns the log module for Lua.
// This function should be registered with L.PreloadModule("log", log.Loader)
//
// @luamodule log
//
// Example usage in Lua:
//
//	local log = require("log")
//	log.info("Application started")
//	log.warn("Low memory warning")
//	log.error("Failed to connect to database")
//
//	-- Or get a logger object
//	local logger = log.logger()
//	logger:info("Using logger object")
//	local contextLogger = logger:with("request_id", "abc123")
//	contextLogger:info("Request processing")
func Loader(L *lua.LState) int {
	// Register Logger type
	registerLoggerType(L)

	// Create module table
	mod := L.SetFuncs(L.NewTable(), exports)

	// Push module onto stack
	L.Push(mod)
	return 1
}

// exports: maps Lua function names to Go implementations (module-level functions)
var exports = map[string]lua.LGFunction{
	"debug":  moduleDebug,
	"info":   moduleInfo,
	"warn":   moduleWarn,
	"error":  moduleError,
	"fatal":  moduleFatal,
	"logger": moduleLogger,
}

// Module-level functions that use the default logger

// moduleDebug: logs a debug-level message using the default logger.
//
// @luafunc debug
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Supports three patterns:
//
//	log.debug("msg", "key", "value")              -- string-primitive pairs
//	log.debug("msg", {key = "value"})             -- flatten table
//	log.debug("msg", "context", {nested = "data"}) -- JSON encode table
//
// Example:
//
//	log.debug("Processing item", "item_id", 42, "status", "pending")
func moduleDebug(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logDebugImpl(logger, L, 1)
	}
	return logDebugImpl(GetDefaultLogger(), L, 1)
}

// moduleInfo: logs an info-level message using the default logger.
//
// @luafunc info
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	log.info("Server started", "port", 8080)
func moduleInfo(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logInfoImpl(logger, L, 1)
	}
	return logInfoImpl(GetDefaultLogger(), L, 1)
}

// moduleWarn: logs a warning-level message using the default logger.
//
// @luafunc warn
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	log.warn("Retry attempt failed", "attempt", 3, "max_retries", 5)
func moduleWarn(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logWarnImpl(logger, L, 1)
	}
	return logWarnImpl(GetDefaultLogger(), L, 1)
}

// moduleError: logs an error-level message using the default logger.
//
// @luafunc error
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	log.error("Database connection failed", "error", err_msg, "retry_in", 5)
func moduleError(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logErrorImpl(logger, L, 1)
	}
	return logErrorImpl(GetDefaultLogger(), L, 1)
}

// moduleFatal: logs a fatal-level message using the default logger and exits.
//
// @luafunc fatal
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	log.fatal("Critical system failure", "component", "database")
func moduleFatal(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	if logger, ok := ud.Value.(*log.Logger); ok {
		return logFatalImpl(logger, L, 1)
	}
	return logFatalImpl(GetDefaultLogger(), L, 1)
}

// moduleLogger: returns the default logger object.
//
// @luafunc logger
// @luareturn log.Logger logger The default logger object
//
// Example:
//
//	local logger = log.logger()
//	logger:info("Using logger object")
func moduleLogger(L *lua.LState) int {
	ud := getDefaultLoggerUserData(L)
	L.Push(ud)
	return 1
}

// Logger methods

// loggerDebug: logs a debug-level message.
//
// @luamethod log.Logger debug
// @luaparam self log.Logger The logger object
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	logger:debug("Processing item", "item_id", 42)
func loggerDebug(L *lua.LState) int {
	logger := checkLogger(L, 1)
	return logDebugImpl(logger, L, 2)
}

// loggerInfo: logs an info-level message.
//
// @luamethod log.Logger info
// @luaparam self log.Logger The logger object
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	logger:info("Server started", "port", 8080)
func loggerInfo(L *lua.LState) int {
	logger := checkLogger(L, 1)
	return logInfoImpl(logger, L, 2)
}

// loggerWarn: logs a warning-level message.
//
// @luamethod log.Logger warn
// @luaparam self log.Logger The logger object
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	logger:warn("High memory usage", "memory_used", 8.5)
func loggerWarn(L *lua.LState) int {
	logger := checkLogger(L, 1)
	return logWarnImpl(logger, L, 2)
}

// loggerError: logs an error-level message.
//
// @luamethod log.Logger error
// @luaparam self log.Logger The logger object
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	logger:error("Database connection failed", "error", err_msg)
func loggerError(L *lua.LState) int {
	logger := checkLogger(L, 1)
	return logErrorImpl(logger, L, 2)
}

// loggerFatal: logs a fatal-level message and exits.
//
// @luamethod log.Logger fatal
// @luaparam self log.Logger The logger object
// @luaparam msg string The message to log
// @luaparam ... any Optional fields: key-value pairs, a table, or key-table pairs
//
// Example:
//
//	logger:fatal("Critical failure", "component", "database")
func loggerFatal(L *lua.LState) int {
	logger := checkLogger(L, 1)
	return logFatalImpl(logger, L, 2)
}

// loggerWith: creates a new logger with additional fields.
//
// @luamethod log.Logger with
// @luaparam self log.Logger The logger object
// @luaparam ... any Key-value pairs to add to the logger context
// @luareturn log.Logger logger A new logger with the additional fields
//
// Example:
//
//	local logger = log.logger()
//	local contextLogger = logger:with("request_id", "abc123", "user", "john")
//	contextLogger:info("User logged in")
//	logger:info("Other action")  -- original logger unchanged
func loggerWith(L *lua.LState) int {
	logger := checkLogger(L, 1)
	fields := extractFields(L, 2)

	// Create new logger with additional fields
	newLogger := logger.With(fields...)

	// Wrap in UserData and return
	ud := wrapLogger(L, newLogger)
	L.Push(ud)
	return 1
}

// Implementation functions

// logDebugImpl: implementation of debug logging
func logDebugImpl(logger *log.Logger, L *lua.LState, startIdx int) int {
	msg := L.CheckString(startIdx)
	fields := extractFields(L, startIdx+1)
	logger.Debug(msg, fields...)
	return 0
}

// logInfoImpl: implementation of info logging
func logInfoImpl(logger *log.Logger, L *lua.LState, startIdx int) int {
	msg := L.CheckString(startIdx)
	fields := extractFields(L, startIdx+1)
	logger.Info(msg, fields...)
	return 0
}

// logWarnImpl: implementation of warn logging
func logWarnImpl(logger *log.Logger, L *lua.LState, startIdx int) int {
	msg := L.CheckString(startIdx)
	fields := extractFields(L, startIdx+1)
	logger.Warn(msg, fields...)
	return 0
}

// logErrorImpl: implementation of error logging
func logErrorImpl(logger *log.Logger, L *lua.LState, startIdx int) int {
	msg := L.CheckString(startIdx)
	fields := extractFields(L, startIdx+1)
	logger.Error(msg, fields...)
	return 0
}

// logFatalImpl: implementation of fatal logging
func logFatalImpl(logger *log.Logger, L *lua.LState, startIdx int) int {
	msg := L.CheckString(startIdx)
	fields := extractFields(L, startIdx+1)
	logger.Fatal(msg, fields...)
	return 0
}

// extractFields: extracts key-value pairs from Lua stack starting at the given index.
// Supports three patterns:
// 1. String-primitive pairs: log.info("msg", "key", "value", "key2", 42)
// 2. Single table: log.info("msg", {key = "value", key2 = 42})
// 3. String-table pairs: log.info("msg", "context", {nested = "data"})
func extractFields(L *lua.LState, startIdx int) []interface{} {
	top := L.GetTop()
	if startIdx > top {
		return nil
	}

	fields := make([]interface{}, 0)

	i := startIdx
	for i <= top {
		arg := L.Get(i)

		// Case 1: Single table argument - flatten first-level keys
		if tbl, ok := arg.(*lua.LTable); ok && i == startIdx && top == startIdx {
			// Only one argument and it's a table - flatten it
			tbl.ForEach(func(key lua.LValue, val lua.LValue) {
				if keyStr, ok := key.(lua.LString); ok {
					// Convert value using translator
					var goVal any
					if err := translator.FromLua(L, val, &goVal); err == nil {
						fields = append(fields, string(keyStr), goVal)
					}
				}
			})
			return fields
		}

		// Case 2: String followed by table - JSON encode the table
		if i+1 <= top {
			nextArg := L.Get(i + 1)
			if keyStr, ok := arg.(lua.LString); ok {
				if tbl, ok := nextArg.(*lua.LTable); ok {
					// String-table pair: encode table as JSON
					fields = append(fields, string(keyStr), tableToJSON(L, tbl))
					i += 2 // Skip both arguments
					continue
				}
			}
		}

		// Case 3: String-primitive pairs (default behavior)
		if keyStr, ok := arg.(lua.LString); ok {
			fields = append(fields, string(keyStr))
		} else {
			// Convert non-string arguments using translator
			var goVal any
			if err := translator.FromLua(L, arg, &goVal); err == nil {
				fields = append(fields, goVal)
			}
		}

		i++
	}

	return fields
}

// tableToJSON: converts a Lua table to a JSON string for logging
func tableToJSON(L *lua.LState, tbl *lua.LTable) string {
	// Use the glua translator to convert Lua table to Go
	var result any
	if err := translator.FromLua(L, tbl, &result); err != nil {
		return fmt.Sprintf("{\"error\": \"failed to convert: %v\"}", err)
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"failed to marshal: %v\"}", err)
	}
	return string(jsonBytes)
}
