package dataproviders

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/version"
	"time"
)

type ApiVersionProvider struct {
	client *kubernetes.Clientset
	apiVersion *version.Info
}

func NewApiVersionProvider(client *kubernetes.Clientset) *ApiVersionProvider{
	instance := new(ApiVersionProvider)
	instance.client = client
	go instance.startWatching()
	return instance
}

func (prov *ApiVersionProvider) GetData() map[string]interface{}{
	return map[string]interface{} {"apiVersion": prov.apiVersion.String()}
}

func (prov *ApiVersionProvider) startWatching(){	
	version,_ := prov.client.DiscoveryClient.ServerVersion()
	prov.apiVersion = version
	for range time.Tick(time.Second * 30) {
		version,_ := prov.client.DiscoveryClient.ServerVersion()
        prov.apiVersion = version
    }
}