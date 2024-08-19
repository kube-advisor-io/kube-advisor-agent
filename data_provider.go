package main 

import (
	"k8s.io/client-go/kubernetes"
	"github.com/bobthebuilderberlin/kube-advisor-agent/dataproviders"
)

type DataProvider interface {
	GetData() map[string]interface{}
}

func getAllDataProviders(client *kubernetes.Clientset) []DataProvider {
   return []DataProvider{dataproviders.NewApiVersionProvider(client)}
}