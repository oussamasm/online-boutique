# E2E Tests Integration - test-auto Namespace Setup

**Updated:** June 7, 2026

## Overview

The E2E test suite has been updated to run in a separate `test-auto` namespace for **complete isolation** from production or default resources. This prevents test interference with live services.

## What Changed

### 1. **Go Code Updates**
- ✅ Added `testNamespace` variable (defaults to `test-auto`)
- ✅ Reads `TEST_NAMESPACE` environment variable for override capability
- ✅ Added `ensureNamespaceExists()` function to create namespace automatically
- ✅ Updated all test functions to use `testNamespace` instead of hardcoded "default"
- ✅ All pod, service, deployment, and event tests now scoped to test-auto namespace

### 2. **New Files Added**

#### `setup-test-namespace.yaml`
RBAC configuration including:
- `test-auto` Namespace
- `e2e-tests` ServiceAccount
- ClusterRole for read-only access to cluster resources
- ClusterRoleBinding for cluster-wide permissions
- Role for namespace-specific permissions (create/delete capabilities)
- RoleBinding for namespace-scoped access

#### `setup-test-namespace.sh`
Automated setup script that:
- Checks kubectl connectivity
- Creates namespace and RBAC resources
- Verifies all components are installed
- Provides next-step guidance

### 3. **Updated Files**

#### `main_test.go`
- Added `testNamespace` global variable
- Added `ensureNamespaceExists()` helper function
- Updated 5 test functions to use `testNamespace`:
  - TestPodLivenessProbes
  - TestPodReadinessStatus
  - TestDeploymentStatus
  - TestServiceAvailability
  - (plus init() function modifications)

#### `advanced_tests.go`
- Updated 4 test functions to use `testNamespace`:
  - TestFrontendServiceHealth
  - TestIngressStatus
  - TestCrashingPods
  - TestEventLog
- Updated PVC listing to use `testNamespace`

#### `Makefile`
- Added `TEST_NAMESPACE` variable (default: `test-auto`)
- Added `setup` target to initialize namespace
- Added `setup-clean` target to remove namespace
- Updated all test targets with `TEST_NAMESPACE` environment variable
- Updated help documentation

#### `Dockerfile`
- Added `ENV TEST_NAMESPACE=test-auto`
- Updated CMD to use shell for environment variable support

#### `run-tests.sh`
- Added `TEST_NAMESPACE` environment variable support
- Updated all test functions to pass `TEST_NAMESPACE`
- Updated output messages to show active namespace

#### `QUICKSTART.md`
- Added namespace setup section (automated + manual options)
- Added prerequisite check steps
- Added deployment instructions for test-auto namespace
- Added namespace-specific command examples
- Added troubleshooting for namespace issues

#### `README.md`
- Added default namespace information
- Added Setup section with multiple options
- Added `TEST_NAMESPACE` environment variable documentation
- Updated test descriptions with namespace scope
- Added namespace troubleshooting section
- Updated CI/CD examples with namespace setup

---

## Quick Start

### Step 1: Setup Namespace
```bash
cd e2e-tests

# Option A: Automated (Recommended)
chmod +x setup-test-namespace.sh
./setup-test-namespace.sh

# Option B: Using Make
make setup

# Option C: Manual kubectl
kubectl apply -f setup-test-namespace.yaml
```

### Step 2: Deploy Microservices to test-auto
```bash
# Deploy the online-boutique microservices
kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto

# Wait for pods to be ready
kubectl get pods -n test-auto -w
```

### Step 3: Run Tests
```bash
# All tests (in test-auto namespace)
make test

# Specific test category
make test-basic          # Cluster + node tests
make test-advanced       # Ingress + PV + events tests
make test-pods          # Pod + deployment + service tests

# Verbose output
make test-verbose

# Continuous monitoring
make monitor            # Every 5 minutes (default)
make monitor MONITOR_INTERVAL=60  # Every 60 seconds

# Custom namespace
make test TEST_NAMESPACE=default   # Override to use default namespace
```

---

## Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `TEST_NAMESPACE` | Kubernetes namespace for tests | `test-auto` | `TEST_NAMESPACE=prod make test` |
| `TIMEOUT` | Test execution timeout | `60s` | `TIMEOUT=120s make test` |
| `MONITOR_INTERVAL` | Monitoring check interval | `300s` | `MONITOR_INTERVAL=60 make monitor` |
| `KUBECONFIG` | Path to kubeconfig | `~/.kube/config` | `KUBECONFIG=/path/to/config` |

---

## File Structure

```
e2e-tests/
├── main_test.go                  # Updated: namespace support
├── advanced_tests.go             # Updated: namespace support
├── Makefile                       # Updated: namespace targets
├── Dockerfile                     # Updated: TEST_NAMESPACE env
├── run-tests.sh                   # Updated: namespace support
├── go.mod                         # Unchanged: dependencies
├── go.sum                         # Unchanged: dependencies
├── README.md                       # Updated: namespace docs
├── QUICKSTART.md                  # Updated: setup instructions
├── .gitignore                     # Unchanged
├── setup-test-namespace.yaml      # NEW: RBAC configuration
└── setup-test-namespace.sh        # NEW: setup automation
```

---

