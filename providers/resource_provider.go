package providers

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/resources"
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ResourceProvider[T resources.Resource] struct {
	Version       int32
	Resource      *schema.GroupVersionResource
	ResourcesList *ResourcesList
}

func GetPodResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Pod] {
	return getResourceProvider[resources.Pod](
		&schema.GroupVersionResource{Group: "", Resource: "pods", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetDeploymentResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Deployment] {
	return getResourceProvider[resources.Deployment](
		&schema.GroupVersionResource{Group: "apps", Resource: "deployments", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetStatefulsetResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Statefulset] {
	return getResourceProvider[resources.Statefulset](
		&schema.GroupVersionResource{Group: "apps", Resource: "statefulsets", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetServiceResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Service] {
	return getResourceProvider[resources.Service](
		&schema.GroupVersionResource{Group: "", Resource: "services", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetIngressResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Ingress] {
	return getResourceProvider[resources.Ingress](
		&schema.GroupVersionResource{Group: "networking.k8s.io", Resource: "ingresses", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func GetNodeResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Node] {
	return getResourceProvider[resources.Node](
		&schema.GroupVersionResource{Group: "", Resource: "nodes", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}
func GetNamespaceResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourceProvider[resources.Namespace] {
	return getResourceProvider[resources.Namespace](
		&schema.GroupVersionResource{Group: "", Resource: "namespaces", Version: "v1"},
		dynamicClient,
		ignoredNamespaces,
		1,
	)
}

func getResourceProvider[T resources.Resource](resource *schema.GroupVersionResource, dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string, version int32) *ResourceProvider[T] {
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
