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
	"github.com/thomas-maurice/glua/pkg/luareg"
	lua "github.com/yuin/gopher-lua"
)

// loggerInjectedKey: Lua global name used to store an injected *log.Logger.
const loggerInjectedKey = "__log_injected_logger__"

// loggerClassName: Lua type name for the Logger class.
const loggerClassName = "log.Logger"

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

// InjectLogger: injects a pre-configured logger instance into the Lua state.
// This is OPTIONAL - if you don't inject a logger, a default one will be created automatically.
// Use this only when you want to add pre-set fields (like request ID, user ID, etc.) from Go.
//
// Example usage:
//
//	// Get default logger and add fields
//	logger := logmodule.GetDefaultLogger().With("request_id", "abc123", "user", "john")
//	logmodule.InjectLogger(L, logger)
func InjectLogger(L *lua.LState, logger *log.Logger) {
	// Store a bare UserData. If the class metatable is already registered on L
	// (i.e., Loader has been called), attach it now. Otherwise leave it bare —
	// Loader will re-wrap after registering the metatable.
	ud := L.NewUserData()
	ud.Value = logger
	mt := L.GetTypeMetatable(loggerClassName)
	if mt != lua.LNil {
		L.SetMetatable(ud, mt)
	}
	L.SetGlobal(loggerInjectedKey, ud)
}

// newLoggerClass: creates a fresh *luareg.Class[*log.Logger] with all methods registered.
// Called once per build() invocation so each Module gets its own Class instance.
func newLoggerClass() *luareg.Class[*log.Logger] {
	cls := luareg.NewClass[*log.Logger](loggerClassName, "structured logger object")
	// Methods use (receiver, *lua.LState, msg string) — the LState escape hatch
	// allows reading extra variadic key-value fields beyond msg, while keeping
	// msg visible to reflection for stub generation.
	cls.Method("debug", loggerDebugMethod, "log at debug level",
		luareg.Args("msg"))
	cls.Method("info", loggerInfoMethod, "log at info level",
		luareg.Args("msg"))
	cls.Method("warn", loggerWarnMethod, "log at warn level",
		luareg.Args("msg"))
	cls.Method("error", loggerErrorMethod, "log at error level",
		luareg.Args("msg"))
	cls.Method("fatal", loggerFatalMethod, "log at fatal level",
		luareg.Args("msg"))
	// with is fully variadic (key-value pairs only, no fixed params) — keep the
	// *lua.LState-only escape hatch; stub shows Logger:with() with no params,
	// which is the acceptable degradation documented below.
	cls.Method("with", loggerWithMethod, "return a child logger with extra fields")
	return cls
}

// getActiveLogger: returns the active logger for the given Lua state.
// Looks for an injected logger in globals; falls back to the package default.
func getActiveLogger(L *lua.LState) *log.Logger {
	lv := L.GetGlobal(loggerInjectedKey)
	if ud, ok := lv.(*lua.LUserData); ok {
		if logger, ok := ud.Value.(*log.Logger); ok {
			return logger
		}
	}
	return GetDefaultLogger()
}

// wrapLogger: wraps a *log.Logger as a Lua UserData with the class metatable.
// Requires the class metatable to be already registered on L (i.e., PushTo has run).
func wrapLogger(L *lua.LState, logger *log.Logger) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = logger
	L.SetMetatable(ud, L.GetTypeMetatable(loggerClassName))
	return ud
}

// Logger method implementations.
// debug/info/warn/error/fatal use (receiver, *lua.LState, msg string): the
// framework reads msg from stack position 2; optional key-value fields at
// positions 3..N are read via extractFields using the LState escape hatch.
// with uses (receiver, *lua.LState) only — fully variadic, no fixed params.

// loggerDebugMethod: logs a debug-level message.
func loggerDebugMethod(l *log.Logger, L *lua.LState, msg string) {
	fields := extractFields(L, 3)
	l.Debug(msg, fields...)
}

// loggerInfoMethod: logs an info-level message.
func loggerInfoMethod(l *log.Logger, L *lua.LState, msg string) {
	fields := extractFields(L, 3)
	l.Info(msg, fields...)
}

// loggerWarnMethod: logs a warn-level message.
func loggerWarnMethod(l *log.Logger, L *lua.LState, msg string) {
	fields := extractFields(L, 3)
	l.Warn(msg, fields...)
}

// loggerErrorMethod: logs an error-level message.
func loggerErrorMethod(l *log.Logger, L *lua.LState, msg string) {
	fields := extractFields(L, 3)
	l.Error(msg, fields...)
}

// loggerFatalMethod: logs a fatal-level message and exits.
func loggerFatalMethod(l *log.Logger, L *lua.LState, msg string) {
	fields := extractFields(L, 3)
	l.Fatal(msg, fields...)
}

