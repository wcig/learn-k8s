package main

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

func main() {
	// create config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		inClusterConfig, err2 := rest.InClusterConfig()
		if err2 != nil {
			klog.Fatal(err2)
		}
		config = inClusterConfig
	}

	// create clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	// create controller
	factory := informers.NewSharedInformerFactory(clientSet, time.Minute*30)
	serviceInformer := factory.Core().V1().Services().Informer()
	serviceLister := factory.Core().V1().Services().Lister()
	ingressInformer := factory.Networking().V1().Ingresses().Informer()
	ingressLister := factory.Networking().V1().Ingresses().Lister()
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	controller, err := NewController(clientSet, serviceInformer, ingressInformer, serviceLister, ingressLister, queue)
	if err != nil {
		klog.Fatal(err)
	}

	// start controller
	stop := make(chan struct{})
	defer close(stop)
	factory.Start(stop)
	factory.WaitForCacheSync(stop)
	controller.Run(stop)
}
