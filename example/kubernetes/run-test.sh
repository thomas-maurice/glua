#!/bin/bash

set -e

CLUSTER_NAME="glua-test"
KUBECONFIG_FILE="/tmp/glua-test-kubeconfig-$$"

# Cleanup function
cleanup() {
	echo ""
	echo "=== Cleaning up ==="

	# Delete Kind cluster
	if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
		kind delete cluster --name "$CLUSTER_NAME"
	fi

	# Remove temporary kubeconfig
	if [ -f "$KUBECONFIG_FILE" ]; then
		rm -f "$KUBECONFIG_FILE"
	fi
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo "=== Setting up Kind cluster ==="

# Create Kind cluster
kind create cluster --name "$CLUSTER_NAME" --wait 5m

# Write kubeconfig to temporary file
kind get kubeconfig --name "$CLUSTER_NAME" > "$KUBECONFIG_FILE"
export KUBECONFIG="$KUBECONFIG_FILE"

echo ""
echo "=== Running Kubernetes client test ==="
echo ""

# Build and run the example
go run main.go test.lua

echo ""
echo "=== Test complete ==="
