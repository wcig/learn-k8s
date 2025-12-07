// 为包中任何类型生成深拷贝方法，可以在局部tag中覆盖此默认行为
// +k8s:deepcopy-gen=package
// 指定资源的group
// +groupName=samplecontroller.k8s.io

// Package v1alpha1 is the v1alpha1 version of the API.
package v1alpha1
