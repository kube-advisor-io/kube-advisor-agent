package main

import (
	"slices"

	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/providers"
	"k8s.io/client-go/kubernetes"
)

// An interface for generic data providers that are not based on resource lists. Example: K8sVersionProvider.
type DataProvider interface {
	GetName() string
	GetData() map[string]interface{}
}

func getAllDataProviders(client *kubernetes.Clientset, config config.Config) *[]DataProvider {

	dataProviders := &[]DataProvider{
		providers.NewApiVersionProvider(client),
	}
	filteredDataProviders := []DataProvider{}
	for _, dataProvider := range *dataProviders {
		if !slices.Contains(config.DisabledProviders, dataProvider.GetName()){
			filteredDataProviders = append(filteredDataProviders, dataProvider)
		}
	}

	return &filteredDataProviders
}
