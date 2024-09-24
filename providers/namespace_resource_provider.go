package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type NamespaceResourceProvider struct {
	ResourceProviderBase
}

func GetNamespaceResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *NamespaceResourceProvider {
	resource := &schema.GroupVersionResource{Group: "", Resource: "namespaces", Version: "v1"}
	return &NamespaceResourceProvider{
		ResourceProviderBase: ResourceProviderBase{
			Resource: resource,
			ResourcesList: GetResourcesListInstance(
				dynamicClient,
				resource,
				ignoredNamespaces,
			)},
	}
}

func (rp *NamespaceResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (prov *NamespaceResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *NamespaceResourceProvider) GetParsedItems() []interface{} {
	var result []interface{} = []interface{}{}
	for _, namespace := range rp.ResourcesList.Resources {
		var namespaceParsed Namespace
		mapstructure.Decode(namespace, &namespaceParsed)
		result = append(result, namespaceParsed)
	}
	return result
}
