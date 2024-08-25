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

func (prov *NakedPodsProvider) GetName() string {
	return "naked_pods_provider"
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
