package main

import (
	"slices"

	"github.com/bobthebuilderberlin/kube-advisor-agent/dataproviders"
	"k8s.io/client-go/kubernetes"
)

type DataProvider interface {
	GetName() string
	GetData() map[string]interface{}
}

func getAllDataProviders(client *kubernetes.Clientset, disabledDataProviders []string) *[]DataProvider {
	dataProviders := &[]DataProvider{
		dataproviders.NewApiVersionProvider(client),
		dataproviders.NewNakedPodsProvider(client),
		dataproviders.NewResourcelessPodsProvider(client),
		dataproviders.NewLabellessResourcesProvider(client),
	}
	filteredDataProviders := []DataProvider{}
	for _, dataProvider := range *dataProviders {
		if !slices.Contains(disabledDataProviders, dataProvider.GetName()){
			filteredDataProviders = append(filteredDataProviders, dataProvider)
		}
	}

	return &filteredDataProviders
}
