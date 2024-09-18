package dataproviders

import (
	"context"
	"slices"
	"sync"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

var (
	resourcesListLock     = &sync.Mutex{}
	resourcesListInstance *ResourcesList
)

type ResourcesList struct {
	ResourceList
	dynamicClient *dynamic.DynamicClient
	Resources []*map[string]interface{}
}

func GetResourcesListInstance(dynamicClient *dynamic.DynamicClient, ignoredNamespaces []string) *ResourcesList {
	if resourcesListInstance == nil {
		resourcesListLock.Lock()
		defer resourcesListLock.Unlock()
		if resourcesListInstance == nil {
			resourcesListInstance = new(ResourcesList)
			resourcesListInstance.dynamicClient = dynamicClient
			resourcesListInstance.ignoredNamespaces = ignoredNamespaces
			resourcesListInstance.startWatching()
		} else {
			log.Trace("Single instance already created.")
		}
	} else {
		log.Trace("Single instance already created.")
	}

	return resourcesListInstance
}

func (rl *ResourcesList) startWatching() {
	err := rl.updateResources()
	if err != nil {
		log.Error(err)
		return
	}
	go rl.watchResources()
}

func (rl *ResourcesList) watchResources() {
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		log.Info("Starting watching resources...")
		return rl.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Resource: "pods", Version: "v1"}).Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds:  &timeOut,
			ResourceVersion: rl.latestResourceVersion,
		})
	}

	watcher, _ := toolsWatch.NewRetryWatcher(rl.latestResourceVersion, &cache.ListWatch{WatchFunc: watchFunc})

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			log.Errorf("Error: Object: %v", event.Object)
		case watch.Modified:
			resource := event.Object.(runtime.Unstructured).UnstructuredContent()
			_, _, resourceVersion := getNamespaceNameAndResourceVer(&resource)
			rl.latestResourceVersion = resourceVersion
			rl.updateResource(&resource)
		case watch.Deleted:
			resource := event.Object.(runtime.Unstructured).UnstructuredContent()
			namespace, name, resourceVersion := getNamespaceNameAndResourceVer(&resource)
			rl.latestResourceVersion = resourceVersion
			rl.deleteResource(name, namespace)
		case watch.Added:
			resource := event.Object.(runtime.Unstructured).UnstructuredContent()
			_, _, resourceVersion := getNamespaceNameAndResourceVer(&resource)
			rl.latestResourceVersion = resourceVersion
			rl.addResource(resource)
		}
	}
}

func (rl *ResourcesList) updateResources() error {
	resources, err := rl.dynamicClient.Resource(schema.GroupVersionResource{Group: "", Resource: "pods", Version: "v1"}).List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
	if err != nil {
		return err
	}
	log.Info("Got resources ", resources)
	rl.Resources = *new([]*map[string]interface{})
	for _, resource := range resources.Items {
		rl.addResource(resource.Object)
	}
	rl.latestResourceVersion = resources.Object["metadata"].(map[string]interface{})["resourceVersion"].(string)
	return nil
}

func (rl *ResourcesList) addResource(resource map[string]interface{}) {
	log.Info("Found resource ", resource)
	namespace, name, _ := getNamespaceNameAndResourceVer(&resource)
	// kind := resource["Kind"].(string)
	if slices.Contains(rl.ignoredNamespaces, namespace){
		log.Trace("Ignoring the resource", name, ", since it is in the ignored namespace " + namespace)
		return
	}
	// var owner string
	// if len(resource["Namespace"].(map[string]interface{})) != 0 {
	// 	owner = resource.OwnerReferences[0].Kind
	// }
	rl.Resources = append(rl.Resources, &resource)
}

func (rl *ResourcesList) deleteResource(name, namespace string) {
	log.Trace("Resource ", name, " in namespace ", namespace, " was deleted")
	for index, resource := range rl.Resources {
		currentResourceNamespace, currentResourceName, _ := getNamespaceNameAndResourceVer(resource)
		if currentResourceName == name && currentResourceNamespace == namespace {
			rl.Resources = removeFromResourcesSlice(rl.Resources, index)
		}
	}
}

func (rl *ResourcesList) updateResource(newResource *map[string]interface{}) {
	newResourceNamespace, newResourceName, _ := getNamespaceNameAndResourceVer(newResource)
	for index, existingResource := range rl.Resources {
	    existingResourceNamespace, exisitingResourceName, _ := getNamespaceNameAndResourceVer(existingResource)
		if exisitingResourceName == newResourceName&& existingResourceNamespace == newResourceNamespace {
			rl.Resources[index] = newResource
			break
		}
	}
	log.Trace("Updated resource ", newResourceName, " in namespace ", newResourceNamespace)
}

func removeFromResourcesSlice(s []*map[string]interface{}, i int) []*map[string]interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func getNamespaceNameAndResourceVer(resource *map[string]interface{}) (string, string, string){
	metadata := (*resource)["metadata"].(map[string]interface{})
	namespace := metadata["namespace"].(string)
	name := metadata["name"].(string)
	resourceVersion := metadata["resourceVersion"].(string)
	return namespace, name, resourceVersion
}

// func resourceFromDeployment(pod *appsv1.Deployment) *Resource {
// 	resource := &Resource{TypeMeta: pod.TypeMeta, ObjectMeta: pod.ObjectMeta}
// 	resource.Kind = "Deployment"
// 	return resource
// }
