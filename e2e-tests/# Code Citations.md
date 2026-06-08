# Code Citations

## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef:
  apiGroup: rbac
```


## License: unknown
https://github.com/mikerain/k8s-book/blob/1351c5bbeeb5ef5d1492ba542dfd999bdca2a480/user/user.md

```
---

## **E2E Tests Project - Comprehensive Analysis**

### **Project Overview**

This is a **Kubernetes End-to-End Testing Suite** written in Go that monitors the health and accessibility of your online-boutique microservices cluster. It tests cluster connectivity, node status, pod health, ingress configuration, persistent volumes, and service availability.

---

## **FILE-BY-FILE BREAKDOWN**

### **1. `go.mod` - Go Module Dependencies**
**Purpose:** Defines project dependencies and Go version  
**What it does:**
- Declares Go 1.21 as minimum version
- Defines three main Kubernetes client libraries:
  - `k8s.io/api` - Kubernetes API types
  - `k8s.io/apimachinery` - Kubernetes common utilities
  - `k8s.io/client-go` - Official Kubernetes client for Go

**Integration:** Already properly configured—no changes needed. Dependencies will auto-download when you run tests.

---

### **2. `main_test.go` - Core Test Functions**
**Purpose:** Contains basic cluster health and pod monitoring tests  
**What it does:** Initializes the Kubernetes client and runs 6 core test functions

#### **Detailed Function Breakdown:**

