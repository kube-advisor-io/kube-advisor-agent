package dataproviders

import (
	"context"
	"fmt"

	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type GeneralInfoProvider struct {
	podsList *PodsList
	client *kubernetes.Clientset
}

func NewGeneralInfoProvider(client *kubernetes.Clientset, config config.Config) *GeneralInfoProvider {
	instance := new(GeneralInfoProvider)
	instance.podsList = GetPodsListInstance(client, config.IgnoredNamespaces)
	instance.client = client
	return instance
}

func (prov *GeneralInfoProvider) GetName() string {
	return "general_info_provider"
}

func (npp *GeneralInfoProvider) GetData() map[string]interface{} {
	namespacesList, err := npp.client.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println("error getting namespaces:", err)
	}
	nodeList, err := npp.client.CoreV1().Nodes().List(context.Background(), v1.ListOptions{})
	if err != nil {
		fmt.Println("error getting nodes:", err)
	}


	return map[string]interface{}{
		"podsCount": len(npp.podsList.Pods),
		"namespacesCount": len(namespacesList.Items),
		"nodesCount": len(nodeList.Items),
	}
}
