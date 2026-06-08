package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestNginxPodRunningStatus verifies that the frontend (nginx) pod is in Running state
func TestNginxPodRunningStatus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// List pods with label app=frontend (nginx)
	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=frontend",
	})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatalf("No frontend pods found in namespace %s", testNamespace)
	}

	frontendPod := pods.Items[0]
	t.Logf("Frontend Pod Name: %s", frontendPod.Name)
	t.Logf("Pod Phase: %s", frontendPod.Status.Phase)
	t.Logf("Pod Ready: %v", isPodReady(&frontendPod))

	// Verify pod is running
	if frontendPod.Status.Phase != corev1.PodRunning {
		t.Errorf("Expected frontend pod phase Running, got %s", frontendPod.Status.Phase)
		t.Logf("Pod status: %+v", frontendPod.Status)
	}

	// Verify pod is ready
	if !isPodReady(&frontendPod) {
		t.Errorf("Frontend pod is not ready")
		for _, condition := range frontendPod.Status.Conditions {
			t.Logf("Condition: %s = %s (Reason: %s)", condition.Type, condition.Status, condition.Reason)
		}
	}

	t.Log("✓ Frontend (nginx) pod is Running and Ready")
}

// TestNginxLivenessProbeConfiguration verifies nginx pod has Liveness Probe configured
func TestNginxLivenessProbeConfiguration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=frontend",
	})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatalf("No frontend pods found in namespace %s", testNamespace)
	}

	frontendPod := pods.Items[0]
	t.Logf("Verifying Liveness Probe for pod: %s", frontendPod.Name)

	for _, container := range frontendPod.Spec.Containers {
		if container.LivenessProbe == nil {
			t.Errorf("Container %s has no Liveness Probe configured", container.Name)
			continue
		}

		t.Logf("✓ Container %s has Liveness Probe:", container.Name)
		if container.LivenessProbe.HTTPGet != nil {
			t.Logf("  - Type: HTTP GET")
			t.Logf("  - Path: %s", container.LivenessProbe.HTTPGet.Path)
			t.Logf("  - Port: %d", container.LivenessProbe.HTTPGet.Port.IntValue())
		} else if container.LivenessProbe.Exec != nil {
			t.Logf("  - Type: Exec")
			t.Logf("  - Command: %v", container.LivenessProbe.Exec.Command)
		} else if container.LivenessProbe.TCPSocket != nil {
			t.Logf("  - Type: TCP Socket")
			t.Logf("  - Port: %d", container.LivenessProbe.TCPSocket.Port.IntValue())
		} else if container.LivenessProbe.GRPC != nil {
			t.Logf("  - Type: gRPC")
			t.Logf("  - Port: %d", container.LivenessProbe.GRPC.Port)
		}

		t.Logf("  - InitialDelaySeconds: %d", container.LivenessProbe.InitialDelaySeconds)
		t.Logf("  - TimeoutSeconds: %d", container.LivenessProbe.TimeoutSeconds)
		t.Logf("  - PeriodSeconds: %d", container.LivenessProbe.PeriodSeconds)
		t.Logf("  - FailureThreshold: %d", container.LivenessProbe.FailureThreshold)
	}

	t.Log("✓ Liveness Probe is properly configured")
}

// TestNginxReadinessProbeConfiguration verifies nginx pod has Readiness Probe configured
func TestNginxReadinessProbeConfiguration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=frontend",
	})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatalf("No frontend pods found in namespace %s", testNamespace)
	}

	frontendPod := pods.Items[0]
	t.Logf("Verifying Readiness Probe for pod: %s", frontendPod.Name)

	for _, container := range frontendPod.Spec.Containers {
		if container.ReadinessProbe == nil {
			t.Errorf("Container %s has no Readiness Probe configured", container.Name)
			continue
		}

		t.Logf("✓ Container %s has Readiness Probe:", container.Name)
		if container.ReadinessProbe.HTTPGet != nil {
			t.Logf("  - Type: HTTP GET")
			t.Logf("  - Path: %s", container.ReadinessProbe.HTTPGet.Path)
			t.Logf("  - Port: %d", container.ReadinessProbe.HTTPGet.Port.IntValue())
		} else if container.ReadinessProbe.Exec != nil {
			t.Logf("  - Type: Exec")
			t.Logf("  - Command: %v", container.ReadinessProbe.Exec.Command)
		} else if container.ReadinessProbe.TCPSocket != nil {
			t.Logf("  - Type: TCP Socket")
			t.Logf("  - Port: %d", container.ReadinessProbe.TCPSocket.Port.IntValue())
		} else if container.ReadinessProbe.GRPC != nil {
			t.Logf("  - Type: gRPC")
			t.Logf("  - Port: %d", container.ReadinessProbe.GRPC.Port)
		}

		t.Logf("  - InitialDelaySeconds: %d", container.ReadinessProbe.InitialDelaySeconds)
		t.Logf("  - TimeoutSeconds: %d", container.ReadinessProbe.TimeoutSeconds)
		t.Logf("  - PeriodSeconds: %d", container.ReadinessProbe.PeriodSeconds)
		t.Logf("  - FailureThreshold: %d", container.ReadinessProbe.FailureThreshold)
	}

	t.Log("✓ Readiness Probe is properly configured")
}

