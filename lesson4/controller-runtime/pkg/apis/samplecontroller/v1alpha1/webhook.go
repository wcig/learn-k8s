package v1alpha1

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// 标记接口（可选，用于代码即文档）
var _ admission.CustomValidator = &FooValidator{}
var _ admission.CustomDefaulter = &FooAnnotator{}

// +kubebuilder:webhook:path=/mutate-samplecontroller-k8s-io-v1alpha1-foo,mutating=true,failurePolicy=fail,sideEffects=None,groups=samplecontroller.k8s.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.samplecontroller.k8s.io,admissionReviewVersions=v1

type FooAnnotator struct{}

// Default 在对象被apiserver持久化前调用（CREATE/UPDATE）
func (a *FooAnnotator) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	log.Info("Annotated Foo")

	foo, ok := obj.(*Foo)
	if !ok {
		return fmt.Errorf("expected a Foo but got a %T", obj)
	}

	if foo.Spec.Replicas == nil {
		foo.Spec.Replicas = pointer.Int32(1)
	}
	if foo.Annotations == nil {
		foo.Annotations = map[string]string{}
	}
	foo.Annotations["example-mutating-admission-webhook"] = "foo"
	return nil
}

// +kubebuilder:webhook:path=/validate-samplecontroller-k8s-io-v1alpha1-foo,mutating=false,failurePolicy=fail,sideEffects=None,groups=samplecontroller.k8s.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.samplecontroller.k8s.io,admissionReviewVersions=v1

type FooValidator struct{}

func (v *FooValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return v.validate(ctx, obj)
}

func (v *FooValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	return v.validate(ctx, newObj)
}

func (v *FooValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return v.validate(ctx, obj)
}

func (v *FooValidator) validate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := logf.FromContext(ctx)
	log.Info("Validating Foo")

	foo, ok := obj.(*Foo)
	if !ok {
		return nil, fmt.Errorf("expected a Foo but got a %T", obj)
	}

	key := "example-mutating-admission-webhook"
	anno, found := foo.Annotations[key]
	if !found {
		return nil, fmt.Errorf("missing annotation %s", key)
	}
	if anno != "foo" {
		return nil, fmt.Errorf("annotation %s did not have value %q", key, "foo")
	}

	return nil, nil
}
