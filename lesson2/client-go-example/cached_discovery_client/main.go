package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/tools/clientcmd"

	"client-go-example/utils"
)

// cachedDiscoveryClient: 本地缓存的发现客户端
func main() {
	// 加载config: 使用默认kubeConfigPath
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	utils.CheckErr(err)

	cachedDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(
		config,
		".cache/discovery",
		".cache/http",
		time.Hour)
	utils.CheckErr(err)

	apiGroupList, apiResourceList, err := cachedDiscoveryClient.ServerGroupsAndResources()
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
}
