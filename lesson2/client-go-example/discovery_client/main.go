package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"

	"client-go-example/utils"
)

// discoveryClient: 发现客户端
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	// 创建discoveryClient
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	utils.CheckErr(err)

	// 查询所有apiGroup/apiResource
	apiGroupList, apiResourceList, err := discoveryClient.ServerGroupsAndResources()
	utils.CheckErr(err)

	for _, d := range apiGroupList {
		fmt.Printf("[apiGroupList] name: %s, versions: %v\n", d.Name, d.Versions)
	}

	for _, d := range apiResourceList {
		gv, err := schema.ParseGroupVersion(d.GroupVersion)
		utils.CheckErr(err)
		for _, apiResource := range d.APIResources {
			fmt.Printf("[apiResourceList] name: %s, group: %s, version: %s, kind: %s\n",
				apiResource.Name, gv.Group, gv.Version, apiResource.Kind)
		}
	}

	// Output:
	// [apiGroupList] name: , versions: [{v1 v1}]
	// [apiGroupList] name: apiregistration.k8s.io, versions: [{apiregistration.k8s.io/v1 v1}]
	// [apiGroupList] name: apps, versions: [{apps/v1 v1}]
	// [apiGroupList] name: events.k8s.io, versions: [{events.k8s.io/v1 v1}]
	// [apiGroupList] name: authentication.k8s.io, versions: [{authentication.k8s.io/v1 v1}]
	// [apiGroupList] name: authorization.k8s.io, versions: [{authorization.k8s.io/v1 v1}]
	// [apiGroupList] name: autoscaling, versions: [{autoscaling/v2 v2} {autoscaling/v1 v1}]
	// [apiGroupList] name: batch, versions: [{batch/v1 v1}]
	// [apiGroupList] name: certificates.k8s.io, versions: [{certificates.k8s.io/v1 v1}]
	// [apiGroupList] name: networking.k8s.io, versions: [{networking.k8s.io/v1 v1}]
	// [apiGroupList] name: policy, versions: [{policy/v1 v1}]
	// [apiGroupList] name: rbac.authorization.k8s.io, versions: [{rbac.authorization.k8s.io/v1 v1}]
	// [apiGroupList] name: storage.k8s.io, versions: [{storage.k8s.io/v1 v1}]
	// [apiGroupList] name: admissionregistration.k8s.io, versions: [{admissionregistration.k8s.io/v1 v1}]
	// [apiGroupList] name: apiextensions.k8s.io, versions: [{apiextensions.k8s.io/v1 v1}]
	// [apiGroupList] name: scheduling.k8s.io, versions: [{scheduling.k8s.io/v1 v1}]
	// [apiGroupList] name: coordination.k8s.io, versions: [{coordination.k8s.io/v1 v1}]
	// [apiGroupList] name: node.k8s.io, versions: [{node.k8s.io/v1 v1}]
	// [apiGroupList] name: discovery.k8s.io, versions: [{discovery.k8s.io/v1 v1}]
	// [apiGroupList] name: flowcontrol.apiserver.k8s.io, versions: [{flowcontrol.apiserver.k8s.io/v1 v1} {flowcontrol.apiserver.k8s.io/v1beta3 v1beta3}]
	// [apiResourceList] name: horizontalpodautoscalers, group: autoscaling, version: v1, kind: HorizontalPodAutoscaler
	// [apiResourceList] name: horizontalpodautoscalers/status, group: autoscaling, version: v1, kind: HorizontalPodAutoscaler
	// [apiResourceList] name: horizontalpodautoscalers, group: autoscaling, version: v2, kind: HorizontalPodAutoscaler
	// [apiResourceList] name: horizontalpodautoscalers/status, group: autoscaling, version: v2, kind: HorizontalPodAutoscaler
	// [apiResourceList] name: certificatesigningrequests, group: certificates.k8s.io, version: v1, kind: CertificateSigningRequest
	// [apiResourceList] name: certificatesigningrequests/approval, group: certificates.k8s.io, version: v1, kind: CertificateSigningRequest
	// [apiResourceList] name: certificatesigningrequests/status, group: certificates.k8s.io, version: v1, kind: CertificateSigningRequest
	// [apiResourceList] name: ingressclasses, group: networking.k8s.io, version: v1, kind: IngressClass
	// [apiResourceList] name: ingresses, group: networking.k8s.io, version: v1, kind: Ingress
	// [apiResourceList] name: ingresses/status, group: networking.k8s.io, version: v1, kind: Ingress
	// [apiResourceList] name: networkpolicies, group: networking.k8s.io, version: v1, kind: NetworkPolicy
	// [apiResourceList] name: leases, group: coordination.k8s.io, version: v1, kind: Lease
	// [apiResourceList] name: flowschemas, group: flowcontrol.apiserver.k8s.io, version: v1beta3, kind: FlowSchema
	// [apiResourceList] name: flowschemas/status, group: flowcontrol.apiserver.k8s.io, version: v1beta3, kind: FlowSchema
	// [apiResourceList] name: prioritylevelconfigurations, group: flowcontrol.apiserver.k8s.io, version: v1beta3, kind: PriorityLevelConfiguration
	// [apiResourceList] name: prioritylevelconfigurations/status, group: flowcontrol.apiserver.k8s.io, version: v1beta3, kind: PriorityLevelConfiguration
	// [apiResourceList] name: bindings, group: , version: v1, kind: Binding
	// [apiResourceList] name: componentstatuses, group: , version: v1, kind: ComponentStatus
	// [apiResourceList] name: configmaps, group: , version: v1, kind: ConfigMap
	// [apiResourceList] name: endpoints, group: , version: v1, kind: Endpoints
	// [apiResourceList] name: events, group: , version: v1, kind: Event
	// [apiResourceList] name: limitranges, group: , version: v1, kind: LimitRange
	// [apiResourceList] name: namespaces, group: , version: v1, kind: Namespace
	// [apiResourceList] name: namespaces/finalize, group: , version: v1, kind: Namespace
	// [apiResourceList] name: namespaces/status, group: , version: v1, kind: Namespace
	// [apiResourceList] name: nodes, group: , version: v1, kind: Node
	// [apiResourceList] name: nodes/proxy, group: , version: v1, kind: NodeProxyOptions
	// [apiResourceList] name: nodes/status, group: , version: v1, kind: Node
	// [apiResourceList] name: persistentvolumeclaims, group: , version: v1, kind: PersistentVolumeClaim
	// [apiResourceList] name: persistentvolumeclaims/status, group: , version: v1, kind: PersistentVolumeClaim
	// [apiResourceList] name: persistentvolumes, group: , version: v1, kind: PersistentVolume
	// [apiResourceList] name: persistentvolumes/status, group: , version: v1, kind: PersistentVolume
	// [apiResourceList] name: pods, group: , version: v1, kind: Pod
	// [apiResourceList] name: pods/attach, group: , version: v1, kind: PodAttachOptions
	// [apiResourceList] name: pods/binding, group: , version: v1, kind: Binding
	// [apiResourceList] name: pods/ephemeralcontainers, group: , version: v1, kind: Pod
	// [apiResourceList] name: pods/eviction, group: , version: v1, kind: Eviction
	// [apiResourceList] name: pods/exec, group: , version: v1, kind: PodExecOptions
	// [apiResourceList] name: pods/log, group: , version: v1, kind: Pod
	// [apiResourceList] name: pods/portforward, group: , version: v1, kind: PodPortForwardOptions
	// [apiResourceList] name: pods/proxy, group: , version: v1, kind: PodProxyOptions
	// [apiResourceList] name: pods/status, group: , version: v1, kind: Pod
	// [apiResourceList] name: podtemplates, group: , version: v1, kind: PodTemplate
	// [apiResourceList] name: replicationcontrollers, group: , version: v1, kind: ReplicationController
	// [apiResourceList] name: replicationcontrollers/scale, group: , version: v1, kind: Scale
	// [apiResourceList] name: replicationcontrollers/status, group: , version: v1, kind: ReplicationController
	// [apiResourceList] name: resourcequotas, group: , version: v1, kind: ResourceQuota
	// [apiResourceList] name: resourcequotas/status, group: , version: v1, kind: ResourceQuota
	// [apiResourceList] name: secrets, group: , version: v1, kind: Secret
	// [apiResourceList] name: serviceaccounts, group: , version: v1, kind: ServiceAccount
	// [apiResourceList] name: serviceaccounts/token, group: , version: v1, kind: TokenRequest
	// [apiResourceList] name: services, group: , version: v1, kind: Service
	// [apiResourceList] name: services/proxy, group: , version: v1, kind: ServiceProxyOptions
	// [apiResourceList] name: services/status, group: , version: v1, kind: Service
	// [apiResourceList] name: csidrivers, group: storage.k8s.io, version: v1, kind: CSIDriver
	// [apiResourceList] name: csinodes, group: storage.k8s.io, version: v1, kind: CSINode
	// [apiResourceList] name: csistoragecapacities, group: storage.k8s.io, version: v1, kind: CSIStorageCapacity
	// [apiResourceList] name: storageclasses, group: storage.k8s.io, version: v1, kind: StorageClass
	// [apiResourceList] name: volumeattachments, group: storage.k8s.io, version: v1, kind: VolumeAttachment
	// [apiResourceList] name: volumeattachments/status, group: storage.k8s.io, version: v1, kind: VolumeAttachment
	// [apiResourceList] name: events, group: events.k8s.io, version: v1, kind: Event
	// [apiResourceList] name: mutatingwebhookconfigurations, group: admissionregistration.k8s.io, version: v1, kind: MutatingWebhookConfiguration
	// [apiResourceList] name: validatingadmissionpolicies, group: admissionregistration.k8s.io, version: v1, kind: ValidatingAdmissionPolicy
	// [apiResourceList] name: validatingadmissionpolicies/status, group: admissionregistration.k8s.io, version: v1, kind: ValidatingAdmissionPolicy
	// [apiResourceList] name: validatingadmissionpolicybindings, group: admissionregistration.k8s.io, version: v1, kind: ValidatingAdmissionPolicyBinding
	// [apiResourceList] name: validatingwebhookconfigurations, group: admissionregistration.k8s.io, version: v1, kind: ValidatingWebhookConfiguration
	// [apiResourceList] name: priorityclasses, group: scheduling.k8s.io, version: v1, kind: PriorityClass
	// [apiResourceList] name: controllerrevisions, group: apps, version: v1, kind: ControllerRevision
	// [apiResourceList] name: daemonsets, group: apps, version: v1, kind: DaemonSet
	// [apiResourceList] name: daemonsets/status, group: apps, version: v1, kind: DaemonSet
	// [apiResourceList] name: deployments, group: apps, version: v1, kind: Deployment
	// [apiResourceList] name: deployments/scale, group: apps, version: v1, kind: Scale
	// [apiResourceList] name: deployments/status, group: apps, version: v1, kind: Deployment
	// [apiResourceList] name: replicasets, group: apps, version: v1, kind: ReplicaSet
	// [apiResourceList] name: replicasets/scale, group: apps, version: v1, kind: Scale
	// [apiResourceList] name: replicasets/status, group: apps, version: v1, kind: ReplicaSet
	// [apiResourceList] name: statefulsets, group: apps, version: v1, kind: StatefulSet
	// [apiResourceList] name: statefulsets/scale, group: apps, version: v1, kind: Scale
	// [apiResourceList] name: statefulsets/status, group: apps, version: v1, kind: StatefulSet
	// [apiResourceList] name: selfsubjectreviews, group: authentication.k8s.io, version: v1, kind: SelfSubjectReview
	// [apiResourceList] name: tokenreviews, group: authentication.k8s.io, version: v1, kind: TokenReview
	// [apiResourceList] name: customresourcedefinitions, group: apiextensions.k8s.io, version: v1, kind: CustomResourceDefinition
	// [apiResourceList] name: customresourcedefinitions/status, group: apiextensions.k8s.io, version: v1, kind: CustomResourceDefinition
	// [apiResourceList] name: runtimeclasses, group: node.k8s.io, version: v1, kind: RuntimeClass
	// [apiResourceList] name: apiservices, group: apiregistration.k8s.io, version: v1, kind: APIService
	// [apiResourceList] name: apiservices/status, group: apiregistration.k8s.io, version: v1, kind: APIService
	// [apiResourceList] name: endpointslices, group: discovery.k8s.io, version: v1, kind: EndpointSlice
	// [apiResourceList] name: flowschemas, group: flowcontrol.apiserver.k8s.io, version: v1, kind: FlowSchema
	// [apiResourceList] name: flowschemas/status, group: flowcontrol.apiserver.k8s.io, version: v1, kind: FlowSchema
	// [apiResourceList] name: prioritylevelconfigurations, group: flowcontrol.apiserver.k8s.io, version: v1, kind: PriorityLevelConfiguration
	// [apiResourceList] name: prioritylevelconfigurations/status, group: flowcontrol.apiserver.k8s.io, version: v1, kind: PriorityLevelConfiguration
	// [apiResourceList] name: cronjobs, group: batch, version: v1, kind: CronJob
	// [apiResourceList] name: cronjobs/status, group: batch, version: v1, kind: CronJob
	// [apiResourceList] name: jobs, group: batch, version: v1, kind: Job
	// [apiResourceList] name: jobs/status, group: batch, version: v1, kind: Job
	// [apiResourceList] name: localsubjectaccessreviews, group: authorization.k8s.io, version: v1, kind: LocalSubjectAccessReview
	// [apiResourceList] name: selfsubjectaccessreviews, group: authorization.k8s.io, version: v1, kind: SelfSubjectAccessReview
	// [apiResourceList] name: selfsubjectrulesreviews, group: authorization.k8s.io, version: v1, kind: SelfSubjectRulesReview
	// [apiResourceList] name: subjectaccessreviews, group: authorization.k8s.io, version: v1, kind: SubjectAccessReview
	// [apiResourceList] name: clusterrolebindings, group: rbac.authorization.k8s.io, version: v1, kind: ClusterRoleBinding
	// [apiResourceList] name: clusterroles, group: rbac.authorization.k8s.io, version: v1, kind: ClusterRole
	// [apiResourceList] name: rolebindings, group: rbac.authorization.k8s.io, version: v1, kind: RoleBinding
	// [apiResourceList] name: roles, group: rbac.authorization.k8s.io, version: v1, kind: Role
	// [apiResourceList] name: poddisruptionbudgets, group: policy, version: v1, kind: PodDisruptionBudget
	// [apiResourceList] name: poddisruptionbudgets/status, group: policy, version: v1, kind: PodDisruptionBudget
}
