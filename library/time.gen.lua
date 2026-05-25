---@meta time

---@class time
local time = {}

--- returns the current Unix timestamp
---@return number
function time.now() end

--- parses a time string with a Go layout, raises on error
---@param timestr string
---@param layout string
---@return number
function time.parse(timestr, layout) end

--- parses an RFC3339 time string, raises on error
---@param timestr string
---@return number
function time.parse_rfc3339(timestr) end

--- formats a Unix timestamp with a Go layout
---@param timestamp number
---@param layout string
---@return string
function time.format(timestamp, layout) end

--- adds seconds to a Unix timestamp
---@param timestamp number
---@param seconds number
---@return number
function time.add(timestamp, seconds) end

--- returns the difference in seconds between two timestamps (t1 - t2)
---@param t1 number
---@param t2 number
---@return number
function time.diff(t1, t2) end

--- pauses execution for the given number of seconds
---@param seconds number
function time.sleep(seconds) end

--- converts a Unix timestamp to an os.date-compatible table
---@param timestamp number
---@return gopher-lua.LTable
function time.to_osdate(timestamp) end

--- converts an os.date-compatible table to a Unix timestamp, raises on invalid input
---@param date_table gopher-lua.LTable
---@return number
function time.from_osdate(date_table) end

return time
