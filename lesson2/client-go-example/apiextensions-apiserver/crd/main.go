package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/pointer"

	"client-go-example/utils"
)

// k8s.io/apiextensions-apiserver: crd扩展
func main() {
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// create crd clientSet
	clientSet, err := clientset.NewForConfig(config)
	utils.CheckErr(err)

	// create crd
	result, err := clientSet.ApiextensionsV1().CustomResourceDefinitions().Create(context.Background(), crd, metav1.CreateOptions{})
	utils.CheckErr(err)
	fmt.Println("result crd name:", result.GetName())

	// delete crd
	err = clientSet.ApiextensionsV1().CustomResourceDefinitions().Delete(context.Background(), crd.GetName(), metav1.DeleteOptions{})
	utils.CheckErr(err)

	// Output:
	// result crd name: examples.mygroup.mydomain
}

var crd = &v1.CustomResourceDefinition{
	ObjectMeta: metav1.ObjectMeta{
		Name: "examples.mygroup.mydomain",
	},
	Spec: v1.CustomResourceDefinitionSpec{
		Group: "mygroup.mydomain",
		Names: v1.CustomResourceDefinitionNames{
			Plural:   "examples",
			Singular: "example",
			Kind:     "Example",
		},
		Scope: v1.NamespaceScoped,
		Versions: []v1.CustomResourceDefinitionVersion{
			{
				Name:    "v1alpha1",
				Served:  true,
				Storage: true,
				Schema: &v1.CustomResourceValidation{
					OpenAPIV3Schema: &v1.JSONSchemaProps{
						Type: "object",
						Properties: map[string]v1.JSONSchemaProps{
							"spec": {
								Type: "object",
								Properties: map[string]v1.JSONSchemaProps{
									"deploymentName": {
										Type: "string",
									},
									"replicas": {
										Type:    "integer",
										Minimum: pointer.Float64(1),
										Maximum: pointer.Float64(10),
									},
								},
							},
							"status": {
								Type: "object",
								Properties: map[string]v1.JSONSchemaProps{
									"availableReplicas": {
										Type: "integer",
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
