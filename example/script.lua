local k8s = require("kubernetes")

print(
	"╔════════════════════════════════════════════════════════════╗"
)
print("║              Lua Script: Pod Processing                    ║")
print(
	"╚════════════════════════════════════════════════════════════╝"
)

-- ============================================================================
-- Helper Functions
-- ============================================================================

-- dump: recursively print tables for debugging
local function dump(o, indent)
	indent = indent or 0
	local prefix = string.rep("  ", indent)

	if type(o) == "table" then
		local s = "{\n"
		for k, v in pairs(o) do
			local key = type(k) == "number" and k or '"' .. k .. '"'
			s = s .. prefix .. "  [" .. key .. "] = " .. dump(v, indent + 1) .. ",\n"
		end
		return s .. prefix .. "}"
	else
		return tostring(o)
	end
end

-- ============================================================================
-- Main Processing
-- ============================================================================

---@type corev1.Pod
local pod = originalPod

-- Verify pod structure
assert(pod, "Pod is nil!")
assert(pod.metadata, "Pod metadata is nil!")
assert(pod.spec, "Pod spec is nil!")
assert(pod.spec.containers, "Pod containers is nil!")
assert(pod.spec.containers[1], "First container is nil!")

print("\n[Lua] Pod Basic Info:")
print("  Kind: " .. pod.kind)
print("  API Version: " .. pod.apiVersion)
print("  Name: " .. pod.metadata.name)
print("  Namespace: " .. pod.metadata.namespace)

-- ============================================================================
-- Feature 1: Parse Kubernetes Timestamp
-- ============================================================================
print("\n[Lua] Parsing creationTimestamp...")
local timestampStr = pod.metadata.creationTimestamp
print("  Input: " .. timestampStr)

local timestamp, err = k8s.parse_time(timestampStr)
if err then
	error("Failed to parse timestamp: " .. err)
end

print(string.format("  Output: %d (Unix timestamp)", timestamp))
assert(timestamp > 0, "Timestamp should be positive")

-- Export for Go verification
parsedTimestamp = timestamp

-- ============================================================================
-- Feature 2: Parse CPU Resources
-- ============================================================================
print("\n[Lua] Parsing CPU resources...")
local container = pod.spec.containers[1]
print("  Container: " .. container.name)

-- Parse CPU limits
local cpuLimitStr = container.resources.limits["cpu"]
print("  CPU limit (raw): " .. cpuLimitStr)

local cpuMillis, err = k8s.parse_cpu(cpuLimitStr)
if err then
	error("Failed to parse CPU: " .. err)
end

print(string.format("  CPU limit (parsed): %d millicores", cpuMillis))
assert(cpuMillis > 0, "CPU should be positive")

-- Parse CPU requests
local cpuRequestStr = container.resources.requests["cpu"]
local cpuReqMillis, err = k8s.parse_cpu(cpuRequestStr)
if err then
	error("Failed to parse CPU request: " .. err)
end

print(string.format("  CPU request (parsed): %d millicores", cpuReqMillis))

-- Export for Go verification
parsedCPUMillis = cpuMillis

-- ============================================================================
-- Feature 3: Parse Memory Resources
-- ============================================================================
print("\n[Lua] Parsing memory resources...")

-- Parse memory limits
local memLimitStr = container.resources.limits["memory"]
print("  Memory limit (raw): " .. memLimitStr)

local memBytes, err = k8s.parse_memory(memLimitStr)
if err then
	error("Failed to parse memory: " .. err)
end

local memMB = memBytes / (1024 * 1024)
print(string.format("  Memory limit (parsed): %d bytes (%.2f MB)", memBytes, memMB))

-- Parse memory requests
local memRequestStr = container.resources.requests["memory"]
local memReqBytes, err = k8s.parse_memory(memRequestStr)
if err then
	error("Failed to parse memory request: " .. err)
end

local memReqMB = memReqBytes / (1024 * 1024)
print(string.format("  Memory request (parsed): %d bytes (%.2f MB)", memReqBytes, memReqMB))

-- Export for Go verification
parsedMemoryBytes = memBytes

-- ============================================================================
-- Feature 4: Iterate Over Arrays (Environment Variables)
-- ============================================================================
print("\n[Lua] Processing environment variables...")
if container.env and #container.env > 0 then
	for i, envVar in ipairs(container.env) do
		print(string.format("  %d. %s=%s", i, envVar.name, envVar.value))
	end
else
	print("  No environment variables")
end

-- ============================================================================
-- Feature 5: Iterate Over Maps (Labels)
-- ============================================================================
print("\n[Lua] Processing labels...")
if pod.metadata.labels then
	for key, value in pairs(pod.metadata.labels) do
		print(string.format("  %s: %s", key, value))
	end
else
	print("  No labels")
end

-- ============================================================================
-- Feature 6: Data Validation
-- ============================================================================
print("\n[Lua] Validating data...")
local validations = {
	{ name = "Pod name not empty", check = pod.metadata.name ~= "" },
	{ name = "Container name not empty", check = container.name ~= "" },
	{ name = "Container image not empty", check = container.image ~= "" },
	{ name = "CPU limit > 0", check = cpuMillis > 0 },
	{ name = "Memory limit > 0", check = memBytes > 0 },
	{ name = "Timestamp > 0", check = timestamp > 0 },
}

local passed = 0
local failed = 0
for _, validation in ipairs(validations) do
	if validation.check then
		print("  ✓ " .. validation.name)
		passed = passed + 1
	else
		print("  ✗ " .. validation.name)
		failed = failed + 1
	end
end

print(string.format("\n  Results: %d passed, %d failed", passed, failed))
if failed > 0 then
	error("Validation failed!")
end

-- ============================================================================
-- Feature 7: Table Modification (Pass-through for round-trip test)
-- ============================================================================
print("\n[Lua] Preparing pod for round-trip conversion...")

-- In a real scenario, you might modify the pod here:
-- pod.metadata.labels["processed-by-lua"] = "true"
-- pod.spec.containers[1].env[#pod.spec.containers[1].env + 1] = {name = "PROCESSED", value = "true"}

-- For this demo, we pass it through unchanged to verify round-trip integrity
modifiedPod = pod
print("  ✓ Pod ready for Go conversion")

-- ============================================================================
-- Summary
-- ============================================================================
print(
	"\n╔════════════════════════════════════════════════════════════╗"
)
print("║         Lua Script Completed Successfully ✓                ║")
print(
	"╚════════════════════════════════════════════════════════════╝"
)
