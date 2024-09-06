package main

import (
	"slices"

	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"github.com/bobthebuilderberlin/kube-advisor-agent/dataproviders"
	"k8s.io/client-go/kubernetes"
)

type DataProvider interface {
	GetName() string
	GetData() map[string]interface{}
	GetVersion() int32
}

func getAllDataProviders(client *kubernetes.Clientset, config config.Config) *[]DataProvider {
	dataProviders := &[]DataProvider{
		dataproviders.NewApiVersionProvider(client),
		dataproviders.NewNakedPodsProvider(client, config),
		dataproviders.NewResourcelessPodsProvider(client, config),
		dataproviders.NewLabellessResourcesProvider(client, config),
		dataproviders.NewGeneralInfoProvider(client, config),
	}
	filteredDataProviders := []DataProvider{}
	for _, dataProvider := range *dataProviders {
		if !slices.Contains(config.DisabledProviders, dataProvider.GetName()) {
			filteredDataProviders = append(filteredDataProviders, dataProvider)
		}
	}

	return &filteredDataProviders
}
