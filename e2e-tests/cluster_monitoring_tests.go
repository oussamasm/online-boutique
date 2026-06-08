package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestClusterStateSnapshot provides a comprehensive snapshot of cluster health
func TestClusterStateSnapshot(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║              CLUSTER HEALTH CHECK SNAPSHOT                        ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")

	// 1. API Server Accessibility
	t.Log("\n[1/4] API Server Accessibility Check")
	t.Log("──────────────────────────────────────────────────────────────────")
	version, err := clientset.Discovery().ServerVersion()
	if err != nil {
		t.Errorf("✗ API Server is NOT accessible: %v", err)
	} else {
		t.Logf("✓ API Server is accessible")
		t.Logf("  - Version: %s", version.GitVersion)
		t.Logf("  - Platform: %s/%s", version.Platform, version.GoVersion)
	}

	// 2. Node Status
	t.Log("\n[2/4] Node Status Check")
	t.Log("──────────────────────────────────────────────────────────────────")
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("✗ Failed to list nodes: %v", err)
	}

	readyCount := 0
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				status := "✗"
				if condition.Status == corev1.ConditionTrue {
					status = "✓"
					readyCount++
				}
				t.Logf("%s Node: %s (Ready: %s)", status, node.Name, condition.Status)
				break
			}
		}
	}
	t.Logf("  Total Nodes: %d | Ready: %d", len(nodes.Items), readyCount)

	// 3. Namespace Status
	t.Log("\n[3/4] Namespace Status Check")
	t.Log("──────────────────────────────────────────────────────────────────")
	namespaces, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("✗ Failed to list namespaces: %v", err)
	}

	t.Logf("Total Namespaces: %d", len(namespaces.Items))
	for _, ns := range namespaces.Items {
		if ns.Name == "default" || ns.Name == testNamespace {
			t.Logf("  - %s (Status: %s)", ns.Name, ns.Status.Phase)
		}
	}

	// 4. Pod Summary
	t.Log("\n[4/4] Pod Summary")
	t.Log("──────────────────────────────────────────────────────────────────")
	allPods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("✗ Failed to list pods: %v", err)
	}

	runningCount := 0
	failedCount := 0
	pendingCount := 0

	for _, pod := range allPods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			runningCount++
		case corev1.PodFailed:
			failedCount++
		case corev1.PodPending:
			pendingCount++
		}
	}

	t.Logf("Total Pods: %d", len(allPods.Items))
	t.Logf("  - Running: %d ✓", runningCount)
	t.Logf("  - Pending: %d ⏳", pendingCount)
	t.Logf("  - Failed: %d ✗", failedCount)

	t.Log("\n╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║                    SNAPSHOT COMPLETE                             ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")
}

// TestPodLifecycleMonitoring monitors pod creation, readiness, and health
func TestPodLifecycleMonitoring(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║              POD LIFECYCLE MONITORING                             ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods: %v", err)
	}

	t.Logf("\nMonitoring %d pods in namespace: %s\n", len(pods.Items), testNamespace)

	for _, pod := range pods.Items {
		t.Logf("Pod: %s", pod.Name)
		t.Logf("  Phase: %s", pod.Status.Phase)
		t.Logf("  QOS Class: %s", pod.Status.QOSClass)
		t.Logf("  Created: %v", pod.CreationTimestamp.Format("2006-01-02 15:04:05"))

		// Show age
		age := time.Since(pod.CreationTimestamp.Time).Seconds()
		t.Logf("  Age: %.0f seconds", age)

		// Check pod conditions
		t.Logf("  Conditions:")
		for _, condition := range pod.Status.Conditions {
			status := "✓"
			if condition.Status != corev1.ConditionTrue {
				status = "✗"
			}
			t.Logf("    %s %s: %s (Reason: %s)", status, condition.Type, condition.Status, condition.Reason)
		}

		// Check container status
		t.Logf("  Containers:")
		for _, containerStatus := range pod.Status.ContainerStatuses {
			readyStr := "✓"
			if !containerStatus.Ready {
				readyStr = "✗"
			}
			t.Logf("    %s %s (Ready: %v, Restarts: %d)", readyStr, containerStatus.Name, containerStatus.Ready, containerStatus.RestartCount)

			if containerStatus.State.Running != nil {
				t.Logf("      State: Running (Started: %v)", containerStatus.State.Running.StartedAt.Format("2006-01-02 15:04:05"))
			} else if containerStatus.State.Waiting != nil {
				t.Logf("      State: Waiting (Reason: %s)", containerStatus.State.Waiting.Reason)
			} else if containerStatus.State.Terminated != nil {
				t.Logf("      State: Terminated (ExitCode: %d, Reason: %s)", containerStatus.State.Terminated.ExitCode, containerStatus.State.Terminated.Reason)
			}
		}
		t.Log()
	}
}