```
┌─ init() ─────────────────────────────────────────────────────────────┐
│ • Initializes the Kubernetes client connection                       │
│ • Reads KUBECONFIG from environment or ~/.kube/config               │
│ • Creates clientset for interacting with API                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestClusterAccessibility() ─────────────────────────────────────────┐
│ Purpose: Verify cluster is reachable                                 │
│ Logic:                                                               │
│   1. Create 10-second timeout context                               │
│   2. Call clientset.Discovery().ServerVersion()                     │
│   3. Log Kubernetes version if successful                           │
│ Output: "✓ Cluster accessible. Kubernetes version: v1.28.0"        │
│ Fail if: API server unreachable or network error                   │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNodesStatus() ──────────────────────────────────────────────────┐
│ Purpose: Check all nodes are in Ready state                         │
│ Logic:                                                               │
│   1. List all nodes using clientset.CoreV1().Nodes()               │
│   2. Loop through each node                                         │
│   3. Extract NodeReady condition status from node.Status.Conditions │
│   4. Log each node: "Node: name | Status: phase | Ready: true/false"│
│   5. Fail if ANY node status ≠ "True"                              │
│ Output Example:                                                      │
│   Found 3 nodes                                                      │
│   Node: node-1 | Status: Ready | Ready: True                        │
│   Node: node-2 | Status: Ready | Ready: True                        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodLivenessProbes() ────────────────────────────────────────────┐
│ Purpose: Verify liveness/readiness probes are configured             │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, iterate through containers                       │
│   3. Check: container.LivenessProbe != nil                          │
│   4. Check: container.ReadinessProbe != nil                         │
│   5. Log with symbols: ✓ (configured) or ⚠ (missing)              │
│ Output: Shows which containers are missing probe configurations     │
│ Warning: Logs alert but doesn't fail (informational)               │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPodReadinessStatus() ───────────────────────────────────────────┐
│ Purpose: Check actual readiness state of running pods                │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each pod, call isPodReady() helper function                │
│   3. Log: "Pod: namespace/name | Status: Ready/NotReady"            │
│   4. Collect running but not-ready pods                             │
│   5. Log warning if any pods running but not ready                  │
│ Key Helper: isPodReady() checks pod.Status.Conditions for Ready=True│
└─────────────────────────────────────────────────────────────────────┘

┌─ TestNamespaceAccessibility() ───────────────────────────────────────┐
│ Purpose: Verify all namespaces are accessible                       │
│ Logic:                                                               │
│   1. List all namespaces via clientset.CoreV1().Namespaces()       │
│   2. Loop and log each: "namespace-name (Status: Active)"           │
│   3. Fail if zero namespaces found (unusual cluster state)          │
│ Output: Confirms cluster has namespaces and they're accessible      │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestDeploymentStatus() ─────────────────────────────────────────────┐
│ Purpose: Check deployment replica readiness                         │
│ Logic:                                                               │
│   1. List deployments in "default" namespace                        │
│   2. Compare: spec.Replicas vs status.ReadyReplicas                │
│   3. Log: "Deployment: name | Replicas: ready/desired"             │
│   4. Fail if any deployment has mismatched replica counts           │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **3. `advanced_tests.go` - Advanced Monitoring Functions**
**Purpose:** Contains specialized tests for frontend, ingress, volumes, and events

#### **Detailed Advanced Functions:**

```
┌─ TestFrontendServiceHealth() ────────────────────────────────────────┐
│ Purpose: Specific health check for frontend service                 │
│ Logic:                                                               │
│   1. Get frontend service from "default" namespace                  │
│   2. Extract label selector from service.Spec.Selector              │
│   3. List pods matching those labels                                │
│   4. For each pod, call isPodReady()                                │
│   5. Log each pod state with restart counts                         │
│   6. Fail if any pod not ready                                      │
│ Output: Shows frontend pods and their container status details      │
│ Use Case: Verify frontend can receive traffic                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestIngressStatus() ────────────────────────────────────────────────┐
│ Purpose: Verify Ingress is properly configured and has IP assigned  │
│ Logic:                                                               │
│   1. List ingresses in "default" namespace                          │
│   2. For each ingress:                                              │
│      a. Log ingress name, class (nginx), rules                      │
│      b. Check status.LoadBalancer.Ingress (external IP/hostname)   │
│      c. List TLS configuration and secrets used                     │
│      d. Print all rules and backend services                        │
│   3. Warn if no IP/hostname assigned yet                            │
│ Output Example:                                                      │
│   Ingress: default/frontend-ingress | Class: nginx                 │
│   ✓ IP: 192.168.1.139                                              │
│   TLS configured for 2 hosts                                        │
│   Secret: frontend-tls | Hosts: [192.168.1.139.nip.io, ...]        │
│ Key Check: Verifies TLS secret is referenced correctly              │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestPersistentVolumes() ────────────────────────────────────────────┐
│ Purpose: Check storage system health (PV/PVC)                       │
│ Logic:                                                               │
│   1. List all PersistentVolumes in cluster                          │
│   2. Log: "PV: name | Status: Bound/Available | Capacity: size"    │
│   3. List PersistentVolumeClaims in "default" namespace             │
│   4. Log: "PVC: name | Status: Bound | VolumeName: pv-name"        │
│ Warning: Logs if PVs/PVCs don't exist (not an error for this app)  │
│ Output: Shows redis-cart volume and any application storage        │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestCrashingPods() ─────────────────────────────────────────────────┐
│ Purpose: Identify unhealthy pods with restart loops                 │
│ Logic:                                                               │
│   1. List all pods in "default" namespace                           │
│   2. For each container, check RestartCount                         │
│   3. If RestartCount > 3, mark pod as "crashing"                   │
│   4. Log: "pod-name (container-name) - N restarts"                  │
│ Threshold: 3 restarts = alert (configurable in production)         │
│ Output: List of pods that have restarted too many times             │
│ Use Case: Early warning of application issues                       │
└─────────────────────────────────────────────────────────────────────┘

┌─ TestEventLog() ─────────────────────────────────────────────────────┐
│ Purpose: Show recent cluster events for troubleshooting              │
│ Logic:                                                               │
│   1. List events from "default" namespace                           │
│   2. Filter events: LastTimestamp > now() - 5 minutes              │
│   3. Log each: "Type: Reason | Namespace/Object | Count: N"        │
│ Use Case: Detect pod failures, image pull errors, etc.             │
│ Output Example:                                                      │
│   Normal: BackOff | default/cartservice | Count: 5                 │
│   Warning: ImagePullBackOff | default/frontend | Count: 2          │
└─────────────────────────────────────────────────────────────────────┘
```

---

### **4. `Makefile` - Build & Test Automation**
**Purpose:** Provides 10+ convenient commands to run tests

**Available Commands:**
```bash
make deps              # Download Go dependencies
make test              # Run all tests
make test-verbose      # Run with detailed output
make test-basic        # Run only cluster/node/namespace tests
make test-advanced     # Run ingress/PV/events tests
make test-pods         # Run pod/deployment/service tests
make test-all          # Run everything
make monitor           # Run tests repeatedly every 5 minutes
make docker-build      # Build Docker image
make docker-run        # Run tests in Docker container
make clean             # Clean up test logs and cache
```

**Integration:** Add these to your CI/CD pipeline:
```yaml
# Example GitHub Actions
- name: Run E2E Tests
  run: cd e2e-tests && make test-verbose
