package main

import (
	"flag"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	externalclient "code-generator/pkg/generated/clientset/versioned"
	externalinformers "code-generator/pkg/generated/informers/externalversions"
	externallister "code-generator/pkg/generated/listers/samplecontroller/v1alpha1"
)

type Controller struct {
	indexer  externallister.FooLister
	queue    workqueue.TypedRateLimitingInterface[string]
	informer cache.Controller
}

func NewController(queue workqueue.TypedRateLimitingInterface[string], indexer externallister.FooLister, informer cache.Controller) *Controller {
	return &Controller{
		informer: informer,
		indexer:  indexer,
		queue:    queue,
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	err := c.syncToStdout(key)
	c.handleErr(err, key)
	return true
}

func (c *Controller) syncToStdout(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		klog.Errorf("invalid resource key: %s", key)
		return nil
	}

	foo, err := c.indexer.Foos(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("Foo %s does not exist anymore", key)
			return nil
		}
		klog.Errorf("Fetching foo with key %s from store failed with %v", key, err)
		return err
	}
	fmt.Printf("Sync/Add/Update for foo %s/%s\n", foo.GetNamespace(), foo.GetName())
	return nil
}

func (c *Controller) handleErr(err error, key string) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("Dropping foo %q out of the queue: %v", key, err)
}

func (c *Controller) Run(workers int, stopCh chan struct{}) {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	klog.Info("Starting Foo controller")
	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Info("Stopping Pod controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

// go run main.go -kubeconfig=/Users/yangbo/.kube/config
func main() {
	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		klog.Fatal(err)
	}

	clientset, err := externalclient.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	factory := externalinformers.NewSharedInformerFactory(clientset, time.Second*30)
	informer := factory.Samplecontroller().V1alpha1().Foos()
	queue := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]())
	_, err = informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})
	if err != nil {
		klog.Fatal(err)
	}

	controller := NewController(queue, informer.Lister(), informer.Informer())
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)
	select {}
}
