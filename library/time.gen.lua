---@meta time

---@class time
local time = {}

--- returns the current Unix timestamp
---@return number timestamp seconds since the Unix epoch (UTC)
function time.now() end

--- parses a time string with a Go layout, raises on error
---@param timestr string the timestamp string to parse
---@param layout string a Go reference-time layout, e.g. "2006-01-02 15:04:05"
---@return number timestamp seconds since the Unix epoch (UTC)
function time.parse(timestr, layout) end

--- parses an RFC3339 time string, raises on error
---@param timestr string an RFC3339-formatted timestamp, e.g. "2025-10-03T16:39:00Z"
---@return number timestamp seconds since the Unix epoch (UTC)
function time.parse_rfc3339(timestr) end

--- formats a Unix timestamp with a Go layout
---@param timestamp number seconds since the Unix epoch (UTC)
---@param layout string a Go reference-time layout, e.g. "2006-01-02 15:04:05"
---@return string out timestamp formatted per layout, in UTC
function time.format(timestamp, layout) end

--- adds seconds to a Unix timestamp
---@param timestamp number seconds since the Unix epoch (UTC)
---@param seconds number the number of seconds to add; negative subtracts
---@return number timestamp timestamp + seconds
function time.add(timestamp, seconds) end

--- returns the difference in seconds between two timestamps (t1 - t2)
---@param t1 number the minuend timestamp, seconds since the Unix epoch
---@param t2 number the subtrahend timestamp, seconds since the Unix epoch
---@return number seconds t1 - t2, in seconds
function time.diff(t1, t2) end

--- pauses execution for the given number of seconds
---@param seconds number how long to sleep; fractional values are supported
function time.sleep(seconds) end

--- converts a Unix timestamp to an os.date-compatible table
---@param timestamp number seconds since the Unix epoch (UTC)
---@return table date table with year, month, day, hour, min, sec, wday, yday, isdst fields, in UTC
function time.to_osdate(timestamp) end

--- converts an os.date-compatible table to a Unix timestamp, raises on invalid input
---@param date_table table table with required year, month, day and optional hour, min, sec (default 0)
---@return number timestamp seconds since the Unix epoch, interpreting the fields as UTC
function time.from_osdate(date_table) end

return time
