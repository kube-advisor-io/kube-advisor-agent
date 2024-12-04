package providers

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
	resourcesListLock                                                      = &sync.Mutex{}
	resourcesListInstances map[*schema.GroupVersionResource]*ResourcesList = map[*schema.GroupVersionResource]*ResourcesList{}
	listTimeout                                                            = int64(60)
)

// A generic list or K8s resources which is updated continously through watching that resource type on the K8s Api.
type ResourcesList struct {
	ignoredNamespaces     []string
	latestResourceVersion string
	dynamicClient         *dynamic.DynamicClient
	resource              *schema.GroupVersionResource
	Resources             []*map[string]interface{}
}

// Only one instance per resource type should exist.
// The per-resource singleton is returned here.
func GetResourcesListInstance(
	dynamicClient *dynamic.DynamicClient,
	resource *schema.GroupVersionResource,
	ignoredNamespaces []string,
) *ResourcesList {
	if resourcesListInstances[resource] == nil {
		resourcesListLock.Lock()
		defer resourcesListLock.Unlock()
		if resourcesListInstances[resource] == nil {
			resourcesListInstances[resource] = new(ResourcesList)
			resourcesListInstances[resource].dynamicClient = dynamicClient
			resourcesListInstances[resource].ignoredNamespaces = ignoredNamespaces
			resourcesListInstances[resource].resource = resource
			resourcesListInstances[resource].startWatching()
		} else {
			log.Trace("Single instance already created.")
		}
	} else {
		log.Trace("Single instance already created.")
	}

	return resourcesListInstances[resource]
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
		log.Trace("Starting watching resources...")
		return rl.dynamicClient.Resource(*rl.resource).Watch(context.Background(), metav1.ListOptions{
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
	resources, err := rl.dynamicClient.Resource(*rl.resource).List(context.Background(), metav1.ListOptions{TimeoutSeconds: &listTimeout})
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
	log.Trace("Found resource ", resource)
	namespace, name, _ := getNamespaceNameAndResourceVer(&resource)
	// kind := resource["Kind"].(string)
	if slices.Contains(rl.ignoredNamespaces, namespace) {
		log.Trace("Ignoring the resource", name, ", since it is in the ignored namespace "+namespace)
		return
	}
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
		if exisitingResourceName == newResourceName && existingResourceNamespace == newResourceNamespace {
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

func getNamespaceNameAndResourceVer(resource *map[string]interface{}) (string, string, string) {
	metadata := (*resource)["metadata"].(map[string]interface{})
	var namespace string
	if val, ok := metadata["namespace"]; ok {
		namespace = val.(string)
	}
	name := metadata["name"].(string)
	resourceVersion := metadata["resourceVersion"].(string)
	return namespace, name, resourceVersion
}
