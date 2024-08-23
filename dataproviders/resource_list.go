package dataproviders

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
}

type ResourceList struct {
	client                *kubernetes.Clientset
	latestResourceVersion string
}

type ResourceInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	owner     string
}

func toResourceInfo(resources []*Resource) []*ResourceInfo {
	podInfos := []*ResourceInfo{}
	for _, resource := range resources {
		podInfos = append(podInfos, &ResourceInfo{Name: resource.Name, Namespace: resource.Namespace, Kind: resource.Kind})
	}
	return podInfos
}
