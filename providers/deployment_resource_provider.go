package providers

import (
	"github.com/go-viper/mapstructure/v2"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type DeploymentResourceProvider struct {
	ResourceProviderBase
}

func GetDeploymentResourceProvider(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *DeploymentResourceProvider {
	resource := &schema.GroupVersionResource{Group: "apps", Resource: "deployments", Version: "v1"}
	return &DeploymentResourceProvider{
		ResourceProviderBase: ResourceProviderBase{
			Resource: resource,
			ResourcesList: GetResourcesListInstance(
				dynamicClient,
				resource,
				ignoredNamespaces,
			)},
	}
}

func (rp *DeploymentResourceProvider) GetResource() *schema.GroupVersionResource {
	return rp.Resource
}

func (prov *DeploymentResourceProvider) GetVersion() int32 {
	return 1
}

func (rp *DeploymentResourceProvider) GetParsedItems() []interface{} {
	var result []interface{} = []interface{}{}
	for _, deployment := range rp.ResourcesList.Resources {
		var deploymentParsed Deployment
		mapstructure.Decode(deployment, &deploymentParsed)
		result = append(result, deploymentParsed)
	}
	return result
}
