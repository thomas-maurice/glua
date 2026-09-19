-- End-to-end test over a realistic Pod-shaped table built in pure Lua. This
-- is the test that proves the Translator round trip actually works end to
-- end (Lua table -> Go interface{} -> jsonpath reflection walk -> Go
-- interface{} -> Lua value), which is this module's only real risk.
local jsonpath = require("jsonpath")

local pod = {
  apiVersion = "v1",
  kind = "Pod",
  metadata = {
    name = "web-1",
    namespace = "default",
    labels = {app = "web"},
  },
  spec = {
    containers = {
      {name = "web", image = "nginx:1.25", ports = {{containerPort = 80}}},
      {name = "sidecar", image = "envoy:1.28", ports = {{containerPort = 9901}}},
    },
  },
  status = {
    phase = "Running",
  },
}

local images = jsonpath.query(pod, "{.spec.containers[*].image}")
assert(#images == 2, "should find both container images")
assert(images[1] == "nginx:1.25" and images[2] == "envoy:1.28", "images should be in container order")

local name = jsonpath.first(pod, "{.metadata.name}", nil)
assert(name == "web-1", "should resolve the pod name")

assert(jsonpath.exists(pod, "{.spec.containers[?(@.name==\"sidecar\")]}") == true,
  "exists should find a container by name via a filter")
assert(jsonpath.exists(pod, "{.spec.containers[?(@.name==\"missing\")]}") == false,
  "exists should not find a container that is not present")

local summary = jsonpath.render(pod, "{.metadata.name} ({.status.phase})")
assert(summary == "web-1 (Running)", "render should compose a human-readable summary")

return true
