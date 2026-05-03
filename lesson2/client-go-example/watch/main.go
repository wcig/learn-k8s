package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"

	"client-go-example/utils"
)

// clientSet watch: 监控资源变化, 相当于以下操作
// 1) kubectl proxy --port=8080 &
// 2) curl 'http://localhost:8080/apis/apps/v1/namespaces/default/deployments?watch=true'
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// 创建clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	utils.CheckErr(err)

	// 构造deploy
	go createUpdateDeleteDeploy(clientSet)

	time.Sleep(time.Minute)

	// watch deploy
	// watchDeploy(clientSet)
}

func createUpdateDeleteDeploy(cs *kubernetes.Clientset) {
	deploymentsClient := cs.AppsV1().Deployments(corev1.NamespaceDefault)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx-deployment",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx-deployment",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx:alpine",
							Env: []corev1.EnvVar{
								{
									Name:  "PAAS_APP_NAME",
									Value: "nginx-deployment",
								},
								{
									Name:  "PAAS_NAMESPACE",
									Value: "default",
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("512Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	// create deployment
	time.Sleep(10 * time.Second)
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	utils.CheckErr(err)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	// modify deployment replicas
	time.Sleep(10 * time.Second)
	fmt.Println("Updating deployment...")
	deployment.Spec.Replicas = pointer.Int32(2)
	result, err = deploymentsClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
	utils.CheckErr(err)
	fmt.Printf("Updated deployment %q.\n", result.GetObjectMeta().GetName())

	// delete deployment
	time.Sleep(10 * time.Second)
	fmt.Println("Deleting deployment...")
	err = deploymentsClient.Delete(context.Background(), deployment.GetObjectMeta().GetName(), metav1.DeleteOptions{})
	utils.CheckErr(err)
	fmt.Printf("Deleted deployment %q.\n", deployment.GetObjectMeta().GetName())
}

func watchDeploy(cs *kubernetes.Clientset) {
	watch, err := cs.AppsV1().Deployments("default").Watch(context.Background(), metav1.ListOptions{})
	utils.CheckErr(err)

	fmt.Println(">> start watch deployments")
	for event := range watch.ResultChan() {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if ok {
			fmt.Printf(">> watch event type: %8s, deployment name: %s, replicas: %d, generation: %d, resourceVersion: %s\n",
				event.Type, deployment.GetObjectMeta().GetName(), *deployment.Spec.Replicas,
				deployment.GetObjectMeta().GetGeneration(), deployment.GetObjectMeta().GetResourceVersion())
		}
	}

	// Output:
	// >> start watch deployments
	// >> watch event type:    ADDED, deployment name: nginx, replicas: 3, generation: 1, resourceVersion: 759
	// Creating deployment...
	// Created deployment "nginx-deployment".
	// >> watch event type:    ADDED, deployment name: nginx-deployment, replicas: 1, generation: 1, resourceVersion: 51565
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 1, generation: 1, resourceVersion: 51567
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 1, generation: 1, resourceVersion: 51571
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 1, generation: 1, resourceVersion: 51577
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 1, generation: 1, resourceVersion: 51585
	// Updating deployment...
	// Updated deployment "nginx-deployment".
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51600
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51601
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51606
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51612
	// >> watch event type: MODIFIED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51620
	// Deleting deployment...
	// Deleted deployment "nginx-deployment".
	// >> watch event type:  DELETED, deployment name: nginx-deployment, replicas: 2, generation: 2, resourceVersion: 51637
}