```

---

### **5. `Dockerfile` - Container Configuration**
**Purpose:** Package tests to run in isolated environment

**What it does:**
- Builds on `golang:1.21-alpine` (lightweight)
- Sets working directory to `/app`
- Downloads dependencies via `go mod download`
- Copies all source files
- Default command: `go test -v ./...`

**Integration:** Use for:
- CI/CD pipelines (no local Go installation needed)
- Kubernetes Job/CronJob for automatic monitoring
- Isolated testing environment

---

### **6. `run-tests.sh` - Bash Test Runner Script**
**Purpose:** Convenient shell wrapper around test commands

**What it does:**
```bash
./run-tests.sh all           # Run all tests, save logs
./run-tests.sh basic         # Run basic cluster tests
./run-tests.sh advanced      # Run advanced tests
./run-tests.sh monitor       # Run continuously with 5-min interval
```

**Features:**
- Auto-creates `logs/` directory
- Timestamps all test output files
- Supports custom timeout: `TIMEOUT=120s ./run-tests.sh all`
- Supports custom monitor interval: `MONITOR_INTERVAL=60 ./run-tests.sh monitor`

**Integration:** Add to crontab for periodic monitoring:
```bash
*/5 * * * * cd /path/to/e2e-tests && MONITOR_INTERVAL=300 ./run-tests.sh monitor >> /var/log/e2e-monitoring.log
```

---

### **7. `README.md` & `QUICKSTART.md` - Documentation**
**Purpose:** User guides and reference

**Key sections:**
- Installation steps
- Running individual tests
- Docker usage
- Test output interpretation
- Troubleshooting guide

---

### **8. `.gitignore` - Git Configuration**
**Purpose:** Exclude test artifacts from version control

**Excludes:**
- Go binaries (`*.exe`, `*.dll`, `*.so`)
- Test logs (`logs/`, `*.log`)
- IDE files (`.vscode/`, `.idea/`)
- Kubeconfig files (`.kubeconfig`, `kubeconfig-*`)

---

## **INTEGRATION GUIDE - How to Integrate into Your Project**

### **Step 1: Verify Prerequisites**
```powershell
# Check Go installation
go version  # Should show Go 1.21+

# Check kubectl/kubeconfig
kubectl cluster-info
echo $KUBECONFIG  # Should point to your cluster config
```

### **Step 2: Setup & Install Dependencies**
```bash
cd online-boutique/e2e-tests
go mod download
go mod tidy
```

### **Step 3: Run Tests Locally**
```bash
# Basic sanity check
go test -v -run TestClusterAccessibility

# Run all tests
make test-verbose

# Run with monitoring
make monitor MONITOR_INTERVAL=60
```

### **Step 4: Integrate into CI/CD Pipeline**

**GitHub Actions Example:**
```yaml
name: E2E Tests
on:
  schedule:
    - cron: '*/10 * * * *'  # Run every 10 minutes
  push:
    branches: [main]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21
      
      - name: Run E2E Tests
        run: |
          cd e2e-tests
          make deps
          make test-verbose
```

### **Step 5: Deploy as Kubernetes CronJob**
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
          containers:
          - name: e2e-tests
            image: k8s-e2e-tests:latest
            volumeMounts:
            - name: kubeconfig
              mountPath: /root/.kube
          volumes:
          - name: kubeconfig
            configMap:
              name: kubeconfig
          restartPolicy: OnFailure
```

### **Step 6: Create Required RBAC (Service Account)**
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: e2e-tests
  namespace: default
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: e2e-tests
rules:
- apiGroups: [""]
  resources: ["nodes", "pods", "services", "namespaces", "events", "persistentvolumes", "persistentvolumeclaims"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments", "statefulsets"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["networking.k8s.io"]
  resources: ["ingresses"]
  verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: e2e-tests
roleRef:
  apiGroup: rbac
```

