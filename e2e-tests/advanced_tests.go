package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestFrontendServiceHealth tests if the frontend service and its pods are healthy
func TestFrontendServiceHealth(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	namespace := testNamespace

	// Get the frontend service
	service, err := clientset.CoreV1().Services(namespace).Get(ctx, "frontend", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("Failed to get frontend service: %v", err)
	}

	t.Logf("✓ Frontend service found | ClusterIP: %s", service.Spec.ClusterIP)

	// Get pods for the frontend service
	selector := metav1.FormatLabelSelector(metav1.SetAsLabelSelector(service.Spec.Selector))
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		t.Fatalf("Failed to list frontend pods: %v", err)
	}

	if len(pods.Items) == 0 {
		t.Fatal("No frontend pods found")
	}

	t.Logf("Found %d frontend pod(s)\n", len(pods.Items))

	allReady := true
	for _, pod := range pods.Items {
		ready := isPodReady(&pod)
		t.Logf("Pod: %s | Ready: %v | Phase: %s", pod.Name, ready, pod.Status.Phase)

		if !ready {
			allReady = false
			// Print container statuses for debugging
			for _, containerStatus := range pod.Status.ContainerStatuses {
				t.Logf("  Container: %s | Ready: %v | RestartCount: %d",
					containerStatus.Name, containerStatus.Ready, containerStatus.RestartCount)
			}
		}
	}

	if !allReady {
		t.Error("Not all frontend pods are ready")
	}
}

// TestIngressStatus tests if the Ingress is properly configured
func TestIngressStatus(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list ingresses: %v", err)
	}

	if len(ingresses.Items) == 0 {
		t.Logf("⚠ No ingresses found in namespace %s", namespace)
		return
	}

	t.Logf("Found %d ingress(es)\n", len(ingresses.Items))

	for _, ingress := range ingresses.Items {
		t.Logf("Ingress: %s/%s | Class: %s", ingress.Namespace, ingress.Name, *ingress.Spec.IngressClassName)

		if len(ingress.Status.LoadBalancer.Ingress) == 0 {
			t.Logf("  ⚠ Ingress has no assigned IP/Hostname yet")
		} else {
			for _, lb := range ingress.Status.LoadBalancer.Ingress {
				if lb.IP != "" {
					t.Logf("  ✓ IP: %s", lb.IP)
				}
				if lb.Hostname != "" {
					t.Logf("  ✓ Hostname: %s", lb.Hostname)
				}
			}
		}

		// Print TLS info
		if len(ingress.Spec.TLS) > 0 {
			t.Logf("  TLS configured for %d host(s)", len(ingress.Spec.TLS))
			for _, tlsConfig := range ingress.Spec.TLS {
				t.Logf("    Secret: %s | Hosts: %v", tlsConfig.SecretName, tlsConfig.Hosts)
			}
		}

		// Print rules
		for _, rule := range ingress.Spec.Rules {
			t.Logf("  Rule: %s", rule.Host)
			if rule.HTTP != nil {
				for _, path := range rule.HTTP.Paths {
					t.Logf("    Path: %s -> Service: %s:%d", path.Path, path.Backend.Service.Name, path.Backend.Service.Port.Number)
				}
			}
		}
	}
}

// TestPersistentVolumes tests persistent volume status
func TestPersistentVolumes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list PVs: %v", err)
	}

	if len(pvs.Items) == 0 {
		t.Logf("⚠ No persistent volumes found")
		return
	}

	t.Logf("Found %d PV(s)\n", len(pvs.Items))

	for _, pv := range pvs.Items {
		t.Logf("PV: %s | Status: %s | Capacity: %v", pv.Name, pv.Status.Phase, pv.Spec.Capacity)
	}

	// Check PVCs
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(testNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list PVCs: %v", err)
	}

	if len(pvcs.Items) > 0 {
		t.Logf("\nFound %d PVC(s) in namespace %s\n", len(pvcs.Items), testNamespace)
		for _, pvc := range pvcs.Items {
			t.Logf("PVC: %s | Status: %s | Volume: %s", pvc.Name, pvc.Status.Phase, pvc.Spec.VolumeName)
		}
	}
}

// TestCrashingPods checks for pods with high restart counts
func TestCrashingPods(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list pods: %v", err)
	}

	crashingPods := []string{}

	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.RestartCount > 3 {
				msg := fmt.Sprintf("%s/%s (container: %s) - %d restarts",
					pod.Namespace, pod.Name, containerStatus.Name, containerStatus.RestartCount)
				crashingPods = append(crashingPods, msg)
			}
		}
	}

	if len(crashingPods) > 0 {
		t.Logf("⚠ Pods with high restart counts:\n")
		for _, pod := range crashingPods {
			t.Logf("  - %s", pod)
		}
	} else {
		t.Log("✓ No pods with excessive restart counts")
	}
}

// TestEventLog checks for recent cluster events
func TestEventLog(t *testing.T) {
	ensureNamespaceExists(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	namespace := testNamespace
	events, err := clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	// Filter for recent events (last 5 minutes)
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	recentEvents := []corev1.Event{}

	for _, event := range events.Items {
		if event.LastTimestamp.After(fiveMinutesAgo) {
			recentEvents = append(recentEvents, event)
		}
	}

	if len(recentEvents) == 0 {
		t.Log("✓ No recent events in the last 5 minutes")
		return
	}

	t.Logf("Found %d recent events:\n", len(recentEvents))
	for _, event := range recentEvents {
		t.Logf("  %s: %s | %s/%s | Count: %d", event.Type, event.Reason, event.Namespace, event.InvolvedObject.Name, event.Count)
	}
}
