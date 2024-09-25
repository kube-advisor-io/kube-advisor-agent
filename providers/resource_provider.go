package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ResourceProvider[T Resource] struct {
	Version       int32
	Resource      *schema.GroupVersionResource
	ResourcesList *ResourcesList
}

func GetPodResourceProvider2(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[Pod] {
	return getResourceProvider[Pod](
		&schema.GroupVersionResource{Group: "", Resource: "pods", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetDeploymentResourceProvider2(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[Deployment] {
	return getResourceProvider[Deployment](
		&schema.GroupVersionResource{Group: "apps", Resource: "deployments", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetStatefulsetResourceProvider2(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[Statefulset] {
	return getResourceProvider[Statefulset](
		&schema.GroupVersionResource{Group: "apps", Resource: "statefulsets", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetNodeResourceProvider2(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[Node] {
	return getResourceProvider[Node](
		&schema.GroupVersionResource{Group: "", Resource: "nodes", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}
func GetNamespaceResourceProvider2(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[Namespace] {
	return getResourceProvider[Namespace](
		&schema.GroupVersionResource{Group: "", Resource: "namespaces", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func getResourceProvider[T Resource](resource *schema.GroupVersionResource, dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string, version int32) *ResourceProvider[T] {
	return &ResourceProvider[T]{
		Version:  version,
		Resource: resource,
		ResourcesList: GetResourcesListInstance(
			dynamicClient,
			resource,
			ignoredNamespaces,
		)}

}

func (rp *ResourceProvider[T]) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (rp *ResourceProvider[T]) GetVersion() int32 {
	return rp.Version
}

func (rp *ResourceProvider[T]) GetParsedItems() []interface{} {
	var result []interface{} = []interface{}{}
	for _, resource := range rp.ResourcesList.Resources {
		var parsed T
		mapstructure.Decode(resource, &parsed)
		result = append(result, parsed)
	}
	return result
}
