
-- Example 5: Timestamp Operations
-- Parse and manipulate Kubernetes timestamps

local k8s = require("kubernetes")

---@type corev1.Pod
local pod = myPod

print("=== Timestamp Operations ===")
print("Pod: " .. pod.metadata.name)
print()

-- Parse creation timestamp
local creationTime = pod.metadata.creationTimestamp
print("Creation Timestamp (string): " .. creationTime)

local timestamp, err = k8s.parse_time(creationTime)
if err then
    error("Failed to parse timestamp: " .. err)
end

print("Creation Timestamp (Unix): " .. timestamp)
print()

-- Calculate age
local currentTime = os.time()
local ageSeconds = currentTime - timestamp
local ageMinutes = ageSeconds / 60
local ageHours = ageMinutes / 60
local ageDays = ageHours / 24

print("Pod Age:")
print(string.format("  %d seconds", ageSeconds))
print(string.format("  %.2f minutes", ageMinutes))
print(string.format("  %.2f hours", ageHours))
print(string.format("  %.2f days", ageDays))
print()

-- Format timestamp back
local formatted, err = k8s.format_time(timestamp)
if err then
    error("Failed to format timestamp: " .. err)
end

print("Formatted Timestamp: " .. formatted)
print()

-- Create custom timestamp
local oneDayAgo = currentTime - (24 * 60 * 60)
local oneDayAgoFormatted, err = k8s.format_time(oneDayAgo)
if not err then
    print("24 hours ago: " .. oneDayAgoFormatted)
end

-- Validate age (example: warn if pod is older than 7 days)
if ageDays > 7 then
    print("\n⚠️  WARNING: Pod is older than 7 days (" .. string.format("%.1f", ageDays) .. " days)")
else
    print("\n✓ Pod age is acceptable")
end
