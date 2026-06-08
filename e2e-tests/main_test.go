package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var clientset *kubernetes.Clientset
var testNamespace string

// Initialize the Kubernetes client
func init() {
	// Get test namespace from environment or use default
	testNamespace = os.Getenv("TEST_NAMESPACE")
	if testNamespace == "" {
		testNamespace = "test-auto"
	}
	log.Printf("Using namespace: %s", testNamespace)
	// Get kubeconfig path from environment or use default
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get home directory: %v", err)
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to build kubeconfig: %v", err)
	}

	var clientErr error
	clientset, clientErr = kubernetes.NewForConfig(config)
	if clientErr != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", clientErr)
	}
}

// ensureNamespaceExists creates the test namespace if it doesn't exist
func ensureNamespaceExists(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ns, err := clientset.CoreV1().Namespaces().Get(ctx, testNamespace, metav1.GetOptions{})
	if err == nil {
		t.Logf("Namespace %s already exists", testNamespace)
		return
	}

	// Namespace doesn't exist, create it
	newNS := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}

	_, err = clientset.CoreV1().Namespaces().Create(ctx, newNS, metav1.CreateOptions{})
	if err != nil {
		t.Logf("Warning: Could not create namespace %s: %v", testNamespace, err)
	}
}

// TestClusterAccessibility tests if the cluster is accessible
func TestClusterAccessibility(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to get the API server version
	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		t.Fatalf("Failed to connect to cluster: %v", err)
	}

	if versionInfo == nil {
		t.Fatal("Server version info is nil")
	}

	t.Logf("✓ Cluster accessible. Kubernetes version: %s", versionInfo.GitVersion)
}

// TestNodesStatus tests the status of all nodes in the cluster
func TestNodesStatus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodes: %v", err)
	}

	if len(nodes.Items) == 0 {
		t.Fatal("No nodes found in the cluster")
	}

	t.Logf("Found %d nodes\n", len(nodes.Items))

	allHealthy := true
	for _, node := range nodes.Items {
		status := "Unknown"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				status = string(condition.Status)
				break
			}
		}

		t.Logf("Node: %s | Status: %s | Ready: %s", node.Name, node.Status.Phase, status)

		if status != "True" {
			allHealthy = false
		}
	}

	if !allHealthy {
		t.Error("Some nodes are not in Ready state")
	}
}

// TestPodLivenessProbes tests if pods have liveness probes configured
func TestPodLivenessProbes(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods in namespace %s: %v", namespace, err)
	}

	if len(pods.Items) == 0 {
		t.Logf("No pods found in namespace %s", namespace)
		return
	}

	t.Logf("Found %d pods in namespace %s\n", len(pods.Items), namespace)

	for _, pod := range pods.Items {
		t.Logf("Pod: %s/%s | Phase: %s", pod.Namespace, pod.Name, pod.Status.Phase)

		for _, container := range pod.Spec.Containers {
			hasLiveness := container.LivenessProbe != nil
			hasReadiness := container.ReadinessProbe != nil

			if !hasLiveness {
				t.Logf("  ⚠ Container %s: No liveness probe configured", container.Name)
			} else {
				t.Logf("  ✓ Container %s: Liveness probe configured", container.Name)
			}

			if !hasReadiness {
				t.Logf("  ⚠ Container %s: No readiness probe configured", container.Name)
			} else {
				t.Logf("  ✓ Container %s: Readiness probe configured", container.Name)
			}
		}
	}
}

// TestPodReadinessStatus tests the actual readiness status of pods
func TestPodReadinessStatus(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods in namespace %s: %v", namespace, err)
	}

	if len(pods.Items) == 0 {
		t.Logf("No pods found in namespace %s", namespace)
		return
	}

	notReadyPods := []string{}
	
	for _, pod := range pods.Items {
		ready := isPodReady(&pod)
		status := "Not Ready"
		if ready {
			status = "Ready"
		}

		t.Logf("Pod: %s/%s | Status: %s", pod.Namespace, pod.Name, status)

		if !ready && pod.Status.Phase == corev1.PodRunning {
			notReadyPods = append(notReadyPods, fmt.Sprintf("%s/%s", pod.Namespace, pod.Name))
		}
	}

	if len(notReadyPods) > 0 {
		t.Logf("Pods not ready (but running): %v", notReadyPods)
	}
}

// TestNamespaceAccessibility tests if we can access different namespaces
func TestNamespaceAccessibility(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list namespaces: %v", err)
	}

	t.Logf("Found %d namespaces\n", len(namespaces.Items))
	for _, ns := range namespaces.Items {
		t.Logf("  - %s (Status: %s)", ns.Name, ns.Status.Phase)
	}

	if len(namespaces.Items) == 0 {
		t.Error("No namespaces found in the cluster")
	}
}

// TestDeploymentStatus tests the status of deployments
func TestDeploymentStatus(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list deployments in namespace %s: %v", namespace, err)
	}

	if len(deployments.Items) == 0 {
		t.Logf("No deployments found in namespace %s", namespace)
		return
	}

	t.Logf("Found %d deployments in namespace %s\n", len(deployments.Items), namespace)

	for _, deployment := range deployments.Items {
		ready := deployment.Status.ReadyReplicas == deployment.Status.Replicas
		status := "Not Ready"
		if ready {
			status = "Ready"
		}

		t.Logf("Deployment: %s/%s | Ready Replicas: %d/%d | Status: %s",
			deployment.Namespace,
			deployment.Name,
			deployment.Status.ReadyReplicas,
			deployment.Status.Replicas,
			status)

		if !ready {
			t.Logf("  ⚠ Deployment %s is not fully ready", deployment.Name)
		}
	}
}

// TestServiceAvailability tests the availability of services
func TestServiceAvailability(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list services in namespace %s: %v", namespace, err)
	}

	if len(services.Items) == 0 {
		t.Logf("No services found in namespace %s", namespace)
		return
	}

	t.Logf("Found %d services in namespace %s\n", len(services.Items), namespace)

	for _, service := range services.Items {
		t.Logf("Service: %s/%s | Type: %s | ClusterIP: %s | Ports: %d",
			service.Namespace,
			service.Name,
			service.Spec.Type,
			service.Spec.ClusterIP,
			len(service.Spec.Ports))

		for _, port := range service.Spec.Ports {
			t.Logf("  - Port: %d | TargetPort: %v | Protocol: %s", port.Port, port.TargetPort, port.Protocol)
		}
	}
}

// Helper function to check if a pod is ready
func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
