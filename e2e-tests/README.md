# Kubernetes E2E Tests (Go)

Comprehensive end-to-end tests for monitoring Kubernetes cluster health, including cluster accessibility, node status, pod liveness/readiness probes, and service availability.

**Default Namespace:** `test-auto` (isolated testing environment)  
**Configurable:** Use `TEST_NAMESPACE` environment variable to test different namespaces

## Prerequisites

- Go 1.21 or higher
- kubectl configured and pointing to your cluster
- KUBECONFIG environment variable set or `~/.kube/config` available
- `test-auto` namespace created (see Setup section)

## Setup - Create test-auto Namespace

Before running tests, create the test-auto namespace with proper RBAC:

**Option A: Automated Setup (Recommended)**
```bash
cd e2e-tests
chmod +x setup-test-namespace.sh
./setup-test-namespace.sh

# Or using Make
make setup
```

**Option B: Manual Setup**
```bash
kubectl apply -f setup-test-namespace.yaml
```

**Option C: With specific namespace**
```bash
# Edit setup-test-namespace.yaml and change "test-auto" to your namespace
kubectl apply -f setup-test-namespace.yaml
```

## Installation

```bash
cd e2e-tests
go mod download
go mod tidy
```

## Running Tests

### Run all tests (in test-auto namespace)
```bash
go test -v ./...
# Or with Make
make test
```

### Run tests in different namespace
```bash
# Use default namespace
TEST_NAMESPACE=default go test -v ./...
# Or with Make
make test TEST_NAMESPACE=default
```

### Run specific test
```bash
go test -v -run TestClusterAccessibility
go test -v -run TestNodesStatus
go test -v -run TestPodLivenessProbes
```

### Run tests with timeout
```bash
go test -v -timeout 60s ./...
# Or with Make
make test TIMEOUT=120s
```

### Run tests with output logging
```bash
go test -v -run TestClusterAccessibility -v 2>&1 | tee test-results.log
```

### Using Makefile Commands
```bash
make test                          # Run all tests in test-auto
make test-verbose                  # Run with verbose output
make test-basic                    # Run basic cluster tests
make test-advanced                 # Run advanced tests (ingress, PV, events)
make test-pods                     # Run pod/deployment/service tests
make monitor                       # Continuous monitoring
make monitor MONITOR_INTERVAL=60   # Monitor every 60 seconds
```

## Available Tests

### Cluster-wide Tests (All Namespaces)
- **TestClusterAccessibility** - Verifies connection to the Kubernetes API server
- **TestNodesStatus** - Checks all nodes and their Ready status
- **TestNamespaceAccessibility** - Lists all namespaces and their status
- **TestPersistentVolumes** - Checks PV and PVC status (cluster-wide)

### Namespace-scoped Tests (test-auto by default)
- **TestPodLivenessProbes** - Verifies liveness probes are configured
- **TestPodReadinessStatus** - Checks actual readiness status of pods
- **TestCrashingPods** - Identifies pods with high restart counts
- **TestDeploymentStatus** - Checks deployment replicas and ready status
- **TestServiceAvailability** - Lists services and their configuration
- **TestFrontendServiceHealth** - Specific test for frontend service health
- **TestIngressStatus** - Validates Ingress configuration and IP assignment
- **TestEventLog** - Shows recent cluster events

## Test Output Example

```
=== RUN   TestClusterAccessibility
--- PASS: TestClusterAccessibility (0.50s)
    main_test.go:31: ✓ Cluster accessible. Kubernetes version: v1.28.0

=== RUN   TestNodesStatus
--- PASS: TestNodesStatus (0.30s)
    main_test.go:41: Found 3 nodes
    main_test.go:55: Node: node-1 | Status: Ready | Ready: True
    main_test.go:55: Node: node-2 | Status: Ready | Ready: True

=== RUN   TestPodReadinessStatus (in namespace: test-auto)
--- PASS: TestPodReadinessStatus (0.40s)
    main_test.go:95: Pod: test-auto/frontend-abc123 | Status: Ready
    main_test.go:95: Pod: test-auto/cartservice-xyz789 | Status: Ready
```

## Continuous Monitoring

To run tests periodically in test-auto namespace:

```bash
# Run tests every 5 minutes
make monitor

# Or set custom interval (in seconds)
make monitor MONITOR_INTERVAL=60

# Or with script
./run-tests.sh monitor
```

## Docker Usage

Build a Docker image for running tests:

```bash
make docker-build

# Run tests in Docker (uses test-auto namespace)
make docker-run

# Run tests in specific namespace
docker run -e TEST_NAMESPACE=default \
  -v ~/.kube/config:/root/.kube/config:ro \
  k8s-e2e-tests:latest
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `KUBECONFIG` | Path to kubeconfig file | `~/.kube/config` |
| `TEST_NAMESPACE` | Kubernetes namespace for tests | `test-auto` |
| `TIMEOUT` | Test timeout | `60s` |
| `MONITOR_INTERVAL` | Monitoring interval in seconds | `300` |

## Troubleshooting

### Tests timing out
- Increase the timeout: `go test -timeout 120s` or `make test TIMEOUT=120s`
- Check cluster connectivity: `kubectl cluster-info`

### Permission denied errors
- Ensure namespace exists: `kubectl get namespace test-auto`
- Reinstall RBAC: `make setup`
- Check permissions: `kubectl auth can-i get pods --as=system:serviceaccount:test-auto:e2e-tests -n test-auto`

### Namespace doesn't exist
```bash
# Create the namespace and RBAC
make setup

# Or manually
kubectl apply -f setup-test-namespace.yaml
```

### No pods/services found
- Verify resources are deployed in test-auto: `kubectl get pods -n test-auto`
- Deploy microservices: `kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto`
- Check namespace: `echo $TEST_NAMESPACE`

### Clean up and reset
```bash
make clean              # Clean logs and test cache
make setup-clean        # Remove test-auto namespace (destructive)
```

## Integration with CI/CD

### GitHub Actions Example
```yaml
name: K8s E2E Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: cd e2e-tests && make setup
      - run: cd e2e-tests && kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto
      - run: cd e2e-tests && make test-verbose
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
```

### Modifying Tests for Your Environment

Edit the tests to match your cluster setup:

```go
// Change namespace
namespace := "your-namespace"

// Change service names
service, err := clientset.CoreV1().Services(namespace).Get(ctx, "your-service", metav1.GetOptions{})

// Add custom checks
// ... your custom logic
```

## Performance Metrics

The tests also report:
- Cluster connectivity time
- Number of nodes, pods, and services
- Pod readiness status
- Deployment replica status
- Service configuration details

## Best Practices

1. Run tests regularly (e.g., every 5 minutes) to catch issues early
2. Store test logs for historical analysis
3. Set up alerts on test failures
4. Customize tests for your specific services
5. Use test results to create dashboards

## Contributing

To add more tests:

1. Create a new function in `main_test.go` or `advanced_tests.go`
2. Follow the naming convention: `Test<Feature>`
3. Use the global `clientset` for API calls
4. Add context with timeout for all operations
5. Log results with `t.Logf()` for debugging
