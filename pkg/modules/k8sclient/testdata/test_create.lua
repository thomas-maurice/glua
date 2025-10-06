-- Test: create operation
--
-- Tests creating a new Kubernetes resource.

local cm = {
	apiVersion = "v1",
	kind = "ConfigMap",
	metadata = {
		name = "new-config",
		namespace = "default"
	},
	data = {
		foo = "bar"
	}
}

local created, err = client.create(cm)

if err then
	error("Failed to create: " .. err)
end

if created.metadata.name ~= "new-config" then
	error("Expected name 'new-config', got " .. tostring(created.metadata.name))
end

if created.data.foo ~= "bar" then
	error("Expected data.foo 'bar', got " .. tostring(created.data.foo))
end

return true
