#!/bin/bash

# E2E Test Suite Automation Script for RKE2 Kubernetes
# This script automates all testing requirements:
# 1. Cluster accessibility and node status
# 2. Pod state verification
# 3. Liveness and Readiness probe testing
# 4. Multi-namespace pod comparison
# 5. Pod cleanup

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
TEST_NAMESPACE="${TEST_NAMESPACE:-test-auto}"
TIMEOUT="${TIMEOUT:-60s}"

# Functions
print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║${NC} $1"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
}

print_step() {
    echo -e "${YELLOW}▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# 1. Setup test-auto namespace
setup_namespace() {
    print_header "STEP 1: Setting Up test-auto Namespace"
    
    print_step "Creating namespace with RBAC..."
    kubectl apply -f setup-test-namespace.yaml
    
    # Verify namespace
    if kubectl get namespace $TEST_NAMESPACE >/dev/null 2>&1; then
        print_success "Namespace $TEST_NAMESPACE created"
    else
        print_error "Failed to create namespace"
        return 1
    fi
    
    # Verify ServiceAccount
    if kubectl get sa -n $TEST_NAMESPACE e2e-tests >/dev/null 2>&1; then
        print_success "ServiceAccount created"
    else
        print_error "Failed to create ServiceAccount"
        return 1
    fi
    
    print_success "Namespace setup complete"
    echo
}

# 2. Deploy microservices
deploy_services() {
    print_header "STEP 2: Deploying Microservices to $TEST_NAMESPACE"
    
    print_step "Deploying microservices..."
    kubectl apply -f ../release/kubernetes-manifests.yaml -n $TEST_NAMESPACE
    
    print_step "Waiting for pods to be ready (this may take 1-2 minutes)..."
    
    # Wait for all pods to be ready
    kubectl wait --for=condition=Ready pod --all -n $TEST_NAMESPACE --timeout=300s || true
    
    # Get pod status
    echo
    print_step "Pod Status:"
    kubectl get pods -n $TEST_NAMESPACE -o wide
    
    print_success "Deployment complete"
    echo
}

# 3. Run cluster accessibility tests
test_cluster_accessibility() {
    print_header "STEP 3: Testing Cluster Accessibility"
    
    print_step "Running cluster accessibility tests..."
    go test -v -timeout $TIMEOUT -run "^TestClusterAccessibility$" ./...
    
    print_success "Cluster accessibility tests passed"
    echo
}

# 4. Run node status tests
test_node_status() {
    print_header "STEP 4: Testing Node Status"
    
    print_step "Running node status tests..."
    go test -v -timeout $TIMEOUT -run "^TestNodesStatus$" ./...
    
    print_success "Node status tests passed"
    echo
}

# 5. Run pod state tests
test_pod_state() {
    print_header "STEP 5: Testing Pod State"
    
    print_step "Testing pod state in $TEST_NAMESPACE..."
    go test -v -timeout $TIMEOUT -run "TestAllPodsInTestAutoNamespace|TestPodsInBothNamespaces" ./...
    
    print_success "Pod state tests passed"
    echo
}

# 6. Run nginx-specific tests
test_nginx_pod() {
    print_header "STEP 6: Testing nginx Pod (Frontend Service)"
    
    print_step "Testing nginx pod running status..."
    go test -v -timeout $TIMEOUT -run "TestNginxPodRunningStatus" ./...
    
    print_step "Testing nginx liveness probe configuration..."
    go test -v -timeout $TIMEOUT -run "TestNginxLivenessProbeConfiguration" ./...
    
    print_step "Testing nginx readiness probe configuration..."
    go test -v -timeout $TIMEOUT -run "TestNginxReadinessProbeConfiguration" ./...
    
    print_step "Testing nginx probe functionality..."
    go test -v -timeout $TIMEOUT -run "TestNginxProbesFunctioning" ./...
    
    print_step "Checking nginx pod restarts..."
    go test -v -timeout $TIMEOUT -run "TestNginxPodRestarts" ./...
    
    print_success "nginx pod tests passed"
    echo
}

# 7. Run comprehensive cluster monitoring
test_cluster_monitoring() {
    print_header "STEP 7: Running Comprehensive Cluster Monitoring"
    
    print_step "Capturing cluster state snapshot..."
    go test -v -timeout $TIMEOUT -run "TestClusterStateSnapshot" ./...
    
    print_step "Monitoring pod lifecycle..."
    go test -v -timeout $TIMEOUT -run "TestPodLifecycleMonitoring" ./...
    
    print_step "Comparing pod health..."
    go test -v -timeout $TIMEOUT -run "TestPodHealthComparison" ./...
    
    print_step "Analyzing probe configuration..."
    go test -v -timeout $TIMEOUT -run "TestDetailedProbeAnalysis" ./...
    
    print_success "Cluster monitoring complete"
    echo
}

# 8. Run all tests
test_all() {
    print_header "STEP 8: Running All Tests"
    
    print_step "Executing all tests in verbose mode..."
    go test -v -timeout $TIMEOUT ./...
    
    print_success "All tests passed"
    echo
}

# 9. Show pod cleanup options
show_cleanup_options() {
    print_header "STEP 9: Pod Cleanup (Optional)"
    
    print_step "Simulating cleanup options..."
    go test -v -timeout $TIMEOUT -run "TestPodCleanupSimulation" ./...
    
    echo
    print_step "Cleanup commands available:"
    echo "  # Delete all pods in $TEST_NAMESPACE:"
    echo "  kubectl delete pods --all -n $TEST_NAMESPACE"
    echo
    echo "  # Delete specific pod:"
    echo "  kubectl delete pod <pod-name> -n $TEST_NAMESPACE"
    echo
    echo "  # Delete entire namespace:"
    echo "  kubectl delete namespace $TEST_NAMESPACE"
    echo
}

# 10. Show final summary
show_summary() {
    print_header "TEST EXECUTION SUMMARY"
    
    echo "Test Suite Execution Completed!"
    echo
    echo "Verification Steps:"
    echo "  1. ✓ Cluster accessibility verified"
    echo "  2. ✓ Node status checked"
    echo "  3. ✓ Pod states verified"
    echo "  4. ✓ Liveness probes tested"
    echo "  5. ✓ Readiness probes tested"
    echo "  6. ✓ Pod health compared"
    echo "  7. ✓ Cluster monitoring complete"
    echo
    echo "Namespace: $TEST_NAMESPACE"
    echo
    echo "View pod status:"
    echo "  kubectl get pods -n $TEST_NAMESPACE -o wide"
    echo
    echo "View pod details:"
    echo "  kubectl describe pod <pod-name> -n $TEST_NAMESPACE"
    echo
    echo "View logs:"
    echo "  kubectl logs <pod-name> -n $TEST_NAMESPACE"
    echo
    print_success "All E2E tests completed successfully!"
}

# Main execution
main() {
    print_header "E2E TEST SUITE - KUBERNETES RKE2 AUTOMATION"
    
    if [ "$1" == "setup" ]; then
        setup_namespace
    elif [ "$1" == "deploy" ]; then
        deploy_services
    elif [ "$1" == "test-cluster" ]; then
        test_cluster_accessibility
        test_node_status
    elif [ "$1" == "test-pods" ]; then
        test_pod_state
        test_nginx_pod
    elif [ "$1" == "test-monitoring" ]; then
        test_cluster_monitoring
    elif [ "$1" == "test-all" ]; then
        test_all
    elif [ "$1" == "cleanup" ]; then
        show_cleanup_options
    elif [ "$1" == "full" ]; then
        setup_namespace
        deploy_services
        test_cluster_accessibility
        test_node_status
        test_pod_state
        test_nginx_pod
        test_cluster_monitoring
        show_cleanup_options
        show_summary
    else
        echo "Usage: $0 {setup|deploy|test-cluster|test-pods|test-monitoring|test-all|cleanup|full}"
        echo
        echo "Options:"
        echo "  setup           - Create test-auto namespace with RBAC"
        echo "  deploy          - Deploy microservices to test-auto"
        echo "  test-cluster    - Test cluster accessibility and node status"
        echo "  test-pods       - Test pod states and nginx pod health"
        echo "  test-monitoring - Run comprehensive cluster monitoring"
        echo "  test-all        - Run all tests"
        echo "  cleanup         - Show pod cleanup options"
        echo "  full            - Execute complete automation (setup → deploy → test → cleanup)"
        echo
        echo "Environment Variables:"
        echo "  TEST_NAMESPACE  - Kubernetes namespace (default: test-auto)"
        echo "  TIMEOUT         - Test timeout (default: 60s)"
        exit 1
    fi
}

main "$@"
