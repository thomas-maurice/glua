---@meta log

---@class log.Logger
local Logger = {}

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function Logger:debug(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function Logger:info(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function Logger:warn(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function Logger:error(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function Logger:fatal(msg, ...) end

---@param ... any Key-value pairs to add to the logger context
---@return log.Logger log.Logger A new logger with the additional fields
function Logger:with(...) end

---@class log
---@field Logger log.Logger
local log = {}

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function log.debug(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function log.info(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function log.warn(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function log.error(msg, ...) end

---@param msg string The message to log
---@param ... any Optional fields: key-value pairs, a table, or key-table pairs
function log.fatal(msg, ...) end

---@return log.Logger log.Logger The default logger object
function log.logger() end

log.Logger = Logger

return log