## Test Scope by Namespace

### Cluster-wide Tests (Any Namespace)
- `TestClusterAccessibility` - API server connectivity
- `TestNodesStatus` - Node health
- `TestNamespaceAccessibility` - All namespaces
- `TestPersistentVolumes` - Cluster-wide storage

### test-auto Namespace Tests (Default)
- `TestPodLivenessProbes` - Pod probe configuration
- `TestPodReadinessStatus` - Pod readiness
- `TestDeploymentStatus` - Deployment replicas
- `TestServiceAvailability` - Services in namespace
- `TestFrontendServiceHealth` - Frontend service
- `TestIngressStatus` - Ingress configuration
- `TestCrashingPods` - Pod restart counts
- `TestEventLog` - Namespace events

---

## RBAC Permissions

The setup creates:

**ClusterRole (e2e-tests):**
- Read-only: nodes, persistentvolumes, cluster events
- Read/list/watch: pods, deployments, statefulsets, ingresses, services

**Role (e2e-tests-namespace in test-auto):**
- Full: pods, services, deployments, ingresses (create/delete for testing)
- Read/list/watch: persistentvolumeclaims, configmaps, secrets

**ServiceAccount:**
- `test-auto:e2e-tests` - Credentials for test execution

---

## CI/CD Integration Examples

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
  namespace: test-auto
spec:
  schedule: "*/5 * * * *"  # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: e2e-tests
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            env:
            - name: TEST_NAMESPACE
              value: "test-auto"
          restartPolicy: OnFailure
```

---

## Verification Steps

### Verify Namespace Setup
```bash
# Check namespace exists
kubectl get namespace test-auto

# Check ServiceAccount
kubectl get serviceaccount e2e-tests -n test-auto

# Check RBAC
kubectl get clusterrole e2e-tests
kubectl get clusterrolebinding e2e-tests
kubectl get role -n test-auto
kubectl get rolebinding -n test-auto

# Test permissions
kubectl auth can-i get pods --as=system:serviceaccount:test-auto:e2e-tests -n test-auto
```

### Verify Microservices Deployment
```bash
# Check pods in test-auto
kubectl get pods -n test-auto

# Check services
kubectl get svc -n test-auto

# Check deployments
kubectl get deployments -n test-auto

# Check ingress
kubectl get ingress -n test-auto
```

### Run Verification Test
```bash
make test-basic
# Should pass: TestClusterAccessibility, TestNodesStatus, TestNamespaceAccessibility
```

---

## Troubleshooting

### Namespace doesn't exist
```bash
make setup
# Or manually:
kubectl apply -f setup-test-namespace.yaml
```

### Permission denied errors
```bash
# Reapply RBAC
kubectl apply -f setup-test-namespace.yaml

# Verify permissions
kubectl auth can-i list pods --as=system:serviceaccount:test-auto:e2e-tests -n test-auto
```

### Tests can't find resources
```bash
# Make sure microservices are deployed
kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto

# Wait for pods
kubectl get pods -n test-auto -w
```

### Test failures with "namespace not found"
```bash
# The ensureNamespaceExists() function will auto-create it
# But you can manually create it:
make setup

# Or verify with:
kubectl get namespace test-auto
```

### Run tests in different namespace
```bash
# Override TEST_NAMESPACE
TEST_NAMESPACE=default make test

# Or with Make:
make test TEST_NAMESPACE=default
```

### Clean up everything
```bash
# Remove namespace (destructive - deletes all resources)
make setup-clean

# Clean test cache and logs
make clean

# Clean Docker image
make docker-clean
```

---

## Summary of Commands

```bash
# Setup
make setup                          # Create test-auto namespace
make setup-clean                    # Remove test-auto namespace

# Run Tests
make test                           # All tests in test-auto
make test-verbose                   # Verbose output
make test-basic                     # Basic cluster tests
make test-advanced                  # Advanced tests
make test-pods                      # Pod tests

# Monitoring
make monitor                        # Every 5 minutes
make monitor MONITOR_INTERVAL=60    # Every 60 seconds

# Docker
make docker-build                   # Build image
make docker-run                     # Run in container

# Utilities
make clean                          # Clean logs
make help                           # Show all options
make deps                           # Download dependencies
```

---

## Benefits of test-auto Namespace

✅ **Isolation** - Tests don't affect default or production namespaces  
✅ **Safety** - Easy cleanup with `kubectl delete namespace test-auto`  
✅ **Flexibility** - Use `TEST_NAMESPACE=default` to test other namespaces  
✅ **CI/CD Ready** - Repeatable, clean state for each test run  
✅ **Multi-environment** - Run tests on dev, staging, and prod with same code  
✅ **Automatic Setup** - Namespace created automatically on first test run if needed  

---

## Next Steps

1. ✅ Run `make setup` to initialize namespace
2. ✅ Deploy microservices: `kubectl apply -f ../release/kubernetes-manifests.yaml -n test-auto`
3. ✅ Run tests: `make test-verbose`
4. ✅ Set up continuous monitoring: `make monitor`
5. ✅ Integrate with CI/CD pipeline
6. ✅ Add alerting for test failures (Slack, PagerDuty, etc.)

---

**Status:** ✅ Ready for production use
