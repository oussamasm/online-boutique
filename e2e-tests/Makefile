.PHONY: help deps test test-verbose test-all test-basic test-advanced test-pods monitor clean docker-build docker-run docker-clean setup setup-clean

# Variables
TIMEOUT ?= 60s
MONITOR_INTERVAL ?= 300
TEST_NAMESPACE ?= test-auto
GO := go
DOCKER_IMAGE ?= k8s-e2e-tests
DOCKER_TAG ?= latest

help:
	@echo "Kubernetes E2E Tests - Makefile Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  setup             Setup test-auto namespace with RBAC"
	@echo "  setup-clean       Remove test-auto namespace (WARNING: destructive)"
	@echo "  deps              Download dependencies"
	@echo "  test              Run all tests (default namespace: test-auto)"
	@echo "  test-verbose      Run all tests with verbose output"
	@echo "  test-basic        Run basic cluster tests"
	@echo "  test-advanced     Run advanced tests (ingress, PV, events)"
	@echo "  test-pods         Run pod and service tests"
	@echo "  monitor           Run tests continuously"
	@echo "  clean             Clean up logs and artifacts"
	@echo "  docker-build      Build Docker image"
	@echo "  docker-run        Run tests in Docker"
	@echo "  docker-clean      Remove Docker image"
	@echo ""
	@echo "Variables:"
	@echo "  TEST_NAMESPACE    Kubernetes namespace for tests (default: test-auto)"
	@echo "  TIMEOUT           Test timeout (default: 60s)"
	@echo "  MONITOR_INTERVAL  Monitoring interval in seconds (default: 300)"
	@echo ""
	@echo "Examples:"
	@echo "  make setup                                  # Setup test-auto namespace"
	@echo "  make test                                   # Run all tests in test-auto"
	@echo "  make test TEST_NAMESPACE=default            # Run tests in default namespace"
	@echo "  make test TIMEOUT=120s                      # Run with 2 minute timeout"
	@echo "  make monitor MONITOR_INTERVAL=60            # Monitor every 60 seconds"
	@echo "  make docker-run                             # Run in Docker"

deps:
	@echo "[*] Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "[✓] Dependencies downloaded"

setup:
	@echo "[*] Setting up test-auto namespace..."
	@if [ -f setup-test-namespace.sh ]; then \
		chmod +x setup-test-namespace.sh && \
		./setup-test-namespace.sh; \
	else \
		kubectl apply -f setup-test-namespace.yaml; \
	fi
	@echo "[✓] Namespace setup complete"

setup-clean:
	@echo "[!] WARNING: This will delete the test-auto namespace and all resources!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		kubectl delete namespace test-auto --ignore-not-found=true; \
		echo "[✓] Namespace deleted"; \
	else \
		echo "[*] Operation cancelled"; \
	fi

test:
	@echo "[*] Running E2E tests (timeout: $(TIMEOUT), namespace: $(TEST_NAMESPACE))..."
	TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -timeout $(TIMEOUT) ./...

test-verbose:
	@echo "[*] Running E2E tests with verbose output (namespace: $(TEST_NAMESPACE))..."
	TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -v -timeout $(TIMEOUT) ./...

test-basic:
	@echo "[*] Running basic cluster tests (namespace: $(TEST_NAMESPACE))..."
	TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -v -timeout $(TIMEOUT) -run "^TestCluster|^TestNodes|^TestNamespace" ./...

test-advanced:
	@echo "[*] Running advanced tests (namespace: $(TEST_NAMESPACE))..."
	TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -v -timeout $(TIMEOUT) -run "^TestFrontend|^TestIngress|^TestPersistent|^TestCrashing|^TestEvent" ./....

test-pods:
	@echo "[*] Running pod and service tests (namespace: $(TEST_NAMESPACE))..."
	TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -v -timeout $(TIMEOUT) -run "^TestPod|^TestDeployment|^TestService" ./....

test-all: test-verbose
	@echo "[✓] All tests completed"

monitor:
	@echo "[*] Starting continuous monitoring (interval: $(MONITOR_INTERVAL)s, namespace: $(TEST_NAMESPACE))..."
	@echo "[*] Press Ctrl+C to stop"
	@mkdir -p logs
	@counter=0; \
	while true; do \
		counter=$$((counter+1)); \
		echo ""; \
		echo "========== Run #$$counter at $$(date) =========="; \
		TEST_NAMESPACE=$(TEST_NAMESPACE) $(GO) test -v -timeout $(TIMEOUT) ./... | tee -a logs/monitoring-$$(date +%Y%m%d).log; \
		echo "Next run in $(MONITOR_INTERVAL)s..."; \
		sleep $(MONITOR_INTERVAL); \
	done

clean:
	@echo "[*] Cleaning up..."
	rm -rf logs/*.log
	$(GO) clean -testcache
	@echo "[✓] Cleanup completed"

docker-build:
	@echo "[*] Building Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "[✓] Docker image built"

docker-run:
	@echo "[*] Running tests in Docker (namespace: $(TEST_NAMESPACE))..."
	docker run -e TEST_NAMESPACE=$(TEST_NAMESPACE) -v ~/.kube/config:/root/.kube/config:ro $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-run-monitor:
	@echo "[*] Running tests in Docker with monitoring..."
	docker run -it -v ~/.kube/config:/root/.kube/config:ro \
		-e "MONITOR_INTERVAL=$(MONITOR_INTERVAL)" \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		/bin/sh -c 'apk add bash && chmod +x run-tests.sh && ./run-tests.sh monitor'

docker-clean:
	@echo "[*] Removing Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)..."
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || echo "Image not found"
	@echo "[✓] Docker cleanup completed"

lint:
	@echo "[*] Running Go linter..."
	$(GO) fmt ./...
	$(GO) vet ./...
	@echo "[✓] Linting completed"

coverage:
	@echo "[*] Running tests with coverage..."
	$(GO) test -cover ./...

coverage-html:
	@echo "[*] Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "[✓] Coverage report generated: coverage.html"

.DEFAULT_GOAL := help
