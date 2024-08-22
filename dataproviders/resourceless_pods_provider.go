package dataproviders

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type ResourcelessPodsProvider struct {
	podsList *PodsList
}

func NewResourcelessPodsProvider(client *kubernetes.Clientset) *ResourcelessPodsProvider {
	instance := new(ResourcelessPodsProvider)
	instance.podsList = GetPodsListInstance(client)
	return instance
}

func (npp *ResourcelessPodsProvider) GetData() map[string]interface{} {
	resourcelessPods := []*corev1.Pod{}
	for _, pod := range npp.podsList.Pods {
		for _, container := range pod.Spec.Containers {
			if len(container.Resources.Limits) == 0 && len(container.Resources.Requests) == 0 {
				resourcelessPods = append(resourcelessPods, pod)
				break
			}
		}
	}
	return map[string]interface{}{
		"resourcelessPods": toPodInfo(resourcelessPods),
	}
}