// loggerWithMethod: creates a child logger with extra fields.
// Fully variadic: stack positions 2..N are key-value pairs.
// Uses the *lua.LState-only escape hatch so no fixed params are enforced.
func loggerWithMethod(l *log.Logger, L *lua.LState) *log.Logger {
	fields := extractFields(L, 2)
	return l.With(fields...)
}

// Module-level shorthand functions use (L *lua.LState, msg string): the
// framework reads msg from stack position 1; extra key-value fields at
// positions 2..N are extracted via the LState escape hatch.

// moduleLuaDebug: logs at debug level on the active logger.
func moduleLuaDebug(L *lua.LState, msg string) {
	logger := getActiveLogger(L)
	fields := extractFields(L, 2)
	logger.Debug(msg, fields...)
}

// moduleLuaInfo: logs at info level on the active logger.
func moduleLuaInfo(L *lua.LState, msg string) {
	logger := getActiveLogger(L)
	fields := extractFields(L, 2)
	logger.Info(msg, fields...)
}

// moduleLuaWarn: logs at warn level on the active logger.
func moduleLuaWarn(L *lua.LState, msg string) {
	logger := getActiveLogger(L)
	fields := extractFields(L, 2)
	logger.Warn(msg, fields...)
}

// moduleLuaError: logs at error level on the active logger.
func moduleLuaError(L *lua.LState, msg string) {
	logger := getActiveLogger(L)
	fields := extractFields(L, 2)
	logger.Error(msg, fields...)
}

// moduleLuaFatal: logs at fatal level on the active logger.
func moduleLuaFatal(L *lua.LState, msg string) {
	logger := getActiveLogger(L)
	fields := extractFields(L, 2)
	logger.Fatal(msg, fields...)
}

// moduleLuaLogger: returns the active logger as a UserData object.
func moduleLuaLogger(L *lua.LState) *log.Logger {
	return getActiveLogger(L)
}

// extractFields: extracts key-value pairs from the Lua stack starting at startIdx.
// Supports three patterns:
//  1. String-primitive pairs: log.info("msg", "key", "value", "key2", 42)
//  2. Single table: log.info("msg", {key = "value"})
//  3. String-table pairs: log.info("msg", "context", {nested = "data"})
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
			tbl.ForEach(func(key lua.LValue, val lua.LValue) {
				if keyStr, ok := key.(lua.LString); ok {
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
					fields = append(fields, string(keyStr), tableToJSON(L, tbl))
					i += 2
					continue
				}
			}
		}

		// Case 3: String-primitive pairs (default)
		if keyStr, ok := arg.(lua.LString); ok {
			fields = append(fields, string(keyStr))
		} else {
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
	var result any
	if err := translator.FromLua(L, tbl, &result); err != nil {
		return fmt.Sprintf("{\"error\": \"failed to convert: %v\"}", err)
	}
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("{\"error\": \"failed to marshal: %v\"}", err)
	}
	return string(jsonBytes)
}

// build: constructs the module definition. Reused by Loader and Register.
// Each call creates a fresh Module and Class to avoid duplicate registration panics.
func build() *luareg.Module {
	m := luareg.NewModule("log", "structured logger")
	cls := newLoggerClass()
	m.RegisterClass(cls)
	// Module-level shorthand functions use (L *lua.LState, msg string): msg is
	// visible to reflection for stub generation, and L is the escape hatch for
	// reading extra variadic key-value fields beyond msg.
	m.Fn("debug", moduleLuaDebug, "log on the default logger at debug level",
		luareg.Args("msg"))
	m.Fn("info", moduleLuaInfo, "log on the default logger at info level",
		luareg.Args("msg"))
	m.Fn("warn", moduleLuaWarn, "log on the default logger at warn level",
		luareg.Args("msg"))
	m.Fn("error", moduleLuaError, "log on the default logger at error level",
		luareg.Args("msg"))
	m.Fn("fatal", moduleLuaFatal, "log on the default logger at fatal level",
		luareg.Args("msg"))
	// logger returns the active *log.Logger; return type is auto-wrapped by luareg.
	m.Fn("logger", moduleLuaLogger, "return the default logger")
	return m
}

// Loader: gopher-lua module loader. Use with L.PreloadModule("log", log.Loader).
func Loader(L *lua.LState) int {
	n := build().PushTo(L)

	// After PushTo the class metatable for "log.Logger" is registered on L.
	// Re-wrap any previously injected logger so it gets the proper metatable.
	if lv := L.GetGlobal(loggerInjectedKey); lv != lua.LNil {
		if ud, ok := lv.(*lua.LUserData); ok {
			if logger, ok2 := ud.Value.(*log.Logger); ok2 {
				wrapped := wrapLogger(L, logger)
				L.SetGlobal(loggerInjectedKey, wrapped)
			}
		}
	}

	return n
}

// Register: adds this module to reg for stub generation.
func Register(reg *luareg.Registry) {
	build().Register(reg)
}
