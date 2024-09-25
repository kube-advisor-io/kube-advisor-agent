package main

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/providers"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ResourceProvider interface {
	GetResource() *schema.GroupVersionResource
	GetParsedItems() []interface{}
	GetVersion() int32
}

func getAllResourceProviders(dynamicClient *dynamic.DynamicClient, config config.Config) *[]ResourceProvider {
	return &[]ResourceProvider{
		providers.GetNamespaceResourceProvider2(dynamicClient, config.IgnoredNamespaces),
		providers.GetNodeResourceProvider2(dynamicClient, config.IgnoredNamespaces),
		providers.GetPodResourceProvider2(dynamicClient, config.IgnoredNamespaces),
		providers.GetDeploymentResourceProvider2(dynamicClient, config.IgnoredNamespaces),
		providers.GetStatefulsetResourceProvider2(dynamicClient, config.IgnoredNamespaces),
	}
}
