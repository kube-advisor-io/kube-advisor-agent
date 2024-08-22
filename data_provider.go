package main

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/dataproviders"
	"k8s.io/client-go/kubernetes"
)

type DataProvider interface {
	GetData() map[string]interface{}
}

func getAllDataProviders(client *kubernetes.Clientset) *[]DataProvider {
	return &[]DataProvider{
		dataproviders.NewApiVersionProvider(client),
		dataproviders.NewNakedPodsProvider(client),
		dataproviders.NewResourcelessPodsProvider(client),
	}
}
