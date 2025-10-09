#!/bin/bash
set -e

echo "=== Testing webhook mutation ==="

# Create test pod with mutation label
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Pod
metadata:
  name: test-mutation
  namespace: default
  labels:
    thomas.maurice/mutate: "true"
spec:
  containers:
  - name: nginx
    image: nginx:alpine
EOF

echo "Waiting for pod to be created..."
sleep 2

echo ""
echo "=== Pod annotations ==="
kubectl get pod test-mutation -n default -o jsonpath='{.metadata.annotations}' | jq '.'

echo ""
echo "Cleaning up test pod..."
kubectl delete pod test-mutation -n default

echo "✓ Test completed"
