package dataproviders

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"k8s.io/client-go/kubernetes"
)

type NakedPodsProvider struct {
	podsList *PodsList
}

func NewNakedPodsProvider(client *kubernetes.Clientset, config config.Config) *NakedPodsProvider {
	instance := new(NakedPodsProvider)
	instance.podsList = GetPodsListInstance(client, config.IgnoredNamespaces)
	return instance
}

func (prov *NakedPodsProvider) GetName() string {
	return "nakedPodsProvider"
}

func (npp *NakedPodsProvider) GetData() map[string]interface{} {
	nakedPods := []*Resource{}
	for _, pod := range npp.podsList.Pods {
		if len(pod.OwnerReferences) == 0 {
			nakedPods = append(nakedPods, resourceFromPod(pod))
		}
	}
	return map[string]interface{}{
		"nakedPods": toResourceInfo(nakedPods),
	}
}
