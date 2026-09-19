---@meta log

---@class log.Logger
local Logger = {}

--- log at debug level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function Logger:debug(msg) end

--- log at info level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function Logger:info(msg) end

--- log at warn level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function Logger:warn(msg) end

--- log at error level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function Logger:error(msg) end

--- log at fatal level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields. Terminates the process after logging
function Logger:fatal(msg) end

--- return a child logger with extra fields
---@return log.Logger logger a new Logger that always includes the given key-value fields
function Logger:with() end

---@class log
---@field Logger log.Logger
local log = {}

--- log on the default logger at debug level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function log.debug(msg) end

--- log on the default logger at info level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function log.info(msg) end

--- log on the default logger at warn level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function log.warn(msg) end

--- log on the default logger at error level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields
function log.error(msg) end

--- log on the default logger at fatal level
---@param msg string the message to log; optional trailing key-value pairs or a single table add structured fields. Terminates the process after logging
function log.fatal(msg) end

--- return the default logger
---@return log.Logger logger the active Logger: the one injected via InjectLogger, or the package default
function log.logger() end

log.Logger = Logger

return log
