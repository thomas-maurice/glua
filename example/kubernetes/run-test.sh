#!/bin/bash

set -e

echo "=== Setting up Kind cluster ==="

# Create Kind cluster
kind create cluster --name glua-test --wait 5m

# Export kubeconfig
export KUBECONFIG="$(kind get kubeconfig --name glua-test)"

echo ""
echo "=== Running Kubernetes client test ==="
echo ""

# Build and run the example
go run main.go test.lua

echo ""
echo "=== Cleaning up Kind cluster ==="

# Delete Kind cluster
kind delete cluster --name glua-test

echo ""
echo "=== Test complete ==="
