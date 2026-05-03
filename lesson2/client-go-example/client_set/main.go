package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"

	"client-go-example/utils"
)

// clientSet: 内置对象client
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// 创建clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	utils.CheckErr(err)

	// 打印kube-system命名空间pod
	printKubeSystemPods(clientSet)

	// deploy增删改查
	crudDeploy(clientSet)
}

func printKubeSystemPods(cs *kubernetes.Clientset) {
	podList, err := cs.CoreV1().Pods("kube-system").List(context.Background(), metav1.ListOptions{})
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

func crudDeploy(clientSet *kubernetes.Clientset) {
	// create deploy
	log.Println("create deploy start")
	ctx := context.Background()
	namespace, name, appKey := "default", fmt.Sprintf("nginx-%s", time.Now().Format(time.DateOnly)), "app"
	deploy := newDeployment(namespace, name, appKey)
	createDeploy, err := clientSet.AppsV1().Deployments(namespace).Create(ctx, deploy, metav1.CreateOptions{})
	utils.CheckErr(err)
	log.Println("create deploy end")

	// update deploy
	log.Println("update deploy start")
	updateDeploy := createDeploy.DeepCopy()
	updateDeploy.Spec.Replicas = pointer.Int32(2)
	updateDeploy, err = clientSet.AppsV1().Deployments(namespace).Update(ctx, updateDeploy, metav1.UpdateOptions{})
	utils.CheckErr(err)
	log.Println("update deploy end")

	// update deploy with patch
	log.Println("update deploy with patch start")
	patchBody := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]string{
				"app": name,
			},
		},
	}
	patchBytes, err := json.Marshal(patchBody)
	utils.CheckErr(err)
	updateDeploy, err = clientSet.AppsV1().Deployments(namespace).Patch(ctx, name, types.MergePatchType, patchBytes, metav1.PatchOptions{})
	utils.CheckErr(err)
	log.Println("update deploy with patch end")

	// query deploy
	log.Println("query deploy start")
	getDeploy, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	utils.CheckErr(err)
	log.Printf("query deploy end, body:%s/%s\n", getDeploy.Namespace, getDeploy.Name)

	time.Sleep(time.Second * 10)

	// delete deploy
	log.Println("delete deploy start")
	err = clientSet.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	utils.CheckErr(err)
	log.Println("delete deploy end")

	// Output:
	// 2026/05/03 20:14:24 create deploy start
	// 2026/05/03 20:14:24 create deploy end
	// 2026/05/03 20:14:24 update deploy start
	// 2026/05/03 20:14:24 update deploy end
	// 2026/05/03 20:14:24 update deploy with patch start
	// 2026/05/03 20:14:24 update deploy with patch end
	// 2026/05/03 20:14:24 query deploy start
	// 2026/05/03 20:14:24 query deploy end, body:default/nginx-2026-05-03
	// 2026/05/03 20:14:34 delete deploy start
	// 2026/05/03 20:14:34 delete deploy end
}

func newDeployment(namespace, name, appKey string) *appsv1.Deployment {
	deployLabels := map[string]string{
		appKey: name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    deployLabels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: deployLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: deployLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:alpine",
						},
					},
				},
			},
		},
	}
}
