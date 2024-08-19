package dataproviders

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

type PodInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	owner     string
}

type NakedPodsProvider struct {
	client *kubernetes.Clientset
	pods   []*PodInfo
}

func NewNakedPodsProvider(client *kubernetes.Clientset) *NakedPodsProvider {
	instance := new(NakedPodsProvider)
	instance.client = client
	instance.pods = []*PodInfo{}
	go instance.startWatching()
	return instance
}

func (npp *NakedPodsProvider) startWatching() {
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		timeOut := int64(60)
		log.Info("Starting watching pods...")
		return npp.client.CoreV1().Pods("").Watch(context.Background(), metav1.ListOptions{TimeoutSeconds: &timeOut})
	}

	watcher, _ := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFunc})
	
	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			fmt.Printf("Error: Object: %v", event.Object)
		case watch.Deleted:
			pod := event.Object.(*corev1.Pod)
			npp.deletePod(pod.GetName(), pod.GetNamespace())
		case watch.Added:
			pod := event.Object.(*corev1.Pod)
			var owner string
			if len(pod.OwnerReferences) != 0 {
				owner = pod.OwnerReferences[0].Kind
			}
			npp.addPod(pod.GetName(), pod.GetNamespace(), owner)
		}
	}
}

func (npp *NakedPodsProvider) addPod(name, namespace, owner string) {
	log.Info("Found pod ", name, " in namespace ", namespace, " with owner ", owner)
	npp.pods = append(npp.pods, &PodInfo{Name: name, Namespace: namespace, owner: owner})
}

func (npp *NakedPodsProvider) deletePod(name, namespace string) {
	for index, podInfo := range npp.pods {
		if podInfo.Name == name && podInfo.Namespace == namespace {
			npp.pods = append(npp.pods[:index], npp.pods[index+1:]...)
		}
	}
}

func (npp *NakedPodsProvider) GetData() map[string]interface{} {
	nakedPods := []*PodInfo{}
	for _, podInfo := range npp.pods {
		if podInfo.owner == "" {
			nakedPods = append(nakedPods, podInfo)
		}
	}
	return map[string]interface{}{"nakedPods": nakedPods}
}
