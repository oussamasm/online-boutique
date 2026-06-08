#!/bin/bash

# Setup script for test-auto namespace
# Usage: ./setup-test-namespace.sh

set -e

NAMESPACE="test-auto"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================="
echo "Setting up test-auto namespace"
echo "=========================================="

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Please install kubectl."
    exit 1
fi

# Check cluster connectivity
echo "[*] Checking cluster connectivity..."
if ! kubectl cluster-info &> /dev/null; then
    echo "❌ Cannot connect to cluster. Please configure kubectl."
    exit 1
fi
echo "✓ Connected to cluster"

# Create namespace and RBAC
echo "[*] Creating namespace: $NAMESPACE..."
kubectl apply -f "$SCRIPT_DIR/setup-test-namespace.yaml"

if [ $? -eq 0 ]; then
    echo "✓ Namespace created successfully"
else
    echo "❌ Failed to create namespace"
    exit 1
fi

# Verify namespace exists
echo "[*] Verifying namespace..."
if kubectl get namespace "$NAMESPACE" &> /dev/null; then
    echo "✓ Namespace verified"
else
    echo "❌ Namespace verification failed"
    exit 1
fi

# Verify ServiceAccount
echo "[*] Verifying ServiceAccount..."
if kubectl get serviceaccount e2e-tests -n "$NAMESPACE" &> /dev/null; then
    echo "✓ ServiceAccount verified"
else
    echo "❌ ServiceAccount verification failed"
    exit 1
fi

# Verify RBAC bindings
echo "[*] Verifying RBAC bindings..."
if kubectl get clusterrolebinding e2e-tests &> /dev/null; then
    echo "✓ ClusterRoleBinding verified"
else
    echo "❌ ClusterRoleBinding verification failed"
    exit 1
fi

if kubectl get rolebinding e2e-tests-namespace -n "$NAMESPACE" &> /dev/null; then
    echo "✓ RoleBinding verified"
else
    echo "❌ RoleBinding verification failed"
    exit 1
fi

echo ""
echo "=========================================="
echo "✓ Setup completed successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Deploy microservices to the test-auto namespace:"
echo "   kubectl apply -f kubernetes-manifests.yaml -n test-auto"
echo ""
echo "2. Run E2E tests:"
echo "   TEST_NAMESPACE=test-auto go test -v ./..."
echo ""
echo "3. Run with Makefile:"
echo "   make test NAMESPACE=test-auto"
echo ""
echo "4. Or simply use default (test-auto):"
echo "   make test"
echo ""
echo "View namespace resources:"
echo "  kubectl get all -n test-auto"
echo "  kubectl describe namespace test-auto"
echo ""
