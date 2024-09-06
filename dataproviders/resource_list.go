package dataproviders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type Resource struct {
	metav1.TypeMeta
	metav1.ObjectMeta
	corev1.PodSpec
}

type ResourceList struct {
	client                *kubernetes.Clientset
	ignoredNamespaces     []string
	latestResourceVersion string
}

type ResourceInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Node      string `json:"node"`
	owner     string
}

func toResourceInfo(resources []*Resource) []*ResourceInfo {
	podInfos := []*ResourceInfo{}
	for _, resource := range resources {
		podInfos = append(podInfos, &ResourceInfo{Name: resource.Name, Namespace: resource.Namespace, Kind: resource.Kind, Node: resource.NodeName})
	}
	return podInfos
}
