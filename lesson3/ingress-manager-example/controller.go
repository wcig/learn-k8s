package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	networkLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
)

const (
	workerNum                          = 1
	maxRetry                           = 5
	annotationIngressHttpKey           = "ingress/http"             // true: enable ingress, false: disable ingress
	annotationIngressHttpPathPrefixKey = "ingress/http-path-prefix" // http path prefix
	host                               = "example.com"
)

type Controller struct {
	client          kubernetes.Interface
	serviceInformer cache.SharedIndexInformer
	ingressInformer cache.SharedIndexInformer
	serviceLister   coreLister.ServiceLister
	ingressLister   networkLister.IngressLister
	queue           workqueue.TypedRateLimitingInterface[string]
}

func NewController(client kubernetes.Interface, serviceInformer, ingressInformer cache.SharedIndexInformer, serviceLister coreLister.ServiceLister, ingressLister networkLister.IngressLister, queue workqueue.TypedRateLimitingInterface[string]) (*Controller, error) {
	c := &Controller{
		client:          client,
		serviceInformer: serviceInformer,
		ingressInformer: ingressInformer,
		serviceLister:   serviceLister,
		ingressLister:   ingressLister,
		queue:           queue,
	}

	_, err := serviceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.addService(obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			c.updateService(old, new)
		},
		// No need DeleteFunc, because the ownerReference of ingress is a service and will be deleted along with the service cascade.
	})
	if err != nil {
		return nil, fmt.Errorf("service informer add event handler err: %w", err)
	}

	_, err = ingressInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			c.deleteIngress(obj)
		},
	})
	if err != nil {
		return nil, fmt.Errorf("ingress informer add event handler err: %w", err)
	}

	return c, nil
}

func (c *Controller) Run(stop chan struct{}) {
	klog.Info("Controller start")
	for i := 0; i < workerNum; i++ {
		go wait.Until(c.runWorker, time.Second, stop)
	}
	<-stop
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.syncService(key)
	c.handleErr(err, key)
	return true
}

func (c *Controller) handleErr(err error, key string) {
	if err == nil {
		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < maxRetry {
		klog.Infof("Error sync service %q: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	runtime.HandleError(err)
	klog.Infof("Dropping service %q out of the queue: %v", key, err)
}

func (c *Controller) addService(obj interface{}) {
	c.enqueue(obj)
}

func (c *Controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		klog.Errorf("Enqueue obj err: %v", err)
		return
	}
	c.queue.Add(key)
}

func (c *Controller) updateService(old interface{}, new interface{}) {
	oldSvc, ok := old.(*corev1.Service)
	if !ok {
		klog.Errorf("Update service type assert old service failed")
		return
	}
	newSvc, ok := new.(*corev1.Service)
	if !ok {
		klog.Errorf("Update service type assert new service failed")
		return
	}
	oldVal, ok1 := oldSvc.GetAnnotations()[annotationIngressHttpKey]
	newVal, ok2 := newSvc.GetAnnotations()[annotationIngressHttpKey]
	if ok1 && ok2 && oldVal != newVal {
		c.enqueue(new)
	}
	if ok1 && !ok2 || !ok1 && ok2 {
		c.enqueue(new)
	}
}

func (c *Controller) deleteIngress(obj interface{}) {
	ingress, ok := obj.(*networkv1.Ingress)
	if !ok {
		klog.Errorf("Delete ingress type assert ingress failed")
		return
	}

	owner := metav1.GetControllerOf(ingress)
	if owner == nil || owner.Kind != "Service" {
		return
	}

	key := cache.NewObjectName(ingress.Namespace, owner.Name).String()
	c.queue.Add(key)
}

func (c *Controller) syncService(key string) error {
	klog.Infof("Sync service %q", key)
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return fmt.Errorf("sync service %q split key err: %w", key, err)
	}

	service, err := c.serviceLister.Services(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("sync service %q get service err: %w", key, err)
	}

	enableIngres := false
	if val, ok := service.GetAnnotations()[annotationIngressHttpKey]; ok {
		enableIngres, err = strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("sync service %q parse service annotation %q err: %w", key, val, err)
		}
	}
	httpPathPrefix, ok := service.GetAnnotations()[annotationIngressHttpPathPrefixKey]
	if !ok || !strings.HasPrefix(httpPathPrefix, "/") {
		httpPathPrefix = "/" + httpPathPrefix
	}

	_, err = c.ingressLister.Ingresses(namespace).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("sync service %q get ingress err: %w", key, err)
	}
	ingressExist := !errors.IsNotFound(err)

	if enableIngres && !ingressExist {
		klog.Infof("Sync service %q create ingress", key)
		newIngress := constructIngress(service, httpPathPrefix)
		_, err = c.client.NetworkingV1().Ingresses(namespace).Create(context.Background(), newIngress, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("sync service %q create ingress err: %w", key, err)
		}
	} else if !enableIngres && ingressExist {
		klog.Infof("Sync service %q delete ingress", key)
		err = c.client.NetworkingV1().Ingresses(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf("sync service %q delete ingress err: %w", key, err)
		}
	}
	return nil
}

func constructIngress(service *corev1.Service, httpPathPrefix string) *networkv1.Ingress {
	pathType := networkv1.PathTypePrefix
	port := service.Spec.Ports[0].Port
	ingress := &networkv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: service.Namespace,
			Name:      service.Name,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(service, corev1.SchemeGroupVersion.WithKind("Service")),
			},
		},
		Spec: networkv1.IngressSpec{
			IngressClassName: pointer.String("cloud-provider-kind"),
			Rules: []networkv1.IngressRule{
				{
					Host: host,
					IngressRuleValue: networkv1.IngressRuleValue{
						HTTP: &networkv1.HTTPIngressRuleValue{
							Paths: []networkv1.HTTPIngressPath{
								{
									Path:     httpPathPrefix,
									PathType: &pathType,
									Backend: networkv1.IngressBackend{
										Service: &networkv1.IngressServiceBackend{
											Name: service.Name,
											Port: networkv1.ServiceBackendPort{
												Number: port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}
