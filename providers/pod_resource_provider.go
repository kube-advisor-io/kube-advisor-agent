package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type PodResourceProvider struct {
	ResourceProviderBase
}

func GetPodResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *PodResourceProvider {
	resource := &schema.GroupVersionResource{Group: "", Resource: "pods", Version: "v1"}
	return &PodResourceProvider{
		ResourceProviderBase: ResourceProviderBase{
			Resource: resource,
			ResourcesList: GetResourcesListInstance(
				dynamicClient,
				resource,
				ignoredNamespaces,
			)},
	}
}

func (rp *PodResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (prov *PodResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *PodResourceProvider) GetParsedItems() []interface{} {
	var result []interface{}
	for _, pod := range rp.ResourcesList.Resources {
		var podParsed Pod
		mapstructure.Decode(pod, &podParsed)
		result = append(result, podParsed)
	}
	return result
}