// TestPodHealthComparison compares pod health metrics across deployments
func TestPodHealthComparison(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║              POD HEALTH COMPARISON                                ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")

	// Get test-auto pods
	testPods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods in test-auto: %v", err)
	}

	type PodMetrics struct {
		Name          string
		Phase         corev1.PodPhase
		Ready         bool
		Restarts      int32
		Age           int64
		ContainerReady map[string]bool
	}

	metrics := make([]PodMetrics, 0)

	for _, pod := range testPods.Items {
		pm := PodMetrics{
			Name:           pod.Name,
			Phase:          pod.Status.Phase,
			Ready:          isPodReady(&pod),
			Restarts:       getTotalRestarts(&pod),
			Age:            int64(time.Since(pod.CreationTimestamp.Time).Seconds()),
			ContainerReady: make(map[string]bool),
		}

		for _, cs := range pod.Status.ContainerStatuses {
			pm.ContainerReady[cs.Name] = cs.Ready
		}

		metrics = append(metrics, pm)
	}

	// Display table
	t.Log("\n┌─────────────────────────┬──────────┬───────┬──────────┬─────────┐")
	t.Log("│ Pod Name                │ Phase    │ Ready │ Restarts │ Age (s) │")
	t.Log("├─────────────────────────┼──────────┼───────┼──────────┼─────────┤")

	for _, pm := range metrics {
		readyStr := "✓"
		if !pm.Ready {
			readyStr := "✗"
			_ = readyStr
		}
		readyStatus := "Yes"
		if !pm.Ready {
			readyStatus = "No"
		}
		phase := string(pm.Phase)
		t.Logf("│ %-23s │ %-8s │ %s    │ %-8d │ %-7d │", pm.Name, phase, readyStatus, pm.Restarts, pm.Age)
	}

	t.Log("└─────────────────────────┴──────────┴───────┴──────────┴─────────┘")

	// Summary
	t.Log("\nSummary:")
	readyCount := 0
	for _, pm := range metrics {
		if pm.Ready {
			readyCount++
		}
	}
	t.Logf("  - Total Pods: %d", len(metrics))
	t.Logf("  - Ready Pods: %d ✓", readyCount)
	t.Logf("  - Not Ready: %d ✗", len(metrics)-readyCount)

	highRestarts := 0
	for _, pm := range metrics {
		if pm.Restarts > 3 {
			highRestarts++
			t.Logf("  - Pod %s has %d restarts (⚠ High)", pm.Name, pm.Restarts)
		}
	}
	if highRestarts == 0 {
		t.Log("  - No pods with excessive restarts ✓")
	}
}

