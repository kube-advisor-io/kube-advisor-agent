package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type DeploymentResourceProvider struct {
	resource        *schema.GroupVersionResource
	podResourceList *ResourcesList
}

func GetDeploymentResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *DeploymentResourceProvider {
	resource := &schema.GroupVersionResource{Group: "apps", Resource: "deployments", Version: "v1"}
	return &DeploymentResourceProvider{
		resource: resource,
		podResourceList: GetResourcesListInstance(
			dynamicClient,
			resource,
			ignoredNamespaces,
		),
	}
}

func (rp *DeploymentResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.resource
}

func (prov *DeploymentResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *DeploymentResourceProvider) GetParsedItems() []interface{} {
	var result []interface{}
	for _, deployment := range rp.podResourceList.Resources {
		var podParsed Deployment
		mapstructure.Decode(deployment, &podParsed)
		result = append(result, podParsed)
	}
	return result
}
