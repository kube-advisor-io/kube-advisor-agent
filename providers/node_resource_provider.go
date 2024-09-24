package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type NodeResourceProvider struct {
	ResourceProviderBase
}

func GetNodeResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *NodeResourceProvider {
	resource := &schema.GroupVersionResource{Group: "", Resource: "nodes", Version: "v1"}
	return &NodeResourceProvider{
		ResourceProviderBase: ResourceProviderBase{
			Resource: resource,
			ResourcesList: GetResourcesListInstance(
				dynamicClient,
				resource,
				ignoredNamespaces,
			)},
	}
}

func (rp *NodeResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (prov *NodeResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *NodeResourceProvider) GetParsedItems() []interface{} {
	var result []interface{} = []interface{}{}
	for _, node := range rp.ResourcesList.Resources {
		var nodeParsed Node
		mapstructure.Decode(node, &nodeParsed)
		result = append(result, nodeParsed)
	}
	return result
}
