package dataproviders

import (
	corev1 "k8s.io/api/core/v1"
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
	nakedPods := []*corev1.Pod{}
	for _, pod := range npp.podsList.Pods {
		if len(pod.OwnerReferences) == 0 {
			nakedPods = append(nakedPods, pod)
		}
	}
	return map[string]interface{}{
		"nakedPods": toPodInfo(nakedPods),
	}
}
