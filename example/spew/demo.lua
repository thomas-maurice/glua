local spew = require("spew")

print("=== Spew Module - Colored JSON Output Demo ===\n")

-- Simple object
print("1. Simple object:")
spew.dump({
    name = "Alice",
    age = 30,
    active = true
})

-- Nested structure
print("\n2. Nested structure:")
spew.dump({
    user = {
        id = "12345",
        name = "Bob",
        email = "bob@example.com",
        preferences = {
            theme = "dark",
            notifications = true
        }
    },
    metadata = {
        created_at = "2025-10-13T10:00:00Z",
        updated_at = "2025-10-13T11:30:00Z"
    }
})

-- Array
print("\n3. Array:")
spew.dump({
    "apple",
    "banana",
    "cherry",
    "date"
})

-- Array of objects
print("\n4. Array of objects:")
spew.dump({
    {name = "Product 1", price = 19.99, in_stock = true},
    {name = "Product 2", price = 29.99, in_stock = false},
    {name = "Product 3", price = 39.99, in_stock = true}
})

-- Complex nested structure
print("\n5. Complex nested structure:")
spew.dump({
    api_response = {
        status = "success",
        data = {
            users = {
                {id = 1, name = "Alice", roles = {"admin", "user"}},
                {id = 2, name = "Bob", roles = {"user"}},
                {id = 3, name = "Charlie", roles = {"moderator", "user"}}
            },
            pagination = {
                page = 1,
                per_page = 10,
                total = 3
            }
        },
        meta = {
            request_id = "abc-123-xyz",
            timestamp = 1697203200,
            version = "1.0.0"
        }
    }
})

-- Using sdump to get string
print("\n6. Using sdump (returns string):")
local json_str = spew.sdump({
    message = "This is a string",
    number = 42,
    boolean = true,
    nested = {
        key = "value"
    }
})
print("Returned string:")
print(json_str)

print("\n=== Demo Complete ===")
