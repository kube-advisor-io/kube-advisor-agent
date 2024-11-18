package main

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/providers"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// emits a list of resources in addition to its version and resource type
type ResourceProvider interface {
	GetResource() *schema.GroupVersionResource
	GetParsedItems() []interface{}
	GetVersion() int32
}

// returns instances of all existing resource providers
func getAllResourceProviders(dynamicClient *dynamic.DynamicClient, config config.Config) *[]ResourceProvider {
	return &[]ResourceProvider{
		providers.GetNamespaceResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetNodeResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetPodResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetDeploymentResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetServiceResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetIngressResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetStatefulsetResourceProvider(dynamicClient, config.IgnoredNamespaces),
	}
}
