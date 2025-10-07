#!/bin/bash
# run-test.sh: Run k8sclient example with Kind cluster

set -e

CLUSTER_NAME="glua-k8sclient-test"

echo "=== Setting up Kind cluster for k8sclient example ==="

# Check if Kind is installed
if ! command -v kind &> /dev/null; then
    echo "Error: kind is not installed. Please install it first:"
    echo "  https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    echo "Error: kubectl is not installed. Please install it first:"
    echo "  https://kubernetes.io/docs/tasks/tools/"
    exit 1
fi

# Cleanup function
cleanup() {
    echo ""
    echo "=== Cleaning up ==="
    kind delete cluster --name "$CLUSTER_NAME" || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Delete existing cluster if it exists
if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "Deleting existing cluster: $CLUSTER_NAME"
    kind delete cluster --name "$CLUSTER_NAME"
fi

# Create Kind cluster
echo "Creating Kind cluster: $CLUSTER_NAME"
kind create cluster --name "$CLUSTER_NAME" --wait 5m

# Wait for cluster to be ready
echo "Waiting for cluster to be ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=300s

# Build the example program
echo ""
echo "=== Building example program ==="
go build -o ../../bin/k8sclient-example .

# Run the test script
echo ""
echo "=== Running k8sclient test script ==="
../../bin/k8sclient-example -script test.lua

echo ""
echo "=== Test completed successfully! ==="
