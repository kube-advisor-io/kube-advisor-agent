package dataproviders

import (
	"github.com/bobthebuilderberlin/kube-advisor-agent/config"
	"k8s.io/client-go/kubernetes"
)

type LabellessResourcesProvider struct {
	podsList        *PodsList
	deploymentsList *DeploymentsList
}

func NewLabellessResourcesProvider(client *kubernetes.Clientset, config config.Config) *LabellessResourcesProvider {
	instance := new(LabellessResourcesProvider)
	instance.podsList = GetPodsListInstance(client, config.IgnoredNamespaces)
	instance.deploymentsList = GetDeploymentsListInstance(client, config.IgnoredNamespaces)
	return instance
}

func (prov *LabellessResourcesProvider) GetName() string {
	return "labellessResourcesProvider"
}

func (prov *LabellessResourcesProvider) GetVersion() int32 {
	return 1
}

func (lrp *LabellessResourcesProvider) GetData() map[string]interface{} {
	labellessResources := []*Resource{}
	for _, pod := range lrp.podsList.Pods {
		if len(pod.Labels) == 0 {
			labellessResources = append(labellessResources, resourceFromPod(pod))
		}
	}
	for _, depyloyment := range lrp.deploymentsList.Deployments {
		if len(depyloyment.Labels) == 0 {
			labellessResources = append(labellessResources, resourceFromDeployment(depyloyment))
		}
	}
	return map[string]interface{}{
		"items": toResourceInfo(labellessResources),
	}
}
