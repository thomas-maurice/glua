-- Copyright (c) 2024-2025 Thomas Maurice
--
-- Permission is hereby granted, free of charge, to any person obtaining a copy
-- of this software and associated documentation files (the "Software"), to deal
-- in the Software without restriction, including without limitation the rights
-- to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
-- copies of the Software, and to permit persons to whom the Software is
-- furnished to do so, subject to the following conditions:
--
-- The above copyright notice and this permission notice shall be included in all
-- copies or substantial portions of the Software.
--
-- THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
-- IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
-- FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
-- AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
-- LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
-- OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
-- SOFTWARE.

-- Test: integration test for match_gvk function
--
-- Verifies match_gvk works correctly with multiple resource types.

local k8s = require("kubernetes")

-- Test Pod (core/v1)
local pod = {
	apiVersion = "v1",
	kind = "Pod",
}
local pod_matcher = {group = "", version = "v1", kind = "Pod"}
if not k8s.match_gvk(pod, pod_matcher) then
	error("Pod should match v1/Pod")
end

-- Test Deployment (apps/v1)
local deployment = {
	apiVersion = "apps/v1",
	kind = "Deployment",
}
local deployment_matcher = {group = "apps", version = "v1", kind = "Deployment"}
if not k8s.match_gvk(deployment, deployment_matcher) then
	error("Deployment should match apps/v1/Deployment")
end

-- Test Service (core/v1)
local service = {
	apiVersion = "v1",
	kind = "Service",
}
local service_matcher = {group = "", version = "v1", kind = "Service"}
if not k8s.match_gvk(service, service_matcher) then
	error("Service should match v1/Service")
end

-- Test StatefulSet (apps/v1)
local statefulset = {
	apiVersion = "apps/v1",
	kind = "StatefulSet",
}
local statefulset_matcher = {group = "apps", version = "v1", kind = "StatefulSet"}
if not k8s.match_gvk(statefulset, statefulset_matcher) then
	error("StatefulSet should match apps/v1/StatefulSet")
end

-- Test wrong kind
local service_matcher_wrong = {group = "", version = "v1", kind = "Service"}
if k8s.match_gvk(pod, service_matcher_wrong) then
	error("Pod should NOT match v1/Service")
end

-- Test wrong version
local deployment_matcher_wrong = {group = "apps", version = "v1beta1", kind = "Deployment"}
if k8s.match_gvk(deployment, deployment_matcher_wrong) then
	error("Deployment should NOT match apps/v1beta1/Deployment")
end

return true