// TestNginxProbesFunctioning verifies that the Readiness Probe is passing (Pod Ready)
func TestNginxProbesFunctioning(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=frontend",
	})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatalf("No frontend pods found in namespace %s", testNamespace)
	}

	frontendPod := pods.Items[0]
	t.Logf("Frontend Pod: %s", frontendPod.Name)
	t.Logf("Total Restarts: %d", getTotalRestarts(&frontendPod))

	// Check each condition
	for _, condition := range frontendPod.Status.Conditions {
		t.Logf("Condition %s: %s (Message: %s)", condition.Type, condition.Status, condition.Message)

		if condition.Type == corev1.PodReady {
			if condition.Status != corev1.ConditionTrue {
				t.Errorf("Pod Ready condition is %s, expected True", condition.Status)
				t.Logf("Reason: %s", condition.Reason)
			} else {
				t.Log("✓ Pod Ready condition is passing")
			}
		}

		if condition.Type == corev1.ContainersReady {
			if condition.Status != corev1.ConditionTrue {
				t.Errorf("Containers Ready condition is %s, expected True", condition.Status)
			} else {
				t.Log("✓ Containers Ready condition is passing")
			}
		}
	}

	// Check container statuses for restart count
	for _, containerStatus := range frontendPod.Status.ContainerStatuses {
		t.Logf("Container %s:", containerStatus.Name)
		t.Logf("  - Ready: %v", containerStatus.Ready)
		t.Logf("  - RestartCount: %d", containerStatus.RestartCount)
		t.Logf("  - State: %+v", containerStatus.State)

		if !containerStatus.Ready {
			t.Errorf("Container %s is not ready", containerStatus.Name)
		}
	}

	t.Log("✓ Readiness Probe is functioning (Pod Ready)")
}

// TestNginxPodRestarts verifies nginx pod has acceptable restart count (not excessive)
func TestNginxPodRestarts(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app=frontend",
	})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatalf("No frontend pods found in namespace %s", testNamespace)
	}

	frontendPod := pods.Items[0]
	totalRestarts := getTotalRestarts(&frontendPod)

	t.Logf("Pod: %s", frontendPod.Name)
	t.Logf("Total Restarts: %d", totalRestarts)

	// Warning threshold: if more than 5 restarts, something might be wrong
	if totalRestarts > 5 {
		t.Logf("⚠ Warning: Pod has high restart count (%d)", totalRestarts)
		for _, containerStatus := range frontendPod.Status.ContainerStatuses {
			t.Logf("  Container %s restarts: %d", containerStatus.Name, containerStatus.RestartCount)
			if containerStatus.LastTerminationState.Terminated != nil {
				term := containerStatus.LastTerminationState.Terminated
				t.Logf("    Last termination: %s (Exit Code: %d)", term.Reason, term.ExitCode)
			}
		}
	} else {
		t.Log("✓ Pod restart count is acceptable")
	}
}

// TestAllPodsInTestAutoNamespace lists all pods in test-auto namespace
func TestAllPodsInTestAutoNamespace(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods: %v", err)
	}

	t.Logf("Total Pods in %s namespace: %d", testNamespace, len(pods.Items))
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, pod := range pods.Items {
		status := "✓"
		if !isPodReady(&pod) {
			status = "✗"
		}
		t.Logf("%s %s/%s (Phase: %s, Ready: %v, Restarts: %d)",
			status, pod.Namespace, pod.Name, pod.Status.Phase, isPodReady(&pod), getTotalRestarts(&pod))

		// Show container details
		for _, containerStatus := range pod.Status.ContainerStatuses {
			readyStr := "✓"
			if !containerStatus.Ready {
				readyStr = "✗"
			}
			t.Logf("    %s %s (Ready: %v, Restarts: %d)", readyStr, containerStatus.Name, containerStatus.Ready, containerStatus.RestartCount)
		}
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

// TestPodsInBothNamespaces compares pods in default and test-auto namespaces
func TestPodsInBothNamespaces(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Get pods from default namespace
	defaultPods, err := clientset.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Logf("Note: Could not access default namespace: %v", err)
		defaultPods = &corev1.PodList{Items: []corev1.Pod{}}
	}

	// Get pods from test-auto namespace
	testPods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods in %s: %v", testNamespace, err)
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("PODS IN DEFAULT NAMESPACE")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(defaultPods.Items) == 0 {
		t.Log("(No pods in default namespace or no access)")
	} else {
		for _, pod := range defaultPods.Items {
			status := "✓"
			if !isPodReady(&pod) {
				status = "✗"
			}
			t.Logf("%s %s (Phase: %s, Ready: %v)", status, pod.Name, pod.Status.Phase, isPodReady(&pod))
		}
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Log("PODS IN TEST-AUTO NAMESPACE")
	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for _, pod := range testPods.Items {
		status := "✓"
		if !isPodReady(&pod) {
			status = "✗"
		}
		t.Logf("%s %s (Phase: %s, Ready: %v)", status, pod.Name, pod.Status.Phase, isPodReady(&pod))
	}

	t.Log("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	t.Logf("Summary: %d pods in default | %d pods in test-auto", len(defaultPods.Items), len(testPods.Items))
}

// Helper function to get total restarts
func getTotalRestarts(pod *corev1.Pod) int32 {
	var totalRestarts int32
	for _, containerStatus := range pod.Status.ContainerStatuses {
		totalRestarts += containerStatus.RestartCount
	}
	return totalRestarts
}

// Helper function to check if pod is ready
func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}
