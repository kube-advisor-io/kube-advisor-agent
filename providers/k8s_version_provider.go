package providers

import (
	"time"

	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/kubernetes"
)

// Emits the used Kubernetes Version
type K8sVersionProvider struct {
	client     *kubernetes.Clientset
	k8sVersion *version.Info
}

func NewApiVersionProvider(client *kubernetes.Clientset) *K8sVersionProvider {
	instance := new(K8sVersionProvider)
	instance.client = client
	go instance.startWatching()
	return instance
}

func (prov *K8sVersionProvider) GetName() string {
	return "k8sVersionProvider"
}

func (prov *K8sVersionProvider) GetData() map[string]interface{} {
	return map[string]interface{}{"k8sVersion": prov.k8sVersion.String()}
}

func (prov *K8sVersionProvider) startWatching() {
	version, _ := prov.client.DiscoveryClient.ServerVersion()
	prov.k8sVersion = version
	for range time.Tick(time.Second * 120) {
		version, _ := prov.client.DiscoveryClient.ServerVersion()
		prov.k8sVersion = version
	}
}
