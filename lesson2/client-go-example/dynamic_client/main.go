package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"client-go-example/utils"
)

// dynamicClient: 任意自定义资源client
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// 创建dynamicClient
	dynamicClient, err := dynamic.NewForConfig(config)
	utils.CheckErr(err)

	// 打印kube-system命名空间pod
	printKubeSystemPods(dynamicClient)
}

func printKubeSystemPods(dc *dynamic.DynamicClient) {
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}
	unstructuredList, err := dc.Resource(podGVR).
		Namespace("kube-system").
		List(context.Background(), metav1.ListOptions{})
	utils.CheckErr(err)

	podList := &corev1.PodList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredList.UnstructuredContent(), podList)
	utils.CheckErr(err)

	for _, d := range podList.Items {
		fmt.Printf("namespace: %v, name: %v, status: %v\n", d.Namespace, d.Name, d.Status.Phase)
	}

	// Output:
	// namespace: kube-system, name: coredns-7c65d6cfc9-974kp, status: Running
	// namespace: kube-system, name: coredns-7c65d6cfc9-bmzfd, status: Running
	// namespace: kube-system, name: etcd-1c2w-control-plane, status: Running
	// namespace: kube-system, name: kindnet-9xjqd, status: Running
	// namespace: kube-system, name: kindnet-d2n24, status: Running
	// namespace: kube-system, name: kindnet-nrhnz, status: Running
	// namespace: kube-system, name: kube-apiserver-1c2w-control-plane, status: Running
	// namespace: kube-system, name: kube-controller-manager-1c2w-control-plane, status: Running
	// namespace: kube-system, name: kube-proxy-6rj7b, status: Running
	// namespace: kube-system, name: kube-proxy-k9hgj, status: Running
	// namespace: kube-system, name: kube-proxy-n8zj4, status: Running
	// namespace: kube-system, name: kube-scheduler-1c2w-control-plane, status: Running
}
