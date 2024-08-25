package dataproviders

import (
	"k8s.io/client-go/kubernetes"
)

type LabellessResourcesProvider struct {
	podsList        *PodsList
	deploymentsList *DeploymentsList
}

func NewLabellessResourcesProvider(client *kubernetes.Clientset) *LabellessResourcesProvider {
	instance := new(LabellessResourcesProvider)
	instance.podsList = GetPodsListInstance(client)
	instance.deploymentsList = GetDeploymentsListInstance(client)
	return instance
}

func (prov *LabellessResourcesProvider) GetName() string {
	return "labelless_resources_provider"
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
		"labellessResources": toResourceInfo(labellessResources),
	}
}
