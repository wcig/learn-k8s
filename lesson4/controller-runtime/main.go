package main

import (
	"context"
	"os"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"controller-runtime/pkg/apis/samplecontroller/v1alpha1"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

type reconciler struct {
	client.Client
	scheme *runtime.Scheme
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("foo", req.NamespacedName)
	log.V(1).Info("reconciling foo")

	var foo v1alpha1.Foo
	if err := r.Get(ctx, req.NamespacedName, &foo); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Foo does not exist anymore")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Foo reconcile err")
		return ctrl.Result{}, err
	}

	log.Info("Sync/Add/Update for foo")
	return ctrl.Result{}, nil
}

func main() {
	ctrl.SetLogger(zap.New())

	// 1. 确保目录存在
	certDir := "/tmp/k8s-webhook-server/serving-certs"
	if err := os.MkdirAll(certDir, 0700); err != nil {
		setupLog.Error(err, "unable to create cert dir")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// in a real controller, we'd create a new scheme for this
	err = v1alpha1.AddToScheme(mgr.GetScheme())
	if err != nil {
		setupLog.Error(err, "unable to add scheme")
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Foo{}).
		Complete(&reconciler{
			Client: mgr.GetClient(),
			scheme: mgr.GetScheme(),
		})
	if err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	err = ctrl.NewWebhookManagedBy(mgr).
		For(&v1alpha1.Foo{}).
		// TODO: 证书问题需要解决
		// WithDefaulter(&v1alpha1.FooAnnotator{}).
		// WithValidator(&v1alpha1.FooValidator{}).
		Complete()
	if err != nil {
		setupLog.Error(err, "unable to create webhook")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