// TestDetailedProbeAnalysis provides detailed analysis of all probes
func TestDetailedProbeAnalysis(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Log("╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║              DETAILED PROBE ANALYSIS                              ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods: %v", err)
	}

	podsWithoutProbes := make([]string, 0)
	probeCount := 0

	for _, pod := range pods.Items {
		t.Logf("\nPod: %s", pod.Name)
		t.Logf("──────────────────────────────────────────────────────────────")

		hasLiveness := false
		hasReadiness := false

		for _, container := range pod.Spec.Containers {
			t.Logf("Container: %s", container.Name)

			// Liveness Probe
			if container.LivenessProbe != nil {
				hasLiveness = true
				probeCount++
				t.Logf("  ✓ Liveness Probe: %s", describeProbe(container.LivenessProbe))
			} else {
				t.Logf("  ✗ Liveness Probe: NOT CONFIGURED")
			}

			// Readiness Probe
			if container.ReadinessProbe != nil {
				hasReadiness = true
				probeCount++
				t.Logf("  ✓ Readiness Probe: %s", describeProbe(container.ReadinessProbe))
			} else {
				t.Logf("  ✗ Readiness Probe: NOT CONFIGURED")
			}

			// Startup Probe
			if container.StartupProbe != nil {
				probeCount++
				t.Logf("  ✓ Startup Probe: %s", describeProbe(container.StartupProbe))
			} else {
				t.Logf("  - Startup Probe: Not configured (optional)")
			}
		}

		if !hasLiveness || !hasReadiness {
			podsWithoutProbes = append(podsWithoutProbes, fmt.Sprintf("%s (Missing: %s%s)", 
				pod.Name,
				map[bool]string{true: "", false: "Liveness "}[hasLiveness],
				map[bool]string{true: "", false: "Readiness"}[hasReadiness]))
		}
	}

	t.Log("\n╔═══════════════════════════════════════════════════════════════════╗")
	t.Logf("Total Probes Found: %d", probeCount)
	if len(podsWithoutProbes) > 0 {
		t.Log("\n⚠ Pods without required probes:")
		for _, podName := range podsWithoutProbes {
			t.Logf("  - %s", podName)
		}
	} else {
		t.Log("✓ All pods have required probes configured")
	}
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")
}

// Helper function to describe a probe
func describeProbe(probe *corev1.Probe) string {
	if probe == nil {
		return "None"
	}

	probeType := "Unknown"
	if probe.HTTPGet != nil {
		probeType = fmt.Sprintf("HTTP GET %s:%d%s", probe.HTTPGet.Scheme, probe.HTTPGet.Port.IntValue(), probe.HTTPGet.Path)
	} else if probe.Exec != nil {
		probeType = fmt.Sprintf("Exec %v", probe.Exec.Command)
	} else if probe.TCPSocket != nil {
		probeType = fmt.Sprintf("TCP Socket :%d", probe.TCPSocket.Port.IntValue())
	} else if probe.GRPC != nil {
		probeType = fmt.Sprintf("gRPC Port:%d", probe.GRPC.Port)
	}

	return fmt.Sprintf("%s (InitialDelay: %ds, Timeout: %ds, Period: %ds, FailureThreshold: %d)",
		probeType, probe.InitialDelaySeconds, probe.TimeoutSeconds, probe.PeriodSeconds, probe.FailureThreshold)
}

// TestPodCleanupSimulation simulates pod cleanup after tests
func TestPodCleanupSimulation(t *testing.T) {
	t.Log("╔═══════════════════════════════════════════════════════════════════╗")
	t.Log("║              POD CLEANUP SIMULATION                               ║")
	t.Log("╚═══════════════════════════════════════════════════════════════════╝")

	t.Log("\nThis is a SIMULATION - no actual pods will be deleted.\n")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods: %v", err)
	}

	t.Logf("Pods available for cleanup in namespace %s:", testNamespace)
	t.Log("──────────────────────────────────────────────────────────────")

	for i, pod := range pods.Items {
		t.Logf("%d. %s (Phase: %s, Age: %v)", i+1, pod.Name, pod.Status.Phase, time.Since(pod.CreationTimestamp.Time))
	}

	t.Log("\nCleanup Instructions:")
	t.Log("──────────────────────────────────────────────────────────────")
	t.Log("To clean up pods individually:")
	for _, pod := range pods.Items {
		t.Logf("  kubectl delete pod %s -n %s", pod.Name, testNamespace)
	}

	t.Log("\nTo clean up all pods in namespace:")
	t.Logf("  kubectl delete pods --all -n %s", testNamespace)

	t.Log("\nTo delete the entire namespace (including all pods):")
	t.Logf("  kubectl delete namespace %s", testNamespace)

	t.Log("\n✓ Cleanup simulation complete - no actual deletions performed")
}
