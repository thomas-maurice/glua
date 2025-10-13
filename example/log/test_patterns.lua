local log = require("log")

print("=== Logging Pattern Examples ===\n")

-- Pattern 1: String-primitive pairs (traditional)
print("1. String-primitive pairs:")
print("   log.info(\"msg\", \"key\", \"value\", \"key2\", 42)")
log.info("User action", "user_id", "12345", "action", "login", "ip", "192.168.1.1")

-- Pattern 2: Single table (flatten first level)
print("\n2. Single table - flattens first-level keys:")
print("   log.info(\"msg\", {key = \"value\", key2 = 42})")
log.info("HTTP Request", {
    method = "GET",
    path = "/api/users",
    status = 200,
    duration_ms = 45,
    user_agent = "curl/7.68.0"
})

-- Pattern 3: String-table pairs (JSON encode nested data)
print("\n3. String-table pair - JSON encodes the table:")
print("   log.info(\"msg\", \"context\", {nested = {data = \"here\"}})")
log.info("Complex operation", "metadata", {
    user = {
        id = "12345",
        name = "alice",
        roles = {"admin", "developer", "user"}
    },
    request = {
        method = "POST",
        headers = {
            ["Content-Type"] = "application/json",
            ["Authorization"] = "Bearer xxx"
        },
        body = {
            action = "create",
            resource = "document",
            data = {
                title = "New Document",
                tags = {"important", "urgent"}
            }
        }
    },
    metrics = {
        processing_time_ms = 123,
        db_queries = 5,
        cache_hits = 12
    }
})

-- Pattern 4: Mixed usage
print("\n4. Mixed - string-primitives + string-table:")
print("   log.info(\"msg\", \"simple\", \"value\", \"complex\", {...})")
log.info("Mixed example",
    "request_id", "abc-123",
    "status", "success",
    "error_context", {
        code = 500,
        message = "Internal server error",
        stack_trace = {"frame1", "frame2", "frame3"}
    },
    "duration_ms", 250
)

-- Using with logger object
print("\n5. With logger object (same patterns work):")
local logger = log.logger():with("service", "api", "version", "1.0.0")

logger:info("Service started", {port = 8080, workers = 4})
logger:info("Database connected", "db", {
    host = "localhost",
    port = 5432,
    database = "myapp",
    pool_size = 10
})

print("\n=== All patterns demonstrated ===")
