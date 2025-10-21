---@meta time

---@class time
local time = {}

---@return number timestamp Current Unix timestamp (seconds since epoch)
function time.now() end

---@param string timestr The time string to parse
---@param string layout The Go time layout format (e.g., "2006-01-02 15:04:05")
---@return number timestamp Unix timestamp, or nil on error
---@return string|nil err Error message if parsing failed
function time.parse(timestr, layout) end

---@param string timestr The RFC3339 time string (e.g., "2024-03-15T14:30:00Z")
---@return number timestamp Unix timestamp, or nil on error
---@return string|nil err Error message if parsing failed
function time.parse_rfc3339(timestr) end

---@param number timestamp Unix timestamp
---@param string layout The Go time layout format (e.g., "2006-01-02 15:04:05")
---@return string formatted Formatted time string
function time.format(timestamp, layout) end

---@param number timestamp Unix timestamp
---@param number seconds Number of seconds to add (can be negative)
---@return number new_timestamp New Unix timestamp
function time.add(timestamp, seconds) end

---@param number time1 First Unix timestamp
---@param number time2 Second Unix timestamp
---@return number seconds Difference in seconds (time1 - time2)
function time.diff(time1, time2) end

---@param number seconds Number of seconds to sleep
function time.sleep(seconds) end

---@param number timestamp Unix timestamp
---@return table date_table Table with year, month, day, hour, min, sec, wday, yday, isdst
function time.to_osdate(timestamp) end

---@param table date_table Table with year, month, day, hour, min, sec (other fields optional)
---@return number Unix timestamp
function time.from_osdate(date_table) end

return time
