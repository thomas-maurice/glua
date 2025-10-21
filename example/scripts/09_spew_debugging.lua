
---@diagnostic disable-next-line: undefined-global
-- Script: Spew Debugging
-- Description: Demonstrates the spew module for deep inspection and debugging of Lua tables

local spew = require("spew")

print("=== Spew Module Demo ===\n")

-- 1. Basic spew usage
print("[1] Basic spew.dump() - prints to stdout:")
local simple = {
    service = "my-app",
    version = "1.0.0",
    port = 8080
}
spew.dump(simple)

-- 2. Using sdump to capture output
print("\n[2] Using spew.sdump() - returns string:")
local data = {
    status = "running",
    uptime = 3600,
    features = {"auth", "api", "metrics"}
}

local dumpStr = spew.sdump(data)
print("Captured spew output:")
print(dumpStr)

-- 3. Inspecting pod structure
print("\n[3] Inspecting pod metadata structure:")
local podMeta = {
    name = myPod.metadata.name,
    namespace = myPod.metadata.namespace,
    labels = myPod.metadata.labels,
    creationTimestamp = myPod.metadata.creationTimestamp
}

spew.dump(podMeta)

-- 4. Deep inspection of container resources
print("\n[4] Deep inspection of container resources:")
for i, container in ipairs(myPod.spec.containers) do
    print(string.format("\nContainer #%d:", i))
    local containerInfo = {
        name = container.name,
        image = container.image,
        resources = container.resources
    }
    spew.dump(containerInfo)
end

-- 5. Comparing structures
print("\n[5] Using spew to compare data structures:")
local config1 = {
    timeout = 30,
    retries = 3,
    endpoints = {"/health", "/metrics"}
}

local config2 = {
    timeout = 60,  -- Different!
    retries = 3,
    endpoints = {"/health", "/metrics"}
}

print("Config 1:")
local dump1 = spew.sdump(config1)
print(dump1)

print("\nConfig 2:")
local dump2 = spew.sdump(config2)
print(dump2)

if dump1 ~= dump2 then
    print("✓ Configs are different (timeout value changed)")
end

-- 6. Debugging complex nested structures
print("\n[6] Debugging deeply nested pod structure:")
local debugInfo = {
    kind = myPod.kind,
    metadata = {
        name = myPod.metadata.name,
        namespace = myPod.metadata.namespace,
        labels = myPod.metadata.labels
    },
    spec = {
        containerCount = #myPod.spec.containers,
        firstContainer = {
            name = myPod.spec.containers[1].name,
            image = myPod.spec.containers[1].image,
            hasEnv = myPod.spec.containers[1].env ~= nil,
            envCount = myPod.spec.containers[1].env and #myPod.spec.containers[1].env or 0
        }
    }
}

print("Debug dump of pod structure:")
spew.dump(debugInfo)

-- 7. Using spew for error debugging
print("\n[7] Using spew for error context:")
local function validateConfig(config)
    if not config.required_field then
        local errorContext = {
            error = "missing required_field",
            received = config,
            timestamp = os.time()
        }
        print("ERROR CONTEXT:")
        spew.dump(errorContext)
        return false
    end
    return true
end

-- Test with invalid config
local invalidConfig = {some_field = "value"}
validateConfig(invalidConfig)

-- 8. Array vs Map detection
print("\n[8] Spew array vs map detection:")

print("\nArray (consecutive integer keys):")
local array = {10, 20, 30, 40}
spew.dump(array)

print("\nMap (string keys):")
local map = {
    first = 10,
    second = 20,
    third = 30
}
spew.dump(map)

print("\nSparse array (treated as map):")
local sparse = {}
sparse[1] = "first"
sparse[5] = "fifth"
sparse[10] = "tenth"
spew.dump(sparse)

print("\n=== Spew Debugging Complete ===")
print("Tip: Use spew.dump() for quick stdout debugging")
print("     Use spew.sdump() when you need to log or manipulate the output")
