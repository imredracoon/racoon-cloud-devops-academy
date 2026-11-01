package main

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestProblemi(t *testing.T) {
	casi := []struct {
		nome string
		pod  corev1.Pod
		n    int
	}{
		{"sano", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{{Name: "api", Ready: true}}}}, 0},
		{"crashloop", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodRunning,
			ContainerStatuses: []corev1.ContainerStatus{{Name: "api", RestartCount: 7,
				State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "CrashLoopBackOff"}}}}}}, 3},
		{"pending", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodPending}}, 1},
		{"in creazione", corev1.Pod{Status: corev1.PodStatus{Phase: corev1.PodPending,
			ContainerStatuses: []corev1.ContainerStatus{{Name: "api",
				State: corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{Reason: "ContainerCreating"}}}}}}, 1},
	}
	for _, tc := range casi {
		t.Run(tc.nome, func(t *testing.T) {
			if got := problemi(tc.pod, 5); len(got) != tc.n {
				t.Errorf("problemi = %v, ne voglio %d", got, tc.n)
			}
		})
	}
}
