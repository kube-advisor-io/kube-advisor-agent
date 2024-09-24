package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type StatefulsetResourceProvider struct {
	ResourceProviderBase
}

func GetStatefulsetResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *StatefulsetResourceProvider {
	resource := &schema.GroupVersionResource{Group: "apps", Resource: "statefulsets", Version: "v1"}
	return &StatefulsetResourceProvider{
		ResourceProviderBase: ResourceProviderBase{
			Resource: resource,
			ResourcesList: GetResourcesListInstance(
				dynamicClient,
				resource,
				ignoredNamespaces,
			)},
	}
}

func (rp *StatefulsetResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (prov *StatefulsetResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *StatefulsetResourceProvider) GetParsedItems() []interface{} {
	var result []interface{} = []interface{}{}
	for _, statefulset := range rp.ResourcesList.Resources {
		var statefulsetParsed Statefulset
		mapstructure.Decode(statefulset, &statefulsetParsed)
		result = append(result, statefulsetParsed)
	}
	return result
}
