# Quick Start Guide - Kubernetes E2E Tests

## 1. Prerequisites

```bash
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

## 2. Setup test-auto Namespace (Recommended)

The E2E tests default to running in a separate `test-auto` namespace for isolation.

**Option A: Automated Setup**
```bash
cd e2e-tests

# Make setup script executable and run
chmod +x setup-test-namespace.sh
./setup-test-namespace.sh

# Or use Make
make setup
```

**Option B: Manual Setup**
```bash
# Apply RBAC and namespace resources
kubectl apply -f setup-test-namespace.yaml
```

**Verify Setup:**
```bash
kubectl get namespace test-auto
kubectl get serviceaccount e2e-tests -n test-auto
kubectl get clusterrolebinding e2e-tests
```

## 3. Deploy Microservices to Test Namespace

Before running tests, deploy your microservices to the test-auto namespace:

```bash
# Deploy from kubernetes-manifests
kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto

# Verify deployment
kubectl get all -n test-auto
kubectl get pods -n test-auto -w  # Watch pods come up
```

## 4. Install Dependencies

```bash
cd e2e-tests
go mod download
go mod tidy
```

## 5. Run Tests

### Option A: Run all tests
```bash
go test -v ./...
# Or with Make
make test
```

### Option B: Run with script (Linux/Mac)
```bash
chmod +x run-tests.sh
./run-tests.sh all
```

### Option C: Run specific test
```bash
go test -v -run TestClusterAccessibility
go test -v -run TestNodesStatus
go test -v -run TestPodLivenessProbes
```

### Option D: Run in different namespace
```bash
# Use default namespace instead
TEST_NAMESPACE=default make test

# Or with script
TEST_NAMESPACE=default ./run-tests.sh all
```

## 6. Monitor Cluster Continuously

```bash
# Run tests every 5 minutes in test-auto namespace
make monitor

# Or set custom interval (in seconds)
MONITOR_INTERVAL=60 make monitor

# Or with script
MONITOR_INTERVAL=60 ./run-tests.sh monitor
```

## 7. Docker Execution

```bash
# Build Docker image
make docker-build

# Run tests in Docker (uses test-auto namespace)
make docker-run

# Run tests in specific namespace
docker run -e TEST_NAMESPACE=default \
  -v ~/.kube/config:/root/.kube/config:ro \
  k8s-e2e-tests:latest
```

## 8. Understanding Output

```
✓ = Component is healthy/configured
⚠ = Component is missing/not configured
PASS = Test passed
FAIL = Test failed
```

Example:
```
✓ Cluster accessible. Kubernetes version: v1.28.0
✓ Liveness probe configured
⚠ No readiness probe configured
```

## 9. Test Summary

| Test Name | Namespace | Purpose |
|-----------|-----------|---------|
| TestClusterAccessibility | Cluster-wide | Can reach API server |
| TestNodesStatus | Cluster-wide | All nodes in Ready state |
| TestPodLivenessProbes | test-auto | Liveness probes configured |
| TestPodReadinessStatus | test-auto | Pods in Ready state |
| TestDeploymentStatus | test-auto | Replicas match ready count |
| TestServiceAvailability | test-auto | Services accessible |
| TestFrontendServiceHealth | test-auto | Frontend service & pods healthy |
| TestIngressStatus | test-auto | Ingress has IP/hostname assigned |
| TestPersistentVolumes | Cluster-wide | PV/PVC status |
| TestCrashingPods | test-auto | No excessive restarts |
| TestEventLog | test-auto | Recent cluster events |

## 10. Useful Commands

```bash
# View test-auto namespace resources
kubectl get all -n test-auto

# View logs from tests
tail -f logs/monitoring-$(date +%Y%m%d).log

# Clean up test logs
make clean

# Remove test-auto namespace (WARNING: destructive)
make setup-clean

# Run tests with longer timeout
make test TIMEOUT=120s

# Run only basic tests
make test-basic

# Run only advanced tests (ingress, PV, events)
make test-advanced

# Run pod/service tests
make test-pods
```

## 11. Integration with CI/CD

### GitHub Actions
```yaml
- name: Setup E2E Tests
  run: |
    cd e2e-tests
    make setup
    kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto

- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

### Kubernetes CronJob
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: e2e-tests
spec:
  schedule: "*/5 * * * *"  # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: e2e-tests
          namespace: test-auto
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            env:
            - name: TEST_NAMESPACE
              value: "test-auto"
          restartPolicy: OnFailure
```

## 12. Troubleshooting

**Tests can't connect to cluster:**
```bash
# Check kubeconfig
kubectl cluster-info

# Verify kubeconfig path
echo $KUBECONFIG
# If empty, ensure ~/.kube/config exists
```

**Namespace errors:**
```bash
# Verify namespace exists
kubectl get namespace test-auto

# Check RBAC permissions
kubectl auth can-i list pods --as=system:serviceaccount:test-auto:e2e-tests -n test-auto
```

**Permission denied:**
```bash
# Reinstall RBAC
kubectl apply -f setup-test-namespace.yaml
```
# Check RBAC permissions
kubectl auth can-i get pods

# Check with specific user/group
kubectl auth can-i get pods --as=system:serviceaccount:default:default
```

**No resources found:**
```bash
# Check if resources exist
kubectl get pods -n default
kubectl get nodes
kubectl get services -n default
```

## 8. Integration with CI/CD

### GitLab CI
```yaml
e2e-tests:
  image: golang:1.21
  script:
    - cd e2e-tests
    - go mod download
    - go test -v ./...
```

### Jenkins
```groovy
stage('E2E Tests') {
    steps {
        dir('e2e-tests') {
            sh 'go mod download'
            sh 'go test -v ./...'
        }
    }
}
```

### Kubernetes CronJob
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: k8s-e2e-tests
spec:
  schedule: "*/5 * * * *"  # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: e2e-tester
          containers:
          - name: tests
            image: k8s-e2e-tests:latest
          restartPolicy: OnFailure
```

## 9. Next Steps

1. **Customize tests** - Modify `main_test.go` and `advanced_tests.go` for your services
2. **Add alerts** - Integrate test failures with monitoring (Prometheus, Datadog, etc.)
3. **Export metrics** - Add test duration/status metrics for dashboards
4. **Set up logging** - Configure centralized log aggregation
5. **Create dashboards** - Visualize test results over time

## 10. Useful Commands

```bash
# Run tests with verbose output
go test -v ./...

# Run tests and save logs
go test -v ./... | tee test-results.log

# Run only failed tests again
go test -v -run FAILED ./...

# Run with code coverage
go test -cover ./...

# Run tests in parallel (faster)
go test -parallel 4 -v ./...

# Run with timeout
go test -timeout 120s -v ./...

# List all available tests
go test -list ".*" ./...
```

## Support

For issues or questions:
1. Check test logs: `cat logs/e2e-tests-*.log`
2. Verify cluster status: `kubectl get nodes` and `kubectl get pods`
3. Review Kubernetes API errors in test output
4. Check cert-manager/ingress controller logs if testing TLS
