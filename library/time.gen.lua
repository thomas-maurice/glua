---@meta time

---@class time
local time = {}

---@return timestamp number Current Unix timestamp (seconds since epoch)
function time.now() end

---@param timestr string The time string to parse
---@param layout string The Go time layout format (e.g., "2006-01-02 15:04:05")
---@return timestamp number Unix timestamp, or nil on error
---@return err string|nil Error message if parsing failed
function time.parse(timestr, layout) end

---@param timestr string The RFC3339 time string (e.g., "2024-03-15T14:30:00Z")
---@return timestamp number Unix timestamp, or nil on error
---@return err string|nil Error message if parsing failed
function time.parse_rfc3339(timestr) end

---@param timestamp number Unix timestamp
---@param layout string The Go time layout format (e.g., "2006-01-02 15:04:05")
---@return formatted string Formatted time string
function time.format(timestamp, layout) end

---@param timestamp number Unix timestamp
---@param seconds number Number of seconds to add (can be negative)
---@return new_timestamp number New Unix timestamp
function time.add(timestamp, seconds) end

---@param time1 number First Unix timestamp
---@param time2 number Second Unix timestamp
---@return seconds number Difference in seconds (time1 - time2)
function time.diff(time1, time2) end

---@param seconds number Number of seconds to sleep
function time.sleep(seconds) end

---@param timestamp number Unix timestamp
---@return date_table table Table with year, month, day, hour, min, sec, wday, yday, isdst
function time.to_osdate(timestamp) end

---@param date_table table Table with year, month, day, hour, min, sec (other fields optional)
---@return timestamp number Unix timestamp
function time.from_osdate(date_table) end

return time
