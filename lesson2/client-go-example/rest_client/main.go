package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"client-go-example/utils"
)

// restClient: 基础client
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// 设置API路径 (这里查询的pod为无组名资源组使用/api而不是/apis)
	config.APIPath = "api"

	// 设置资源组和版本, 对应GVR中的GV
	config.GroupVersion = &corev1.SchemeGroupVersion

	// 设置编解码器
	config.NegotiatedSerializer = scheme.Codecs

	// 初始化Client
	restClient, err := rest.RESTClientFor(config)
	utils.CheckErr(err)

	// 构造接收pod列表对象
	result := &corev1.PodList{}

	// 查询pod列表
	err = restClient.Get().
		Namespace("kube-system").                                                // 命名空间
		Resource("pods").                                                        // 资源对象, 对应GVR的R
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec). // 参数及序列化工具
		Do(context.Background()).                                                // 发送请求
		Into(result)                                                             // 写入返回值
	utils.CheckErr(err)

	// 列出pod列表
	for _, d := range result.Items {
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
