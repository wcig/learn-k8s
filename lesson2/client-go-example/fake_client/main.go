package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"client-go-example/utils"
)

// fake包: 主要用于mock测试
func main() {
	const (
		pod1Name = "test-pod-1"
		pod2Name = "test-pod-2"
		ns       = "default"
	)

	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod1Name,
			Namespace: ns,
		},
	}
	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod2Name,
			Namespace: ns,
		},
	}

	// create fake client with pod1
	cs := fake.NewClientset(pod1)

	// create pod2
	_, err := cs.CoreV1().Pods(ns).Create(context.Background(), pod2, metav1.CreateOptions{})
	utils.CheckErr(err)

	// list pod
	podList, err := cs.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{})
	utils.CheckErr(err)

	// print pod list
	for _, d := range podList.Items {
		fmt.Printf("namespace: %v, name: %v\n", d.Namespace, d.Name)
	}

	// Output:
	// namespace: default, name: test-pod-1
	// namespace: default, name: test-pod-2
}
