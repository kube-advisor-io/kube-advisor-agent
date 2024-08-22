package dataproviders

import (
	"k8s.io/client-go/kubernetes"
)

type NakedPodsProvider struct {
	podsList *PodsList
}

func NewNakedPodsProvider(client *kubernetes.Clientset) *NakedPodsProvider {
	instance := new(NakedPodsProvider)
	instance.podsList = GetPodsListInstance(client)
	return instance
}

func (npp *NakedPodsProvider) GetData() map[string]interface{} {
	nakedPods := []*PodInfo{}
	for _, podInfo := range npp.podsList.Pods {
		if podInfo.owner == "" {
			nakedPods = append(nakedPods, podInfo)
		}
	}
	return map[string]interface{}{
		"nakedPods": nakedPods,
	}
}