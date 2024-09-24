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
		providers.GetPodResourceProvider(dynamicClient, config.IgnoredNamespaces),
		providers.GetDeploymentResourceProvider(dynamicClient, config.IgnoredNamespaces),
	}
}
