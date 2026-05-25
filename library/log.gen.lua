---@meta log

---@class log.Logger
local Logger = {}

--- log at debug level
---@param msg string
function Logger:debug(msg) end

--- log at info level
---@param msg string
function Logger:info(msg) end

--- log at warn level
---@param msg string
function Logger:warn(msg) end

--- log at error level
---@param msg string
function Logger:error(msg) end

--- log at fatal level
---@param msg string
function Logger:fatal(msg) end

--- return a child logger with extra fields
---@return log.Logger
function Logger:with() end

---@class log
---@field Logger log.Logger
local log = {}

--- log on the default logger at debug level
---@param msg string
function log.debug(msg) end

--- log on the default logger at info level
---@param msg string
function log.info(msg) end

--- log on the default logger at warn level
---@param msg string
function log.warn(msg) end

--- log on the default logger at error level
---@param msg string
function log.error(msg) end

--- log on the default logger at fatal level
---@param msg string
function log.fatal(msg) end

--- return the default logger
---@return log.Logger
function log.logger() end

log.Logger = Logger

return log
