package dataproviders

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

var (
	deploymentsListLock     = &sync.Mutex{}
	deploymentsListInstance *DeploymentsList
)

type DeploymentsList struct {
	ResourceList
	Deployments []*appsv1.Deployment
}

func GetDeploymentsListInstance(client *kubernetes.Clientset) *DeploymentsList {
	if deploymentsListInstance == nil {
		deploymentsListLock.Lock()
		defer deploymentsListLock.Unlock()
		if deploymentsListInstance == nil {
			deploymentsListInstance = new(DeploymentsList)
			deploymentsListInstance.client = client
			deploymentsListInstance.startWatching()
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return deploymentsListInstance
}

func (dl *DeploymentsList) startWatching() {
	err := dl.updateDeployments()
	if err != nil {
		log.Error(err)
		return
	}
	go dl.watchDeployments()
}

func (dl *DeploymentsList) watchDeployments() {
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		log.Info("Starting watching deployments...")
		return dl.client.AppsV1().Deployments("").Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds:  &timeOut,
			ResourceVersion: dl.latestResourceVersion,
		})
	}

	watcher, _ := toolsWatch.NewRetryWatcher(dl.latestResourceVersion, &cache.ListWatch{WatchFunc: watchFunc})

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			fmt.Printf("Error: Object: %v", event.Object)
		case watch.Modified:
			deployment := event.Object.(*appsv1.Deployment)
			dl.latestResourceVersion = deployment.ObjectMeta.ResourceVersion
			dl.updateDeployment(deployment)
		case watch.Deleted:
			deployment := event.Object.(*appsv1.Deployment)
			dl.latestResourceVersion = deployment.ObjectMeta.ResourceVersion
			dl.deleteDeployment(deployment.GetName(), deployment.GetNamespace())
		case watch.Added:
			deployment := event.Object.(*appsv1.Deployment)
			dl.latestResourceVersion = deployment.ObjectMeta.ResourceVersion
			dl.addDeployment(deployment)
		}
	}
}

func (dl *DeploymentsList) updateDeployments() error {
	deployments, err := dl.client.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
	if err != nil {
		return err
	}
	dl.Deployments = []*appsv1.Deployment{}
	for _, deployment := range deployments.Items {
		dl.addDeployment(&deployment)
	}
	dl.latestResourceVersion = deployments.ListMeta.ResourceVersion
	return nil
}

func (dl *DeploymentsList) addDeployment(deployment *appsv1.Deployment) {
	var owner string
	if len(deployment.OwnerReferences) != 0 {
		owner = deployment.OwnerReferences[0].Kind
	}
	log.Info("Found deployment ", deployment.Name, " in namespace ", deployment.Namespace, " with owner ", owner)
	dl.Deployments = append(dl.Deployments, deployment)
}

func (dl *DeploymentsList) deleteDeployment(name, namespace string) {
	log.Info("Deployment ", name, " in namespace ", namespace, " was deleted")
	for index, deployment := range dl.Deployments {
		if deployment.Name == name && deployment.Namespace == namespace {
			dl.Deployments = removeFromDeploymentsSlice(dl.Deployments, index)
		}
	}
}

func (pl *DeploymentsList) updateDeployment(deployment *appsv1.Deployment) {
	for index, exisitingDeployment := range pl.Deployments {
		if exisitingDeployment.Name == deployment.Name && exisitingDeployment.Namespace == deployment.Namespace {
			pl.Deployments[index] = deployment
			break
		}
	}
	log.Info("Updated deployment ", deployment.Name, " in namespace ", deployment.Namespace)
}

func removeFromDeploymentsSlice(s []*appsv1.Deployment, i int) []*appsv1.Deployment {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func resourceFromDeployment(pod *appsv1.Deployment) *Resource {
	resource := &Resource{TypeMeta: pod.TypeMeta, ObjectMeta: pod.ObjectMeta}
	resource.Kind = "Deployment"
	return resource
}
