package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"client-go-example/utils"
)

// SharedInformerFactory: informer + lister
func main() {
	// create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// create clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	utils.CheckErr(err)

	// create informer factory (重同步周期)
	informerFactory := informers.NewSharedInformerFactory(clientSet, time.Second*30)

	// get informer/lister
	deployInformer := informerFactory.Apps().V1().Deployments()
	informer := deployInformer.Informer()
	lister := deployInformer.Lister()

	// event handler
	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    addDeployment,
		UpdateFunc: updateDeployment,
		DeleteFunc: deleteDeployment,
	})
	utils.CheckErr(err)

	// watch cache sync
	stopCh := make(chan struct{})
	defer close(stopCh)
	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	// list deploy
	deployments, err := lister.Deployments("").List(labels.Everything())
	utils.CheckErr(err)
	for i, d := range deployments {
		fmt.Printf("lister: %d, namespace: %v, name: %v, replicas: %d\n", i, d.Namespace, d.Name, *d.Spec.Replicas)
	}

	// wait
	<-stopCh
}

func addDeployment(obj interface{}) {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return
	}
	fmt.Println("add deployment:", deploy.Name)
}

func updateDeployment(old, new interface{}) {
	oldDeploy, ok := old.(*v1.Deployment)
	if !ok {
		return
	}
	newDeploy, ok := new.(*v1.Deployment)
	if !ok {
		return
	}
	fmt.Println("update deployment:", oldDeploy.Name, newDeploy.Name)
}

func deleteDeployment(obj interface{}) {
	deploy, ok := obj.(*v1.Deployment)
	if !ok {
		return
	}
	fmt.Println("delete deployment:", deploy.Name)
}
