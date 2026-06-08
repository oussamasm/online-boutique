#!/bin/bash

# Kubernetes E2E Tests Runner Script
# Usage: ./run-tests.sh [all|basic|advanced|monitor]

set -e

TEST_TYPE=${1:-all}
TIMEOUT=${TIMEOUT:-60s}
TEST_NAMESPACE=${TEST_NAMESPACE:-test-auto}
LOG_DIR="logs"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

# Create logs directory
mkdir -p "$LOG_DIR"

echo "=========================================="
echo "Kubernetes E2E Tests"
echo "=========================================="
echo "Test Type: $TEST_TYPE"
echo "Timeout: $TIMEOUT"
echo "Namespace: $TEST_NAMESPACE"
echo "Log File: $LOG_DIR/e2e-tests-$TIMESTAMP.log"
echo "=========================================="

run_all_tests() {
    echo -e "\n[*] Running all tests in namespace: $TEST_NAMESPACE..."
    TEST_NAMESPACE=$TEST_NAMESPACE go test -v -timeout "$TIMEOUT" ./... 2>&1 | tee "$LOG_DIR/e2e-tests-$TIMESTAMP.log"
}

run_basic_tests() {
    echo -e "\n[*] Running basic cluster tests in namespace: $TEST_NAMESPACE..."
    TEST_NAMESPACE=$TEST_NAMESPACE go test -v -timeout "$TIMEOUT" -run "^TestCluster|^TestNodes|^TestNamespace" ./... 2>&1 | tee "$LOG_DIR/basic-tests-$TIMESTAMP.log"
}

run_advanced_tests() {
    echo -e "\n[*] Running advanced tests in namespace: $TEST_NAMESPACE..."
    TEST_NAMESPACE=$TEST_NAMESPACE go test -v -timeout "$TIMEOUT" -run "^TestFrontend|^TestIngress|^TestPersistent|^TestCrashing|^TestEvent" ./... 2>&1 | tee "$LOG_DIR/advanced-tests-$TIMESTAMP.log"
}

run_pod_tests() {
    echo -e "\n[*] Running pod health tests in namespace: $TEST_NAMESPACE..."
    TEST_NAMESPACE=$TEST_NAMESPACE go test -v -timeout "$TIMEOUT" -run "^TestPod|^TestDeployment|^TestService" ./... 2>&1 | tee "$LOG_DIR/pod-tests-$TIMESTAMP.log"
}

run_continuous_monitoring() {
    INTERVAL=${MONITOR_INTERVAL:-300}  # Default 5 minutes
    echo -e "\n[*] Starting continuous monitoring in namespace: $TEST_NAMESPACE (interval: ${INTERVAL}s)..."
    
    counter=0
    while true; do
        counter=$((counter+1))
        echo -e "\n============ Run #$counter at $(date) ============"
        go test -v -timeout "$TIMEOUT" ./... 2>&1 | tee -a "$LOG_DIR/monitoring-$(date +%Y%m%d).log"
        
        echo "Next run in ${INTERVAL}s... (Press Ctrl+C to stop)"
        sleep "$INTERVAL"
    done
}

# Parse command line arguments
case "$TEST_TYPE" in
    all)
        run_all_tests
        ;;
    basic)
        run_basic_tests
        ;;
    advanced)
        run_advanced_tests
        ;;
    pods)
        run_pod_tests
        ;;
    monitor)
        run_continuous_monitoring
        ;;
    *)
        echo "Usage: ./run-tests.sh [all|basic|advanced|pods|monitor]"
        echo ""
        echo "Options:"
        echo "  all       - Run all tests (default)"
        echo "  basic     - Run basic cluster accessibility tests"
        echo "  advanced  - Run advanced tests (ingress, PV, events, etc.)"
        echo "  pods      - Run pod and service health tests"
        echo "  monitor   - Run tests continuously (set MONITOR_INTERVAL for interval in seconds)"
        echo ""
        echo "Environment variables:"
        echo "  TIMEOUT - Test timeout (default: 60s)"
        echo "  MONITOR_INTERVAL - Monitoring interval in seconds (default: 300)"
        exit 1
        ;;
esac

echo -e "\n=========================================="
echo "Tests completed at $(date)"
echo "Log saved to: $LOG_DIR/e2e-tests-$TIMESTAMP.log"
echo "=========================================="
