package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
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

	externalClientset, err := externalclient.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	factory := externalinformers.NewSharedInformerFactory(externalClientset, time.Second*30)
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

	// 无 leader election
	// go controller.Run(1, stop)
	// select {}

	// leader election
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	lock, err := getLock(clientset)
	if err != nil {
		klog.Fatal(err)
	}
	// 60秒租期 + 5秒重试 + 15秒续租超时：抢到锁才启动控制器，丢锁就停，保证集群里始终只有一个活跃控制器，故障时自动切换。
	leaderelection.RunOrDie(context.Background(), leaderelection.LeaderElectionConfig{
		Lock:          lock,
		LeaseDuration: 60 * time.Second, // 锁的最长有效期，抢到锁后60s内必须续租，否则视为宕机
		RenewDeadline: 15 * time.Second, // 续租超时时间，若15s内味完成续租，主动放弃锁（防止脑裂）
		RetryPeriod:   5 * time.Second,  // 重试间隔，每5s尝试抢锁/续租（退避jitter在0～5s）
		Callbacks: leaderelection.LeaderCallbacks{
			// 抢到锁，启动控制器
			OnStartedLeading: func(ctx context.Context) {
				klog.Info("leader election win")
				go controller.Run(1, stop)
			},
			// 抢到锁，停止控制器
			OnStoppedLeading: func() {
				klog.Info("leader election lost")
			},
			// 旁观模式，每当新Leader产生，其他副本打印日志（用于观测切换）
			OnNewLeader: func(identity string) {
				klog.Infof("leader election %s win", identity)
			},
		},
	})
	select {}
}

func getLock(client *kubernetes.Clientset) (resourcelock.Interface, error) {
	lockName := "code-generator"
	lockNamespace := "kube-system"
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	return resourcelock.New(
		resourcelock.LeasesResourceLock,
		lockNamespace,
		lockName,
		client.CoreV1(),
		client.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: hostname,
		},
	)
}
